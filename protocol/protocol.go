package protocol

type ReqLogin struct {
	UserName string
	Password string
}

type ResqLogin struct {
	UserName string
	NickName string
	PicName  string
}
