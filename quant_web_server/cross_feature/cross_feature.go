package cross_feature

import (
	"database/sql"
	"fmt"
	"math"
	"quant_web_server/feature_extraction"
	"quant_web_server/goal_feature"
	"sort"

	_ "github.com/go-sql-driver/mysql"
)

type Cross_feature struct {
	File_name             string
	Username              string
	Password              string
	Processed_data        string
	Cross_feature_data    string
	cross_feature_data_db *sql.DB
	Divided_point         float64
	price_change_list     []float64
	price_change_abs_list []float64
}

func (C *Cross_feature) init() {
	dsn := C.Username + ":" + C.Password + "@tcp(127.0.0.1:3306)/" + C.Cross_feature_data
	var err error
	C.cross_feature_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
}

func (C *Cross_feature) Make_all_list() {
	feature_list := []int{1}
	goal_list := []int{1}
	for i := 0; i < len(feature_list); i++ {
		feature_time_period := feature_list[i]
		f_e := feature_extraction.Feature_extraction{File_name: C.File_name, Username: C.Username, Password: C.Password, Processed_data: C.Processed_data, Cross_feature_data: C.Cross_feature_data, Time_period: feature_time_period}
		f_e.Get_all_info()
	}
	for j := 0; j < len(goal_list); j++ {
		goal_feature_period := goal_list[j]
		g_e := goal_feature.Goal_feature{File_name: C.File_name, Username: C.Username, Password: C.Password, Processed_data: C.Processed_data, Cross_feature_data: C.Cross_feature_data, Goal_feature_time_period: goal_feature_period}
		g_e.Get_all_info()
	}
	fmt.Println("存入目标变量")
	C.Get_price_change()
	fmt.Println("存入目标变量完成")
}

func (C *Cross_feature) create_new_table() {
	create_sql := "create table " + C.File_name + "big_or_not_p1 " + "(id int PRIMARY KEY AUTO_INCREMENT,price_change double,big_or_not int)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	C.cross_feature_data_db.Exec(create_sql)
}

type price_unit struct {
	start_price float64
	end_price   float64
}

func (C *Cross_feature) Get_price_change() {
	C.init()
	C.create_new_table()
	sql_select := "select bids_start,bids_end from " + C.File_name + "goal_feature_p1"
	data, err := C.cross_feature_data_db.Query(sql_select)
	if err != nil {
		fmt.Println("ERR", err)
	}
	var p_u price_unit
	for data.Next() {
		data.Scan(&p_u.start_price, &p_u.end_price)
		C.price_change_list = append(C.price_change_list, p_u.end_price-p_u.start_price)
		C.price_change_abs_list = append(C.price_change_abs_list, math.Abs(p_u.end_price-p_u.start_price))
	}
	sort.Float64s(sort.Float64Slice(C.price_change_abs_list))
	judge_point := C.price_change_abs_list[int(math.Floor(float64(len(C.price_change_abs_list))*C.Divided_point))]
	sql_insert := "insert into " + C.File_name + "big_or_not_p1 (price_change,big_or_not) values(?,?)"
	stmt_insert, err := C.cross_feature_data_db.Prepare(sql_insert)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	for i := 0; i < len(C.price_change_abs_list); i++ {
		// fmt.Println(i)
		if math.Abs(C.price_change_list[i]) > judge_point {
			if C.price_change_list[i] > 0 {
				stmt_insert.Exec(C.price_change_list[i], 1)
			} else {
				stmt_insert.Exec(C.price_change_list[i], -1)
			}
		} else {
			stmt_insert.Exec(C.price_change_list[i], 0)
		}

	}

}
