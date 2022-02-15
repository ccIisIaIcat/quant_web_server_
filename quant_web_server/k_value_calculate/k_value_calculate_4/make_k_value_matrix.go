package main

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type K_matrix struct {
	File_name        string //必填
	Username         string //必填
	Password         string //必填
	K_value_data     string //必填
	Goal_label       string //必填，目标变量表名
	K_value_data_db  *sql.DB
	label            string //构建新列表的名称
	dividied_point   []float64
	data_long_matrix [][]int  //用于处理预测变量级别
	info_sql         string   //用于构建新列表的sql语句
	table_info       []string //各列的名称列表
	period_list      []int    //预测变量的时间段
	goal_list        []int    //存储对应目标变量的分类
	id_gap           int      //预测变量和结果变量的id差值
	pa_              float64  //大订单比率
	period_num       int      //多参数数目
}

func (K *K_matrix) init() {
	K.data_long_matrix = make([][]int, 52)
	if K.id_gap == 0 {
		K.id_gap = 10
	}
	if K.Goal_label == "" {
		K.Goal_label = K.File_name + "goal_matrix_p10"
	}
	if len(K.period_list) == 0 {
		K.period_list = []int{3, 5, 10, 20}
	}
	dsn := K.Username + ":" + K.Password + "@tcp(127.0.0.1:3306)/" + K.K_value_data
	var err error
	K.K_value_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = K.K_value_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	K.info_sql = fmt.Sprintf("(id int PRIMARY KEY AUTO_INCREMENT,q_agg_aver_p%v double,q_agg_aver_p%v double,q_agg_aver_p%v double,q_agg_aver_p%v double,pq_agg_sum_p%v double,pq_agg_sum_p%v double,pq_agg_sum_p%v double,pq_agg_sum_p%v double,price_change_p%v double,price_change_p%v double,price_change_p%v double,price_change_p%v double,counts_sum_p%v double,counts_sum_p%v double,counts_sum_p%v double,counts_sum_p%v double,counts_m1_sum_p%v double,counts_m1_sum_p%v double,counts_m1_sum_p%v double,counts_m1_sum_p%v double,counts_m1_proportion_p%v double,counts_m1_proportion_p%v double,counts_m1_proportion_p%v double,counts_m1_proportion_p%v double,m_pq_sum_p%v double,m_pq_sum_p%v double,m_pq_sum_p%v double,m_pq_sum_p%v double,m_pq_proportion_p%v double,m_pq_proportion_p%v double,m_pq_proportion_p%v double,m_pq_proportion_p%v double,R_aver_p%v double,R_aver_p%v double,R_aver_p%v double,R_aver_p%v double,orderbook_abs_ratio_p%v double,orderbook_abs_ratio_p%v double,orderbook_abs_ratio_p%v double,orderbook_abs_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_varance_ratio_p%v double,orderbook_varance_ratio_p%v double,orderbook_varance_ratio_p%v double,orderbook_varance_ratio_p%v double)", K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3], K.period_list[0], K.period_list[1], K.period_list[2], K.period_list[3])
	K.table_info = K.make_table_info()
}

func (K *K_matrix) make_table_info() []string {
	answer := make([]string, 0)
	for i := 0; i < len(K.period_list); i++ {
		answer = append(answer, fmt.Sprintf("price_change_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("q_agg_aver_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("pq_agg_sum_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("counts_sum_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("counts_m1_sum_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("counts_m1_proportion_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("m_pq_sum_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("m_pq_proportion_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("R_aver_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_abs_ratio_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_volum_ratio_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_average_ratio_p%v", K.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_varance_ratio_p%v", K.period_list[i]))
	}
	return answer
}

type all_info_with_order struct {
	price_change_1            int
	q_agg_aver_1              int
	pq_agg_sum_1              int
	counts_sum_1              int
	counts_m1_sum_1           int
	counts_m1_proportion_1    int
	m_pq_sum_1                int
	m_pq_proportion_1         int
	R_aver_1                  int
	orderbook_abs_ratio_1     int
	orderbook_volum_ratio_1   int
	orderbook_average_ratio_1 int
	orderbook_varance_ratio_1 int
	price_change_2            int
	q_agg_aver_2              int
	pq_agg_sum_2              int
	counts_sum_2              int
	counts_m1_sum_2           int
	counts_m1_proportion_2    int
	m_pq_sum_2                int
	m_pq_proportion_2         int
	R_aver_2                  int
	orderbook_abs_ratio_2     int
	orderbook_volum_ratio_2   int
	orderbook_average_ratio_2 int
	orderbook_varance_ratio_2 int
	price_change_3            int
	q_agg_aver_3              int
	pq_agg_sum_3              int
	counts_sum_3              int
	counts_m1_sum_3           int
	counts_m1_proportion_3    int
	m_pq_sum_3                int
	m_pq_proportion_3         int
	R_aver_3                  int
	orderbook_abs_ratio_3     int
	orderbook_volum_ratio_3   int
	orderbook_average_ratio_3 int
	orderbook_varance_ratio_3 int
	price_change_4            int
	q_agg_aver_4              int
	pq_agg_sum_4              int
	counts_sum_4              int
	counts_m1_sum_4           int
	counts_m1_proportion_4    int
	m_pq_sum_4                int
	m_pq_proportion_4         int
	R_aver_4                  int
	orderbook_abs_ratio_4     int
	orderbook_volum_ratio_4   int
	orderbook_average_ratio_4 int
	orderbook_varance_ratio_4 int
}

func (K *K_matrix) make_data_long_matrix() {
	sql_query := "select "
	for i := 0; i < len(K.table_info); i++ {
		sql_query += K.table_info[i]
		if i != len(K.table_info)-1 {
			sql_query += ","
		}
	}
	sql_query += " from " + K.File_name + "feature_matrix"
	data, err := K.K_value_data_db.Query(sql_query)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	var a_i_w_o all_info_with_order
	for data.Next() {
		data.Scan(&a_i_w_o.price_change_1, &a_i_w_o.q_agg_aver_1, &a_i_w_o.pq_agg_sum_1, &a_i_w_o.counts_sum_1, &a_i_w_o.counts_m1_sum_1, &a_i_w_o.counts_m1_proportion_1, &a_i_w_o.m_pq_sum_1, &a_i_w_o.m_pq_proportion_1, &a_i_w_o.R_aver_1, &a_i_w_o.orderbook_abs_ratio_1, &a_i_w_o.orderbook_volum_ratio_1, &a_i_w_o.orderbook_average_ratio_1, &a_i_w_o.orderbook_varance_ratio_1, &a_i_w_o.price_change_2, &a_i_w_o.q_agg_aver_2, &a_i_w_o.pq_agg_sum_2, &a_i_w_o.counts_sum_2, &a_i_w_o.counts_m1_sum_2, &a_i_w_o.counts_m1_proportion_2, &a_i_w_o.m_pq_sum_2, &a_i_w_o.m_pq_proportion_2, &a_i_w_o.R_aver_2, &a_i_w_o.orderbook_abs_ratio_2, &a_i_w_o.orderbook_volum_ratio_2, &a_i_w_o.orderbook_average_ratio_2, &a_i_w_o.orderbook_varance_ratio_2, &a_i_w_o.price_change_3, &a_i_w_o.q_agg_aver_3, &a_i_w_o.pq_agg_sum_3, &a_i_w_o.counts_sum_3, &a_i_w_o.counts_m1_sum_3, &a_i_w_o.counts_m1_proportion_3, &a_i_w_o.m_pq_sum_3, &a_i_w_o.m_pq_proportion_3, &a_i_w_o.R_aver_3, &a_i_w_o.orderbook_abs_ratio_3, &a_i_w_o.orderbook_volum_ratio_3, &a_i_w_o.orderbook_average_ratio_3, &a_i_w_o.orderbook_varance_ratio_3, &a_i_w_o.price_change_4, &a_i_w_o.q_agg_aver_4, &a_i_w_o.pq_agg_sum_4, &a_i_w_o.counts_sum_4, &a_i_w_o.counts_m1_sum_4, &a_i_w_o.counts_m1_proportion_4, &a_i_w_o.m_pq_sum_4, &a_i_w_o.m_pq_proportion_4, &a_i_w_o.R_aver_4, &a_i_w_o.orderbook_abs_ratio_4, &a_i_w_o.orderbook_volum_ratio_4, &a_i_w_o.orderbook_average_ratio_4, &a_i_w_o.orderbook_varance_ratio_4)
		for i := 0; i < reflect.ValueOf(a_i_w_o).NumField(); i++ {
			K.data_long_matrix[i] = append(K.data_long_matrix[i], int(reflect.ValueOf(a_i_w_o).Field(i).Int()))
		}
	}
}

type big_or_nor_ struct {
	big_or_nor int
}

func (K *K_matrix) make_goal_list() {
	sql_query := "select big_or_not from " + K.Goal_label + fmt.Sprintf(" where id>%v", K.id_gap)
	data, err := K.K_value_data_db.Query(sql_query)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	var b_o_n big_or_nor_
	num := 0
	for data.Next() {
		data.Scan(&b_o_n.big_or_nor)
		num += b_o_n.big_or_nor

		K.goal_list = append(K.goal_list, b_o_n.big_or_nor)
	}
	K.pa_ = float64(num) / float64(len(K.goal_list))
	fmt.Println(K.goal_list)
	fmt.Println(len(K.goal_list))
}

func (K *K_matrix) creat_new_table() {
	sql := "create table " + K.File_name + "k_value" + K.info_sql + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	K.K_value_data_db.Exec(sql)
	K.label = K.File_name + "k_value"
	fmt.Println("建表成功：", K.label)
}

func (K *K_matrix) calculate_k_value() {

}

func (K *K_matrix) calculate_two_list_k_value(list_suq int) []float64 {

	sample_list := K.data_long_matrix[list_suq]
	//calculate P(B)
	answer_list := make([]float64, len(K.period_list)+1)
	//calculate P(A|B)
	tool_list := make([]float64, len(K.period_list)+1)
	return_list := make([]float64, len(K.period_list)+1)
	for j := 0; j < len(K.period_list)+1; j++ {
		for i := 0; i < len(K.goal_list); i++ {
			if sample_list[i] == j+1 {
				answer_list[j] += 1
			}
			if sample_list[i] == j+1 && K.goal_list[i] == 1 {
				tool_list[j] += 1
			}
		}
		answer_list[j] = answer_list[j] / float64(len(K.goal_list))
		tool_list[j] = tool_list[j] / float64(len(K.goal_list))
		return_list[j] = math.Abs(tool_list[j]/K.pa_/answer_list[j] - 1)
	}
	return return_list

}

func (K *K_matrix) calculate_two_list_k_value_2(list_suq int) []float64 {
	sample_list := K.data_long_matrix[list_suq]
	answer_list := make([]float64, K.period_num+1)
	tool_list := make([]float64, K.period_num+1)
	return_list := make([]float64, K.period_num+1)
	for j := 0; j < K.period_num+1; j++ {
		for i := 0; i < len(K.goal_list); i++ {
			if sample_list[i] == j+1 {
				answer_list[j] += 1
			}
			if sample_list[i] == j+1 && K.goal_list[i] == 1 {
				tool_list[j] += 1
			}
		}
		answer_list[j] = answer_list[j] / float64(len(K.goal_list))
		tool_list[j] = tool_list[j] / float64(len(K.goal_list))
		return_list[j] = math.Abs(tool_list[j]/K.pa_/answer_list[j] - 1)
	}
	return return_list

}

func (K *K_matrix) Start() {
	K.init()
	K.creat_new_table()
	K.make_data_long_matrix()
	K.make_goal_list()
	K.make_table_info()
}

func (K *K_matrix) Start_2() {
	K.init()
	K.creat_new_table()
	K.make_data_long_matrix()
	K.make_goal_list()

}

func main() {
	fmt.Println("lalala")
	k_m := K_matrix{File_name: "btcusdt2022_01_01_19h05m11s", Username: "root", Password: "", K_value_data: "classify_data", period_num: 10}
	k_m.init()
	k_m.make_table_info()
	k_m.make_data_long_matrix()
	k_m.make_goal_list()
	for i := 0; i < 52; i++ {
		fmt.Println(k_m.table_info[i])
		fmt.Println(k_m.calculate_two_list_k_value_2(i))
	}
	// k_m.creat_new_table()
}
