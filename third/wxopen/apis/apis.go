package apis

import "github.com/windrivder/gopkg/third"

const (
	// 微信 api 服务器地址
	WXServerUrl = "https://api.weixin.qq.com"
)

var (
	GetAccessTokenAPI = third.Api{
		Name:        "获取 access_token",
		Description: "通过 code 获取 access_token",
		Request:     "GET https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code",
		Method:      "GET",
		Path:        "/sns/oauth2/access_token",
		See:         "https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html",
		FuncName:    "GetAccessToken",
		GetParams: []third.Param{
			{Name: `appid`, Type: `string`},
			{Name: `secret`, Type: `string`},
			{Name: `code`, Type: `string`},
			{Name: `grant_type`, Type: `string`},
		},
	}

	RefreshAccessTokenAPI = third.Api{
		Name:        "刷新 access_token",
		Description: "通过 access_token 获取到的 refresh_token 参数刷新",
		Request:     "GET https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN",
		Method:      "GET",
		Path:        "/sns/oauth2/refresh_token",
		See:         "https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html",
		GetParams: []third.Param{
			{Name: `appid`, Type: `string`},
			{Name: `grant_type`, Type: `string`},
			{Name: `refresh_token`, Type: `string`},
		},
	}

	GetUserInfoAPI = third.Api{
		Name:        "获取用户信息",
		Description: "获取用户个人信息",
		Request:     "GET https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID",
		Method:      "GET",
		Path:        "/sns/userinfo",
		See:         "https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Authorized_API_call_UnionID.html",
		GetParams: []third.Param{
			{Name: `access_token`, Type: `string`},
			{Name: `openid`, Type: `string`},
			{Name: `lang`, Type: `string`},
		},
	}
)
