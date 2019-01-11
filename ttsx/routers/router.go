package routers

import (
	"ttsx/ttsx/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//登陆过滤
	beego.InsertFilter("/goods/*",beego.BeforeExec,filterFunc)
	//首页
    beego.Router("/", &controllers.GoodsController{},"get:ShowIndex")
    //注册页面
    beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	//用户注册激活
	beego.Router("/active", &controllers.UserController{},"get:ActiveUser")
    //登陆页面
    beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //退出页面
    beego.Router("/logout",&controllers.UserController{},"get:ShowLogout")
    //展示用户中心
    beego.Router("/goods/usercenterinfo", &controllers.UserController{},"get:ShowUserCenterInfo")
    //用户订单
    beego.Router("/goods/usercenterorder", &controllers.UserController{},"get:ShowUserCenterOrder")
    //收货地址
    beego.Router("/goods/usercentersite", &controllers.UserController{},"get:ShowUserCenterSite")
    //用户添加系统地址
    beego.Router("/goods/addSite", &controllers.UserController{},"post:HandleAddSite")
    //商品详情
    beego.Router("/detail",&controllers.GoodsController{},"get:ShowDetail")
   //商品用表
    beego.Router("/list",&controllers.GoodsController{},"get:ShowList")
    //购物车
    beego.Router("/goods/addcart",&controllers.CartController{},"post:HandleAddCart;get:ShowAddCart")
    //搜索
    beego.Router("/search",&controllers.GoodsController{},"post:HandleSearch")

}

func filterFunc(ctx *context.Context){
	userName:=ctx.Input.Session("userName")
	if userName ==nil{
	ctx.Redirect(302,"/login")
	return
	}
}