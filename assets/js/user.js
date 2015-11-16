var width = $(window).width();
var height = $(window).height();

function optimize(t) {
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
}

function studentLogin() {
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/user/student/login",
		data: {
			username: $("#username").val(),
			password: $("#password").val(),
		},
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.location.href = data.url;
			} else {
				alert(data.message);
			}
		}
	});
}

function studentRegister() {
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/user/student/register",
		data: {
			username: $("#username").val(),
			password: $("#password").val(),
		},
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.location.href = data.url;
			} else {
				alert(data.message);
			}
		}
	});
}

function teacherLogin() {
	$.ajax({
		type: "POST",
		async: false,
		url: "/Reservation/user/teacher/login",
		data: {
			username: $("#username").val(),
			password: $("#password").val(),
		},
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.location.href = data.url;
			} else {
				alert(data.message);
			}
		}
	});
}

function logout() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/Reservation/user/logout",
		data: {},
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.location.href = data.url;
			}
		},
	});
}
