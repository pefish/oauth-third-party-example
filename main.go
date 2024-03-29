package main

import (
	"fmt"
	"github.com/pefish/go-application"
	"github.com/pefish/go-config"
	"github.com/pefish/go-http"
	"github.com/pefish/go-logger"
	"github.com/pefish/go-redis"
	"oauth-third-party-example/service"
	"os"
	"runtime/debug"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
			os.Exit(1)
		}
		os.Exit(0)
	}()

	go_config.Config.LoadYamlConfig(go_config.Configuration{
		ConfigEnvName: `GO_CONFIG`,
		SecretEnvName: `GO_SECRET`,
	})

	go_http.Http.SetTimeout(20 * time.Second)

	service.LoginSvc.Init().SetHealthyCheck(nil)

	// 处理日志
	env := go_config.Config.GetString(`env`)
	go_application.Application.Debug = env == `local` || env == `dev`
	if go_application.Application.Debug {
		loggerInstance := go_logger.Log4goClass{}
		go_logger.Logger.InitWithLogger(&loggerInstance, service.LoginSvc.GetName(), `debug`)
	} else {
		loggerInstance := go_logger.LogrusClass{}
		go_logger.Logger.InitWithLogger(&loggerInstance, service.LoginSvc.GetName(), `info`)
	}

	go_redis.RedisHelper.ConnectWithMap(go_config.Config.GetMap(`redis`))
	defer go_redis.RedisHelper.Close()

	service.LoginSvc.SetHost(go_config.Config.GetString(`host`))
	service.LoginSvc.SetPort(go_config.Config.GetUint64(`port`))
	service.LoginSvc.Run()
}
