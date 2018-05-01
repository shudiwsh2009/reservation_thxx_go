/*
 * 栏目树
 */
var Tree = function(o){
        this.root = document.getElementById(o);
        this.labels = this.root.getElementsByTagName('label');
        var that = this;
        this.int();
        Tree.prototype.addEvent(this.root,'click',function(e){that.treeshow(e)});
    }
    Tree.prototype = {
        int:function(){//初始化页面，加载时给有孩子结点的元素动态添加图标
            for(var i=0;i<this.labels.length;i++){
                var span = document.createElement('span');
                span.style.cssText ='display:inline-block;height:18px;vertical-align:middle;width:16px;cursor:pointer;';
                span.innerHTML = ' '
                span.className = 'add';
                if(this.nextnode(this.labels[i].nextSibling)&&this.nextnode(this.labels[i].nextSibling).nodeName == 'UL')
                    this.labels[i].parentNode.insertBefore(span,this.labels[i]);
                else
                    this.labels[i].className = 'rem'
            }
        },
        treeshow:function(e){
            e = e||window.event;
            var target = e.target||e.srcElement;
            var tp = this.nextnode(target.parentNode.nextSibling);
            switch(target.nodeName){
                case 'A'://点击A标签展开和收缩树形目录，并改变其样式
                    if(tp&&tp.nodeName == 'UL'){
                        if(tp.style.display != 'block' ){
                            tp.style.display = 'block';
                            this.prevnode(target.parentNode.previousSibling).className = 'ren'
                        }else{
                            tp.style.display = 'none';
                            this.prevnode(target.parentNode.previousSibling).className = 'add'
                        }    
                    }
                    break;
                case 'SPAN'://点击图标只展开或者收缩
                    var ap = this.nextnode(this.nextnode(target.nextSibling).nextSibling);
                    if(ap.style.display != 'block' ){
                        ap.style.display = 'block';
                        target.className = 'ren'
                    }else{
                        ap.style.display = 'none';
                        target.className = 'add'
                    }
                    break;
            }
        },
        addEvent:function(el,name,fn){//绑定事件
            if(el.addEventListener) return el.addEventListener(name,fn,false);
            return el.attachEvent('on'+name,fn);
        },
        nextnode:function(node){//寻找下一个兄弟并剔除空的文本节点
            if(!node)return ;
            if(node.nodeType == 1)
                return node;
            if(node.nextSibling)
                return this.nextnode(node.nextSibling);
        },
        prevnode:function(node){//寻找上一个兄弟并剔除空的文本节点
            if(!node)return ;
            if(node.nodeType == 1)
                return node;
            if(node.previousSibling)
                return prevnode(node.previousSibling);
        }
    }
    tree = new Tree("root");//实例化应用

/*
 * 图片上传
 */
function imageUploaded() {
    $(".imageUploaded").css("background","#70dbfa");
    $(".imageUploaded").text("选择成功");
}
