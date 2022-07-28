package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
	"userma-lx/protocol"
)

const (
	magicNum = "token-tpy"
)

func createRandomNumber() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func getCurTime() string {
	now := time.Now()
	dateStr := fmt.Sprintf("%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	return dateStr
}

func GetToken(userName string) string {
	var token string
	token += magicNum + "-"
	token += userName + "-"
	token += getCurTime() + "-"
	token += createRandomNumber()
	return MD5(token)
}

func GetFileName(fileName string, ext string) string {
	h := md5.New()
	h.Write([]byte(fileName + strconv.FormatInt(time.Now().Unix(), 10)))
	return hex.EncodeToString(h.Sum(nil)) + ext
}

// CheckAndCreateFileName
func CheckAndCreateFileName(oldName string) (newName string, isLegal bool) {
	ext := path.Ext(oldName)
	if strings.ToLower(ext) == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		newName = GetFileName(oldName, ext)
		isLegal = true
	}
	return newName, isLegal
}

func Encode(data protocol.RPCdata) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		log.Println("fail to encode: ", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (protocol.RPCdata, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data protocol.RPCdata
	if err := decoder.Decode(&data); err != nil {
		log.Println("fail to decode: ", err)
		return protocol.RPCdata{}, err
	}
	return data, nil
}
