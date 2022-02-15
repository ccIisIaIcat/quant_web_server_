package local_ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func Local_web(wb_name string, info_gather_this *chan []byte) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	http.HandleFunc(wb_name, func(rw http.ResponseWriter, r *http.Request) {
		//建立一个websocket通道
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			fmt.Println("报错！", err)
		}
		//传送数据
		for {
			msg := <-*info_gather_this
			// fmt.Println(string(msg))
			//fmt.Println(string(msg))
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("报错！", err)
			}

		}
	})
}
