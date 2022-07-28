package rpc

//
//import (
//	"encoding/gob"
//	"errors"
//	"log"
//	"time"
//	"userma-lx/config"
//	"userma-lx/dao"
//	"userma-lx/protocol"
//	"userma-lx/rpc"
//	"userma-lx/utils"
//)
//
//func main() {
//	gob.Register(protocol.ReqLogin{})
//	server := rpc.NewServer(config.TcpServerAddr)
//	server.Register("TestDoSignUp", TestDoSignUp)
//	go server.Run()
//	time.Sleep(1 * time.Second)
//
//	client, err := rpc.NewClient(50, config.TcpServerAddr)
//	if err != nil {
//		log.Println(err)
//	}
//	costumer := protocol.ReqLogin{
//		UserName: "testRPC",
//		Password: "12345",
//	}
//	var newerr error
//	err = client.Call("TestDoSignUp", costumer, &newerr)
//	log.Println("call 返回结果 ", err)
//	log.Println("输出结果：", newerr)
//}
//
//func TestDoSignUp(req protocol.ReqLogin) error {
//	username := req.UserName
//	password := req.Password
//	password = utils.MD5(password)
//	ok := dao.GetUserByNameBool(username)
//	if !ok {
//		log.Println("the account exists")
//		return errors.New("the account exists")
//	}
//	log.Println(username, password)
//	err := dao.InsertUser(username, password)
//	return err
//}
