let app = {
    layerTime:1500,
    socketHost:'ws://' + window.location.host + "/admin/index/socket",
    post:function (url,data,func) {
        One.loader('show');
        $.ajax({
            type : "POST",
            dataType : "json",
            url : url,
            data:data,
            traditional:true,//防止深度序列化
            success: function (result) {
                One.loader('hide');
                if(func){
                    func(result);
                }else {
                    swal({title:result.message,text:false,type:(result.status === true) ? "success" :"error",timer:window.layerTime});
                    setTimeout(function(){
                        if(result && result.url) window.location.href = result.url
                    }, app.layerTime+500);
                }
            },
            error: function(data) {
                One.loader('hide');
                swal("请求错误","网络断开或服务器已停止","error")
            }
        });
    },
    cookie_menu:function () {
        $(".cookie_ajax").click(function () {
            var href = $(this).data("href");
            if(!href){
                swal("操作提示","当前分组无节点","error");
            }else {
                app.post($(this).data("action"),{id:$(this).data("id")},function (res) {
                    if(res && res.status === true && href){
                        window.location.href = href
                    }else {
                        swal("操作提示","当前分组无更多功能了","error")
                    }
                })
            }
        })
    },
    switch_status:function () { // 表格单个switch控制
        $(".switch_status").change(function() {
            const value = $(this).is(':checked') === true ? 1 : 0;
            const url   = $(this).data('action');
            const ids   = [$(this).data('ids')];
            const object=$(this);
            app.post(url,{ids:ids,status:value},function (res) {
                if(res.status === false){
                    object.prop("checked", !value);
                    One.helpers('notify', {align: 'center',type: "danger", icon: 'fa fa-check mr-1', message: res.message});
                }else {
                    One.helpers('notify', {align: 'center',"success": "danger", icon: 'fa fa-check mr-1', message: res.message});
                }
            });
        });
    },
    switch_status_all:function () { // 表格多个switch控制
        $(".switch_status_all").click(function () {
            const url   = $(this).data('action');
            const value   = $(this).data('value');
            let id_value = [];
            $(".checkList").each(function(){
                if($(this).is(':checked') === true){
                    id_value.push($(this).val());
                }
            });
            app.post(url,{ids:id_value,status:value});
        });
    },
    ajax_delete:function () { // 表格单个删除
        $(".ajax_delete").click(function () {
            const url   = $(this).data('action');
            const ids   = [$(this).data('ids')];
            swal({
                    title: "确定删除吗?",
                    text:"确定删除后，你将无法恢复该！" ,
                    type: "warning",
                    showCancelButton: true,
                    confirmButtonColor: "#DD6B55",
                    confirmButtonText: "确定操作",
                    closeOnConfirm: false,
                    showLoaderOnConfirm: true,
                },
                function(){app.post(url,{ids:ids})})
        });
    },
    ajax_deletes:function () { // 表格多个删除
        $(".ajax_deletes").click(function () {
            const url   = $(this).data('action');
            let id_value = [];
            $(".checkList").each(function(){
                if($(this).is(':checked') === true){
                    id_value.push($(this).val());
                }
            });
            swal({
                    title: "确定删除吗?",
                    text:"确定删除后，你将无法恢复该！" ,
                    type: "warning",
                    showCancelButton: true,
                    confirmButtonColor: "#DD6B55",
                    confirmButtonText: "确定操作",
                    closeOnConfirm: false,
                    showLoaderOnConfirm: true,
                },
                function(){
                    app.post(url,{ids:id_value})
                })
        });
    },
    checkAllTable:function () { // 表格全选反选
        $("#checkAll").on('click', function() {
            $("tbody .checkList:checkbox").prop("checked", $(this).prop('checked'));
        });
        $("tbody .checkList:checkbox").on('click', function() {
            if ($("tbody .checkList:checkbox").length === $("tbody .checkList:checked").length) {
                $("#checkAll").prop("checked", true);
            } else {
                $("#checkAll").prop("checked", false);
            }
        })
    }
};
$(document).ready(function(){
    app.cookie_menu();
});
