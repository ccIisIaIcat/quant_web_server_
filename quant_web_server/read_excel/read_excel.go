package read_excel

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"

	"strconv"
)

//读取float64类型的数据
func D_ex(path_ string, sheet string) map[string]([]float64) {
	f, err := excelize.OpenFile(path_)
	if err != nil {
		fmt.Println("数据导入失败：", err)
	}
	rows := f.GetRows(sheet)
	N_size := len(rows) - 1
	P_size := len(rows[0]) - 1
	Y_data := make([]float64, N_size)
	X_data := make([]float64, N_size*P_size)
	temp := float64(0)
	for i := 1; i <= N_size; i++ {
		temp, _ = strconv.ParseFloat(rows[i][0], 64)
		Y_data[i-1] = temp
		for j := 1; j <= P_size; j++ {
			temp, _ = strconv.ParseFloat(rows[i][j], 64)
			X_data[(i-1)*P_size+j-1] = temp
		}
	}
	answer := make(map[string][]float64, 2)
	answer["Y"] = Y_data
	answer["X"] = X_data

	return answer

}

//读取string类型的数据
func D_ex_s(path_ string, sheet string) map[string]([]string) {
	f, err := excelize.OpenFile(path_)
	if err != nil {
		fmt.Println("数据导入失败：", err)
	}
	rows := f.GetRows(sheet)
	N_size := len(rows) - 1
	P_size := len(rows[0]) - 1
	Y_data := make([]string, N_size)
	X_data := make([]string, N_size*P_size)
	temp := ""
	for i := 1; i <= N_size; i++ {
		temp = rows[i][0]
		Y_data[i-1] = temp
		for j := 1; j <= P_size; j++ {
			temp = rows[i][j]
			X_data[(i-1)*P_size+j-1] = temp
		}
	}
	answer := make(map[string][]string, 2)
	answer["Y"] = Y_data
	answer["X"] = X_data

	return answer

}
