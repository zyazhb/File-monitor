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
	_Client    *rpc.Client
	_RpcMsg    chan *rpc.Call
	_CanCancel chan bool
)

func init() {
	_RpcMsg = make(chan *rpc.Call, 1024)
	_CanCancel = make(chan bool)
}

//rpcreport inotify上报日志信息
func rpcreporti(event fsnotify.Event, filehash string, serverip string) {
	_Client, _ = rpc.Dial("udp", serverip)
	var resp string
	//上报内容
	args := ReportEvent{AgentIP: geAgentIP(), FileName: event.Name, FileEvent: event.Op.String(), FileHash: filehash}

	if nil == _Client {
		for {
			if rpcReconnect(serverip, _Client) {
				break
			}

			time.Sleep(5000 * time.Millisecond)
		}
	}
	_Client.Call(RPCReportEvent, args, &resp)
	_Client = nil
}

//rpcreport fanotify上报日志信息
func rpcreportfan(FileName, FileEvent, filehash string, serverip string) {
	_Client, _ = rpc.Dial("udp", serverip)
	var resp string
	//上报内容
	args := ReportEvent{AgentIP: geAgentIP(), FileName: FileName, FileEvent: FileEvent, FileHash: filehash}

	if nil == _Client {
		for {
			if rpcReconnect(serverip, _Client) {
				break
			}

			time.Sleep(5000 * time.Millisecond)
		}
	}
	_Client.Call(RPCReportEvent, args, &resp)
	_Client = nil
}

// rpcReconnect 重新连接rpc服务端
func rpcReconnect(serverip string, _Client *rpc.Client) bool {
	_Client, err := rpc.Dial("udp", serverip)
	if err != nil {
		log.Printf("ReDial RPC Server Error: %s", err)
		return false
	}

	log.Printf("ReDial RPC Server Sucess")
	_Client = nil
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
