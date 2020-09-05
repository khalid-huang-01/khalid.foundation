package async

import (
	"bryson.foundation/kbuildresource/models"
	"sync"
)

type RequestHandler interface {
	// 在请求进入执行队列之前，做一些前置处理，诸如参数校验，默认值填充等工作，values里面保存了在处理链上需要传递的一些中间信息，避免重复计算
	// 这一步，无论是同步还是异步都需要做，是一些比较轻量级的操作
	PreExec(requestDTO interface{}, requestType string, values map[string]interface{}) error
	// 根据请求参数封装出request，并将request进行缓存； 记得把dto里面的instanceName注入到request中
	MakeRequest(requestDTO interface{}, requestType string, values map[string]interface{}) (*models.Request, error)
	// 异步执行请求，limiChan用于控制 并发量，wg用于通知外部请求执行完成
	AsyncExec(request *models.Request, limitChan <-chan struct{}, wg *sync.WaitGroup)
	// 在请求进入执行队列之后需要做的一些操作
	PostAsyncExec(request *models.Request, requestType string, values map[string]interface{}) error
	// 同步执行请求
	SyncExec(requestDTO interface{}, requestType string, values map[string]interface{}) (interface{}, error)
	// 给requestDTO 设置请求的instanceName
	SetInstanceName(requestDTO interface{}, instanceName string)
}