package rpc

import (
	"log"
	"net"
	"userma-lx/protocol"
	"userma-lx/utils"
)

type TcpClient struct {
	pool chan net.Conn
}

func NewClient(numConn int, addr string) (TcpClient, error) {
	pool := make(chan net.Conn, numConn)
	for i := 0; i < numConn; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println("init client failed")
			return TcpClient{nil}, err
		}
		pool <- conn
	}
	return TcpClient{pool}, nil
}

func (c *TcpClient) Call(rpcName string, req interface{}, resp interface{}) error {
	conn := c.getConn()
	defer c.releaseConn(conn)
	reqData := protocol.RPCdata{
		Name: rpcName,
		Args: req,
		Err:  "",
	}
	cReqTransport := NewTransport(conn)
	cReqBuff, err := utils.Encode(reqData)
	if err != nil {
		return err
	}
	err = cReqTransport.Send(cReqBuff)
	if err != nil {
		return err
	}
	cRespBuff, err := cReqTransport.Read()
	if err != nil {
		return err
	}
	respData, err := utils.Decode(cRespBuff)
	if err != nil {
		return err
	}
	//if respData.Args == nil {
	//	log.Println("no out for the call function")
	//	return errors.New("no respArgs")
	//}
	resp = respData
	return nil
}

func (c *TcpClient) getConn() net.Conn {
	select {
	case conn := <-c.pool:
		return conn
	}
}

func (c *TcpClient) releaseConn(conn net.Conn) {
	select {
	case c.pool <- conn:
		return
	}
}
