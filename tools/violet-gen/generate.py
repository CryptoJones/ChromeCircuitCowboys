#!/usr/bin/env python3
"""Violet Lotus content generator for Chrome Circuit Cowboys.

DEV-TIME ONLY. A BBS door game is an isolated process with no network access, so
this never runs inside the game. It reads the live content manifest dumped by
cowboy/loredump_test.go, asks Violet Lotus (the local uncensored creative model
on ronin28) to enrich every room description, zone-lore line, flavor-NPC line and
item lore line in the established cyberpunk-noir voice, and writes cowboy/lore.json
— which the game embeds and overlays at build time (cowboy/lore.go). Re-runnable:
already-generated keys are skipped unless --force.

Usage:
  python3 generate.py --manifest manifest.json --out ../../cowboy/lore.json
Env (with defaults matching the ronin28 OMI note):
  VIOLET_URL  default http://172.16.27.183:8080/v1/chat/completions
  VIOLET_KEY  default 0123456789abcdef
  VIOLET_MODEL default violet-lotus-12b
"""
import argparse, json, os, re, sys, time, urllib.request, urllib.error

URL = os.environ.get("VIOLET_URL", "http://172.16.27.183:8080/v1/chat/completions")
KEY = os.environ.get("VIOLET_KEY", "0123456789abcdef")
MODEL = os.environ.get("VIOLET_MODEL", "violet-lotus-12b")

SYSTEM = (
    "You are a writer for CHROME CIRCUIT COWBOYS, a cyberpunk BBS door game set in "
    "Noche City: rain-slick neon streets, corporate arcologies, gutter gangs, an "
    "underground descent (the Neon Wasteland down to the Geo-Anchor Vault) and an "
    "ascent through the Net (cyberspace, ICE, breaches). Voice: terse, gritty, "
    "noir, second-person-implied, street-slang ('choom', 'scrip', 'chrome', "
    "'sleeve', 'runner'). Sensory and concrete, never purple. Match the existing "
    "lines EXACTLY in tone and register. No markdown, no headers, no emoji, no "
    "quotation marks around the text. Plain prose only. Output STRICT JSON exactly "
    "as instructed — an object, nothing else, no prose before or after."
)

def chat(prompt, temperature=0.8, max_tokens=2048, retries=4):
    body = json.dumps({
        "model": MODEL,
        "messages": [{"role": "system", "content": SYSTEM},
                     {"role": "user", "content": prompt}],
        "temperature": temperature,
        "max_tokens": max_tokens,
    }).encode()
    last = None
    for attempt in range(retries):
        try:
            req = urllib.request.Request(URL, data=body, headers={
                "Content-Type": "application/json",
                "Authorization": "Bearer " + KEY,
            })
            with urllib.request.urlopen(req, timeout=300) as r:
                data = json.loads(r.read())
            return data["choices"][0]["message"]["content"]
        except (urllib.error.URLError, KeyError, json.JSONDecodeError, TimeoutError) as e:
            last = e
            time.sleep(2 * (attempt + 1))
    raise RuntimeError(f"Violet Lotus call failed after {retries} tries: {last}")

def chat_json(prompt, temperature=0.8, max_tokens=4096, tries=4):
    """Call Violet Lotus and parse a JSON object, retrying on bad/truncated JSON.
    Raises only after exhausting tries — callers skip the item and move on."""
    last = None
    for i in range(tries):
        try:
            return extract_json(chat(prompt, temperature=temperature,
                                     max_tokens=max_tokens))
        except (ValueError, json.JSONDecodeError) as e:
            last = e
            # nudge toward stricter output on retry
            prompt = prompt + "\n\nIMPORTANT: reply with ONLY the JSON object, complete and valid."
    raise last

def extract_json(text):
    """Pull the first balanced {...} object out of a model reply."""
    text = text.strip()
    # strip ```json fences if present
    text = re.sub(r"^```(?:json)?", "", text).strip()
    text = re.sub(r"```$", "", text).strip()
    start = text.find("{")
    if start < 0:
        raise ValueError("no JSON object in reply")
    depth, instr, esc = 0, False, False
    for i in range(start, len(text)):
        c = text[i]
        if instr:
            if esc: esc = False
            elif c == "\\": esc = True
            elif c == '"': instr = False
            continue
        if c == '"': instr = True
        elif c == "{": depth += 1
        elif c == "}":
            depth -= 1
            if depth == 0:
                return json.loads(text[start:i+1])
    raise ValueError("unbalanced JSON object")

def clean(s):
    return " ".join(str(s).replace("\r", " ").replace("\n", " ").split()).strip()

def chunked(seq, n):
    for i in range(0, len(seq), n):
        yield seq[i:i+n]

# ---- generators -----------------------------------------------------------

def gen_rooms(man, out, force):
    rooms = man["rooms"]
    by_zone = {}
    for r in rooms:
        by_zone.setdefault(r["zone_key"], []).append(r)
    done = out.setdefault("rooms", {})
    for zk, group in by_zone.items():
        zone_name = group[0]["zone_name"]
        lore = man["zone_lore"].get(zk, [])
        for batch in chunked(group, 8):
            todo = [r for r in batch if force or r["id"] not in done]
            if not todo:
                continue
            entries = "\n".join(
                f'- id "{r["id"]}" [{r["name"]}]: {clean(r["desc"])}' for r in todo)
            plot = ("\nZone plot context: " + " ".join(lore)) if lore else ""
            prompt = (
                f"AREA: {zone_name}.{plot}\n\n"
                "Rewrite each room description below to be more vivid and atmospheric "
                "(2-4 punchy sentences each), in the game's voice. PRESERVE every "
                "proper noun, character name, faction, mob, vendor/medic/cache hint, "
                "and any directional/exit cue from the original — only enrich the "
                "prose around them. Do not invent new named characters or new exits.\n\n"
                f"ROOMS:\n{entries}\n\n"
                'Return STRICT JSON: {"<id>": "<new description>", ...} with one key '
                "per room id above and nothing else."
            )
            try:
                res = chat_json(prompt)
            except Exception as e:
                print(f"  rooms {zk}: SKIP batch ({e})", flush=True)
                continue
            for r in todo:
                if r["id"] in res and clean(res[r["id"]]):
                    done[r["id"]] = clean(res[r["id"]])
            save(out)
            print(f"  rooms {zk}: +{len([r for r in todo if r['id'] in res])} "
                  f"({len(done)}/{len(rooms)})", flush=True)

def gen_zone_lore(man, out, force):
    done = out.setdefault("zone_lore", {})
    names = man["zone_names"]
    for zk, lines in man["zone_lore"].items():
        if zk in done and not force:
            continue
        zone_name = names.get(zk, zk)
        prompt = (
            f"AREA: {zone_name}. The locals here gossip to a runner who uses TALK. "
            "Below are the existing backstory lines for this area — they ARE the canon "
            "(characters, factions, plot beats). Keep all of them as-is in spirit, then "
            "EXPAND to a set of 5 distinct spoken lines total: rumors, warnings, lore, "
            "and color, each 1-2 sentences, in the gritty street voice. Each line is a "
            "standalone thing an NPC might say. Do not contradict the canon below.\n\n"
            "EXISTING LINES:\n" + "\n".join(f"- {l}" for l in lines) + "\n\n"
            'Return STRICT JSON: {"lines": ["...", "...", "...", "...", "..."]}'
        )
        try:
            res = chat_json(prompt, temperature=0.9)
        except Exception as e:
            print(f"  zone_lore {zk}: SKIP ({e})", flush=True)
            continue
        got = [clean(x) for x in res.get("lines", []) if clean(x)]
        if got:
            done[zk] = got
        save(out)
        print(f"  zone_lore {zk}: {len(got)} lines", flush=True)

def gen_npc(man, out, force):
    done = out.setdefault("room_npc", {})
    for rid, npc in man["room_npc"].items():
        if rid in done and not force:
            continue
        lines = npc["lines"]
        # Detect bilingual NPCs (Spanish greeters, Chinese taunts w/ " // " gloss).
        bilingual = any(" // " in l for l in lines)
        nonascii = any(ord(c) > 127 for l in lines for c in l)
        rule = ""
        if bilingual:
            rule = ("These lines are in another language with an English gloss after "
                    "' // '. Keep that EXACT format: 'foreign text // english gloss'. "
                    "Match the same language as the existing lines. ")
        elif nonascii:
            rule = ("Keep the same language as the existing lines. ")
        prompt = (
            f'NPC: {npc["speaker"]}. Below are their existing spoken lines — same '
            "persona, same attitude, same voice. EXPAND their repertoire to 6 distinct "
            f"lines total in character. {rule}Each line stands alone.\n\n"
            "EXISTING LINES:\n" + "\n".join(f"- {l}" for l in lines) + "\n\n"
            'Return STRICT JSON: {"lines": ["...", ...]}'
        )
        try:
            res = chat_json(prompt, temperature=0.9)
        except Exception as e:
            print(f"  room_npc {rid}: SKIP ({e})", flush=True)
            continue
        got = [clean(x) for x in res.get("lines", []) if clean(x)]
        if got:
            done[rid] = got
        save(out)
        print(f"  room_npc {rid}: {len(got)} lines", flush=True)

def gen_items(man, out, force):
    done = out.setdefault("items", {})
    items = man["items"]
    for batch in chunked(items, 10):
        todo = [it for it in batch if force or it["name"] not in done]
        if not todo:
            continue
        entries = "\n".join(
            f'- "{it["name"]}" ({it["effect"]}): {it["desc"]}' for it in todo)
        prompt = (
            "Write a single vivid lore sentence (street-tech flavor, where it came "
            "from / how it feels to use) for each item below. Keep it consistent with "
            "the item's real mechanical effect. One sentence each, no stats.\n\n"
            f"ITEMS:\n{entries}\n\n"
            'Return STRICT JSON: {"<name>": "<lore sentence>", ...}'
        )
        try:
            res = chat_json(prompt)
        except Exception as e:
            print(f"  items: SKIP batch ({e})", flush=True)
            continue
        for it in todo:
            if it["name"] in res and clean(res[it["name"]]):
                done[it["name"]] = clean(res[it["name"]])
        save(out)
        print(f"  items: {len(done)}/{len(items)}", flush=True)

OUT_PATH = None
def save(out):
    tmp = OUT_PATH + ".tmp"
    with open(tmp, "w") as f:
        json.dump(out, f, indent=2, ensure_ascii=False)
    os.replace(tmp, OUT_PATH)

def main():
    global OUT_PATH
    ap = argparse.ArgumentParser()
    ap.add_argument("--manifest", required=True)
    ap.add_argument("--out", required=True)
    ap.add_argument("--force", action="store_true")
    ap.add_argument("--only", choices=["rooms", "zone_lore", "npc", "items"])
    args = ap.parse_args()
    OUT_PATH = args.out
    with open(args.manifest) as f:
        man = json.load(f)
    out = {}
    if os.path.exists(args.out):
        with open(args.out) as f:
            out = json.load(f)
    steps = [("zone_lore", gen_zone_lore), ("npc", gen_npc),
             ("items", gen_items), ("rooms", gen_rooms)]
    for name, fn in steps:
        if args.only and args.only != name:
            continue
        print(f"== {name} ==", flush=True)
        fn(man, out, args.force)
    save(out)
    print("done ->", args.out, flush=True)

if __name__ == "__main__":
    main()
