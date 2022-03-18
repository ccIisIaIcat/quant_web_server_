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

type K_line_info struct {
	Instrument_id     string             //币对
	Infogather        chan []byte        //用于存放数据的channal
	Period_length     string             //k线的区间长度，只可取一些特殊值
	websocket_address string             //websocket地址
	connect           *websocket.Conn    //websocket连接
	ctx               context.Context    //更为优雅的停止向channel中输入数据
	cancel            context.CancelFunc //更为优雅的停止向channel中输入数据
	wg                sync.WaitGroup     //用于优雅的结束进程
	logger            *log.Logger        //用于输出日志文件

}

func (K *K_line_info) return_string_period() string {
	switch K.Period_length {
	case "1Y":
		return "candle1Y"
	case "6M":
		return "candle6M"
	case "3M":
		return "candle3M"
	case "1M":
		return "candle1M"
	case "1W":
		return "candle1W"
	case "5D":
		return "candle5D"
	case "3D":
		return "candle3D"
	case "2D":
		return "candle2D"
	case "1D":
		return "candle1D"
	case "12H":
		return "candle12H"
	case "6H":
		return "candle6H"
	case "4H":
		return "candle4H"
	case "2H":
		return "candle2H"
	case "1H":
		return "candle1H"
	case "30m":
		return "candle30m"
	case "15m":
		return "candle15m"
	case "5m":
		return "candle5m"
	case "3m":
		return "cadle3m"
	case "1m":
		return "candle1m"

	default:
		K.logger.Panic("所输入时间区间不存在")
		return "err"
	}
}

//初始化
func (K *K_line_info) init() {
	K.wg.Add(1)
	K.Infogather = make(chan []byte, 100)
	K.ctx, K.cancel = context.WithCancel(context.Background())
	K.websocket_address = "wss://ws.okx.com:8443/ws/v5/public"
	file := "./log日志/" + time.Now().Format("20060102") + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	K.logger = log.New(logFile, "[ok_quant:spot_goods_k_line_info]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出

}

//建立wevsocket连接
func (K *K_line_info) start_websocket() {
	dialer := websocket.Dialer{}
	var err error
	K.connect, _, err = dialer.Dial(K.websocket_address, nil)
	if nil != err {
		K.logger.Println(err)
		return
	}
	K.logger.Println("websocket连接已建立")
}

//结束websocket连接
func (K *K_line_info) end_websocket() {
	K.connect.Close()
}

//向websocket订阅一个服务
func (K *K_line_info) subscribe_a_server() {

	i_f_subscribe := info_sender{Op: "subscribe", Args: []args_struct{{Channel: K.return_string_period(), InstID: K.Instrument_id}}}
	message_send, err := json.Marshal(i_f_subscribe)
	if err != nil {
		K.logger.Println("转为json文件失败：", err)
	}
	err = K.connect.WriteMessage(websocket.TextMessage, message_send)
	if err != nil {
		K.logger.Println("订阅失败", err)
	} else {
		K.logger.Println("订阅已申请,订阅json:", string(message_send))
	}
}

//向websocket取消一个订阅
func (K *K_line_info) unsubscribe_a_server() {
	i_f_unsubscribe := info_sender{Op: "unsubscribe", Args: []args_struct{{Channel: K.return_string_period(), InstID: K.Instrument_id}}}
	message_send, err := json.Marshal(i_f_unsubscribe)
	if err != nil {
		K.logger.Println("转为json文件失败：", err)
	}
	err = K.connect.WriteMessage(websocket.TextMessage, message_send)
	if err != nil {
		K.logger.Println("订阅取消失败", err)
	} else {
		K.logger.Println("订阅已申请取消,取消json:", string(message_send))
	}
}

//将收到的信息存放于特定channel中
func (K *K_line_info) save_in_channel() {
LOOP:
	for {
		messageType, messageData, err := K.connect.ReadMessage()
		if err != nil {
			K.logger.Println("websocket_get_err:", err)
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			K.Infogather <- messageData
		case websocket.BinaryMessage: //二进制数据
			K.Infogather <- messageData
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
			K.connect.WriteMessage(websocket.PongMessage, messageData)
		case websocket.PongMessage: //Pong
		default:
			K.logger.Println("unkown_message_type:", messageType)
		}
		select {
		case <-K.ctx.Done():
			break LOOP
		default:
		}
	}
}

//用于结束通道
func (K *K_line_info) End_channel() {
	K.cancel()
	K.wg.Done()
	K.logger.Println("通道已关闭")
	K.unsubscribe_a_server()
	K.logger.Println("websocket握手已关闭")
}

//开始tiker信息的获取
func (K *K_line_info) Start_ticker_info() {
	K.init()
	K.start_websocket()
	K.subscribe_a_server()
	K.save_in_channel()
}
