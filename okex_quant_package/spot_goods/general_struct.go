package spot_goods

//信息传递结构体
type info_sender struct {
	Op   string        `json:"op"`
	Args []args_struct `json:"args"`
}

type args_struct struct {
	Channel string `json:"channel"`
	InstID  string `json:"instId"`
}

//orderbook结构体
type orderbook_infogather_struct struct {
	Args   args_struct      `json:"arg"`
	Action string           `json:"action"`
	Data   []orderbook_data `json:"data"`
}

//orderbook子结构体
type orderbook_data struct {
	Asks     [][]string `json:"asks"`
	Bids     [][]string `json:"bids"`
	Ts       int64      `json:"ts"`
	Checksum int        `json:"checksum"`
}

//orderbook用于存储在mysql里的结构体
type orderbook_outputer struct {
	T     int64             `json:"T"`
	B_dic map[string]string `json:"bids"`
	A_dic map[string]string `json:"asks"`
}
