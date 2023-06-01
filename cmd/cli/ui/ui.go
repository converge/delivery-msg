package ui

import (
	"delivery-msg/internal/domain"
	"delivery-msg/pkg"
	"encoding/json"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nats-io/nats.go"
)

func CreateUITable(rows []table.Row) table.Model {

	columns := []table.Column{
		{Title: "Tracking Code", Width: 36},
		{Title: "Source", Width: 12},
		{Title: "Destination", Width: 14},
		{Title: "Status", Width: 14},
		{Title: "Created", Width: 20},
		{Title: "Modified", Width: 20},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func UpdateTable(data *[]domain.Delivery, natsData *nats.Msg) ([]table.Row, error) {

	var newData domain.Delivery
	err := json.Unmarshal(natsData.Data, &newData)
	if err != nil {
		return nil, err
	}

	for k, v := range *data {
		if v.TrackingCode == newData.TrackingCode {
			(*data)[k] = newData
		}
	}

	rows, err := TransformDbDataToRows(*data)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func TransformDbDataToRows(data []domain.Delivery) ([]table.Row, error) {
	var rows []table.Row
	for _, v := range data {
		rows = append(rows, table.Row{
			v.TrackingCode,
			v.SourceAddress,
			v.DestinationAddress,
			v.Status,
			v.Created,
			v.Modified,
		})
	}

	return rows, nil
}

func WaitForActivity(newCh chan *nats.Msg) tea.Cmd {
	return func() tea.Msg {
		return pkg.NatsListener(<-newCh)
	}
}
