package read_mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type My_mysql struct {
	Username          string
	Password          string
	Original_data     string
	Processed_data    string
	original_data_db  *sql.DB
	processed_data_db *sql.DB
}

//mysql初始化
func (M *My_mysql) Init() {
	dsn := M.Username + ":" + M.Password + "@tcp(127.0.0.1:3306)/" + M.Original_data
	dsn_2 := M.Username + ":" + M.Password + "@tcp(127.0.0.1:3306)/" + M.Processed_data
	var err error
	//原始数据数据库初始化
	M.original_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = M.original_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	//预处理数据数据库初始化
	M.processed_data_db, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = M.processed_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}

}

//查看original_data内的全部列表信息
func (M *My_mysql) Show_original_data_table_info() []string {
	sql := "show tables;"
	sql_info, _ := M.original_data_db.Query(sql)
	var table string
	reply := make([]string, 0)
	for sql_info.Next() {
		sql_info.Scan(&table)
		reply = append(reply, table)
	}
	return reply
}

func (M *My_mysql) Show_processed_data_table_info() []string {
	sql := "show tables;"
	sql_info, _ := M.processed_data_db.Query(sql)
	var table string
	reply := make([]string, 0)
	for sql_info.Next() {
		sql_info.Scan(&table)
		reply = append(reply, table)
	}
	return reply
}

//关闭数据库
func (M *My_mysql) Close_mysql() {
	M.original_data_db.Close()
	M.processed_data_db.Close()
}

// func main() {
// 	m_m := My_mysql{Username: "root", Password: "", Original_data: "original_data", Processed_data: "processed_data"}
// 	m_m.Init()
// 	fmt.Println(m_m.Show_processed_data_table_info())
// }
