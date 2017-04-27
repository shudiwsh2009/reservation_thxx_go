var width = $(window).width();
var height = $(window).height();
var reservations;
var reservationGroups;

function getReservationById(id) {
	for (var i = 0; i < reservations.length; i++) {
		if (reservations[i].id === id) {
			return reservations[i];
		}
	}
	return null;
}

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/api/student/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				console.log(data);
				reservations = data.payload.reservations;
				refreshDataTable(reservations);
				optimizeTable();
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function viewGroupedReservations() {
	$.ajax({
			type: "GET",
			async: false,
			url: "/api/student/reservation/view/group",
			dataType: "json",
			success: function(data) {
				if (data.status === "OK") {
					console.log(data);
					reservations = data.payload.reservations;
					reservationGroups = data.payload.reservation_groups;
					refreshDataTableForGroups(reservationGroups);
				} else {
					alert(data.err_msg);
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
		var id = reservations[i].id;
		$("#col_time").append("<div class='table_cell' id='cell_time_" + id + "'>"
			+ reservations[i].start_time.substr(2) + "-"
			+ reservations[i].end_time.split(" ")[1] + "</div>");
		$("#col_teacher").append("<div class='table_cell' id='cell_teacher_" + id + "' onclick='getTeacher(\"" + id +
			"\")'>" + reservations[i].teacher_fullname + "</div>");
		if (reservations[i].status === 1) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + id
				+ "'><button type='button' id='cell_status_b_" + id + "' onclick='makeReservation(\"" + id
				+ "\")'>预约</button></div>");
		} else if (reservations[i].status === 2) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + id
				+ "'><button type='button' id='cell_status_b_" + id + "' disabled='true'>已预约</button>"
				+ "</div>");
		} else if (reservations[i].status === 3) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + id
				+ "'><button type='button' id='cell_status_b_" + id + "' onclick='getFeedback(\"" + id
				+ "\")'>反馈</button></div>");
		}
	}
}

function refreshDataTableForGroups(reservationGroups) {
	for (var i = 0; i < reservationGroups.length; i++) {
		var group = reservationGroups[i];
		$("#page_maintable").append("\
			<div class='has_children' data-name='" + group.date + "' style='overflow: hidden;'>\
				<div class='children_title'>" + group.date + " " + group.teacher_num + "个咨询师开设" + group.total_reservation_count
					 + "个咨询，" + group.available_reservation_count + "个空闲" + "</div>\
				<div class='table_col children' id='col_time_" + group.date + "' style='background-color:white;'>\
					<div class='table_head table_cell'>时间</div>\
				</div>\
				<div class='table_col children' id='col_teacher_" + group.date + "' style='background-color:white;'>\
					<div class='table_head table_cell'>咨询师</div>\
				</div>\
				<div class='table_col children' id='col_status_" + group.date + "' style='background-color:white;'>\
					<div class='table_head table_cell'>状态</div>\
				</div>\
				<div class='clearfix children'></div>\
			</div>\
		");
		for (var j = 0; j < group.reservations.length; j++) {
			var id = group.reservations[j].id;
			$("#col_time_" + group.date).append("<div class='table_cell' id='cell_time_" + id + "'>"
				+ group.reservations[j].start_time.substr(2) + "-"
				+ group.reservations[j].end_time.split(" ")[1] + "</div>");
			$("#col_teacher_" + group.date).append("<div class='table_cell' id='cell_teacher_" + id + "' onclick='getTeacher(\"" + id +
				"\")'>" + group.reservations[j].teacher_fullname + "</div>");
			if (group.reservations[j].status === 1) {
				$("#col_status_" + group.date).append("<div class='table_cell' id='cell_status_" + id
					+ "'><button type='button' id='cell_status_b_" + id + "' onclick='makeReservation(\"" + id
					+ "\")'>预约</button></div>");
			} else if (group.reservations[j].status === 2) {
				$("#col_status_" + group.date).append("<div class='table_cell' id='cell_status_" + id
					+ "'><button type='button' id='cell_status_b_" + id + "' disabled='true'>已预约</button>"
					+ "</div>");
			} else if (group.reservations[j].status === 3) {
				$("#col_status_" + group.date).append("<div class='table_cell' id='cell_status_" + id
					+ "'><button type='button' id='cell_status_b_" + id + "' onclick='getFeedback(\"" + id
					+ "\")'>反馈</button></div>");
			}
		}
	}
	$(function() {
    $('.has_children').click(function() {
      $(this).addClass('highlight').children('.table_col').show().end()
          .siblings().removeClass('highlight').children('.table_col').hide();
			optimizeGroup($(this).attr('data-name'));
    });
  });
}

function optimize(t){
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
}

function optimizeTable() {
	$("#col_time").width(width * 0.48);
	$("#col_teacher").width(width * 0.22);
	$("#col_status").width(width * 0.24);
	$("#col_time").css("margin-left", width * 0.02 + "px");

	for (var i = 0; i < reservations.length; ++i) {
		var id = reservations[i].id;
		var maxHeight = Math.max(
			$("#cell_time_" + id).height(),
			$("#cell_teacher_" + id).height(),
			$("#cell_status_" + id).height()
		);

		$("#cell_time_" + id).height(maxHeight);
		$("#cell_teacher_" + id).height(maxHeight);
		$("#cell_status_" + id).height(maxHeight);
		if (i % 2 == 1) {
			$("#cell_time_" + id).css("background-color", "white");
			$("#cell_teacher_" + id).css("background-color", "white");
			$("#cell_status_" + id).css("background-color", "white");
		} else {
			$("#cell_time_" + id).css("background-color", "#f3f3ff");
			$("#cell_teacher_" + id).css("background-color", "#f3f3ff");
			$("#cell_status_" + id).css("background-color", "#f3f3ff");
		}
	}
}

function optimizeGroup(groupName) {
	$("#col_time_" + groupName).width(width * 0.48);
	$("#col_teacher_" + groupName).width(width * 0.22);
	$("#col_status_" + groupName).width(width * 0.24);
	$("#col_time_" + groupName).css("margin-left", width * 0.02 + "px");

	for (var i = 0; i < reservationGroups.length; ++i) {
		if (reservationGroups[i].date !== groupName) {
			continue;
		}
		for (var j = 0; j < reservationGroups[i].reservations.length; ++j) {
			var id = reservationGroups[i].reservations[j].id;
			var maxHeight = Math.max(
				$("#cell_time_" + id).height(),
				$("#cell_teacher_" + id).height(),
				$("#cell_status_" + id).height()
			);

			$("#cell_time_" + id).height(maxHeight);
			$("#cell_teacher_" + id).height(maxHeight);
			$("#cell_status_" + id).height(maxHeight);
			if (j % 2 == 1) {
				$("#cell_time_" + id).css("background-color", "white");
				$("#cell_teacher_" + id).css("background-color", "white");
				$("#cell_status_" + id).css("background-color", "white");
			} else {
				$("#cell_time_" + id).css("background-color", "#f3f3ff");
				$("#cell_teacher_" + id).css("background-color", "#f3f3ff");
				$("#cell_status_" + id).css("background-color", "#f3f3ff");
			}
		}
	}
}

function makeReservation(id) {
	$("body").append("\
		<div class='yuyue_stu_pre'>\
			确定预约后请准确填写个人信息，方便咨询中心老师与你取得联系。\
			<br>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();makeReservationData(\"" + id + "\");'>\
				立即预约</button>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();'>暂不预约</button>\
		</div>\
	");
	optimize(".yuyue_stu_pre");
}

function makeReservationData(id) {
	$("body").append("\
		<div class='yuyue_stu' id='make_reservation_data_" + id + "' style='text-align:left;height:370px;overflow:scroll'>\
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
			<button type='button' onclick='makeReservationConfirm(\"" + id + "\");'>确定</button>\
			<button type='button' onclick='$(\".yuyue_stu\").remove();'>取消</button>\
		</div>\
	");
	optimize(".yuyue_stu");
}

function makeReservationConfirm(id) {
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
		reservation_id: id,
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
		url: "/api/student/reservation/make",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				makeReservationSuccess(id);
			} else {
				alert(data.err_msg);
			}
		},
	});
}

function makeReservationSuccess(id) {
	$(".yuyue_stu").remove();
	$("#cell_status_b_" + id).attr("disabled", "true");
	$("#cell_status_b_" + id).text("已预约");
	$("body").append("\
		<div class='yuyue_stu_success'>\
			你已预约成功，<br>\
			请关注短信提醒。<br>\
			<button type='button' onclick='$(\".yuyue_stu_success\").remove();viewReservations();'>确定</button>\
		</div>\
	");
	optimize(".yuyue_stu_success");
}

function getFeedback(id) {
    $("body").append("\
		<div class='fankui_stu_pre'>\
			请输入预约学号\
			<br>\
			<input id='fb_student_username'/>\
			<br>\
			<button type='button' onclick='getFeedbackConfirm(\"" + id + "\");'>确定</button>\
			<button type='button' onclick='$(\".fankui_stu_pre\").remove();'>取消</button>\
		</div>\
	");
    optimize(".fankui_stu_pre");
}

function getFeedbackConfirm(id) {
	var username = $("#fb_student_username").val();
	if (username === "") {
		alert("学号为空");
		return;
	}
	var payload = {
		reservation_id: id,
		username: username,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/student/reservation/feedback/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			console.log(data);
			if (data.status === "OK") {
				showFeedback(id, data.payload.feedback, username);
			} else {
				alert(data.err_msg);
			}
		},
	});
}

function showFeedback(id, feedback, username) {
    $(".fankui_stu_pre").remove();
    $('body').append('\
			<div class="fankui_stu" id="fankui_stu_'+ id +'" style="text-align:left;font-size:8px;position:absolute;height:606px;top:110px;left:5px;">\
			<div style="text-align:center;font-size:23px">咨询效果反馈问卷</div><br>\
				姓名：<input id="fullname"/><br>\
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
				<button type="button" onclick="submitFeedback(\'' + id + '\',\'' + username + '\')">提交反馈</button>\
				<button type="button" onclick="$(\'.fankui_stu\').remove();">取消</button></div>\
			</div>\
		');
		$(".fankui_stu").css("top", $(document).scrollTop() + 100);
	  $("#fullname").val(feedback.fullname);
	  $("#problem").val(feedback.problem);
	  $("#score").val(feedback.score);
	  $("#feedback").val(feedback.feedback);
	  for (var i = 1; i <= 12 && i <= feedback.choices.length; i++){
	      var t = feedback.choices[i - 1];
	      $("#fb_q"+i).val(t);
	  }
}

function submitFeedback(id, username) {
	var choices = "";
  for (var i = 1; i <= 12; i++) {
      choices += $("#fb_q" + i).val();
  }
	var fullname = $("#fullname").val();
	if (fullname === "") {
		alert("姓名为空");
		return;
	}
	var problem = $("#problem").val();
	if (problem === "") {
		alert("咨询问题为空");
		return;
	}
	var score = $("#score").val();
	if (score === "") {
		alert("分数为空");
		return;
	}
	var feedback = $("#feedback").val();
	if (feedback === "") {
		alert("建议为空");
		return;
	}
	var payload = {
		reservation_id: id,
    fullname: fullname,
    problem: problem,
    choices: choices,
		score: score,
    feedback: feedback,
    username: username,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/student/reservation/feedback/submit",
		data: payload,
		traditional: true,
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
	$(".fankui_stu").remove();
	$("body").append("\
		<div class='fankui_stu_success'>\
			您已成功提交反馈！<br>\
			<button type='button' onclick='$(\".fankui_stu_success\").remove();'>确定</button>\
		</div>\
	");
	optimize(".fankui_stu_success");
}

function getTeacher(id) {
	var payload = {
		reservation_id: id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/student/reservation/teacher/get",
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
		<div class='yuyue_stu_pre' style='text-align: left;'>\
			姓　　名：" + teacher.fullname + "<br>\
			性　　别：" + teacher.gender + "<br>\
			专业背景：" + teacher.major + "<br>\
			学　　历：" + teacher.academic + "<br>\
			资　　质：" + teacher.aptitude + "<br>\
			可咨询的问题：<br>\
			" + teacher.problem + "\
			<br>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();'>关闭</button>\
		</div>\
	");
	optimize(".yuyue_stu_pre");
}
