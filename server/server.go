package main

import (
	"net"
	"net/rpc"
	"log"
	"net/http"
	"time"
	"path/filepath"
	"os"

	"github.com/gorilla/mux"
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

func RpcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    log.Println(w, "Category: %v\n", vars["category"])
}


type spaHandler struct {
	staticPath string
	indexPath  string
}
// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
        // if we failed to get the absolute path respond with a 400 bad request
        // and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

    // check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
        // if we got an error (that wasn't that the file doesn't exist) stating the
        // file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/rpc/{key}", RpcHandler)
	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
    router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}


	DbInit()
	go RpcServer()
	log.Fatal(srv.ListenAndServe())
}