package dubbogo

import (
	"fmt"
	"github.com/AlexStocks/goext/log"
	"github.com/AlexStocks/goext/net"
	"github.com/AlexStocks/goext/time"
	log "github.com/AlexStocks/log4go"
	"github.com/dubbo/dubbo-go/jsonrpc"
	"github.com/dubbo/dubbo-go/registry"
	jerrors "github.com/juju/errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"
)

var (
	survivalTimeout int = 3e9
	servo           *jsonrpc.Server
)

const (
	APP_CONF_FILE     string = "APP_CONF_FILE"
	APP_LOG_CONF_FILE string = "APP_LOG_CONF_FILE"
)

var (
	conf *ServerConfig
)

type (
	ServerConfig struct {
		// pprof
		Pprof_Enabled bool `default:"false" yaml:"pprof_enabled"  json:"pprof_enabled,omitempty"`
		Pprof_Port    int  `default:"10086"  yaml:"pprof_port" json:"pprof_port,omitempty"`

		// transport & registry
		Transport  string `default:"http"  yaml:"transport" json:"transport,omitempty"`
		NetTimeout string `default:"100ms"  yaml:"net_timeout" json:"net_timeout,omitempty"` // in ms
		netTimeout time.Duration
		// application
		Application_Config registry.ApplicationConfig `yaml:"application_config" json:"application_config,omitempty"`
		// Registry_Address  string `default:"192.168.35.3:2181"`
		Registry_Config registry.RegistryConfig  `yaml:"registry_config" json:"registry_config,omitempty"`
		Service_List    []registry.ServiceConfig `yaml:"service_list" json:"service_list,omitempty"`
		Server_List     []registry.ServerConfig  `yaml:"server_list" json:"server_list,omitempty"`
	}
)

func InitDubbo() *jsonrpc.Server {
	err := configInit()
	if err != nil {
		log.Error("configInit() = error{%#v}", err)
		panic(err)
	}
	initProfiling()

	servo = initServer()
	return servo
}

func DubboServiceRun() {
	servo.Start()
	initSignal()
}

func initServerConf() *ServerConfig {
	var (
		err      error
		confFile string
	)

	confFile = os.Getenv(APP_CONF_FILE)
	if confFile == "" {
		panic(fmt.Sprintf("application configure file name is nil"))
		return nil
	}
	if path.Ext(confFile) != ".yml" {
		panic(fmt.Sprintf("application configure file name{%v} suffix must be .yml", confFile))
		return nil
	}

	conf = &ServerConfig{}
	confFileStream, err := ioutil.ReadFile(confFile)
	if err != nil {
		panic(fmt.Sprintf("ioutil.ReadFile(file:%s) = error:%s", confFile, jerrors.ErrorStack(err)))
		return nil
	}
	err = yaml.Unmarshal(confFileStream, conf)
	if err != nil {
		panic(fmt.Sprintf("yaml.Unmarshal() = error:%s", jerrors.ErrorStack(err)))
		return nil
	}
	if conf.netTimeout, err = time.ParseDuration(conf.NetTimeout); err != nil {
		panic(fmt.Sprintf("time.ParseDuration(NetTimeout:%#v) = error:%s", conf.NetTimeout, err))
		return nil
	}
	if conf.Registry_Config.Timeout, err = time.ParseDuration(conf.Registry_Config.TimeoutStr); err != nil {
		panic(fmt.Sprintf("time.ParseDuration(Registry_Config.Timeout:%#v) = error:%s", conf.Registry_Config.TimeoutStr, err))
		return nil
	}

	gxlog.CInfo("config{%#v}\n", conf)

	return conf
}

func configInit() error {
	var (
		confFile string
	)

	initServerConf()

	confFile = os.Getenv(APP_LOG_CONF_FILE)
	if confFile == "" {
		panic(fmt.Sprintf("log configure file name is nil"))
		return nil
	}
	if path.Ext(confFile) != ".xml" {
		panic(fmt.Sprintf("log configure file name{%v} suffix must be .xml", confFile))
		return nil
	}

	log.LoadConfiguration(confFile)

	return nil
}

func initServer() *jsonrpc.Server {
	var (
		err            error
		serverRegistry *registry.ZkProviderRegistry
		srv            *jsonrpc.Server
	)

	if conf == nil {
		panic(fmt.Sprintf("conf is nil"))
		return nil
	}

	// registry
	serverRegistry, err = registry.NewZkProviderRegistry(
		registry.ApplicationConf(conf.Application_Config),
		registry.RegistryConf(conf.Registry_Config),
		registry.BalanceMode(registry.SM_RoundRobin),
		registry.ServiceTTL(conf.netTimeout),
	)
	if err != nil || serverRegistry == nil {
		panic(fmt.Sprintf("fail to init registry.Registy, err:%s", jerrors.ErrorStack(err)))
		return nil
	}

	// provider
	srv = jsonrpc.NewServer(
		jsonrpc.Registry(serverRegistry),
		jsonrpc.ConfList(conf.Server_List),
		jsonrpc.ServiceConfList(conf.Service_List),
	)

	return srv
}

func uninitServer() {
	if servo != nil {
		servo.Stop()
	}
	log.Close()
}

func initProfiling() {
	if !conf.Pprof_Enabled {
		return
	}
	const (
		PprofPath = "/debug/pprof/"
	)
	var (
		err  error
		ip   string
		addr string
	)

	ip, err = gxnet.GetLocalIP()
	if err != nil {
		panic("cat not get local ip!")
	}
	addr = ip + ":" + strconv.Itoa(conf.Pprof_Port)
	log.Info("App Profiling startup on address{%v}", addr+PprofPath)

	go func() {
		log.Info(http.ListenAndServe(addr, nil))
	}()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		log.Info("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		// reload()
		default:
			go gxtime.Future(survivalTimeout, func() {
				log.Warn("app exit now by force...")
				os.Exit(1)
			})

			// 要么fastFailTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			uninitServer()
			fmt.Println("provider app exit now...")
			return
		}
	}
}
