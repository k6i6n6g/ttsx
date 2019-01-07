package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ttsx/ttsx/models"
	"regexp"
	"github.com/astaxie/beego/utils"
	"strconv"
)

type UserController struct {
	beego.Controller
}
//展示注册页面
func (this *UserController)ShowRegister(){
	this.TplName="register.html"
}
//处理登陆页面
func (this *UserController)HandleRegister(){
	user_name:=this.GetString("user_name")
	pwd:=this.GetString("pwd")
	cpwd:=this.GetString("cpwd")
	email:=this.GetString("email")
	if email==""||cpwd==""||pwd==""||user_name==""{
		this.Data["errmsg"]="信息收集不全"
		this.TplName="register.html"
		return
	}
	//获取邮箱正则
	reg,err:=regexp.Compile(`^[A-Za-z\d]+([-_.][A-Za-z\d]+)*@([A-Za-z\d]+[-.])+[A-Za-z\d]{2,4}$`)
	if err!=nil{
		this.Data["errmsg"]="正则创建失败"
		this.TplName="register.html"
		return
	}
	//获取邮箱格式是否正确
	res:=reg.MatchString(email)
	if res==false{
		this.Data["errmsg"]="邮箱格式不正确"
		this.TplName="register.html"
		return
	}
	//比较密码对不对
	if pwd!=cpwd{
		this.Data["errmsg"]="密码对比格式不正确"
		this.TplName="register.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Pwd=pwd
	user.Email=email
	user.UserName=user_name
	_,err=o.Insert(&user)
	if err !=nil{
		this.Data["errmsg"]="导入信息错误"
		this.TplName="register.html"
		return
	}
	//this.Redirect("/login",30
	//注册成功时侯发送激活邮件  发送的邮箱                  邮箱的密钥                  服务器地址      端口属性
  config:=`{"username":"858744867@qq.com","password":"mthtelphapjvbejj","host":"smtp.qq.com","port":587}`
  //邮箱对象  邮件管理器
  emailSend:=utils.NewEMail(config)
  emailSend.From="858744867@qq.com"
  emailSend.To=[]string{email}
  //题目
  emailSend.Subject="天天生鲜用户激活"
  //内容     发了一个链接
  emailSend.HTML=`<a href="http://192.168.109.137:8000/active?userId=`+strconv.Itoa(user.Id)+`">点击激活</a>`
  //发送
  err=emailSend.Send()

 //在页面上显示
  this.Ctx.WriteString("注册成功，请前往邮箱激活")
}
//激活用户
func  (this *UserController)ActiveUser()  {
	//获取用户id
	userId,err:=this.GetInt("userId")
	if err!=nil{
		this.Data["errmsg"]="获取id错误"
		this.TplName="register.html"
		return
	}
	//更新usrId对应用户的active字段
	o:=orm.NewOrm()
	var user models.User
	user.Id=userId
	err=o.Read(&user)
	if err!=nil{
		this.Data["errmsg"]="激活失败，用户不存在"
		this.TplName="register.html"
		return
	}
	user.Active=1
	_,err=o.Update(&user)
	if err!=nil{
		this.Data["errmsg"]="激活失败，根新用户有问题"
		this.TplName="register.html"
		return
	}
	this.Redirect("/login",302)


}
//展示登陆页面
func (this *UserController)ShowLogin(){
//获取登陆页面的username
	userName:=this.Ctx.GetCookie("userName")
	if userName==""{
		this.Data["userName"]=""
		this.Data["checked"]=""
	}else {
		this.Data["userName"]=userName
		this.Data["checked"]="checked"
	}
	this.TplName="login.html"
}
//退出登陆
func(this *UserController)ShowLogout(){
	this.DelSession("userName")
	this.Redirect("/",302)
}
//处理登陆页面
func (this *UserController)HandleLogin(){
	username:=this.GetString("username")
	pwd:=this.GetString("pwd")
	if pwd==""||username==""{
		this.Data["errmsg"]="用户名和密码输入为空"
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.UserName=username
	user.Pwd=pwd
	err:=o.Read(&user,"username")
	if err!=nil{
		this.Data["errmsg"]="用户名和密码读取错误"
		this.TplName="login.html"
		return
	}
	if pwd!=user.Pwd{
		beego.Error("密码不匹配")
		this.TplName="login.html"
		return
	}
	if user.Active==0{
		beego.Error("用户名为激活，请先去邮箱激活")
		this.TplName="login.html"
		return
	}
	//记住用户名cookie
	remember:=this.GetString("remember")
	if remember=="on"{
		this.Ctx.SetCookie("userName",username,3600)
	}else{
		this.Ctx.SetCookie("userName",username,-1)

	}
	//记住用户名
	this.SetSession("userName",username)
	this.Redirect("/",302)

}
//封装函数用于获取用户名
func GetUser(this *UserController){
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Data["userName"]=""
	}else{
		this.Data["userName"]=userName.(string)
	}
}
//展示用户中心详情
func (this *UserController) ShowUserCenterInfo(){
//获取用户名
	GetUser(this)
//获取当前用户的默认联系方式和默认地址
	o:=orm.NewOrm()
	var receiver models.Receiver
	//查询默认地址
	userName:=this.GetSession("userName")
	qs:=o.QueryTable("Receiver").RelatedSel("User").Filter("User__UserName",userName.(string))
	//获取默认的地址
	qs.Filter("IsDefault",true).One(&receiver)
	this.Data["receiver"]=receiver

	this.Layout="layout.html"
	this.TplName="user_center_info.html"
}
//展示用户中心订单
func(this *UserController)ShowUserCenterOrder(){
	GetUser(this)
	//这俩个是个组合 。进行页面拼接的
	this.Layout="layout.html"
	this.TplName="user_center_order.html"
}
//展示用户中心地址
func(this *UserController)ShowUserCenterSite(){
	GetUser(this)
    userName:=this.GetSession("userName")
    //获取信息  获取当前用户的默认地址信息
	o:=orm.NewOrm()
	var receiver models.Receiver
	//获取当前用户所有收件人
	qs:=o.QueryTable("Receiver").RelatedSel("User").Filter("User__UserName",userName.(string))
	qs.Filter("IsDefault",true).One(&receiver)
	//传递给前段
	this.Data["receiver"]=receiver

    this.Layout="layout.html"
	this.TplName="user_center_site.html"
}
//插入用户收件信息
func(this *UserController)HandleAddSite(){
	receiverName:=this.GetString("receiverName")
	addr:=this.GetString("addr")
	zipCode:=this.GetString("zipCode")
	phone:=this.GetString("phone")
	if phone==""||zipCode==""||addr==""||receiverName==""{
		beego.Error("输入信息不完整，请从新输入")
		this.Redirect("/goods/usercentersite",302)
		return
	}
	//电话号码1校验
	//邮箱格式校验
	//处理数据
	o:=orm.NewOrm()
	var receiver models.Receiver
	//给插入对象赋值
	receiver.Name=receiverName
	receiver.Phone=phone
	receiver.ZipCode=zipCode
	receiver.Addr=addr

	//获取对象
	userName:=this.GetSession("userName")
	//查询数据库，获取userduixiang
	var user models.User
	user.UserName=userName.(string)
	o.Read(&user,"UserName")
	receiver.User=&user

	//每次插入的地址为默认地址，
	var oldReceiver models.Receiver
	//每次当前用户是否有默认地址，查询当前用户的所有收件人的地址
	qs:=o.QueryTable("Receiver").RelatedSel("User").Filter("User__Id",user.Id)
	//查询是否有默认地址
	err:=qs.Filter("IsDefault",true).One(&oldReceiver)
	//如果查询到默认地址，把默认地址更新为非默认地址
	if err==nil{
		oldReceiver.IsDefault=false
		o.Update(&oldReceiver)
	}
	receiver.IsDefault=true
	o.Insert(&receiver)
	this.Redirect("/goods/usercentersite",302)
}

