package async

import (
	"bryson.foundation/kbuildresource/async/handler"
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/models"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type RequestController struct {
	limitChan chan struct{} //用于控制并发，可以用协程池来做
	stopCh chan struct{} // 控制器停止通道
	requestChannel chan *models.Request
}

var r *RequestController

func init() {
	r = &RequestController{
		limitChan:      make(chan struct{}, 2000),
		stopCh:         make(chan struct{}),
		requestChannel: make(chan *models.Request, 2000),
	}
	go r.StartUp()
}

func GetRequestController() *RequestController {
	return r
}

// 启动请求控制器
func (r *RequestController) StartUp() {
	wg := sync.WaitGroup{} // 用于知道所有request的执行情况
	for request := range r.requestChannel {
		r.limitChan <- struct{}{}
		wg.Add(1)
		logrus.Info("INFO: receive request %s and start handle", request.Name)
		requestHandler := getHandlerFromRequestType(request.RequestType)
		go requestHandler.AsyncExec(request, r.limitChan, &wg)
	}
	wg.Wait()
	logrus.Info("INFO: finish all request")
	close(r.stopCh) // 通知shutdown函数继续执行
}

func (r *RequestController) Shutdown() {
	logrus.Info("INFO: shutdown requestController")
	// sleep 一小段时间，保证收到的请求都入channel了
	time.Sleep(2 * time.Second)
	close(r.requestChannel) // 关闭requestChannel,促使requestHandler里面的for range循环可以在遍历完成之后结束
	<-r.stopCh // 等待requestHandle处理完成的信号，当close(r.stopCh)时可以结束
}

// 请求管理器对外提供的接收请求的接口
// @Param requestDTO interface{} 请求传输对象，主要包含用户传入的参数
// @Param requestType string 请求类型，用于分派请求到对应的处理器
// return interface{} 请求处理的返回结果
func (r *RequestController) AcceptRequest(requestDTO interface{}, requestType string) (interface{}, error) {
	requestHandler := getHandlerFromRequestType(requestType)
	values := make(map[string]interface{}, 0)
	err := requestHandler.PreExec(requestDTO, requestType, values)
	if err != nil {
		return requestDTO, err
	}
	request, err := requestHandler.MakeRequest(requestDTO, requestType, values)
	if err != nil {
		logrus.Error("ERROR: MakeRequest failed, try to use SyncExec")
		return requestHandler.SyncExec(requestDTO, requestType, values)
	}
	go r.sendRequestToChannel(request)
	err = requestHandler.PostAsyncExec(request, requestType, values)
	if err != nil {
		logrus.Error("ERROR: posyAsyncExec failed, err: ", err)
	}
	return requestDTO, nil
}

func (r *RequestController) sendRequestToChannel(request *models.Request) {
	r.requestChannel <- request
}

//func (r *RequestController) TakeOverRequest(deadInstanceName string, newInstanceName string) error {
//	log.Infof("INFO: start takeover request of %s", deadInstanceName)
//	request, err :=
//}

func getHandlerFromRequestType(requestType string) handler.RequestHandler {
	s := strings.Split(requestType, "_")
	switch s[0] {
	case common.BuildJobPrefix:
		return handler.GetBuildJobHandler()
	default:
		return nil
	}
}