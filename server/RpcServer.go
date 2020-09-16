package main

import (
	"net"
	"net/http"
	"net/rpc"
	"log"
)

//MonitorServer uncomment
type MonitorServer struct{
}
//ReportEvent 该方法向外暴露ReportEvent
func (ms *MonitorServer) ReportEvent(message string, resp *string) error {
	log.Println(message)
	*resp = message
	return nil //返回类型
}


func RpcServer(){
	//1、初始化指针数据类型
	MonitorServer := new(MonitorServer) //初始化指针数据类型

	//2、调用net/rpc包的功能将服务对象进行注册
	err := rpc.Register(MonitorServer)
	if err != nil {
		panic(err.Error())
	}

	//3、通过该函数把mathUtil中提供的服务注册到HTTP协议上，方便调用者可以利用http的方式进行数据传递
	rpc.HandleHTTP()

	//4、在特定的端口进行监听
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err.Error())
	}
	log.Println("Server up")
	http.Serve(listen, nil)
	log.Println("Server down")
}




