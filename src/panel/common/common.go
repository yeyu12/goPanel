package common

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func RetJson(g *gin.Context, code int32, msg string, data interface{}) {
	g.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	})

	g.Abort()
}

func GetEnvDefaultBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val != "" {
		defaultValue, _ = strconv.ParseBool(val)
	}

	return defaultValue
}

func GetEnvDefaultString(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		defaultValue = val
	}

	return defaultValue
}

func GetEnvDefaultInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val != "" {
		defaultValue, _ = strconv.Atoi(val)
	}

	return defaultValue
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

// 目录是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

// 创建目录
func CreatePath(path string) bool {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Error("%s", err)
		return false
	}

	return true
}

// 获取当前工作目录
func GetCurrentDir() string {
	fpt, err := filepath.Abs("")
	if err != nil {
		log.Println(err)
	}

	return fpt
}
