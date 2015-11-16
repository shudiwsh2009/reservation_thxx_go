var width=$(window).width();
var height=$(window).height();
var teacher;
var reservations;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/Reservation/teacher/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				console.log(data);
				reservations = data.reservations;
				teacher = data.teacher_info;
				refreshDataTable(reservations);
				optimize();
			} else {
				alert(data.message);
			}
		}
	});
}

function refreshDataTable(reservations) {
	$("#page_maintable")[0].innerHTML = "\
		<div class='table_col' id='col_select'>\
			<div class='table_head table_cell' id='head_select'>\
				<button onclick='$(\".checkbox\").click();' style='padding:0px;'>全选</button>\
			</div>\
		</div>\
		<div class='table_col' id='col_time'>\
			<div class='table_head table_cell'>时间</div>\
		</div>\
		<div class='table_col' id='col_teacher_fullname'>\
			<div class='table_head table_cell'>咨询师</div>\
		</div>\
		<div class='table_col' id='col_teacher_mobile'>\
			<div class='table_head table_cell'>咨询师手机</div>\
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
			+ i + ")'>" + reservations[i].start_time.split(" ")[0].substr(2) + "<br>" 
			+ reservations[i].start_time.split(" ")[1] + "-" + reservations[i].end_time.split(" ")[1] + "</div>");
		$("#col_teacher_fullname").append("<div class='table_cell' id='cell_teacher_fullname_"
			+ i + "'>" + reservations[i].teacher_fullname + "</div>");
		$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_"
			+ i + "'>" + reservations[i].teacher_mobile + "</div>");
		if (reservations[i].status === "AVAILABLE") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>未预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>" 
				+ "<button type='button' id='cell_student_view_" + i + "' disabled='true'>查看"
				+ "</button></div>");
		} else if (reservations[i].status === "RESERVATED") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>已预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>" 
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='getStudent(" + i + ");'>查看"
				+ "</button></div>");
		} else if (reservations[i].status === "FEEDBACK") {
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
	$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_add'></div>");
	$("#col_status").append("<div class='table_cell' id='cell_status_add'></div>");
	$("#col_student").append("<div class='table_cell' id='cell_student_add'></div>");
}

function optimize(t) {
	$("#col_select").width(20);
	$("#col_time").width(80);
	$("#col_teacher_fullname").width(44);
	$("#col_teacher_mobile").width(76);
	$("#col_status").width(40);
	$("#col_student").width(40);
	$("#col_select").css("margin-left", (width - 300) / 2 + "px");
	for (var i = 0; i < reservations.length; ++i) {
		var maxHeight = Math.max(
				$("#cell_select_" + i).height(),
				$("#cell_time_" + i).height(),
				$("#cell_teacher_fullname_" + i).height(),
				$("#cell_teacher_mobile_" + i).height(),
				$("#cell_status_" + i).height(),
				$("#cell_student_" + i).height()
			);
		$("#cell_select_" + i).height(maxHeight);
		$("#cell_time_" + i).height(maxHeight);
		$("#cell_teacher_fullname_" + i).height(maxHeight);
		$("#cell_teacher_mobile_" + i).height(maxHeight);
		$("#cell_status_" + i).height(maxHeight);
		$("#cell_student_" + i).height(maxHeight);

		if (i % 2 == 1) {
			$("#cell_select_" + i).css("background-color", "white");
			$("#cell_time_" + i).css("background-color", "white");
			$("#cell_teacher_fullname_" + i).css("background-color", "white");
			$("#cell_teacher_mobile_" + i).css("background-color", "white");
			$("#cell_status_" + i).css("background-color", "white");
			$("#cell_student_" + i).css("background-color", "white");
		} else {
			$("#cell_select_" + i).css("background-color", "#f3f3ff");
			$("#cell_time_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_fullname_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_mobile_" + i).css("background-color", "#f3f3ff");
			$("#cell_status_" + i).css("background-color", "#f3f3ff");
			$("#cell_student_" + i).css("background-color", "#f3f3ff");
		}
	}
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");

	var s = 28;
	if (t === "add") {
		s = 68;
	}
	$("#cell_select_add").height(s);
	$("#cell_time_add").height(s);
	$("#cell_teacher_fullname_add").height(s);
	$("#cell_teacher_mobile_add").height(s);
	$("#cell_status_add").height(s);
	$("#cell_student_add").height(s);
	$(".table_head").height($("#head_select").height());
}

function addReservation() {
	$("#cell_time_add")[0].onclick = "";
	$("#cell_time_add")[0].innerHTML = "<input type='text' id='input_date' style='width: 60px'></input><br>"
		+ "<input style='width:15px' id='start_hour'></input>时<input style='width:15px' id='start_minute'></input>分<br>"
		+ "<input style='width:15px' id='end_hour'></input>时<input style='width:15px' id='end_minute'></input>分";
	$("#cell_teacher_fullname_add")[0].innerHTML = "<input id='teacher_fullname' style='width:80px' value='" + teacher.teacher_fullname + "'></input>";
	$("#cell_teacher_mobile_add")[0].innerHTML = "<input id='teacher_mobile' style='width:120px' value='" + teacher.teacher_mobile + "'></input>";
	$("#cell_status_add")[0].innerHTML = "<button type='button' onclick='addReservationConfirm();'>确认</button>";
	$("#cell_student_add")[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
	$("#input_date").DatePicker({
		format: "YY-m-dd",
		date: $("#input_date").val(),
		current: $("#input_date").val(),
		starts: 1,
		position: "r",
		onBeforeShow: function() {
			$("#input_date").DatePickerSetDate($("#input_date").val(), true);
		},
		onChange: function(formated, dates) {
			$("#input_date").val(formated);
			$("#input_date").val($("#input_date").val().substr(4, 10));
			$("#input_date").DatePickerHide();
		},
	});
	optimize("add");
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
	var payload = {
		start_time: startTime,
		end_time: endTime,
		teacher_fullname: $("#teacher_fullname").val(),
		teacher_mobile: $("#teacher_mobile").val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/reservation/add",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function editReservation(index) {
	$("#cell_time_" + index).height(68);
	$("#cell_time_" + index)[0].onclick = "";
	$("#cell_time_" + index)[0].innerHTML = "<input type='text' id='input_date' style='width:60px'></input><br>"
		+ "<input style='width:15px' id='start_hour'></input>时<input style='width:15px' id='start_minute'></input>分"
		+ "<input style='width:15px' id='end_hour'></input>时<input style='width:15px' id='end_minute'></input>分";
	$("#cell_teacher_fullname_" + index)[0].innerHTML = "<input id='teacher_fullname" + index + "' style='width:80px' "
		+ "value='" + reservations[index].teacher_fullname + "''></input>";
	$("#cell_teacher_mobile_" + index)[0].innerHTML = "<input id='teacher_mobile" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_mobile + "'></input>";
	$("#cell_status_" + index)[0].innerHTML = "<button type='button' onclick='editReservationConfirm(" + index + ");'>确认</button>";
	$("#cell_student_" + index)[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
	$("#input_date").DatePicker({
		format: "YY-m-dd",
		date: $("#input_date").val(),
		current: $("#input_date").val(),
		starts: 1,
		position: "r",
		onBeforeShow: function() {
			$("#input_date").DatePickerSetDate($("#input_date").val(), true);
		},
		onChange: function(formated, dates) {
			$("#input_date").val(formated);
			$("#input_date").val($("#input_date").val().substr(4, 10));
			$("#input_date").DatePickerHide();
		},
	});
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
	var payload = {
		reservation_id: reservations[index].reservation_id,
		start_time: startTime,
		end_time: endTime,
		teacher_fullname: $("#teacher_fullname" + index).val(),
		teacher_mobile: $("#teacher_mobile" + index).val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/reservation/edit",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function removeReservations() {
	$("body").append("\
		<div class='delete_teacher_pre'>\
			确认删除选中的咨询记录？\
			<br>\
			<button type='button' onclick='$(\".delete_teacher_pre\").remove();removeReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".delete_teacher_pre\").remove();'>取消</button>\
		</div>\
	");
	optimize(".delete_teacher_pre");
}

function removeReservationsConfirm() {
	var reservationIds = [];
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			reservationIds.push(reservations[i].reservation_id);
		}
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/reservation/remove",
		data: payload,
		traditional: true,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function cancelReservations() {
	$("body").append("\
		<div class='cancel_teacher_pre'>\
			确认取消选中的预约？\
			<br>\
			<button type='button' onclick='$(\".cancel_teacher_pre\").remove();cancelReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".cancel_teacher_pre\").remove();'>取消</button>\
		</div>\
	");
	optimize(".cancel_teacher_pre");
}

function cancelReservationsConfirm() {
	var reservationIds = [];
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			reservationIds.push(reservations[i].reservation_id);
		}
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/reservation/cancel",
		data: payload,
		traditional: true,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function getFeedback(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/reservation/feedback/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showFeedback(index, data.feedback);
			} else {
				alert(data.message);
			}
		},
	});
}

function showFeedback(index, feedback) {
	$("body").append("\
		<div class='fankui_tch' id='feedback_table_" + index + "' style='font-size:11px;text-align:left;top:100px;height:300px;width:200px;left:100px'>\
			咨询师反馈表<br>\
			问题评估：<br>\
			<textarea id='problem' style='width:180px;height:80px'></textarea><br>\
			咨询记录：<br>\
			<textarea id='record' style='width:180px;height:80px'></textarea><br>\
			<button type='button' onclick='submitFeedback(" + index + ");'>提交</button>\
			<button type='button' onclick='$(\".fankui_tch\").remove();'>取消</button>\
		</div>\
	");
	$("#problem").val(feedback.problem);
	$("#record").val(feedback.record);
	optimize(".fankui_tch");
}

function submitFeedback(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
		problem: $("#problem").val(),
		record: $("#record").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/reservation/feedback/submit",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				successFeedback();
			} else {
				alert(data.message);
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
		reservation_id: reservations[index].reservation_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/teacher/student/get",
		data: payload, 
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showStudent(data.student_info);
			} else {
				alert(data.message);
			}
		}
	});
}

function showStudent(student) {
	$("body").append("\
		<div class='admin_chakan' style='text-align: left'>\
			学号：" + student.student_username + "<br>\
			姓名：" + student.student_fullname + "<br>\
			性别：" + student.student_gender + "<br>\
			出生日期：" + student.student_birthday + "<br>\
			系别：" + student.student_school + "<br>\
			年级：" + student.student_grade + "<br>\
			现住址：" + student.student_current_address + "<br>\
			家庭住址：" + student.student_family_address + "<br>\
			联系电话：" + student.student_mobile + "<br>\
			Email：" + student.student_email + "<br>\
			咨询经历：" + (student.student_experience_time ? "时间：" + student.student_experience_time + " 地点：" + student.student_experience_location + " 咨询师：" + stduent.student_experience_teacher : "无") + "<br>\
			父亲年龄：" + student.student_father_age + " 职业：" + student.student_father_job + " 学历：" + student.student_father_edu + "<br>\
			母亲年龄：" + student.student_mother_age + " 职业：" + student.student_mother_job + " 学历：" + student.student_mother_edu + "<br>\
			父母婚姻状况：" + student.student_parent_marriage + "<br>\
			近三个月里发生的有重大意义的事：" + student.student_significant + "<br>\
			需要接受帮助的主要问题：" + student.student_problem + "<br>\
			<br>\
			已绑定的咨询师：<span id='binded_teacher_username'>" + student.student_binded_teacher_username + "</span>&nbsp;\
				<span id='binded_teacher_fullname'>" + student.student_binded_teacher_fullname + "</span><br>\
			<br>\
			<button type='button' onclick='$(\".admin_chakan\").remove();'>关闭</button>\
		</div>\
	");
	optimize(".admin_chakan");
}