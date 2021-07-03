package apis

import (
	"context"
	"io/ioutil"

	"github.com/windrivder/gopkg/encoding/jsonx"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/third/wxpay"
)

type UnifiedOrderAmount struct {
	Total int `json:"total"` // 订单总金额，单位为分。
}

type UnifiedOrderReq struct {
	AppId      string `json:"appid"`
	MchId      string `json:"mchid"`
	Desc       string `json:"description"`
	OutTradeNo string `json:"out_trade_no"`
	TimeExpire string `json:"time_expire,omitempty"`
	Attach     string `json:"attach,omitempty"`    // 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用
	NotifyUrl  string `json:"notify_url"`          // 通知URL必须为直接可访问的URL，不允许携带查询串
	GoodsTag   string `json:"goods_tag,omitempty"` // 订单优惠标记

	// 订单金额 amount
	Amount UnifiedOrderAmount `json:"amount"`
}

type UnifiedOrderReply struct {
	PrepayId string `json:"prepay_id"`
}

func UnifiedOrder(ctx context.Context, wxPay *wxpay.WxPay, params UnifiedOrderReq) (*UnifiedOrderReply, error) {
	params.AppId = wxPay.Options.AppId
	params.MchId = wxPay.Options.MchId

	url := WxPayServerUrl + AppUnifiedOrderAPI.Path
	resp, err := wxPay.Client.Post(ctx, url, params)
	if err != nil {
		return nil, errorx.Wrap(err, "wxpay unified order error")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reply := &UnifiedOrderReply{}
	if err := jsonx.Decode(bytes, reply); err != nil {
		return nil, err
	}

	return reply, nil
}
