package data_process_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Data_process struct {
	File_name               string         //所收集数据的文件名
	Username                string         //数据库用户名
	Password                string         //密码
	Original_data           string         //原始数据数据库
	Processed_data          string         //处理数据数据库
	db_1                    *sql.DB        //读取文件的sql对象
	db_2                    *sql.DB        //存入数据的sql对象
	channel_read_total      chan []byte    //用于读取非订单簿全部类型数据的chan
	channel_read_mark_price chan []byte    //用于从channel_read_total中分流出最新标记价格的信息
	channel_read_best_offer chan []byte    //用于从channel_read_total中分流出最优订单信息
	channel_read_agg_trade  chan []byte    //用于从channel_read_total中分流出聚合交易信息
	wg                      sync.WaitGroup //用于优雅地停止进程
	max_id                  int            //用于输出信息
}

type judge_tool struct {
	The_type string `json:"e"`
}

type recevier struct {
	id      int
	context []byte
}

type time_ struct {
	T_send int64 `json:"E"`
	T_make int64 `json:"T"`
}

func (D *Data_process) init() {
	D.wg.Add(1)
	D.channel_read_total = make(chan []byte, 10)
	D.channel_read_mark_price = make(chan []byte, 10)
	D.channel_read_best_offer = make(chan []byte, 10)
	D.channel_read_agg_trade = make(chan []byte, 10)
	dsn := D.Username + ":" + D.Password + "@tcp(127.0.0.1:3306)/" + D.Original_data
	dsn_2 := D.Username + ":" + D.Password + "@tcp(127.0.0.1:3306)/" + D.Processed_data
	var err error
	D.db_1, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("db_1格式错误：", err)
		return
	}
	D.db_2, err = sql.Open("mysql", dsn_2)
	if err != nil {
		fmt.Println("db_2格式错误：", err)
		return
	}
	err = D.db_1.Ping()
	if err != nil {
		fmt.Println("db_1建立链接出错：")
		panic(err)
	}
	fmt.Println("db_1连接成功！")
	err = D.db_2.Ping()
	if err != nil {
		fmt.Println("db_2建立链接出错：")
		panic(err)
	}
	fmt.Println("db_2连接成功！")

}

//建立最优买卖单和最新报价的sql表
func (D *Data_process) creat_table() (string, string, string) {
	defer D.wg.Done()
	label_1 := D.File_name + "best_offer_price"
	label_2 := D.File_name + "newest_marked_price"
	label_3 := D.File_name + "agg_trade"
	sql_1 := "CREATE TABLE " + label_1 + "(id int PRIMARY KEY AUTO_INCREMENT,T1 bigint,T2 bigint,Bp double,Bq double,Ap double,Aq double)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	sql_2 := "CREATE TABLE " + label_2 + "(id int PRIMARY KEY AUTO_INCREMENT,T1 bigint,T2 bigint,P1 double,P2 Double,P3 Double,R Double)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	sql_3 := "CREATE TABLE " + label_3 + "(id int PRIMARY KEY AUTO_INCREMENT,T1 bigint,T2 bigint,P double,Q Double,Counts int,M bool)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := D.db_2.Prepare(sql_1)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println(label_1 + "建表成功")
	}
	stmt2, err := D.db_2.Prepare(sql_2)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt2.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println(label_2 + "建表成功")
	}
	stmt3, err := D.db_2.Prepare(sql_3)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt3.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println(label_3 + "建表成功")
	}
	return label_1, label_2, label_3
}

//从sql中读取数据并分发给对应的channel
func (D *Data_process) send_message() {
	sql_query := "select id,context from " + D.File_name + " where id=?"
	stmt, err := D.db_1.Prepare(sql_query)
	if err != nil {
		fmt.Println("stmt错误：", err)
	}
	i := 1
	for {
		rowOJB := stmt.QueryRow(i)
		var recevier recevier
		rowOJB.Scan(&recevier.id, &recevier.context)
		if recevier.id == 0 {
			for {
				if len(D.channel_read_best_offer) == 0 && len(D.channel_read_mark_price) == 0 && len(D.channel_read_agg_trade) == 0 {
					time.Sleep(time.Second)
					break
				}
			}
			D.wg.Done()
			break
		}
		var tool judge_tool
		json.Unmarshal(recevier.context, &tool)
		if string(tool.The_type) == "bookTicker" {
			D.channel_read_best_offer <- recevier.context
		} else if string(tool.The_type) == "markPriceUpdate" {
			D.channel_read_mark_price <- recevier.context
		} else if string(tool.The_type) == "aggTrade" {
			D.channel_read_agg_trade <- recevier.context
		}
		i++
	}

}

//处理channel_read_best_offer里的数据(结构体)
type best_offer_price struct {
	T1 int64  `json:"E"`
	T2 int64  `json:"T"`
	Ap string `json:"a"`
	Aq string `json:"A"`
	Bp string `json:"b"`
	Bq string `json:"B"`
}

//处理channel_read_best_offer里的数据(操作)
func (D *Data_process) save_best_offer(label_1 string) {
	SQL := "insert into " + label_1 + " (T1,T2,Ap,Aq,Bp,Bq) values(?,?,?,?,?,?);"
	stmt, err := D.db_2.Prepare(SQL)
	if err != nil {
		fmt.Println("stmt错误", err)
	}
	defer stmt.Close()
	for {
		var bop best_offer_price
		data := <-D.channel_read_best_offer

		json.Unmarshal(data, &bop)
		ap, _ := strconv.ParseFloat(bop.Ap, 64)
		aq, _ := strconv.ParseFloat(bop.Aq, 64)
		bp, _ := strconv.ParseFloat(bop.Bp, 64)
		bq, _ := strconv.ParseFloat(bop.Bq, 64)
		stmt.Exec(bop.T1, bop.T2, ap, aq, bp, bq)
	}

}

//处理channel_read_mark_price里的数据(结构体)
type mark_price struct {
	T1 int64  `json:"E"` //事件时间
	T2 int64  `json:"T"` //下次资金时间
	P1 string `json:"p"` //标记价格
	P2 string `json:"i"` //现货指数价格
	P3 string `json:"P"` //预估结算价格
	R  string `json:"r"` //资金费率
}

//处理channel_read_mark_price里的数据(操作)
func (D *Data_process) save_mark_price(label_2 string) {
	SQL := "insert into " + label_2 + " (T1,T2,P1,P2,P3,R) values(?,?,?,?,?,?)"
	stmt, err := D.db_2.Prepare(SQL)
	if err != nil {
		fmt.Println("stmt错误：", err)
	}
	defer stmt.Close()
	for {

		var mp mark_price
		data := <-D.channel_read_mark_price

		json.Unmarshal(data, &mp)
		p1, _ := strconv.ParseFloat(mp.P1, 64)
		p2, _ := strconv.ParseFloat(mp.P2, 64)
		p3, _ := strconv.ParseFloat(mp.P3, 64)
		r, _ := strconv.ParseFloat(mp.R, 64)
		stmt.Exec(mp.T1, mp.T2, p1, p2, p3, r)
	}

}

//处理channel_read_agg_trade里的数据(结构体)
type agg_trade struct {
	T1 int64  `json:"E"` //事件时间
	T2 int64  `json:"T"` //成交时间
	P  string `json:"p"` //成交价格
	Q  string `json:"q"` //成交量
	C1 int    `json:"f"` //首个交易id
	C2 int    `json:"l"` //最后交易id
	M  bool   `json:"m"` //买方还是卖方
}

//处理channel_read_agg_trade里的数据(操作)
func (D *Data_process) save_agg_trade(label_3 string) {
	SQL := "insert into " + label_3 + " (T1,T2,P,Q,Counts,M) values(?,?,?,?,?,?)"
	stmt, err := D.db_2.Prepare(SQL)
	if err != nil {
		fmt.Println("stmt错误：", err)
	}
	defer stmt.Close()
	for {
		var at agg_trade
		data := <-D.channel_read_agg_trade
		json.Unmarshal(data, &at)
		p, _ := strconv.ParseFloat(at.P, 64)
		q, _ := strconv.ParseFloat(at.Q, 64)
		m2 := 1
		if at.M {
			m2 = 1
		} else {
			m2 = 0
		}
		stmt.Exec(at.T1, at.T2, p, q, (at.C2-at.C1)+1, m2)
	}

}

func (D *Data_process) Process_data() {
	D.init()
	l1, l2, l3 := D.creat_table()
	D.wg.Wait()
	D.wg.Add(1)
	go D.send_message()
	go D.save_best_offer(l1)
	go D.save_mark_price(l2)
	go D.save_agg_trade(l3)
	D.wg.Wait()
	fmt.Println("数据录入完成(dataprocess_1)")
	fmt.Println("数据库：data_processed")
	fmt.Println("表一：", l1)
	fmt.Println("表二：", l2)
	fmt.Println("表三：", l3)
}
