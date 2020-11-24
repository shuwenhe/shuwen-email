package controller

import (
	"net/http"
	"text/template"

	"github.com/shuwenhe/shuwen-shop/model"
	"github.com/shuwenhe/shuwen-shop/service"
	"github.com/shuwenhe/shuwen-shop/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	flag, _ := service.IsLogin(r) // Determine whether you have logged in
	if flag {
		GetPageBooksByPrice(w, r) // Go to homepage
	} else {
		username := r.PostFormValue("username") // Get username
		password := r.PostFormValue("password") // Get password
		user, _ := service.CheckUserNameAndPassword(username, string(password))
		if user.ID > 0 { // 用户名密码正确
			uuid := utils.CreateUUID()    // 生成UUID
			session := &model.Session{ // 创建一个session
			SessionID: uuid,
				UserName:  user.Username,
				UserID:    user.ID,
			}
			service.AddSession(session) // 将session写入数据库,数据库中就有标识
			cookie := http.Cookie{      // 创建一个cookie与session相关联
				Name:     "user",
				Value:    uuid,
				HttpOnly: true,
			}
			// 将cookie发送给浏览器
			http.SetCookie(w, &cookie)
			t := template.Must(template.ParseFiles("views/pages/user/login_success.html"))
			t.Execute(w, user)
		} else {
			t := template.Must(template.ParseFiles("views/pages/user/login.html"))
			t.Execute(w, "用户名或密码不正确")
		}
	}
}

// Logout 注销
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("user")
	if cookie != nil {
		cookieValue := cookie.Value // 获取cookie的value值
		service.DeleteSession(cookieValue)
		cookie.MaxAge = -1        // 设置cookie失效
		http.SetCookie(w, cookie) // 将修改之后的cookie发送给浏览器
	}
	GetPageBooksByPrice(w, r) // 去首页
}

func Regist(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")
	// TODO 密码加密
	// hasedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // 创建用户的时候要加密用户的密码
	user, _ := service.CheckUserName(username)
	if user.ID > 0 {
		t := template.Must(template.ParseFiles("views/pages/user/regist.html"))
		t.Execute(w, "用户名已存在")
	} else {
		service.SaveUser(username, password, email)
		t := template.Must(template.ParseFiles("views/pages/user/regist_success.html"))
		t.Execute(w, "用户名或密码不正确")
	}
}

// CheckUserName 通过发送Ajax请求验证用户名是否存在
func CheckUserName(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	user, _ := service.CheckUserName(username)
	if user.ID > 0 {
		w.Write([]byte("用户名已存在！"))
	} else {
		w.Write([]byte("<font style = 'color:green'>用户名不存在！</font>"))
	}
}
