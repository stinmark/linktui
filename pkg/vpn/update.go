package vpn

import (
	"time"

	"github.com/austinemk/linktui/pkg/config"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

// Update acts as the central router for incoming messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.UIState {
	case StateAddForm:
		return m.handleFormState(msg)
	case StateImportFile:
		return m.handleFilePickerState(msg)
	case StateActionsMenu:
		return m.handleActionsMenuState(msg)
	}

	switch msg := msg.(type) {
	case NMStatusMsg:
		m.NMStatus = bool(msg)
		return m, nil
	case TunnelsLoadedMsg:
		m.Tunnels = msg.Tunnels
		// Map backend tunnels data cleanly to the UI table rows
		m.syncTableRows()
		return m, nil

	case IPInfoMsg:
		m.IPInfo = msg
		return m, nil

	case ActionSuccessMsg:
		return m, FetchTunnelsCmd()

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

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	}

	var cmd tea.Cmd

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

// syncTableRows translates the internal Tunnels state into viewable table rows
// and recalculates viewport/column dimensions dynamically.
func (m *Model) syncTableRows() {
	var rows []table.Row

	// 1. Sync Table & FilePicker Viewport Dimensions
	m.Table.SetWidth(config.ListWidth)
	m.Table.SetHeight(config.ListHeight)
	m.FilePicker.SetHeight(config.ListHeight)

	// 2. Proportionately distribute column widths to fill 100% of ListWidth
	nameWidth := (config.ListWidth * 22) / 45;
	typeWidth := (config.ListWidth * 10) / 45;
	statusWidth := (config.ListWidth * 8) /45;

	m.Table.SetColumns([]table.Column{
		{Title: "", Width: nameWidth},
		{Title: "", Width: typeWidth},
		{Title: "", Width: statusWidth},
	})

	// 3. Build Table Rows
	m.Table.SetRows(nil)
	for _, t := range m.Tunnels {
		status := "Inactive"
		if t.Active {
			status = "Active 󰌆"
		}
		rows = append(rows, table.Row{t.Name, t.Type, status})
	}

	m.Table.SetRows(rows)

	if m.Table.Cursor() >= len(rows) && len(rows) > 0 {
		m.Table.GotoTop()
		m.Cursor = m.Table.Cursor()
	}
}
