var width = $(window).width();
var height = $(window).height();

function optimize(t) {
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
}

function adminLogin() {
	var username = $("#username").val();
	if (username === "") {
		alert("工号为空");
		return;
	}
	var password = $("#password").val();
	if (password === "") {
		alert("密码为空");
		return;
	}
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/user/admin/login",
		data: {
			username: username,
			password: password,
		},
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				window.location.href = data.payload.redirect_url;
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function teacherLogin() {
	var username = $("#username").val();
	if (username === "") {
		alert("工号为空");
		return;
	}
	var password = $("#password").val();
	if (password === "") {
		alert("密码为空");
		return;
	}
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/user/teacher/login",
		data: {
			username: username,
			password: password,
		},
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				window.location.href = data.payload.redirect_url;
			} else {
				alert(data.err_msg);
			}
		}
	});
}

function logout() {
	$.ajax({
		type: "POST",
		async: false,
		url: "/api/user/logout",
		data: {},
		dataType: "json",
		success: function(data) {
			if (data.status === "OK") {
				window.location.href = data.payload.redirect_url;
			}
		},
	});
}
