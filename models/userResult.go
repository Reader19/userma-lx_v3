package models

type UserResult struct {
	Id       int    `db:"uid"`
	UserName string `db:"UserName"`
	NickName string `db:"NickName"`
	PicName  string `db:"PicName"`
	Password string `db:"Password"`
}
