var width=$(window).width();
var height=$(window).height();
var teacher;
var reservations;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/api/teacher/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				reservations = data.payload.reservations;
				teacher = data.payload.teacher;
				refreshDataTable(reservations);
				optimize();
			} else {
				alert(data.err_msg);
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
			+ i + "'>" + reservations[i].teacher_fullname + "<br>"
			+ reservations[i].teacher_fullname_en + "</div>");
		$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_"
			+ i + "'>" + reservations[i].teacher_mobile + "</div>");
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
	$("#cell_time_add")[0].innerHTML = "<input type='text' id='input_date' style='width: 60px'/><br>"
		+ "<input style='width:15px' id='start_hour'/>时<input style='width:15px' id='start_minute'/>分<br>"
		+ "<input style='width:15px' id='end_hour'/>时<input style='width:15px' id='end_minute'/>分";
	$("#cell_teacher_fullname_add")[0].innerHTML = "<input id='teacher_fullname' style='width:80px' value='" + teacher.fullname + "'/><br>"
		+ "<input id='teacher_fullname_en' style='width:80px' value='" + teacher.fullname_en + "'/>";
	$("#cell_teacher_mobile_add")[0].innerHTML = "<input id='teacher_mobile' style='width:120px' value='" + teacher.mobile + "'/>";
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
	if ($("#input_date").val().length == 0 || startHour.length == 0) {
		startTime = "";
	}
	if ($("#input_date").val().length == 0 || endHour.length == 0) {
		endTime = "";
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
	var payload = {
		start_time: startTime,
		end_time: endTime,
		fullname: $("#teacher_fullname").val(),
		fullname_en: $("#teacher_fullname_en").val(),
		mobile: $("#teacher_mobile").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/teacher/reservation/add",
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
	$("#cell_time_" + index).height(68);
	$("#cell_time_" + index)[0].onclick = "";
	$("#cell_time_" + index)[0].innerHTML = "<input type='text' id='input_date' style='width:60px'/><br>"
		+ "<input style='width:15px' id='start_hour'/>时<input style='width:15px' id='start_minute'/>分"
		+ "<input style='width:15px' id='end_hour'/>时<input style='width:15px' id='end_minute'/>分";
	$("#cell_teacher_fullname_" + index)[0].innerHTML = "<input id='teacher_fullname" + index + "' style='width:80px' "
		+ "value='" + reservations[index].teacher_fullname + "'/><br>"
		+ "<input id='teacher_fullname_en" + index + "' style='width:80px' "
		+ "value='" + reservations[index].teacher_fullname_en + "'/>";
	$("#cell_teacher_mobile_" + index)[0].innerHTML = "<input id='teacher_mobile" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_mobile + "'/>";
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
	var payload = {
		reservation_id: reservations[index].id,
		start_time: startTime,
		end_time: endTime,
		fullname: $("#teacher_fullname" + index).val(),
		fullname_en: $("#teacher_fullname_en" + index).val(),
		mobile: $("#teacher_mobile" + index).val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/teacher/reservation/edit",
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

function removeReservations() {
	$("body").append("\
		<div class='delete_teacher_pre'>\
			确认删除选中的咨询记录？\
			<br>\
			<button type='button' onclick='$(\".delete_teacher_pre\").remove();removeReservationsCheck();'>确认</button>\
			<button type='button' onclick='$(\".delete_teacher_pre\").remove();'>取消</button>\
		</div>\
	");
	optimize(".delete_teacher_pre");
}

function removeReservationsCheck() {
	var doubleCheck = false;
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			var status = reservations[i].status;
			if (status === 2 || status === 3) {
				doubleCheck = true;
				break;
			}
		}
	}
	if (doubleCheck) {
		$("body").append("\
			<div class='delete_teacher_pre'>\
				选中的咨询记录中有已经被预约或咨询的时段，请再次确认是否删除？\
				<br>\
				<button type='button' onclick='$(\".delete_teacher_pre\").remove();removeReservationsConfirm();'>确认</button>\
				<button type='button' onclick='$(\".delete_teacher_pre\").remove();'>取消</button>\
			</div>\
		");
		optimize(".delete_teacher_pre");
	} else {
		removeReservationsConfirm();
	}
}

function removeReservationsConfirm() {
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
		url: "/api/teacher/reservation/remove",
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
			reservationIds.push(reservations[i].id);
		}
	}
	var payload = {
		reservation_ids: reservationIds,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/teacher/reservation/cancel",
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
			<button type='button' onclick='$(\".admin_chakan\").remove();'>取消</button>\
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
		url: "/api/teacher/reservation/make",
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
		url: "/api/teacher/reservation/feedback/get",
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
		<div class='fankui_tch' id='feedback_table_" + index + "' style='text-align:left; top:100px; height:420px; width:200px; left:100px'>\
			咨询师反馈表<br>\
			您的姓名：<input id='teacher_fullname'/><br>\
			工作证号：<input id='teacher_username'/><br>\
			来访者姓名：<input id='student_fullname'/><br>\
			来访者问题描述：<br>\
			<textarea id='problem' style='padding:0px;height:60px;'></textarea><br>\
			咨询师提供的问题解决方法：<br>\
			<textarea id='solution' style='padding:0px;height:60px;'></textarea><br>\
			对中心的工作建议：<br>\
			<textarea id='advice' style='padding:0px;height:60px;'></textarea><br>\
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
		url: "/api/teacher/reservation/feedback/submit",
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
	apiGetStudent(index, showStudent);
}

function apiGetStudent(index, succCallback) {
	var payload = {
		reservation_id: reservations[index].id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/teacher/reservation/student/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				succCallback(data.payload.student);
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

function sendSms() {
	var checkedIndex = -1;
	for (var i = 0; i < reservations.length; i++) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			if (checkedIndex == -1) {
				checkedIndex = i;
			} else {
				alert("不能选中多个预约发送短信");
				return;
			}
		}
	}
	if (checkedIndex != -1) {
		if (reservations[checkedIndex].status != 2 && reservations[checkedIndex].status != 3) {
			alert("无法给未预约的咨询发送短信");
			return;
		}
	}
	$("body").append("\
		<div class='send_sms_teacher_pre'>\
			自定义发送短信\
			<br>\
			　手机号：<input id='mobile' type='tel' style='width:300px;'><br><br>\
			短信内容：<textarea id='content' style='width:300px;'></textarea><br>\
			<button type='button' onclick='sendSmsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".send_sms_teacher_pre\").remove();'>取消</button>\
		</div>\
	");
	if (checkedIndex != -1) {
		var setMobile = function(student) {
			$('#mobile').val(student.mobile);
		};
		apiGetStudent(checkedIndex, setMobile);
	}
	$('#content').text('有任何问题欢迎联系学习发展中心，62792453');
	optimize(".send_sms_teacher_pre");
}

function sendSmsConfirm() {
	var mobile = $("#mobile").val();
	if (mobile === "") {
		alert("手机为空");
		return;
	}
	var content = $("#content").val();
	if (content === "") {
		alert("内容为空");
		return;
	}
	var payload = {
		mobile: mobile,
		content: content,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/teacher/sms/send",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				$(".send_sms_teacher_pre").remove();
				alert("发送成功");
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function editSmsSuffix() {
	$("body").append("\
		<div class='edit_sms_suffix_pre'>\
			编辑预约短信后缀\
			<br>\
			中文后缀：<textarea id='sms_suffix' style='width:300px;'></textarea><br><br>\
			英文后缀：<textarea id='sms_suffix_en' style='width:300px;'></textarea><br>\
			<button type='button' onclick='editSmsSuffixConfirm();'>确认</button>\
			<button type='button' onclick='$(\".edit_sms_suffix_pre\").remove();'>取消</button>\
		</div>\
	");
	$('#sms_suffix').text(teacher.sms_suffix);
	$('#sms_suffix_en').text(teacher.sms_suffix_en);
	optimize(".edit_sms_suffix_pre");
}

function editSmsSuffixConfirm() {
	var smsSuffix = $('#sms_suffix').val();
	var smsSuffixEn = $('#sms_suffix_en').val();
	var payload = {
		sms_suffix: smsSuffix,
		sms_suffix_en: smsSuffixEn,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/teacher/sms_suffix/update",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				$(".edit_sms_suffix_pre").remove();
				alert("编辑成功");
				viewReservations();
			} else {
				alert(data.err_msg);
			}
		}
	});
}
