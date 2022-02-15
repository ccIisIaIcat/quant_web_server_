package main

import (
	"log"
	"quant_web_server/config"
	"quant_web_server/web_server"
	"time"
)

//基本信息结构体
type Basic_info struct {
	//mysql
	username           string
	password           string
	original_data      string
	processed_data     string
	cross_feature_data string
	//local_server
	port_number string
	//order_book
	order_book_depth int
	//feature_extraction
	time_period              int
	goal_feature_time_period int
	//big or not
	divided_point float64
	//graph
	time_point_period int
}

//获取基本信息
func (B *Basic_info) get_info() {
	B.username, B.password, B.original_data, B.processed_data, B.cross_feature_data = config.Read_config_mysql()
	B.port_number = config.Read_config_local_server()
	B.order_book_depth = config.Read_order_book_depth()
	B.time_period, B.goal_feature_time_period = config.Read_config_feature_period()
	B.divided_point, B.time_point_period = config.Read_big_or_nor_and_time_point_period()

}

func main() {
	basic_info := Basic_info{}
	basic_info.get_info()
	log.Println("设置读取成功")
	log.Println("开启服务")
	w_s := web_server.Web_server{Username: basic_info.username, Password: basic_info.password, Original_data: basic_info.original_data, Processed_data: basic_info.processed_data, Cross_feature_data: basic_info.cross_feature_data, Port_number: basic_info.port_number, Local_book_deep: basic_info.order_book_depth, Time_period: basic_info.time_period, Goal_feature_time_period: basic_info.goal_feature_time_period, Divided_point: basic_info.divided_point, Time_point_period: basic_info.time_point_period}
	w_s.Start_server()
	time.Sleep(time.Hour)

}
