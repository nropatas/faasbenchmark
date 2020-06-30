package kubeless

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/golang/gddo/httputil/header"
	"github.com/nropatas/httpbench/syncedtrace"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Kubeless struct {
	region     string
	name       string
	ingressUrl string
}

func New() (*Kubeless, error) {
	name := "kubeless"
	region := "eu-central-1"

	viper.SetConfigName("providers")
	viper.SetConfigType("json")
	viper.AddConfigPath("/app/provider/")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "reading providers config")
	}

	ingressUrl := viper.GetString(fmt.Sprintf("%s.ingress-url", name))

	return &Kubeless{region: region, name: name, ingressUrl: ingressUrl}, nil
}

func (kubeless *Kubeless) Name() string {
	return kubeless.name
}

func (kubeless *Kubeless) TLSConfig() *tls.Config {
	return nil
}

func (kubeless *Kubeless) buildFuncInvokeReq(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}
	funcUrl.Scheme = "http"
	funcUrl.Host = kubeless.ingressUrl
	funcUrl.Path = path.Join(funcUrl.Path, funcName)

	req, err := http.NewRequest("POST", funcUrl.String(), ioutil.NopCloser(bytes.NewReader(*body)))

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = qParams.Encode()
	req.Host = "example.com"

	for k, multiH := range *headers {
		for _, h := range multiH {
			req.Header.Set(k, h)
		}
	}

	return req, nil
}

func (kubeless *Kubeless) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func(uniqueId string) (*http.Request, error) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return kubeless.buildFuncInvokeReq(function.Name(), qParams, &localHeaders, body)
	}
}

func (kubeless *Kubeless) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.ConnectDone
}
