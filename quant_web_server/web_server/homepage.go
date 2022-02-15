package web_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quant_web_server/cross_feature"
	"quant_web_server/data_gather"
	"quant_web_server/data_process_1"
	"quant_web_server/feature_extraction"
	"quant_web_server/goal_feature"
	"quant_web_server/graph_analysis_package"
	"quant_web_server/order_book_analysis"
	"quant_web_server/read_mysql"
	"strconv"
)

type Web_server struct {
	//mysql
	Username           string
	Password           string
	Original_data      string
	Processed_data     string
	Cross_feature_data string
	//local_server
	Port_number string
	//local_book_deep
	Local_book_deep int
	//feature_extraction
	Time_period              int
	Goal_feature_time_period int
	//big_point(big_or_not)
	Divided_point float64
	//graph
	Time_point_period int
}

func (W *Web_server) Start_server() {
	http.HandleFunc("/", homepage)
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/data_gather", data_gather_)
	http.HandleFunc("/data_process", data_process)
	http.HandleFunc("/feature_and_more", data_feature_and_more)
	http.HandleFunc("/graph_analysis", graph_analysis)
	http.HandleFunc("/data_review/query", W.data_review_query)
	http.HandleFunc("/data_review/query2", W.data_review_query_2)
	http.HandleFunc("/data_gather/query", W.data_gather_query)
	http.HandleFunc("/data_process/classify", W.data_process_classify)
	http.HandleFunc("/data_process/order_book_analysis", W.data_process_orderbook_analysis)
	http.HandleFunc("/feature_and_more/get_feature", W.get_feature)
	http.HandleFunc("/feature_and_more/goal_feature", W.goal_feature)
	http.HandleFunc("/feature_and_more/cross_feature", W.cross_feature)
	http.HandleFunc("/graph_analysis/price_by_second", W.ga_price_by_second)
	http.HandleFunc("/graph_analysis/big_agg_trade", W.ga_big_agg_trade)
	http.HandleFunc("/graph_analysis/price_detail", W.ga_price_detail)
	http.HandleFunc("/graph_analysis/orderbook", W.ga_orderbook)
	http.ListenAndServe("127.0.0.1:"+W.Port_number, nil)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/homepage.html")
}

func (W *Web_server) data_review_query(w http.ResponseWriter, r *http.Request) {
	new_mysql := read_mysql.My_mysql{Username: W.Username, Password: W.Password, Original_data: W.Original_data, Processed_data: W.Processed_data}
	new_mysql.Init()
	defer new_mysql.Close_mysql()
	data := new_mysql.Show_original_data_table_info()
	new_info, _ := json.Marshal(data)
	w.Write(new_info)
}

func (W *Web_server) data_review_query_2(w http.ResponseWriter, r *http.Request) {
	new_mysql := read_mysql.My_mysql{Username: W.Username, Password: W.Password, Original_data: W.Original_data, Processed_data: W.Processed_data}
	new_mysql.Init()
	defer new_mysql.Close_mysql()
	data := new_mysql.Show_processed_data_table_info()
	new_info, _ := json.Marshal(data)
	w.Write(new_info)
}

func (W *Web_server) data_gather_query(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	time_length, _ := strconv.ParseInt(quert["time_length"][0], 10, 64)
	time_wait, _ := strconv.ParseInt(quert["time_wait"][0], 10, 64)
	symbol := quert["symbol"][0]
	fmt.Println("收到咨询查询请求，时长（秒数）：", time_length, "等待次数：", time_wait, "symbol:", symbol)
	d_g := data_gather.Data_gather{Username: W.Username, Password: W.Password, Original_data: W.Original_data, Symbol: symbol, Long: int(time_length), Num_wait: int(time_wait), Order_book_depth: W.Local_book_deep}
	d_g.Start()
	answer_info := []string{"数据收集已开启"}
	an, _ := json.Marshal(answer_info)
	w.Write(an)
}

func (W *Web_server) data_process_classify(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到原始数据分类请求，文件名：", data_tittle)
	d_p := data_process_1.Data_process{File_name: data_tittle, Username: W.Username, Password: W.Password, Original_data: W.Original_data, Processed_data: W.Processed_data}
	d_p.Process_data()
}

func (W *Web_server) data_process_orderbook_analysis(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到订单簿分析请求，文件名：", data_tittle)
	o_b_a := order_book_analysis.Order_book_analysis{K_: 0.05, File_name: data_tittle, Username: W.Username, Password: W.Password, Original_data: W.Original_data, Processed_data: W.Processed_data}
	o_b_a.Get_order_book_info_and_save()
}

func (W *Web_server) get_feature(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到提取特征请求，文件名：", data_tittle)
	f_e := feature_extraction.Feature_extraction{File_name: data_tittle, Username: W.Username, Password: W.Password, Processed_data: W.Processed_data, Cross_feature_data: W.Cross_feature_data, Time_period: W.Time_period}
	f_e.Get_all_info()
}

func (W *Web_server) goal_feature(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到提取目标变量请求，文件名：", data_tittle)
	g_f := goal_feature.Goal_feature{File_name: data_tittle, Username: W.Username, Password: W.Password, Processed_data: W.Processed_data, Cross_feature_data: W.Cross_feature_data, Goal_feature_time_period: W.Goal_feature_time_period}
	g_f.Get_all_info()
}

func (W *Web_server) cross_feature(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到构造交叉变量请求，文件名：", data_tittle)
	c_f := cross_feature.Cross_feature{File_name: data_tittle, Username: W.Username, Password: W.Password, Processed_data: W.Processed_data, Cross_feature_data: W.Cross_feature_data, Divided_point: W.Divided_point}
	c_f.Make_all_list()
}

func data_gather_(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/get_info.html")
}

func data_process(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/process.html")
}

func data_feature_and_more(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/feature_and_more.html")
}

func graph_analysis(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/info_view.html")
}

func (W *Web_server) ga_price_by_second(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到ga_price_by_second请求，文件名：", data_tittle)
	g_s_p := graph_analysis_package.Price_by_second{File_name: data_tittle, Username: W.Username, Password: W.Password, Cross_feature_data: W.Cross_feature_data}
	w.Write(g_s_p.Start())
}

func (W *Web_server) ga_big_agg_trade(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	fmt.Println("收到ga_big_agg_trade请求，文件名：", data_tittle)
	g_b_a_t := graph_analysis_package.Big_agg_trade{File_name: data_tittle, Username: W.Username, Password: W.Password, Cross_feature_data: W.Cross_feature_data}
	w.Write(g_b_a_t.Start())
}

func (W *Web_server) ga_price_detail(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	data_tittle := quert["data_tittle"][0]
	time_point := quert["time_point"][0]
	fmt.Println("收到ga_price_detail，文件名：", data_tittle, "时间点id：", time_point)
	time_point_int, err := strconv.Atoi(time_point)
	if err != nil {
		fmt.Println("错误，所输入时间点不是整数")
	}
	g_p_d := graph_analysis_package.Price_detial{File_name: data_tittle, Username: W.Username, Password: W.Password, Cross_feature_data: W.Cross_feature_data, Processed_data: W.Processed_data, Time_point_period: W.Time_point_period, Time_point_id: time_point_int}
	w.Write(g_p_d.Get_info_list())

}

func (W *Web_server) ga_orderbook(w http.ResponseWriter, r *http.Request) {
	quert := r.URL.Query()
	fmt.Println(quert)
	data_title := quert["data_tittle"][0]
	time_point := quert["time_point"][0]
	fmt.Println("收到ga_orderbook，文件名：", data_title, "时间点id：", time_point)
	time_point_int, err := strconv.Atoi(time_point)
	if err != nil {
		fmt.Println("错误，所输入时间点不是整数")
	}
	g_o := graph_analysis_package.Order_book{File_name: data_title, Username: W.Username, Password: W.Password, Cross_feature_data: W.Cross_feature_data, Original_data: W.Original_data, Processed_data: W.Processed_data, Time_start_id: time_point_int}
	w.Write(g_o.Start())
}
