package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Make_the_matrix struct {
	File_name             string //必填
	Username              string //必填
	Password              string //必填
	Cross_feature_data    string //必填
	cross_feature_data_db *sql.DB
	K_value_data          string //必填
	K_value_data_db       *sql.DB
	info_sql              string
	period_list           []int
	label                 string
	Max_int               int
	list_length           int
}

func (M *Make_the_matrix) init() {
	if len(M.period_list) == 0 {
		M.period_list = []int{3, 5, 10, 20}
	}
	M.Max_int = M.period_list[len(M.period_list)-1]
	dsn := M.Username + ":" + M.Password + "@tcp(127.0.0.1:3306)/" + M.Cross_feature_data
	dsn_2 := M.Username + ":" + M.Password + "@tcp(127.0.0.1:3306)/" + M.K_value_data
	var err error
	M.cross_feature_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = M.cross_feature_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	M.K_value_data_db, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = M.K_value_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	// M.info_sql = "(id int PRIMARY KEY AUTO_INCREMENT,T_start bigint,p_agg_aver_p3 double,p_agg_aver_p5 double,p_agg_aver_p10 double,p_agg_aver_p20 double,q_agg_aver_p3 double,q_agg_aver_p5 double,q_agg_aver_p10 double,q_agg_aver_p20 double,pq_agg_sum_p3 double,pq_agg_sum_p5 double,pq_agg_sum_p10 double,pq_agg_sum_p20 double,price_change_p3 double,price_change_p5 double,price_change_p10 double,price_change_p20 double,counts_sum_p3 int,counts_sum_p5 int,counts_sum_p10 int,counts_sum_p20 int,counts_m1_sum_p3 int,counts_m1_sum_p5 int,counts_m1_sum_p10 int,counts_m1_sum_p20 int,counts_m1_proportion_p3 double,counts_m1_proportion_p5 double,counts_m1_proportion_p10 double,counts_m1_proportion_p20 double,m_pq_sum_p3 double,m_pq_sum_p5 double,m_pq_sum_p10 double,m_pq_sum_p20 double,m_pq_proportion_p3 double,m_pq_proportion_p5 double,m_pq_proportion_p10 double,m_pq_proportion_p20 double,R_aver_p3 double,R_aver_p5 double,R_aver_p10 double,R_aver_p20 double,orderbook_abs_ratio_p3 double,orderbook_abs_ratio_p5 double,orderbook_abs_ratio_p10 double,orderbook_abs_ratio_p20 double,orderbook_volum_ratio_p3 double,orderbook_volum_ratio_p5 double,orderbook_volum_ratio_p10 double,orderbook_volum_ratio_p20 double,orderbook_average_ratio_p3 double,orderbook_average_ratio_p5 double,orderbook_average_ratio_p10 double,orderbook_average_ratio_p20 double,orderbook_varance_ratio_p3 double,orderbook_varance_ratio_p5 double,orderbook_varance_ratio_p10 double,orderbook_varance_ratio_p20 double)"
	M.info_sql = fmt.Sprintf("(id int PRIMARY KEY AUTO_INCREMENT,T_start bigint,p_agg_aver_p%v double,p_agg_aver_p%v double,p_agg_aver_p%v double,p_agg_aver_p%v double,q_agg_aver_p%v double,q_agg_aver_p%v double,q_agg_aver_p%v double,q_agg_aver_p%v double,pq_agg_sum_p%v double,pq_agg_sum_p%v double,pq_agg_sum_p%v double,pq_agg_sum_p%v double,price_change_p%v double,price_change_p%v double,price_change_p%v double,price_change_p%v double,counts_sum_p%v int,counts_sum_p%v int,counts_sum_p%v int,counts_sum_p%v int,counts_m1_sum_p%v int,counts_m1_sum_p%v int,counts_m1_sum_p%v int,counts_m1_sum_p%v int,counts_m1_proportion_p%v double,counts_m1_proportion_p%v double,counts_m1_proportion_p%v double,counts_m1_proportion_p%v double,m_pq_sum_p%v double,m_pq_sum_p%v double,m_pq_sum_p%v double,m_pq_sum_p%v double,m_pq_proportion_p%v double,m_pq_proportion_p%v double,m_pq_proportion_p%v double,m_pq_proportion_p%v double,R_aver_p%v double,R_aver_p%v double,R_aver_p%v double,R_aver_p%v double,orderbook_abs_ratio_p%v double,orderbook_abs_ratio_p%v double,orderbook_abs_ratio_p%v double,orderbook_abs_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_volum_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_average_ratio_p%v double,orderbook_varance_ratio_p%v double,orderbook_varance_ratio_p%v double,orderbook_varance_ratio_p%v double,orderbook_varance_ratio_p%v double)", M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3], M.period_list[0], M.period_list[1], M.period_list[2], M.period_list[3])
}

func get_time_strat(tittle string) int64 {
	signal := 0
	for i := 0; i < len(tittle); i++ {
		if tittle[i] < 97 {
			signal = i
			break
		}
	}
	t, _ := time.ParseInLocation("2006_01_02_15h04m05s", tittle[signal:], time.Local)
	return t.UnixNano() / 1e6
}

func (M *Make_the_matrix) creat_table() string {
	label := M.File_name + "data_matrix"
	sql := "CREATE TABLE " + label + M.info_sql + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := M.K_value_data_db.Prepare(sql)
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
	M.label = label
	return label
}

type info_box struct {
	id                      int
	T_start                 int64   //开始时间
	Ap_aver                 float64 //卖方最优报价价格均值
	Aq_aver                 float64 //卖方最优报价出售量均值
	Bp_aver                 float64 //买方最优报价价格均值
	Bq_aver                 float64 //买方最优报价出售量均值
	P_agg_aver              float64 //时间段内归集交易平均价格
	Q_agg_sum               float64 //时间段内归集交易总交易数量
	PQ_agg_sum              float64 //时间段内归集交易总交易价值量
	Counts_sum              int     //总的交易笔数
	Counts_m1_sum           int     //总交易笔数中买入单所占比重
	M_PQ_sum                float64 //买入单价值量占总交易价值量的比值
	R_aver                  float64 //平均资金费率，和杠杆市场与现货市场币值的价差相关
	orderbook_abs_ratio     float64 //指定深度订单簿买卖双方出价区间比值
	orderbook_volum_ratio   float64 //总体币的数量总量的比值
	orderbook_average_ratio float64 //指定双方买卖均值距离交易价值距离比
	orderbook_varance_ratio float64 //买卖双方指定深度价格分布方差比值
}

type table_length struct {
	length int
}

func (M *Make_the_matrix) make_update_sql(data_period int, info_list_1 []float64, info_list_2 []int, id int) string {
	sql_1 := "update " + M.label + " set "
	number := strconv.Itoa(data_period)
	sql_2 := fmt.Sprintf(" p_agg_aver_p"+number+"=%f"+", q_agg_aver_p"+number+"=%f"+", pq_agg_sum_p"+number+"=%f"+", counts_sum_p"+number+"=%v"+", counts_m1_sum_p"+number+"=%v"+", m_pq_sum_p"+number+"=%f"+", r_aver_p"+number+"=%f"+", orderbook_abs_ratio_p"+number+"=%f"+", orderbook_volum_ratio_p"+number+"=%f"+", orderbook_average_ratio_p"+number+"=%f"+", orderbook_varance_ratio_p"+number+"=%f"+" where id=%v;", info_list_1[0], info_list_1[1], info_list_1[2], info_list_2[0], info_list_2[1], info_list_1[3], info_list_1[4], info_list_1[5], info_list_1[6], info_list_1[7], info_list_1[8], id)
	sql := sql_1 + sql_2
	return sql
}

func (M *Make_the_matrix) full_the_null_table() {
	sql_query := "select id,T_start from " + M.File_name + "feature_extraction_p1" + " where id=?;"
	stmt_query, err := M.cross_feature_data_db.Prepare(sql_query)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_query.Close()
	sql_insert := "insert into " + M.label + " (T_start) values(?)"
	stmt_insert, err := M.K_value_data_db.Prepare(sql_insert)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_insert.Close()
	var i_f info_box
	data_id := 1
	id := M.Max_int + 1
	fmt.Println("表格开始初始化")
	for {
		data := stmt_query.QueryRow(id)
		data.Scan(&i_f.id, &i_f.T_start)
		stmt_insert.Exec(i_f.T_start)
		id++
		if (data_id - i_f.id) == 0 {
			M.list_length = data_id - M.Max_int + 1
			break
		}
		data_id = i_f.id
	}
	fmt.Println("表格初始化完成，列表总长度：", M.list_length+1)

}

func (M *Make_the_matrix) save_period_data_1(data_period int) {
	// sql_query := "select id,T_start,P_agg_sum,Q_agg_sum,PQ_agg_sum,Counts_sum,Counts_m1_sum,M_PQ_sum,R_aver,orderbook_abs_ratio,orderbook_volum_ratio,orderbook_average_ratio,orderbook_varance_ratio from " + M.File_name + "feature_extraction_p1" + " where id>?;"
	sql_query_2 := "select id,T_start,P_agg_sum,Q_agg_sum,PQ_agg_sum,Counts_sum,Counts_m1_sum,M_PQ_sum,R_aver,orderbook_abs_ratio,orderbook_volum_ratio,orderbook_average_ratio,orderbook_varance_ratio from " + M.File_name + "feature_extraction_p1" + " where id>=?;"

	id := M.Max_int + 1
	var i_f info_box
	for i := 0; i < M.list_length; i++ {
		fmt.Println("original_feature_p", data_period, ":", i+1, "/", M.list_length)
		p_agg_aver := float64(0)
		q_agg_aver := float64(0)
		pq_agg_sum := float64(0)
		counts_sum := int(0)
		counts_m1_sum := int(0)
		m_pq_sum := float64(0)
		r_aver := float64(0)
		orderbook_abs_ratio := float64(0)
		orderbook_volum_ratio := float64(0)
		orderbook_average_ratio := float64(0)
		orderbook_varance_ratio := float64(0)
		r_aver_num := 0
		data_2, err := M.cross_feature_data_db.Query(sql_query_2, id-data_period)
		signal := 0
		for data_2.Next() {
			signal++
			data_2.Scan(&i_f.id, &i_f.T_start, &i_f.P_agg_aver, &i_f.Q_agg_sum, &i_f.PQ_agg_sum, &i_f.Counts_sum, &i_f.Counts_m1_sum, &i_f.M_PQ_sum, &i_f.R_aver, &i_f.orderbook_abs_ratio, &i_f.orderbook_volum_ratio, &i_f.orderbook_average_ratio, &i_f.orderbook_varance_ratio)
			p_agg_aver += i_f.P_agg_aver
			q_agg_aver += i_f.Q_agg_sum
			pq_agg_sum += i_f.PQ_agg_sum
			counts_sum += i_f.Counts_sum
			counts_m1_sum += i_f.Counts_m1_sum
			m_pq_sum += i_f.M_PQ_sum
			orderbook_abs_ratio += i_f.orderbook_abs_ratio
			orderbook_volum_ratio += i_f.orderbook_volum_ratio
			orderbook_average_ratio += i_f.orderbook_average_ratio
			orderbook_varance_ratio += i_f.orderbook_varance_ratio
			r_aver += i_f.R_aver
			if i_f.R_aver != 0 {
				r_aver_num++
			}
			if signal == data_period {
				break
			}
		}
		data_2.Close()
		p_agg_aver = p_agg_aver / float64(data_period)
		q_agg_aver = q_agg_aver / float64(data_period)
		orderbook_abs_ratio = orderbook_abs_ratio / float64(data_period)
		orderbook_average_ratio = orderbook_average_ratio / float64(data_period)
		orderbook_volum_ratio = orderbook_volum_ratio / float64(data_period)
		orderbook_varance_ratio = orderbook_varance_ratio / float64(data_period)
		if r_aver_num != 0 {
			r_aver = r_aver / float64(r_aver_num)
		}
		info_list_1 := []float64{p_agg_aver, q_agg_aver, pq_agg_sum, m_pq_sum, r_aver, orderbook_abs_ratio, orderbook_volum_ratio, orderbook_average_ratio, orderbook_varance_ratio}
		info_list_2 := []int{counts_sum, counts_m1_sum}
		sql_update := M.make_update_sql(data_period, info_list_1, info_list_2, id-M.Max_int)
		stmt_update, err := M.K_value_data_db.Prepare(sql_update)
		if err != nil {
			fmt.Println("ERR!", err)
		}
		stmt_update.Exec()
		id++
	}
}

type price_change struct {
	bids_start float64
	bids_end   float64
}

func (p *price_change) return_change() float64 {
	return (p.bids_end - p.bids_start)
}

func (M *Make_the_matrix) make_update_sql_2(price_change float64, data_period int, id int) string {
	str_1 := "update " + M.label + " set price_change_p"
	str_2 := fmt.Sprintf("%v=%v", data_period, price_change)
	str_3 := fmt.Sprintf(" where id=%v;", id)
	str := str_1 + str_2 + str_3
	return str
}

func (M *Make_the_matrix) save_period_data_2(data_period int) {
	sql_query := "select bids_start,bids_end from " + M.File_name + "goal_feature_p1 " + "where id = ?"
	stmt_query, err := M.cross_feature_data_db.Prepare(sql_query)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_query.Close()
	id := M.Max_int + 1
	var p_c_1 price_change
	var p_c_2 price_change
	for i := 0; i < M.list_length; i++ {
		fmt.Println("price_change_p", data_period, ":", i+1, "/", M.list_length)
		data_1 := stmt_query.QueryRow(id - data_period)
		data_2 := stmt_query.QueryRow(id - 1)
		data_1.Scan(&p_c_1.bids_start, &p_c_1.bids_end)
		data_2.Scan(&p_c_2.bids_start, &p_c_2.bids_end)
		p_c_1.bids_end = p_c_2.bids_end
		sql_update := M.make_update_sql_2(p_c_1.return_change(), data_period, id-M.Max_int)
		M.K_value_data_db.Exec(sql_update)
		id++
	}

}

type proportion struct {
	pq_agg_sum    float64
	m_pq_sum_p    float64
	counts_sum    float64
	counts_m1_sum float64
}

func (p *proportion) return_proportion() (float64, float64) {
	answer_1 := float64(0)
	answer_2 := float64(0)
	if p.pq_agg_sum != 0 {
		answer_1 = p.m_pq_sum_p / p.pq_agg_sum
	}
	if p.counts_sum != 0 {
		answer_2 = p.counts_m1_sum / p.counts_sum
	}
	return answer_1, answer_2
}

func (M *Make_the_matrix) save_period_data_3(data_period int) {
	tool_sql_1 := fmt.Sprintf("pq_agg_sum_p%v,m_pq_sum_p%v,counts_sum_p%v,counts_m1_sum_p%v", data_period, data_period, data_period, data_period)
	sql_query := "select " + tool_sql_1 + " from " + M.label + " where id=?"
	fmt.Println(sql_query)
	stmt_query, err := M.K_value_data_db.Prepare(sql_query)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_query.Close()
	var p proportion
	for id := 1; id <= M.list_length; id++ {
		fmt.Println("proportion_p", data_period, ":", id, "/", M.list_length)
		data := stmt_query.QueryRow(id)
		data.Scan(&p.pq_agg_sum, &p.m_pq_sum_p, &p.counts_sum, &p.counts_m1_sum)
		m_pq_proportion, counts_m1_proportion := p.return_proportion()
		sql_update := "update " + M.label + fmt.Sprintf(" set m_pq_proportion_p%v=%v,counts_m1_proportion_p%v=%v where id=%v;", data_period, m_pq_proportion, data_period, counts_m1_proportion, id)
		M.K_value_data_db.Exec(sql_update)
	}

}

func (M *Make_the_matrix) Get_start() {
	M.init()
	M.creat_table()
	M.full_the_null_table()
	for i := 0; i < len(M.period_list); i++ {
		M.save_period_data_1(M.period_list[i])
		M.save_period_data_2(M.period_list[i])
		M.save_period_data_3(M.period_list[i])
	}
}

func main() {
	m_k_t := Make_the_matrix{File_name: "btcusdt2022_01_01_19h05m11s", Username: "root", Password: "", Cross_feature_data: "cross_feature_data", K_value_data: "classify_data"}
	m_k_t.Get_start()

}
