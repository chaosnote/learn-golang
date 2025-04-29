package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	name     = "chris"
	password = "123456"
	addr     = "192.168.0.236:3306"
	db_name  = "game_dev"
)

func main() {
	var cmd = fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		name,
		password,
		addr,
		db_name,
	)
	log.Println(cmd)

	db, err := sql.Open("mysql", cmd)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`
        SELECT
            TABLE_NAME,
            COLUMN_NAME,
            DATA_TYPE,
			CHARACTER_MAXIMUM_LENGTH,
            COLUMN_KEY,
            IS_NULLABLE,
            COLUMN_DEFAULT
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = ?
        ORDER BY TABLE_NAME, ORDINAL_POSITION
    `, db_name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	schema := make(map[string][]map[string]interface{})
	for rows.Next() {
		var tableName, columnName, dataType, columnKey, isNullable, columnDefault sql.NullString
		var characterMaximumLength sql.NullInt64
		if err := rows.Scan(&tableName, &columnName, &dataType, &characterMaximumLength, &columnKey, &isNullable, &columnDefault); err != nil {
			panic(err)
		}

		tableInfo := schema[tableName.String]
		tableInfo = append(tableInfo, map[string]interface{}{
			"COLUMN_NAME":              columnName.String,
			"DATA_TYPE":                dataType.String,
			"CHARACTER_MAXIMUM_LENGTH": characterMaximumLength.Int64,
			"COLUMN_KEY":               columnKey.String,
			"IS_NULLABLE":              isNullable.String,
			"COLUMN_DEFAULT":           columnDefault.String,
		})
		schema[tableName.String] = tableInfo
	}

	// 生成 CREATE TABLE 語句，包含長度
	for tableName, columns := range schema {
		fmt.Printf("CREATE TABLE IF NOT EXISTS `%s` (\n", tableName)
		for i, col := range columns {
			fmt.Printf("    %s %s", col["COLUMN_NAME"], col["DATA_TYPE"])
			if length, ok := col["CHARACTER_MAXIMUM_LENGTH"].(int64); ok && length > 0 && (col["DATA_TYPE"] == "varchar" || col["DATA_TYPE"] == "char" || col["DATA_TYPE"] == "text" || col["DATA_TYPE"] == "tinytext" || col["DATA_TYPE"] == "mediumtext" || col["DATA_TYPE"] == "longtext") {
				fmt.Printf("(%d)", length)
			}
			if col["COLUMN_KEY"] == "PRI" {
				fmt.Print(" PRIMARY KEY")
			}
			if col["IS_NULLABLE"] == "NO" {
				fmt.Print(" NOT NULL")
			}
			if defaultValue, ok := col["COLUMN_DEFAULT"].(string); ok {
				fmt.Printf(" DEFAULT '%s'", defaultValue)
			}
			if i < len(columns)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Println(");")
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
