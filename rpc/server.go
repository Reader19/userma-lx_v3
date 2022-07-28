package rpc

import (
	"io"
	"log"
	"net"
	"reflect"
	"userma-lx/protocol"
	"userma-lx/utils"
)

type TcpServer struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *TcpServer {
	server := TcpServer{addr: addr, funcs: make(map[string]reflect.Value)}
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
	respArgs := out[0].Interface()
	log.Println(len(out))
	var er string
	if err, ok := out[1].Interface().(error); ok {
		er = err.Error()
	}
	return protocol.RPCdata{req.Name, respArgs, er}
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
		defer conn.Close()
		go func() {
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
				encResp, err := utils.Encode(resp)
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
