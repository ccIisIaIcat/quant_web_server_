package graph_analysis_package

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Order_book struct {
	File_name             string //必填
	Username              string //必填
	Password              string //必填
	Original_data         string //必填
	original_data_db      *sql.DB
	Cross_feature_data    string //必填
	cross_feature_data_db *sql.DB
	Processed_data        string //必填
	processed_data_db     *sql.DB
	Time_start_id         int //起始点的时间点(必填)
	time_start            int64
	info_send_sub         info_sender_new
	time_start_new_id     int
}

type info_sender_new struct {
	time_start          int64
	Id_time             float64   `json:"Id_time"`
	Order_book_id       int       `json:"Order_book_id"`
	Abs                 float64   `json:"Abs"`
	Volum               float64   `json:"Volum"`
	Average             float64   `json:"Average"`
	Varance             float64   `json:"Varance"`
	Bids_price_list     []float64 `json:"Bids_price_list"`
	Bids_quantity_list  []float64 `json:"Bids_quantity_list"`
	Bids_info_list_pre  []int     `json:"Bids_info_list_pre"`
	Bids_info_list_next []int     `json:"Bids_info_list_next"`
	Asks_price_list     []float64 `json:"Asks_price_list"`
	Asks_quantity_list  []float64 `json:"Asks_quantity_list"`
	Asks_info_list_pre  []int     `json:"Asks_info_list_pre"`
	Asks_info_list_next []int     `json:"Asks_info_list_next"`
}

func (O *Order_book) init() {
	dsn_original := O.Username + ":" + O.Password + "@tcp(127.0.0.1:3306)/" + O.Original_data
	dsn_cross_feature := O.Username + ":" + O.Password + "@tcp(127.0.0.1:3306)/" + O.Cross_feature_data
	dsn_processed_data := O.Username + ":" + O.Password + "@tcp(127.0.0.1:3306)/" + O.Processed_data
	var err error
	O.original_data_db, err = sql.Open("mysql", dsn_original)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = O.original_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	O.cross_feature_data_db, err = sql.Open("mysql", dsn_cross_feature)
	if err != nil {
		fmt.Println("db格式错误：", err)
	}
	err = O.cross_feature_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	O.processed_data_db, err = sql.Open("mysql", dsn_processed_data)
	if err != nil {
		fmt.Println("db格式错误：", err)
	}
	err = O.processed_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}

}

func (O *Order_book) get_info_from_cross_feature_data() {
	sql_select := "select T_start,orderbook_abs_ratio,orderbook_volum_ratio,orderbook_average_ratio,orderbook_varance_ratio from " + O.File_name + "feature_extraction_p1 where id=?"
	data := O.cross_feature_data_db.QueryRow(sql_select, O.Time_start_id)
	data.Scan(&O.info_send_sub.time_start, &O.info_send_sub.Abs, &O.info_send_sub.Volum, &O.info_send_sub.Average, &O.info_send_sub.Varance)
	O.time_start = int64(O.info_send_sub.time_start)
}

type info_temp struct {
	id           int
	time         int64
	context_temp []byte
}

type specific_info struct {
	Time int               `json:"T"`
	Bids map[string]string `json:"bids"`
	Asks map[string]string `json:"asks"`
}

func (O *Order_book) get_info_from_original_data() {
	sql_select := "select id,Time from " + O.File_name + "_local_order_book_analysis "
	data, err := O.processed_data_db.Query(sql_select)
	if err != nil {
		fmt.Println("ERR!__", err)
	}
	var i_t info_temp
	for data.Next() {
		data.Scan(&i_t.id, &i_t.time)
		if i_t.time > O.time_start {
			O.time_start_new_id = i_t.id
			O.info_send_sub.Id_time = float64(O.Time_start_id) + float64((int(i_t.time)-int(O.time_start)))/1000
			break
		}
	}
	sql_select_2 := "select context from " + O.File_name + "_local_order_book where id=?"
	context_data := O.original_data_db.QueryRow(sql_select_2, O.time_start_new_id)
	context_data.Scan(&i_t.context_temp)
	var s_i specific_info
	json.Unmarshal(i_t.context_temp, &s_i)
	bids_string_map := s_i.Bids
	asks_string_map := s_i.Asks
	O.info_send_sub.Bids_price_list = O.return_float_list("b", bids_string_map)
	O.info_send_sub.Asks_price_list = O.return_float_list("a", asks_string_map)
	O.info_send_sub.Bids_quantity_list = make([]float64, 11)
	O.info_send_sub.Asks_quantity_list = make([]float64, 11)
	for i := 0; i < len(O.info_send_sub.Bids_price_list); i++ {
		str_b := strconv.FormatFloat(O.info_send_sub.Bids_price_list[i], 'f', -1, 64)
		float_b, _ := strconv.ParseFloat(bids_string_map[str_b], 64)

		O.info_send_sub.Bids_quantity_list[i] = float_b
		str_a := strconv.FormatFloat(O.info_send_sub.Asks_price_list[i], 'f', -1, 64)
		float_a, _ := strconv.ParseFloat(asks_string_map[str_a], 64)
		O.info_send_sub.Asks_quantity_list[i] = float_a
	}
	O.compare_pre_and_next(O.get_info_from_original_data_sub(O.time_start_new_id-1), O.get_info_from_original_data_sub(O.time_start_new_id), O.get_info_from_original_data_sub(O.time_start_new_id+1))

}

func (O *Order_book) compare_pre_and_next(pre_orderbook tool_orderbook, now_orderbook tool_orderbook, next_orderbook tool_orderbook) {
	// pre
	pre_bids_price_list := pre_orderbook.Bids_price_list
	now_bids_price_list := now_orderbook.Bids_price_list
	pre_bids_quantity_list := pre_orderbook.Bids_quantity_list
	now_bids_quantity_list := now_orderbook.Bids_quantity_list
	O.info_send_sub.Bids_info_list_pre = make([]int, len(pre_bids_price_list))
	pre_asks_price_list := pre_orderbook.Asks_price_list
	now_asks_price_list := now_orderbook.Asks_price_list
	pre_asks_quantity_list := pre_orderbook.Asks_quantity_list
	now_asks_quantity_list := now_orderbook.Asks_quantity_list
	O.info_send_sub.Asks_info_list_pre = make([]int, len(pre_asks_price_list))
	signal := 0
	// 处理bids_pre
	j := 0
	for i := 0; i < len(now_bids_price_list); i++ {
		if pre_bids_price_list[j] == now_bids_price_list[i] {
			if pre_bids_quantity_list[j] > now_bids_quantity_list[i] {
				// -2表示同一价格的订单数量减少
				O.info_send_sub.Bids_info_list_pre[i] = -2
			} else if pre_bids_quantity_list[j] < now_bids_quantity_list[i] {
				// -3表示同一价格的订单数量增加
				O.info_send_sub.Bids_info_list_pre[i] = -3
			}
			if j < 10 {
				j++
			}
		} else {
			signal = 0
			for num := j; num < len(now_bids_price_list); num++ {
				if pre_bids_price_list[num] == now_bids_price_list[i] {
					j = num
					i--
					signal = 1
					break
				}
			}
			if signal == 0 {
				// -1表示同一价格的订单数量增加
				O.info_send_sub.Bids_info_list_pre[i] = -1
			}
		}
	}
	// 处理asks_pre
	j = 0
	for i := 0; i < len(now_asks_price_list); i++ {
		if pre_asks_price_list[j] == now_asks_price_list[i] {
			if pre_asks_quantity_list[j] > now_asks_quantity_list[i] {
				// -2表示同一价格的订单数量减少
				O.info_send_sub.Asks_info_list_pre[i] = -2
			} else if pre_asks_quantity_list[j] < now_asks_quantity_list[i] {
				// -3表示同一价格的订单数量增加
				O.info_send_sub.Asks_info_list_pre[i] = -3
			}
			if j < 10 {
				j++
			}
		} else {
			signal = 0
			for num := j; num < len(now_asks_price_list); num++ {
				if pre_asks_price_list[num] == now_asks_price_list[i] {
					j = num
					i--
					signal = 1
					break
				}
			}
			if signal == 0 {
				// -1表示同一价格的订单将增加
				O.info_send_sub.Asks_info_list_pre[i] = -1
			}
		}
	}

	// next
	next_bids_price_list := next_orderbook.Bids_price_list
	next_bids_quantity_list := next_orderbook.Bids_quantity_list
	O.info_send_sub.Bids_info_list_next = make([]int, len(next_bids_price_list))
	next_asks_price_list := next_orderbook.Asks_price_list
	next_asks_quantity_list := next_orderbook.Asks_quantity_list
	O.info_send_sub.Asks_info_list_next = make([]int, len(next_asks_price_list))
	// 处理bids_next
	j = 0
	for i := 0; i < len(now_bids_price_list); i++ {
		if next_bids_price_list[j] == now_bids_price_list[i] {
			if next_bids_quantity_list[j] > now_bids_quantity_list[i] {
				// 2表示同一价格的订单下一节点将减少
				O.info_send_sub.Bids_info_list_next[i] = 2
			} else if next_bids_quantity_list[j] < now_bids_quantity_list[i] {
				// 3表示同一价格的订单下一节点将增加
				O.info_send_sub.Bids_info_list_next[i] = 3
			}
			if j < 10 {
				j++
			}
		} else {
			signal = 0
			for num := j; num < len(now_bids_price_list); num++ {
				if next_bids_price_list[num] == now_bids_price_list[i] {
					j = num
					i--
					signal = 1
					break
				}
			}
			if signal == 0 {
				// 1表示同一价格的订单下一节点将消失
				O.info_send_sub.Bids_info_list_next[i] = 1
			}
		}
	}

	j = 0
	for i := 0; i < len(now_asks_price_list); i++ {
		if next_asks_price_list[j] == now_asks_price_list[i] {
			if next_asks_quantity_list[j] > now_asks_quantity_list[i] {
				// 2表示同一价格的订单数量将减少
				O.info_send_sub.Asks_info_list_next[i] = 2
			} else if next_asks_quantity_list[j] < now_asks_quantity_list[i] {
				// 3表示同一价格的订单下一节点将增加
				O.info_send_sub.Asks_info_list_next[i] = 3
			}
			if j < 10 {
				j++
			}
		} else {
			signal = 0
			for num := j; num < len(now_asks_price_list); num++ {
				if next_asks_price_list[num] == now_asks_price_list[i] {
					j = num
					i--
					signal = 1
					break
				}
			}
			if signal == 0 {
				// 1表示同一价格的订单下一节点将消失
				O.info_send_sub.Asks_info_list_next[i] = 1
			}
		}
	}

}

type tool_orderbook struct {
	Bids_price_list    []float64
	Bids_quantity_list []float64
	Asks_price_list    []float64
	Asks_quantity_list []float64
}

func (O *Order_book) get_info_from_original_data_sub(original_id int) tool_orderbook {
	var answer_struct tool_orderbook
	sql_select := "select context from " + O.File_name + "_local_order_book where id=?"
	context_data := O.original_data_db.QueryRow(sql_select, original_id)
	var i_t info_temp
	context_data.Scan(&i_t.context_temp)
	var s_i specific_info
	json.Unmarshal(i_t.context_temp, &s_i)
	bids_string_map := s_i.Bids
	asks_string_map := s_i.Asks
	answer_struct.Bids_price_list = O.return_float_list("b", bids_string_map)
	answer_struct.Asks_price_list = O.return_float_list("a", asks_string_map)
	answer_struct.Bids_quantity_list = make([]float64, 11)
	answer_struct.Asks_quantity_list = make([]float64, 11)
	for i := 0; i < len(O.info_send_sub.Bids_price_list); i++ {
		str_b := strconv.FormatFloat(answer_struct.Bids_price_list[i], 'f', -1, 64)
		float_b, _ := strconv.ParseFloat(bids_string_map[str_b], 64)

		answer_struct.Bids_quantity_list[i] = float_b
		str_a := strconv.FormatFloat(answer_struct.Asks_price_list[i], 'f', -1, 64)
		float_a, _ := strconv.ParseFloat(asks_string_map[str_a], 64)
		answer_struct.Asks_quantity_list[i] = float_a
	}

	return answer_struct
}

func (O *Order_book) return_float_list(signal string, string_map map[string]string) []float64 {
	temp_list := make([]float64, len(string_map))
	num := 0
	for price := range string_map {
		floatvar, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("该string不可转为float64", err)
		}
		temp_list[num] = floatvar
		num++
	}
	if signal == "a" {
		sort.Sort(sort.Float64Slice(temp_list))
	} else if signal == "b" {
		sort.Sort(sort.Reverse(sort.Float64Slice(temp_list)))
	} else {
		fmt.Println("signal err")
	}

	return temp_list[:11]
}

func (O *Order_book) Start() []byte {
	O.init()
	O.get_info_from_cross_feature_data()
	O.get_info_from_original_data()
	answer, err := json.Marshal(O.info_send_sub)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(O.info_send_sub)
	return answer
}
