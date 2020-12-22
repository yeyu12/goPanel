package common

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type StringUtils string

func (s *StringUtils) Set(v string) {
	if v != "" {
		*s = StringUtils(v)
	} else {
		s.Clear()
	}
}

func (s *StringUtils) Clear() {
	*s = StringUtils(0x1E)
}

func (s StringUtils) Exist() bool {
	return string(s) != string(0x1E)
}

func (s StringUtils) Bool() (bool, error) {
	v, err := strconv.ParseBool(s.String())
	return bool(v), err
}

func (s StringUtils) Float32() (float32, error) {
	v, err := strconv.ParseFloat(s.String(), 32)
	return float32(v), err
}

func (s StringUtils) Float64() (float64, error) {
	return strconv.ParseFloat(s.String(), 64)
}

func (s StringUtils) Int() (int, error) {
	v, err := strconv.ParseInt(s.String(), 10, 32)
	return int(v), err
}

func (s StringUtils) Int8() (int8, error) {
	v, err := strconv.ParseInt(s.String(), 10, 8)
	return int8(v), err
}

func (s StringUtils) Int16() (int16, error) {
	v, err := strconv.ParseInt(s.String(), 10, 16)
	return int16(v), err
}

func (s StringUtils) Int32() (int32, error) {
	v, err := strconv.ParseInt(s.String(), 10, 32)
	return int32(v), err
}

func (s StringUtils) Int64() (int64, error) {
	v, err := strconv.ParseInt(s.String(), 10, 64)
	return int64(v), err
}

func (s StringUtils) Uint() (uint, error) {
	v, err := strconv.ParseUint(s.String(), 10, 32)
	return uint(v), err
}

func (s StringUtils) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(s.String(), 10, 8)
	return uint8(v), err
}

func (s StringUtils) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(s.String(), 10, 16)
	return uint16(v), err
}

func (s StringUtils) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(s.String(), 10, 32)
	return uint32(v), err
}

func (s StringUtils) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(s.String(), 10, 64)
	return uint64(v), err
}

func (s StringUtils) ToTitleLower() string {
	str := strings.ToLower(s.String()[:1]) + s.String()[1:]
	return str
}

func (s StringUtils) ToTitleUpper() string {
	str := strings.ToUpper(s.String()[:1]) + s.String()[1:]
	return str
}

func (s StringUtils) RegexpSQLVal() (bool, error) {
	r := "^[0-9a-zA-Z\\s-_\u4E00-\u9FA5'=@.?]+$"
	b, err := regexp.MatchString(r, s.String())
	return bool(b), err
}

func (s StringUtils) ContainsNum() (bool, error) {
	r := "^.*\\d+.*+$"
	b, err := regexp.MatchString(r, s.String())
	return bool(b), err
}

func (s StringUtils) ContainsBool(sep string) bool {
	index := strings.Index(s.String(), sep)
	return index > -1
}

func (s StringUtils) RegexpSQLSgin() (bool, error) {
	r := "^[<>=!?]+$"
	b, err := regexp.MatchString(r, s.String())
	return bool(b), err
}

func (s StringUtils) String() string {
	if s.Exist() {
		return string(s)
	}
	return ""
}

func (s StringUtils) MD5() string {
	m := md5.New()
	m.Write([]byte(s.String()))
	return hex.EncodeToString(m.Sum(nil))
}

func (s StringUtils) SHA1() string {
	sha := sha1.New()
	sha.Write([]byte(s.String()))
	return hex.EncodeToString(sha.Sum(nil))
}

func (s StringUtils) SHA256() string {
	sha := sha256.New()
	sha.Write([]byte(s.String()))
	return hex.EncodeToString(sha.Sum(nil))
}

func (s StringUtils) SHA512() string {
	sha := sha512.New()
	sha.Write([]byte(s.String()))
	return hex.EncodeToString(sha.Sum(nil))
}

func (s StringUtils) HMAC_SHA1(key string) string {
	mc := hmac.New(sha1.New, []byte(key))
	mc.Write([]byte(s.String()))
	return hex.EncodeToString(mc.Sum(nil))
}

func (s StringUtils) HMAC_SHA256(key string) string {
	mc := hmac.New(sha256.New, []byte(key))
	mc.Write([]byte(s.String()))
	return hex.EncodeToString(mc.Sum(nil))
}

func (s StringUtils) HMAC_SHA512(key string) string {
	mc := hmac.New(sha512.New, []byte(key))
	mc.Write([]byte(s.String()))
	return hex.EncodeToString(mc.Sum(nil))
}

func (s StringUtils) Base64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(s.String()))
}

func (s StringUtils) Base64Decode() (string, error) {
	v, err := base64.StdEncoding.DecodeString(s.String())
	return string(v), err
}

// GZIP压缩
func (s StringUtils) GzipEncode() (string, error) {
	var b bytes.Buffer
	gW, err := gzip.NewWriterLevel(&b, gzip.HuffmanOnly)
	defer gW.Close()
	if err != nil {
		return "", err
	}

	_, err = gW.Write([]byte(s.String()))
	if err != nil {
		return "", err
	}
	err = gW.Flush()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// GZIP 解压
func (s StringUtils) GzipDecode() (string, error) {
	gRead, err := gzip.NewReader(strings.NewReader(s.String()))
	if err != nil {
		return "", err
	}
	defer gRead.Close()
	rBuf, err := ioutil.ReadAll(gRead)
	if err != nil {
		return "", err
	}

	return string(rBuf), nil
}
