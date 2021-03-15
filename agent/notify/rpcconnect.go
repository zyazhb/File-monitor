package notify

import (
	"log"
	"net/rpc"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	// RPCReportEvent 上报信息的方法
	RPCReportEvent = "MonitorServer.ReportEvent"
)

//ReportEvent 传递的数据结构 一定要与server/rpcserver统一
type ReportEvent struct {
	AgentIP   string
	FileName  string
	FileEvent string
	FileHash  string
}

var (
	// _Client    *rpc.Client
	_RpcMsg    chan *rpc.Call
	_CanCancel chan bool
)

func init() {
	_RpcMsg = make(chan *rpc.Call, 1024)
	_CanCancel = make(chan bool)
}

//rpcreport inotify上报日志信息
func rpcreporti(event fsnotify.Event, filehash string, serverip string) {
	client, _ := rpc.DialHTTP("tcp", serverip)
	var resp string
	//上报内容
	args := ReportEvent{AgentIP: geAgentIP(), FileName: event.Name, FileEvent: event.Op.String(), FileHash: filehash}

	if nil == client {
		for {
			if rpcReconnect(serverip,client) {
				break
			}

			time.Sleep(5000 * time.Millisecond)
		}
	}
	client.Call(RPCReportEvent, args, &resp)
}

//rpcreport fanotify上报日志信息
func rpcreportfan(FileName, FileEvent, filehash string, serverip string) {
	client, _ := rpc.DialHTTP("tcp", serverip)
	var resp string
	//上报内容
	args := ReportEvent{AgentIP: geAgentIP(), FileName: FileName, FileEvent: FileEvent, FileHash: filehash}

	if nil == client {
		for {
			if rpcReconnect(serverip, client) {
				break
			}

			time.Sleep(5000 * time.Millisecond)
		}
	}
	client.Call(RPCReportEvent, args, &resp)
}

// rpcReconnect 重新连接rpc服务端
func rpcReconnect(serverip string, NewClient *rpc.Client) bool {
	client, err := rpc.DialHTTP("tcp", serverip)
	if err != nil {
		log.Printf("ReDial RPC Server Error: %s", err)
		return false
	}

	NewClient = client
	log.Printf("ReDial RPC Server Sucess")

	return true
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
