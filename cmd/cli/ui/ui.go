package ui

import (
	"delivery-msg/internal/domain"
	"delivery-msg/pkg"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"time"
)

//func UpdateData(deliveryData *pkg.DeliveryData, data domain.Delivery) error {
//
//	timeOutputLayout := "02/01/2006 15:04:05"
//	created, err := time.Parse(time.RFC3339, data.Created)
//	if err != nil {
//		return err
//	}
//	modified, err := time.Parse("2006-01-02 15:04:05", data.Modified)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//
//	formatCreated := created.Format(timeOutputLayout)
//	formatModified := modified.Format(timeOutputLayout)
//
//	for k, v := range *deliveryData {
//		if v.TrackingCode
//
//	}
//	deliveryData[data.TrackingCode] = domain.Delivery{
//		SourceAddress:      data.SourceAddress,
//		DestinationAddress: data.DestinationAddress,
//		Status:             data.Status,
//		Created:            formatCreated,
//		Modified:           formatModified,
//	}
//	return nil
//}

func PopulateInitialData(dbData []domain.Delivery) (pkg.DeliveryData, []table.Row, error) {

	var deliveryData = make(map[string]domain.Delivery)
	var rows []table.Row

	for _, v := range dbData {

		created, err := time.Parse(time.RFC3339, v.Created)
		if err != nil {
			return nil, nil, err
		}

		modified, err := time.Parse(time.RFC3339, v.Created)
		if err != nil {
			return nil, nil, err
		}

		// model data
		deliveryData[v.TrackingCode] = domain.Delivery{
			SourceAddress:      v.SourceAddress,
			DestinationAddress: v.DestinationAddress,
			Status:             v.Status,
			Created:            created.Format(pkg.TimeOutputLayout),
			Modified:           modified.Format(pkg.TimeOutputLayout),
		}

		// presentation data
		rows = append(rows, table.Row{
			v.TrackingCode,
			v.SourceAddress,
			v.DestinationAddress,
			v.Status,
			created.Format(pkg.TimeOutputLayout),
			modified.Format(pkg.TimeOutputLayout),
		})
	}

	return deliveryData, rows, nil

}

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
