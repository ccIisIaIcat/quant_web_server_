package spot_goods

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Mark_price_info struct {
	Instrument_id     string             //币对
	Infogather        chan []byte        //用于存放数据的channal
	websocket_address string             //websocket地址
	connect           *websocket.Conn    //websocket连接
	ctx               context.Context    //更为优雅的停止向channel中输入数据
	cancel            context.CancelFunc //更为优雅的停止向channel中输入数据
	wg                sync.WaitGroup     //用于优雅的结束进程
	logger            *log.Logger        //用于输出日志文件
}

//初始化
func (M *Mark_price_info) init() {
	M.wg.Add(1)
	M.Infogather = make(chan []byte, 100)
	M.ctx, M.cancel = context.WithCancel(context.Background())
	M.websocket_address = "wss://ws.okx.com:8443/ws/v5/public"
	file := "./log日志/" + time.Now().Format("20060102") + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	M.logger = log.New(logFile, "[ok_quant:spot_goods_mark_price_info]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出

}

//建立wevsocket连接
func (M *Mark_price_info) start_websocket() {
	dialer := websocket.Dialer{}
	var err error
	M.connect, _, err = dialer.Dial(M.websocket_address, nil)
	if nil != err {
		M.logger.Println(err)
		return
	}
	M.logger.Println("websocket连接已建立")
}

//结束websocket连接
func (M *Mark_price_info) end_websocket() {
	M.connect.Close()
}

//向websocket订阅一个服务
func (M *Mark_price_info) subscribe_a_server() {
	i_f_subscribe := info_sender{Op: "subscribe", Args: []args_struct{{Channel: "mark-price", InstID: M.Instrument_id}}}
	message_send, err := json.Marshal(i_f_subscribe)
	if err != nil {
		M.logger.Println("转为json文件失败：", err)
	}
	err = M.connect.WriteMessage(websocket.TextMessage, message_send)
	if err != nil {
		M.logger.Println("订阅失败", err)
	} else {
		M.logger.Println("订阅已申请,订阅json:", string(message_send))
	}
}

//向websocket取消一个订阅
func (M *Mark_price_info) unsubscribe_a_server() {
	i_f_unsubscribe := info_sender{Op: "unsubscribe", Args: []args_struct{{Channel: "mark-price", InstID: M.Instrument_id}}}
	message_send, err := json.Marshal(i_f_unsubscribe)
	if err != nil {
		M.logger.Println("转为json文件失败：", err)
	}
	err = M.connect.WriteMessage(websocket.TextMessage, message_send)
	if err != nil {
		M.logger.Println("订阅取消失败", err)
	} else {
		M.logger.Println("订阅已申请取消,取消json:", string(message_send))
	}
}

//将收到的信息存放于特定channel中
func (M *Mark_price_info) save_in_channel() {
LOOP:
	for {
		messageType, messageData, err := M.connect.ReadMessage()
		if err != nil {
			M.logger.Println("websocket_get_err:", err)
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			M.Infogather <- messageData
		case websocket.BinaryMessage: //二进制数据
			M.Infogather <- messageData
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
			M.connect.WriteMessage(websocket.PongMessage, messageData)
		case websocket.PongMessage: //Pong
		default:
			M.logger.Println("unkown_message_type:", messageType)
		}
		select {
		case <-M.ctx.Done():
			break LOOP
		default:
		}
	}
}

//用于结束通道
func (M *Mark_price_info) End_channel() {
	M.cancel()
	M.wg.Done()
	M.logger.Println("通道已关闭")
	M.unsubscribe_a_server()
	M.logger.Println("websocket握手已关闭")
}

//开始tiker信息的获取
func (M *Mark_price_info) Start_ticker_info() {
	M.init()
	M.start_websocket()
	M.subscribe_a_server()
	M.save_in_channel()
}
