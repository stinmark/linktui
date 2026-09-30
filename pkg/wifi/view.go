package wifi

import (
	"fmt"

	"github.com/austinemk/linktui/pkg/config"

	"charm.land/lipgloss/v2"
)

func (m Model) View() string {
	if !m.NMStatus {
		return config.Styles.LogBox.Render("NewtorkManager service is not responding. Ensure 'NewtorkManager' is installed and bluetooth is running")
	}
	// 1. Build the base background (table + adapter + hints + errors)
	var bgSegments []string

	if m.Scanning {
		bgSegments = append(bgSegments, m.ScanningBlock())
	} else {
		bgSegments = append(bgSegments, m.SavedBlock())
	}

	bgSegments = append(bgSegments, m.adapterBlock())
	bgSegments = append(bgSegments, m.HintsBlock())

	if m.Err != nil {
		bgSegments = append(bgSegments, config.LogBlock(m.Err.Error()))
	}

	background := lipgloss.JoinVertical(lipgloss.Left, bgSegments...)

	// 2. Overlay popup on top of background (never appended to segments)
	if m.UIState == StatePasswordInput {
		return config.PlaceOverlay(background, m.PasswordBlock())
	}
	if m.UIState == StateSavedActionsMenu {
		return config.PlaceOverlay(background, m.OptionsBlock())
	}

	return background
}
func (m Model) adapterBlock() string {
	linkStat := false
	if m.Adapter.State == "Connected" {
		linkStat = true
	}

	var lines []string

	// Only include device info if window width is larger than SmallWidth (40)
	if config.WindowWidth > config.SmallWidth {
		lines = append(lines, fmt.Sprintf("device: %s", m.Adapter.Interface))
	}

	lines = append(lines, fmt.Sprintf(" connected: %s", map[bool]string{true: "", false: ""}[linkStat]))
	lines = append(lines, fmt.Sprintf(" power: %s", map[bool]string{true: "󰤨 ", false: "󰤭 "}[m.Adapter.Enabled]))

	// Only include status info if window width is larger than SmallWidth (40)
	if config.WindowWidth > config.SmallWidth {
		lines = append(lines, fmt.Sprintf(" status: %s", map[bool]string{true: "scanning", false: "saved"}[m.Scanning]))
	}

	// Join items with padding so they stay legible in tight spaces
	return config.Truncate(config.Styles.AdapterInfo.Render(lipgloss.JoinHorizontal(lipgloss.Center, lines...)), config.TruncateWidth)
}

func (m Model) ScanningBlock() string {
	title := config.Styles.Heading.Render("Nearby Access Points")
	table := m.Table.View()
	return lipgloss.NewStyle().Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			table,
		),
	)
}

func (m Model) SavedBlock() string {
	title := config.Styles.Heading.Render("󰆓 Saved networks")
	table := m.Table.View()

	return lipgloss.NewStyle().Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			table,
		),
	)
}

func (m Model) PasswordBlock() string {
		return config.Styles.PopupStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, 
	"password", string(m.PassInput.View()),
	),)
}

func (m Model) OptionsBlock() string {
	options := []string{"autoconnect/off", "forget"}
	var menuLines []string
	/*menuLines = append(menuLines, lipgloss.NewStyle().Render(fmt.Sprintf("%s options", m.SelectedSaved.Name)))*/
	for i, opt := range options {
		if m.MenuCursor == i {
			menuLines = append(menuLines, config.Styles.HighlightText.Render(opt))
		} else {
			menuLines = append(menuLines, opt)
		}
	}
	popup := config.Styles.PopupStyle.Render(lipgloss.JoinVertical(lipgloss.Center, menuLines...))
	return popup
}

func (m Model) HintsBlock() string {
	actionsHints := ""

	switch m.UIState {
	case StatePasswordInput:
		actionsHints = "esc: close | enter: submit"
	case StateSavedActionsMenu:
		actionsHints = "j/k: nav | backspace: back "
	case StateNormal:
		actionsHints = "j/k: nav | p: power"
		if m.Scanning {
			actionsHints = actionsHints + " | enter: connect | s: scan off"
		} else {
			actionsHints = actionsHints + "| enter: options| s: scan"
		}
	}

	hints := actionsHints + " | q: quit"

	if(config.WindowWidth <= config.TinyWidth) {
		return ""
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		config.DividerBorder(),
		config.Truncate( config.Styles.Hints.Render(hints), config.TruncateWidth),
	)
}
