package service

import (
	"errors"
	"log"
	"time"
	"userma-lx/config"
	"userma-lx/dao"
	"userma-lx/protocol"
	"userma-lx/redis"
	"userma-lx/utils"
)

func GetUserInfonByname(username string) protocol.ResqLogin {
	userinfo, err := GetUserInfoCache(username)
	if err != nil {
		userinfo = dao.GetUserInfoByName(username)
	}
	return userinfo
}

func DoLogin(req protocol.ReqLogin) (protocol.ResqLogin, string, error) {
	username := req.UserName
	password := req.Password
	password = utils.MD5(password)
	log.Println("client:", username, password)
	ok := redis.VerifyLogin(username, password)
	var user protocol.ResqLogin
	var err error
	if ok {
		user, err = redis.GetUserResp(username)
		if err != nil {
			user, err = checkMysql(username, password)
			log.Println("fail dologin in the redis")
		} else {
			log.Println("success dologin in the redis")
		}
	} else {
		user, err = checkMysql(username, password)
	}
	if err != nil {
		return protocol.ResqLogin{}, "", errors.New("incorrect username or password")
	}
	log.Println("success dologin in the mysql")
	token := utils.GetToken(username)
	setCache(req, user, token)
	return user, token, nil

}

func checkMysql(username string, password string) (protocol.ResqLogin, error) {
	userinfo := dao.GetUserByAccount(username, password)
	log.Println("server:", userinfo.UserName, userinfo.Password)
	if password != userinfo.Password {
		return protocol.ResqLogin{}, errors.New("incorrect username or password")
	} else {
		user := protocol.ResqLogin{
			UserName: userinfo.UserName,
			NickName: userinfo.NickName,
			PicName:  userinfo.PicName,
		}
		return user, nil
	}
}

func setCache(req protocol.ReqLogin, resp protocol.ResqLogin, token string) {
	log.Println("开始进行redis缓存设置：")
	redis.SetUserReq(req.UserName, utils.MD5(req.Password), time.Duration(config.MaxExTimeRedis))
	redis.SetToken(req.UserName, token, time.Duration(config.MaxExTimeRedis))
	redis.SetUserResp(resp.UserName, resp.NickName, resp.PicName)
	log.Println("结束缓存redis缓存设置：")
}

func DoSignUp(req protocol.ReqLogin) error {
	username := req.UserName
	password := req.Password
	password = utils.MD5(password)
	ok := dao.GetUserByNameBool(username)
	if !ok {
		log.Println("the account exists")
		return errors.New("the account exists")
	}
	log.Println(username, password)
	err := dao.InsertUser(username, password)
	return err
}

func VerifyToken(username string, token string) bool {
	ok := redis.VerifyToken(username, token)
	return ok
}

func GetUserInfoCache(username string) (protocol.ResqLogin, error) {
	resp, err := redis.GetUserResp(username)
	return resp, err
}

func UpdateNickname(username string, nickname string) (protocol.ResqLogin, error) {
	err := dao.UpdateNickName(username, nickname)
	if err != nil {
		userinfo := dao.GetUserInfoByName(username)
		return userinfo, err
	}
	err, ok, resp := redis.SetNickname(username, nickname)
	if !ok {
		userinfo := dao.GetUserInfoByName(username)
		return userinfo, err
	}
	return resp, err
}

func UploadFile(username string, filename string) (protocol.ResqLogin, error) {
	err := dao.UpdatePicName(username, filename)
	if err != nil {
		userinfo := dao.GetUserInfoByName(username)
		return userinfo, err
	}
	err, ok, resp := redis.SetPicName(username, filename)
	if !ok {
		userinfo := dao.GetUserInfoByName(username)
		return userinfo, err
	}
	return resp, err
}
