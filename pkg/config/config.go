// Package config for app configurations
package config

// Enforced window constraint thresholds
var (
	WindowWidth  = 75
	WindowHeight = 25
)

var (
	SmallWidth = 30
)

var (
	MinWindowWidth = 20
	MinWindowHeight = 15
	MaxWindowWidth = 80
	MaxWindowHeight = 30
)

// All dependent layout items MUST be vars so they can be recalculated
var (
	ContentHeight         = WindowHeight - 5
	ContentWidth          = WindowWidth - 4
	)

// Header
var (
	HeaderSpacing = (WindowWidth - 20) / 8
)

// Popup box layout variables
var (
	PopupWidth  = (ContentWidth * 3) / 5
	PopupHeight = ContentHeight / 2
)


// booleans
 var (
	 ShowHints = true
)


// RecalculateDimensions updates all internal layout metrics based on current WindowWidth / WindowHeight
func RecalculateDimensions() {
	ContentHeight = max(WindowHeight - 5, 8)

	ContentWidth = WindowWidth - 4
	HeaderSpacing = (WindowWidth - 20) / 8	

	PopupWidth = (ContentWidth * 3) / 5
	PopupHeight = ContentHeight / 3
}

func Truncate(s string, max int) string {
	runes := []rune(s)

	if len(runes) <= max {
		return s
	}

	return string(runes[:max-3]) + "..."
}
