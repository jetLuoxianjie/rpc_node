package transport

import (
	"encoding/binary"
	"io"
	"net"
)

//前4个传输数据长度
const headerLen = 4

// Transport 传输结构体
type Transport struct {
	conn net.Conn
}

// NewTransport 创建一个传输
func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

// Send 发生数据
func (t *Transport) Send(data []byte) error {
	//我们将需要4个字节，然后加上数据的len
	buf := make([]byte, headerLen+len(data))
	binary.BigEndian.PutUint32(buf[:headerLen], uint32(len(data)))
	copy(buf[headerLen:], data)
	_, err := t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// Read 读数据
func (t *Transport) Read() ([]byte, error) {
	header := make([]byte, headerLen)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
