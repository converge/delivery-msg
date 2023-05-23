package repositories

import (
	"database/sql"
	"delivery-msg/internal/domain"
	"delivery-msg/pb"
	"github.com/charmbracelet/log"
)

type PostgreSQLRepository struct {
	conn *sql.DB
}

func NewPostgreSQLRepository(conn *sql.DB) PostgreSQLRepository {
	return PostgreSQLRepository{
		conn: conn,
	}
}

func (pgRepo *PostgreSQLRepository) CreateDelivery(
	sourceAddress string,
	destinationAddress string,
	uuid string,
	created string,
	modified string,
) error {

	sqlStmt := "INSERT INTO delivery_service.public.delivery (tracking_code, source_address, destination_address, status, created, modified) VALUES ($1, $2, $3, $4, $5, $6)"

	result, err := pgRepo.conn.Exec(
		sqlStmt,
		uuid,
		sourceAddress,
		destinationAddress,
		pb.StatusEnum_CONFIRMED.String(),
		created,
		modified,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof(
		"Created delivery: tracking_code=%s, source_address=%s, destination_address=%s, rows affected=%d",
		uuid,
		sourceAddress,
		destinationAddress,
		rowsAffected,
	)
	return nil
}

func (pgRepo *PostgreSQLRepository) UpdateDelivery(trackingCode string, modified string, status string) (domain.Delivery, error) {

	sqlStmt := `
   UPDATE delivery_service.public.delivery 
      SET status=$1, 
          modified=$2 
    WHERE tracking_code=$3
RETURNING tracking_code, source_address, destination_address, status, created, modified
 `

	row := pgRepo.conn.QueryRow(sqlStmt, status, modified, trackingCode)

	var delivery domain.Delivery
	err := row.Scan(
		&delivery.TrackingCode,
		&delivery.SourceAddress,
		&delivery.DestinationAddress,
		&delivery.Status,
		&delivery.Created,
		&delivery.Modified,
	)
	if err != nil {
		return domain.Delivery{}, err
	}

	log.Infof(
		"Updated delivery status: tracking_code=%s, status=%s",
		trackingCode,
		status,
	)

	return delivery, nil
}

func (pgRepo *PostgreSQLRepository) GetAllDeliveries() ([]domain.Delivery, error) {

	sqlStmt := "SELECT tracking_code, source_address, destination_address, status, created, modified FROM delivery_service.public.delivery"

	rows, err := pgRepo.conn.Query(sqlStmt)
	if err != nil {
		return nil, err
	}

	var result []domain.Delivery

	for rows.Next() {
		deliveryRow := domain.Delivery{}
		err = rows.Scan(
			&deliveryRow.TrackingCode,
			&deliveryRow.SourceAddress,
			&deliveryRow.DestinationAddress,
			&deliveryRow.Status,
			&deliveryRow.Created,
			&deliveryRow.Modified,
		)

		result = append(result, deliveryRow)

		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
