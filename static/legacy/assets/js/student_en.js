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
		data: {
			"language": "en_us",
			"student": getUrlVars()['student'],
		},
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
        data: {
            "language": "en_us",
            "student": getUrlVars()['student'],
        },
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
			<div class='table_head table_cell'>Time</div>\
		</div>\
		<div class='table_col' id='col_teacher' style='background-color:white;'>\
			<div class='table_head table_cell'>Advisor</div>\
		</div>\
		<div class='table_col' id='col_status' style='background-color:white;'>\
			<div class='table_head table_cell'>Status</div>\
		</div>\
		<div class='clearfix'></div>\
	";
	for (var i = 0; i < reservations.length; ++i) {
		var id = reservations[i].id;
		$("#col_time").append("<div class='table_cell' id='cell_time_" + id + "'>"
			+ reservations[i].start_time.substr(2) + "-"
			+ reservations[i].end_time.split(" ")[1] + "</div>");
        $("#col_teacher").append("<div class='table_cell' id='cell_teacher_" + id
            + "'><button type='button' id='cell_teacher_b_" + id + "' onclick='getTeacher(\"" + id
            + "\")'>" + reservations[i].teacher_fullname_en + "</button></div>");
		if (reservations[i].status === 1) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + id
				+ "'><button type='button' id='cell_status_b_" + id + "' onclick='makeReservation(\"" + id
				+ "\")'>Available</button></div>");
		} else if (reservations[i].status === 2) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + id
				+ "'><button type='button' id='cell_status_b_" + id + "' disabled='true'>Booked</button>"
				+ "</div>");
		} else if (reservations[i].status === 3) {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + id
				+ "'><button type='button' id='cell_status_b_" + id + "' onclick='getFeedback(\"" + id
				+ "\")'>Feedback</button></div>");
		}
	}
}

function refreshDataTableForGroups(reservationGroups) {
	for (var i = 0; i < reservationGroups.length; i++) {
		var group = reservationGroups[i];
		$("#page_maintable").append("\
			<div class='has_children' data-name='" + group.date + "' style='overflow: hidden;'>\
				<div class='children_title'>" + group.date + " " + group.teacher_num + " consultants and " + group.total_reservation_count
					 + " sessions, " + group.available_reservation_count + " available" + "</div>\
				<div class='table_col children' id='col_time_" + group.date + "' style='background-color:white;'>\
					<div class='table_head table_cell'>Time</div>\
				</div>\
				<div class='table_col children' id='col_teacher_" + group.date + "' style='background-color:white;'>\
					<div class='table_head table_cell'>Advisor</div>\
				</div>\
				<div class='table_col children' id='col_status_" + group.date + "' style='background-color:white;'>\
					<div class='table_head table_cell'>Status</div>\
				</div>\
				<div class='clearfix children'></div>\
			</div>\
		");
		for (var j = 0; j < group.reservations.length; j++) {
			var id = group.reservations[j].id;
			$("#col_time_" + group.date).append("<div class='table_cell' id='cell_time_" + id + "'>"
				+ group.reservations[j].start_time.substr(2) + "-"
				+ group.reservations[j].end_time.split(" ")[1] + "</div>");
            $("#col_teacher_" + group.date).append("<div class='table_cell' id='cell_teacher_" + id
                + "'><button type='button' id='cell_teacher_b_" + id + "' onclick='getTeacher(\"" + id
                + "\")'>" + group.reservations[j].teacher_fullname_en + "</button></div>");
			if (group.reservations[j].status === 1) {
				$("#col_status_" + group.date).append("<div class='table_cell' id='cell_status_" + id
					+ "'><button type='button' id='cell_status_b_" + id + "' onclick='makeReservation(\"" + id
					+ "\")'>Available</button></div>");
			} else if (group.reservations[j].status === 2) {
				$("#col_status_" + group.date).append("<div class='table_cell' id='cell_status_" + id
					+ "'><button type='button' id='cell_status_b_" + id + "' disabled='true'>Booked</button>"
					+ "</div>");
			} else if (group.reservations[j].status === 3) {
				$("#col_status_" + group.date).append("<div class='table_cell' id='cell_status_" + id
					+ "'><button type='button' id='cell_status_b_" + id + "' onclick='getFeedback(\"" + id
					+ "\")'>Feedback</button></div>");
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
			Please fill in your personal information accurately for further contacts from CSLD staff.\
			<br>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();makeReservationData(\"" + id + "\");'>\
				Book Now</button>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();'>Not Now</button>\
		</div>\
	");
	optimize(".yuyue_stu_pre");
}

function makeReservationData(id) {
	$("body").append("\
		<div class='yuyue_stu' id='make_reservation_data_" + id + "' style='text-align:left;height:370px;overflow:scroll'>\
			<div style='text-align:center;font-size:23px'>Advising Application Form</div><br>\
			Name          :<input id='fullname'/><br>\
			Gender        :<select id='gender'><option value=''>choose male/female</option><option value='male'>male</option><option value='female'>female</option></select><br>\
			Student ID No.:<input id='username'/><br>\
			Department    :<input id='school'/><br>\
			Nationality   :<input id='hometown'/><br>\
			Mobile No.    :<input id='mobile'/><br>\
			Email         :<input id='email'/><br>\
			Have you ever accepted any advising service on academic development or career development or any mental health consulting ?<select id='experience'><option value=''>choose yes/no</option><option value='是'>yes</option><option value='否'>no</option></select><br>\
			Please briefly tell us your problems to be solved.<br>\
			<textarea id='problem'></textarea><br>\
			<button type='button' onclick='makeReservationConfirm(\"" + id + "\");'>Confirm</button>\
			<button type='button' onclick='$(\".yuyue_stu\").remove();'>Cancel</button>\
		</div>\
	");
	optimize(".yuyue_stu");
}

function makeReservationConfirm(id) {
	var fullname = $("#fullname").val();
	if (fullname === "") {
		alert("Name cannot be empty!");
		return;
	}
	var gender = $("#gender").val();
	if (gender === "") {
		alert("Gender cannot be empty!");
		return;
	}
	var username = $("#username").val();
	if (username === "") {
		alert("Student ID No. cannot be empty!");
		return;
	}
	var school = $("#school").val();
	if (school === "") {
		alert("Department cannot be empty!");
		return;
	}
	var hometown = $("#hometown").val();
	if (hometown === "") {
		alert("Nationality cannot be empty!");
		return;
	}
	var mobile = $("#mobile").val();
	if (mobile === "") {
		alert("Mobile No. cannot be empty!");
		return;
	}
	var email = $("#email").val();
	if (email === "") {
		alert("Email cannot be empty!");
		return;
	}
	var experience = $("#experience").val();
	if (experience === "") {
		alert("Consulting Experience cannot be empty!");
		return;
	}
	var problem = $("#problem").val();
	if (problem === "") {
		alert("Main Problem cannot be empty!");
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
	$("#cell_status_b_" + id).text("Booked");
	$("body").append("\
		<div class='yuyue_stu_success'>\
			Success! Please pay attention to SMS alerts!<br>\
			<button type='button' onclick='$(\".yuyue_stu_success\").remove();viewReservations();'>Confirm</button>\
		</div>\
	");
	optimize(".yuyue_stu_success");
}

function getFeedback(id) {
    $("body").append("\
		<div class='fankui_stu_pre'>\
			Please fill in your Student ID No. for booking.\
			<br>\
			<input id='fb_student_username'/>\
			<br>\
			<button type='button' onclick='getFeedbackConfirm(\"" + id + "\");'>Confirm</button>\
			<button type='button' onclick='$(\".fankui_stu_pre\").remove();'>Cancel</button>\
		</div>\
	");
    optimize(".fankui_stu_pre");
}

function getFeedbackConfirm(id) {
	var username = $("#fb_student_username").val();
	if (username === "") {
		alert("Student ID No. cannot be empty!");
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
			<div style="text-align:center;font-size:23px">Advising Feedback Questionnaire</div><br>\
				Name:<input id="fullname"/><br>\
				What was your main problem?<br><textarea id="problem" style="width:250px;margin-left:20px"></textarea><br>\
				The advisor and I reached a consensus on the objectives of the advising.<select id="fb_q1"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I came up with some new observations on myself.<select id="fb_q2"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I had some new ideas on how to solve my problem.<select id="fb_q3"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I have learned some policies and information related to my problem.<select id="fb_q4"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I am very clear about my next step.<select id="fb_q5"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I have learned some ways to know myself.<select id="fb_q6"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I have learned some methods on how to get more information.<select id="fb_q7"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I have learned how to improve my ability.<select id="fb_q8"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I have realized that I need to be responsible for my study and development.<select id="fb_q9"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I am going to apply what I have learned from the advising to my life.<select id="fb_q10"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I become more confident in solving problems through the advising.<select id="fb_q11"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				I like my advisor and I would love to make an appointment in the future.<select id="fb_q12"><option value="">please choose</option><option value="A">Strongly agree</option><option value="B">Neutral</option><option value="C">Disagree</option></select><br>\
				Please rate this advising session (0-100):<input id="score" style="width:50px;"/><br>\
				Additional Comments:<br><textarea id="feedback" style="width:250px;margin-left:20px"></textarea><br>\
				<div style="text-align:center;">\
				<button type="button" onclick="submitFeedback(\'' + id + '\',\'' + username + '\')">Submit</button>\
				<button type="button" onclick="$(\'.fankui_stu\').remove();">Cancel</button></div>\
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
		alert("Name cannot be empty!");
		return;
	}
	var problem = $("#problem").val();
	if (problem === "") {
		alert("Main Problem cannot be empty!");
		return;
	}
	var score = $("#score").val();
	if (score === "") {
		alert("Score cannot be empty!");
		return;
	}
	var feedback = $("#feedback").val();
	if (feedback === "") {
		alert("Advice cannot be empty!");
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
			Success!<br>\
			<button type='button' onclick='$(\".fankui_stu_success\").remove();'>Confirm</button>\
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
			Name:" + teacher.fullname_en + "<br>\
			Gender:" + teacher.gender_en + "<br>\
			Major:" + teacher.major_en + "<br>\
			Academic Degree:" + teacher.academic_en + "<br>\
			Qualifications:" + teacher.aptitude_en + "<br>\
			Consultable Questions:<br>\
			" + teacher.problem_en + "\
			<br>\
			<button type='button' onclick='$(\".yuyue_stu_pre\").remove();'>OK</button>\
		</div>\
	");
	optimize(".yuyue_stu_pre");
}
