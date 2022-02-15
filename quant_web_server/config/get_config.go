package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

func Read_config_mysql() (string, string, string, string, string) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("文件读取错误", err)
		os.Exit(1)
	}
	username := cfg.Section("mysql").Key("username").String()
	password := cfg.Section("mysql").Key("password").String()
	if password == "0" {
		password = ""
	}
	original_data := cfg.Section("mysql").Key("original_data").String()
	processed_data := cfg.Section("mysql").Key("processed_data").String()
	cross_feature_data := cfg.Section("mysql").Key("cross_feature_data").String()
	return username, password, original_data, processed_data, cross_feature_data
}

func Read_config_local_server() string {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("文件读取错误", err)
		os.Exit(1)
	}
	port_number := cfg.Section("local_server").Key("port_number").String()
	return port_number
}

func Read_config_feature_period() (int, int) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("文件读取错误", err)
		os.Exit(1)
	}
	time_period := cfg.Section("feature_extraction").Key("time_period").MustInt()
	goal_feature_time_period := cfg.Section("feature_extraction").Key("goal_feature_time_period").MustInt()
	return time_period, goal_feature_time_period
}

func Read_order_book_depth() int {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("文件读取错误", err)
		os.Exit(1)
	}
	book_deep := cfg.Section("orderbook_depth").Key("book_deep").MustInt()
	return book_deep
}

func Read_big_or_nor_and_time_point_period() (float64, int) {
	cfg, err := ini.Load(("./config/config.ini"))
	if err != nil {
		fmt.Println("文件读取错误", err)
		os.Exit(1)
	}
	divided_point := cfg.Section("big_point").Key("divided_point").MustFloat64()
	time_point_period := cfg.Section("big_point").Key("detail_price_period").MustInt()
	return divided_point, time_point_period
}
