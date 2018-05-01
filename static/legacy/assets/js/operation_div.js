$(document).ready(function(){
	$("#li_login").click(function(){
		$("#li_main").removeClass("active");
		$("#li_login").addClass("active");
		$("#start-button").stop();
		$("#start-text").stop();
		$("#start-tsinghua").stop();
		
		$("#start-button").animate({
		    opacity:'0',
		    left:'-300px'
		  },500);
		$("#start-tsinghua").animate({
		    left:'-300px',
		    bottom:'-60px',
		    opacity:'0.8'
		  },500);
		$("#start-text").animate({
		    left:'-300px',
		    bottom:'-20px',
		    opacity:'0'
		  },500);
		$("#head-xs").hide();	
		$("#login-form").show();
		
	})

	$("#li_main").click(function(){
		$("#li_main").addClass("active");
		$("#li_login").removeClass("active");
		$("#start-button").stop();
		$("#start-text").stop();
		$("#start-tsinghua").stop();
		
		$("#start-tsinghua").animate({
		    left:'0px',
		    bottom:'0px',
		    opacity:'1'
		  },500);
		$("#start-text").animate({
		    left:'0px',
		    bottom:'0px',
		    opacity:'1'
		  },500);
		$("#login-form").hide();
		$("#head-xs").show();
		$("#start-button").animate({left:'0px'},500);
		$("#start-button").animate({opacity:'1'},600);
		
	})

	$("#start-button").click(function(){
		$("#li_main").removeClass("active");
		$("#li_login").addClass("active");
		$("#start-button").stop();
		$("#start-text").stop();
		$("#start-tsinghua").stop();
		$("#start-button").animate({
		    opacity:'0',
		    left:'-300px'
		  },500);				
		$("#start-tsinghua").animate({
		    left:'-300px',
		    bottom:'-60px',
		    opacity:'0.8'
		  },500);
		$("#start-text").animate({
		    left:'-300px',
		    bottom:'-20px',
		    opacity:'0'
		  },500);
		$("#login-form").fadeIn(200);
	})

	$("#start-button-xs").click(function(){	
		$("#li_main").removeClass("active");
		$("#li_login").addClass("active");
		$("#head-xs").hide();
		$("#login-form").show();
	})




	

    $(".datepicker").datetimepicker();
	$('.datepicker').datepicker({ dayNamesMin: ['日', '一', '二', '三', '四', '五', '六'] }); 
    var dayNamesMin = $('.datepicker').datepicker('option', 'dayNamesMin'); 
    $('.datepicker').datepicker('option', 'dayNamesMin', ['日', '一', '二', '三', '四', '五', '六'] );

    $('.datepicker').datepicker({ monthNames:['一月','二月','三月','四月','五月','六月','七月','八月','九月','十月','十一月','十二月']});
    var monthNames = $('.datepicker').datepicker('option', 'monthNames');
    $('.datepicker').datepicker('option', 'monthNames', ['一月','二月','三月','四月','五月','六月','七月','八月','九月','十月','十一月','十二月']);

    $(".btn-delete").css("opacity","0");

    // 判断图片加载状况，加载完成后回调
	$(".notice-pic").each(function(){
    	$(this).children(".notice-pic-pic").css("height","200px");
    	$(this).children(".notice-pic-pic").children("a").children("img").css("height","100%");
    	$(this).css("width",$(this).children(".notice-pic-pic").children("a").children("img").css("width"));
    	$(this).children(".notice-pic-pic").children(".btn-delete").css("width",$(this).children(".notice-pic-pic").children("a").children("img").css("width"));
    	$(this).children(".notice-pic-pic").children(".btn-delete").css("position","absolute");
    })

    
    
    $(".notice-pic").children(".notice-pic-pic").mouseover(function(){
    	$(this).parent().children("div").children(".btn-delete").animate({
		    opacity: 1
		},500);
    })
    $(".btn-delete").mouseover(function(){
    	$(this).animate({
		    opacity: 1
		},500);
    })
    $(".notice-pic").mouseleave(function(){
    	$(this).children(".btn-delete").stop(1,1);
    	$(this).children(".btn-delete").css("opacity","0");
    })

    // $("input[type='radio'][name='tpl']:checked").each(function(){
    // 	var val=$(this).val();
    // 	if(val==0){
    // 		$("#bg-meeting").show();
    // 		$("#screen-meeting").show();
    // 		$("#content-div").hide();
    // 		$("#content-div1").hide();
    // 	}
    // 	else if(val==1){
    // 		$("#bg-meeting").hide();
    // 		$("#screen-meeting").hide();
    // 		$("#content-div").hide();
    // 		$("#content-div1").hide();
    // 	}
    // 	else{
    // 		$("#bg-meeting").hide();
    // 		$("#screen-meeting").hide();
    // 		$("#content-div").show();
    // 		$("#content-div1").show();
    // 	}   	
    // })


    $("#optionsRadios1").click(function(){
    	$("#bg-meeting").removeClass('hide');
		$("#screen-meeting").removeClass('hide');
		$("#content-div").addClass('hide');
		$("#content-div1").addClass('hide');
    })
    $("#optionsRadios2").click(function(){
    	$("#bg-meeting").addClass('hide');
		$("#screen-meeting").addClass('hide');
		$("#content-div").addClass('hide');
		$("#content-div1").addClass('hide');
    })
    $("#optionsRadios3").click(function(){
    	$("#bg-meeting").addClass('hide');
		$("#screen-meeting").addClass('hide');
		$("#content-div").removeClass('hide');
		$("#content-div1").removeClass('hide');
    })


    // 判断图片加载的

    
})