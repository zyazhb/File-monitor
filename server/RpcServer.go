package main

import (
	"main/model"

	"io/ioutil"
	"log"

	"github.com/google/logger"

	"net"
	"net/http"
	"net/rpc"
)

const (
	// RPCReportEvent 上报信息的方法
	RPCReportEvent = "MonitorServer.ReportEvent"
)

//ReportEvent 传递的数据结构 一定要与agent/notify/rpcconnect统一
type ReportEvent struct {
	FileName  string
	FileEvent string
	FileHash  string
}

//MonitorServer uncomment
type MonitorServer struct {
}

//ReportEvent 该方法向外暴露ReportEvent
func (ms *MonitorServer) ReportEvent(event *ReportEvent, resp *string) error {
	if event.FileHash != "" {
		logger.Infof("\033[1;33m [*]%s hash:%x\n", event.FileName, event.FileHash)
	}
	logger.Infof("\033[1;33m [*]%s file:%s\033[0m", event.FileEvent, event.FileName)
	model.RPCDbInsert(event.FileName, event.FileEvent, event.FileHash)
	*resp = event.FileName
	return nil //返回类型
}

// RPCServer Rpc服务端
func RPCServer() {
	//1、初始化指针数据类型
	MonitorServer := new(MonitorServer) //初始化指针数据类型
	model.RPCDbInit()                   //初始化log数据库

	//2、调用net/rpc包的功能将服务对象进行注册
	err := rpc.Register(MonitorServer)
	if err != nil {
		panic(err.Error())
	}

	//3、通过该函数把mathUtil中提供的服务注册到HTTP协议上，方便调用者可以利用http的方式进行数据传递
	rpc.HandleHTTP()
	logger.Init("RpcLogger", true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	//4、在特定的端口进行监听
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err.Error())
	}
	logger.Info("Server up")
	http.Serve(listen, nil)
	logger.Warning("Server down")
}
