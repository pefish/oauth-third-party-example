package service

import (
	"github.com/pefish/go-core/api-channel-builder"
	"github.com/pefish/go-core/api-strategy"
	"github.com/pefish/go-core/service"
	"oauth-login-consent/constant"
	"oauth-login-consent/controller"
)

type LoginSvcClass struct {
	service.BaseServiceClass
}

var LoginSvc = LoginSvcClass{}

func (this *LoginSvcClass) Init(opts ...interface{}) service.InterfaceService {
	this.SetName(`oauth三方登录示例`)
	api_strategy.ParamValidateStrategy.SetErrorCode(constant.PARAM_ERROR)

	this.SetRoutes(map[string]*api_channel_builder.Route{
		`login callback`: {
			Description: "主页",
			Path:        "/home",
			Method:      "GET",
			Strategies: []api_channel_builder.StrategyRoute{

			},
			ParamType: api_strategy.ALL_TYPE,
			Params: controller.LoginCallbackParam{

			},
			Controller: controller.LoginController.LoginCallback,
		},
		`login get`: {
			Description: "登录页面",
			Path:        "/login",
			Method:      "GET",
			Strategies: []api_channel_builder.StrategyRoute{

			},
			ParamType: api_strategy.ALL_TYPE,
			Controller: controller.LoginController.LoginGet,
		},
	})
	return this
}
