package fission

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

type Fission struct {
	region     string
	name       string
	ingressUrl string
}

func New() (*Fission, error) {
	name := "fission"
	region := "eu-central-1"

	viper.SetConfigName("providers")
	viper.SetConfigType("json")
	viper.AddConfigPath("/app/provider/")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "reading providers config")
	}

	ingressUrl := viper.GetString(fmt.Sprintf("%s.ingress-url", name))

	return &Fission{region: region, name: name, ingressUrl: ingressUrl}, nil
}

func (fission *Fission) Name() string {
	return fission.name
}

func (fission *Fission) TLSConfig() *tls.Config {
	return nil
}

func (fission *Fission) buildFuncInvokeReq(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}
	funcUrl.Scheme = "http"
	funcUrl.Host = fission.ingressUrl
	funcUrl.Path = path.Join(funcUrl.Path, funcName)

	req, err := http.NewRequest("POST", funcUrl.String(), ioutil.NopCloser(bytes.NewReader(*body)))

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = qParams.Encode()

	for k, multiH := range *headers {
		for _, h := range multiH {
			req.Header.Set(k, h)
		}
	}

	return req, nil
}

func (fission *Fission) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func(uniqueId string) (*http.Request, error) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return fission.buildFuncInvokeReq(function.Name(), qParams, &localHeaders, body)
	}
}

func (fission *Fission) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.ConnectDone
}
