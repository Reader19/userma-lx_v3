package dao

import (
	"log"
	"userma-lx/models"
	"userma-lx/protocol"
)

//func GetUser() models.UserResult {
//	rows, _ := DB.Query("select * from users limit 1")
//	log.Println(rows)
//	var user models.UserResult
//	for rows.Next() {
//		err := rows.Scan(&user.Id, &user.UserName, &user.NickName, &user.PicName, &user.Password)
//		if err != nil {
//			log.Println("bad")
//		}
//	}
//
//	log.Println(user)
//	return user
//}

func GetUserInfoByName(username string) protocol.RespProfile {
	row := DB.QueryRowx("select * from users where UserName=?", username)
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
	row := DB.QueryRowx("select * from users where UserName=?", username)
	if row.Err() != nil {
		log.Println(row.Err())
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
	err := DB.QueryRow("select UserName from users where UserName=?", username).Scan(&name)
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
	_, err := DB.Exec("INSERT INTO users(UserName, NickName, PicName, Password) VALUES (?,?,?,?)", username, username, defaultPicName, password)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateNickName(username string, nickname string) error {
	_, err := DB.Exec("update users set NickName=? where UserName=?", nickname, username)
	if err != nil {
		log.Println("fail update: ", err)
	}
	return err
}

func UpdatePicName(username string, picname string) error {
	_, err := DB.Exec("update users set PicName=? where UserName=?", picname, username)
	if err != nil {
		log.Println("fail update picname: ", err)
	}
	return err
}
