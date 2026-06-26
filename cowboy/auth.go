package cowboy

import "golang.org/x/crypto/bcrypt"

// Character password auth (#55/#56). Passwords are stored as bcrypt hashes,
// never plaintext. The server does the actual prompting on the I/O side; these
// are the world-side store helpers it drives via events.

// CharAuth describes a name's auth state for the login flow.
type CharAuth struct {
	Exists      bool
	HasPassword bool
}

// AuthInfo returns whether a character exists and whether it has a password set
// (legacy characters exist with no password → trigger the set-password
// migration, #56).
func (w *World) AuthInfo(name string) CharAuth {
	sp, ok, _ := w.store.Load(name)
	if !ok {
		return CharAuth{}
	}
	return CharAuth{Exists: true, HasPassword: sp.PasswordHash != ""}
}

// CheckPassword verifies a plaintext password against a saved character's hash.
func (w *World) CheckPassword(name, plain string) bool {
	sp, ok, _ := w.store.Load(name)
	if !ok || sp.PasswordHash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(sp.PasswordHash), []byte(plain)) == nil
}

// SetPassword hashes and stores a password for a character (create or migrate).
// If the character is currently online its live hash is updated too.
func (w *World) SetPassword(name, plain string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if p := w.onlineByName(name); p != nil {
		p.passwordHash = string(h)
		w.save(p)
		return nil
	}
	sp, ok, err := w.store.Load(name)
	if err != nil {
		return err
	}
	if !ok {
		return nil // nothing to set on a non-existent character
	}
	sp.PasswordHash = string(h)
	return w.store.Save(sp)
}
