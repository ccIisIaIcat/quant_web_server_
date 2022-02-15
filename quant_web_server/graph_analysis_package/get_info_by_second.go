package graph_analysis_package

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Price_by_second struct {
	File_name          string    //所收集数据的文件名
	Username           string    //数据库用户名
	Password           string    //密码
	Cross_feature_data string    //交叉数据数据库
	db                 *sql.DB   //读取文件的sql对象
	answer_list        []float64 //用于存储读到的数据
}

func (P *Price_by_second) init() {
	dsn := P.Username + ":" + P.Password + "@tcp(127.0.0.1:3306)/" + P.Cross_feature_data
	var err error
	P.db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = P.db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	fmt.Println("db连接成功！")
}

type bids_start struct {
	Price float64 `json:"Price"`
}

func (P *Price_by_second) select_the_price() {
	select_sql := "select " + "bids_start" + " from " + P.File_name + "goal_feature_p1"
	data, err := P.db.Query(select_sql)
	if err != nil {
		fmt.Println("数据读取错误！")
	}
	var b_s bids_start
	for data.Next() {
		data.Scan(&b_s.Price)
		P.answer_list = append(P.answer_list, b_s.Price)
	}

}

func (P *Price_by_second) send_message() []byte {
	new, err := json.Marshal(P.answer_list)
	if err != nil {
		fmt.Println(err)
	}
	return new
}

func (P *Price_by_second) Start() []byte {
	P.init()
	P.select_the_price()
	return P.send_message()
}
