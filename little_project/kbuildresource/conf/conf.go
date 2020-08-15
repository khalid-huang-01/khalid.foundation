package conf

import (
	"bryson.foundation/kbuildresource/utils"
	"github.com/astaxie/beego"
)

var Conf struct{
	SQLPWD string
	SQLCONN string
	SQLUser string
	InstanceName string
}

func init() {
	Conf.SQLCONN = beego.AppConfig.String("sqlconn")
	Conf.SQLPWD = beego.AppConfig.String("sqlpwd")
	Conf.SQLUser = beego.AppConfig.String("sqluser")

	Conf.InstanceName = utils.CreateRandomString(8)
}