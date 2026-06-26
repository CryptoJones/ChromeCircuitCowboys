package cowboy

import (
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

// SQLiteStore persists characters in a pure-Go SQLite file (separate from the
// BBS database — the game owns its own state).
type SQLiteStore struct{ db *sql.DB }

// OpenSQLite opens (or creates) the character database at path.
func OpenSQLite(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS cowboy_player (
		name          TEXT PRIMARY KEY COLLATE NOCASE,
		class         TEXT NOT NULL DEFAULT '',
		level         INTEGER NOT NULL,
		xp            INTEGER NOT NULL,
		eddies        INTEGER NOT NULL,
		hp            INTEGER NOT NULL,
		maxhp         INTEGER NOT NULL,
		body          INTEGER NOT NULL,
		reflexes      INTEGER NOT NULL,
		intelligence  INTEGER NOT NULL,
		weapon_bonus  INTEGER NOT NULL,
		weapon_name   TEXT NOT NULL,
		ram           INTEGER NOT NULL DEFAULT 0,
		deck_bonus    INTEGER NOT NULL DEFAULT 0,
		room          TEXT NOT NULL,
		inv_json      TEXT NOT NULL,
		quests_json   TEXT NOT NULL DEFAULT '{}',
		stash_json    TEXT NOT NULL DEFAULT '{}',
		stat_points   INTEGER NOT NULL DEFAULT 0,
		done_json     TEXT NOT NULL DEFAULT '{}',
		clan          TEXT NOT NULL DEFAULT ''
	)`); err != nil {
		db.Close()
		return nil, err
	}
	// Idempotent migrations for DBs created before a column existed (each errors
	// harmlessly when the column is already present).
	_, _ = db.Exec(`ALTER TABLE cowboy_player ADD COLUMN stash_json TEXT NOT NULL DEFAULT '{}'`)
	_, _ = db.Exec(`ALTER TABLE cowboy_player ADD COLUMN stat_points INTEGER NOT NULL DEFAULT 0`)
	_, _ = db.Exec(`ALTER TABLE cowboy_player ADD COLUMN done_json TEXT NOT NULL DEFAULT '{}'`)
	_, _ = db.Exec(`ALTER TABLE cowboy_player ADD COLUMN clan TEXT NOT NULL DEFAULT ''`)
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS cowboy_mail (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		to_name   TEXT NOT NULL COLLATE NOCASE,
		from_name TEXT NOT NULL,
		body      TEXT NOT NULL
	)`); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

// PushMail queues a message for a recipient.
func (s *SQLiteStore) PushMail(to, from, body string) error {
	_, err := s.db.Exec(`INSERT INTO cowboy_mail (to_name, from_name, body) VALUES (?,?,?)`, to, from, body)
	return err
}

// PopMail returns and deletes the recipient's queued mail (oldest first).
func (s *SQLiteStore) PopMail(to string) ([]Mail, error) {
	rows, err := s.db.Query(`SELECT from_name, body FROM cowboy_mail WHERE to_name = ? COLLATE NOCASE ORDER BY id`, to)
	if err != nil {
		return nil, err
	}
	var out []Mail
	for rows.Next() {
		var m Mail
		if err := rows.Scan(&m.From, &m.Body); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, m)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	_, err = s.db.Exec(`DELETE FROM cowboy_mail WHERE to_name = ? COLLATE NOCASE`, to)
	return out, err
}

// Close releases the database.
func (s *SQLiteStore) Close() error { return s.db.Close() }

// Load fetches a saved character by name.
func (s *SQLiteStore) Load(name string) (*SavedPlayer, bool, error) {
	var sp SavedPlayer
	var invJSON, questsJSON, stashJSON, doneJSON string
	err := s.db.QueryRow(`SELECT name, class, level, xp, eddies, hp, maxhp, body, reflexes,
		intelligence, weapon_bonus, weapon_name, ram, deck_bonus, room, inv_json, quests_json, stash_json, stat_points, done_json, clan
		FROM cowboy_player WHERE name = ? COLLATE NOCASE`, name).
		Scan(&sp.Name, &sp.Class, &sp.Level, &sp.XP, &sp.Eddies, &sp.HP, &sp.MaxHP, &sp.Body,
			&sp.Reflexes, &sp.Intelligence, &sp.WeaponBonus, &sp.WeaponName, &sp.RAM, &sp.DeckBonus, &sp.Room, &invJSON, &questsJSON, &stashJSON, &sp.StatPoints, &doneJSON, &sp.Clan)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	sp.Inv = map[string]int{}
	_ = json.Unmarshal([]byte(invJSON), &sp.Inv)
	sp.Stash = map[string]int{}
	_ = json.Unmarshal([]byte(stashJSON), &sp.Stash)
	sp.Quests = map[string]int{}
	_ = json.Unmarshal([]byte(questsJSON), &sp.Quests)
	sp.Done = map[string]int{}
	_ = json.Unmarshal([]byte(doneJSON), &sp.Done)
	return &sp, true, nil
}

// Top returns up to n characters ranked by level then XP (for the leaderboard).
// Only the scalar fields are read; Inv/Quests are left nil.
func (s *SQLiteStore) Top(n int) ([]SavedPlayer, error) {
	rows, err := s.db.Query(`SELECT name, class, level, xp, eddies
		FROM cowboy_player ORDER BY level DESC, xp DESC, name LIMIT ?`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SavedPlayer
	for rows.Next() {
		var sp SavedPlayer
		if err := rows.Scan(&sp.Name, &sp.Class, &sp.Level, &sp.XP, &sp.Eddies); err != nil {
			return nil, err
		}
		out = append(out, sp)
	}
	return out, rows.Err()
}

// Save upserts a character.
func (s *SQLiteStore) Save(sp *SavedPlayer) error {
	inv, _ := json.Marshal(sp.Inv)
	qjson, _ := json.Marshal(sp.Quests)
	stash, _ := json.Marshal(sp.Stash)
	done, _ := json.Marshal(sp.Done)
	_, err := s.db.Exec(`INSERT INTO cowboy_player
		(name, class, level, xp, eddies, hp, maxhp, body, reflexes, intelligence, weapon_bonus, weapon_name, ram, deck_bonus, room, inv_json, quests_json, stash_json, stat_points, done_json, clan)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(name) DO UPDATE SET
		  class=excluded.class, level=excluded.level, xp=excluded.xp, eddies=excluded.eddies, hp=excluded.hp,
		  maxhp=excluded.maxhp, body=excluded.body, reflexes=excluded.reflexes,
		  intelligence=excluded.intelligence, weapon_bonus=excluded.weapon_bonus,
		  weapon_name=excluded.weapon_name, ram=excluded.ram, deck_bonus=excluded.deck_bonus,
		  room=excluded.room, inv_json=excluded.inv_json, quests_json=excluded.quests_json, stash_json=excluded.stash_json,
		  stat_points=excluded.stat_points, done_json=excluded.done_json, clan=excluded.clan`,
		sp.Name, sp.Class, sp.Level, sp.XP, sp.Eddies, sp.HP, sp.MaxHP, sp.Body, sp.Reflexes,
		sp.Intelligence, sp.WeaponBonus, sp.WeaponName, sp.RAM, sp.DeckBonus, sp.Room, string(inv), string(qjson), string(stash), sp.StatPoints, string(done), sp.Clan)
	return err
}
