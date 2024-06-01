package utils

import (
	"net/http"
	"time"
)

func GenerateUniqueFileName() string {
	currentTime := time.Now()

	return currentTime.Format("20060102150405") + ".jpg"
}

func GenerateUniqueTableName() string {
	currentTime := time.Now()

	// 使用当前时间生成表名
	return "Table_" + currentTime.Format("20060102150405")
}

// AddCoresHeader 添加CORES跨域头
func AddCoresHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
}
