package main

import (
	"protocol"

	"io/ioutil"
	"log"
	"net/rpc"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	_CLIENT     *rpc.Client
	_RPC_MSG    chan *rpc.Call
	_CAN_CANCEL chan bool
)

func init() {
	_RPC_MSG = make(chan *rpc.Call, 1024)
	_CAN_CANCEL = make(chan bool)
}

// GetAllFile 获取目录中所有文件
func GetAllFile(pathname string) ([]string, error) {
	filenames := []string{}

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return filenames, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			log.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			log.Printf("[+]Find file: " + pathname + fi.Name())
			filenames = append(filenames, pathname+fi.Name())
		}
	}
	return filenames, nil
}

// rpcconnect 连接rpc服务端
func rpcconnect() {
	client, err := rpc.DialHTTP("tcp", "localhost:8083")
	if err != nil {
		panic(err.Error())
	}

	_CLIENT = client
}

// rpcReconnect 重新连接rpc服务端
func rpcReconnect() bool {
	client, err := rpc.DialHTTP("tcp", "localhost:8083")
	if err != nil {
		log.Printf("ReDial RPC Server Error: %s", err)
		return false
	}

	_CLIENT = client
	log.Printf("ReDial RPC Server Sucess")

	return true
}

// rpcreport 上报日志信息
func rpcreport(event fsnotify.Event, filehash []byte) {
	// var resp *string //返回值
	// err := client.Call(, event, &resp)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// log.Println(*resp)

	var resp string
	args := protocol.ReportEvent{FileName: event.Name, FileEvent: event.Op.String(), FileHash: filehash}

	if nil == _CLIENT {
		for {
			if rpcReconnect() {
				break
			}

			time.Sleep(5000 * time.Millisecond)
		}
	}

	_CLIENT.Call(protocol.RPCReportEvent, args, &resp)
}

// func loop() {
// 	for {
// 		select {
// 		case rpcMsg, ok := <- _RPC_MSG:
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
