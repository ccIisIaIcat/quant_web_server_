package local_order_book_2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"quant_web_server/get_info"
	"strconv"
)

type My_orderbook struct {
	Outputer    *chan []byte //必填
	Symbol      string       //必填
	pa          []string
	info_gather chan []byte
	url         string
	Max_length  int //建议填入,默认为100
	A_max       float64
	B_min       float64
	Judge       bool
	n_wb        get_info.Ask_a_wb
}

type order_photo struct {
	LastUpdateId int        `json:"lastUpdateId"`
	E            int        `json:"E"` // 消息时间
	T            int64      `json:"T"` // 撮合引擎时间
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type order_updata struct {
	T      int64      `json:"T"`
	U_1    int        `json:"U"` //新增的第一个ID
	U_2    int        `json:"u"` //新增的最后一个ID
	B_list [][]string `json:"b"`
	A_list [][]string `json:"a"`
}

type Out_put struct {
	T     int64               `json:"T"`
	B_dic map[float64]float64 `json:"bids"`
	A_dic map[float64]float64 `json:"asks"`
}

type Out_put_2 struct {
	T     int64             `json:"T"`
	B_dic map[string]string `json:"bids"`
	A_dic map[string]string `json:"asks"`
}

//初始化
func (M *My_orderbook) init() {
	if M.Max_length == 0 {
		M.Max_length = 100
	}
	M.Judge = false
	M.info_gather = make(chan []byte, 100)
	M.pa = []string{M.Symbol + "@depth@100ms"}
	M.url = "https://fapi.binance.com/fapi/v1/depth?symbol=" + M.Symbol + "&limit=100"

}

func (M *My_orderbook) End_local_order_book() {
	M.Judge = true
}

//更新订单簿A
func (M *My_orderbook) update_orderbook_A(original_orderbook_list *map[float64]float64, update_infomation *[][]string) {
	updata_new := turn_stringmatric_to_float64_list(*update_infomation)
	for i := 0; i < len(*update_infomation); i++ {
		if updata_new[i][1] == 0 {
			delete(*original_orderbook_list, updata_new[i][0])
		} else {
			if len(*original_orderbook_list) < M.Max_length {
				(*original_orderbook_list)[updata_new[i][0]] = updata_new[i][1]
				if updata_new[i][0] > M.A_max {
					M.A_max = updata_new[i][0]
				}
			} else if updata_new[i][0] < M.A_max {
				delete((*original_orderbook_list), M.A_max)
				(*original_orderbook_list)[updata_new[i][0]] = updata_new[i][1]
				M.A_max = find_max_in_map(*original_orderbook_list)
			}
		}
	}
}

//更新订单簿B
func (M *My_orderbook) update_orderbook_B(original_orderbook_list *map[float64]float64, update_infomation *[][]string) {
	updata_new := turn_stringmatric_to_float64_list(*update_infomation)
	for i := 0; i < len(*update_infomation); i++ {
		if updata_new[i][1] == 0 {
			delete(*original_orderbook_list, updata_new[i][0])
		} else {
			if len(*original_orderbook_list) < M.Max_length {
				(*original_orderbook_list)[updata_new[i][0]] = updata_new[i][1]
				if updata_new[i][0] < M.B_min {
					M.B_min = updata_new[i][0]
				}
			} else if updata_new[i][0] > M.B_min {
				delete((*original_orderbook_list), M.B_min)
				(*original_orderbook_list)[updata_new[i][0]] = updata_new[i][1]
				M.B_min = find_min_in_map(*original_orderbook_list)
			}
		}
	}
}

func find_max_in_map(data map[float64]float64) float64 {
	max := float64(0)
	for k := range data {
		if k >= max {
			max = k
		}
	}
	return max
}

func find_min_in_map(data map[float64]float64) float64 {
	min := float64(10000000000)
	for k := range data {
		if k <= min {
			min = k
		}
	}
	return min
}

func (M *My_orderbook) get_original_order_book() {
	M.n_wb = get_info.Ask_a_wb{Info_type: get_info.Symbol_Depth_addition, Params: M.pa}
	M.n_wb.Start_wb(&M.info_gather)
	data := <-M.info_gather
	fmt.Println("深度增量的websocket已开启", string(data))
	re, err := http.Get(M.url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(re.Body)
	var order_book_new order_photo
	json.Unmarshal(body, &order_book_new)
	// M.A_max = order_book_new.Asks[M.Max_length-1]
}

func turn_stringmatric_to_float64_list(data [][]string) [][]float64 {
	answer := make([][]float64, len(data))
	for i := 0; i < len(data); i++ {
		temp_arr := make([]float64, len(data[i]))
		for j := 0; j < len(data[i]); j++ {
			temp_arr[j], _ = strconv.ParseFloat(data[i][j], 64)
		}
		answer[i] = temp_arr
	}
	return answer
}

func turn_stringmatric_to_float64_map(data [][]string) map[float64]float64 {
	answer := make(map[float64]float64, 0)
	for i := 0; i < len(data); i++ {
		a, _ := strconv.ParseFloat(data[i][0], 64)
		b, _ := strconv.ParseFloat(data[i][1], 64)
		answer[a] = b
	}
	return answer
}

func turn_float64_map_to_string_map(data map[float64]float64) map[string]string {
	answer := make(map[string]string, 0)
	for k, v := range data {
		a := strconv.FormatFloat(k, 'g', 10, 64)
		b := strconv.FormatFloat(v, 'g', 10, 64)
		answer[a] = b
	}
	return answer
}

//开启服务
func (M *My_orderbook) Start_serve() {
	M.init()
	n_wb := get_info.Ask_a_wb{Info_type: get_info.Symbol_Depth_addition, Params: M.pa}
	n_wb.Start_wb(&M.info_gather)
	data := <-M.info_gather
	fmt.Println("深度增量的websocket已开启", string(data))
	re, err := http.Get(M.url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(re.Body)
	var order_book_new order_photo
	json.Unmarshal(body, &order_book_new)
	var out_puter Out_put
	a_list := turn_stringmatric_to_float64_list(order_book_new.Asks)
	b_list := turn_stringmatric_to_float64_list(order_book_new.Bids)
	out_puter.A_dic = turn_stringmatric_to_float64_map(order_book_new.Asks)
	out_puter.B_dic = turn_stringmatric_to_float64_map(order_book_new.Bids)
	LU := order_book_new.LastUpdateId
	out_puter.T = order_book_new.T
	M.A_max = a_list[M.Max_length-1][0]
	M.B_min = b_list[M.Max_length-1][0]
	var order_up order_updata

	for {
		if M.Judge {
			fmt.Println("order_book通道已关闭")
			break
		}
		data = <-M.info_gather
		out_puter.T = order_up.T
		json.Unmarshal(data, &order_up)
		if order_up.U_2 > LU {
			M.update_orderbook_B(&out_puter.B_dic, &order_up.B_list)
			M.update_orderbook_A(&out_puter.A_dic, &order_up.A_list)
			var temp Out_put_2
			temp.T = order_up.T
			temp.A_dic = turn_float64_map_to_string_map(out_puter.A_dic)
			temp.B_dic = turn_float64_map_to_string_map(out_puter.B_dic)
			anan, _ := json.Marshal(temp)

			(*M.Outputer) <- anan
		}

	}
}

// func main() {
// 	my_chan := make(chan []byte, 10)
// 	symbol := "btcusdt"
// 	mo := My_orderbook{Symbol: symbol, Outputer: &my_chan}
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	a := make([]interface{}, 0)
// 	go mo.Start_serve()
// 	for {
// 		data := <-my_chan
// 		json.Unmarshal(data, &a)
// 	}
// 	//wg.Wait()
// }
