package handler

import (
	"bryson.foundation/kbuildresource/buildjob"
	"bryson.foundation/kbuildresource/cache"
	"bryson.foundation/kbuildresource/common"
	"bryson.foundation/kbuildresource/dto"
	"bryson.foundation/kbuildresource/models"
	"bryson.foundation/kbuildresource/utils"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type BuildJobHandler struct {

}

var handler *BuildJobHandler

func init() {
	handler = &BuildJobHandler{}
}

func GetBuildJobHandler() *BuildJobHandler {
	return handler
}

func (b *BuildJobHandler) PreExec(requestDTO interface{}, requestType string, values map[string]interface{}) error {
	buildJobDTO, ok := requestDTO.(*dto.BuildJobDTO)
	if !ok {
		logrus.Error("ERROR: BuildJobHandler PreExec requestDTO is not a type of dto.BuildJobDTO")
		return fmt.Errorf("buildJobHandler PreExec requestDTO is not a type of dto.BuildJobDTO")
	}
	switch requestType {
	case common.BuildJobCreateRequestType:
		// 重新命名
		if buildJobDTO.ReName {
			buildJobDTO.Name = buildJobDTO.Name + "-" + utils.CreateRandomString(5)
		}
		return buildjob.VerifyBuildJobDTO(buildJobDTO)
	default:
		return fmt.Errorf("invalid reqeusType %s",requestType)
	}
}

func (b *BuildJobHandler) MakeRequest(requestDTO interface{}, requestType string, values map[string]interface{}) (*models.Request, error) {
	buildJobDTO, ok := requestDTO.(*dto.BuildJobDTO)
	if !ok {
		logrus.Error("ERROR: BuildJobHandler PreExec requestDTO is not a type of dto.BuildJobDTO")
		return nil, fmt.Errorf("buildJobHandler PreExec requestDTO is not a type of dto.BuildJobDTO")
	}
	buildJobDTOJsonData, err := json.Marshal(buildJobDTO)
	if err != nil {
		return nil, fmt.Errorf("marshal requestDTO failed")
	}
	request := &models.Request{
		Name:        buildJobDTO.Name,
		Status:      common.RequestStatusPending,
		RequestType: requestType,
		RequestDTO:   string(buildJobDTOJsonData),
	}
	// 放入cache中并返回
	err = cache.AddRequest(request)
	if err != nil {
		return nil, err
	}
	logrus.Info("INFO: create request success and add to cache")
	return request, nil
}

func (b *BuildJobHandler) AsyncExec(request *models.Request, limitChan <-chan struct{}, wg *sync.WaitGroup) {
	time.Sleep(10 * time.Second)
	defer func() {
		<-limitChan
		wg.Done()
	}()
	switch request.RequestType {
	case common.BuildJobCreateRequestType:
		err := transferRequestStatus(request, common.RequestStatusExecuting)
		if err != nil {
			logrus.Error("ERROR: AsyncExec failed, err: ", err)
			return
		}
		buildJobDTO := &dto.BuildJobDTO{}
		err = json.Unmarshal([]byte(request.RequestDTO), buildJobDTO)
		if err != nil {
			logrus.Error("ERROR: AsyncExec failed, err: ", err)
			err = transferRequestStatus(request, common.RequestStatusFailed)
			logrus.Error("ERROR: AsyncExec failed, err: ", err)
			return
		}
		err = buildjob.CreatePod(buildJobDTO)
		if err != nil {
			logrus.Error("ERROR: AsyncExec failed, err: ", err)
			err = transferRequestStatus(request, common.RequestStatusFailed)
			logrus.Error("ERROR: AsyncExec failed, err: ", err)
			return
		}
		err = transferRequestStatus(request, common.RequestStatusSuccess)
		if err != nil {
			logrus.Error("ERROR: AsyncExec failed, err: ", err)
			return
		}
		logrus.Info("INFO: finish AsyncExec")
	default:
		return
	}
}

func (b *BuildJobHandler) PostAsyncExec(request *models.Request, requestType string, values map[string]interface{}) error {
	return nil
}

func (b *BuildJobHandler) SyncExec(requestDTO interface{}, requestType string, values map[string]interface{}) (interface{}, error) {
	buildJobDTO , ok := requestDTO.(*dto.BuildJobDTO)
	if !ok {
		logrus.Error("ERROR: BuildJobHandler PreExec requestDTO is not a type of dto.BuildJobDTO")
		return nil, fmt.Errorf("buildJobHandler PreExec requestDTO is not a type of dto.BuildJobDTO")
	}
	err := buildjob.CreatePod(buildJobDTO)
	if err != nil {
		return nil, err
	}
	return buildJobDTO, nil
}

func transferRequestStatus(request *models.Request, status string) error {
	if request.Status == status {
		return nil
	}
	if status == common.RequestStatusExecuting {
		logrus.Info("INFO: transfer request status to executing")
		request.Status = common.RequestStatusExecuting
		return cache.UpdateRequest(request)
	}
	if status == common.RequestStatusFailed {
		// 更新状态，删除cache,并转到dao层
		logrus.Info("INFO: transfer request status to failed")
		request.Status = common.RequestStatusFailed
		err := cache.DeleteRequest(request)
		if err != nil {
			return err
		}
		// 转到dao层，
		_, err = models.AddRequest(request)
		return err
	}
	if status == common.RequestStatusSuccess {
		//直接从cache里面删除
		logrus.Info("INFO: finish request and delete from redis")
		return cache.DeleteRequest(request)
	}
	return nil
}

