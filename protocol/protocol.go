package protocol

type ReqLogin struct {
	UserName string
	Password string
}

type RespProfile struct {
	UserName string
	NickName string
	PicName  string
}

type RPCdata struct {
	Name string
	Args interface{}
	Err  string
}

type RespLogin struct {
	UserName string
	NickName string
	PicName  string
	Token    string
}

type ReqSetNickName struct {
	UserName string
	NickName string
}

type ReqUploadFile struct {
	UserName string
	FileName string
}
type ReqVerifyToken struct {
	UserName string
	Token    string
}
