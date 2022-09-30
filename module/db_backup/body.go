package db_backup

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"regexp"
	"service-hub/global"
	"service-hub/util"
	"strings"
)

type body struct {
	bDb *sql.DB
	sDb *sql.DB
}

func (r body) run() {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Errorf("同步数据表异常: %s", err)
		}
	}()

	isConn := r.testConnection()
	if !isConn {
		panic("数据库未连接")
	}

	unCreatedTables := r.getUnCreatedTables()
	if unCreatedTables == nil {
		panic("表对比失败")
	}

	// 表结构同步(row_date)
	syncTableStruct := r.syncTableStruct(unCreatedTables)
	if !syncTableStruct {
		panic("表结构同步失败")
	}

	// 数据同步(page)
	r.syncData()

}

func (r body) syncData() {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Errorf("同步表数据发生异常 %s", err)
		}
	}()
	queryTableSQL := "show tables"
	allTableRows, err := r.bDb.Query(queryTableSQL)
	if err != nil {
		global.Log.Errorf("查询bDb所有表错误 %s", err)
		return
	}
	allTable := getRowsSimpleData(allTableRows)
	// 表名称
	allTableName := make([]string, len(allTable))
	for index := range allTable {
		allTableName = append(allTableName, allTable[index].(string))
	}
	dbName := global.Global.DbBackup.BDb.Database
	queryColumnSqlTemplate := "select * from information_schema.columns where table_schema = '%s' and table_name = '%s'"
	for _, tableName := range allTableName {
		// 查询表字段SQL
		queryColumnSQL := fmt.Sprintf(queryColumnSqlTemplate, dbName, tableName)
		tableColumnRows, err := r.bDb.Query(queryColumnSQL)
		if err != nil {
			global.Log.Errorf("查询bDb表字段错误 %s %s", err, queryColumnSQL)
			continue
		}
		tableColumn := getRowsData(tableColumnRows)
		// 存放表字段
		tableColumns := make([]string, len(tableColumn))
		for index := range tableColumn {
			columnName := tableColumn[index]["COLUMN_NAME"].(string)
			tableColumns = append(tableColumns, columnName)
		}

		// 获取记录条数

	}

}

func (r body) syncTableStruct(tables []string) bool {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Errorf("同步表结构发生异常 %s", err)
		}
	}()

	if len(tables) == 0 || tables == nil {
		return true
	}
	queryTableCreateSql := "show create table %s"
	alterTableRowDateSQL := "alter table %s add column row_date varchar(10) comment '备份时间'"
	for _, elem := range tables {
		tableCreateRows, err := r.bDb.Query(fmt.Sprintf(queryTableCreateSql, elem))
		if err != nil {
			global.Log.Errorf("获取表结构异常 %s", err)
			return false
		}
		// 原生SQL
		nativeTableCreateSQL := getRowsData(tableCreateRows)[0]["Create Table"].(string)

		// 原生SQL分割()
		begin := strings.Index(nativeTableCreateSQL, "(")
		end := strings.LastIndex(nativeTableCreateSQL, ")")
		// before center end
		nativeTableBeforeSQL, nativeTableCenterSQL, nativeTableEndSQL :=
			nativeTableCreateSQL[:begin+1], nativeTableCreateSQL[begin+1:end], nativeTableCreateSQL[end:]
		// 处理 centerSQL
		nativeTableCenterSQLArray := strings.Split(nativeTableCenterSQL, ",")
		newNativeTableCenterSQLArray := make([]string, 0)
		for _, value := range nativeTableCenterSQLArray {
			if strings.Contains(value, "PRIMARY KEY") ||
				strings.Contains(value, "UNIQUE KEY") {
				continue
			}
			value = strings.ReplaceAll(value, "AUTO_INCREMENT", "")
			newNativeTableCenterSQLArray = append(newNativeTableCenterSQLArray, value)
		}
		newNativeTableCenterSQL := strings.Join(newNativeTableCenterSQLArray, ",")

		// endSQL 处理
		nativeTableEndSQL = strings.ReplaceAll(nativeTableEndSQL, "ENGINE=MyISAM", "")
		nativeTableEndSQL = strings.ReplaceAll(nativeTableEndSQL, "ENGINE=InnoDB", "")
		nativeTableEndSQL = strings.ReplaceAll(nativeTableEndSQL, "ROW_FORMAT=FIXED", "")
		nativeTableEndSQL = regexp.MustCompile(`AUTO_INCREMENT=(\d){1,}`).ReplaceAllString(nativeTableEndSQL, ``)
		nativeTableEndSQL = regexp.MustCompile(`ROW_FORMAT=[A-Z]{1,} `).ReplaceAllString(nativeTableEndSQL, ``)
		// 新原生SQL
		newNativeTableCreateSQL := fmt.Sprintf("%s%s%s", nativeTableBeforeSQL, newNativeTableCenterSQL, nativeTableEndSQL)
		_, err = r.sDb.Exec(newNativeTableCreateSQL)
		if err != nil {
			global.Log.Errorf("%s 新SQL执行错误\n%s\n%s", elem, newNativeTableCreateSQL, err)
			return false
		}
		// 增加 row_date
		_, err = r.bDb.Exec(fmt.Sprintf(alterTableRowDateSQL, elem))
		if err != nil {
			global.Log.Errorf("%s 追加row_date错误: %s", elem, err)
			return false
		}
	}

	return true
}

func (r body) getUnCreatedTables() []string {
	queryTableSQL := "show tables"

	bDbTableRows, err := r.bDb.Query(queryTableSQL)
	if err != nil {
		global.Log.Errorf("查询bDb所有表错误 %s", err)
		return nil
	}
	bDbTable := getRowsSimpleData(bDbTableRows)
	bDbTables := make([]string, len(bDbTable))
	for index := range bDbTable {
		bDbTables = append(bDbTables, bDbTable[index].(string))
	}

	sDbTableRows, err := r.sDb.Query(queryTableSQL)
	if err != nil {
		global.Log.Errorf("查询sDb所有表错误 %s", err)
		return nil
	}
	sDbTable := getRowsSimpleData(sDbTableRows)
	sDbTables := make([]string, len(bDbTable))
	for index := range sDbTable {
		sDbTables = append(sDbTables, sDbTable[index].(string))
	}

	dbNameMap := make(map[string]bool, len(sDbTables))
	for _, elem := range sDbTables {
		dbNameMap[elem] = true
	}

	unCreatedTableName := make([]string, 0)
	for _, elem := range bDbTables {
		if !dbNameMap[elem] {
			unCreatedTableName = append(unCreatedTableName, elem)
		}
	}

	return unCreatedTableName
}

func (r body) testConnection() bool {
	var err error

	if err = r.bDb.Ping(); err != nil {
		return false
	}

	if err = r.sDb.Ping(); err != nil {
		return false
	}

	return true
}

func getRowsData(rows *sql.Rows) []map[string]interface{} {
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
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
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var res = make([]interface{}, 0)
	for rows.Next() {
		var name interface{}
		_ = rows.Scan(&name)
		res = append(res, util.Byte2Str(name.([]byte)))
	}
	return res
}
