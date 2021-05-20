package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func getDatabaseConn(
	dbAddr, dbName, dbUser, dbPassword string,
) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		dbUser, dbPassword,
		dbAddr, dbName,
	)
	log.Println(dsn)
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
	config appConfig, params pkgQueryParams,
) ([]pkgRow, error) {
	ctx := context.Background()
	conn, err := config.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	args := []interface{}{}
	conditions := []string{}
	if params.ownerId != 0 {
		conditions = append(conditions, "owner_id=?")
		args = append(args, params.ownerId)
	}
	if len(params.name) != 0 {
		conditions = append(conditions, "name=?")
		args = append(args, params.name)
	}
	if len(params.version) != 0 {
		conditions = append(conditions, "version=?")
		args = append(args, params.version)
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no query conditions found")
	}

	query := fmt.Sprintf(
		"SELECT * FROM packages WHERE %s",
		strings.Join(conditions, " AND "),
	)

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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pkgResults, nil
}
