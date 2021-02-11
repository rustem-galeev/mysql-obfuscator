package obfuscating

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"obfuscator/config"
	"obfuscator/encoding"
	"strings"
)

const (
	MysqlDriverName     = "mysql"
	selectBySlicesQuery = "SELECT * FROM %v ORDER BY %v LIMIT %v OFFSET %v;"
)

func ObfuscateSchema(model map[string][]Column, processId string,
	originalDbConnInfo, destinationDbConnInfo ConnectionInfo) {
	originalDb, err := openDbConnection(originalDbConnInfo)
	if err != nil {
		writeError(processId, err)
		return
	}
	destinationDb, err := openDbConnection(destinationDbConnInfo)
	if err != nil {
		writeError(processId, err)
		return
	}

	tables, err := getSortedTables(originalDb, originalDbConnInfo.Schema)
	if err != nil {
		writeError(processId, err)
		return
	}

	for _, table := range tables {
		println(table + " copying started")

		err := createTableCopy(originalDb, destinationDb, table)
		if err != nil {
			writeError(processId, err)
			return
		}
		err = obfuscateTable(model[table], table, originalDb, destinationDb)
		if err != nil {
			writeError(processId, err)
			return
		}

		increaseFinished(processId)

		println(table + " copying finished")
	}
}

func obfuscateTable(model []Column, tableName string, originalDb, destinationDb *sql.DB) error {
	//locking writing to table by all sessions until unlocking below
	lockTableQuery := fmt.Sprintf("LOCK TABLES %v READ;", tableName)
	_, err := originalDb.Exec(lockTableQuery)
	if err != nil {
		return err
	}

	orderByValues := getOrderByValues(model)
	i := 0
	for {
		limit := config.GetConfig().Obfuscator.SliceSize
		selectQuery := fmt.Sprintf(selectBySlicesQuery, tableName, orderByValues, limit, limit*i)
		values, err := getValues(selectQuery, originalDb)
		if err != nil {
			return err
		}
		if values == nil || len(values) == 0 {
			break
		}

		err = obfuscateSlice(values, model, tableName, destinationDb)
		if err != nil {
			return err
		}
		i++
	}

	unlockTablesQuery := "UNLOCK TABLES;"
	_, err = originalDb.Exec(unlockTablesQuery)
	if err != nil {
		return err
	}
	return nil
}

func getValues(selectQuery string, db *sql.DB) ([]map[string]*interface{}, error) {
	rows, err := db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rawResult := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i, _ := range rawResult {
		values[i] = &rawResult[i]
	}

	var result []map[string]*interface{}
	for rows.Next() {
		row := make(map[string]*interface{})
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for i, column := range columns {
			row[column] = values[i].(*interface{})
		}
		result = append(result, row)

		rawResult = make([]interface{}, len(columns))
		values = make([]interface{}, len(columns))
		for i, _ := range rawResult {
			values[i] = &rawResult[i]
		}
	}
	return result, nil
}

func obfuscateSlice(data []map[string]*interface{}, model []Column, tableName string, db *sql.DB) error {
	if len(data) == 0 {
		return nil
	}

	valuesTemplate, columnNames := getInsertsTemplate(model, len(data))
	insertQuery := "INSERT INTO " + tableName + " (" + columnNames + ") VALUES " + valuesTemplate + ";"
	var params = make([]interface{}, len(model)*len(data))
	i := 0
	for _, row := range data {
		for _, column := range model {
			valueToInsert := row[column.Name]
			if column.NeedToObfuscate {
				obfuscatedValue, err := encoding.ObfuscateValue(valueToInsert, column.Type)
				if err != nil {
					log.Printf("Error: Encoding value was failed. Table: %v, Column: %v. %v",
						tableName, column.Name, err.Error())
					params[i] = valueToInsert
					err = nil
				} else {
					params[i] = obfuscatedValue
				}
			} else {
				params[i] = valueToInsert
			}
			i++
		}
	}
	_, err := db.Exec(insertQuery, params...)
	if err != nil {
		return err
	}
	return nil
}

func getInsertsTemplate(columns []Column, rowsCount int) (valuesTemplate string, columnsTemplate string) {
	var columnNames []string
	var valueParams []string
	for _, column := range columns {
		columnNames = append(columnNames, column.Name)
		valueParams = append(valueParams, "?")
	}
	valueSlice := "(" + strings.Join(valueParams, ",") + ")"
	var valueSlices []string
	for i := 0; i < rowsCount; i++ {
		valueSlices = append(valueSlices, valueSlice)
	}
	valuesTemplate = strings.Join(valueSlices, ",")
	columnsTemplate = strings.Join(columnNames, ",")
	return
}

func getOrderByValues(columns []Column) string {
	var orderByValues []string
	for _, column := range columns {
		//UNI key doesn't guarantee that there's all unique columns are showed
		if column.IsPrimaryKey {
			orderByValues = append(orderByValues, column.Name)
		}
	}
	//checking that at least one column is primary key is carried out during getting columns info and validating model
	orderByValuesString := strings.Join(orderByValues, ",")
	return orderByValuesString
}

func openDbConnection(connInfo ConnectionInfo) (*sql.DB, error) {
	url := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&interpolateParams=true",
		connInfo.User, connInfo.Password, connInfo.Host, connInfo.Schema)
	db, err := sql.Open(MysqlDriverName, url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.GetConfig().Db.MaxOpenConnections)
	return db, nil
}
