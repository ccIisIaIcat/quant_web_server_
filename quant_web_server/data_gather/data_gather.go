package data_gather

import (
	"database/sql"
	"fmt"
	"quant_web_server/get_info"
	"quant_web_server/local_order_book_2"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Data_gather struct {
	Symbol             string      //币种(必填)
	Original_data      string      //存入数据库路径(必填)
	Username           string      //数据库用户名(必填)
	Password           string      //数据库密码(必填)
	Order_book_depth   int         //订单簿深度
	chan_length        int         //channel 长度
	Num_wait           int         //介于本地订单簿需要一定时间才能稳定，设置一个等待本地订单簿的更新次数(必填)
	Add_depth          chan []byte //深度增量
	Best_pending_price chan []byte //最优挂单价格
	Newest_mark_peice  chan []byte //最新标记价格
	Agg_Trade          chan []byte //最新归集交易
	Local_book         chan []byte //本地订单簿
	command_list       [][]string  //存放拼接完成的ws字符串pa，对应顺序是，深度增量，最优挂单价格，最新标记价格，最新归集交易
	db                 *sql.DB     //mysql数据库指针
	Long               int         //检测时长的总秒数(必填)
	//一些进程结构体名称
	ask_1   get_info.Ask_a_wb
	ask_2   get_info.Ask_a_wb
	ask_3   get_info.Ask_a_wb
	ask_4   get_info.Ask_a_wb
	lo_bool local_order_book_2.My_orderbook
	label_1 string
}

func (D *Data_gather) init() {
	D.chan_length = 10
	D.command_list = make([][]string, 4)
	//完成对应字符串拼接：1、深度增量
	D.command_list[0] = []string{D.Symbol + "@depth@100ms"}
	//2、最优挂单价格
	D.command_list[1] = []string{D.Symbol + "@bookTicker"}
	//3、最新标记价格
	D.command_list[2] = []string{D.Symbol + "@markPrice"}
	//4、最新归集交易
	D.command_list[3] = []string{D.Symbol + "@aggTrade"}
	//初始化五个channel
	D.Add_depth = make(chan []byte, D.chan_length)
	D.Best_pending_price = make(chan []byte, D.chan_length)
	D.Newest_mark_peice = make(chan []byte, D.chan_length)
	D.Agg_Trade = make(chan []byte, D.chan_length)
	D.Local_book = make(chan []byte, D.chan_length)
	//初始化数据库db
	dsn := D.Username + ":" + D.Password + "@tcp(127.0.0.1:3306)/" + D.Original_data
	var err error
	D.db, err = sql.Open("mysql", dsn) //defer db.Close() // 注意这行代码要写在上面err判断的下面
	if err != nil {
		fmt.Println("mysql建立链接出错：", err)
		return
	}
	err = D.db.Ping()
	if err != nil {
		fmt.Println("mysql建立链接出错：")
		panic(err)
	}
	fmt.Println("mysql连接成功！")

}

//开启本地订单簿服务
func (D *Data_gather) strat_local_order() {
	D.lo_bool = local_order_book_2.My_orderbook{Symbol: D.Symbol, Outputer: &D.Local_book, Max_length: D.Order_book_depth}
	go D.lo_bool.Start_serve()
}

//开启另外四个服务
func (D *Data_gather) start_all_servers() {
	D.ask_1 = get_info.Ask_a_wb{Info_type: get_info.Symbol_Depth_addition, Params: D.command_list[0]}
	D.ask_2 = get_info.Ask_a_wb{Info_type: get_info.Symbol_BookTicker, Params: D.command_list[1]}
	D.ask_3 = get_info.Ask_a_wb{Info_type: get_info.Mark_Price, Params: D.command_list[2]}
	D.ask_4 = get_info.Ask_a_wb{Info_type: get_info.Agg_Trade, Params: D.command_list[3]}
	go D.ask_1.Start_wb(&D.Add_depth)
	go D.ask_2.Start_wb(&D.Best_pending_price)
	go D.ask_3.Start_wb(&D.Newest_mark_peice)
	go D.ask_4.Start_wb(&D.Agg_Trade)
}

//在数据库中建立两个个以当前时间特征为名的table
func (D *Data_gather) creat_new_table() (string, string) {
	aa := time.Now().Format("2006_01_02_15h04m05s")
	label_1 := D.Symbol + aa
	D.label_1 = label_1
	sql_l := "CREATE TABLE " + label_1 + "(id int PRIMARY KEY AUTO_INCREMENT, context Blob)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	label_2 := D.Symbol + aa + "_local_order_book"
	sql_2 := "CREATE TABLE " + label_2 + "(id int PRIMARY KEY AUTO_INCREMENT, context Blob)" + "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	stmt, err := D.db.Prepare(sql_l)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Printf("User Table successfully ....\n")
	}
	stmt2, err := D.db.Prepare(sql_2)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt2.Exec()
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Printf("User Table successfully ....\n")
	}
	return label_1, label_2
}

//在数据库中插入数据insert context
func (D *Data_gather) insert_data(data []byte, label_new string) {
	sqlStr := "insert into " + label_new + " (context) values(?)"
	_, err := D.db.Exec(sqlStr, data)
	if err != nil {
		fmt.Println("数据导入失败！", err)
	}

}

func (D *Data_gather) Start() string {
	D.init()
	D.strat_local_order()
	for i := 0; i < D.Num_wait; i++ {
		data_local := <-D.Local_book
		fmt.Println("已更新次数", i+1, "订单簿长度：", len(data_local))
	}
	D.start_all_servers()
	var wg sync.WaitGroup
	wg.Add(1)
	l1, l2 := D.creat_new_table()
	go func() {
		for i := 0; i < D.Long; i++ {
			time.Sleep(time.Second)
			fmt.Println("录入进度：", i+1, "/", D.Long)
		}
		fmt.Println("录入编号：", l1)
		D.ask_1.End_Conn()
		D.ask_2.End_Conn()
		D.ask_3.End_Conn()
		D.ask_4.End_Conn()
		D.lo_bool.End_local_order_book()
		wg.Done()
		fmt.Println("录入完成")
	}()
	go func() {
		for {
			data_local := <-D.Local_book
			D.insert_data(data_local, l2)
		}
	}()
	go func() {
		for {
			data_2 := <-D.Add_depth
			D.insert_data(data_2, l1)
		}
	}()
	go func() {
		for {
			data_3 := <-D.Best_pending_price
			D.insert_data(data_3, l1)
		}
	}()
	go func() {
		for {
			data_4 := <-D.Newest_mark_peice
			D.insert_data(data_4, l1)
		}
	}()
	go func() {
		for {
			data_5 := <-D.Agg_Trade
			D.insert_data(data_5, l1)
		}
	}()
	wg.Wait()
	D.Data_process()
	return D.label_1
}

func (D *Data_gather) Data_process() {
	// data_process_1.Process_data(D.label_1)
	// data_process_2_3.Process_data_2(D.label_1)
}
