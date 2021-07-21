package apis

import (
	"io/ioutil"
	"net/url"

	"github.com/windrivder/gopkg/encoding/jsonx"
	"github.com/windrivder/gopkg/encoding/urlx"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/logx"
	"github.com/windrivder/gopkg/third/wxopen"
)

type WechatResponse struct {
	ErrCode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Openid       string `json:"openid,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Unionid      string `json:"unionid,omitempty"`
	WechatResponse
}

var defaultAccessToken = AccessToken{
	AccessToken:  "ACCESS_TOKEN",
	ExpiresIn:    7200,
	RefreshToken: "REFRESH_TOKEN",
	Openid:       "OPENID",
	Scope:        "SCOPE",
	Unionid:      "o6_bmasdasdsad6_2sgVt7hMZOPfL",
}

func GetAccessToken(p *wxopen.Platform, code string) (*AccessToken, error) {
	params := url.Values{}
	params.Add("appid", p.Options.AppId)
	params.Add("secret", p.Options.AppSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	u, err := urlx.BuildURL(WXServerUrl+GetAccessTokenAPI.Path, params)
	if err != nil {
		return nil, errorx.Wrap(err, "build url error")
	}
	resp, err := p.Client.Get(u, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "get wechat access token error")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accessToken AccessToken
	err = jsonx.Decode(data, &accessToken)
	if err != nil {
		return nil, err
	}

	if p.Options.Debug {
		return &defaultAccessToken, nil
	}

	return &accessToken, nil
}

type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

func RefreshAccessToken(p *wxopen.Platform, refresh_token string) (*RefreshToken, error) {
	params := url.Values{}
	params.Add("appid", p.Options.AppId)
	params.Add("grant_type", "authorization_code")
	params.Add("refresh_token", refresh_token)

	u, err := urlx.BuildURL(WXServerUrl+RefreshAccessTokenAPI.Path, params)
	if err != nil {
		return nil, errorx.Wrap(err, "build url error")
	}
	resp, err := p.Client.Get(u, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "refresh wechat access token error")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var refreshToken RefreshToken
	err = jsonx.Decode(data, &refreshToken)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

type UserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int64    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
	WechatResponse
}

var defaultUserInfo = UserInfo{
	Openid:     "OPENID",
	Nickname:   "dolabox",
	Sex:        1,
	Province:   "四川省",
	City:       "成都市",
	Country:    "未知",
	Headimgurl: "https://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
	Privilege:  []string{"PRIVILEGE1", "PRIVILEGE2"},
	Unionid:    "o6_bmasdasdsad6_2sgVt7hMZOPfL",
}

func (a UserInfo) MarshalBinary() (data []byte, err error) {
	return jsonx.Encode(a)
}

func GetUserInfo(p *wxopen.Platform, access_token, openid string) (*UserInfo, error) {
	params := url.Values{}
	params.Add("access_token", access_token)
	params.Add("openid", openid)
	params.Add("lang", "zh_CN")

	u, err := urlx.BuildURL(WXServerUrl+GetUserInfoAPI.Path, params)
	if err != nil {
		return nil, errorx.Wrap(err, "build url error")
	}

	logx.Info().Msgf("request get user info: %+v", u)

	resp, err := p.Client.Get(u, nil)
	if err != nil {
		return nil, errorx.Wrap(err, "get wechat user info error")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user UserInfo
	err = jsonx.Decode(data, &user)
	if err != nil {
		return nil, err
	}

	if p.Options.Debug {
		return &defaultUserInfo, nil
	}

	return &user, nil
}
