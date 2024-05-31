package impl

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"net/http"
)

var avatarFilename string

// StartHttpServer 开启http服务
func StartHttpServer() error {
	mux := http.NewServeMux()
	// 添加CORS中间件
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // 允许前端应用的域名
	})

	// 注册处理函数
	mux.HandleFunc("/uploadAvatarFile", uploadAvatarFile)
	mux.HandleFunc("/changeAvatarFile", changeAvatarFile)
	mux.HandleFunc("/uploadData", uploadData)
	mux.HandleFunc("/verifyLoginData", verifyLoginData)
	mux.HandleFunc("/locateAvatar", locateAvatar)
	mux.HandleFunc("/locateUserAvatar", locateUserAvatar)
	mux.HandleFunc("/fetchUserData", fetchUserData)
	mux.HandleFunc("/addUser", addUser)
	mux.HandleFunc("/deleteUserByName", deleteUserByName)
	mux.HandleFunc("/getOriginalUserData", getOriginalUserData)
	mux.HandleFunc("/getOriginalAdminData", getOriginalAdminData)
	mux.HandleFunc("/changeData", changeData)
	mux.HandleFunc("/changeAdminData", changeAdminData)
	mux.HandleFunc("/getOriginalAvatarName", getOriginalAvatarName)
	mux.HandleFunc("/getOriginalAdminAvatarName", getOriginalAdminAvatarName)
	mux.HandleFunc("/getAvatarFilename", getAvatarFilename)
	mux.HandleFunc("/deleteAdmin", deleteAdmin)
	mux.HandleFunc("/changeLoginPrivilege", changeLoginPrivilege)
	mux.HandleFunc("/verifyUserLoginData", verifyUserLoginData)
	mux.HandleFunc("/getLoginPrivilege", getLoginPrivilege)

	// 使用CORS中间件包装处理器
	handler := c.Handler(mux)
	// http.HandleFunc("/upload", uploadFile)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return fmt.Errorf("http.ListenAndServe()函数执行错误: %v", err)
	}

	return nil
}
