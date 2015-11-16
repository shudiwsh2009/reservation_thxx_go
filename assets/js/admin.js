var width = $(window).width();
var height = $(window).height();
var reservations;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/Reservation/admin/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				console.log(data);
				reservations = data.reservations;
				refreshDataTable(reservations);
				optimize();
			} else {
				alert(data.message);
			}
		}
	});
}

function queryReservations() {
	var payload = {
		from_time: $("#query_date").val(),
	};
	$.ajax({
		type: "GET",
		async: false,
		url: "/Reservation/admin/reservation/view/monthly",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				console.log(data);
				reservations = data.reservations;
				refreshDataTable(reservations);
				optimize();
			}
		},
	});
}

function exportReservations() {
	var payload = {
		from_time: $("#export_date").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/reservation/export",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.open(data.url);
			} else {
				alert(data.message);
			}
		}
	});
}

function queryStudent() {
	var payload = {
		student_username: $("#query_student").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/student/query",
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
			+ i + "'>" + reservations[i].teacher_fullname + "</div>");
		$("#col_teacher_username").append("<div class='table_cell' id='cell_teacher_username_"
			+ i + "'>" + reservations[i].teacher_username + "</div>");
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
	$("#col_teacher_username").append("<div class='table_cell' id='cell_teacher_username_add'></div>");
	$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_add'></div>");
	$("#col_status").append("<div class='table_cell' id='cell_status_add'></div>");
	$("#col_student").append("<div class='table_cell' id='cell_student_add'></div>");
}

function optimize(t) {
	$("#col_select").width(40);
	$("#col_time").width(405);
	$("#col_teacher_fullname").width(120);
	$("#col_teacher_username").width(160);
	$("#col_teacher_mobile").width(160);
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
				$("#cell_status_" + i).height(),
				$("#cell_student_" + i).height()
			);
		$("#cell_select_" + i).height(maxHeight);
		$("#cell_time_" + i).height(maxHeight);
		$("#cell_teacher_fullname_" + i).height(maxHeight);
		$("#cell_teacher_username_" + i).height(maxHeight);
		$("#cell_teacher_mobile_" + i).height(maxHeight);
		$("#cell_status_" + i).height(maxHeight);
		$("#cell_student_" + i).height(maxHeight);

		if (i % 2 == 1) {
			$("#cell_select_" + i).css("background-color", "white");
			$("#cell_time_" + i).css("background-color", "white");
			$("#cell_teacher_fullname_" + i).css("background-color", "white");
			$("#cell_teacher_username_" + i).css("background-color", "white");
			$("#cell_teacher_mobile_" + i).css("background-color", "white");
			$("#cell_status_" + i).css("background-color", "white");
			$("#cell_student_" + i).css("background-color", "white");
		} else {
			$("#cell_select_" + i).css("background-color", "#f3f3ff");
			$("#cell_time_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_fullname_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_username_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_mobile_" + i).css("background-color", "#f3f3ff");
			$("#cell_status_" + i).css("background-color", "#f3f3ff");
			$("#cell_student_" + i).css("background-color", "#f3f3ff");
		}
	}
	$("#cell_select_add").height(28);
	$("#cell_time_add").height(28);
	$("#cell_teacher_fullname_add").height(28);
	$("#cell_teacher_username_add").height(28);
	$("#cell_teacher_mobile_add").height(28);
	$("#cell_status_add").height(28);
	$("#cell_student_add").height(28);

	$(".table_head").height($("#head_select").height());
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
	$("#page_maintable").css("margin-left", 0.5 * ($(window).width() 
		- (40 + 405 + 120 + 85 + 85 + 320)) + "px");
}

function addReservation() {
	$("#cell_time_add")[0].onclick = "";
	$("#cell_time_add")[0].innerHTML = "<input type='text' id='input_date' style='width: 80px'></input>日，"
		+ "<input style='width:20px' id='start_hour'></input>时<input style='width:20px' id='start_minute'></input>分"
		+ "至<input style='width:20px' id='end_hour'></input>时<input style='width:20px' id='end_minute'></input>分";
	$("#cell_teacher_fullname_add")[0].innerHTML = "<input id='teacher_fullname' style='width:60px'></input>"
		+ "<button type='button' onclick='searchTeacher();'>搜索</button>";
	$("#cell_teacher_username_add")[0].innerHTML = "<input id='teacher_username' style='width:120px'></input>";
	$("#cell_teacher_mobile_add")[0].innerHTML = "<input id='teacher_mobile' style='width:120px'></input>";
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
	var payload = {
		start_time: startTime,
		end_time: endTime,
		teacher_username: $("#teacher_username").val(),
		teacher_fullname: $("#teacher_fullname").val(),
		teacher_mobile: $("#teacher_mobile").val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/reservation/add",
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
	$("#cell_time_" + index)[0].onclick = "";
	$("#cell_time_" + index)[0].innerHTML = "<input type='text' id='input_date' style='width: 80px'></input>日，"
		+ "<input style='width:20px' id='start_hour'></input>时<input style='width:20px' id='start_minute'></input>分"
		+ "至<input style='width:20px' id='end_hour'></input>时<input style='width:20px' id='end_minute'></input>分";
	$("#cell_teacher_fullname_" + index)[0].innerHTML = "<input id='teacher_fullname" + index + "' style='width:60px' "
		+ "value='" + reservations[index].teacher_fullname + "''></input>"
		+ "<button type='button' onclick='searchTeacher();'>搜索</button>";
	$("#cell_teacher_username_" + index)[0].innerHTML = "<input id='teacher_username" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_username + "'></input>";
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
		teacher_username: $("#teacher_username" + index).val(),
		teacher_fullname: $("#teacher_fullname" + index).val(),
		teacher_mobile: $("#teacher_mobile" + index).val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/reservation/edit",
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

function searchTeacher(index) {
	var payload = {
		teacher_username: $("#teacher_username" + (index === undefined ? "" : index)).val(),
		teacher_fullname: $("#teacher_fullname" + (index === undefined ? "" : index)).val(),
		teacher_mobile: $("#teacher_mobile" + (index === undefined ? "" : index)).val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/user/teacher/search",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				$("#teacher_username" + (index === undefined ? "" : index)).val(data.teacher_username);
				$("#teacher_fullname" + (index === undefined ? "" : index)).val(data.teacher_fullname);
				$("#teacher_mobile" + (index === undefined ? "" : index)).val(data.teacher_mobile);
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
			reservationIds.push(reservations[i].reservation_id);
		}
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/reservation/remove",
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
			reservationIds.push(reservations[i].reservation_id);
		}
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/reservation/cancel",
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
		url: "/Reservation/admin/reservation/feedback/get",
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
		<div class='fankui_tch' id='feedback_table_" + index + "' style='text-align:left; top:100px; height:300px; width:400px; left:100px'>\
			咨询师反馈表<br>\
			问题评估：<br>\
			<textarea id='problem' style='width:350px;height:80px'></textarea><br>\
			咨询记录：<br>\
			<textarea id='record' style='width:350px;height:80px'></textarea><br>\
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
		url: "/Reservation/admin/reservation/feedback/submit",
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
		url: "/Reservation/admin/student/get",
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
				<span id='binded_teacher_fullname'>" + student.student_binded_teacher_fullname + "</span>\
				<button type='button' onclick='unbindStudent(" + student.student_username + ");'>解绑</button><br>\
			请输入匹配咨询师工号：<input id='teacher_username' type='text'></input>\
			<button type='button' onclick='bindStudent(" + student.student_username + ");'>绑定</button><br>\
			<br>\
			<button type='button' onclick='exportStudent(" + student.student_username + ");'>导出</button>\
			<button type='button' onclick='$(\".admin_chakan\").remove();'>关闭</button>\
		</div>\
	");
	optimize(".admin_chakan");
}

function exportStudent(studentUsername) {
	var payload = {
		student_username: studentUsername,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/student/export",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.open(data.url);
			} else {
				alert(data.message);
			}
		},
	});
}

function unbindStudent(studentUsername) {
	var payload = {
		student_username: studentUsername,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/student/unbind",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				$("#binded_teacher_username").text(data.student_info.student_binded_teacher_username);
				$("#binded_teacher_fullname").text(data.student_info.student_binded_teacher_fullname);
			} else {
				alert(data.message);
			}
		},
	});
}

function bindStudent(studentUsername) {
	var payload = {
		student_username: studentUsername,
		teacher_username: $("#teacher_username").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/admin/student/bind",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				$("#binded_teacher_username").text(data.student_info.student_binded_teacher_username);
				$("#binded_teacher_fullname").text(data.student_info.student_binded_teacher_fullname);
			} else {
				alert(data.message);
			}
		},
	});
}