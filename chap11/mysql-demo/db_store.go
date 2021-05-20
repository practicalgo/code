package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func getDatabaseConn(
	dbAddr, dbName, dbUser, dbPassword string,
) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
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

	columnNames := []string{
		"owner_id", "name", "version", "object_store_id",
	}
	valuesPlaceholder := []string{"?", "?", "?", "?"}
	args := []interface{}{row.OwnerId, row.Name, row.Version, row.ObjectStoreId}

	if row.RepoURL.Valid {
		columnNames = append(columnNames, "repo_url")
		valuesPlaceholder = append(valuesPlaceholder, "?")
		args = append(args, row.RepoURL.String)
	}
	query := fmt.Sprintf(
		"INSERT INTO packages (%s) VALUES (%s);", strings.Join(columnNames, ","),
		strings.Join(valuesPlaceholder, ","),
	)

	result, err := conn.ExecContext(
		ctx, query, args...,
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
			&pkg.ObjectStoreId, &pkg.RepoURL, &pkg.Created,
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
