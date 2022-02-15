function ticker_info_around_point() {
    var title = document.getElementById("info_set_name").value
    var time_point = document.getElementById("time_point").value
    if (title == "" || time_point == "") {alert("输入信息不全")}
    else{
        price_detail()
    }

}


function price_detail() {
    var title = document.getElementById("info_set_name").value
    var time_point = document.getElementById("time_point").value
    url_ = "./graph_analysis/price_detail"+"?data_tittle="+title+"&time_point="+time_point
    var client = new HttpClient();
    
    client.get(url_, function(response){
    var json_ = eval('(' + response + ')');
        var data_around_point_p_q_b = []
        var data_around_point_p_q_a = []
        for (i=0;i<json_["T"].length;i++) {
            data_around_point_p_q_b.push([json_["T"][i],json_["Bp"][i],json_["Bq"][i]])
            data_around_point_p_q_a.push([json_["T"][i],json_["Ap"][i],json_["Aq"][i]])
        }
        console.log(json_["T"])
        option.series[option.series.length]= {
            color:'rgba(0, 100, 0,0.1)',
            type:"scatter",
            data:data_around_point_p_q_b,
            
            symbolSize: function (data_around_point_p_q_b) { 
                return data_around_point_p_q_b[2]*bubble_size; 
            }
            
        }
        option.series[option.series.length]= {
            color: 'rgba(256,0,0,0.1)',
            type:"scatter",
            data:data_around_point_p_q_a,
            
            symbolSize: function (data_around_point_p_q_b) { 
                return data_around_point_p_q_b[2]*bubble_size; 
            }
        }
        main_graph.setOption(option)
        main_graph.resize();
    })
    
}

function bubble_skip() {
    bubble_size = document.getElementById("bubble_size").value
    main_graph.setOption(option)
    main_graph.resize();
}

function bubble_bigger() {
    bubble_size = bubble_size*1.2
    main_graph.setOption(option)
    main_graph.resize();
}

function bubble_smaller() {
    bubble_size = bubble_size*0.8
    main_graph.setOption(option)
    main_graph.resize();

}

