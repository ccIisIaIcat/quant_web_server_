option = {
    series: [
        {
            type:'line',
            data: [[1,70],[2,60],[3,40],[4,45],[5,45],[6,56],[7,77],[8,38],[9,79],[10,50]],
            
            markLine: {
                silent: true,
                data: [{
                    yAxis: 65
                }, {
                    yAxis: 70
                }, {
                    yAxis: 68
                }],
                lineStyle: {
                    normal: {
                    type: 'solid',
                },
            }
        }

        },
        {
            type:"scatter",
            data:[[3.5,88,10],[6.5,68,20]],
            symbolSize: function (data) { 
                return data[2]; 
            }, 
        },
        {
            type:"scatter",
            data:[[5,98,7],[7,78,15]],
            symbolSize: function (data) { 
                return data[2]; 
            }, 
        }

    ]
}

// option.series[0] = {
//         type:'line',
//         data: [70,60,40,45,45,56,77,38,79,50]
// }


main_graph.setOption(option)

function get_info_title(){
    var title = document.getElementById("info_set_name").value
    if (title == "") {alert("输入信息不全")}
    else{
        url_ = "./graph_analysis/price_by_second"+"?data_tittle="+title
        var client = new HttpClient();
        client.get(url_, function(response){
        var json_ = eval('(' + response + ')');
        var x_data = []
        new_data = []
        max_min_start_end[0] = 0
        max_min_start_end[1] = 9999999
        for (i = max_min_start_end[2]; i < Number(period)+  Number(max_min_start_end[2]); i++) { 

            new_data.push([i,json_[i]])
            x_data.push(i)
            if (json_[i] > max_min_start_end[0]) {
                max_min_start_end[0] = json_[i]
            }
            if (json_[i] < max_min_start_end[1]) {
                max_min_start_end[1] = json_[i]
            }
        }
        option = {
            xAxis: {
                data: x_data,
                min:max_min_start_end[2],
                max:Number(period)+  Number(max_min_start_end[2])-1
            },
            yAxis: {
                min:max_min_start_end[1],
                max:max_min_start_end[0]
            },
            series: [
                {
                    type:'line',
                    data: new_data
                }
            ]
        }
        
        main_graph.setOption(option)
        main_graph.resize();
    })
    }
}