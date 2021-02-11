package obfuscating

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//returns tables sorted by possibility to insert without foreign keys violations
func getSortedTables(db *sql.DB, schemaName string) ([]string, error) {
	tablesWithDependencies, err := getTablesWithDependencies(db, schemaName)
	if err != nil {
		return nil, err
	}

	//delete link to itself
	for table := range tablesWithDependencies {
		delete(tablesWithDependencies[table], table)
	}

	var result []string
	for {
		hasMutualDependencies := true
		var dependenciesToDelete []string
		for name, dependencies := range tablesWithDependencies {
			if len(dependencies) == 0 {
				result = append(result, name)
				delete(tablesWithDependencies, name)
				dependenciesToDelete = append(dependenciesToDelete, name)
				hasMutualDependencies = false
			}
		}

		for otherTableName := range tablesWithDependencies {
			for _, dependencyToDelete := range dependenciesToDelete {
				delete(tablesWithDependencies[otherTableName], dependencyToDelete)
			}
		}

		if hasMutualDependencies {
			return nil, fmt.Errorf("the schema has mutual dependencies")
		}
		if len(tablesWithDependencies) == 0 {
			break
		}
	}
	return result, nil
}

func getTablesWithDependencies(db *sql.DB, schemaName string) (map[string]map[string]bool, error) {
	tableNames, err := getTables(db, schemaName)
	if err != nil {
		return nil, err
	}
	tables := make(map[string]map[string]bool)
	for _, tableName := range tableNames {
		dependencies, err := getDependencies(db, schemaName, tableName)
		if err != nil {
			return nil, err
		}

		dependenciesMap := make(map[string]bool)
		// for optimal deleting
		for _, dependency := range dependencies {
			dependenciesMap[dependency] = true
		}

		tables[tableName] = dependenciesMap
	}
	return tables, nil
}

func getTables(db *sql.DB, schemaName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%v'; ", schemaName))
	if err != nil {
		return nil, err
	}
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func getDependencies(db *sql.DB, schemaName, tableName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT referenced_table_name FROM information_schema.key_column_usage"+
		" WHERE  referenced_table_schema = '%v' AND table_name = '%v'; ", schemaName, tableName))
	if err != nil {
		return nil, err
	}
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
