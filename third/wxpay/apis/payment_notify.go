package apis

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/windrivder/gopkg/encoding/jsonx"
	"github.com/windrivder/gopkg/third/wxpay"
)

var (
	SuccessPaymentNofityReply = PaymentNofityReply{
		Code:    "SUCCESS",
		Message: "成功",
	}
	FailPaymentNofityReply = PaymentNofityReply{
		Code:    "FAIL",
		Message: "失败",
	}
)

type PaymentNotifyResourceReq struct {
	Algorithm      string `json:"algorithm"`       // 加密算法类型
	Ciphertext     string `json:"ciphertext"`      // 数据密文
	AssociatedData string `json:"associated_data"` // 附加数据
	OriginalType   string `json:"original_type"`   // 原始类型
	Nonce          string `json:"nonce"`           // 随机串
}

type PaymentNofityReq struct {
	Id           string                   `json:"id"`            // 通知id
	CreateTime   string                   `json:"create_time"`   // 通知创建时间
	EventType    string                   `json:"event_type"`    // 通知类型
	ResourceType string                   `json:"resource_type"` // 通知数据类型
	Resource     PaymentNotifyResourceReq `json:"resource"`      // 资源
	Summary      string                   `json:"summary"`       // 回调摘要
}

type PaymentNofityReply struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 解密后的资源对象

type PaymentNofityResourcePayer struct {
	Openid string `json:"openid"`
}

type PaymentNofityResourceAmount struct {
	Total         int    `json:"total"`
	PayerTotal    int    `json:"payer_total"`
	Currency      string `json:"currency"`
	PayerCurrency string `json:"payer_currency"`
}

type PaymentNofityResourceSceneInfo struct {
	DeviceId string `json:"device_id"`
}

type PaymentNofityResource struct {
	Appid          string                         `json:"appid"`
	Mchid          string                         `json:"mchid"`
	OutTradeNo     string                         `json:"out_trade_no"`
	TransactionId  string                         `json:"transaction_id"`
	TradeType      string                         `json:"trade_type"`
	TradeState     string                         `json:"trade_state"`
	TradeStateDesc string                         `json:"trade_state_desc"`
	BankType       string                         `json:"bank_type"`
	Attach         string                         `json:"attach"`
	SuccessTime    string                         `json:"success_time"`
	Payer          PaymentNofityResourcePayer     `json:"payer"`
	Amount         PaymentNofityResourceAmount    `json:"amount"`
	SceneInfo      PaymentNofityResourceSceneInfo `json:"scene_info"`
}

func PaymentNofityDecrypt(wxPay *wxpay.WxPay, req PaymentNotifyResourceReq) (*PaymentNofityResource, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher([]byte(wxPay.Options.ApiKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, []byte(req.Nonce), decodedCiphertext, []byte(req.AssociatedData))
	if err != nil {
		return nil, err
	}

	reply := &PaymentNofityResource{}
	if err := jsonx.Decode(plaintext, reply); err != nil {
		return nil, err
	}

	return reply, nil
}
