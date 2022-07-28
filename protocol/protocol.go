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

type RPCdata struct {
	Name string
	Args interface{}
	Err  string
}
