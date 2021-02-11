package obfuscating

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"obfuscator/encoding"
	"strings"
)

type RawColumn struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default *string
	Extra   *string
}

const (
	primaryKeyWord = "PRI"
)

func getColumnsInfo(db *sql.DB, tableName string) ([]Column, error) {
	rawColumns, err := showColumns(db, tableName)
	if err != nil {
		return nil, err
	}
	hasPrimaryKey := false
	var columns []Column
	for _, rawColumn := range rawColumns {
		column := Column{}
		column.Name = rawColumn.Field
		column.Type = rawColumn.Type
		column.NeedToObfuscate = needToObfuscate(rawColumn)
		if rawColumn.Key == primaryKeyWord {
			column.IsPrimaryKey = true
			hasPrimaryKey = true
		}
		columns = append(columns, column)
	}
	if !hasPrimaryKey {
		return nil, fmt.Errorf("table %v hasn't primary key", tableName)
	}
	return columns, nil
}

func showColumns(db *sql.DB, tableName string) ([]RawColumn, error) {
	rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM %v", tableName))
	if err != nil {
		return nil, err
	}
	var columns []RawColumn
	for rows.Next() {
		column := &RawColumn{}
		err = rows.Scan(&column.Field, &column.Type, &column.Null, &column.Key, &column.Default, &column.Extra)
		if err != nil {
			return nil, err
		}
		columns = append(columns, *column)
	}
	return columns, nil
}

func needToObfuscate(column RawColumn) bool {
	//can take PRI, UNI, MUL values. We don't obfuscate columns with them not to violate PRIMARY KEY, FOREIGN KEY,
	//UNIQUE constraints and to save indexes structure
	//(all indexes in a usual case will be copied due to using SHOW CREATE TABLE statement(tested on BTREE, HASH))
	if column.Key != "" {
		return false
	}
	t := column.Type
	if t == encoding.TinyintType || t == encoding.SmallintType || t == encoding.MediumintType || t == encoding.IntType || t == encoding.BigintType ||
		t == encoding.UTinyintType || t == encoding.USmallintType || t == encoding.UMediumintType || t == encoding.UIntType || t == encoding.UBigintType ||
		t == encoding.FloatType || t == encoding.DoubleType || strings.HasPrefix(t, encoding.DecimalType) ||
		strings.HasPrefix(t, encoding.CharType) || strings.HasPrefix(t, encoding.VarcharType) ||
		t == encoding.TinytextType || t == encoding.TextType || t == encoding.MediumtextType || t == encoding.LongtextType {
		return true
	}
	return false
}
