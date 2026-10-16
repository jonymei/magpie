package cli

import (
	"testing"

	"github.com/yetone/magpie/internal/catalog"
	"github.com/yetone/magpie/internal/settings"
)

// magpie model compact off runs Codex's and Claude Code's conversations to
// the model's whole window, on compacts them at the working window again
// (Max on Discord found no way to it but the app's Settings); anything else
// is refused.
func TestModelCompactCmd(t *testing.T) {
	groupsHome(t)
	catalog.Changed = nil // no agent's files are rewritten here
	if err := modelCmd([]string{"compact"}); err != nil {
		t.Fatal(err)
	}
	if err := modelCmd([]string{"compact", "off"}); err != nil || !settings.Load().FullContext {
		t.Fatal(err, settings.Load().FullContext)
	}
	if err := modelCmd([]string{"compact"}); err != nil {
		t.Fatal(err)
	}
	if err := modelCmd([]string{"compact", "on"}); err != nil || settings.Load().FullContext {
		t.Fatal(err, settings.Load().FullContext)
	}
	if err := modelCmd([]string{"full-context", "full"}); err != nil || !settings.Load().FullContext {
		t.Fatal(err, settings.Load().FullContext)
	}
	if err := modelCmd([]string{"compact", "maybe"}); err == nil {
		t.Fatal("took maybe")
	}
	// 400 is 400 tokens, not the 400K it was typed for: refused, and the
	// setting it would have been compacted every turn at stands
	before := settings.Load()
	if err := modelCmd([]string{"compact", "400"}); err == nil {
		t.Fatal("took 400")
	}
	if s := settings.Load(); s.FullContext != before.FullContext || s.CompactAt != before.CompactAt {
		t.Fatalf("a refused size changed the setting: %+v", s)
	}
}
