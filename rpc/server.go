package rpc

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"time"
	"userma-lx/config"
	"userma-lx/dao"
	"userma-lx/protocol"
	"userma-lx/redis"
	"userma-lx/utils"
)

type TcpServer struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *TcpServer {
	gob.Register(protocol.ReqLogin{})
	gob.Register(protocol.RespProfile{})
	gob.Register(protocol.ReqUploadFile{})
	gob.Register(protocol.ReqVerifyToken{})
	gob.Register(protocol.RespLogin{})
	gob.Register(protocol.ReqSetNickName{})
	gob.Register(protocol.RPCdata{})
	server := TcpServer{addr: addr, funcs: make(map[string]reflect.Value)}
	server.Register("DoLogin", DoLogin)
	server.Register("DoSignUp", DoSignUp)
	server.Register("GetUserInfoByName", GetUserInfoByName)
	server.Register("VerifyToken", VerifyToken)
	server.Register("UpdateNickName", UpdateNickName)
	server.Register("UploadFile", UploadFile)
	return &server
}

func (s *TcpServer) Register(fName string, fFunc interface{}) {
	if _, ok := s.funcs[fName]; ok {
		log.Println("already register func")
		return
	}
	s.funcs[fName] = reflect.ValueOf(fFunc)
}

func (s *TcpServer) Execute(req protocol.RPCdata) protocol.RPCdata {
	//f, ok := s.funcs[req.Name]
	//if !ok {
	//	log.Println(req.Name, " function not register in server")
	//	return protocol.RPCdata{req.Name, nil, "function not register in server"}
	//}
	//log.Printf("func %s is called\n", req.Name)
	//inArgs := make([]reflect.Value, len(req.Args))
	//for i := range req.Args {
	//	inArgs[i] = reflect.ValueOf(req.Args[i])
	//}
	//out := f.Call(inArgs)
	////last one is error
	//respArgs := make([]interface{}, len(out)-1)
	//for i := 0; i < len(out)-1; i++ {
	//	respArgs[i] = out[i].Interface()
	//}
	//var er string
	//if err, ok := out[len(out)-1].Interface().(error); ok {
	//	er = err.Error()
	//}
	//return protocol.RPCdata{req.Name, respArgs, er}
	f, ok := s.funcs[req.Name]
	if !ok {
		log.Println(req.Name, " function not register in server")
		return protocol.RPCdata{req.Name, nil, "function not register in server"}
	}
	log.Printf("func %s is called\n", req.Name)
	inArgs := make([]reflect.Value, 1)
	inArgs[0] = reflect.ValueOf(req.Args)
	out := f.Call(inArgs)
	//last one is error
	respArgs := make([]interface{}, len(out))
	//[0]:data/[1]:error/[2]:error
	for i := 0; i < len(out); i++ {
		respArgs[i] = out[i].Interface()
	}
	var er error = errors.New("")
	log.Println(len(out))

	err := out[len(out)-1].Interface()
	if err != nil {
		log.Println(err)
		er = err.(error)
	}

	//if len(out) == 2 {
	//	if err, ok := out[1].Interface().(error); ok {
	//		er = err
	//	}
	//	if respArgs[0] != nil {
	//		log.Println(respArgs[0])
	//		er = errors.New("1server fun failed")
	//	}
	//} else {
	//	if err, ok := out[2].Interface().(error); ok {
	//		er = err
	//	}
	//	if respArgs[1] != nil {
	//		log.Println(respArgs[1])
	//		er = errors.New("2server fun failed")
	//	}
	//}
	jsonResp, _ := json.Marshal(respArgs[0])
	return protocol.RPCdata{req.Name, jsonResp, er.Error()}

}

func (s *TcpServer) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println("fail to listen, err: ", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("server fail to accept. err: ", err)
			continue
		}
		go func() {
			defer conn.Close()
			connTransport := NewTransport(conn)
			for {
				req, err := connTransport.Read()
				if err != nil {
					if err != io.EOF {
						log.Println("read err: ", err)
						return
					}
				}

				decReq, err := utils.Decode(req)
				if err != nil {
					return
				}
				resp := s.Execute(decReq)
				log.Println("server1")
				encResp, err := utils.Encode(resp)
				log.Println("server1")
				if err != nil {
					return
				}
				err = connTransport.Send(encResp)
				if err != nil {
					return
				}
			}
		}()
	}
}

//login
func DoLogin(req protocol.ReqLogin) (protocol.RespLogin, error) {
	username := req.UserName
	password := req.Password
	password = utils.MD5(password)
	log.Println("client:", username, password)
	ok := redis.VerifyLogin(username, password)
	var user protocol.RespProfile
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
		return protocol.RespLogin{}, errors.New("incorrect username or password")
	}
	log.Println("success dologin in the mysql")
	token := utils.GetToken(username)
	setCache(req, user, token)
	resp := protocol.RespLogin{
		user.UserName,
		user.NickName,
		user.PicName,
		token,
	}
	return resp, nil

}

func checkMysql(username string, password string) (protocol.RespProfile, error) {
	userinfo := dao.GetUserByAccount(username, password)
	log.Println("server:", userinfo.UserName, userinfo.Password)
	if password != userinfo.Password {
		return protocol.RespProfile{}, errors.New("incorrect username or password")
	} else {
		user := protocol.RespProfile{
			UserName: userinfo.UserName,
			NickName: userinfo.NickName,
			PicName:  userinfo.PicName,
		}
		return user, nil
	}
}

func setCache(req protocol.ReqLogin, resp protocol.RespProfile, token string) {
	log.Println("开始进行redis缓存设置：")
	redis.SetUserReq(req.UserName, utils.MD5(req.Password), time.Duration(config.MaxExTimeRedis))
	redis.SetToken(req.UserName, token, time.Duration(config.MaxExTimeRedis))
	redis.SetUserResp(resp.UserName, resp.NickName, resp.PicName)
	log.Println("结束缓存redis缓存设置：")
}

//sign up
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

func GetUserInfoByName(username string) (protocol.RespProfile, error) {
	userinfo, err := GetUserInfoCache(username)
	if err != nil {
		userinfo = dao.GetUserInfoByName(username)
	}
	return userinfo, nil
}

func GetUserInfoCache(username string) (protocol.RespProfile, error) {
	resp, err := redis.GetUserResp(username)
	return resp, err
}

//
func VerifyToken(userinfo protocol.ReqVerifyToken) (bool, error) {
	username := userinfo.UserName
	token := userinfo.Token
	ok := redis.VerifyToken(username, token)
	return ok, nil
}

//
func UpdateNickName(name protocol.ReqSetNickName) (protocol.RespProfile, error) {
	username := name.UserName
	nickname := name.NickName
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

//
func UploadFile(name protocol.ReqUploadFile) (protocol.RespProfile, error) {
	username := name.UserName
	filename := name.FileName
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
