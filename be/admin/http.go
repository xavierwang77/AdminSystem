package admin

import (
	DB "admin_system/db"
	"admin_system/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// UploadAvatarFile 处理上传的头像s
func UploadAvatarFile(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	//解析文件上传请求
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "无法解析表单", http.StatusBadRequest)
		return
	}

	//获取文件句柄和头信息
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "无法获取文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//在后端指定目录中创建文件
	avatarFilename = utils.GenerateUniqueFileName()
	f, err := os.OpenFile("images/"+avatarFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "无法创建文件", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	//将文件内容复制到后端创建的文件中
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "无法复制文件内容", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "文件上传成功")
}

// ChangeAvatarFile 处理更改用户头像
func ChangeAvatarFile(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	// 解析文件上传请求
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "无法解析表单", http.StatusBadRequest)
		return
	}

	// 获取文件句柄和头信息
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "无法获取文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 在后端指定目录中创建文件
	avatarFilename = utils.GenerateUniqueFileName()
	f, err := os.OpenFile("images/"+avatarFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "无法创建文件", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 将文件内容复制到后端创建的文件中
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "无法复制文件内容", http.StatusInternalServerError)
		return
	}
}

func GetAvatarFilename(w http.ResponseWriter, r *http.Request) {
	filename := avatarFilename

	//返回响应
	response := map[string]string{"filename": filename}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UploadData 处理提交的注册表数据
func UploadData(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var registerData registerData
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 检查用户名是否已经存在于管理员信息表中
	exists, err := DB.UsernameExistsInAdminTable(db, registerData.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "用户名已存在", http.StatusConflict)
		return
	}

	tableName := "admin" //该表用户存储管理员信息

	//自动建表（用户）
	adminTableName, err := DB.CreateTableIfNotExists(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	registerData.adminTableName = adminTableName

	//若用户未上传头像，则使用默认头像
	if len(registerData.AvatarFilename) == 0 {
		registerData.AvatarFilename = "v2-6afa72220d29f045c15217aa6b275808_hd.jpg"
	}

	//将数据插入到生成的表格中
	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (name, email, password, phone_number, avatar_filename, table_name, remarks) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		tableName), registerData.Name, registerData.Email, registerData.Password, registerData.PhoneNumber, registerData.AvatarFilename, registerData.adminTableName, registerData.Remarks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//返回响应
	response := map[string]string{"message": "Register successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VerifyLoginData 处理管理员登录数据
func VerifyLoginData(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginData loginData
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询数据库中是否存在匹配的用户名和密码
	tableName := "admin" // 你的用户表
	var storedPassword string
	err = db.QueryRow(fmt.Sprintf("SELECT password FROM %s WHERE name = $1", tableName), loginData.Name).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// 检查密码是否匹配
	if loginData.Password != storedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// 登录成功，返回成功响应
	response := map[string]string{"message": "Login successful"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LocateAvatar 定位管理员头像路径
func LocateAvatar(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName adminName
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询数据库中是否存在匹配的用户名
	tableName := "admin" // 你的用户表
	var avatarFilename string
	err = db.QueryRow(fmt.Sprintf("SELECT avatar_filename FROM %s WHERE name = $1", tableName), adminName.Name).Scan(&avatarFilename)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 构建头像文件的完整路径
	avatarPath := "images/" + avatarFilename // 假设文件在 images/ 目录下

	// 打开头像文件
	file, err := os.Open(avatarPath)
	if err != nil {
		http.Error(w, "Error opening avatar file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// // 构建头像文件的完整URL
	// avatarURL := "http://localhost:8080/images/" + avatarFilename // 修改为适当的URL

	// // 返回头像文件URL作为响应
	// response := map[string]string{"avatarURL": avatarURL}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)

	// 返回头像文件数据作为响应
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", avatarFilename))
	http.ServeFile(w, r, avatarPath)
}

// GetOriginalAdminData 处理前端获取管理员原始数据
func GetOriginalAdminData(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName adminName
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询指定用户名的用户数据
	tableName := "admin" // 你的用户表
	var user showUser
	err = db.QueryRow(fmt.Sprintf("SELECT name, email, password, phone_number, remarks FROM %s WHERE name = $1", tableName), adminName.Name).
		Scan(&user.Name, &user.Email, &user.Password, &user.PhoneNumber, &user.Remarks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回用户数据作为 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ChangeAdminData 处理更改管理员数据
func ChangeAdminData(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeAdminData changeAdminData
	if err := json.NewDecoder(r.Body).Decode(&changeAdminData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 使用 UPDATE 语句更新数据
	tableName := "admin"
	_, err = db.Exec(fmt.Sprintf("UPDATE %s SET name=$1, email=$2, password=$3, phone_number=$4, avatar_filename=$5, remarks=$6 WHERE name=$7",
		tableName), changeAdminData.Name, changeAdminData.Email, changeAdminData.Password, changeAdminData.PhoneNumber, changeAdminData.AvatarFilename, changeAdminData.Remarks, changeAdminData.OriginalAdminName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//返回响应
	response := map[string]string{"message": "Change successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetOriginalAdminAvatarName 处理获取管理员头像文件名
func GetOriginalAdminAvatarName(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName adminName
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//将用户原始头像文件名提取出来
	tableName := "admin"
	var originalAvatarName string
	err = db.QueryRow(fmt.Sprintf("SELECT avatar_filename FROM %s WHERE name = $1", tableName), adminName.Name).Scan(&originalAvatarName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	//返回响应
	response := map[string]string{"originalAvatarName": originalAvatarName}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteAdmin 处理注销管理员
func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	utils.AddCoresHeader(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName adminName
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := DB.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "admin"
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), adminName.Name).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 删除用户表
	_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", userTableName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//在管理员数据表中删除管理员信息
	_, err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE name = $1", tableName), adminName.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//返回响应
	response := map[string]string{"message": "Delete Admin successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
