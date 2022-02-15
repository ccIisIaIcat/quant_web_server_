function start_point() {
    max_min_start_end[2] = document.getElementById("start_point").value
    get_info_title()
}

function change_start_point_pre_3() {
    max_min_start_end[2] = Number(max_min_start_end[2])-Number(period)
    get_info_title()
}

function change_start_point_next_3() {
    max_min_start_end[2] = Number(max_min_start_end[2])+Number(period)
    get_info_title()
}

function change_start_point_pre_2() {
    max_min_start_end[2] = Number(max_min_start_end[2])-parseInt((Number(period)*0.5))
    get_info_title()
}

function change_start_point_next_2() {
    max_min_start_end[2] = Number(max_min_start_end[2])+parseInt((Number(period)*0.5))
    get_info_title()
}

function change_start_point_pre_1() {
    max_min_start_end[2] = Number(max_min_start_end[2])-parseInt((Number(period)*0.1))
    get_info_title()
}

function change_start_point_next_1() {
    max_min_start_end[2] = Number(max_min_start_end[2])+parseInt((Number(period)*0.1))
    get_info_title()
}

