package app

import (
	"fmt"

	"github.com/austinemk/linktui/pkg/config"
	"github.com/austinemk/linktui/pkg/vpn"
	"github.com/austinemk/linktui/pkg/wifi"

	tea "charm.land/bubbletea/v2"
)

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.TerminalWidth = windowMsg.Width
		m.TerminalHeight = windowMsg.Height

		// 1. Pick base width/height (Terminal size by default, Config override if specified)
		targetWidth := windowMsg.Width
		if config.HasCustomWidth {
			targetWidth = config.WindowWidth
		}

		targetHeight := windowMsg.Height
		if config.HasCustomHeight {
			targetHeight = config.WindowHeight
		}

		// 2. Cap dimensions with min(target, MaxWindow)
		if targetWidth > config.MaxWindowWidth {
			targetWidth = config.MaxWindowWidth
		}
		if targetHeight > config.MaxWindowHeight {
			targetHeight = config.MaxWindowHeight
		}

		// 3. Check against minimum required dimensions
		if targetWidth < config.MinWindowWidth || targetHeight < config.MinWindowHeight {
			m.SizeError = fmt.Sprintf("⚠️ Terminal screen too small!\n\n  Current: %dx%d\n  Minimum required: %dx%d",
				targetWidth, targetHeight, config.MinWindowWidth, config.MinWindowHeight)
			return m, nil
		}
		m.SizeError = ""

		// 4. Update global sizes and recalculate grid dependents
		config.WindowWidth = targetWidth
		config.WindowHeight = targetHeight
		config.RecalculateDimensions()
		config.InitStyles()

		// Pass computed dimensions down to active view models
		windowMsg.Width = targetWidth
		windowMsg.Height = targetHeight
		msg = windowMsg
	}

	if m.SizeError != "" {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		oldTab := m.ActiveTab

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			isWifiInput := m.ActiveTab == WifiTab && m.WifiView.UIState == wifi.StatePasswordInput
			isVpnInput := m.ActiveTab == VpnTab && m.VpnView.UIState == vpn.StateAddForm
			if !isWifiInput && !isVpnInput {
				return m, tea.Quit
			}

		case "tab", "pagedown", "pgdown":
			m.ActiveTab = (m.ActiveTab + 1) % 3

		case "shift+tab", "pageup", "pgup":
			m.ActiveTab = (m.ActiveTab - 1 + 3) % 3
		}

		if oldTab != m.ActiveTab {
			initCmd := m.lazyLoadTab(m.ActiveTab)
			if initCmd != nil {
				cmds = append(cmds, initCmd)
			}
			return m, tea.Batch(cmds...)
		}
	}

	switch m.ActiveTab {
	case WifiTab:
		m.WifiView, cmd = m.WifiView.Update(msg)
		cmds = append(cmds, cmd)
	case BluetoothTab:
		m.BtView, cmd = m.BtView.Update(msg)
		cmds = append(cmds, cmd)
	case VpnTab:
		m.VpnView, cmd = m.VpnView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
