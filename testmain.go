package testmain

//import (
//	"log"
//	"net/http"
//	"time"
//	"userma-lx/config"
//	"userma-lx/router"
//	"userma-lx/rpc"
//)
//
//func main() {
//	tcpServer := rpc.NewServer(config.TcpServerAddr)
//	go tcpServer.Run()
//	time.Sleep(1 * time.Second)
//
//	server := http.Server{
//		Addr: "localhost:8080",
//	}
//	router.Router()
//	err := server.ListenAndServe()
//	if err != nil {
//		log.Println(err)
//	}
//
//}
