package spot_goods

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Websocket_general struct {
	Instrument_id        string                    //币对
	Infogather_set       map[string](*chan []byte) //用于存放数据的channal的指针地址
	Infogather_orderbook chan map[string]string    //用于存放处理完的orderbook信息
	websocket_address    string                    //websocket地址
	connect              *websocket.Conn           //websocket连接
	ctx                  context.Context           //更为优雅的停止向channel中输入数据
	cancel               context.CancelFunc        //更为优雅的停止向channel中输入数据
	wg                   sync.WaitGroup            //用于优雅的结束进程
	logger               *log.Logger               //用于输出日志文件
}
