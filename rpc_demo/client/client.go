package client

import (
	"errors"
	"net"
	"reflect"
	"rpc_demo/dataserial"
	"rpc_demo/transport"
)

// Client rpc客户端结构体
type Client struct {
	conn net.Conn
}

// NewClient 创建一个新的客户端
func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

// CallRPC 调用rpc
func (c *Client) CallRPC(rpcName string, fPtr interface{}) {
	container := reflect.ValueOf(fPtr).Elem()
	f := func(req []reflect.Value) []reflect.Value {
		//创建一个传输的
		cReqTransport := transport.NewTransport(c.conn)
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		// 填充输入参数
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}

		// 请求rpc接口
		reqRPC := dataserial.RPCdata{Name: rpcName, Args: inArgs}
		b, err := dataserial.Encode(reqRPC)
		if err != nil {
			panic(err)
		}

		//发送请求
		err = cReqTransport.Send(b)
		if err != nil {
			return errorHandler(err)
		}
		// 等待 接受 远程服务的响应
		rsp, err := cReqTransport.Read()
		if err != nil {
			return errorHandler(err)
		}
		rspDecode, _ := dataserial.Decode(rsp)
		if rspDecode.Err != "" {
			return errorHandler(errors.New(rspDecode.Err))
		}

		if len(rspDecode.Args) == 0 {
			rspDecode.Args = make([]interface{}, container.Type().NumOut())
		}
		// 解包
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 {
				if rspDecode.Args[i] == nil { // 符合为nil 就设置为 Zero
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rspDecode.Args[i])
				}
			} else { // 最后一个是错误处理
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}

		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), f))
}
