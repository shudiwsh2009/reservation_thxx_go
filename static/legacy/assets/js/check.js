$(".scan").click(function(){
  $(this).parent().find(".lecturer_photo").click();
})
$(".lecturer_photo").change(function(){
  var new_url = $(this).val();
  $(this).parent().find(".photo_url").val(new_url);
})
$(".scan-file").click(function(){
  $(this).parent().find(".upload_file").click();
})
$(".upload_file").change(function(){
  var new_url = $(this).val();
  $(this).parent().find(".file_url").val(new_url);
})
$(".scan_pic_intro").click(function(){
  $(this).parent().find(".pic_intro").click();
})
$(".pic_intro").change(function(){
  var new_url = $(this).val();
  $(this).parent().find(".pic_intro_url").val(new_url);
})
$(".scan_pic_exam_desc").click(function(){
  $(this).parent().find(".pic_exam_desc").click();
})
$(".pic_exam_desc").change(function(){
  var new_url = $(this).val();
  $(this).parent().find(".pic_exam_desc_url").val(new_url);
})
$(".scan_pic_exam_tab").click(function(){
  $(this).parent().find(".pic_exam_tab").click();
})
$(".pic_exam_tab").change(function(){
  var new_url = $(this).val();
  $(this).parent().find(".pic_exam_tab_url").val(new_url);
})
$(".btn-submit").click(function(){
  var subject = $(this).parent().parent().find(".subject");
  var time_periods = $(this).parent().parent().find(".time_periods");
  if(!subject.val()){  //没有填写预约人姓名
    alert("请填写辅导科目！");
  }
  else if(!time_periods.val()){  //没有填写辅导内容
    alert("请填写可选时间段！");
  }
  else{
    $("#submit-lecture").click();
  }
})
$(".btn-submit-user").click(function(){
  var name = $(this).parent().parent().find(".name");
  var username = $(this).parent().parent().find(".username");
  var password = $(this).parent().parent().find(".password");
  if(!name.val()){  //没有填写讲师姓名
    alert("请填写讲师姓名！");
  }
  else if(!username.val()){  //没有填写用户名
    alert("请填写用户名！");
  }
  else if(!password.val()){  //没有填写密码
    alert("请填写密码！");
  }
  else{
    $("#submit-user").click();
  }
})
$(".btn-upload-file").click(function(){
  var filename = $(this).parent().parent().find(".filename");
  var file_url = $(this).parent().parent().find(".file_url");
  if(!filename.val()){  //没有填写文件名
    alert("请填写文件名！");
  }
  else if(!file_url.val()){  //没有选择文件
    alert("请选择文件！");
  }
  else{
    $("#upload-file").click();
  }
})

$(".btn-submit-appointment").click(function(){
  var app_name = $(this).parent().parent().find(".app_name");
  var app_phone = $(this).parent().parent().find(".app_phone");
  var app_email = $(this).parent().parent().find(".app_email");
  var app_wechat = $(this).parent().parent().find(".app_wechat");
  var depart = $(this).parent().parent().find(".depart");
  var class_name = $(this).parent().parent().find(".class_name");
  var stu_num = $(this).parent().parent().find(".stu_num");
  var room = $(this).parent().parent().find(".room");
  var content = $(this).parent().parent().find(".content");
  if(!app_name.val()){  //没有填写预约人姓名
    alert("请填写预约人姓名！");
  }
  else if(!app_phone.val()){  //没有填写联系电话
    alert("请填写联系电话！");
  }
  else if(!app_email.val()){  //没有填写邮箱
    alert("请填写邮箱！");
  }
  else if(!depart.val()){  //没有填写院系
    alert("请填写院系！");
  }
  else if(!class_name.val()){  //没有填写班级
    alert("请填写班级！");
  }
  else if(!stu_num.val()){  //没有填写人数
    alert("请填写人数！");
  }
  else if(!content.val()){  //没有填写辅导内容
    alert("请填写辅导内容！");
  }
  else if(stu_num.val() < 8){  //人数不足8人
    alert("参加人数不足8人时无法预约！");
  }
  else{
	$(this).parent().find(".submit-appointment").click();
    //$("#submit-appointment").click();
  }
})
$(".btn-add-appointment").click(function(){
  var subject = $(this).parent().parent().find(".subject");
  var app_name = $(this).parent().parent().find(".app_name");
  var app_phone = $(this).parent().parent().find(".app_phone");
  var app_email = $(this).parent().parent().find(".app_email");
  var app_wechat = $(this).parent().parent().find(".app_wechat");
  var depart = $(this).parent().parent().find(".depart");
  var class_name = $(this).parent().parent().find(".class_name");
  var stu_num = $(this).parent().parent().find(".stu_num");
  var time_periods = $(this).parent().parent().find(".time_periods");
  var content = $(this).parent().parent().find(".content");
  if(!subject.val()){  //没有填写辅导科目
    alert("请填写辅导科目！");
  }
  else if(!app_name.val()){  //没有填写预约人姓名
    alert("请填写预约人姓名！");
  }
  else if(!app_phone.val()){  //没有填写联系电话
    alert("请填写联系电话！");
  }
  else if(!app_email.val()){  //没有填写邮箱
    alert("请填写邮箱！");
  }
  else if(!depart.val()){  //没有填写院系
    alert("请填写院系！");
  }
  else if(!class_name.val()){  //没有填写班级
    alert("请填写班级！");
  }
  else if(!stu_num.val()){  //没有填写人数
    alert("请填写人数！");
  }
  else if(!time_periods.val()){  //没有填写开始时间
    alert("请填写可选时间段！");
  }
  else if(!content.val()){  //没有填写辅导内容
    alert("请填写辅导内容！");
  }
  else if(stu_num.val() < 8){  //人数不足8人
    alert("参加人数不足8人时无法预约！");
  }
  else{
    $("#add-appointment").click();
  }
})
$(".check").click(function (){
  if($(this).text() == "编辑"){
    $(this).text("收起");
  }
  else{
    $(this).text("编辑");
  }
})
$(".tab_collapse").click(function (){
  if($(this).text() == "预约"){
    $(this).text("收起");
  }
  else{
    $(this).text("预约");
  }
})
