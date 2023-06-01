package main

import (
	"database/sql"
	"delivery-msg/cmd/cli/ui"
	"delivery-msg/config"
	"delivery-msg/internal/domain"
	"delivery-msg/internal/handlers"
	"delivery-msg/internal/repositories"
	"delivery-msg/internal/services"
	"delivery-msg/pkg"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"
)

type model struct {
	table        table.Model
	data         *[]domain.Delivery
	natsReceiver chan *nats.Msg
	baseStyle    lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		ui.WaitForActivity(m.natsReceiver),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case pkg.NatsListener:
		rows, err := ui.UpdateTable(m.data, msg)
		m.table.SetRows(rows)
		if err != nil {
			log.Error(err)
			return nil, nil
		}
		return m, ui.WaitForActivity(m.natsReceiver)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.baseStyle.Render(m.table.View()) + "\n"
}

func main() {

	// initialize model
	m := model{}
	// initialize style
	m.baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	var rows []table.Row
	m.table = ui.CreateUITable(rows)

	// get all delivery data from DB, populate this data into rows []table.Row.
	// bellow we will subscribe to a specific subject and wait for updates, when a new update is received,
	// it will update the model data reference to follow the updated data that was received via NATS message.
	cfg := config.ReadConfig()
	if cfg.DatabaseUrl == "" {
		log.Error("DATABASE_URL is not set")
		return
	}

	pgConn, err := sql.Open("pgx", cfg.DatabaseUrl)
	if err != nil {
		log.Error(err)
		return
	}
	defer pgConn.Close()
	dbRepository := repositories.NewPostgreSQLRepository(pgConn)

	// initialize NATS connection
	nc, err := nats.Connect(pkg.NATSHostPost, nats.UserInfo(cfg.NATSUser, cfg.NATSPassword))
	if err != nil {
		log.Error(err)
		return
	}
	defer nc.Close()

	natsClient, err := services.NewNATSClient(nc)
	if err != nil {
		fmt.Println(err)
	}
	deliveryHandler := handlers.NewDeliveryHandler(&dbRepository, &natsClient)

	dbData, err := deliveryHandler.GetDelivery()
	if err != nil {
		log.Error(err)
		return
	}
	// set initial data reference, this entity will be used as a reference to update the model data
	m.data = dbData

	rows, err = ui.TransformDbDataToRows(*dbData)
	if err != nil {
		log.Error(err)
		return
	}
	m.table.SetRows(rows)

	// enable NATS async subscription
	m.natsReceiver = make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe(pkg.NATSSubject, m.natsReceiver)
	if err != nil {
		log.Error(err)
		return
	}
	defer sub.Drain()

	if _, err = tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
