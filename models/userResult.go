package models

type UserResult struct {
	Id       int    `db:"id"`
	UserName string `db:"username"`
	NickName string `db:"nickname"`
	PicName  string `db:"picname"`
	Password string `db:"password"`
}
