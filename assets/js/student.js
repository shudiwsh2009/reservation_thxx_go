var width = $(window).width();
var height = $(window).height();
var reservations;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/reservation/student/view",
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

function refreshDataTable(reservations) {
	$("#page_maintable")[0].innerHTML = "\
		<div class='table_col' id='col_time' style='background-color:white;'>\
			<div class='table_head table_cell'>时间</div>\
		</div>\
		<div class='table_col' id='col_teacher' style='background-color:white;'>\
			<div class='table_head table_cell'>咨询师</div>\
		</div>\
		<div class='table_col' id='col_status' style='background-color:white;'>\
			<div class='table_head table_cell'>状态</div>\
		</div>\
		<div class='clearfix'></div>\
	";			
	for (var i = 0; i < reservations.length; ++i) {
		$("#col_time").append("<div class='table_cell' id='cell_time_" + i + "'>"
			+ reservations[i].start_time.substr(2) + "-"
			+ reservations[i].end_time.split(" ")[1] + "</div>");
		$("#col_teacher").append("<div class='table_cell' id='cell_teacher_" + i + "'>"
			+ reservations[i].teacher_fullname + "</div>");
		if (reservations[i].status === "AVAILABLE") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i
				+ "'><button type='button' id='cell_status_b_" + i + "' onclick='makeReservation(" + i
				+ ")'>预约</button></div>");
		} else if (reservations[i].status === "RESERVATED") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i
				+ "'><button type='button' id='cell_status_b_" + i + "' disabled='true'>已预约</button>"
				+ "</div>");
		} else if (reservations[i].status === "FEEDBACK") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i
				+ "'><button type='button' id='cell_status_b_" + i + "' onclick='getFeedback(" + i
				+ ")'>反馈</button></div>");
		}
	}
}

function optimize(t){
	$("#col_time").width(width * 0.48);
	$("#col_teacher").width(width * 0.22);
	$("#col_status").width(width * 0.24);
	$("#col_time").css("margin-left", width * 0.02 + "px");

	for (var i = 0; i < reservations.length; ++i) {
		var maxHeight = Math.max(
			$("#cell_time_" + i).height(),
			$("#cell_teacher_" + i).height(),
			$("#cell_status_" + i).height()
		);

		$("#cell_time_" + i).height(maxHeight);
		$("#cell_teacher_" + i).height(maxHeight);
		$("#cell_status_" + i).height(maxHeight);
		if (i % 2 == 1) {
			$("#cell_time_" + i).css("background-color", "white");
			$("#cell_teacher_" + i).css("background-color", "white");
			$("#cell_status_" + i).css("background-color", "white");
		} else {
			$("#cell_time_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_" + i).css("background-color", "#f3f3ff");
			$("#cell_status_" + i).css("background-color", "#f3f3ff");
		}
	}
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
}

function makeReservation(index) {
	$("body").append("\
		<div class='yuyue_stu_pre'>\
			确定预约后请准确填写个人信息，方便心理咨询中心老师与你取得联系。\
			<br>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();makeReservationData(" + index + ");'>\
				立即预约</button>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();'>暂不预约</button>\
		</div>\
	");
	optimize(".yuyue_stu_pre");
}

function makeReservationData(index) {
	$("body").append("\
		<div class='yuyue_stu' id='make_reservation_data_" + index + "' style='text-align:left;height:370px;overflow:scroll'>\
			<div style='text-align:center;font-size:23px'>学生信息登记表</div><br>\
			姓　　名：<input id='student_fullname' value='" + student.student_fullname + "'><br>\
			性　　别：<select id='student_gender'><option value=''>请选择</option><option value='男'>男</option><option value='女'>女</option></select><br>\
			出生日期：<input id='student_birthday' value='" + student.student_birthday + "'><br>\
			系　　别：<input id='student_school' value='" + student.student_school + "'><br>\
			年　　级：<input id='student_grade' value='" + student.student_grade + "'><br>\
			现在住址：<input id='student_current_address' value='" + student.student_current_address + "'><br>\
			家庭住址：<input id='student_family_address' value='" + student.student_family_address + "'><br>\
			联系电话：<input id='student_mobile' value='" + student.student_mobile + "'><br>\
			邮　　箱：<input id='student_email' value='" + student.student_email + "'><br>\
			咨询经历：<br>\
			时间：<input id='student_experience_time' value='" + student.student_experience_time + "'><br>\
			地点：<input id='student_experience_location' style='width:60px' value='" + student.student_experience_location + "'>\
			咨询师：<input id='student_experience_teacher' style='width:60px' value='" + student.student_experience_teacher + "'><br>\
			父　　亲<br>\
			年龄：<input id='student_father_age' style='width:20px' value='" + student.student_father_age + "'> 职业：<input id='student_father_job' style='width:40px' value='" + student.student_father_job + "'> 学历：<input id='student_father_edu' style='width:40px' value='" + student.student_father_edu + "'><br>\
			母　　亲<br>\
			年龄：<input id='student_mother_age' style='width:20px' value='" + student.student_mother_age + "'> 职业：<input id='student_mother_job' style='width:40px' value='" + student.student_mother_job + "'> 学历：<input id='student_mother_edu' style='width:40px' value='" + student.student_mother_edu + "'><br>\
			父母婚姻状况：<select id='student_parent_marriage'><option value=''>请选择</option><option value='良好'>良好</option><option value='一般'>一般</option><option value='离婚'>离婚</option><option value='再婚'>再婚</option></select><br>\
			在近三个月里，是否发生了对你有重大意义的事（如亲友的死亡、法律诉讼、失恋等）？<br>\
			<textarea id='student_significant'></textarea><br>\
			你现在需要接受帮助的主要问题是什么？<br>\
			<textarea id='student_problem'></textarea><br>\
			<button type='button' onclick='makeReservationConfirm(" + index + ");'>确定</button>\
			<button type='button' onclick='$(\".yuyue_stu\").remove();'>取消</button>\
		</div>\
	");
	$("#student_gender").val(student.student_gender);
	$("#student_parent_marriage").val(student.student_parent_marriage);
	$("#student_significant").val(student.student_significant);
	$("#student_problem").val(student.student_problem);
	optimize(".yuyue_stu");
}

function makeReservationConfirm(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
		student_fullname: $("#student_fullname").val(),
		student_gender: $("#student_gender").val(),
		student_birthday: $("#student_birthday").val(),
		student_school: $("#student_school").val(),
		student_grade: $("#student_grade").val(),
		student_current_address: $("#student_current_address").val(),
		student_family_address: $("#student_family_address").val(),
		student_mobile: $("#student_mobile").val(),
		student_email: $("#student_email").val(),
		student_experience_time: $("#student_experience_time").val(),
		student_experience_location: $("#student_experience_location").val(),
		student_experience_teacher: $("#student_experience_teacher").val(),
		student_father_age: $("#student_father_age").val(),
		student_father_job: $("#student_father_job").val(),
		student_father_edu: $("#student_father_edu").val(),
		student_mother_age: $("#student_mother_age").val(),
		student_mother_job: $("#student_mother_job").val(),
		student_mother_edu: $("#student_mother_edu").val(),
		student_parent_marriage: $("#student_parent_marriage").val(),
		student_significant: $("#student_significant").val(),
		student_problem: $("#student_problem").val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/reservation/student/make",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				makeReservationSuccess(index);
			} else {
				alert(data.message);
			}
		},
	});
}

function makeReservationSuccess(index) {
	$(".yuyue_stu").remove();
	$("#cell_status_b_" + index).attr("disabled", "true");
	$("#cell_status_b_" + index).text("已预约");
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
		reservation_id: reservations[index].reservation_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/reservation/student/feedback/get",
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
		<div class='fankui_stu' id='feedback_table_" + index + "'\
			style='text-align:left; font-size:8px; position: absolute; height:350px;top:100px;left:5px;'>\
			学生反馈表<br>\
			1、你是否得到了你所希望的咨询？<br>\
			<select id='feedback_q1'>\
				<option value='4'>肯定是的</option>\
				<option value='3'>基本上是的</option>\
				<option value='2'>没有</option>\
				<option value='1'>肯定没有</option>\
			</select><br>\
			2、咨询在多大程度上满足了你的需要？<br>\
			<select id='feedback_q2'>\
				<option value='4'>几乎全部需要得到满足</option>\
				<option value='3'>大部分需要得到满足</option>\
				<option value='2'>仅一小部分需要得到满足</option>\
				<option value='1'>需要丝毫没有得到满足</option>\
			</select><br>\
			3、如果一个朋友需要咨询，你会向他或她推荐这位咨询师吗？<br>\
			<select id='feedback_q3'>\
				<option value='4'>肯定会</option>\
				<option value='3'>会</option>\
				<option value='2'>不会</option>\
				<option value='1'>肯定不会</option>\
			</select><br>\
			4、总体来讲，你对你接受的咨询有多满意？<br>\
			<select id='feedback_q4'>\
				<option value='4'>非常满意</option>\
				<option value='3'>大部分满意</option>\
				<option value='2'>无所谓，或不太满意</option>\
				<option value='1'>非常不满意</option>\
			</select><br>\
			5、如果你将再次寻求咨询，你会回来找这位咨询师吗？<br>\
			<select id='feedback_q5'>\
				<option value='4'>肯定会</option>\
				<option value='3'>会</option>\
				<option value='2'>不会</option>\
				<option value='1'>肯定不会</option>\
			</select><br>\
			<button type='button' onclick='submitFeedback(" + index + ");'>提交</button>\
			<button type='button' onclick='$(\".fankui_stu\").remove();'>取消</button>\
		</div>\
	");
	for (var i = 1; i <= feedback.scores.length; ++i) {
		$("#feedback_q" + i).val(feedback.scores[i - 1]);
	}
	optimize(".fankui_stu");
}

function submitFeedback(index) {
	var scores = [];
	for (var i = 1; i <= 5; ++i) {
		scores.push(parseInt($("#feedback_q" + i).val()));
	}
	var payload = {
		reservation_id: reservations[index].reservation_id,
		scores: scores,
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/reservation/student/feedback/submit",
		data: payload,
		traditional: true,
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
	$(".fankui_stu").remove();
	$("body").append("\
		<div class='fankui_stu_success'>\
			您已成功提交反馈！<br>\
			<button type='button' onclick='$(\".fankui_stu_success\").remove();'>确定</button>\
		</div>\
	");
	optimize(".fankui_stu_success");
}