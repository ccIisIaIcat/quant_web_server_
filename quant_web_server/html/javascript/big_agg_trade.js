function big_agg_trade(){
    var title = document.getElementById("info_set_name").value
    if (title == "") {alert("输入信息不全")}
    else{
        console.log("向后台发送数据")
        url_ = "./graph_analysis/big_agg_trade"+"?data_tittle="+title
        var client = new HttpClient();
        client.get(url_, function(response){
        var json_ = eval('(' + response + ')');
        scatter_point_set_1 = []
        scatter_point_set_2 = []
        for (i = max_min_start_end[2]; i < Number(period)+  Number(max_min_start_end[2]); i++) { 
            if (json_[i] == 1){
                scatter_point_set_1.push([i,new_data[i-max_min_start_end[2]][1]])
            }else if (json_[i] == -1) {
                scatter_point_set_2.push([i,new_data[i-max_min_start_end[2]][1]])
            }
        }
        console.log(scatter_point_set_1[0])
        console.log(scatter_point_set_2[0])
        option.series[option.series.length]= {
            color: ["#00ff00"],
            type:"scatter",
            data:scatter_point_set_1,
            symbolSize: 15
            
        }
        option.series[option.series.length]={
            color: ["#ff0000"],
            type:"scatter",
            data:scatter_point_set_2,
            symbolSize: 15
        }
        main_graph.setOption(option)
    })
    }
}