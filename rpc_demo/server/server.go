package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"rpc_demo/dataserial"
	"rpc_demo/transport"
)

// RPCServer rpc服务
type RPCServer struct {
	addr  string
	funcs map[string]reflect.Value //使用反射 value ,这个类型可以调用函数
}

// NewServer 创建一个rpc服务
func NewServer(addr string) *RPCServer {
	return &RPCServer{addr: addr, funcs: make(map[string]reflect.Value)}
}

// Register 注册rpc服务
func (s *RPCServer) Register(fnName string, fFunc interface{}) {
	if _, ok := s.funcs[fnName]; ok {
		return
	}
	s.funcs[fnName] = reflect.ValueOf(fFunc)
}

// Execute 如何函数存在就执行函数
func (s *RPCServer) Execute(req dataserial.RPCdata) dataserial.RPCdata {
	//根据函数名字调用具体的方法
	f, ok := s.funcs[req.Name]
	if !ok {
		//没有该方法
		e := fmt.Sprintf("func %s not Registered", req.Name)
		log.Println(e)
		return dataserial.RPCdata{Name: req.Name, Args: nil, Err: e}
	}

	log.Printf("func %s is called\n", req.Name)
	//填充参数
	inArgs := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		inArgs[i] = reflect.ValueOf(req.Args[i])
	}

	// 调用方法
	out := f.Call(inArgs)
	//打包响应数据
	resArgs := make([]interface{}, len(out)-1)
	for i := 0; i < len(out)-1; i++ {
		resArgs[i] = out[i].Interface()
	}
	// 看打包的数据有没有报错的
	var er string
	if _, ok := out[len(out)-1].Interface().(error); ok {
		er = out[len(out)-1].Interface().(error).Error()
	}
	return dataserial.RPCdata{Name: req.Name, Args: resArgs, Err: er}
}

// rpc服务运行
func (s *RPCServer) Run() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("listen on %s err: %v\n", s.addr, err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}
		go func() {
			connTransport := transport.NewTransport(conn)
			for {
				// 读数据
				req, err := connTransport.Read()
				if err != nil {
					if err != io.EOF {
						log.Printf("read err: %v\n", err)
						return
					}
				}

				//解码数据
				decReq, err := dataserial.Decode(req)
				if err != nil {
					log.Printf("Error Decoding the Payload err: %v\n", err)
					return
				}
				//运行
				resP := s.Execute(decReq)
				// 编码运行后的数据
				b, err := dataserial.Encode(resP)
				if err != nil {
					log.Printf("Error Encoding the Payload for response err: %v\n", err)
					return
				}
				//发送响应返回
				err = connTransport.Send(b)
				if err != nil {
					log.Printf("transport write err: %v\n", err)
				}
			}
		}()
	}
}
