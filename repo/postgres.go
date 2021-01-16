package repo

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"stpCommon/model"
	"sync"
)

type PostgresConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	DBName     string
	SSLMode    string
	SchemaName string
}

type PostgresRepo struct {
	sync.Mutex
	db         *sql.DB
	schemaName string
}

func NewPostgresRepo(cfg *PostgresConfig) (*PostgresRepo, error) {
	logrus.Debug("connecting to postgres database")
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Debug("connected to postgres to postgres database")

	repo := PostgresRepo{
		Mutex:      sync.Mutex{},
		db:         db,
		schemaName: cfg.SchemaName,
	}

	return &repo, nil
}

func getQueryForExportTable(schemaName string) string {
	res := fmt.Sprintf("select table_schema, table_name, columns_list, keys_list, order_column_name, "+
		"order_column_value, timestamp_column_name, operation_column_name, topic_name  "+
		"from %s.exports e where executing = 0 limit 1", schemaName)
	logrus.Debug("generated query from getting export is: ", res)
	return res
}

func getQueryForExportExecute(schemaName string) string {
	res := fmt.Sprintf("update %s.exports set executing = 1 where table_schema =$1 and table_name = $2", schemaName)
	return res
}

func getQueryForExportStatus(schemaName string) string {
	res := fmt.Sprintf("update %s.exports set executing = 0, order_column_value = $3 where table_schema =$1 and table_name = $2", schemaName)
	return res
}

func (r *PostgresRepo) GetExport() (*model.Export, error) {
	r.Lock()
	defer r.Unlock()

	var export model.Export
	query := getQueryForExportTable(r.schemaName)

	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error("could not start transaction")
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(query).Scan(
		&export.TableSchema,
		&export.TableName,
		&export.ColumnsList,
		&export.KeysList,
		&export.OrderColumnName,
		&export.OrderColumnFromValue,
		&export.TimestampColumnName,
		&export.OperationColumnName,
		&export.TopicName)
	if err == sql.ErrNoRows {
		logrus.Debug("no rows selected for export")
		export.Command = model.ExportSkipCommand
		return &export, nil
	}

	if err != nil {
		logrus.Error("error getting export: ", err.Error())
		return nil, err
	}

	updateQuery := getQueryForExportExecute(r.schemaName)
	_, err = tx.Exec(updateQuery, export.TableSchema, export.TableName)
	if err != nil {
		logrus.Error("could not set update flag to 1, ", err.Error())
		return nil, err
	}

	tx.Commit()

	return &export, nil
}

func (r *PostgresRepo) SetExportStatus(status *model.ExportStatus) error {
	query := getQueryForExportStatus(r.schemaName)

	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error("could not start transaction")
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, status.TableSchema, status.TableName, status.OrderColumnToValue)
	if err != nil {
		logrus.Error("could not set update flag to 1, ", err.Error())
		return err
	}

	tx.Commit()
	return nil
}
