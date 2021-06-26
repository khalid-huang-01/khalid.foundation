package activestandby

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/deprecated/scheme"
	 "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"
)

func Run(readyzAdaptor *ReadyzAdaptor) {
	// 获取clientgo
	var clientset *kubernetes.Clientset
	k8sConfig := "D:\\workspace\\gocode\\gomodule\\local-conf\\config"
	config, err := clientcmd.BuildConfigFromFlags("", k8sConfig)
	if err != nil {
		klog.Infof("error %v", err)
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		klog.Infof("error %v", err)
	}

	coreBroadcaster := record.NewBroadcaster()
	coreRecorder := coreBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "ha-server"})

	leaderElection := &componentbaseconfig.LeaderElectionConfiguration{
		LeaderElect:       false,
		LeaseDuration:     metav1.Duration{Duration: 15 * time.Second},
		RenewDeadline:     metav1.Duration{Duration: 10 * time.Second},
		RetryPeriod:       metav1.Duration{Duration: 2 * time.Second},
		ResourceLock:      "endpointsleases",
		ResourceNamespace: "ha-test",
		ResourceName:      "cloudcorelease",
	}
	id, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	id = id + "_" + string(uuid.NewUUID())
	rl, err := resourcelock.New(
		leaderElection.ResourceLock,
		leaderElection.ResourceNamespace,
		leaderElection.ResourceName,
		clientset.CoreV1(),
		clientset.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: coreRecorder,
		},
	)
	if err != nil {
		panic(err)
	}
	leaderElectionConfig := &leaderelection.LeaderElectionConfig{
		Lock:            rl,
		ReleaseOnCancel: true, // 用这个实现服务快速释放锁
		LeaseDuration:   leaderElection.LeaseDuration.Duration,
		RenewDeadline:   leaderElection.RenewDeadline.Duration,
		RetryPeriod:     leaderElection.RetryPeriod.Duration,
		Callbacks:       leaderelection.LeaderCallbacks{},
		WatchDog:        nil,
		Name:            "cloudcore",
	}
	
	leaderElectionConfig.Callbacks = leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			StartMainHTTPServer()
			// TODO set podreadinessgate
		},
		OnStoppedLeading: func() {
			klog.Errorf("leaderElection lost ")
			// TODO set PodReadinessGate
			//triggerGracefulShutdown()
		},
		OnNewLeader:      nil,
	}

	leaderElector, err := leaderelection.NewLeaderElector(*leaderElectionConfig)
	if err != nil {
		klog.Errorf("error: %v", err)
		return
	}
	readyzAdaptor.SetLeaderElection(leaderElector)
	ctx, cancel := context.WithCancel(context.Background())
	go leaderElector.Run(ctx)

	// 监听信号
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	select {
	case s := <-c:
		klog.Infof("Get os signal %v", s.String())
		//Cleanup each modules
		cancel()
	}

	for {
		if leaderElector.IsLeader() {
			time.Sleep(1)
		} else {
			break
		}
	}

}

func triggerGracefulShutdown()  {
	klog.Errorln("Trigger graceful shutdown!")
	p, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		klog.Errorf("Failed to find self process: %v", err)
	}
	err = p.Signal(os.Interrupt)
	if err != nil {
		klog.Errorf("Failed to trigger graceful shutdown: %v", err)
	}
}