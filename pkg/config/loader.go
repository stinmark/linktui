package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// AppConfig mirrors the config.toml file structure
type AppConfig struct {
	Window struct {
		Width  int `toml:"width"`
		Height int `toml:"height"`
		Hints bool `toml:"hints"`
	} `toml:"window"`

	Colors struct {
		Foreground          string `toml:"foreground"`
		Background          string `toml:"background"`
		Accent              string `timl:"accent"`
		Highlight           string `toml:"highlight"`
		HighlightBackground string `toml:"Highlight_background"`
		Muted               string `toml:"muted"`
		Border              string `toml:"border"`
		PopupBackground     string `toml:"popup_background"`
		LogBackground       string `toml:"log_background"`
		DividerBackground   string `toml:"divider_background"`
		Cursor              string `toml:"cursor"`
	} `toml:"colors"`
}

var (
	HasCustomWidth  bool
	HasCustomHeight bool
)

// LoadConfig opens the TOML file, maps configurations, and initializes styles
func LoadConfig(filePath string) error {
	cfg := AppConfig{}

	// 1. Assign current file values as fallbacks in case file elements are missing
	cfg.Window.Width = WindowWidth
	cfg.Window.Height = WindowHeight
	cfg.Window.Hints = ShowHints
	cfg.Colors.Foreground = ColorForeground
	cfg.Colors.Background = ColorBackground
	cfg.Colors.Border = ColorBorder
	cfg.Colors.Accent = ColorAccent
	cfg.Colors.Highlight = ColorHighlight
	cfg.Colors.HighlightBackground = ColorHighlightBackground
	cfg.Colors.PopupBackground = ColorPopupBackground
	cfg.Colors.Muted = ColorMuted
	cfg.Colors.DividerBackground = ColorDivider
	cfg.Colors.LogBackground = ColorLogBackground
	cfg.Colors.Cursor = ColorCursor

	// 2. Decode the TOML file if found and track custom window overrides
	if _, err := os.Stat(filePath); err == nil {
		if meta, err := toml.DecodeFile(filePath, &cfg); err == nil {
			HasCustomWidth = meta.IsDefined("window", "width")
			HasCustomHeight = meta.IsDefined("window", "height")
		}
	}

	// 3. Set global window constraints if specified in config
	if HasCustomWidth {
		WindowWidth = cfg.Window.Width
	}
	if HasCustomHeight {
		WindowHeight = cfg.Window.Height
	}

	// Recalculate dependent grid layout variables
	RecalculateDimensions()

	// 4. Update other booleans
	ShowHints = cfg.Window.Hints

	// 5. Update style colors
	ColorForeground = cfg.Colors.Foreground
	ColorBackground = cfg.Colors.Background
	ColorBorder = cfg.Colors.Border
	ColorLogBackground = cfg.Colors.LogBackground
	ColorAccent = cfg.Colors.Accent
	ColorHighlight = cfg.Colors.Highlight
	ColorHighlightBackground = cfg.Colors.HighlightBackground
	ColorPopupBackground = cfg.Colors.PopupBackground
	ColorMuted = cfg.Colors.Muted
	ColorDivider = cfg.Colors.DividerBackground
	ColorCursor = cfg.Colors.Cursor

	// 5. Build Lipgloss styles
	InitStyles()

	return nil
}


