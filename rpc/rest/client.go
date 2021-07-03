package rest

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/encoding/jsonx"
	"github.com/windrivder/gopkg/vars"
)

type (
	Client = httpclient.Client
)

func NewClient(o Options) (*Client, error) {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(o.ClientTimeout * time.Second),
	)

	return client, nil
}

func PostJSON(client *Client, url string, data typex.GenericType) (*http.Response, error) {
	body, err := jsonx.Encode(data)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set(vars.HeaderContentType, vars.MIMEApplicationJSONCharsetUTF8)
	headers.Set(vars.HeaderAccept, vars.MIMEApplicationJSONCharsetUTF8)

	return client.Post(url, bytes.NewBuffer(body), headers)
}
