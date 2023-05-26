package main

import (
	"context"
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
	deliveryData *pkg.DeliveryData
	baseStyle    lipgloss.Style
}

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

	//var rows []table.Row
	deliveryData, _, err := ui.PopulateInitialData(dbData)
	if err != nil {
		return nil, err
	}
	//m.deliveryData = new(pkg.DeliveryData)
	m.deliveryData = &deliveryData
	//x.data = DeliveryData

	return deliveryData, nil
}

func fetchBackground(m *model) tea.Cmd {

	data, err := fetchData(m)
	if err != nil {
		fmt.Println(err) // todo:
	}

	m.deliveryData = &data
	return func() tea.Msg {
		return m.deliveryData
	}

}

func (m model) Init() tea.Cmd {
	// this will be called later by bubble tea
	//return fetchBackground(&m)
	return nil
}

func checkForUpdate(m *model, ctx context.Context) error {

	//NATSMsg, err := m.sub.NextMsg(60 * time.Second)
	//fmt.Println("Waiting for message...")
	NATSMsg, err := m.sub.NextMsgWithContext(ctx)
	if err != nil {
		return err
	}
	//fmt.Println("done...")
	var data domain.Delivery
	err = json.Unmarshal(NATSMsg.Data, &data)
	if err != nil {
		return err
	}

	for k, _ := range *m.deliveryData {
		if k == data.TrackingCode {
			(*m.deliveryData)[k] = data
		}
	}

	//err = ui.UpdateData(*m.deliveryData, data)
	//if err != nil {
	//	return err
	//}
	var rows []table.Row
	for k, v := range *m.deliveryData {
		rows = append(rows, table.Row{k, v.SourceAddress, v.DestinationAddress, v.Status, v.Created, v.Modified})
	}

	m.table.SetRows(rows)

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "u":
			err := checkForUpdate(&m, context.Background())
			if err != nil {
				log.Error(err)
				return nil, nil
			}
		}
		//case *pkg.DeliveryData:
		//	// todo: fix it
		//	m.deliveryData = msg
		//	var rows []table.Row
		//	for k, v := range *m.deliveryData {
		//		rows = append(rows, table.Row{k, v.SourceAddress, v.DestinationAddress, v.Status, v.Created, v.Modified})
		//	}
		//
		//	m.table.SetRows(rows)
		//	time.Sleep(1 * time.Second)
		//	return m, fetchBackground(&m)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.baseStyle.Render(m.table.View()) + "\n"
}

func main() {

	m := model{}
	m.baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	var rows []table.Row
	t := ui.CreateUITable(rows)
	m.table = t

	cfg := config.ReadConfig()

	data, err := fetchData(&m)
	if err != nil {
		return
	}
	for k, v := range data {
		rows = append(rows, table.Row{k, v.SourceAddress, v.DestinationAddress, v.Status, v.Created, v.Modified})
	}
	m.table.SetRows(rows)

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

	if _, err = tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
	//wg.Done()
	//ctx.Done()
	//}()

	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			break
	//		default:
	//			rows, err := checkForUpdate(&m, ctx)
	//			if err != nil {
	//				return
	//			}
	//			m.table.SetRows(rows)
	//			time.Sleep(3 * time.Second)
	//			m.table.UpdateViewport()
	//		}
	//	}
	//}()

	//wg.Wait()

	//// go routine NATS MSG
	//wg.Add(1)
	//go func() {
	//
	//	defer wg.Done()
	//
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			cancel()
	//		default:
	//			checkForUpdate(&m, ctx)
	//			time.Sleep(3 * time.Second)
	//		}
	//	}
	//}()
	//
	//<-runningUI
	//wg.Wait()

}
