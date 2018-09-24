package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Away0x/api_server/config"
	"github.com/Away0x/api_server/model"
	v "github.com/Away0x/api_server/pkg/version"
	"github.com/Away0x/api_server/router"
	"github.com/Away0x/api_server/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "C", "", "apiserver config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	// 解析命令行参数
	pflag.Parse()

	// api version
	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}

	// 初始化配置
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// 设置 gin 的运行模式
	// gin 有 3 种运行模式：debug、release 和 test，其中 debug 模式会打印很多 debug 信息
	gin.SetMode(viper.GetString("runmode"))

	// 创建 gin.Engine
	g := gin.New()

	// 加载路由和中间件
	router.Load(
		// Cores.
		g,

		// Middlwares.
		//middleware.Logging(),
		middleware.RequestId(),
	)

	// 在启动 HTTP 端口前 go 一个 pingServer 协程，启动 HTTP 端口后，
	// 该协程不断地 ping /sd/health 路径，如果失败次数超过一定次数，则终止 HTTP 服务器进程
	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	// 如果提供了 TLS 证书和私钥则启动 HTTPS 端口 (8081)
	// Start to listening the incoming requests.
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	// HTTP 端口 8080
	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
