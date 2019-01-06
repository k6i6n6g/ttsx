package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//主要放表的设计

type User struct {
	Id int
	UserName string `orm:"unique;size(100)"`   //用户名
	Pwd string `orm:"size(100)"`			   //密码
	Email string							   //邮箱
	Power int  `orm:"default(0)"` //0标识普通用户  1标识管理员用户  用户权限
	Active int `orm:"default(0)"` //0标识未激活，1标识激活    是否激活 邮箱激活

	Receivers []*Receiver `orm:"reverse(many)"`
}

type Receiver struct {
	Id int
	Name string			//收件人名字
	ZipCode string		//收件人邮编
	Addr string			//地址
	Phone string		//收件人联系方式
	IsDefault bool `orm:"default(false)"` //是否未默认收件人

	User *User `orm:"rel(fk)"`

}


func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/ttsx?charset=utf8")
	//注册表
	orm.RegisterModel(new(User),new(Receiver))
	//跑起来
	orm.RunSyncdb("default",false,true)
}