package main

import (
	"okex_quant_package/spot_goods"
	"time"
)

func main() {
	t_i := spot_goods.Orderbook_info{Instrument_id: "BTC-USDT"}
	go t_i.Start_ticker_info()
	time.Sleep(time.Second * 10)
	// time.Sleep(time.Second)
	// for i := 0; i < 10; i++ {
	// 	new_data := <-t_i.Infogather
	// 	fmt.Println(string(new_data))
	// 	new_data = <-k_l_i.Infogather
	// 	fmt.Println(string(new_data))
	// }
	// t_i.End_channel()
	// k_l_i.End_channel()

}
