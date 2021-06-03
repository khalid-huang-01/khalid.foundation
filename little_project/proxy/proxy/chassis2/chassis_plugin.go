package chassis2

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/prometheus/common/log"
	_ "khalid.foundation/proxy/proxy/chassis2/panel"
	registry2 "khalid.foundation/proxy/proxy/chassis2/registry"
)

// 整体的go-chassis使用流程：
// 1. 首先安装一个控制面,panel
// 2. 注册一个服务发现中心
// 3. 安装负载均衡算法
// 4. 整个控制面初始化
// 5. 接收到请求后，创建一条责任链，组装负载均衡、传输处理、回调函数

// Install installs go-chassis plugins
func init() {
	// service discovery
	opt := registry.Options{}
	registry.DefaultServiceDiscoveryService = registry2.NewEdgeServiceDiscovery(opt)
	// load balance
	loadbalancer.InstallStrategy(loadbalancer.StrategyRandom, func() loadbalancer.Strategy {
		return &loadbalancer.RandomStrategy{}
	})

	loadbalancer.InstallStrategy(loadbalancer.StrategyRoundRobin, func() loadbalancer.Strategy {
		return &loadbalancer.RoundRobinStrategy{}
	})

	// control panel
	config.GlobalDefinition = &model.GlobalCfg{
		Panel: model.ControlPanel{
			Infra: "edge",
		},
		Ssl: make(map[string]string),
	}
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	if err := control.Init(opts); err != nil {
		log.Errorf("failed to init control: %v", err)
	}
	// init archaius
	if err := archaius.Init(); err != nil {
		log.Errorf("failed to init archaius: %v", err)
	}
}