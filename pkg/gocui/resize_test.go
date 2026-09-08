package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResizeUsesEventDimensions reproduces the stuck-resize bug: a resize
// event carries the terminal's authoritative new size, which can disagree with
// the screen's cached Screen.Size() for a frame or two. flush() must lay out at
// the event's size; otherwise it caches the stale screen size in maxX/maxY and
// every later flush becomes a no-op, leaving the UI at the old size until the
// app is restarted.
func TestResizeUsesEventDimensions(t *testing.T) {
	g := newTestGui(t)

	// The headless screen stays at its created size (80x24).
	screenW, screenH := Screen.Size()
	assert.Equal(t, 80, screenW)
	assert.Equal(t, 24, screenH)

	// Record the size the manager saw during layout.
	var layoutW, layoutH int
	g.SetManager(ManagerFunc(func(g *Gui) error {
		layoutW, layoutH = g.Size()
		return nil
	}))

	// A resize event reporting the terminal has become 120x50, even though the
	// screen's cached size still reports 80x24.
	g.onResize(&GocuiEvent{Type: eventResize, Width: 120, Height: 50})

	assert.NoError(t, g.flush())

	assert.Equal(t, 120, layoutW, "layout must run with the event's width, not the stale screen size")
	assert.Equal(t, 50, layoutH, "layout must run with the event's height, not the stale screen size")

	// The pending resize is consumed, so a later flush (with no new resize)
	// falls back to the screen's cached size, whatever it is.
	layoutW, layoutH = -1, -1
	assert.NoError(t, g.flush())
	assert.Equal(t, screenW, layoutW)
	assert.Equal(t, screenH, layoutH)
}

// TestResizePendingClearedWhenScreenMatches: after the screen has caught up to
// the event's size, the pending resize must not override subsequent flushes.
func TestResizePendingOverridesStaleScreenSize(t *testing.T) {
	g := newTestGui(t)

	var layoutW, layoutH int
	g.SetManager(ManagerFunc(func(g *Gui) error {
		layoutW, layoutH = g.Size()
		return nil
	}))

	g.onResize(&GocuiEvent{Type: eventResize, Width: 100, Height: 40})
	assert.NoError(t, g.flush())
	assert.Equal(t, 100, layoutW)
	assert.Equal(t, 40, layoutH)

	// The event was consumed; a resize with dimensions equal to the current
	// screen size still clears the pending value and matches.
	g.onResize(&GocuiEvent{Type: eventResize, Width: 80, Height: 24})
	assert.NoError(t, g.flush())
	assert.Equal(t, 80, layoutW)
	assert.Equal(t, 24, layoutH)
}
