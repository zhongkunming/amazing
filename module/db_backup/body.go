package db_backup

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"service-hub/global"
	"unsafe"
)

type body struct {
	db1 *sql.DB
	db2 *sql.DB
}

func (r body) run() {
	marshal, _ := json.Marshal(global.Global)
	global.Log.Infof("%s", Byte2Str(marshal))
	println(Byte2Str(marshal))
	r.loadDb()
	//dsn := "root:LANKE678@tcp(114.116.112.75:3306)/erp?charset=utf8mb4&parseTime=True&loc=Local"
	//b := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	//bDb, _ := gorm.Open(mysql.Open(b), &gorm.Config{})
	//
	//b2 := "root:123456@tcp(127.0.0.1:3306)/test1?charset=utf8mb4&parseTime=True&loc=Local"
	//b2Db, _ := gorm.Open(mysql.Open(b2), &gorm.Config{})
	//
	//// 查询所有表名
	//tx := bDb.Raw("show tables")
	////var tableName string
	////tx.Scan(&tableName)
	//rows, _ := tx.Rows()
	//var tableNames = make([]string, 0)
	//_ = tx.ScanRows(rows, &tableNames)
	//tableNames = tableNames[1:]
	//
	//// 查询2表所有表
	//tx2 := b2Db.Raw("show tables")
	////var tableName string
	////tx.Scan(&tableName)
	//rows2, _ := tx2.Rows()
	//var tableNames2 = make([]string, 0)
	//_ = tx.ScanRows(rows2, &tableNames2)
	//tableNames2 = tableNames2[1:]
	//
	//println()
	//
	//db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	////rows, _ := db.Query("show tables")
	//rows, _ := db.Query("show variables")
	//
	//types, _ := rows.ColumnTypes()
	//var rowParam = make([]interface{}, len(types))
	//var rowValue = make([]interface{}, len(types))
	//for i, colType := range types {
	//	rowValue[i] = reflect.New(colType.ScanType())
	//	rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface()
	//}
	//res := make([]map[string]interface{}, 0)
	//for rows.Next() {
	//	_ = rows.Scan(rowParam...)
	//	record := make(map[string]interface{})
	//	for i, colType := range types {
	//		if rowValue[i] == nil {
	//			// 如果是 nil 用空字符串代替
	//			record[colType.Name()] = ""
	//		} else {
	//			//record[colType.Name()] = rowValue[i]
	//			record[colType.Name()] = Byte2Str(rowValue[i].([]byte))
	//		}
	//	}
	//	res = append(res, record)
	//}
	//marshal, _ := json.Marshal(res)
	//fmt.Println(Byte2Str(marshal))

	// 必须要把 rows 里的内容读完，或者显式调用 Close() 方法，
	// 否则在 defer 的 rows.Close() 执行之前，连接永远不会释放
	//var tablesName = make([]string, 0)
	//for rows.Next() {
	//	var name string
	//	_ = rows.Scan(&name)
	//	tablesName = append(tablesName, name)
	//}
}

func (r body) loadDb() {

	//var err error
	//db1Config := global.Global.DbBackup.Db1
	//global.Log.Infof("数据库1配置, %v", db1Config)
	//r.db1, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
	//	db1Config.Username, db1Config.Passwd, db1Config.Host, db1Config.Database))
	//if err != nil {
	//	global.Log.Errorf("连接数据库1失败, %s", err)
	//	return
	//}
	//
	//db2Config := global.Global.DbBackup.Db2
	//global.Log.Infof("数据库2配置, %v", db1Config)
	//r.db2, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
	//	db2Config.Username, db2Config.Passwd, db2Config.Host, db2Config.Database))
	//if err != nil {
	//	global.Log.Errorf("连接数据库2失败, %s", err)
	//	return
	//}
}

func qmap(db *gorm.DB) {
	tx := db.Raw("show variables")

	rows, _ := tx.Rows()
	res := make([]map[string]interface{}, 0)

	types, _ := rows.ColumnTypes()
	var rowParam = make([]interface{}, len(types))
	var rowValue = make([]interface{}, len(types))
	for i, colType := range types {
		rowValue[i] = reflect.New(colType.ScanType())
		rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface()
	}

	for rows.Next() {
		_ = rows.Scan(rowParam...)
		record := make(map[string]interface{})
		for i, colType := range types {
			if rowValue[i] == nil {
				// 如果是 nil 用空字符串代替
				record[colType.Name()] = ""
			} else {
				//record[colType.Name()] = rowValue[i]
				record[colType.Name()] = Byte2Str(rowValue[i].([]byte))
			}
		}
		res = append(res, record)
	}
	marshal, _ := json.Marshal(res)
	fmt.Println(Byte2Str(marshal))
}

func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
