package wifi

import (
	"time"

	"github.com/austinemk/linktui/pkg/config"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// 1. State-based Structural Intercepts
	switch m.UIState {
	case StatePasswordInput:
		return m.handlePasswordInput(msg)
	case StateSavedActionsMenu:
		return m.handleSavedActionsMenu(msg)
	}

	// 2. Normal State Core Navigation Loop
	switch msg := msg.(type) {
	case NMStatusMsg:
		m.NMStatus = bool(msg)
	case InfoLoadedMsg:
		return m.handleInfoLoaded(msg)

	case ScanFinishedMsg:
		return m.handleScanFinished(msg)

	case TickMsg:
		return m.handleTick()

	case AdapterToggledMsg, ActionSuccessMsg:
		return m.handleAdapterOrActionSuccess()

	case ErrMsg:
		m.Err = msg
		m.LogID++
		return m, func() tea.Msg {
			time.Sleep(4 * time.Second) // Display duration before auto-removal
			return ClearLogMsg{ID: m.LogID}
		}
	case ClearLogMsg:
		if msg.ID == m.LogID {
			m.Err = nil
		}
		return m, nil

	// V2 Change: KeyMsg is now KeyPressMsg
	case tea.KeyPressMsg:
		return m.handleKeyInput(msg)
	}

	// 3. Fallback to sub-component updates
	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)

	return m, cmd
}

func (m *Model) syncTableRows() {
	var rows []table.Row

	// 1. Sync table viewport width to latest ListWidth
	m.Table.SetWidth(config.ContentWidth)
	m.Table.SetHeight(config.ContentHeight - 5)


	if m.Scanning {
		m.Table.SetRows(nil)
		col1 := (config.ContentWidth * 2) / 45
		col2 := (config.ContentWidth * 31) / 45
		col3 := (config.ContentWidth * 5) /45
		col4 := (config.ContentWidth * 2)/ 45

		m.Table.SetColumns([]table.Column{
			{Title: "", Width: col1},
			{Title: "", Width: col2},
			{Title: "", Width: col3},
			{Title: "", Width: col4},
		})

		for _, ap := range m.ActiveAPs {
			activeMark := " "
			if ap.IsActive {
				activeMark = ""
			}
			rows = append(rows, table.Row{
				RenderSignal(ap.Strength, ap.Security),
				ap.SSID,
				ap.Security,
				activeMark,
			})
		}
	} else {
		m.Table.SetRows(nil)
		
		// Proportional distribution for saved networks table
		nameWidth := (config.ContentWidth * 30) / 70
		autoWidth := (config.ContentWidth * 5) / 70
		uuidWidth := (config.ContentWidth * 24) /70 

		m.Table.SetColumns([]table.Column{
			{Title: "", Width: nameWidth},
			{Title: "", Width: autoWidth},
			{Title: "", Width: uuidWidth},
		})

		for _, prof := range m.Saved {
			autoStr := " "
			if prof.AutoConnect {
				autoStr = "󰁪"
			}
			uuidShort := ""
			if len(prof.UUID) >= 8 {
				uuidShort = prof.UUID[:8]
			}
			rows = append(rows, table.Row{
				prof.Name,
				autoStr,
				uuidShort,
			})
		}
	}

	m.Table.SetRows(rows)

	if m.Table.Cursor() >= len(rows) && len(rows) > 0 {
		m.Table.GotoTop()
		m.Cursor = m.Table.Cursor()
	}
}
