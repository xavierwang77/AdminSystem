package user

import (
	DB "admin_system/db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// VerifyUserLoginData 处理用户登陆数据
func VerifyUserLoginData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginData userLoginData
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

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	tableName := "admin" // 你的用户表
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

// LocateUserAvatar 定位用户头像路径
func LocateUserAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo changeInfo
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
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
	tableName := "admin" // 你的用户表
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

// FetchUserData 处理主页获取用户数据
func FetchUserData(w http.ResponseWriter, r *http.Request) {
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
	tableName := "admin" // 你的用户表
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

	var users []user
	for rows.Next() {
		var user user
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

// AddUser 处理添加用户信息
func AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var addForm addForm
	if err := json.NewDecoder(r.Body).Decode(&addForm); err != nil {
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

	// 查询管理员数据表中是否存在匹配的用户名并把该管理员名下的用户数据表名提取出来
	adminTableName := "admin" // 你的用户表
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

// DeleteUserByName 处理删除用户数据
func DeleteUserByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var deleteInfo deleteInfo
	if err := json.NewDecoder(r.Body).Decode(&deleteInfo); err != nil {
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
	tableName := "admin" // 你的用户表
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

// GetOriginalUserData 处理前端获取用户原始数据
func GetOriginalUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo changeInfo
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
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
	tableName := "admin" // 你的用户表
	var userTableName string
	err = db.QueryRow(fmt.Sprintf("SELECT table_name FROM %s WHERE name = $1", tableName), changeInfo.AdminName).Scan(&userTableName)
	if err != nil {
		http.Error(w, "Administrator not found", http.StatusNotFound)
		return
	}

	// 查询指定用户名的用户数据
	var user showUser
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

// ChangeUserData 处理更改用户数据
func ChangeUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeData userChangeData
	if err := json.NewDecoder(r.Body).Decode(&changeData); err != nil {
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
	var data []userChangeData
	for rows.Next() {
		var item userChangeData
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

// GetOriginalAvatarName 处理获取用户头像文件名
func GetOriginalAvatarName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo changeInfo
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
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

// GetLoginPrivilege 处理获取登陆权限
func GetLoginPrivilege(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var changeInfo changeInfo
	if err := json.NewDecoder(r.Body).Decode(&changeInfo); err != nil {
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

// ChangeLoginPrivilege 处理更改用户登录权限
func ChangeLoginPrivilege(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//解析前端传入的JSON数据
	var deleteInfo deleteInfo
	if err := json.NewDecoder(r.Body).Decode(&deleteInfo); err != nil {
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
	tableName := "admin" // 你的用户表
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

	// 返回响应
	response := map[string]string{"message": "Login privilege changed successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
