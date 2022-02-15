package graph_analysis_package

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Price_detial struct {
	File_name          string    //所收集数据的文件名
	Username           string    //数据库用户名
	Password           string    //密码
	Cross_feature_data string    //交叉数据数据库
	cross_feature_db   *sql.DB   //读取交叉数据的sql对象
	Processed_data     string    //预处理数据数据库
	processed_data_db  *sql.DB   //读取与处理数据的sql对象
	answer_list        []float64 //用于存储读到的数据
	Time_point_id      int       //数据点id
	Time_point_period  int       //图像展示时长单边范围
	time_point_time    int64     //时间点

}

func (P *Price_detial) init() {
	dsn_cross_feature := P.Username + ":" + P.Password + "@tcp(127.0.0.1:3306)/" + P.Cross_feature_data
	dsn_processed_data := P.Username + ":" + P.Password + "@tcp(127.0.0.1:3306)/" + P.Processed_data
	var err error
	P.processed_data_db, err = sql.Open("mysql", dsn_processed_data)
	if err != nil {
		fmt.Println("processed_data_db格式错误：", err)
		return
	}
	err = P.processed_data_db.Ping()
	if err != nil {
		fmt.Println("processed_data_db建立链接出错：")
		panic(err)
	}
	fmt.Println("processed_data_db连接成功！")

	P.cross_feature_db, err = sql.Open("mysql", dsn_cross_feature)
	if err != nil {
		fmt.Println("cross_feature_db格式错误：", err)
		return
	}
	err = P.cross_feature_db.Ping()
	if err != nil {
		fmt.Println("cross_feature_db建立链接出错：")
		panic(err)
	}
	fmt.Println("cross_feature_db连接成功！")
}

type time_point struct {
	time_id int
	time    int64
}

func (P *Price_detial) get_time_point_time() {
	sql_select := "select T_start from " + P.File_name + "feature_extraction_p1" + " where id=" + strconv.Itoa(P.Time_point_id)
	data := P.cross_feature_db.QueryRow(sql_select)
	var t_p time_point
	data.Scan(&t_p.time)
	P.time_point_time = t_p.time

}

type get_info struct {
	T1 int64
	Bp float64
	Bq float64
	Ap float64
	Aq float64
}

func (P *Price_detial) get_info_list() []byte {
	x_data := make([]float64, 0)
	y_price_b := make([]float64, 0)
	y_quantity_b := make([]float64, 0)
	y_price_a := make([]float64, 0)
	y_quantity_a := make([]float64, 0)
	sql_select := "select T1,Bp,Bq,Ap,Aq from " + P.File_name + "best_offer_price where T1>" + strconv.Itoa(int(P.time_point_time)-P.Time_point_period*1000) + "&&T1<" + strconv.Itoa(int(P.time_point_time)+P.Time_point_period*1000)
	data, err := P.processed_data_db.Query(sql_select)
	if err != nil {
		fmt.Println("查询出错：", err)
	}
	var g_f get_info
	for data.Next() {
		data.Scan(&g_f.T1, &g_f.Bp, &g_f.Bq, &g_f.Ap, &g_f.Aq)
		x_data = append(x_data, float64((g_f.T1-P.time_point_time))/1000+float64(P.Time_point_id))
		y_price_b = append(y_price_b, g_f.Bp)
		y_price_a = append(y_price_a, g_f.Ap)
		y_quantity_b = append(y_quantity_b, g_f.Bq)
		y_quantity_a = append(y_quantity_a, g_f.Aq)

	}
	i_s := info_sender{T: x_data, Bp: y_price_b, Bq: y_quantity_b, Ap: y_price_a, Aq: y_quantity_a}
	my_info, err := json.Marshal(i_s)
	return my_info

}

type info_sender struct {
	T  []float64 `json:"T"`
	Bp []float64 `json:"Bp"`
	Bq []float64 `json:"Bq"`
	Ap []float64 `json:"Ap"`
	Aq []float64 `json:"Aq"`
}

func (P *Price_detial) Get_info_list() []byte {
	P.init()
	P.get_time_point_time()
	// fmt.Println(P.get_info_list())
	return P.get_info_list()

}
