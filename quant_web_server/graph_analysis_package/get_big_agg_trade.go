package graph_analysis_package

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Big_agg_trade struct {
	File_name          string  //所收集数据的文件名
	Username           string  //数据库用户名
	Password           string  //密码
	Cross_feature_data string  //交叉数据数据库
	db                 *sql.DB //读取文件的sql对象
	answer_list        []int   //用于存储读到的数据
}

func (B *Big_agg_trade) init() {
	dsn := B.Username + ":" + B.Password + "@tcp(127.0.0.1:3306)/" + B.Cross_feature_data
	var err error
	B.db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = B.db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	fmt.Println("db连接成功！")
	B.answer_list = make([]int, 0)
}

type info_check struct {
	big_or_not int
}

func (B *Big_agg_trade) select_info_and_process() {
	sql_select := "select big_or_not from " + B.File_name + "big_or_not_p1"
	data, err := B.db.Query(sql_select)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	var i_f info_check
	for data.Next() {
		data.Scan(&i_f.big_or_not)
		B.answer_list = append(B.answer_list, i_f.big_or_not)
	}
}

func (B *Big_agg_trade) send_message() []byte {
	new, err := json.Marshal(B.answer_list)
	if err != nil {
		fmt.Println(err)
	}
	return new
}

func (B *Big_agg_trade) Start() []byte {
	B.init()
	B.select_info_and_process()
	return B.send_message()
}
