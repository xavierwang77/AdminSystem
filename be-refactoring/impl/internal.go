package impl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func generateUniqueFileName() string {
	currentTime := time.Now()

	return currentTime.Format("20060102150405") + ".jpg"
}

func generateUniqueTableName() string {
	currentTime := time.Now()

	// 使用当前时间生成表名
	return "Table_" + currentTime.Format("20060102150405")
}

// 连接到PG数据库
func connectDB() (*sql.DB, error) {
	connStr := "host=192.168.85.128 port=5432 user=postgres dbname=postgres password=766515 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 自动建表函数（用户）
func createTableIfNotExists(db *sql.DB) (string, error) {
	tableName := generateUniqueTableName() //使用当前时间戳作为表名

	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			email VARCHAR(255),
			password VARCHAR(255),
			phone_number VARCHAR(255),
			avatar_filename VARCHAR(255),
			remarks VARCHAR(255),
			login_privilege VARCHAR(255)
		);
	`, tableName)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return "", err
	}

	return tableName, nil
}

// 自动建表函数（管理员）
func createAdminTableIfNotExists(db *sql.DB) (string, error) {
	tableName := generateUniqueTableName() //使用当前时间戳作为表名

	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			email VARCHAR(255),
			password VARCHAR(255),
			phone_number VARCHAR(255),
			avatar_filename VARCHAR(255),
			table_name VARCHAR(255),
			remarks VARCHAR(255)
		);
	`, tableName)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return "", err
	}

	return tableName, nil
}

// 处理上传的头像
func uploadAvatarFile(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	avatarFilename = generateUniqueFileName()
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

// 处理更改用户头像
func changeAvatarFile(w http.ResponseWriter, r *http.Request) {
	// 添加跨域头，允许所有域名的请求，可以根据需要进行配置
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	avatarFilename = generateUniqueFileName()
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

func getAvatarFilename(w http.ResponseWriter, r *http.Request) {
	filename := avatarFilename

	//返回响应
	response := map[string]string{"filename": filename}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理提交的注册表数据
func uploadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var registerData REGISTERDATA
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 检查用户名是否已经存在于管理员信息表中
	exists, err := usernameExistsInAdminTable(db, registerData.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "用户名已存在", http.StatusConflict)
		return
	}

	tableName := "table_20231006215428" //该表用户存储管理员信息

	//自动建表（用户）
	adminTableName, err := createTableIfNotExists(db)
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

// 检查用户名是否已经存在于管理员信息表中
func usernameExistsInAdminTable(db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM table_20231006215428 WHERE name = $1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// 处理管理员登录数据
func verifyLoginData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginData LOGINDATA
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询数据库中是否存在匹配的用户名和密码
	tableName := "table_20231006215428" // 你的用户表
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

// 处理用户登陆数据
func verifyUserLoginData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginData USERLOGINDATA
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), loginData.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询数据库中是否存在匹配的用户名
	var storedPassword string
	var storedLoginPrivilege string
	err = db.QueryRow(fmt.Sprintf("SELECT password, login_privilege FROM %s WHERE name = $1", userTableName), loginData.Name).Scan(&storedPassword, &storedLoginPrivilege)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// 检查loginPrivilege是否为"false"
	if storedLoginPrivilege == "false" {
		http.Error(w, "Login privilege is false", http.StatusUnauthorized)
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

// 定位管理员头像路径
func locateAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName ADMINNAME
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询数据库中是否存在匹配的用户名
	tableName := "table_20231006215428" // 你的用户表
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

// 定位用户头像路径
func locateUserAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo CHANGEINFO
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), changeInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询用户数据表中是否存在匹配的用户名
	var avatarFilename string
	err = db.QueryRow(fmt.Sprintf("SELECT avatar_filename FROM %s WHERE name = $1", userTableName), changeInfo.UserName).Scan(&avatarFilename)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 构建头像文件的完整URL
	avatarPath := "images/" + avatarFilename

	// 打开头像文件
	file, err := os.Open(avatarPath)
	if err != nil {
		http.Error(w, "Error opening avatar file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 返回头像文件数据作为响应
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", avatarFilename))
	http.ServeFile(w, r, avatarPath)
}

// func locateHomeUserAvatar(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
//         http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
//         return
//     }

// 	//解析前端传入的JSON数据
// 	var userAvatarFileName USERAVATARFILENAME
//     if err := json.NewDecoder(r.Body).Decode(&userAvatarFileName); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

// 	// 构建头像文件的完整URL
// 	avatarPath := "images/" + userAvatarFileName.Name

// 	// 打开头像文件
//     file, err := os.Open(avatarPath)
//     if err != nil {
//         http.Error(w, "Error opening avatar file", http.StatusInternalServerError)
//         return
//     }
//     defer file.Close()

// 	// 返回头像文件数据作为响应
//     w.Header().Set("Content-Type", "application/octet-stream")
//     w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", avatarFilename))
//     http.ServeFile(w, r, avatarPath)
// }

// 处理主页获取用户数据
func fetchUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName ADMINNAME
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), adminName.Name).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询并按 id 升序排序用户数据
	query := fmt.Sprintf("SELECT name, email, password, phone_number, avatar_filename, remarks, login_privilege FROM %s ORDER BY id", userTableName)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []USER
	for rows.Next() {
		var user USER
		err := rows.Scan(&user.Name, &user.Email, &user.Password, &user.PhoneNumber, &user.AvatarFilename, &user.Remarks, &user.LoginPrivilege)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// 返回用户数据作为 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 处理添加用户信息
func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var addForm ADDFORM
	if err := json.NewDecoder(r.Body).Decode(&addForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	adminTableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", adminTableName), addForm.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	//若未上传用户头像，则使用默认头像
	if len(addForm.AvatarFilename) == 0 {
		addForm.AvatarFilename = "v2-6afa72220d29f045c15217aa6b275808_hd.jpg"
	}

	//将数据插入到生成的表格中
	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (name, email, password, phone_number, avatar_filename, remarks, login_privilege) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		userTableName), addForm.Name, addForm.Email, addForm.Password, addForm.PhoneNumber, addForm.AvatarFilename, addForm.Remarks, "true")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//返回响应
	response := map[string]string{"message": "Add successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理删除用户数据
func deleteUserByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var deleteInfo DELETEINFO
	if err := json.NewDecoder(r.Body).Decode(&deleteInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), deleteInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 在用户数据表中删除具有指定名称的行
	_, err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE name = $1", userTableName), deleteInfo.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//返回响应
	response := map[string]string{"message": "Delete successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理前端获取管理员原始数据
func getOriginalAdminData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName ADMINNAME
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询指定用户名的用户数据
	tableName := "table_20231006215428" // 你的用户表
	var user SHOWUSER
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

// 处理前端获取用户原始数据
func getOriginalUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo CHANGEINFO
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), changeInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询指定用户名的用户数据
	var user SHOWUSER
	err = db.QueryRow(fmt.Sprintf("SELECT name, email, password, phone_number, remarks FROM %s WHERE name = $1", userTableName), changeInfo.UserName).
		Scan(&user.Name, &user.Email, &user.Password, &user.PhoneNumber, &user.Remarks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回用户数据作为 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 处理更改管理员数据
func changeAdminData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeAdminData CHANGEADMINDATA
	if err := json.NewDecoder(r.Body).Decode(&changeAdminData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 使用 UPDATE 语句更新数据
	tableName := "table_20231006215428"
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

// 处理更改用户数据
func changeData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeData CHANGEDATA
	if err := json.NewDecoder(r.Body).Decode(&changeData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428"
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), changeData.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询表并按id值排序
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY id", userTableName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 从查询结果中获取数据，并存储在切片中
	var data []CHANGEDATA
	for rows.Next() {
		var item CHANGEDATA
		err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.Password, &item.PhoneNumber, &item.AvatarFilename, &item.Remarks, &item.LoginPrivilege)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, item)
	}

	// 使用 UPDATE 语句更新数据
	_, err = db.Exec(fmt.Sprintf("UPDATE %s SET name=$1, email=$2, password=$3, phone_number=$4, avatar_filename=$5, remarks=$6, login_privilege=$7 WHERE name=$8",
		userTableName), changeData.Name, changeData.Email, changeData.Password, changeData.PhoneNumber, changeData.AvatarFilename, changeData.Remarks, changeData.LoginPrivilege, changeData.OriginalUserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 手动恢复数据的排序顺序
	for i, item := range data {
		_, err = db.Exec(fmt.Sprintf("UPDATE %s SET id=$1 WHERE name=$2", userTableName), i+1, item.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//返回响应
	response := map[string]string{"message": "Change successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理获取管理员头像文件名
func getOriginalAdminAvatarName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName ADMINNAME
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//将用户原始头像文件名提取出来
	tableName := "table_20231006215428"
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

// 处理获取用户头像文件名
func getOriginalAvatarName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo CHANGEINFO
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428"
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), changeInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	//将用户原始头像文件名提取出来
	var originalAvatarName string
	err = db.QueryRow(fmt.Sprintf("SELECT avatar_filename FROM %s WHERE name = $1", userTableName), changeInfo.UserName).Scan(&originalAvatarName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	//返回响应
	response := map[string]string{"originalAvatarName": originalAvatarName}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理获取登陆权限
func getLoginPrivilege(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo CHANGEINFO
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428"
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), changeInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	//将用户原始头像文件名提取出来
	var loginPrivilege string
	err = db.QueryRow(fmt.Sprintf("SELECT login_privilege FROM %s WHERE name = $1", userTableName), changeInfo.UserName).Scan(&loginPrivilege)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	//返回响应
	response := map[string]string{"loginPrivilege": loginPrivilege}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 处理注销管理员
func deleteAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var adminName ADMINNAME
	if err := json.NewDecoder(r.Body).Decode(&adminName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428"
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

// 处理更改用户登录权限
func changeLoginPrivilege(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var deleteInfo DELETEINFO
	if err := json.NewDecoder(r.Body).Decode(&deleteInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 连接到数据库
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "table_20231006215428" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), deleteInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询当前用户的loginPrivilege状态
	var currentPrivilege string
	err = db.QueryRow(fmt.Sprintf("SELECT login_privilege FROM %s WHERE name = $1", userTableName), deleteInfo.UserName).Scan(&currentPrivilege)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 根据当前状态切换loginPrivilege
	var newPrivilege string
	if currentPrivilege == "true" {
		newPrivilege = "false"
	} else if currentPrivilege == "false" {
		newPrivilege = "true"
	}

	// 更新用户数据表中的loginPrivilege
	_, err = db.Exec(fmt.Sprintf("UPDATE %s SET login_privilege = $1 WHERE name = $2", userTableName), newPrivilege, deleteInfo.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// // 重新排序表中的行
	// _, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ORDER BY id", userTableName))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// 返回响应
	response := map[string]string{"message": "Login privilege changed successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
