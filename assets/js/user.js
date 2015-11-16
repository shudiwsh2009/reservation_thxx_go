var width = $(window).width();
var height = $(window).height();

function optimize(t) {
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
}

function login() {
	$.ajax({
		type: "POST",
		async: false,
		url: "/reservation/user/login",
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
		url: "/reservation/user/logout",
		data: {},
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.location.href = data.url;
			}
		},
	});
}
