var width = $(window).width();
var height = $(window).height();
var reservations;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/student/reservation/view",
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
			<div style='text-align:center;font-size:23px'>咨询申请表</div><br>\
			姓　　名：<input id='name'/><br>\
			性　　别：<select id='gender'><option value=''>请选择</option><option value='男'>男</option><option value='女'>女</option></select><br>\
			学　　号：<input id='student_id'/><br>\
			院　　系：<input id='school'/><br>\
			生 源 地：<input id='hometown'/><br>\
			手　　机：<input id='mobile'/><br>\
			邮　　箱：<input id='email'/><br>\
			以前曾做过学习发展咨询、职业咨询或心理咨询吗？<select id='experience'><option value=''>请选择</option><option value='是'>是</option><option value='否'>否</option></select><br>\
			请概括你最想要咨询的问题：<br>\
			<textarea id='problem'></textarea><br>\
			<button type='button' onclick='makeReservationConfirm(" + index + ");'>确定</button>\
			<button type='button' onclick='$(\".yuyue_stu\").remove();'>取消</button>\
		</div>\
	");
	optimize(".yuyue_stu");
}

function makeReservationConfirm(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
		name: $("#name").val(),
		gender: $("#gender").val(),
		student_id: $("#student_id").val(),
		school: $("#school").val(),
        hometown: $("#hometown").val(),
        mobile: $("#mobile").val(),
        email: $("#email").val(),
        experience: $("#experience").val(),
        problem: $("#problem").val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/student/reservation/make",
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
    $("body").append("\
		<div class='fankui_stu_pre'>\
			请输入预约学号\
			<br>\
			<input id='fb_student_id'/>\
			<br>\
			<button type='button' onclick='getFeedbackConfirm(" + index + ");'>确定</button>\
			<button type='button' onclick='$(\".fankui_stu_pre\").remove();'>取消</button>\
		</div>\
	");
    optimize(".fankui_stu_pre");
}

function getFeedbackConfirm(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
        student_id: $("#fb_student_id").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/student/reservation/feedback/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showFeedback(index, data.feedback, $("#fb_student_id").val());
			} else {
				alert(data.message);
			}
		},
	});
}

function showFeedback(index, feedback, studentId) {
    $(".fankui_stu_pre").remove();
    $('body').append('\
		<div class="fankui_stu" id="fankui_stu_'+ index +'" style="text-align:left;font-size:8px;position:absolute;height:606px;top:110px;left:5px;">\
		<div style="text-align:center;font-size:23px">咨询效果反馈问卷</div><br>\
			姓名：<input id="name"/><br>\
			咨询问题：<br><textarea id="problem" style="width:250px;margin-left:20px"></textarea><br>\
			我和咨询师对咨询目标的看法是一致的<select id="fb_q1"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我对自己有新的认识<select id="fb_q2"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我对如何解决面临的问题有了新的思路<select id="fb_q3"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			了解到有关这一问题的政策信息与知识<select id="fb_q4"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我很清楚接下来需要干什么<select id="fb_q5"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我掌握了认识自己的方法<select id="fb_q6"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我掌握了如何获取更多信息的方法<select id="fb_q7"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我掌握了如何提升自身能力的方法<select id="fb_q8"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我意识到要对自己的学习与发展负责<select id="fb_q9"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我将尝试将咨询中的收获应用于生活中<select id="fb_q10"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			通过本次咨询，让我对解决问题更有信心了<select id="fb_q11"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			我喜欢我的咨询师，下次还会来预约咨询<select id="fb_q12"><option value="">请选择</option><option value="A">非常同意</option><option value="B">一般</option><option value="C">不同意</option></select><br>\
			请为本次咨询打分（0~100）：<input id="score" style="width:50px;"/><br>\
			感受和建议：<br><textarea id="feedback" style="width:250px;margin-left:20px"></textarea><br>\
			<div style="text-align:center;">\
			<button type="button" onclick="submitFeedback(' + index + ',' + studentId + ');">提交反馈</button>\
			<button type="button" onclick="$(\'.fankui_stu\').remove();">取消</button></div>\
		</div>\
	');
    $("#name").val(feedback.name);
    $("#problem").val(feedback.problem);
    $("#score").val(feedback.score);
    $("#feedback").val(feedback.feedback);
    for (var i = 1; i <= 12 && i <= feedback.choices.length; i++){
        var t = feedback.choices[i - 1];
        $("#fb_q"+i).val(t);
    }
}

function submitFeedback(index, studentId) {
	var choices = "";
    for (var i = 1; i <= 12; i++) {
        choices += $("#fb_q" + i).val();
    }
	var payload = {
		reservation_id: reservations[index].reservation_id,
        name: $("#name").val(),
        problem: $("#problem").val(),
        choices: choices,
		score: $("#score").val(),
        feedback: $("#feedback").val(),
        student_id: studentId,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/student/reservation/feedback/submit",
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