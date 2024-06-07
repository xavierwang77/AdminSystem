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
	mux.HandleFunc("/api/uploadAvatarFile", admin.UploadAvatarFile)
	mux.HandleFunc("/api/changeAvatarFile", admin.ChangeAvatarFile)
	mux.HandleFunc("/api/uploadData", admin.UploadData)
	mux.HandleFunc("/api/verifyLoginData", admin.VerifyLoginData)
	mux.HandleFunc("/api/locateAvatar", admin.LocateAvatar)
	mux.HandleFunc("/api/locateUserAvatar", user.LocateUserAvatar)
	mux.HandleFunc("/api/fetchUserData", user.FetchUserData)
	mux.HandleFunc("/api/addUser", user.AddUser)
	mux.HandleFunc("/api/deleteUserByName", user.DeleteUserByName)
	mux.HandleFunc("/api/getOriginalUserData", user.GetOriginalUserData)
	mux.HandleFunc("/api/getOriginalAdminData", admin.GetOriginalAdminData)
	mux.HandleFunc("/api/changeData", user.ChangeUserData)
	mux.HandleFunc("/api/changeAdminData", admin.ChangeAdminData)
	mux.HandleFunc("/api/getOriginalAvatarName", user.GetOriginalAvatarName)
	mux.HandleFunc("/api/getOriginalAdminAvatarName", admin.GetOriginalAdminAvatarName)
	mux.HandleFunc("/api/getAvatarFilename", admin.GetAvatarFilename)
	mux.HandleFunc("/api/deleteAdmin", admin.DeleteAdmin)
	mux.HandleFunc("/api/changeLoginPrivilege", user.ChangeLoginPrivilege)
	mux.HandleFunc("/api/verifyUserLoginData", user.VerifyUserLoginData)
	mux.HandleFunc("/api/getLoginPrivilege", user.GetLoginPrivilege)

	// 使用CORS中间件包装处理器
	handler := c.Handler(mux)
	// http.HandleFunc("/api/upload", uploadFile)
	err := http.ListenAndServe(":6270", handler)
	if err != nil {
		fmt.Errorf("http.ListenAndServe()函数执行错误: %v", err)
	}
}
