package apis

import "github.com/windrivder/gopkg/third"

const (
	WxPayServerUrl = "https://api.mch.weixin.qq.com"
)

var (
	AppUnifiedOrderAPI = third.Api{
		Name:        "APP 统一下单接口",
		Description: "微信支付 APP 统一下单接口",
		Path:        "/v3/pay/transactions/app",
		Method:      "POST",
		See:         "https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_2_1.shtml",
		FuncName:    "UnifiedOrder",
	}
)
