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

}

func (r body) syncTableStruct(tables []string) bool {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Errorf("同步表结构发生 %s", err)
		}
	}()

	if len(tables) == 0 || tables == nil {
		return true
	}
	queryTableStructSql := "show create table %s"
	rowDateSQl := "alter table %s add column row_date varchar(10) comment '备份时间'"
	for _, elem := range tables {
		rows, err := r.bDb.Query(fmt.Sprintf(queryTableStructSql, elem))
		if err != nil {
			global.Log.Errorf("获取表结构异常 %s", err)
			return false
		}
		// 获取到原生表创建SQL
		tableSQL := getRowsData(rows)[0]["Create Table"]
		// 原生SQL分割()
		tableStructCreateSQL := tableSQL.(string)
		begin := strings.Index(tableStructCreateSQL, "(")
		end := strings.LastIndex(tableStructCreateSQL, ")")
		beforeSql, centSQL, endSQL := tableStructCreateSQL[:begin+1], tableStructCreateSQL[begin+1:end], tableStructCreateSQL[end:]
		//global.Log.Infof("%s 原生表结构SQL: %s", elem, tableStructCreateSQL)

		sqlArray := strings.Split(centSQL, ",")
		newSqlArray := make([]string, 0)
		for _, elem := range sqlArray {
			if strings.Contains(elem, "PRIMARY KEY") ||
				strings.Contains(elem, "UNIQUE KEY") {
				continue
			}
			elem = strings.ReplaceAll(elem, "AUTO_INCREMENT", "")
			newSqlArray = append(newSqlArray, elem)
		}
		newCentSQl := strings.Join(newSqlArray, ",")

		endSQL = strings.ReplaceAll(endSQL, "ENGINE=MyISAM", "")
		endSQL = strings.ReplaceAll(endSQL, "ENGINE=InnoDB", "")
		endSQL = strings.ReplaceAll(endSQL, "ROW_FORMAT=FIXED", "")
		reg := regexp.MustCompile(`AUTO_INCREMENT=(\d){1,}`)
		endSQL = reg.ReplaceAllString(endSQL, ``)
		reg = regexp.MustCompile(`ROW_FORMAT=[A-Z]{1,} `)
		endSQL = reg.ReplaceAllString(endSQL, ``)

		newTableStructSQl := fmt.Sprintf("%s%s%s", beforeSql, newCentSQl, endSQL)
		//global.Log.Infof("%s 新结构SQL:\n%s", elem, newTableStructSQl)
		_, err = r.sDb.Exec(newTableStructSQl)
		if err != nil {
			global.Log.Errorf("%s 新SQL执行错误\n%s\n%s", elem, newTableStructSQl, err)
			return false
		}
		// 增加 row_date
		_, err = r.bDb.Exec(fmt.Sprintf(rowDateSQl, elem))
		if err != nil {
			global.Log.Errorf("%s 追加row_date错误: %s", elem, err)
			return false
		}
	}

	return true
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
