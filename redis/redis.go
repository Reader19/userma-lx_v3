package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"usermaLX4/config"
	"usermaLX4/protocol"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	initClient()
}

func initClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: config.RedisPoolSize,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("redis does not work")
		return
	}
}

func SetUserReq(username string, password string, maxtime time.Duration) {
	err := rdb.Set(ctx, username+"_pwd", password, maxtime*1e9).Err()
	if err != nil {
		log.Println("fail to set UserCache: ", err)
	} else {
		log.Println("success to set UserCache")
	}

}
func SetToken(username string, token string, maxtime time.Duration) {
	err := rdb.Set(ctx, username+"_tok", token, maxtime*1e9).Err()
	if err != nil {
		log.Println("fail to set token: ", err)
	} else {
		log.Println("success to set token")
	}
}

func SetUserResp(username string, nickname string, picname string) {
	userinfo := map[string]interface{}{
		"NickName": nickname,
		"PicName":  picname,
	}
	err := rdb.HMSet(ctx, username+"_inf", userinfo).Err()
	if err != nil {
		log.Println("fail to set userinfo: ", err)
	} else {
		log.Println("success to set userinfo")
	}
}

func SetNickname(username string, nickname string) (error, bool, protocol.RespProfile) {
	userinfo, err := rdb.HGetAll(ctx, username+"_inf").Result()
	if err != nil {
		log.Println("userInfo not in redis")
		return err, false, protocol.RespProfile{}
	}
	resp := protocol.RespProfile{
		username,
		userinfo["NickName"],
		userinfo["PicName"],
	}
	userinfo["NickName"] = nickname
	err = rdb.HMSet(ctx, username+"_inf", userinfo).Err()
	if err != nil {
		log.Println("update fail")
		return err, true, resp
	}
	resp.NickName = nickname
	return err, true, resp
}

func SetPicName(username string, picname string) (error, bool, protocol.RespProfile) {
	userinfo, err := rdb.HGetAll(ctx, username+"_inf").Result()
	if err != nil {
		log.Println("userInfo not in redis")
		return err, false, protocol.RespProfile{}
	}
	resp := protocol.RespProfile{
		username,
		userinfo["NickName"],
		userinfo["PicName"],
	}
	userinfo["PicName"] = picname
	err = rdb.HMSet(ctx, username+"_inf", userinfo).Err()
	if err != nil {
		log.Println("update fail")
		return err, true, resp
	}
	resp.PicName = picname
	return err, true, resp
}

func GetUserResp(username string) (protocol.RespProfile, error) {
	userinfo, err := rdb.HGetAll(ctx, username+"_inf").Result()
	if err != nil {
		log.Println("fail to get userinfo in server")
		return protocol.RespProfile{}, err
	}
	resp := protocol.RespProfile{
		username,
		userinfo["NickName"],
		userinfo["PicName"],
	}
	return resp, nil
}

func VerifyToken(username string, token string) bool {
	servertoken, err := rdb.Get(ctx, username+"_tok").Result()
	if err != nil {
		log.Println("fail to verify token in server")
		return false
	}
	log.Println("success to verify token in server")
	return token == servertoken
}

func VerifyLogin(username string, password string) bool {
	serverpassword, err := rdb.Get(ctx, username+"_pwd").Result()
	if err != nil {
		log.Println("account not in redis")
		return false
	}
	return password == serverpassword
}
