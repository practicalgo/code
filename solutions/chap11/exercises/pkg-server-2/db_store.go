package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func getDatabaseConn(
	dbAddr, dbName,
	dbUser, dbPassword string,
) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		dbUser, dbPassword,
		dbAddr, dbName,
	)
	return sql.Open("mysql", dsn)
}

func updateDb(config appConfig, row pkgRow) error {
	ctx := context.Background()
	conn, err := config.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	result, err := conn.ExecContext(
		ctx,
		`INSERT INTO packages (owner_id, name, version, object_store_id) 
		VALUES (?,?,?,?);`,
		row.OwnerId, row.Name, row.Version, row.ObjectStoreId,
	)
	if err != nil {
		return err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Couldn't obtain the last insert id")
	} else {
		log.Printf("last Insert id: %v", lastInsertId)
	}
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if nRows != 1 {
		return fmt.Errorf(
			"expected 1 row to be inserted, Got: %v",
			nRows,
		)
	}
	return nil
}

func queryDb(
	config appConfig,
	params pkgQueryParams,

) ([]pkgRow, error) {
	ctx := context.Background()
	conn, err := config.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	args := make([]interface{}, 0)
	query := "SELECT * FROM packages WHERE"
	if len(params.packageName) != 0 {
		query += " name=?"
		args = append(args, params.packageName)
	}
	if len(params.packageVersion) != 0 {
		query += " AND version=?"
		args = append(args, params.packageVersion)
	}
	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pkgResults []pkgRow

	for rows.Next() {
		var pkg pkgRow
		if err := rows.Scan(
			&pkg.OwnerId, &pkg.Name, &pkg.Version,
			&pkg.ObjectStoreId, &pkg.Created,
		); err != nil {
			return nil, err
		}
		pkgResults = append(pkgResults, pkg)
	}
	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pkgResults, nil
}
