package wxpay

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"time"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/windrivder/gopkg/errorx"
)

type Options struct {
	AppId               string        `json:"Appid"`
	MchId               string        `json:"Mchid"`               // 商户 id
	ApiKey              string        `json:"ApiKey"`              // 商户 api key
	MchCertSerialNumber string        `json:"MchCertSerialNumber"` // 商户证书序列号
	PrivateKeyPath      string        `json:"PrivateKeyPath"`      // 商户私钥文件路径
	WechatCertPath      string        `json:"WechatCertPath"`
	WechatCertSerialNo  string        `json:"WechatCertSerialNo"`
	Timeout             time.Duration `json:"Timeout"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("Wxpay", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal wxpay option error")
	}

	return o, err
}

type WxPay struct {
	Options    Options
	Client     *core.Client
	PrivateKey *rsa.PrivateKey
}

func New(o Options) (*WxPay, error) {
	privateKey, err := utils.LoadPrivateKeyWithPath(o.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	wechatCert, err := utils.LoadCertificateWithPath(o.WechatCertPath)
	if err != nil {
		return nil, err
	}

	// 增加客户端配置
	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithMerchant(o.MchId, o.MchCertSerialNumber, privateKey), // 设置商户信息，用于生成签名信息
		option.WithWechatPay([]*x509.Certificate{wechatCert}),           // 设置微信支付平台证书信息，对回包进行校验
		option.WithTimeout(o.Timeout * time.Second),                     // 自行进行超时时间配置
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &WxPay{Options: o, Client: client, PrivateKey: privateKey}, nil
}

var ProviderSet = wire.NewSet(NewOptions, New)
