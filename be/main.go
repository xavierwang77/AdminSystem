package main

import (
	"admin_system/admin"
	"admin_system/user"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// 添加CORS中间件
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // 允许前端应用的域名
	})

	// 注册处理函数
	mux.HandleFunc("/uploadAvatarFile", admin.UploadAvatarFile)
	mux.HandleFunc("/changeAvatarFile", admin.ChangeAvatarFile)
	mux.HandleFunc("/uploadData", admin.UploadData)
	mux.HandleFunc("/verifyLoginData", admin.VerifyLoginData)
	mux.HandleFunc("/locateAvatar", admin.LocateAvatar)
	mux.HandleFunc("/locateUserAvatar", user.LocateUserAvatar)
	mux.HandleFunc("/fetchUserData", user.FetchUserData)
	mux.HandleFunc("/addUser", user.AddUser)
	mux.HandleFunc("/deleteUserByName", user.DeleteUserByName)
	mux.HandleFunc("/getOriginalUserData", user.GetOriginalUserData)
	mux.HandleFunc("/getOriginalAdminData", admin.GetOriginalAdminData)
	mux.HandleFunc("/changeData", user.ChangeUserData)
	mux.HandleFunc("/changeAdminData", admin.ChangeAdminData)
	mux.HandleFunc("/getOriginalAvatarName", user.GetOriginalAvatarName)
	mux.HandleFunc("/getOriginalAdminAvatarName", admin.GetOriginalAdminAvatarName)
	mux.HandleFunc("/getAvatarFilename", admin.GetAvatarFilename)
	mux.HandleFunc("/deleteAdmin", admin.DeleteAdmin)
	mux.HandleFunc("/changeLoginPrivilege", user.ChangeLoginPrivilege)
	mux.HandleFunc("/verifyUserLoginData", user.VerifyUserLoginData)
	mux.HandleFunc("/getLoginPrivilege", user.GetLoginPrivilege)

	// 使用CORS中间件包装处理器
	handler := c.Handler(mux)
	// http.HandleFunc("/upload", uploadFile)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Errorf("http.ListenAndServe()函数执行错误: %v", err)
	}
}
