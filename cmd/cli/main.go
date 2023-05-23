package main

import (
	"database/sql"
	"delivery-msg/cmd/cli/ui"
	"delivery-msg/config"
	"delivery-msg/internal/domain"
	"delivery-msg/internal/repositories"
	"delivery-msg/pkg"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table        table.Model
	sub          *nats.Subscription
	deliveryData map[string]domain.Delivery
	baseStyle    lipgloss.Style
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "u":
			NATSMsg, err := m.sub.NextMsg(60 * time.Second)
			if err != nil {
				log.Error(err)
			}
			var data domain.Delivery
			err = json.Unmarshal(NATSMsg.Data, &data)
			if err != nil {
				log.Error(err)
			}

			m.deliveryData[data.TrackingCode] = data

			err = ui.UpdateData(m.deliveryData, data)
			if err != nil {
				log.Error(err)
			}
			var rows []table.Row
			for k, v := range m.deliveryData {
				rows = append(rows, table.Row{k, v.SourceAddress, v.DestinationAddress, v.Status, v.Created, v.Modified})
			}
			m.table.SetRows(rows)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.baseStyle.Render(m.table.View()) + "\n"
}

func createUITable(rows []table.Row) table.Model {

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
func main() {

	m := model{}
	m.deliveryData = make(map[string]domain.Delivery)
	m.baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	cfg := config.ReadConfig()
	databaseURL := cfg.DatabaseUrl
	if databaseURL == "" {
		log.Error("DATABASE_URL is not set")
		return
	}

	pgConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Error(err)
		return
	}
	defer pgConn.Close()

	dbRepository := repositories.NewPostgreSQLRepository(pgConn)
	dbData, err := dbRepository.GetAllDeliveries()
	if err != nil {
		log.Error(err)
		return
	}

	var rows []table.Row
	deliveryData, rows, err := ui.PopulateInitialData(dbData)
	if err != nil {
		log.Error(err)
		return
	}
	m.deliveryData = deliveryData

	t := createUITable(rows)
	m.table = t

	nc, err := nats.Connect(pkg.NATSHostPost, nats.UserInfo(cfg.NATSUser, cfg.NATSPassword))
	if err != nil {
		log.Error(err)
		return
	}
	defer nc.Close()

	m.sub, err = nc.SubscribeSync("nats_development")
	if err != nil {
		log.Error(err)
		return
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
