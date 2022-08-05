package dao

import (
	"log"
	"usermaLX4/models"
	"usermaLX4/protocol"
)

func GetUserInfoByName(username string) protocol.RespProfile {
	row := DB.QueryRow("select * from users where username=?", username)
	if row.Err() != nil {
		log.Println(row.Err())
	}
	var user models.UserResult
	//row.StructScan(&user)
	row.Scan(&user.Id, &user.UserName, &user.NickName, &user.PicName, &user.Password)
	userinfo := protocol.RespProfile{
		user.UserName,
		user.NickName,
		user.PicName,
	}
	return userinfo

}

func GetUserByAccount(username string, password string) models.UserResult {
	row := DB.QueryRow("select * from users where username=?", username)
	if row.Err() != nil {
		log.Println("数据库筛选错误", row.Err())
	}
	var user models.UserResult
	//row.StructScan(&user)
	row.Scan(&user.Id, &user.UserName, &user.NickName, &user.PicName, &user.Password)
	log.Println(user)
	return user
}

func GetUserByNameBool(username string) bool {
	log.Println(username)
	var name string
	err := DB.QueryRow("select username from users where username=?", username).Scan(&name)
	if err != nil {
		log.Println(err)
		return true
	}
	log.Println("err is false")
	return false
}

func InsertUser(username string, password string) error {
	//Pwd, _ := os.Getwd()
	defaultPicName := "lx.png"
	log.Println(defaultPicName)
	_, err := DB.Exec("INSERT INTO users(username, nickname, picname, password) VALUES (?,?,?,?)", username, username, defaultPicName, password)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateNickName(username string, nickname string) error {
	_, err := DB.Exec("update users set nickname=? where username=?", nickname, username)
	if err != nil {
		log.Println("fail update: ", err)
	}
	return err
}

func UpdatePicName(username string, picname string) error {
	_, err := DB.Exec("update users set picname=? where username=?", picname, username)
	if err != nil {
		log.Println("fail update picname: ", err)
	}
	return err
}
