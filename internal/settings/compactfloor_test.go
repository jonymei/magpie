package settings

import (
	"maps"
	"os"
	"path/filepath"
	"testing"
)

// A threshold under CompactFloor is a unit left off, not a window: 400 typed
// in Settings → Long conversations for 400K read every model's window as 400
// tokens and had Codex and Claude Code compact on every turn. A file that
// already holds one reads as no threshold at all, so the working window
// stands, and a caller's own map is copied out rather than written into.
func TestCompactUnderTheFloorIsNoThreshold(t *testing.T) {
	h := t.TempDir()
	t.Setenv("HOME", h)
	t.Setenv("USERPROFILE", h)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(h, ".config"))
	if err := os.MkdirAll(Dir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(), []byte(`{"compactAt":400,"modelCompacts":{"a/b":400,"a/*":200000}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	s := Load()
	if s.CompactAt != 0 || s.Compact() != WorkingWindow {
		t.Errorf("compactAt 400 read as %d, compacting at %d", s.CompactAt, s.Compact())
	}
	if want := map[string]int{"a/*": 200000}; !maps.Equal(s.ModelCompacts, want) {
		t.Errorf("modelCompacts %v; want %v", s.ModelCompacts, want)
	}

	own := map[string]int{"a/b": 300, "a/*": 500000}
	if err := Save(Settings{CompactAt: 300, ModelCompacts: own}); err != nil {
		t.Fatal(err)
	}
	if s := Load(); s.CompactAt != 0 || !maps.Equal(s.ModelCompacts, map[string]int{"a/*": 500000}) {
		t.Errorf("saved: compactAt %d, modelCompacts %v", s.CompactAt, s.ModelCompacts)
	}
	if want := map[string]int{"a/b": 300, "a/*": 500000}; !maps.Equal(own, want) {
		t.Errorf("the caller's map became %v; want %v", own, want)
	}
}
