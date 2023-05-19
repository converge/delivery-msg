package repositories

import (
	"database/sql"
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

func (pgRepo *PostgreSQLRepository) CreateDelivery(sourceAddress string, destinationAddress string, uuid string) error {

	sqlStmt := "INSERT INTO delivery_service.public.delivery (tracking_code, source_address, destination_address, status, created, modified) VALUES ($1, $2, $3, $4, NOW(), NOW())"

	result, err := pgRepo.conn.Exec(
		sqlStmt,
		uuid,
		sourceAddress,
		destinationAddress,
		pb.StatusEnum_CONFIRMED.String(),
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

func (pgRepo *PostgreSQLRepository) UpdateDelivery(trackingCode string, status string) error {

	sqlStmt := "UPDATE delivery_service.public.delivery SET status=$1, modified=NOW() WHERE tracking_code=$2"

	result, err := pgRepo.conn.Exec(sqlStmt, status, trackingCode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof(
		"Updated delivery status: tracking_code=%s, status=%s, rows affected=%d",
		trackingCode,
		status,
		rowsAffected,
	)

	return nil
}
