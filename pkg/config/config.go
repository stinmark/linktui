// Package config for app configurations
package config

// Enforced window constraint thresholds
var (
	WindowWidth  = 75
	WindowHeight = 25
)

var (
	SmallWidth = 50
	TinyWidth = 25
	TruncateWidth = WindowWidth - 5
)

var (
	MinWindowWidth = 20
	MinWindowHeight = 15
	MaxWindowWidth = 80
	MaxWindowHeight = 30
)

// All dependent layout items MUST be vars so they can be recalculated
var (
	OtherContentHeight = 10
	ListHeight         = WindowHeight - OtherContentHeight
	ListWidth          = WindowWidth - 6
	)

// Header
var (
	HeaderSpacing = (WindowWidth - 20) / 8
)

// Popup box layout variables
var (
	PopupWidth  = (ListWidth * 3) / 5
	PopupHeight = ListWidth / 2
)


// RecalculateDimensions updates all internal layout metrics based on current WindowWidth / WindowHeight
func RecalculateDimensions() {
	ListHeight = WindowHeight - OtherContentHeight
	if ListHeight < 2 {
		ListHeight = 2
	}

	ListWidth = WindowWidth - 6
	HeaderSpacing = (WindowWidth - 20) / 8

	TruncateWidth = WindowWidth - 5


	PopupWidth = (ListWidth * 3) / 5
	PopupHeight = ListHeight / 3
}

func Truncate(s string, max int) string {
	runes := []rune(s)

	if len(runes) <= max {
		return s
	}

	return string(runes[:max-3]) + "..."
}
