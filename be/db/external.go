package db

import (
	"admin_system/utils"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDB 连接到PG数据库
func ConnectDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=admin_system_user dbname=admin_system_db password=766515 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateTableIfNotExists 自动建表函数（用户）
func CreateTableIfNotExists(db *sql.DB) (string, error) {
	tableName := utils.GenerateUniqueTableName() //使用当前时间戳作为表名

	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS admin_system.%s (
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

// CreateAdminTableIfNotExists 自动建表函数（管理员）
func CreateAdminTableIfNotExists(db *sql.DB) (string, error) {
	tableName := utils.GenerateUniqueTableName() //使用当前时间戳作为表名

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

// UsernameExistsInAdminTable 检查用户名是否已经存在于管理员信息表中
func UsernameExistsInAdminTable(db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM admin WHERE name = $1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
