package rpc

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Transport struct {
	conn net.Conn
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

func (t *Transport) Send(data []byte) error {
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := t.conn.Write(buf)
	if err != nil {
		log.Println("conn.write is wrong, err: ", err)
		return err
	}
	return nil
}

func (t *Transport) Read() ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		log.Println("conn.read is wrong, err: ", err)
		return nil, err
	}
	datalen := binary.BigEndian.Uint32(header)
	data := make([]byte, datalen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		log.Println("conn.read is wrong, err: ", err)
		return nil, err
	}
	return data, nil
}
