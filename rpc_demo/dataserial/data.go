package dataserial

import (
	"bytes"
	"encoding/gob"
)

// RPCdata
type RPCdata struct {
	Name string        // 调用方法名称
	Args []interface{} // 请求 或 响应 的数据
	Err  string        // 执行的远程服务的报错信息
}

//编码
func Encode(data RPCdata) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 解码
func Decode(b []byte) (RPCdata, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data RPCdata
	if err := decoder.Decode(&data); err != nil {
		return RPCdata{}, err
	}
	return data, nil
}