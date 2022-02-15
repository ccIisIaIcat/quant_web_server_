

function bids_orderbook() {
    var bids_table = document.getElementById("bids_table")
    var asks_table = document.getElementById("asks_table")
    var title = document.getElementById("info_set_name").value
    var time_point = document.getElementById("time_point").value
    url_ = "./graph_analysis/orderbook"+"?data_tittle="+title+"&time_point="+time_point

    var client = new HttpClient()

    client.get(url_, function(response){
        var json_ = eval('(' + response + ')');

        option_infotable.series[0].data[0].value = json_["Abs"]
        option_infotable.series[0].data[1].value = json_["Volum"]
        option_infotable.series[0].data[2].value = json_["Average"]
        option_infotable.series[0].data[3].value = json_["Varance"]
        console.log(option_infotable)
        sub_grapg_1.setOption(option_infotable)
            
        // console.log(json_)
        for (i=1;i<=bids_table.rows.length;i++) {
            bids_table.rows[i].cells[0].innerHTML = json_["Bids_price_list"][i-1]
            asks_table.rows[i].cells[0].innerHTML = json_["Asks_price_list"][i-1]
            bids_table.rows[i].cells[1].innerHTML = json_["Asks_quantity_list"][i-1]
            asks_table.rows[i].cells[1].innerHTML = json_["Bids_quantity_list"][i-1]
            if (json_["Bids_info_list_pre"][i-1] != 0) {
                if (json_["Bids_info_list_pre"][i-1] == -1) {
                    bids_table.rows[i].cells[0].bgColor = "Lime"
                }
                if (json_["Bids_info_list_pre"][i-1] == -2) {
                    bids_table.rows[i].cells[0].bgColor = "PaleGreen"
                }
                if (json_["Bids_info_list_pre"][i-1] == -3) {
                    bids_table.rows[i].cells[0].bgColor = "Green"
                }
            }else{
                bids_table.rows[i].cells[0].bgColor = "white"
            }
            if (json_["Asks_info_list_pre"][i-1] != 0) {
                if (json_["Asks_info_list_pre"][i-1] == -1) {
                    asks_table.rows[i].cells[0].bgColor = "Lime"
                }
                if (json_["Asks_info_list_pre"][i-1] == -2) {
                    asks_table.rows[i].cells[0].bgColor = "PaleGreen"
                }
                if (json_["Asks_info_list_pre"][i-1] == -3) {
                    asks_table.rows[i].cells[0].bgColor = "Green"
                }
            }else{
                asks_table.rows[i].cells[0].bgColor = "white"
            }
            if (json_["Bids_info_list_next"][i-1] != 0) {
                if (json_["Bids_info_list_next"][i-1] == 1) {
                    bids_table.rows[i].cells[1].bgColor = "Red"
                }
                if (json_["Bids_info_list_next"][i-1] == 2) {
                    bids_table.rows[i].cells[1].bgColor = "Salmon"
                }
                if (json_["Bids_info_list_next"][i-1] == 3) {
                    bids_table.rows[i].cells[1].bgColor = "Brown"
                }
            }else{
                bids_table.rows[i].cells[1].bgColor = "white"
            }
            if (json_["Asks_info_list_next"][i-1] != 0) {
                if (json_["Asks_info_list_next"][i-1] == 1) {
                    asks_table.rows[i].cells[1].bgColor = "Red"
                }
                if (json_["Asks_info_list_next"][i-1] == 2) {
                    asks_table.rows[i].cells[1].bgColor = "Salmon"
                }
                if (json_["Asks_info_list_next"][i-1] == 3) {
                    asks_table.rows[i].cells[1].bgColor = "Brown"
                }
            }
            else{
                asks_table.rows[i].cells[1].bgColor = "white"
            }
            
        }
              

        })
    // for (i=0;i<bids_table.rows.length;i++){
    //     console.log(bids_table.rows[i].value)
    // }
}

