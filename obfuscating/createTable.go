package obfuscating

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type createTableView struct {
	table       string
	createTable string
}

func createTableCopy(originalDb, destinationDb *sql.DB, tableName string) error {
	createTableQuery, err := showCreateTable(originalDb, tableName)
	if err != nil {
		return err
	}

	_, err = destinationDb.Exec(*createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func showCreateTable(db *sql.DB, tableName string) (*string, error) {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE %v;", tableName))
	if err != nil {
		return nil, err
	}
	var views []createTableView
	for rows.Next() {
		var view createTableView
		err = rows.Scan(&view.table, &view.createTable)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	if views == nil || len(views) == 0 {
		return nil, fmt.Errorf("showCreateTable didn't return result")
	}
	if len(views) > 1 {
		return nil, fmt.Errorf("showCreateTable returned more than 1 result")
	}
	return &views[0].createTable, nil
}
