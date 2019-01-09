package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/ttsx/models"
	"github.com/gomodule/redigo/redis"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController)ShowIndex()  {
	userName:=this.GetSession("userName")
	if userName ==nil{
		this.Data["userName"]=""
	}else{
		this.Data["userName"]=userName.(string)
	}
	//获取首页内容
	o:=orm.NewOrm()
	//
	var goodsTypes []models.GoodsType
	o.QueryTable("goodsType").All(&goodsTypes)
	this.Data["goodsTypes"]=goodsTypes
    //
	var goodslunbo []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&goodslunbo)
	this.Data["goodslunbo"]=goodslunbo
    //
	var goodsPro []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&goodsPro)
	this.Data["goodsPro"]=goodsPro


	var goods []map[string]interface{}
	//把所有商品类型插入到大容器中
	for _,v:=range goodsTypes{
		temp:=make(map[string]interface{})
		temp["goodsType"]=v
		//把所有商品插入到大容器中
		goods=append(goods,temp)
	}
	//把类型对应的首页展示商品插入到大容器中
    for _,v:=range goods{
     	//获取到类型对应的所有商品                                                                                                  这是一个key对应的就是上面的v ，v就是类型
     	qs:=o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsSKU","GoodsType").Filter("GoodsType",v["goodsType"])
     	//需要把商品放到map中
     	var goodsText []models.IndexTypeGoodsBanner
     	qs.Filter("DisplayType",0).OrderBy("Index").All(&goodsText)
     	//获取图片商品
     	var goodsImage []models.IndexTypeGoodsBanner
     	qs.Filter("DisplayType",1).OrderBy("Index").All(&goodsImage)
     	//插入到大容器中
     	v["goodsText"]=goodsText
     	v["goodsImage"]=goodsImage

	}
     //把整个大容器都传给前段
	this.Data["goods"]=goods


	this.TplName="index.html"
}

func(this * GoodsController)ShowDetail(){
	goodsId,err:=this.GetInt("goodsId")
	if err!=nil{
		beego.Error("获取iD错误")
		this.Redirect("/",302)
		return
	}
	o:=orm.NewOrm()
	var  goodsSku models.GoodsSKU
	//查询
	o.QueryTable("GoodsSKU").RelatedSel("Goods","GoodsType").Filter("Id",goodsId).One(&goodsSku)
	//获取数据
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All((&goodsTypes))
	var newGoods []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType",goodsSku.GoodsType).OrderBy("Time").Limit(2,0).All(&newGoods)

	this.Data["newGoods"]=newGoods
	this.Data["goodsTypes"]=goodsTypes
	this.Data["goodsSku"]=goodsSku

	//添加历史浏览记录
	userName:=this.GetSession("userName")

	if userName!=nil{
		//redis list
		conn,err:=redis.Dial("tcp","192.168.109.138:6379")
		if err!=nil{
			beego.Error("redis链接失败",err)
			return
			}
		defer  conn.Close()
		conn.Do("lrem","history_"+userName.(string),0,goodsId)
		conn.Do("lpush","history_"+userName.(string),goodsId)
	}


	this.TplName="detail.html"

}