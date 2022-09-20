package db_backup

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"service-hub/global"
	"service-hub/util"
)

type body struct {
	bDb *sql.DB
	sDb *sql.DB
}

func (r body) run() {
	isConn := r.testConnection()
	if !isConn {
		panic("数据库未连接")
	}

	unCreatedTables := r.getUnCreatedTables()
	if unCreatedTables == nil {
		panic("表对比失败")
	}

	// 表结构同步(row_date)

	// 数据同步(page)

}

func (r body) syncTableStruct() bool {

	return false
}

func (r body) getUnCreatedTables() []string {
	tablesSql := "show tables"

	allTableRows, err := r.bDb.Query(tablesSql)
	if err != nil {
		global.Log.Errorf("查询bDb所有表错误 %s", err)
		return nil
	}
	allTable := getRowsSimpleData(allTableRows)
	allTableName := make([]string, len(allTable))
	for index := range allTable {
		tableName := allTable[index].(string)
		//name := *(*string)(unsafe.Pointer(&tableName))
		allTableName = append(allTableName, tableName)
	}

	haveTableRows, err := r.sDb.Query(tablesSql)
	if err != nil {
		global.Log.Errorf("查询sDb所有表错误 %s", err)
		return nil
	}
	haveTable := getRowsSimpleData(haveTableRows)
	haveTableName := make([]string, len(allTable))
	for index := range haveTable {
		tableName := haveTable[index].(string)
		//name := *(*string)(unsafe.Pointer(&tableName))
		haveTableName = append(haveTableName, tableName)
	}

	dbNameMap := make(map[string]bool, len(haveTableName))
	for _, elem := range haveTableName {
		dbNameMap[elem] = true
	}

	unCreatedTableName := make([]string, 0)
	for _, elem := range allTableName {
		if !dbNameMap[elem] {
			unCreatedTableName = append(unCreatedTableName, elem)
		}
	}

	return unCreatedTableName
}

func (r body) testConnection() bool {
	var result string
	var row *sql.Row
	var err error

	row = r.bDb.QueryRow("select 1")
	err = row.Scan(&result)
	if err != nil {
		return false
	}

	row = r.sDb.QueryRow("select 1")
	err = row.Scan(&result)
	if err != nil {
		return false
	}

	return true
}

func getRowsData(rows *sql.Rows) []map[string]interface{} {
	defer rows.Close()
	types, _ := rows.ColumnTypes()
	var rowParam = make([]interface{}, len(types))
	var rowValue = make([]interface{}, len(types))
	for i, colType := range types {
		rowValue[i] = reflect.New(colType.ScanType())
		rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface()
	}
	res := make([]map[string]interface{}, 0)
	for rows.Next() {
		_ = rows.Scan(rowParam...)
		record := make(map[string]interface{})
		for i, colType := range types {
			if rowValue[i] == nil {
				// 如果是 nil 用空字符串代替
				record[colType.Name()] = ""
			} else {
				record[colType.Name()] = util.Byte2Str(rowValue[i].([]byte))
			}
		}
		res = append(res, record)
	}
	return res
}

func getRowsSimpleData(rows *sql.Rows) []interface{} {
	defer rows.Close()
	var res = make([]interface{}, 0)
	for rows.Next() {
		var name interface{}
		_ = rows.Scan(&name)
		res = append(res, util.Byte2Str(name.([]byte)))
	}
	return res
}
