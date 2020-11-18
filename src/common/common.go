package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"goPanel/src/library/snowFlake"
	"net"
	"net/http"
	"os"
	"os/exec"
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
func GenRsaKey(path string, bits int) {
	pubPath := path + "public.pem"
	prvPath := path + "private.pem"

	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create(prvPath)
	if err != nil {
		log.Panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	_ = pem.Encode(privateFile, &privateBlock)

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create(pubPath)
	if err != nil {
		log.Panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	_ = pem.Encode(publicFile, &publicBlock)
}

// 公钥加密
func RsaEncrypt(plainText []byte, path string) ([]byte, error) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//返回密文
	return cipherText, nil
}

// 私钥解密
func RsaDecrypt(cipherText []byte, path string) ([]byte, error) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)

	return plainText, nil
}

// 结构体转map
func StructToJson(data interface{}) (dataMap map[string]interface{}) {
	dataJson, _ := json.Marshal(data)
	_ = json.Unmarshal(dataJson, &dataMap)

	return
}

func portInUse(port int) bool {
	checkStatement := fmt.Sprintf("netstat -anp | grep -q %d ", port)
	output, err := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if err != nil {
		return true // err != nil 说明端口没被占用
	}

	if len(output) > 0 {
		return false
	}

	return true
}

// 返回可用的中继端口
func RetRelayPort(port int) int {
	for i := port; i < 65535; i++ {
		if portInUse(i) {
			return i
		}
	}

	return -1
}

// 创建tcp连接
func ConnTcp(addr string) (*net.TCPConn, error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Error("Client connect error ! " + err.Error())
		return nil, err
	}

	return conn, nil
}

// interface转map
func InterfaceByMapStr(data interface{}) (map[string]interface{}, error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var ret map[string]interface{}
	if err = json.Unmarshal(dataJson, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
