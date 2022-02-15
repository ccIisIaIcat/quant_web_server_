package goal_feature

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Goal_feature struct {
	File_name                string //必填
	Username                 string //必填
	Password                 string //必填
	Processed_data           string
	processed_data_db        *sql.DB
	Cross_feature_data       string
	cross_feature_data_db    *sql.DB
	Time_start_point         int64 //起始点的时间点(必填)
	Goal_feature_time_period int   //时间区间的秒数(必填)
}

type best_offer_price struct {
	id  int
	T_1 int64
	Ap  float64
	Aq  float64
	Bp  float64
	Bq  float64
}

func (G *Goal_feature) init() {
	dsn := G.Username + ":" + G.Password + "@tcp(127.0.0.1:3306)/" + G.Processed_data
	dsn_2 := G.Username + ":" + G.Password + "@tcp(127.0.0.1:3306)/" + G.Cross_feature_data
	var err error
	G.processed_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = G.processed_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	G.cross_feature_data_db, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = G.processed_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	G.Time_start_point = get_time_start(G.File_name)
}

func (G *Goal_feature) Get_all_info() {
	G.init()
	sql_1 := "select id,T1,Ap,Aq,Bp,Bq from " + G.File_name + "best_offer_price" + " where id=?"
	stmt_best_offer_price, err := G.processed_data_db.Prepare(sql_1)
	if err != nil {
		fmt.Println("ERR!", err)
	}
	G.get_all_info_sub(stmt_best_offer_price)
	defer stmt_best_offer_price.Close()

}

func (G *Goal_feature) creat_table() string {
	label := G.File_name + "goal_feature"
	period_str := strconv.Itoa(G.Goal_feature_time_period)
	label = label + "_p" + period_str
	sql := "CREATE TABLE " + label + "(id int PRIMARY KEY AUTO_INCREMENT,T_start bigint,bids_start double,asks_start double,bids_end double,asks_end double,b_b_hat double,b_r_2 double,a_b_hat double,a_r_2 double,max_change int,positive_change int)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := G.cross_feature_data_db.Prepare(sql)
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

func (G *Goal_feature) get_all_info_sub(stmt_best_offer_price *sql.Stmt) {
	new_label := G.creat_table()
	sql_insert := "insert into " + new_label + " (T_start,bids_start,asks_start,bids_end,asks_end,b_b_hat,b_r_2,a_b_hat,a_r_2) values(?,?,?,?,?,?,?,?,?);"
	stmt_insert, err := G.cross_feature_data_db.Prepare(sql_insert)
	if err != nil {
		fmt.Println("stmt_insert错误：", err)
	}
	defer stmt_insert.Close()
	id := 1
	start_time_now := G.Time_start_point
	for {
		bids_start := float64(0)
		asks_start := float64(0)
		bids_end := float64(0)
		asks_end := float64(0)
		b_b_hat := float64(0)
		a_b_hat := float64(0)
		b_r_2 := float64(0)
		a_r_2 := float64(0)
		x_list := make([]int64, 0)
		y_list_asks := make([]float64, 0)
		y_list_bids := make([]float64, 0)
		signal := false
		num := 0
		for {
			data := stmt_best_offer_price.QueryRow(id)
			var b_o_p best_offer_price
			data.Scan(&b_o_p.id, &b_o_p.T_1, &b_o_p.Ap, &b_o_p.Aq, &b_o_p.Bp, &b_o_p.Bq)
			if num == 0 {
				bids_start = b_o_p.Bp
				asks_start = b_o_p.Ap
			}
			if b_o_p.id == 0 {
				signal = true
				if num != 0 {
					bids_end = b_o_p.Bp
					asks_end = b_o_p.Ap
					x_list = append(x_list, b_o_p.T_1)
					y_list_asks = append(y_list_asks, b_o_p.Ap)
					y_list_bids = append(y_list_bids, b_o_p.Bp)
					b_b_hat, b_r_2 = calculate_b_r(x_list, y_list_bids)
					a_b_hat, a_r_2 = calculate_b_r(x_list, y_list_asks)
				}
				break
			}
			if b_o_p.T_1 > start_time_now+int64(G.Goal_feature_time_period*1000) {
				if num != 0 {
					bids_end = b_o_p.Bp
					asks_end = b_o_p.Ap
					x_list = append(x_list, b_o_p.T_1)
					y_list_asks = append(y_list_asks, b_o_p.Ap)
					y_list_bids = append(y_list_bids, b_o_p.Bp)
					b_b_hat, b_r_2 = calculate_b_r(x_list, y_list_bids)
					a_b_hat, a_r_2 = calculate_b_r(x_list, y_list_asks)
				}
				break
			} else {
				id++
				num++
				x_list = append(x_list, b_o_p.T_1)
				y_list_asks = append(y_list_asks, b_o_p.Ap)
				y_list_bids = append(y_list_bids, b_o_p.Bp)
			}

		}
		stmt_insert.Exec(start_time_now, bids_start, asks_start, bids_end, asks_end, b_b_hat, b_r_2, a_b_hat, a_r_2)
		start_time_now += int64(G.Goal_feature_time_period * 1000)
		if signal {
			break
		}

	}

	fmt.Println("目标变量提取完成完成")

}

func get_time_start(tittle string) int64 {
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

//给定两列数组，计算线性回归b估和r方
func calculate_b_r(nums_1 []int64, nums_2 []float64) (float64, float64) {
	if len(nums_1) == 0 {
		return 0, 0
	}
	sum_num := 0
	x_list := make([]float64, 0)
	y_list := make([]float64, 0)
	for i := 0; i < len(nums_2); i++ {
		if nums_2[i] != float64(0) {
			sum_num++
			x_list = append(x_list, float64(nums_1[i]-nums_1[0]))
			y_list = append(y_list, nums_2[i])
		}
	}

	l_xx := float64(0)
	l_xy := float64(0)
	l_yy := float64(0)
	x_sum := float64(0)
	y_sum := float64(0)
	for i := 0; i < sum_num; i++ {
		x_sum += x_list[i]
		y_sum += y_list[i]
	}
	for i := 0; i < sum_num; i++ {
		l_xx += x_list[i] * x_list[i]
		l_xy += x_list[i] * y_list[i]
		l_yy += y_list[i] * y_list[i]
	}
	l_xx = l_xx - x_sum*x_sum/float64(sum_num)
	l_xy = l_xy - x_sum*y_sum/float64(sum_num)
	l_yy = l_yy - y_sum*y_sum/float64(sum_num)
	b_hat := float64(0)
	if l_xx != 0 {
		b_hat = l_xy / l_xx
	}
	R_2 := float64(0)
	if l_xx != 0 && l_yy != 0 {
		R_2 = l_xy * l_xy / l_xx / l_yy
	}

	return b_hat, R_2

}

// func main() {
// 	fmt.Println("hello world!")
// 	g_f := Goal_feature{File_name: "lunausdt2021_12_30_23h26m51s", Username: "root", Password: "", Processed_data: "processed_data", Time_start_point: 1640857534692, Goal_feature_time_period: 1}
// 	g_f.Get_all_info()
// }
