package protocol

const (
	// RPCReportEvent 上报信息的方法
	RPCReportEvent = "MonitorServer.ReportEvent"
)

// ReportEvent 要上传的信息
type ReportEvent struct {
	FileName  string
	FileEvent string
	FileHash  string
}
