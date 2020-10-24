package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/library/snowFlake"
	"io/ioutil"
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

	return
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

// 目录或文件是否存在
func DirOrFileByIsExists(path string) bool {
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
		log.Error(err)
	}

	return fpt
}

// 生成token
func GenToken() (string, error) {
	sf := snowFlake.NewSnowFlake(1, 1)
	token, err := sf.NextID()

	return strconv.Itoa(int(token)), err
}

// 获取加解密文件文件路径
func GetRsaFilePath() string {
	return GetCurrentDir() + "/script/"
}

// rsa
func GenRsaKey(path string) {
	pubPath := path + "public.pem"
	prvPath := path + "private.key"

	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey := pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey := pem.EncodeToMemory(block)

	if !DirOrFileByIsExists(pubPath) {
		if err = ioutil.WriteFile(pubPath, pubkey, 0666); err != nil {
			log.Panic(err)
		}
	}
	if !DirOrFileByIsExists(prvPath) {
		if err = ioutil.WriteFile(prvPath, prvkey, 0666); err != nil {
			log.Panic(err)
		}
	}
}

// 公钥加密
func RsaEncrypt(data, keyBytes []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// 私钥解密
func RsaDecrypt(ciphertext, keyBytes []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 结构体转map
func StructToJson(data interface{}) (dataMap map[string]interface{}) {
	dataJson, _ := json.Marshal(data)
	_ = json.Unmarshal(dataJson, &dataMap)

	return
}
