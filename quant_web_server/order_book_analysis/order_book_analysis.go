package order_book_analysis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Order_book_analysis struct {
	File_name               string         //所收集数据的文件名
	Username                string         //数据库用户名
	Password                string         //密码
	Original_data           string         //原始数据数据库
	Processed_data          string         //处理数据数据库
	db_1                    *sql.DB        //读取文件的sql对象
	db_2                    *sql.DB        //存入数据的sql对象
	label_name              string         //用于存放对象名
	channel_read_total      chan []byte    //用于读取非订单簿全部类型数据的chan
	channel_read_mark_price chan []byte    //用于从channel_read_total中分流出最新标记价格的信息
	channel_read_best_offer chan []byte    //用于从channel_read_total中分流出最优订单信息
	channel_read_agg_trade  chan []byte    //用于从channel_read_total中分流出聚合交易信息
	wg                      sync.WaitGroup //用于优雅地停止进程
	max_id                  int            //用于输出信息
	K_                      float64        //分位点
}

type recevier struct {
	id       int64
	New_info []byte
}

type specific_info struct {
	Time int               `json:"T"`
	Bids map[string]string `json:"bids"`
	Asks map[string]string `json:"asks"`
}

type answer struct {
	Time              int               `json:"Time"`
	Abs_ratio         float64           `json:"abs_ratio"`
	Volum_ratio       float64           `json:"volum_ratio"`
	Average_ratio     float64           `json:"average_ratio"`
	Varance_ratio     float64           `json:"varance_ratio"`
	Hard_point_list_b map[string]string `json:"hard_point_list_b"`
	Hard_point_list_a map[string]string `json:"hard_point_list_a"`
}

func turn_string_map_into_float_map(string_map map[string]string) map[float64]float64 {
	float_map := make(map[float64]float64, 0)
	for k, v := range string_map {
		price, _ := strconv.ParseFloat(k, 64)
		volume, _ := strconv.ParseFloat(v, 64)
		float_map[price] = volume
	}
	return float_map
}

//float32s2 := strconv.FormatFloat(v, 'E', -1, 64)
func turn_float_map_into_string_map(float_map map[float64]float64) map[string]string {
	string_map := make(map[string]string, 0)
	for k, v := range float_map {
		price := strconv.FormatFloat(k, 'E', -1, 64)
		proportion := strconv.FormatFloat(v, 'E', -1, 64)
		string_map[price] = proportion
	}
	return string_map
}

func (O *Order_book_analysis) init() {
	if O.K_ == 0 {
		O.K_ = 0.05
	}
	dsn := O.Username + ":" + O.Password + "@tcp(127.0.0.1:3306)/" + O.Original_data
	dsn_2 := O.Username + ":" + O.Password + "@tcp(127.0.0.1:3306)/" + O.Processed_data
	var err error
	O.db_1, err = sql.Open("mysql", dsn)
	O.db_2, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db_1格式错误：", err)
		return
	}
	O.db_2, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db_2格式错误：", err)
		return
	}
	err = O.db_1.Ping()
	if err != nil {
		fmt.Println("db_1建立链接出错：")
		panic(err)
	}
	fmt.Println("db_1连接成功！")
	err = O.db_2.Ping()
	if err != nil {
		fmt.Println("db_2建立链接出错：")
		panic(err)
	}
	fmt.Println("db_2连接成功！")
	O.label_name = O.creat_table()

}

func (O *Order_book_analysis) calculate_feature(sample map[float64]float64) (float64, float64, float64, float64, float64, float64, map[float64]float64) {
	max_, min_ := float64(0), float64(9999999999)
	sum_1 := float64(0)
	sum_2 := float64(0)
	sum_3 := float64(0)
	sum_tool := float64(0)
	for k, v := range sample {
		if k > max_ {
			max_ = k
		}
		if k < min_ {
			min_ = k
		}
		sum_1 += k
		sum_tool += k * v
		sum_2 += k * v * k * v
		sum_3 += v
	}
	num_sum := sum_3
	average := sum_1 / float64(len(sample))
	varance := (sum_2 - sum_tool*sum_tool/float64(len(sample))) / float64(len(sample))
	length := max_ - min_
	big_orderbook := make(map[float64]float64, 0)
	for k, v := range sample {
		if v > num_sum*O.K_ {
			big_orderbook[k] = v / num_sum
		}
	}
	return min_, max_, length, num_sum, average, varance, big_orderbook

}

func (O *Order_book_analysis) creat_table() string {
	label_1 := O.File_name + "_analysis"
	sql_1 := "CREATE TABLE " + label_1 + "(id int PRIMARY KEY AUTO_INCREMENT,Time bigint,abs_ratio double,volum_ratio double,average_ratio double,varance_ratio double,hard_point_list_b blob,hard_point_list_a blob)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := O.db_2.Prepare(sql_1)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println(label_1 + "建表成功")
	}
	return label_1
}

func (O *Order_book_analysis) insert_data(answer_ *answer) {
	sql_insert := "insert into " + O.label_name + " (Time,abs_ratio,volum_ratio,average_ratio,varance_ratio,hard_point_list_b,hard_point_list_a) values(?,?,?,?,?,?,?);"
	stmt_insert, err := O.db_2.Prepare(sql_insert)
	if err != nil {
		fmt.Println("stmt_insert错误：", err)
	}
	bb, _ := json.Marshal(answer_.Hard_point_list_b)
	aa, _ := json.Marshal(answer_.Hard_point_list_a)
	stmt_insert.Exec(answer_.Time, answer_.Abs_ratio, answer_.Volum_ratio, answer_.Average_ratio, answer_.Varance_ratio, bb, aa)
	defer stmt_insert.Close()
}

func (O *Order_book_analysis) Get_order_book_info_and_save() {
	O.init()
	sql := "select id,context from " + O.File_name + " where id =?"
	stmt_order_book, err := O.db_1.Prepare(sql)
	if err != nil {
		fmt.Println("stmt_order_book错误：", err)
	}
	defer stmt_order_book.Close()
	id := 1
	for {
		data := stmt_order_book.QueryRow(id)
		var re recevier
		data.Scan(&re.id, &re.New_info)
		if re.id == 0 {
			break
		}
		id++
		var s_i specific_info
		json.Unmarshal(re.New_info, &s_i)
		T_, B_, A_ := s_i.Time, s_i.Bids, s_i.Asks
		b_, a_ := turn_string_map_into_float_map(B_), turn_string_map_into_float_map(A_)
		_, b_max, b_1, b_2, b_3, b_4, b_list := O.calculate_feature(b_)
		a_min, _, a_1, a_2, a_3, a_4, a_list := O.calculate_feature(a_)
		var anan answer
		anan.Abs_ratio = b_1 / a_1
		anan.Time = T_
		anan.Volum_ratio = b_2 / a_2
		anan.Average_ratio = (b_max - b_3) / (a_3 - a_min)
		anan.Varance_ratio = b_4 / a_4
		anan.Hard_point_list_b = turn_float_map_into_string_map(b_list)
		anan.Hard_point_list_a = turn_float_map_into_string_map(a_list)
		O.insert_data(&anan)
	}
	fmt.Println("数据录入完成")

}

// func main() {
// 	fmt.Println("hello world")
// 	o_b_a := Order_book_analysis{K_: 0.05, File_name: "lunausdt2021_11_11_18h03m11s_local_order_book", Username: "root", Password: "", Original_data: "original_data", Processed_data: "processed_data"}
// 	o_b_a.get_order_book_info_and_save()

// }
