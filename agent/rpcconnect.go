package main

import (
	"protocol"

	"log"
	"net/rpc"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	_Client    *rpc.Client
	_RpcMsg    chan *rpc.Call
	_CanCancel chan bool
)

func init() {
	_RpcMsg = make(chan *rpc.Call, 1024)
	_CanCancel = make(chan bool)
}

// rpcconnect 连接rpc服务端
func rpcconnect(serverip string) {
	client, err := rpc.DialHTTP("tcp", serverip)
	if err != nil {
		panic(err.Error())
	}

	_Client = client
}

// rpcReconnect 重新连接rpc服务端
func rpcReconnect(serverip string) bool {
	client, err := rpc.DialHTTP("tcp", serverip)
	if err != nil {
		log.Printf("ReDial RPC Server Error: %s", err)
		return false
	}

	_Client = client
	log.Printf("ReDial RPC Server Sucess")

	return true
}

// rpcreport 上报日志信息
func rpcreport(event fsnotify.Event, filehash string, serverip string) {
	// var resp *string //返回值
	// err := client.Call(, event, &resp)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// log.Println(*resp)

	var resp string
	args := protocol.ReportEvent{FileName: event.Name, FileEvent: event.Op.String(), FileHash: filehash}

	if nil == _Client {
		for {
			if rpcReconnect(serverip) {
				break
			}

			time.Sleep(5000 * time.Millisecond)
		}
	}

	_Client.Call(protocol.RPCReportEvent, args, &resp)
}

// func loop() {
// 	for {
// 		select {
// 		case rpcMsg, ok := <- _RpcMsg:
// 			if !ok {
// 				log.Println("Rpc Call error!")
// 			}
// 			rpcMsgHandler(rpcMsg)
// 		}
// 	}
// }

// func rpcMsgHandler(msg *rpc.Call) {
// 	switch msg.ServiceMethod {
// 	case protocol.RPCReportEvent:
// 		reply := msg.Reply.(*string)
// 		log.Printf("Reply message: %s\n", *reply)
// 	default:
// 		log.Fatalln("Can't receiver any reply!")
// 	}
// }
