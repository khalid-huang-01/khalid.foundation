package loadbalancer

import (
	"math/rand"
	"sync"

	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
)

// RoundRobinStrategy is strategy
type RoundRobinStrategy struct {
	instances []*registry.MicroServiceInstance
	key       string
}

func newRoundRobinStrategy() loadbalancer.Strategy {
	return &RoundRobinStrategy{}
}

// 这个用于接受一些额外的信息用于做负载均衡，在RoundRobinStragegy里面，会在
// load_balancer.go的BuildStragegy里面接收实例数量和serviceKey
//ReceiveData receive data
func (r *RoundRobinStrategy) ReceiveData(inv *invocation.Invocation, instances []*registry.MicroServiceInstance, serviceKey string) {
	r.instances = instances
	r.key = serviceKey
}

// 在loadbalancer_hadnler的getEndpoint里会调用具体策略的Pick函数，用于挑选
//Pick return instance
func (r *RoundRobinStrategy) Pick() (*registry.MicroServiceInstance, error) {
	if len(r.instances) == 0 {
		return nil, loadbalancer.ErrNoneAvailableInstance
	}

	i := pick(r.key)
	return r.instances[i%len(r.instances)], nil
}

var rrIdxMap = make(map[string]int)
var mu sync.RWMutex

func pick(key string) int {
	mu.RLock()
	i, ok := rrIdxMap[key]
	if !ok {
		mu.RUnlock()
		mu.Lock()
		i, ok = rrIdxMap[key]
		if !ok {
			i = rand.Int()
			rrIdxMap[key] = i
		}
		rrIdxMap[key]++
		mu.Unlock()
		return i
	}

	mu.RUnlock()
	mu.Lock()
	rrIdxMap[key]++
	mu.Unlock()
	return i
}
