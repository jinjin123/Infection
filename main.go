package main

import (
	"infection/etcd"
	//"infection/machineinfo"
	"infection/util/lib"
	"log"
	"os"
	"strconv"

	//"strings"
	"sync/atomic"
	"systray"
	//"time"
	"infection/transfer"
	"infection/util/icon"
)

var localAddr string

type Info struct {
	Dev        bool
	ClientPort int
}

var Config = Info{
	true,
	8888,
}

type AppConfig struct {
	Url string
}

// confirm lock
type AppConfigMgr struct {
	config atomic.Value
}

var appConfigMgr = &AppConfigMgr{}

func (a *AppConfigMgr) Callback(conf *etcd.Config) {
	appConfig := &AppConfig{}
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(appConfig)
}
func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("freedom")
	mQuit := systray.AddMenuItem("Quit", "Quit freedom")
	start := systray.AddMenuItem("Start", "Start")
	stop := systray.AddMenuItem("Stop", "Stop")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	//loop up the switch signal
	for {
		select {
		case <-start.ClickedCh:
			start.Check()
			stop.Uncheck()
			start.SetTitle("Start")
			systray.SetTooltip("running")
		case <-stop.ClickedCh:
			stop.Check()
			start.Uncheck()
			stop.SetTitle("Stop")
			systray.SetTooltip("stop")
		}
	}
}
func onExit() {
	// clean up here
}
func init() {
	//todo request daemon check

	if !Config.Dev {
		log.Println("已启动free客户端，请在free_" + strconv.Itoa(Config.ClientPort) + ".log查看详细日志")
		f, _ := os.OpenFile("free"+strconv.Itoa(Config.ClientPort)+".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
		log.SetOutput(f)
	}

	localAddr = ":" + strconv.Itoa(Config.ClientPort)
}
func main() {
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)
	//machineinfo.MachineSend(conf.Url)
	go lib.DoUpdate(conf.Url)
	go transfer.InitCfg(conf.Url, localAddr)
	systray.Run(onReady, onExit)
}

//func run(){
//	for {
//		appConfig := appConfigMgr.config.Load().(*AppConfig)
//		fmt.Println("Hostname:", appConfig.Url)
//		fmt.Printf("%v\n", "--------------------")
//		time.Sleep(5 * time.Second)
//	}
//}