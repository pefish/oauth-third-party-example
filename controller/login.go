package controller

import (
	"fmt"
	"github.com/pefish/go-config"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-error"
	"github.com/pefish/go-http"
	"github.com/pefish/go-random"
	"github.com/pefish/go-redis"
	"time"
)

type LoginControllerClass struct {
}

var LoginController = LoginControllerClass{}

type LoginCallbackParam struct {
	Code  *string `json:"code" validate:"omitempty" desc:"授权码"`
	Scope *string `json:"scope" validate:"omitempty" desc:"scope"`
	State *string `json:"state" validate:"omitempty" desc:"state"`

	Error            *string `json:"error" validate:"omitempty" desc:"error"`
	ErrorDescription *string `json:"error_description" validate:"omitempty" desc:"error_description"`
}

func (this *LoginControllerClass) LoginCallback(apiSession *api_session.ApiSessionClass) interface{} {
	params := LoginCallbackParam{}
	apiSession.ScanParams(&params)

	apiSession.Ctx.Header("Content-Type", "text/html; charset=utf-8")

	if params.Error != nil {
		str := `
	<html><title>title</title>
		<body>
			<h1 style="color: red;">Error: %s</h1>
			<h1 style="color: red;">ErrorDescription: %s</h1>
		</body>
	</html>
`
		apiSession.Ctx.Write([]byte(fmt.Sprintf(str, *params.Error, *params.ErrorDescription)))
		return nil
	}

	if params.Code != nil {
		if go_redis.RedisHelper.String.Get(*params.State) == `` {
			go_error.Throw(`state error`, 2000)
		}
		go_redis.RedisHelper.Del(*params.State)
		map_ := go_http.Http.PostMultipartForMap(go_http.PostMultipartParam{
			Url: go_config.Config.GetString(`/oauth/serverUrl`) + `/oauth2/token`,
			Params: map[string]interface{}{
				`grant_type`:   `authorization_code`,
				`client_id`:    go_config.Config.GetString(`/oauth/clientId`),
				`redirect_uri`: `http://127.0.0.1:8000/home`,
				`code`:         *params.Code,
			},
			BasicAuth: &go_http.BasicAuth{
				Username: go_config.Config.GetString(`/oauth/clientId`),
				Password: go_config.Config.GetString(`/oauth/clientSecret`),
			},
		})
		str := `
	<html><title>title</title>
		<body>
			<h2 style="color: red;">access_token: %s</h2>
			<h2 style="color: red;">expires_in: %f</h2>
			<h2 style="color: red;">id_token: %s</h2>
			<h2 style="color: red;">refresh_token: %s</h2>
			<h2 style="color: red;">scope: %s</h2>
			<h2 style="color: red;">token_type: %s</h2>
		</body>
	</html>
`
		apiSession.Ctx.Write([]byte(fmt.Sprintf(
			str,
			map_[`access_token`],
			map_[`expires_in`],
			map_[`id_token`],
			map_[`refresh_token`],
			map_[`scope`],
			map_[`token_type`])))
		return nil
	}

	str := `
	<html><title>title</title>
		<body>
			<h1 style="color: red;">主页</h1>
		</body>
	</html>
`
	apiSession.Ctx.Write([]byte(str))
	return nil
}

func (this *LoginControllerClass) getAuthUrl(state string) string {
	// max_age指定ID token的过期时间
	// scope必须包含openid
	return fmt.Sprintf(
		`%s/oauth2/auth?max_age=0&client_id=%s&redirect_uri=%s&response_type=code&scope=openid+offline+test&state=%s`,
		go_config.Config.GetString(`/oauth/serverUrl`),
		go_config.Config.GetString(`/oauth/clientId`),
		go_config.Config.GetString(`/oauth/callbackUrl`),
		state)
}

func (this *LoginControllerClass) LoginGet(apiSession *api_session.ApiSessionClass) interface{} {
	apiSession.Ctx.Header("Content-Type", "text/html; charset=utf-8")
	state := go_random.Random.GetUniqueIdString()
	authUrl := this.getAuthUrl(state)
	go_redis.RedisHelper.String.Set(state, `1`, 5*time.Minute)
	str := `
	<html><title>title</title>
		<body>
			<a href="%s" style="color: red;">ZG登录</a>
		</body>
	</html>
`
	apiSession.Ctx.Write([]byte(fmt.Sprintf(str, authUrl)))
	return nil
}


