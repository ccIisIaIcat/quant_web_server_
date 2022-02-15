package feature_extraction

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Feature_extraction struct {
	File_name             string //必填
	Username              string //必填
	Password              string //必填
	Processed_data        string
	processed_data_db     *sql.DB
	Cross_feature_data    string
	cross_feature_data_db *sql.DB
	Time_start_point      int64 //起始点的时间点(必填)
	Time_period           int   //时间区间的秒数(必填)
}

type info_box struct {
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
	M_PQ_proportion         float64 //买入单价值量占总交易价值量的比值
	R_aver                  float64 //平均资金费率，和杠杆市场与现货市场币值的价差相关
	orderbook_abs_ratio     float64 //指定深度订单簿买卖双方出价区间比值
	orderbook_volum_ratio   float64 //总体币的数量总量的比值
	orderbook_average_ratio float64 //指定双方买卖均值距离交易价值距离比
	orderbook_varance_ratio float64 //买卖双方指定深度价格分布方差比值
}

func (F *Feature_extraction) init() {
	dsn := F.Username + ":" + F.Password + "@tcp(127.0.0.1:3306)/" + F.Processed_data
	dsn_2 := F.Username + ":" + F.Password + "@tcp(127.0.0.1:3306)/" + F.Cross_feature_data
	var err error
	F.processed_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = F.processed_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	F.cross_feature_data_db, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = F.cross_feature_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	F.Time_start_point = get_time_strat(F.File_name)
}

func (F *Feature_extraction) creat_table() string {
	label := F.File_name + "feature_extraction"
	period_str := strconv.Itoa(F.Time_period)
	label = label + "_p" + period_str
	sql := "CREATE TABLE " + label + "(id int PRIMARY KEY AUTO_INCREMENT,T_start bigint,Ap_aver double,Aq_aver double,Bp_aver double,Bq_aver double,P_agg_sum double,Q_agg_sum double,PQ_agg_sum double,Counts_sum int,Counts_m1_sum int,M_PQ_sum double,R_aver double,orderbook_abs_ratio double,orderbook_volum_ratio double,orderbook_average_ratio double,orderbook_varance_ratio double)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := F.cross_feature_data_db.Prepare(sql)
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
	return label
}

type best_offer_price struct {
	id  int
	T_1 int64
	Ap  float64
	Aq  float64
	Bp  float64
	Bq  float64
}

type agg_trade struct {
	id     int
	T_1    int64
	P      float64
	Q      float64
	Counts int
	M      int //1为主动卖出单，0为主动买入单
}

type newest_marked_price struct {
	id  int
	T_1 int64
	R   float64
}

type local_order_book_analysis struct {
	id            int
	Time          int64
	abs_ratio     float64
	volum_ratio   float64
	average_ratio float64
	varance_ratio float64
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

func (F *Feature_extraction) get_all_info_sub(stmt_best_offer_price *sql.Stmt, stmt_agg_trade *sql.Stmt, stmt_newest_marked_price *sql.Stmt, stmt_local_order_book_analysis *sql.Stmt) {
	new_label := F.creat_table()
	sql_insert := "insert into " + new_label + " (T_start,Ap_aver,Aq_aver,Bp_aver,Bq_aver,P_agg_sum,Q_agg_sum,PQ_agg_sum,Counts_sum,Counts_m1_sum,M_PQ_sum,R_aver,orderbook_abs_ratio,orderbook_volum_ratio,orderbook_average_ratio,orderbook_varance_ratio) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	stmt_insert, err := F.cross_feature_data_db.Prepare(sql_insert)
	if err != nil {
		fmt.Println("stmt_insert错误：", err)
	}
	defer stmt_insert.Close()
	id_1 := 1
	id_2 := 1
	id_3 := 1
	id_4 := 1
	start_time_now := F.Time_start_point
	for {
		signal := map[int]bool{1: false, 2: false, 3: false, 4: false}
		var i_b info_box
		i_b.T_start = start_time_now
		num_1 := 0
		for {
			data_1 := stmt_best_offer_price.QueryRow(id_1)
			var b_o_p best_offer_price
			data_1.Scan(&b_o_p.id, &b_o_p.T_1, &b_o_p.Ap, &b_o_p.Aq, &b_o_p.Bp, &b_o_p.Bq)
			if b_o_p.id == 0 {
				signal[1] = true
				if num_1 != 0 {
					i_b.Ap_aver /= float64(num_1)
					i_b.Aq_aver /= float64(num_1)
					i_b.Bp_aver /= float64(num_1)
					i_b.Bq_aver /= float64(num_1)
				}
				break
			}
			if b_o_p.T_1 > start_time_now+int64(F.Time_period*1000) {
				if num_1 != 0 {
					i_b.Ap_aver /= float64(num_1)
					i_b.Aq_aver /= float64(num_1)
					i_b.Bp_aver /= float64(num_1)
					i_b.Bq_aver /= float64(num_1)
				}
				break
			} else {
				id_1++
				num_1++
				i_b.Ap_aver += b_o_p.Ap
				i_b.Aq_aver += b_o_p.Aq
				i_b.Bp_aver += b_o_p.Bp
				i_b.Bq_aver += b_o_p.Bq
			}
		}
		num_2 := 0
		for {
			data_2 := stmt_agg_trade.QueryRow(id_2)
			var a_t agg_trade
			data_2.Scan(&a_t.id, &a_t.T_1, &a_t.P, &a_t.Q, &a_t.Counts, &a_t.M)
			if a_t.id == 0 {
				signal[2] = true
				if i_b.PQ_agg_sum != 0 {
					i_b.M_PQ_proportion = i_b.M_PQ_proportion / i_b.PQ_agg_sum
				}
				if num_2 != 0 {
					i_b.P_agg_aver /= float64(num_2)
				}
				break
			}
			if a_t.T_1 > start_time_now+int64(F.Time_period*1000) {
				if i_b.PQ_agg_sum != 0 {
					i_b.M_PQ_proportion = i_b.M_PQ_proportion * 1
				}
				if num_2 != 0 {
					i_b.P_agg_aver /= float64(num_2)
				}
				break
			} else {
				id_2++
				num_2++
				i_b.Counts_sum += a_t.Counts
				i_b.P_agg_aver += a_t.P
				i_b.Q_agg_sum += a_t.Q
				i_b.PQ_agg_sum += a_t.P * a_t.Q
				i_b.Counts_m1_sum += a_t.Counts * a_t.M
				i_b.M_PQ_proportion += a_t.P * a_t.Q * float64(a_t.M)
			}

		}
		num_3 := 0
		for {
			data_3 := stmt_newest_marked_price.QueryRow(id_3)
			var n_m_p newest_marked_price
			data_3.Scan(&n_m_p.id, &n_m_p.T_1, &n_m_p.R)
			if n_m_p.id == 0 {
				signal[3] = true
				if num_3 != 0 {
					i_b.R_aver = i_b.R_aver / float64(num_3)
				}
				break
			}

			if n_m_p.T_1 > start_time_now+int64(F.Time_period*1000) {
				if num_3 != 0 {
					i_b.R_aver = i_b.R_aver / float64(num_3)
				}
				break
			} else {
				id_3++
				num_3++
				i_b.R_aver += n_m_p.R
			}
		}
		num_4 := 0
		for {
			data_4 := stmt_local_order_book_analysis.QueryRow(id_4)
			var l_o_b_a local_order_book_analysis
			data_4.Scan(&l_o_b_a.id, &l_o_b_a.Time, &l_o_b_a.abs_ratio, &l_o_b_a.volum_ratio, &l_o_b_a.average_ratio, &l_o_b_a.varance_ratio)
			if l_o_b_a.id == 0 {
				signal[4] = true
				if num_4 != 0 {
					i_b.orderbook_abs_ratio /= float64(num_4)
					i_b.orderbook_volum_ratio /= float64(num_4)
					i_b.orderbook_average_ratio /= float64(num_4)
					i_b.orderbook_varance_ratio /= float64(num_4)
				}
				break
			}
			if l_o_b_a.Time > start_time_now+int64(F.Time_period*1000) {
				if num_4 != 0 {
					i_b.orderbook_abs_ratio /= float64(num_4)
					i_b.orderbook_volum_ratio /= float64(num_4)
					i_b.orderbook_average_ratio /= float64(num_4)
					i_b.orderbook_varance_ratio /= float64(num_4)
				}
				break
			} else {
				id_4++
				num_4++
				i_b.orderbook_abs_ratio += l_o_b_a.average_ratio
				i_b.orderbook_volum_ratio += l_o_b_a.volum_ratio
				i_b.orderbook_average_ratio += l_o_b_a.average_ratio
				i_b.orderbook_varance_ratio += l_o_b_a.varance_ratio
			}
		}
		stmt_insert.Exec(i_b.T_start, i_b.Ap_aver, i_b.Aq_aver, i_b.Bp_aver, i_b.Bq_aver, i_b.P_agg_aver, i_b.Q_agg_sum, i_b.PQ_agg_sum, i_b.Counts_sum, i_b.Counts_m1_sum, i_b.M_PQ_proportion, i_b.R_aver, i_b.orderbook_abs_ratio, i_b.orderbook_volum_ratio, i_b.orderbook_average_ratio, i_b.orderbook_varance_ratio)
		start_time_now += int64(F.Time_period * 1000)
		answer := true
		for i := 1; i < 5; i++ {
			answer = signal[i] && answer
		}
		if answer {
			break
		}
	}
	fmt.Println("特征提取完成")

}

func (F *Feature_extraction) Get_all_info() {
	F.init()
	sql_1 := "select id,T1,Ap,Aq,Bp,Bq from " + F.File_name + "best_offer_price" + " where id=?"
	sql_2 := "select id,T1,P,Q,Counts,M from " + F.File_name + "agg_trade" + " where id=?"
	sql_3 := "select id,T1,R from " + F.File_name + "newest_marked_price" + " where id=?"
	sql_4 := "select id,Time,abs_ratio,volum_ratio,average_ratio,varance_ratio from " + F.File_name + "_local_order_book_analysis" + " where id=?"
	stmt_best_offer_price, err := F.processed_data_db.Prepare(sql_1)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_best_offer_price.Close()
	stmt_agg_trade, err := F.processed_data_db.Prepare(sql_2)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_agg_trade.Close()
	stmt_newest_marked_price, err := F.processed_data_db.Prepare(sql_3)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt_newest_marked_price.Close()
	stmt__local_order_book_analysis, err := F.processed_data_db.Prepare(sql_4)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	defer stmt__local_order_book_analysis.Close()

	F.get_all_info_sub(stmt_best_offer_price, stmt_agg_trade, stmt_newest_marked_price, stmt__local_order_book_analysis)

}

// func main() {
// 	fmt.Println("hello world!")
// 	f_e := Feature_extraction{File_name: "lunausdt2021_12_30_17h45m34s", Username: "root", Password: "", Processed_data: "processed_data", Time_start_point: 1640857534692, Time_period: 1}
// 	f_e.get_all_info()
// }
