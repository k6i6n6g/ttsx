package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/orm"
	"strconv"
	"ttsx/ttsx/models"
)

type CartController struct {
	beego.Controller
}

//添加购物车
func  (this * CartController)HandleAddCart(){
	//获取数据
	count,err1:=this.GetInt("count")
	goodsId,err2:=this.GetInt("goodsId")
	//校验数据
	if err1!=nil ||err2!=nil{
		beego.Error("ajax传递数据失败")
		return
	}
	//处理数据 判断是否有注册用户
	//1.有个数据储存json数据
	resp:=make(map[string]interface{})
	userName:=this.GetSession("userName")
	if userName==nil{
		resp["errno"]=1
		resp["errmsg"]="用户未登录"
		//把容器传递给前段
		this.Data["json"]=resp
		//指定接受方式
		this.ServeJSON()
		return
		}
	//处理数据
	//redis储存
	conn,err:=redis.Dial("tcp","192.168.109.138:6379")
	if err!=nil{
		resp["errno"]=2
		resp["errmsg"]="redis链接失败"
		//把容器传递给前段页面
		this.Data["json"]=resp
		//指定接受方式
		this.ServeJSON()
		return
	}
	defer conn.Close()
	conn.Do("hset","cart_"+userName.(string),goodsId,count)
	resp["errno"]=5
	resp["errmsg"]="ok"
	this.Data["json"]=resp
	this.ServeJSON()

}
//展示购物车
func (this * CartController)ShowAddCart() {
//获取数据 先检验一下用户是否登录了
userName :=this.GetSession("userName")
//username是接口，所以用nil
if userName==nil{
	beego.Error("用户未登录")
	this.Redirect("/",302)
	return
}
//若果是登录状态，需要从redis中获取数据
conn,err:=redis.Dial("tcp","192.168.109.138:6379")
if err!=nil{
	beego.Error("redis链接失败")
	this.Redirect("/",302)
	return
}
defer  conn.Close()
//读取数据
resp,err:=conn.Do("hgetall","cart_"+userName.(string))
if err!=nil{
	beego.Error("1",err)
	this.Redirect("/",302)
	return
}
//把空接口转换成map[string]int类型
cartMap,err:=redis.IntMap(resp,err)
	if err!=nil{
		beego.Error("2",err)
		this.Redirect("/",302)
		return
	}
//循环遍历从redis获取的数据
 o:=orm.NewOrm()
 var goods []map[string]interface{}
 var totalPrice,totalCount int
 for goodsId,value:=range cartMap{
 	temp:=make(map[string]interface{})
 	//查询商品信息     goodsId为string 所以要转换成int类型
 	id,_:=strconv.Atoi(goodsId)
 	var goodsSku models.GoodsSKU
 	goodsSku.Id=id
 	o.Read(&goodsSku)

 	//给每一行容器插入数据
 	temp["goodsSku"]=goodsSku
 	temp["count"]=value
 	//商品单价乘单个价值是一种商品的总价值
 	price:=goodsSku.Price * value
 	temp["price"]=price
 	totalPrice+=price
 	totalCount+=1
 	//把容器放在到大容器中
 	goods=append(goods,temp)
 }
 //大容器中有 goodsSku price count
 this.Data["goods"]=goods
 //所有商品价格总数
 this.Data["totalPrice"]=totalPrice
 //所有商品总数
 this.Data["totalCount"]=totalCount
 this.Data["userName"]=userName
 this.Layout="layout.html"
 this.TplName="cart.html"

}
