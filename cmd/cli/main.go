package main

import (
	"database/sql"
	"delivery-msg/cmd/cli/ui"
	"delivery-msg/config"
	"delivery-msg/internal/domain"
	"delivery-msg/internal/repositories"
	"delivery-msg/pkg"
	"encoding/json"
	"errors"
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
	sub          *nats.Subscription
	deliveryData pkg.DeliveryData
	newChField   chan pkg.DeliveryData
	chReceiver   chan *nats.Msg
	baseStyle    lipgloss.Style
}

type testingDev *nats.Msg

func fetchData(m *model) (pkg.DeliveryData, error) {

	cfg := config.ReadConfig()
	databaseURL := cfg.DatabaseUrl
	if databaseURL == "" {
		log.Error("DATABASE_URL is not set")
		return nil, errors.New("DATABASE_URL is not set")
	}

	pgConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	defer pgConn.Close()

	dbRepository := repositories.NewPostgreSQLRepository(pgConn)
	dbData, err := dbRepository.GetAllDeliveries()
	if err != nil {
		return nil, err
	}

	deliveryData, _, err := ui.PopulateInitialData(dbData)
	if err != nil {
		return nil, err
	}
	m.deliveryData = deliveryData

	return deliveryData, nil
}

func waitForActivity(newCh chan *nats.Msg) tea.Cmd {
	return func() tea.Msg {
		return testingDev(<-newCh)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.chReceiver),
	)
}

func doUpdate(m *model, natsData *nats.Msg) error {

	var data domain.Delivery
	err := json.Unmarshal(natsData.Data, &data)
	if err != nil {
		return nil
	}

	for k, v := range *m.deliveryData {
		if v.TrackingCode == data.TrackingCode {
			(*m.deliveryData)[k] = data
			//v.Status = data.Status
			//v.SourceAddress = data.SourceAddress
			//v.DestinationAddress = data.DestinationAddress
			// todo: modified?
		}
	}

	//err = ui.UpdateData(*m.deliveryData, data)
	//if err != nil {
	//	return err
	//}
	var rows []table.Row
	for _, v := range *m.deliveryData {
		rows = append(rows, table.Row{v.TrackingCode, v.SourceAddress, v.DestinationAddress, v.Status, v.Created, v.Modified})
	}

	m.table.SetRows(rows)

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {

	case testingDev:
		err := doUpdate(&m, msg)
		if err != nil {
			log.Error(err)
			return nil, nil
		}
		return m, waitForActivity(m.chReceiver)

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

	m := model{
		newChField: make(chan pkg.DeliveryData),
	}
	m.baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	var rows []table.Row
	t := ui.CreateUITable(rows)
	m.table = t

	data, err := fetchData(&m)
	if err != nil {
		log.Error(err)
		return
	}
	for _, v := range *data {
		rows = append(rows, table.Row{v.TrackingCode, v.SourceAddress, v.DestinationAddress, v.Status, v.Created, v.Modified})
	}
	m.table.SetRows(rows)

	cfg := config.ReadConfig()
	nc, err := nats.Connect(pkg.NATSHostPost, nats.UserInfo(cfg.NATSUser, cfg.NATSPassword))
	if err != nil {
		log.Error(err)
		return
	}
	defer nc.Close()

	m.chReceiver = make(chan *nats.Msg, 64)
	m.sub, err = nc.ChanSubscribe("nats_development", m.chReceiver)
	if err != nil {
		log.Error(err)
		return
	}

	defer func(sub *nats.Subscription) {
		_, err := sub.Dropped()
		if err != nil {
			fmt.Println(err)
		}
	}(m.sub)

	if _, err = tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
