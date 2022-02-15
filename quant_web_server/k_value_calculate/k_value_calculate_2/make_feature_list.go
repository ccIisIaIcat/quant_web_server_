package main

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Feature_list struct {
	File_name             string //必填
	Username              string //必填
	Password              string //必填
	Cross_feature_data    string //必填
	cross_feature_data_db *sql.DB
	K_value_data          string //必填
	K_value_data_db       *sql.DB
	period_               int
	label                 string
	Max_int               int
	list_length           int
	dividied_point        float64
	data_long_list        []float64
}

func (F *Feature_list) init() {
	if F.period_ == 0 {
		F.period_ = 10
	}
	if F.dividied_point == 0 {
		F.dividied_point = 0.9
	}
	F.Max_int = F.period_
	dsn := F.Username + ":" + F.Password + "@tcp(127.0.0.1:3306)/" + F.Cross_feature_data
	dsn_2 := F.Username + ":" + F.Password + "@tcp(127.0.0.1:3306)/" + F.K_value_data
	var err error
	F.cross_feature_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = F.cross_feature_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	F.K_value_data_db, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = F.K_value_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	F.list_length = F.get_list_length() - 2*F.Max_int
}

type table_length struct {
	length int
}

func (F *Feature_list) get_list_length() int {
	sql := "select COUNT(*) from " + F.File_name + "goal_feature_p1;"
	data := F.cross_feature_data_db.QueryRow(sql)
	var t_l table_length
	data.Scan(&t_l.length)
	return t_l.length
}

func (F *Feature_list) creat_table() string {
	number := strconv.Itoa(F.period_)
	label := F.File_name + "goal_matrix_p" + number
	sql := "CREATE TABLE " + label + "(id int PRIMARY KEY AUTO_INCREMENT,T_start bigint,price_change_abs double,big_or_not int)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := F.K_value_data_db.Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println(label + "建表成功")
	}
	F.label = label
	return label
}

type data_change struct {
	T_start    int64
	data_start float64
	data_end   float64
}

func (d *data_change) price_change() float64 {
	return math.Abs(d.data_end - d.data_start)
}

func (F *Feature_list) save_goal_data_1() {
	sql_query := "select T_start,bids_start,bids_end from " + F.File_name + "goal_feature_p1 where id=?"
	stmt_query, err := F.cross_feature_data_db.Prepare(sql_query)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt_query.Close()
	sql_insert := "insert into " + F.label + "(T_start,price_change_abs) values (?,?)"
	stmt_insert, err := F.K_value_data_db.Prepare(sql_insert)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt_insert.Close()
	var d_c_1 data_change
	var d_c_2 data_change
	for id := F.Max_int + 1; id <= F.list_length-F.Max_int+1; id++ {
		fmt.Println("insert_price_change_abs:", id-F.Max_int-1, "/", F.list_length-2*F.Max_int)
		data_1 := stmt_query.QueryRow(id)
		data_2 := stmt_query.QueryRow(id + F.period_ - 1)
		data_1.Scan(&d_c_1.T_start, &d_c_1.data_start, &d_c_1.data_end)
		data_2.Scan(&d_c_2.T_start, &d_c_2.data_start, &d_c_2.data_end)
		d_c_1.data_end = d_c_2.data_end
		stmt_insert.Exec(d_c_1.T_start, d_c_1.price_change())
		F.data_long_list = append(F.data_long_list, d_c_1.price_change())
	}
}

type get_info struct {
	price_change_abs float64
}

func (F *Feature_list) update_big_or_not() {
	sort.Float64s(sort.Float64Slice(F.data_long_list))
	divide_point := F.data_long_list[int(math.Floor(float64(len(F.data_long_list))*F.dividied_point))]
	sql_query := "select price_change_abs from " + F.label + " where id=?"
	stmt_query, err := F.K_value_data_db.Prepare(sql_query)
	if err != nil {
		fmt.Println(err)
	}
	var g_i get_info
	for id := 1; id <= F.list_length-2*F.Max_int+1; id++ {
		fmt.Println("judge_big_or_not:", id-F.Max_int-1, "/", F.list_length-2*F.Max_int)
		data := stmt_query.QueryRow(id)
		data.Scan(&g_i.price_change_abs)
		fmt.Println(id-F.Max_int, g_i.price_change_abs, divide_point)
		if g_i.price_change_abs >= divide_point {
			F.K_value_data_db.Exec("update " + F.label + fmt.Sprintf(" set big_or_not=1 where id=%v", id))
		} else {
			F.K_value_data_db.Exec("update " + F.label + fmt.Sprintf(" set big_or_not=0 where id=%v", id))
		}

	}

}

func main() {
	f_l := Feature_list{File_name: "btcusdt2022_01_01_19h05m11s", Username: "root", Password: "", Cross_feature_data: "cross_feature_data", K_value_data: "classify_data"}
	f_l.init()
	f_l.creat_table()
	f_l.save_goal_data_1()
	f_l.update_big_or_not()

}
