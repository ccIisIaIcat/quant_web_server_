package get_info

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

//申明一个关于请求的结构体
type Ask_a_wb struct {
	Info_type string
	Params    []string
	conn_this *websocket.Conn
	Judge     bool
}

//开始web服务
func (A *Ask_a_wb) Start_wb(info_gather *chan []byte) {
	A.Judge = false
	go A.get_information(A.Info_type, A.Params, info_gather)
}

//中止web服务
func (A *Ask_a_wb) End_wb() {
	unsubscribe_a_server(A.conn_this, A.Params)
}

//重启web服务
func (A *Ask_a_wb) Restart_wb() {
	subscribe_a_server(A.conn_this, A.Params)
}

//彻底关闭web服务
func (A *Ask_a_wb) End_Conn_Force() {
	fmt.Println("wb通道已强制关闭")
	A.conn_this.Close()
}

func (A *Ask_a_wb) End_Conn() {
	A.Judge = true
}

//建立一个请求头范式
type ask_websocket struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

//申明一些关键变量
var (
	//币安websocket主页面
	BASE_URL string = "wss://fstream.binance.com/ws/"
	//(归集交易）同一价格、同一方向、同一时间(100ms计算)的trade会被聚合为一条.
	Agg_Trade string = "<symbol>@aggTrade"
	//最新标记价格
	Mark_Price string = "<symbol>@markPrice"
	//全场最新标记价格!markPrice@arr
	All_Mark_Price string = "!markPrice@arr"
	//K线图（需要额外输入时间）（0.25秒一更新）
	K_Line string = "<symbol>@kline_<interval>"
	//按Symbol的精简Ticker(symbol的最新成交加，以及二十四小时内的最高价，最低价)(半秒一更新)
	Symbol_Mini_Ticker string = "<symbol>@miniTicker"
	//全市场的精简ticker
	All_Mini_Ticker string = "!miniTicker@arr"
	//按symbol的完整ticker（半秒一更新）
	Symbol_Ticker string = "<symbol>@ticker"
	//全市场的完整ticker(一秒一更新)
	All_Ticker string = "!ticker@arr"
	//按Symbol的最优挂单信息（实时）
	Symbol_BookTicker string = "<symbol>@bookTicker"
	//全场最优挂单信息
	All_BookTicker string = "!bookTicker"
	//有限深度信息（可选5/10/20档）
	Symbol_Depth string = "<symbol>@depth<levels>"
	//增量深度信息
	Symbol_Depth_addition string = "<symbol>@depth@100ms"
)

//通过提供的方法，将信息传到一个[]byte的channel中
func (A *Ask_a_wb) get_information(info_type string, params []string, info_gather *chan []byte) {
	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//生成访问地址
	wss := BASE_URL + Symbol_Depth
	//向服务器发送连接请求，websocket
	connect, _, err := dialer.Dial(wss, nil)
	A.conn_this = connect
	if nil != err {
		log.Println(err)
		return
	}
	defer connect.Close()
	//发送订阅请求
	subscribe_a_server(connect, params)
	//读取订阅请求
	for {
		//判断是否退出循环
		//fmt.Println(A)
		if A.Judge {
			fmt.Println("conn通道已关闭")
			break
		}
		//messageType 消息类型，websocket 标准，messageData 消息数据
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			*info_gather <- messageData
		case websocket.BinaryMessage: //二进制数据
			*info_gather <- messageData
			fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
			connect.WriteMessage(websocket.PongMessage, messageData)
		case websocket.PongMessage: //Pong
		default:

		}
	}

}

//订阅一个服务
func subscribe_a_server(conn *websocket.Conn, params []string) {
	//设计请求结构体
	sbc := ask_websocket{
		Method: "SUBSCRIBE",
		Params: params,
		Id:     1,
	}
	//将请求转为JSON的[]byte
	message_send, err := json.Marshal(sbc)
	if err != nil {
		log.Panicln("请求结构体json解析失败", err)
	}
	//发送websocket请求
	err = conn.WriteMessage(websocket.TextMessage, message_send)
	if err != nil {
		log.Println("订阅失败", err)
	} else {
		fmt.Println("订阅已申请")
	}
}

//取消一个服务
func unsubscribe_a_server(conn *websocket.Conn, params []string) {
	//设计请求结构体
	sbc := ask_websocket{
		Method: "UNSUBSCRIBE",
		Params: params,
		Id:     312,
	}
	//将请求转为JSON的[]byte
	message_send, err := json.Marshal(sbc)
	if err != nil {
		log.Panicln("请求结构体json解析失败", err)
	}
	//发送websocket请求
	err = conn.WriteMessage(websocket.TextMessage, message_send)
	if err != nil {
		log.Println("取消订阅失败", err)
	} else {
		fmt.Println("订阅已取消")
	}
}

/*一个脚本案例
package main

import (
	"fmt"
	"quant/get_info"
	"time"
)

func main() {
	pa := []string{"btcusdt@miniTicker"}
	new_ask := get_info.Ask_a_wb{Info_type: get_info.Symbol_Mini_Ticker, Params: pa}
	infogather := make(chan []byte, 2)
	new_ask.Start_wb(&infogather)
	count := 0
	for {
		count++
		data := <-infogather
		fmt.Println(string(data))
		if count == 6 {
			new_ask.End_wb()
			break
		}
	}
	time.Sleep(time.Second * 3)
	new_ask.Restart_wb()

	for {
		count++
		data := <-infogather
		fmt.Println(string(data))
		if count == 15 {
			new_ask.End_Conn()
			time.Sleep(time.Second)
			break
		}
	}

}

*/
