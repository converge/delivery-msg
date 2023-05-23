package ui

import (
	"delivery-msg/internal/domain"
	"delivery-msg/pkg"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/log"
	"time"
)

// func UpdateData(tackingCode string) {
func UpdateData(deliveryData map[string]domain.Delivery, data domain.Delivery) error {

	timeOutputLayout := "02/01/2006 15:04:05"
	created, err := time.Parse(time.RFC3339, data.Created)
	if err != nil {
		return err
	}
	modified, err := time.Parse("2006-01-02 15:04:05", data.Modified)
	if err != nil {
		log.Error(err)
		return err
	}

	formatCreated := created.Format(timeOutputLayout)
	formatModified := modified.Format(timeOutputLayout)

	deliveryData[data.TrackingCode] = domain.Delivery{
		SourceAddress:      data.SourceAddress,
		DestinationAddress: data.DestinationAddress,
		Status:             data.Status,
		Created:            formatCreated,
		Modified:           formatModified,
	}
	return nil
}

// func populateInitialData(m *model, rows []table.Row) ([]table.Row, error) {
func PopulateInitialData(dbData []domain.Delivery) (map[string]domain.Delivery, []table.Row, error) {

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
