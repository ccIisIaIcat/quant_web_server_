package main

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"sort"

	_ "github.com/go-sql-driver/mysql"
)

type Feature_matrix struct {
	File_name          string //必填
	Username           string //必填
	Password           string //必填
	K_value_data       string //必填
	K_value_data_db    *sql.DB
	label              string      //新建表单名称
	dividied_point     []float64   //分位点
	data_long_matrix   [][]float64 //用于存储全部数据
	period_list        []int       //周期表单
	info_sql           string      //用于创建建表sql
	table_info         []string    //用于存放表列的名称
	classify_point_set [][]float64 //不同列的划分点
}

func (F *Feature_matrix) init() {
	F.classify_point_set = make([][]float64, 52)
	F.data_long_matrix = make([][]float64, 52)
	if len(F.dividied_point) == 0 {
		F.dividied_point = []float64{0.2, 0.4, 0.6, 0.8}
	}
	if len(F.period_list) == 0 {
		F.period_list = []int{3, 5, 10, 20}
	}
	dsn := F.Username + ":" + F.Password + "@tcp(127.0.0.1:3306)/" + F.K_value_data
	var err error
	F.K_value_data_db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db格式错误：", err)
		return
	}
	err = F.K_value_data_db.Ping()
	if err != nil {
		fmt.Println("db建立链接出错：")
		panic(err)
	}
	F.info_sql = fmt.Sprintf("(id int PRIMARY KEY AUTO_INCREMENT,T_start bigint,q_agg_aver_p%v int,q_agg_aver_p%v int,q_agg_aver_p%v int,q_agg_aver_p%v int,pq_agg_sum_p%v int,pq_agg_sum_p%v int,pq_agg_sum_p%v int,pq_agg_sum_p%v int,price_change_p%v int,price_change_p%v int,price_change_p%v int,price_change_p%v int,counts_sum_p%v int,counts_sum_p%v int,counts_sum_p%v int,counts_sum_p%v int,counts_m1_sum_p%v int,counts_m1_sum_p%v int,counts_m1_sum_p%v int,counts_m1_sum_p%v int,counts_m1_proportion_p%v int,counts_m1_proportion_p%v int,counts_m1_proportion_p%v int,counts_m1_proportion_p%v int,m_pq_sum_p%v int,m_pq_sum_p%v int,m_pq_sum_p%v int,m_pq_sum_p%v int,m_pq_proportion_p%v int,m_pq_proportion_p%v int,m_pq_proportion_p%v int,m_pq_proportion_p%v int,R_aver_p%v int,R_aver_p%v int,R_aver_p%v int,R_aver_p%v int,orderbook_abs_ratio_p%v int,orderbook_abs_ratio_p%v int,orderbook_abs_ratio_p%v int,orderbook_abs_ratio_p%v int,orderbook_volum_ratio_p%v int,orderbook_volum_ratio_p%v int,orderbook_volum_ratio_p%v int,orderbook_volum_ratio_p%v int,orderbook_average_ratio_p%v int,orderbook_average_ratio_p%v int,orderbook_average_ratio_p%v int,orderbook_average_ratio_p%v int,orderbook_varance_ratio_p%v int,orderbook_varance_ratio_p%v int,orderbook_varance_ratio_p%v int,orderbook_varance_ratio_p%v int,goal_feature int)", F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3], F.period_list[0], F.period_list[1], F.period_list[2], F.period_list[3])
	F.table_info = F.make_table_info()
}

func (F *Feature_matrix) make_table_info() []string {
	answer := make([]string, 0)
	for i := 0; i < len(F.period_list); i++ {
		answer = append(answer, fmt.Sprintf("price_change_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("q_agg_aver_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("pq_agg_sum_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("counts_sum_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("counts_m1_sum_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("counts_m1_proportion_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("m_pq_sum_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("m_pq_proportion_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("R_aver_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_abs_ratio_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_volum_ratio_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_average_ratio_p%v", F.period_list[i]))
		answer = append(answer, fmt.Sprintf("orderbook_varance_ratio_p%v", F.period_list[i]))
	}
	return answer
}

func (F *Feature_matrix) creat_table() {
	label := F.File_name + "feature_matrix"
	sql := "create table " + label + F.info_sql + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"

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
}

type all_info struct {
	price_change            float64
	q_agg_aver              float64
	pq_agg_sum              float64
	counts_sum              int
	counts_m1_sum           int
	counts_m1_proportion    float64
	m_pq_sum                float64
	m_pq_proportion         float64
	R_aver                  float64
	orderbook_abs_ratio     float64
	orderbook_volum_ratio   float64
	orderbook_average_ratio float64
	orderbook_varance_ratio float64
}

func (F *Feature_matrix) get_all_info_by_order(suq int) {
	sql_query := "select "
	for i := suq * 13; i < (suq+1)*13-1; i++ {
		sql_query += F.table_info[i]
		sql_query += ","
	}
	sql_query += F.table_info[(suq+1)*13-1]
	sql_query += " from " + F.File_name + "data_matrix"
	data, err := F.K_value_data_db.Query(sql_query)
	if err != nil {
		fmt.Println("err", err)
	}
	var a_l all_info
	for data.Next() {
		data.Scan(&a_l.price_change, &a_l.q_agg_aver, &a_l.pq_agg_sum, &a_l.counts_sum, &a_l.counts_m1_sum, &a_l.counts_m1_proportion, &a_l.m_pq_sum, &a_l.m_pq_proportion, &a_l.R_aver, &a_l.orderbook_abs_ratio, &a_l.orderbook_volum_ratio, &a_l.orderbook_average_ratio, &a_l.orderbook_varance_ratio)
		F.data_long_matrix[13*suq+0] = append(F.data_long_matrix[13*suq+0], math.Abs(a_l.price_change))
		F.data_long_matrix[13*suq+1] = append(F.data_long_matrix[13*suq+1], a_l.q_agg_aver)
		F.data_long_matrix[13*suq+2] = append(F.data_long_matrix[13*suq+2], a_l.pq_agg_sum)
		F.data_long_matrix[13*suq+3] = append(F.data_long_matrix[13*suq+3], float64(a_l.counts_sum))
		F.data_long_matrix[13*suq+4] = append(F.data_long_matrix[13*suq+4], float64(a_l.counts_m1_sum))
		F.data_long_matrix[13*suq+5] = append(F.data_long_matrix[13*suq+5], a_l.counts_m1_proportion)
		F.data_long_matrix[13*suq+6] = append(F.data_long_matrix[13*suq+6], a_l.m_pq_sum)
		F.data_long_matrix[13*suq+7] = append(F.data_long_matrix[13*suq+7], a_l.m_pq_proportion)
		F.data_long_matrix[13*suq+8] = append(F.data_long_matrix[13*suq+8], a_l.R_aver)
		F.data_long_matrix[13*suq+9] = append(F.data_long_matrix[13*suq+9], math.Abs(a_l.orderbook_abs_ratio-1))
		F.data_long_matrix[13*suq+10] = append(F.data_long_matrix[13*suq+10], math.Abs(a_l.orderbook_volum_ratio-1))
		F.data_long_matrix[13*suq+11] = append(F.data_long_matrix[13*suq+11], math.Abs(a_l.orderbook_average_ratio-1))
		F.data_long_matrix[13*suq+12] = append(F.data_long_matrix[13*suq+12], math.Abs(a_l.orderbook_varance_ratio-1))
	}
}

func (F *Feature_matrix) get_divide_matrix() {
	for i := 0; i < 4; i++ {
		F.get_all_info_by_order(i)
	}

	for i := 0; i < 52; i++ {
		a, b, c, d := F.divide_list_float(F.data_long_matrix[i])
		F.classify_point_set[i] = append(F.classify_point_set[i], a)
		F.classify_point_set[i] = append(F.classify_point_set[i], b)
		F.classify_point_set[i] = append(F.classify_point_set[i], c)
		F.classify_point_set[i] = append(F.classify_point_set[i], d)
	}

}

func (F *Feature_matrix) divide_list_float(sample_list []float64) (float64, float64, float64, float64) {
	sort.Float64s(sort.Float64Slice(sample_list))
	a := sample_list[int(math.Floor(float64(len(sample_list))*F.dividied_point[0]))]
	b := sample_list[int(math.Floor(float64(len(sample_list))*F.dividied_point[1]))]
	c := sample_list[int(math.Floor(float64(len(sample_list))*F.dividied_point[2]))]
	d := sample_list[int(math.Floor(float64(len(sample_list))*F.dividied_point[3]))]
	return a, b, c, d
}

type all_info_with_order struct {
	price_change_1            float64
	q_agg_aver_1              float64
	pq_agg_sum_1              float64
	counts_sum_1              int
	counts_m1_sum_1           int
	counts_m1_proportion_1    float64
	m_pq_sum_1                float64
	m_pq_proportion_1         float64
	R_aver_1                  float64
	orderbook_abs_ratio_1     float64
	orderbook_volum_ratio_1   float64
	orderbook_average_ratio_1 float64
	orderbook_varance_ratio_1 float64
	price_change_2            float64
	q_agg_aver_2              float64
	pq_agg_sum_2              float64
	counts_sum_2              int
	counts_m1_sum_2           int
	counts_m1_proportion_2    float64
	m_pq_sum_2                float64
	m_pq_proportion_2         float64
	R_aver_2                  float64
	orderbook_abs_ratio_2     float64
	orderbook_volum_ratio_2   float64
	orderbook_average_ratio_2 float64
	orderbook_varance_ratio_2 float64
	price_change_3            float64
	q_agg_aver_3              float64
	pq_agg_sum_3              float64
	counts_sum_3              int
	counts_m1_sum_3           int
	counts_m1_proportion_3    float64
	m_pq_sum_3                float64
	m_pq_proportion_3         float64
	R_aver_3                  float64
	orderbook_abs_ratio_3     float64
	orderbook_volum_ratio_3   float64
	orderbook_average_ratio_3 float64
	orderbook_varance_ratio_3 float64
	price_change_4            float64
	q_agg_aver_4              float64
	pq_agg_sum_4              float64
	counts_sum_4              int
	counts_m1_sum_4           int
	counts_m1_proportion_4    float64
	m_pq_sum_4                float64
	m_pq_proportion_4         float64
	R_aver_4                  float64
	orderbook_abs_ratio_4     float64
	orderbook_volum_ratio_4   float64
	orderbook_average_ratio_4 float64
	orderbook_varance_ratio_4 float64
}

func (F *Feature_matrix) make_new_query_sql() string {
	sql := "select "
	for suq := 0; suq < 4; suq++ {
		for i := suq * 13; i < (suq+1)*13-1; i++ {
			sql += F.table_info[i]
			sql += ","
		}
		sql += F.table_info[(suq+1)*13-1]
		if suq != 3 {
			sql += ","
		}
	}
	sql += " from " + F.File_name + "data_matrix"

	return sql

}

func (F *Feature_matrix) save_classification_info() {
	sql_insert := "insert into " + F.label + " ("
	for suq := 0; suq < 4; suq++ {
		for i := suq * 13; i < (suq+1)*13-1; i++ {
			sql_insert += F.table_info[i]
			sql_insert += ","
		}
		sql_insert += F.table_info[(suq+1)*13-1]
		if suq != 3 {
			sql_insert += ","
		}
	}
	sql_insert += ") values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	stmt_insert, err := F.K_value_data_db.Prepare(sql_insert)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt_insert.Close()

	data_list := make([]float64, 52)
	new_query_sql := F.make_new_query_sql()
	data, err := F.K_value_data_db.Query(new_query_sql)
	if err != nil {
		fmt.Println(err)
	}
	var a_i_w_o all_info_with_order
	num := 1
	for data.Next() {
		fmt.Println(num)
		num++
		data.Scan(&a_i_w_o.price_change_1, &a_i_w_o.q_agg_aver_1, &a_i_w_o.pq_agg_sum_1, &a_i_w_o.counts_sum_1, &a_i_w_o.counts_m1_sum_1, &a_i_w_o.counts_m1_proportion_1, &a_i_w_o.m_pq_sum_1, &a_i_w_o.m_pq_proportion_1, &a_i_w_o.R_aver_1, &a_i_w_o.orderbook_abs_ratio_1, &a_i_w_o.orderbook_volum_ratio_1, &a_i_w_o.orderbook_average_ratio_1, &a_i_w_o.orderbook_varance_ratio_1, &a_i_w_o.price_change_2, &a_i_w_o.q_agg_aver_2, &a_i_w_o.pq_agg_sum_2, &a_i_w_o.counts_sum_2, &a_i_w_o.counts_m1_sum_2, &a_i_w_o.counts_m1_proportion_2, &a_i_w_o.m_pq_sum_2, &a_i_w_o.m_pq_proportion_2, &a_i_w_o.R_aver_2, &a_i_w_o.orderbook_abs_ratio_2, &a_i_w_o.orderbook_volum_ratio_2, &a_i_w_o.orderbook_average_ratio_2, &a_i_w_o.orderbook_varance_ratio_2, &a_i_w_o.price_change_3, &a_i_w_o.q_agg_aver_3, &a_i_w_o.pq_agg_sum_3, &a_i_w_o.counts_sum_3, &a_i_w_o.counts_m1_sum_3, &a_i_w_o.counts_m1_proportion_3, &a_i_w_o.m_pq_sum_3, &a_i_w_o.m_pq_proportion_3, &a_i_w_o.R_aver_3, &a_i_w_o.orderbook_abs_ratio_3, &a_i_w_o.orderbook_volum_ratio_3, &a_i_w_o.orderbook_average_ratio_3, &a_i_w_o.orderbook_varance_ratio_3, &a_i_w_o.price_change_4, &a_i_w_o.q_agg_aver_4, &a_i_w_o.pq_agg_sum_4, &a_i_w_o.counts_sum_4, &a_i_w_o.counts_m1_sum_4, &a_i_w_o.counts_m1_proportion_4, &a_i_w_o.m_pq_sum_4, &a_i_w_o.m_pq_proportion_4, &a_i_w_o.R_aver_4, &a_i_w_o.orderbook_abs_ratio_4, &a_i_w_o.orderbook_volum_ratio_4, &a_i_w_o.orderbook_average_ratio_4, &a_i_w_o.orderbook_varance_ratio_4)
		for i := 0; i < reflect.ValueOf(a_i_w_o).NumField(); i++ {
			if i%13 == 3 || i%13 == 4 {
				data_list[i] = float64(reflect.ValueOf(a_i_w_o).Field(i).Int())
			} else if i%13 >= 9 {
				data_list[i] = math.Abs(reflect.ValueOf(a_i_w_o).Field(i).Float() - 1)
			} else {
				data_list[i] = math.Abs(reflect.ValueOf(a_i_w_o).Field(i).Float())
			}
		}
		values_list := make([]int, 52)
		F.make_divided_list(data_list, &values_list)
		stmt_insert.Exec(values_list[0], values_list[1], values_list[2], values_list[3], values_list[4], values_list[5], values_list[6], values_list[7], values_list[8], values_list[9], values_list[10], values_list[11], values_list[12], values_list[13], values_list[14], values_list[15], values_list[16], values_list[17], values_list[18], values_list[19], values_list[20], values_list[21], values_list[22], values_list[23], values_list[24], values_list[25], values_list[26], values_list[27], values_list[28], values_list[29], values_list[30], values_list[31], values_list[32], values_list[33], values_list[34], values_list[35], values_list[36], values_list[37], values_list[38], values_list[39], values_list[40], values_list[41], values_list[42], values_list[43], values_list[44], values_list[45], values_list[46], values_list[47], values_list[48], values_list[49], values_list[50], values_list[51])

	}

}

func (F *Feature_matrix) make_divided_list(sample_list []float64, null_list *[]int) {
	for i := 0; i < len(sample_list); i++ {
		for j := 0; j < len(F.classify_point_set[i]); j++ {
			if sample_list[i] > F.classify_point_set[i][j] {
				(*null_list)[i] = j + 2
			}
		}
		if (*null_list)[i] == 0 {
			(*null_list)[i] = 1
		}
	}
}

type t_update struct {
	id      int
	T_start int64
}

func (F *Feature_matrix) update_t_start() {
	fmt.Println("更新时间列表")
	sql_select := "select id,T_start from " + F.File_name + "data_matrix"
	fmt.Println(sql_select)
	data, err := F.K_value_data_db.Query(sql_select)
	if err != nil {
		fmt.Println(err)
	}
	var t_u t_update
	for data.Next() {
		data.Scan(&t_u.id, &t_u.T_start)
		sql_update := "update " + F.label + fmt.Sprintf(" set T_start=%v where id=%v", t_u.T_start, t_u.id)
		F.K_value_data_db.Exec(sql_update)
	}
	fmt.Println("更新完毕")
}

func (F *Feature_matrix) Start() {
	F.init()
	F.creat_table()
	F.get_divide_matrix()
	F.save_classification_info()
	// F.update_t_start()
}

func main() {
	f_m := Feature_matrix{File_name: "btcusdt2022_01_01_19h05m11s", Username: "root", Password: "", K_value_data: "classify_data"}
	f_m.Start()
}
