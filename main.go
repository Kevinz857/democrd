package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

var (
	flagSet             = flag.NewFlagSet("democrd", flag.ExitOnError)
	master              = flag.String("master", "", "The address of the kubernetes API server")
	kubeconfig          = flag.String("kubeconfig", "", "Path to a kubeconfig")
	onlyOneSingleHander = make(chan struct{})
	shutdownSignals     = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func setupSignalHandler() (stoCh <-chan struct{}) {
	close(onlyOneSingleHander)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)

	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	return stop
}

func main() {

	flag.Parse()

	//设置一个信号处理，应用于优雅关闭
	stopCh := setupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	mydemoClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}
	//informerFactory 工厂类，这里注入我们通过代码生成的client
	//client主要用于和API server进行通信，实现ListAndWatch
	mydemoInformerFactory := informers.NewSharedInformerFactory(mydemoClient, time.Second*30)

	//生成一个democrd组的Mydemo对象传递给自定义控制器
	controller := NewController(kubeClient, mydemoClient,
		mydemoInformerFactory().V1().Mydemos())

	go mydemoInformerFactory(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}