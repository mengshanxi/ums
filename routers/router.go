package routers

import (
	"ums/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.SetStaticPath("/ums/images", "static/images")
	beego.SetStaticPath("/ums/css", "static/css")
	beego.SetStaticPath("/ums/js", "static/js")
	beego.SetStaticPath("/ums/font", "static/font")
	beego.SetStaticPath("/ums/lib", "static/lib")

	beego.Router("/", &controllers.MainController{})
	beego.Router("/ums", &controllers.MainController{})
	beego.Router("/ums/login", &controllers.UserController{}, "get:Login;post:DoLogin")
	beego.Router("/ums/logout", &controllers.UserController{}, "get:Logout")
	beego.Router("/ums/welcome", &controllers.MainController{}, "get:Welcome")
	beego.Router("/ums/user/add", &controllers.UserController{}, "get:Add;post:DoAdd")
	beego.Router("/ums/user/list", &controllers.UserController{}, "get:Read")
	beego.Router("/ums/users", &controllers.UserController{}, "get:GetUsers")

}
