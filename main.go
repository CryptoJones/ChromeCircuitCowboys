// Command cowboy is the persistent game server for Chrome Circuit Cowboys — a
// multiplayer cyberpunk MUD. It listens on TCP; AdmiralBBS bridges each caller
// in as a "resident" door, so everyone shares one live world. All game state is
// mutated on a single goroutine (events + ticks serialized), so the engine
// stays lock-free and deterministic; this process only owns I/O.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/CryptoJones/ChromeCircuitCowboys/cowboy"
)

// version is this build's release version (compared against the configured
// update feed, if any).
const version = "1.0.5"

func main() {
	addr := flag.String("addr", "127.0.0.1:4000", "TCP listen address for BBS bridge")
	dbPath := flag.String("db", "cowboy.db", "character database path (SQLite)")
	tick := flag.Duration("tick", 2*time.Second, "combat/world tick interval")
	updateURL := flag.String("update-url", os.Getenv("CCC_UPDATE_URL"),
		"forge 'releases/latest' JSON endpoint to check for updates (GitHub/Codeberg/Forgejo shape); empty = no check")
	showVer := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVer {
		println("Chrome Circuit Cowboys " + version)
		return
	}
	checkUpdate(version, *updateURL)

	store, err := cowboy.OpenSQLite(*dbPath)
	if err != nil {
		log.Fatalf("open character db: %v", err)
	}
	defer store.Close()

	world := cowboy.NewWorld(store)
	events := make(chan event, 256)

	// The single world goroutine: every mutation happens here.
	go func() {
		ticker := time.NewTicker(*tick)
		autosave := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		defer autosave.Stop()
		for {
			select {
			case ev := <-events:
				if ev.typ == evShutdown {
					world.SaveAll() // flush everyone before the process exits
					close(ev.done)
					return
				}
				handle(world, ev)
			case <-ticker.C:
				world.Tick()
				// Re-show the prompt ONLY for players who got output this tick
				// (combat/chat/room events). Idle players keep their single
				// prompt instead of it repeating every tick.
				for _, c := range activeConns() {
					if c.player != nil {
						world.PromptIfDirty(c.player)
					}
				}
			case <-autosave.C:
				world.SaveAll() // periodic persistence so a crash loses < 30s
			}
		}
	}()

	// Graceful shutdown: on SIGINT/SIGTERM (e.g. systemctl restart), flush all
	// connected players via the world goroutine before exiting.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Print("shutting down — saving players")
		done := make(chan struct{})
		events <- event{typ: evShutdown, done: done}
		select {
		case <-done:
		case <-time.After(5 * time.Second):
		}
		os.Exit(0)
	}()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Chrome Circuit Cowboys %s listening on %s (tick %s)", version, *addr, *tick)
	for {
		nc, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go serve(nc, events)
	}
}

// checkUpdate asynchronously compares the running version to the latest release
// at updateURL — a forge "releases/latest" JSON endpoint ({"tag_name":"vX.Y.Z"},
// the shape GitHub, Codeberg, and Forgejo all share). The forge is configured,
// never hardcoded: empty URL or any error is a silent no-op and never blocks
// startup.
func checkUpdate(current, updateURL string) {
	if strings.TrimSpace(updateURL) == "" {
		return
	}
	go func() {
		client := &http.Client{Timeout: 6 * time.Second}
		resp, err := client.Get(updateURL)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var rel struct {
			TagName string `json:"tag_name"`
		}
		if json.NewDecoder(resp.Body).Decode(&rel) != nil {
			return
		}
		latest := strings.TrimPrefix(strings.TrimSpace(rel.TagName), "v")
		if latest != "" && latest != current {
			log.Printf("*** UPDATE AVAILABLE: Chrome Circuit Cowboys %s is out (running %s) — %s ***",
				latest, current, updateURL)
		}
	}()
}

// ---- per-connection plumbing ----

type conn struct {
	nc     net.Conn
	outCh  chan string
	player *cowboy.Player
	closed bool // set on the world goroutine during teardown

	mu     sync.Mutex // guards inLine/prompt so output redraws don't race input echo
	inLine []byte     // the caller's in-progress (un-submitted) input
	prompt string     // the current status prompt (managed-prompt redraw)
}

// raw enqueues bytes for the writer. Non-blocking (drop on overflow) so a stalled
// client can't block the world goroutine.
func (c *conn) raw(s string) {
	select {
	case c.outCh <- s:
	default:
	}
}

// out is the simple pre-world sink (login/creation prompts; no input to preserve).
func (c *conn) out(s string) { c.raw(s) }

const clrLine = "\r\x1b[K" // CR + erase-to-end-of-line: wipe the current row

// emit writes async/response content while PRESERVING the caller's in-progress
// input: wipe the current prompt+input row, print the content, then redraw the
// prompt with whatever they'd typed. This is what stops combat lines from
// scrambling text the player is mid-typing.
func (c *conn) emit(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.raw(clrLine + s + c.prompt + string(c.inLine))
}

// setPrompt updates the status prompt and redraws it (with current input).
func (c *conn) setPrompt(p string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prompt = p
	c.raw(clrLine + c.prompt + string(c.inLine))
}

var (
	connMu  sync.Mutex
	connSet = map[*conn]struct{}{}
)

func activeConns() []*conn {
	connMu.Lock()
	defer connMu.Unlock()
	out := make([]*conn, 0, len(connSet))
	for c := range connSet {
		out = append(out, c)
	}
	return out
}

// readHandleHint reads the OPTIONAL reciprocal handshake the host sends right
// after our caps=handle advertisement: ESC ] ABBS;handle=<h> BEL. Returns the
// handle, or "" if nothing arrives in a short window (so a direct nc/telnet
// caller, or a host that pushes no handle, just sees the normal prompt). It
// never consumes bytes that aren't the sentinel.
func readHandleHint(nc net.Conn, r *bufio.Reader) string {
	const pfx = "\x1b]ABBS;handle="
	_ = nc.SetReadDeadline(time.Now().Add(700 * time.Millisecond))
	defer nc.SetReadDeadline(time.Time{})
	peek, err := r.Peek(len(pfx))
	if err != nil || string(peek) != pfx {
		return "" // no sentinel — leave whatever arrived for the prompt
	}
	_, _ = r.Discard(len(pfx))
	var h []byte
	for i := 0; i < 64; i++ {
		b, e := r.ReadByte()
		if e != nil || b == 0x07 {
			break
		}
		h = append(h, b)
	}
	return strings.TrimSpace(string(h))
}

// readArrow consumes an escape sequence after ESC and returns the movement
// direction for an arrow key (Up/Down/Right/Left -> north/south/east/west), or
// "" for any other sequence (which it consumes and ignores). Handles both CSI
// (ESC [ A) and SS3 (ESC O A) forms.
func readArrow(r *bufio.Reader) string {
	b, err := r.ReadByte()
	if err != nil {
		return ""
	}
	if b != '[' && b != 'O' { // lone ESC / Alt-combo — push the byte back as input
		_ = r.UnreadByte()
		return ""
	}
	n := 0
	var final byte
	for i := 0; i < 8; i++ {
		fb, e := r.ReadByte()
		if e != nil {
			return ""
		}
		n++
		if fb >= 0x40 && fb <= 0x7e { // final byte of the sequence
			final = fb
			break
		}
	}
	if n != 1 { // a plain arrow is a single final byte; modified/other keys aren't
		return ""
	}
	switch final {
	case 'A':
		return "north"
	case 'B':
		return "south"
	case 'C':
		return "east"
	case 'D':
		return "west"
	}
	return ""
}

func serve(nc net.Conn, events chan event) {
	c := &conn{nc: nc, outCh: make(chan string, 512)}
	connMu.Lock()
	connSet[c] = struct{}{}
	connMu.Unlock()

	// Writer goroutine: drains outCh to the socket, then closes the socket.
	go func() {
		for s := range c.outCh {
			if _, err := nc.Write([]byte(s)); err != nil {
				break
			}
		}
		nc.Close()
	}()

	r := bufio.NewReader(nc)
	// Advertise our version + the "handle" capability as the very FIRST bytes
	// (ABBS Door Spec §2.2). OSC-framed, so a terminal reached directly just
	// swallows it. Because we asked for "handle", the host pushes the caller's
	// BBS handle back; we read it and default the name prompt to it.
	c.out("\x1b]ABBS;version=" + version + ";caps=handle\x07")
	deflt := readHandleHint(nc, r)
	prompt := "Handle (your runner name): "
	if deflt != "" {
		prompt = "Handle [" + deflt + "] (Enter to use): "
	}
	c.out("\r\n" + prompt)
	name, err := cowboy.ReadLine(r, c.out)
	if err != nil {
		events <- event{typ: evClose, c: c}
		return
	}
	if t := strings.ToLower(strings.TrimSpace(name)); t == "q" || t == "quit" {
		c.out("\r\nNO CARRIER\r\n")
		events <- event{typ: evClose, c: c}
		return
	}
	if strings.TrimSpace(name) == "" {
		name = deflt // hit Enter on the prompt -> use the BBS handle
	}
	if strings.TrimSpace(name) == "" {
		events <- event{typ: evClose, c: c}
		return
	}

	reply := make(chan connectResult, 1)
	events <- event{typ: evConnect, c: c, name: name, reply: reply}
	res := <-reply
	if res.rejected {
		// name already online — give the writer a beat to flush, then close.
		time.Sleep(200 * time.Millisecond)
		events <- event{typ: evClose, c: c}
		return
	}
	if res.needCreate {
		// New runner: run the creation screen on this goroutine (the I/O side),
		// then hand the chosen loadout to the world to build the character.
		spec, err := cowboy.RunCharacterCreation(r, c.out)
		if err != nil {
			events <- event{typ: evClose, c: c}
			return
		}
		reply2 := make(chan connectResult, 1)
		events <- event{typ: evCreate, c: c, name: name, spec: spec, reply: reply2}
		<-reply2
	}

	// In-world input loop with a MANAGED prompt: we own the input buffer here so
	// async output (combat/chat) can wipe-and-redraw it via conn.emit/setPrompt
	// without garbling what the caller is typing.
	for {
		b, err := r.ReadByte()
		if err != nil {
			events <- event{typ: evDisconnect, c: c}
			return
		}
		switch b {
		case '\r', '\n':
			if r.Buffered() > 0 { // swallow a buffered CRLF/LFCR partner (never block)
				if nb, e := r.ReadByte(); e == nil {
					if !((b == '\r' && nb == '\n') || (b == '\n' && nb == '\r')) {
						_ = r.UnreadByte()
					}
				}
			}
			c.mu.Lock()
			line := string(c.inLine)
			c.inLine = c.inLine[:0]
			c.raw("\r\n")
			c.mu.Unlock()
			events <- event{typ: evLine, c: c, line: line}
		case 0x08, 0x7f: // backspace / DEL
			c.mu.Lock()
			if len(c.inLine) > 0 {
				c.inLine = c.inLine[:len(c.inLine)-1]
				c.raw("\b \b")
			}
			c.mu.Unlock()
		case 0x00:
			// ignore
		case 0xff: // telnet IAC — skip it and its command byte
			_, _ = r.ReadByte()
		case 0x1b: // ESC — an arrow key becomes a movement command
			if dir := readArrow(r); dir != "" {
				c.mu.Lock()
				moved := len(c.inLine) == 0 // only when not mid-typing a command
				if moved {
					c.raw("\r\n")
				}
				c.mu.Unlock()
				if moved {
					events <- event{typ: evLine, c: c, line: dir}
				}
			}
		default:
			if b >= 0x20 && b < 0x7f {
				c.mu.Lock()
				c.inLine = append(c.inLine, b)
				c.raw(string(b)) // echo
				c.mu.Unlock()
			}
		}
	}
}

// ---- events (all handled on the world goroutine) ----

type evType int

const (
	evConnect evType = iota
	evCreate
	evLine
	evDisconnect
	evClose
	evShutdown // flush all players, then signal done (graceful shutdown)
)

// connectResult tells the connection goroutine how the world handled a connect:
// rejected (name online), needCreate (new runner — run the creation screen), or
// otherwise an existing character was placed in the world.
type connectResult struct {
	rejected   bool
	needCreate bool
}

type event struct {
	typ   evType
	c     *conn
	name  string
	line  string
	spec  cowboy.CharSpec
	reply chan connectResult
	done  chan struct{} // evShutdown: closed once all players are saved
}

func handle(world *cowboy.World, ev event) {
	switch ev.typ {
	case evConnect:
		if world.Online(ev.name) {
			ev.c.out("\r\nThat runner is already jacked in. Try another handle.\r\n")
			ev.reply <- connectResult{rejected: true}
			return
		}
		if !world.HasCharacter(ev.name) {
			ev.reply <- connectResult{needCreate: true}
			return
		}
		p := world.Connect(ev.name, ev.c.emit)
		world.SetPrompter(p, ev.c.setPrompt)
		ev.c.player = p
		world.Prompt(p)
		ev.reply <- connectResult{}
	case evCreate:
		if world.Online(ev.name) { // lost a race to another connection
			ev.c.out("\r\nThat runner just jacked in elsewhere.\r\n")
			ev.reply <- connectResult{rejected: true}
			return
		}
		p := world.CreateCharacter(ev.name, ev.spec, ev.c.emit)
		world.SetPrompter(p, ev.c.setPrompt)
		ev.c.player = p
		world.Prompt(p)
		ev.reply <- connectResult{}
	case evLine:
		if ev.c.closed || ev.c.player == nil {
			return
		}
		if world.Command(ev.c.player, ev.line) {
			teardown(world, ev.c)
		}
	case evDisconnect, evClose:
		teardown(world, ev.c)
	}
}

func teardown(world *cowboy.World, c *conn) {
	if c.closed {
		return
	}
	c.closed = true
	if c.player != nil {
		world.Disconnect(c.player)
		c.player = nil
	}
	connMu.Lock()
	delete(connSet, c)
	connMu.Unlock()
	close(c.outCh) // ends the writer, which flushes remaining output then closes the socket
}
