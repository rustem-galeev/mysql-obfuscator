package obfuscating

import "fmt"

func ValidateObfuscationModel(model map[string][]Column, dbConnInfo ConnectionInfo) error {
	dbInfo, err := GetSchemaInfo(dbConnInfo)
	if err != nil {
		return err
	}

	if len(model) != len(dbInfo) {
		return fmt.Errorf("count of tables aren't equal in model and in schema")
	}

	for mTableName, mColumnsSlice := range model {
		dColumnsSlice, contains := dbInfo[mTableName]
		if !contains {
			return fmt.Errorf("schema hasn't table %v", mTableName)
		}

		if len(mColumnsSlice) != len(dColumnsSlice) {
			return fmt.Errorf("count of columns of table %v aren't equal in model and in schema", mTableName)
		}

		mColumns := castColumnsSliceToMap(mColumnsSlice)
		dColumns := castColumnsSliceToMap(dColumnsSlice)

		for _, mColumn := range mColumns {
			dColumn, contains := dColumns[mColumn.Name]
			if !contains {
				return fmt.Errorf("table in schema hasn't column %v in table %v", mColumn.Name, mTableName)
			}

			//names equality are guaranteed by the fact that dCoulmn was got from map
			if mColumn.Type != dColumn.Type {
				return fmt.Errorf("Types of column %v in table %v aren't equal"+
					" in model in schema ", mColumn.Name, mTableName)
			}

			if mColumn.NeedToObfuscate && !dColumn.NeedToObfuscate {
				return fmt.Errorf("you cann't obfuscate this column."+
					" Table name: %v, Column name: %v, Type: %v", mTableName, dColumn.Name, dColumn.Type)
			}

			if mColumn.IsPrimaryKey != dColumn.IsPrimaryKey {
				return fmt.Errorf("isPrimaryKey values in model and in schema aren't equal."+
					" Table name: %v, Column name: %v", mTableName, dColumn.Name)
			}
		}
	}

	return nil
}

func GetSchemaInfo(dbConnInfo ConnectionInfo) (map[string][]Column, error) {
	db, err := openDbConnection(dbConnInfo)
	if err != nil {
		return nil, err
	}

	tables, err := getTables(db, dbConnInfo.Schema)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]Column)
	for _, table := range tables {
		columns, err := getColumnsInfo(db, table)
		if err != nil {
			return nil, err
		}
		result[table] = columns
	}
	return result, nil
}

func castColumnsSliceToMap(slice []Column) map[string]Column {
	result := make(map[string]Column)
	for _, c := range slice {
		result[c.Name] = c
	}
	return result
}
