var width = $(window).width();
var height = $(window).height();
var reservations;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/api/admin/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				console.log(data);
				reservations = data.payload.reservations;
				refreshDataTable(reservations);
				optimize();
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function queryReservations() {
	var payload = {
		from_date: $("#query_date").val(),
	};
	$.ajax({
		type: "GET",
		async: false,
		url: "/api/admin/reservation/view/monthly",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				reservations = data.payload.reservations;
				refreshDataTable(reservations);
				optimize();
			}
		},
	});
}

function exportReservations() {
    var reservationIds = [];
    for (var i = 0; i < reservations.length; ++i) {
        if ($("#cell_checkbox_" + i)[0].checked) {
            reservationIds.push(reservations[i].id);
        }
    }
    var payload = {
        reservation_ids: reservationIds,
    };
    $.ajax({
        type: "POST",
        async: false,
        url: "/api/admin/reservation/export",
        data: payload,
        traditional: true,
        dataType: "json",
        success: function(data) {
            if (data.status === "OK") {
                window.open(data.payload.url);
            } else {
                alert(data.err_msg);
            }
        }
    });
}

function exportReservationArrangements() {
    $.post('/api/admin/reservation/export/arrangements', {
        from_date: $('#arrangement_date').val()
    }, function(json, textStatus) {
        if (json.status === 'OK') {
            window.open(json.payload.url);
        } else {
            alert(json.err_msg);
        }
    });
}

function refreshDataTable(reservations) {
	$("#page_maintable")[0].innerHTML = "\
		<div class='table_col' id='col_select'>\
			<div class='table_head table_cell' id='head_select'>\
				<button onclick='$(\".checkbox\").click();'>全选</button>\
			</div>\
		</div>\
		<div class='table_col' id='col_time'>\
			<div class='table_head table_cell'>时间</div>\
		</div>\
		<div class='table_col' id='col_teacher_fullname'>\
			<div class='table_head table_cell'>咨询师</div>\
		</div>\
		<div class='table_col' id='col_teacher_username'>\
			<div class='table_head table_cell'>咨询师编号</div>\
		</div>\
		<div class='table_col' id='col_teacher_mobile'>\
			<div class='table_head table_cell'>咨询师手机</div>\
		</div>\
		<div class='table_col' id='col_teacher_address'>\
			<div class='table_head table_cell'>咨询师地址</div>\
		</div>\
		<div class='table_col' id='col_international_type'>\
        	<div class='table_head table_cell'>支持语言</div>\
        </div>\
		<div class='table_col' id='col_status'>\
			<div class='table_head table_cell'>状态</div>\
		</div>\
		<div class='table_col' id='col_student'>\
			<div class='table_head table_cell'>学生</div>\
		</div>\
		<div class='clearfix'></div>\
	";

	for (var i = 0; i < reservations.length; ++i) {
		$("#col_select").append("<div class='table_cell' id='cell_select_" + i + "'>"
			+ "<input class='checkbox' type='checkbox' id='cell_checkbox_" + i + "'></div>");
		$("#col_time").append("<div class='table_cell' id='cell_time_" + i + "' onclick='editReservation("
			+ i + ")'>" + reservations[i].start_time + "至" + reservations[i].end_time + "</div>");
		$("#col_teacher_fullname").append("<div class='table_cell' id='cell_teacher_fullname_"
			+ i + "' onclick='getTeacher(" + i + ")'>" + reservations[i].teacher_fullname + "<br>"
			+ reservations[i].teacher_fullname_en + "</div>");
		$("#col_teacher_username").append("<div class='table_cell' id='cell_teacher_username_"
			+ i + "' onclick='getTeacher(" + i + ")'>" + reservations[i].teacher_username + "</div>");
		$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_"
			+ i + "' onclick='getTeacher(" + i + ")'>" + reservations[i].teacher_mobile + "</div>");
		$("#col_teacher_address").append("<div class='table_cell' id='cell_teacher_address_"
			+ i + "' onclick='getTeacher(" + i + ")'>" + reservations[i].teacher_address + "<br>"
			+ reservations[i].teacher_address_en + "</div>");
		if (reservations[i].international_type === 0) {
			$("#col_international_type").append("<div class='table_cell' id='cell_international_type_" + i + "'>仅中文</div>");
		} else if (reservations[i].international_type === 1) {
			$("#col_international_type").append("<div class='table_cell' id='cell_international_type_" + i + "'>中英双语</div>");
		}
		if (reservations[i].status === 1) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>未预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>"
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='makeReservation(" + i + ");'>帮约"
				+ "</button></div>");
		} else if (reservations[i].status === 2) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>已预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>"
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='getStudent(" + i + ");'>查看"
				+ "</button></div>");
		} else if (reservations[i].status === 3) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>"
				+ "<button type='button' id='cell_status_feedback_" + i + "' onclick='getFeedback(" + i + ");'>"
				+ "反馈</button></div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>"
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='getStudent(" + i + ");'>查看"
				+ "</button></div>");
		}
	}
	$("#col_select").append("<div class='table_cell' id='cell_select_add'><input type='checkbox'></div>");
	$("#col_time").append("<div class='table_cell' id='cell_time_add' onclick='addReservation();'>点击新增</div>");
	$("#col_teacher_fullname").append("<div class='table_cell' id='cell_teacher_fullname_add'></div>");
	$("#col_teacher_username").append("<div class='table_cell' id='cell_teacher_username_add'></div>");
	$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_add'></div>");
	$("#col_teacher_address").append("<div class='table_cell' id='cell_teacher_address_add'></div>");
	$("#col_international_type").append("<div class='table_cell' id='cell_international_type_add'></div>");
	$("#col_status").append("<div class='table_cell' id='cell_status_add'></div>");
	$("#col_student").append("<div class='table_cell' id='cell_student_add'></div>");
}

function optimize(t) {
	$("#col_select").width(40);
	$("#col_time").width(305);
	$("#col_teacher_fullname").width(120);
	$("#col_teacher_username").width(110);
	$("#col_teacher_mobile").width(110);
	$("#col_teacher_address").width(300);
	$("#col_international_type").width(105);
	$("#col_status").width(85);
	$("#col_student").width(85);
	// $('#col0').css('margin-left',width*0.02+'px')
	for (var i = 0; i < reservations.length; ++i) {
		var maxHeight = Math.max(
				$("#cell_select_" + i).height(),
				$("#cell_time_" + i).height(),
				$("#cell_teacher_fullname_" + i).height(),
				$("#cell_teacher_username_" + i).height(),
				$("#cell_teacher_mobile_" + i).height(),
				$("#cell_teacher_address_" + i).height(),
				$("#cell_international_type_" + i).height(),
				$("#cell_status_" + i).height(),
				$("#cell_student_" + i).height()
			);
		$("#cell_select_" + i).height(maxHeight);
		$("#cell_time_" + i).height(maxHeight);
		$("#cell_teacher_fullname_" + i).height(maxHeight);
		$("#cell_teacher_username_" + i).height(maxHeight);
		$("#cell_teacher_mobile_" + i).height(maxHeight);
		$("#cell_teacher_address_" + i).height(maxHeight);
		$("#cell_international_type_" + i).height(maxHeight);
		$("#cell_status_" + i).height(maxHeight);
		$("#cell_student_" + i).height(maxHeight);

		if (i % 2 == 1) {
			$("#cell_select_" + i).css("background-color", "white");
			$("#cell_time_" + i).css("background-color", "white");
			$("#cell_teacher_fullname_" + i).css("background-color", "white");
			$("#cell_teacher_username_" + i).css("background-color", "white");
			$("#cell_teacher_mobile_" + i).css("background-color", "white");
			$("#cell_teacher_address_" + i).css("background-color", "white");
			$("#cell_international_type_" + i).css("background-color", "white");
			$("#cell_status_" + i).css("background-color", "white");
			$("#cell_student_" + i).css("background-color", "white");
		} else {
			$("#cell_select_" + i).css("background-color", "#f3f3ff");
			$("#cell_time_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_fullname_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_username_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_mobile_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_address_" + i).css("background-color", "#f3f3ff");
			$("#cell_international_type_" + i).css("background-color", "#f3f3ff");
			$("#cell_status_" + i).css("background-color", "#f3f3ff");
			$("#cell_student_" + i).css("background-color", "#f3f3ff");
		}
	}
	$("#cell_select_add").height(50);
	$("#cell_time_add").height(50);
	$("#cell_teacher_fullname_add").height(50);
	$("#cell_teacher_username_add").height(50);
	$("#cell_teacher_mobile_add").height(50);
	$("#cell_teacher_address_add").height(50);
	$("#cell_international_type_add").height(50);
	$("#cell_status_add").height(50);
	$("#cell_student_add").height(50);

	$(".table_head").height($("#head_select").height());
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
	$("#page_maintable").css("margin-left", 0.5 * ($(window).width()
		- (40 + 305 + 120 + 110 + 110 + 300 + 105 + 85 + 85)) + "px");
}

function addReservation() {
	$("#cell_time_add")[0].onclick = "";
	$("#cell_time_add")[0].innerHTML = "<input type='text' id='input_date' style='width: 80px'/>日，"
		+ "<input style='width:20px' id='start_hour'/>时<input style='width:20px' id='start_minute'/>分"
		+ "至<input style='width:20px' id='end_hour'/>时<input style='width:20px' id='end_minute'/>分";
	$("#cell_teacher_fullname_add")[0].innerHTML = "<input id='teacher_fullname' style='width:60px'/>"
		+ "<button type='button' onclick='searchTeacher();'>搜索</button><br>"
		+ "<input id='teacher_fullname_en' style='width:60px'/>";
	$("#cell_teacher_username_add")[0].innerHTML = "<input id='teacher_username' style='width:120px'/>";
	$("#cell_teacher_mobile_add")[0].innerHTML = "<input id='teacher_mobile' style='width:120px'/>";
	$("#cell_teacher_address_add")[0].innerHTML = "<input id='teacher_address' style='width:180px' value='紫荆C楼407室'/><br>"
		+ "<input id='teacher_address_en' style='width:180px' value='Room 407, Zijing C Building'/>";
	$("#cell_international_type_add")[0].innerHTML = "<select id='international_type'><option value='0'>仅中文</option><option value='1'>中英双语</option></select>"
	$("#cell_status_add")[0].innerHTML = "<button type='button' onclick='addReservationConfirm();'>确认</button>";
	$("#cell_student_add")[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
	$("#input_date").datepicker({
        showOtherMonths: true,
        selectOtherMonths: true,
        showButtonPanel: true,
        dateFormat: 'yy-mm-dd',
        showWeek: true,
        firstDay: 1
    });
	optimize();
}

function addReservationConfirm() {
	var startHour = $("#start_hour").val();
	var startMinute = $("#start_minute").val();
	var endHour = $("#end_hour").val();
	var endMinute = $("#end_minute").val();
	var startTime = $("#input_date").val() + " " + (startHour.length < 2 ? "0" : "") + startHour + ":";
	if (startMinute.length == 0) {
		startTime += "00";
	} else if (startMinute.length == 1) {
		startTime += "0" + startMinute;
	} else {
		startTime += startMinute;
	}
	var endTime = $("#input_date").val() + " " + (endHour.length < 2 ? "0" : "") + endHour + ":";
	if (endMinute.length == 0) {
		endTime += "00";
	} else if (endMinute.length == 1) {
		endTime += "0" + endMinute;
	} else {
		endTime += endMinute;
	}
	if ($("#input_date").val().length == 0 || startHour.length == 0) {
		startTime = "";
	}
	if ($("#input_date").val().length == 0 || endHour.length == 0) {
		endTime = "";
	}
	if (startTime === "") {
		alert("开始时间为空");
		return;
	}
	if (endTime === "") {
		alert("结束时间为空");
		return;
	}
	var username = $("#teacher_username").val();
	if (username === "") {
		alert("咨询师编号为空");
		return;
	}
	var fullname = $("#teacher_fullname").val();
	if (fullname === "" ) {
		alert("咨询师姓名为空");
		return;
	}
	var fullname_en = $("#teacher_fullname_en").val();
	if (fullname_en === "" ) {
		alert("咨询师英文姓名为空");
		return;
	}
	var mobile = $("#teacher_mobile").val();
	if (mobile === "") {
		alert("咨询师手机号为空");
		return;
	}
	var address = $("#teacher_address").val();
	if (address === "") {
		alert("咨询师地址为空");
		return;
	}
	var address_en = $("#teacher_address_en").val();
	if (address_en === "") {
		alert("咨询师英文地址为空");
		return;
	}
	var payload = {
		start_time: startTime,
		end_time: endTime,
		username: username,
		fullname: fullname,
		fullname_en: fullname_en,
		mobile: mobile,
		address: address,
		address_en: address_en,
		international_type: $("#international_type").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/add",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				viewReservations();
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function editReservation(index) {
	$("#cell_time_" + index)[0].onclick = "";
	$("#cell_time_" + index)[0].innerHTML = "<input type='text' id='input_date' style='width: 80px'/>日，"
		+ "<input style='width:20px' id='start_hour'/>时<input style='width:20px' id='start_minute'/>分"
		+ "至<input style='width:20px' id='end_hour'/>时<input style='width:20px' id='end_minute'/>分";
	$("#cell_teacher_fullname_" + index)[0].onclick = "";
	$("#cell_teacher_fullname_" + index)[0].innerHTML = "<input id='teacher_fullname" + index + "' style='width:60px' "
		+ "value='" + reservations[index].teacher_fullname + "'/>"
		+ "<button type='button' onclick='searchTeacher(" + index + ");'>搜索</button><br>"
		+ "<input id='teacher_fullname_en" + index + "' style='width:60px' value='" + reservations[index].teacher_fullname_en + "'/>";
	$("#cell_teacher_username_" + index)[0].onclick = "";
	$("#cell_teacher_username_" + index)[0].innerHTML = "<input id='teacher_username" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_username + "'/>";
	$("#cell_teacher_mobile_" + index)[0].onclick = "";
	$("#cell_teacher_mobile_" + index)[0].innerHTML = "<input id='teacher_mobile" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_mobile + "'/>";
	$("#cell_teacher_address_" + index)[0].onclick = "";
	$("#cell_teacher_address_" + index)[0].innerHTML = "<input id='teacher_address" + index + "' style='width:180px' "
		+ "value='" + reservations[index].teacher_address + "'/><br>"
		+ "<input id='teacher_address_en" + index + "' style='width:180px' value='" + reservations[index].teacher_address_en + "'/>";
	$("#cell_international_type_" + index)[0].innerHTML = "<select id='international_type" + index + "'><option value='0'>仅中文</option><option value='1'>中英双语</option></select>";
	$("#cell_status_" + index)[0].innerHTML = "<button type='button' onclick='editReservationConfirm(" + index + ");'>确认</button>";
	$("#cell_student_" + index)[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
	$("#input_date").datepicker({
        showOtherMonths: true,
        selectOtherMonths: true,
        showButtonPanel: true,
        dateFormat: 'yy-mm-dd',
        showWeek: true,
        firstDay: 1
    });
	$("#international_type" + index).val(reservations[index].international_type);
	optimize();
}

function editReservationConfirm(index) {
	var startHour = $("#start_hour").val();
	var startMinute = $("#start_minute").val();
	var endHour = $("#end_hour").val();
	var endMinute = $("#end_minute").val();
	var startTime = $("#input_date").val() + " " + (startHour.length < 2 ? "0" : "") + startHour + ":";
	if (startMinute.length == 0) {
		startTime += "00";
	} else if (startMinute.length == 1) {
		startTime += "0" + startMinute;
	} else {
		startTime += startMinute;
	}
	var endTime = $("#input_date").val() + " " + (endHour.length < 2 ? "0" : "") + endHour + ":";
	if (endMinute.length == 0) {
		endTime += "00";
	} else if (endMinute.length == 1) {
		endTime += "0" + endMinute;
	} else {
		endTime += endMinute;
	}
	if ($("#input_date").val().length == 0 || startHour.length == 0) {
		startTime = "";
	}
	if ($("#input_date").val().length == 0 || endHour.length == 0) {
		endTime = "";
	}
	if (startTime === "") {
		alert("开始时间为空");
		return;
	}
	if (endTime === "") {
		alert("结束时间为空");
		return;
	}
	var username = $("#teacher_username" + index).val();
	if (username === "") {
		alert("咨询师编号为空");
		return;
	}
	var fullname = $("#teacher_fullname" + index).val();
	if (fullname === "" ) {
		alert("咨询师姓名为空");
		return;
	}
	var fullname_en = $("#teacher_fullname_en" + index).val();
	if (fullname_en === "" ) {
		alert("咨询师英文姓名为空");
		return;
	}
	var mobile = $("#teacher_mobile" + index).val();
	if (mobile === "") {
		alert("咨询师手机号为空");
		return;
	}
	var address = $("#teacher_address" + index).val();
	if (address === "") {
		alert("咨询师地址为空");
		return;
	}
	var address_en = $("#teacher_address_en" + index).val();
	if (address_en === "") {
		alert("咨询师英文地址为空");
		return;
	}
	var payload = {
		reservation_id: reservations[index].id,
		start_time: startTime,
		end_time: endTime,
		username: username,
		fullname: fullname,
		fullname_en: fullname_en,
		mobile: mobile,
		address: address,
		address_en: address_en,
		international_type: $("#international_type" + index).val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/edit",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				viewReservations();
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function searchTeacher(index) {
	var payload = {
		username: $("#teacher_username" + (index === undefined ? "" : index)).val(),
		fullname: $("#teacher_fullname" + (index === undefined ? "" : index)).val(),
		fullname_en: $("#teacher_fullname_en" + (index === undefined ? "" : index)).val(),
		mobile: $("#teacher_mobile" + (index === undefined ? "" : index)).val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/teacher/search",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				$("#teacher_username" + (index === undefined ? "" : index)).val(data.payload.teacher.username);
				$("#teacher_fullname" + (index === undefined ? "" : index)).val(data.payload.teacher.fullname);
				$("#teacher_fullname_en" + (index === undefined ? "" : index)).val(data.payload.teacher.fullname_en);
				$("#teacher_mobile" + (index === undefined ? "" : index)).val(data.payload.teacher.mobile);
				$("#teacher_address" + (index === undefined ? "" : index)).val(data.payload.teacher.address);
				$("#teacher_address_en" + (index === undefined ? "" : index)).val(data.payload.teacher.address_en);
				$("#international_type" + (index === undefined ? "" : index)).val(data.payload.teacher.international_type);
			}
		}
	});
}

function removeReservations() {
	$("body").append("\
		<div class='delete_admin_pre'>\
			确认删除选中的咨询记录？\
			<br>\
			<button type='button' onclick='$(\".delete_admin_pre\").remove();removeReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".delete_admin_pre\").remove();'>取消</button>\
		</div>\
	");
	optimize(".delete_admin_pre");
}

function removeReservationsConfirm() {
	var reservationIds = [];
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			reservationIds.push(reservations[i].id);
		}
	}
	if (reservationIds.length === 0) {
		return;
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/remove",
		data: payload,
		traditional: true,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				viewReservations();
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function cancelReservations() {
	$("body").append("\
		<div class='cancel_admin_pre'>\
			确认取消选中的预约？\
			<br>\
			<button type='button' onclick='$(\".cancel_admin_pre\").remove();cancelReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".cancel_admin_pre\").remove();'>取消</button>\
		</div>\
	");
	optimize(".cancel_admin_pre");
}

function cancelReservationsConfirm() {
	var reservationIds = [];
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			reservationIds.push(reservations[i].id);
		}
	}
	if (reservationIds.length === 0) {
		return;
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/cancel",
		data: payload,
		traditional: true,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				viewReservations();
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function makeReservation(index) {
	$("body").append("\
		<div class='admin_chakan' id='make_reservation_data_" + index + "' style='text-align:left'>\
			<div style='text-align:center;font-size:23px'>咨询申请表</div><br>\
			姓　　名：<input id='fullname'/><br>\
			性　　别：<select id='gender'><option value=''>请选择</option><option value='男'>男</option><option value='女'>女</option></select><br>\
			学　　号：<input id='username'/><br>\
			院　　系：<input id='school'/><br>\
			生 源 地：<input id='hometown'/><br>\
			手　　机：<input id='mobile'/><br>\
			邮　　箱：<input id='email'/><br>\
			以前曾做过学习发展咨询、职业咨询或心理咨询吗？<select id='experience'><option value=''>请选择</option><option value='是'>是</option><option value='否'>否</option></select><br>\
			请概括你最想要咨询的问题：<br>\
			<textarea id='problem'></textarea><br>\
			<button type='button' onclick='makeReservationConfirm(\"" + index + "\");'>确定</button>\
			<button type='button' onclick='$(\".yuyue_admin\").remove();'>取消</button>\
		</div>\
	");
	optimize(".admin_chakan");
}

function makeReservationConfirm(index) {
	var fullname = $("#fullname").val();
	if (fullname === "") {
		alert("姓名为空");
		return;
	}
	var gender = $("#gender").val();
	if (gender === "") {
		alert("性别为空");
		return;
	}
	var username = $("#username").val();
	if (username === "") {
		alert("学号为空");
		return;
	}
	var school = $("#school").val();
	if (school === "") {
		alert("院系为空");
		return;
	}
	var hometown = $("#hometown").val();
	if (hometown === "") {
		alert("生源地为空");
		return;
	}
	var mobile = $("#mobile").val();
	if (mobile === "") {
		alert("手机为空");
		return;
	}
	var email = $("#email").val();
	if (email === "") {
		alert("邮箱为空");
		return;
	}
	var experience = $("#experience").val();
	if (experience === "") {
		alert("咨询经历为空");
		return;
	}
	var problem = $("#problem").val();
	if (problem === "") {
		alert("咨询问题为空");
		return;
	}
	var payload = {
		reservation_id: reservations[index].id,
		fullname: fullname,
		gender: gender,
		username: username,
		school: school,
		hometown: hometown,
		mobile: mobile,
		email: email,
		experience: experience,
		problem: problem,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/make",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				makeReservationSuccess(index);
			} else {
				alert(data.err_msg);
			}
		},
	});
}

function makeReservationSuccess(index) {
	$(".admin_chakan").remove();
	$("#cell_student_view_" + index).attr("disabled", "true");
	$("#cell_student_view_" + index).text("查看");
	$("body").append("\
		<div class='yuyue_stu_success'>\
			你已预约成功，<br>\
			请关注短信提醒。<br>\
			<button type='button' onclick='$(\".yuyue_stu_success\").remove();viewReservations();'>确定</button>\
		</div>\
	");
	optimize(".yuyue_stu_success");
}

function getFeedback(index) {
	var payload = {
		reservation_id: reservations[index].id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/feedback/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				showFeedback(index, data.payload.feedback);
			} else {
				alert(data.err_msg);
			}
		},
	});
}

function showFeedback(index, feedback) {
	$("body").append("\
		<div class='fankui_tch' id='feedback_table_" + index + "' style='text-align:left; top:100px; height:420px; width:400px; left:100px'>\
			咨询师反馈表<br>\
			您的姓名：<input id='teacher_fullname'/><br>\
			工作证号：<input id='teacher_username'/><br>\
			来访者姓名：<input id='student_fullname'/><br>\
			来访者问题描述：<br>\
			<textarea id='problem' style='width:350px;height:80px'></textarea><br>\
			咨询师提供的问题解决方法：<br>\
			<textarea id='solution' style='width:350px;height:80px'></textarea><br>\
			对中心的工作建议：<br>\
			<textarea id='advice' style='width:350px;height:80px'></textarea><br>\
			<button type='button' onclick='submitFeedback(" + index + ");'>提交</button>\
			<button type='button' onclick='$(\".fankui_tch\").remove();'>取消</button>\
		</div>\
	");
	$("#teacher_fullname").val(feedback.teacher_fullname);
	$("#teacher_username").val(feedback.teacher_username);
	$("#student_fullname").val(feedback.student_fullname);
	$("#problem").val(feedback.problem);
	$("#solution").val(feedback.solution);
	$("#advice").val(feedback.advice_to_center);
	optimize(".fankui_tch");
}

function submitFeedback(index) {
	var teacherFullname = $("#teacher_fullname").val();
	if (teacherFullname === "" ) {
		alert("咨询师姓名为空");
		return;
	}
	var teacherUsername = $("#teacher_username").val();
	if (teacherUsername === "") {
		alert("咨询师工号为空");
		return;
	}
	var studentFullname = $("#student_fullname").val();
	if (studentFullname === "" ) {
		alert("学生姓名为空");
		return;
	}
	var problem = $("#problem").val();
	if (problem === "") {
		alert("问题描述为空");
		return;
	}
	var solution = $("#solution").val();
	if (solution === "") {
		alert("解决方法为空");
		return;
	}
	var adviceToCenter = $("#advice").val();
	if (adviceToCenter === "") {
		alert("工作建议为空");
		return;
	}
	var payload = {
		reservation_id: reservations[index].id,
		teacher_fullname: teacherFullname,
		teacher_username: teacherUsername,
		student_fullname: studentFullname,
		problem: problem,
		solution: solution,
		advice_to_center: adviceToCenter,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/feedback/submit",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				successFeedback();
			} else {
				alert(data.err_msg);
			}
		},
	});
}

function successFeedback() {
	$(".fankui_tch").remove();
	$("body").append("\
		<div class='fankui_tch_success'>\
			您已成功提交反馈！<br>\
			<button type='button' onclick='$(\".fankui_tch_success\").remove();'>确定</button>\
		</div>\
	");
	optimize(".fankui_tch_success");
}

function getStudent(index) {
	var payload = {
		reservation_id: reservations[index].id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/reservation/student/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				showStudent(data.payload.student);
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function showStudent(student) {
	$("body").append("\
		<div class='admin_chakan' style='text-align: left'>\
			姓名：" + student.fullname + "<br>\
			性别：" + student.gender + "<br>\
			院系：" + student.school + "<br>\
			学号：" + student.username + "<br>\
			手机：" + student.mobile + "<br>\
			生源地：" + student.hometown + "<br>\
			邮箱：" + student.email + "<br>\
			是否使用本系统：" + student.experience + "<br>\
			咨询问题：" + student.problem + "<br>\
			<button type='button' onclick='$(\".admin_chakan\").remove();'>返回</button>\
		</div>\
	");
	optimize(".admin_chakan");
}

function getTeacher(index) {
	var payload = {
		username: reservations[index].teacher_username,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/teacher/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				showTeacher(data.payload.teacher);
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function showTeacher(teacher) {
	$("body").append("\
		<div class='admin_chakan' style='text-align: left'>\
			姓　　名：<input id='show_teacher_fullname' value='" + teacher.fullname + "'><br>\
			Name    ：<input id='show_teacher_fullname_en' value='" + teacher.fullname_en + "'><br>\
			性　　别：<select id='show_teacher_gender'><option value=''>请选择</option><option value='男'>男</option><option value='女'>女</option></select><br>\
			Gender  ：<select id='show_teacher_gender_en'><option value=''>choose male/female</option><option value='male'>male</option><option value='female'>female</option></select><br>\
			专业背景：<input id='show_teacher_major' value='" + teacher.major + "'><br>\
			Major   ：<input id='show_teacher_major_en' value='" + teacher.major_en + "'><br>\
			学　　历：<input id='show_teacher_academic' value='" + teacher.academic + "'><br>\
			Academic：<input id='show_teacher_academic_en' value='" + teacher.academic_en + "'><br>\
			资　　质：<input id='show_teacher_aptitude' value='" + teacher.aptitude + "'><br>\
			Aptitude：<input id='show_teacher_aptitude_en' value='" + teacher.aptitude_en + "'><br>\
			可咨询的问题：<br>\
			<textarea id='show_teacher_problem' style='width:70%; height:60px;'></textarea><br>\
			Consultable Questions：<br>\
			<textarea id='show_teacher_problem_en' style='width:70%; height:60px;'></textarea><br>\
			支持语言：<select id='show_teacher_international_type'><option value='0'>仅中文</option><option value='1'>中英双语</option></select><br>\
			<button type='button' onclick='editTeacher(" + teacher.username + ");'>保存</button>\
			<button type='button' onclick='$(\".admin_chakan\").remove();'>关闭</button>\
			<span id='edit_tip' style='color: red'></span>\
		</div>\
	");
	$('#show_teacher_gender').val(teacher.gender);
	$('#show_teacher_gender_en').val(teacher.gender_en);
	$('#show_teacher_problem').val(teacher.problem);
	$('#show_teacher_problem_en').val(teacher.problem_en);
	$('#show_teacher_international_type').val(teacher.international_type);
	optimize(".admin_chakan");
}

function editTeacher(teacherUsername) {
	var fullname = $('#show_teacher_fullname').val();
	if (fullname === "") {
		alert("姓名为空");
		return;
	}
	var fullname_en = $('#show_teacher_fullname_en').val();
	if (fullname_en === "") {
		alert("英文姓名为空");
		return;
	}
	var gender = $('#show_teacher_gender').val();
	if (gender === "") {
		alert("性别为空");
		return;
	}
	var gender_en = $('#show_teacher_gender_en').val();
	if (gender_en === "") {
		alert("性别（英文）为空");
		return;
	}
	var major = $('#show_teacher_major').val();
	if (major === "") {
		alert("专业背景为空");
		return;
	}
	var major_en = $('#show_teacher_major_en').val();
	if (major_en === "") {
		alert("专业背景（英文）为空");
		return;
	}
	var academic = $('#show_teacher_academic').val();
	if (academic === "") {
		alert("学历为空");
		return;
	}
	var academic_en = $('#show_teacher_academic_en').val();
	if (academic_en === "") {
		alert("学历（英文）为空");
		return;
	}
	var aptitude = $('#show_teacher_aptitude').val();
	if (aptitude === "") {
		alert("资质为空");
		return;
	}
	var aptitude_en = $('#show_teacher_aptitude_en').val();
	if (aptitude_en === "") {
		alert("资质（英文）为空");
		return;
	}
	var problem = $('#show_teacher_problem').val();
	if (problem === "") {
		alert("可咨询问题为空");
		return;
	}
	var problem_en = $('#show_teacher_problem_en').val();
	if (problem_en === "") {
		alert("可咨询问题（英文）为空");
		return;
	}
	var payload = {
		username: teacherUsername,
		fullname: fullname,
		fullname_en: fullname_en,
		gender: gender,
		gender_en: gender_en,
		major: major,
		major_en: major_en,
		academic: academic,
		academic_en: academic_en,
		aptitude: aptitude,
		aptitude_en: aptitude_en,
		problem: problem,
		problem_en: problem_en,
		international_type: $('#show_teacher_international_type').val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/admin/teacher/edit",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				$('#edit_tip').text("更新成功！");
			} else {
				alert(data.err_msg);
			}
		}
	});
}
