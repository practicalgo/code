package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/practicalgo/code/appendix-b/http-server/config"
	"github.com/practicalgo/code/appendix-b/http-server/types"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabaseConn(
	dbAddr, dbName, dbUser, dbPassword string,
) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		dbUser, dbPassword,
		dbAddr, dbName,
	)
	return sql.Open("mysql", dsn)
}

func UpdateDb(ctx context.Context, config *config.AppConfig, row types.PkgRow) error {
	conn, err := config.Db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			config.Logger.Debug().Msg(err.Error())
		}
	}()

	_, spanTx := config.Trace.Start(
		config.SpanCtx, "sql:transaction",
	)
	defer spanTx.End()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO packages (owner_id, name, version, object_store_id) 
		VALUES (?,?,?,?);`,
		row.OwnerId, row.Name, row.Version, row.ObjectStoreId,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			config.Logger.Err(rollbackErr)
		}
		return err
	}
	nRows, err := result.RowsAffected()
	if err != nil {
		return tx.Rollback()
	}
	if nRows != 1 {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			config.Logger.Err(rollbackErr)
		}
		return fmt.Errorf(
			"expected 1 row to be inserted, Got: %v",
			nRows,
		)
	}
	return tx.Commit()
}

func QueryDb(
	config *config.AppConfig, params types.PkgQueryParams,
) ([]types.PkgRow, error) {
	ctx := context.Background()
	conn, err := config.Db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	args := []interface{}{}
	conditions := []string{}
	if params.OwnerId != 0 {
		conditions = append(conditions, "owner_id=?")
		args = append(args, params.OwnerId)
	}
	if len(params.Name) != 0 {
		conditions = append(conditions, "name=?")
		args = append(args, params.Name)
	}
	if len(params.Version) != 0 {
		conditions = append(conditions, "version=?")
		args = append(args, params.Version)
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no query conditions found")
	}

	query := fmt.Sprintf(
		"SELECT * FROM packages WHERE %s",
		strings.Join(conditions, " AND "),
	)
	_, spanQuery := config.Trace.Start(
		config.SpanCtx, "sql:query",
	)
	defer spanQuery.End()

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pkgResults []types.PkgRow

	for rows.Next() {
		var pkg types.PkgRow
		if err := rows.Scan(
			&pkg.OwnerId, &pkg.Name, &pkg.Version,
			&pkg.ObjectStoreId, &pkg.Created,
		); err != nil {
			return nil, err
		}
		pkgResults = append(pkgResults, pkg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pkgResults, nil
}
