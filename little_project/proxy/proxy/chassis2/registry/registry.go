package registry

import (
	"fmt"
	"strings"

	"github.com/go-chassis/go-chassis/core/registry"
	utiltags "github.com/go-chassis/go-chassis/pkg/util/tags"
)

const (
	// EdgeRegistry constants string
	EdgeRegistry = "edge"
)

type instanceList []*registry.MicroServiceInstance

func (I instanceList) Len() int {
	return len(I)
}

func (I instanceList) Less(i, j int) bool {
	return strings.Compare(I[i].InstanceID, I[j].InstanceID) < 0
}

func (I instanceList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

// init initialize the trafficplugin of edge meta registry
func init() { registry.InstallServiceDiscovery(EdgeRegistry, NewEdgeServiceDiscovery) }

// EdgeServiceDiscovery to represent the object of service center to call the APIs of service center
type EdgeServiceDiscovery struct {
	Name string
}

func NewEdgeServiceDiscovery(options registry.Options) registry.ServiceDiscovery {
	return &EdgeServiceDiscovery{
		Name: EdgeRegistry,
	}
}

// GetAllMicroServices Get all MicroService information.
func (esd *EdgeServiceDiscovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	return nil, nil
}

// 这个会在load_balancer.go的BuildStragegy里面调用，用于返回所有的实例，然后给调度算法使用
// FindMicroServiceInstances find micro-service instances (subnets)
func (esd *EdgeServiceDiscovery) FindMicroServiceInstances(consumerID, microServiceName string, tags utiltags.Tags) ([]*registry.MicroServiceInstance, error) {
	fmt.Println("microServiceName: ", microServiceName)
	// 不管传入的micorServicename是什么，都返回127.0.0.1:30002
	hostIP := "127.0.0.1"
	HostPort := 30002

	// 这个地址的EndpointsMap会在loadbalance_handler.go里面的getEndpoint函数里面使用
	// 具体的作用是会把这个key值作为invocation对象，也就是调用链的Protocol变量的值
	// 这里EndPointsMap里面的value值最后会作为loadbalance_handler.go里面的getEndpoint的返回值
	var microServiceInstances instanceList
	microServiceInstances = append(microServiceInstances, &registry.MicroServiceInstance{
		InstanceID:      "",
		ServiceID:       "",
		EndpointsMap: map[string]string{"tcp": fmt.Sprintf("%s:%d", hostIP, HostPort)},
	})

	return microServiceInstances, nil
}

// GetMicroServiceID get microServiceID
func (esd *EdgeServiceDiscovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	return "", nil
}

// GetMicroServiceInstances return instances
func (esd *EdgeServiceDiscovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	return nil, nil
}

// GetMicroService return service
func (esd *EdgeServiceDiscovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	return nil, nil
}

// AutoSync updating the cache manager
func (esd *EdgeServiceDiscovery) AutoSync() {}

// Close close all websocket connection
func (esd *EdgeServiceDiscovery) Close() error { return nil }