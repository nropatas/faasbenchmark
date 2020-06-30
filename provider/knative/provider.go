package knative

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/gddo/httputil/header"
	"github.com/nropatas/httpbench/syncedtrace"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Knative struct {
	region     string
	name       string
	ingressUrl string
}

func New() (*Knative, error) {
	name := "knative"
	region := "eu-central-1"

	viper.SetConfigName("providers")
	viper.SetConfigType("json")
	viper.AddConfigPath("/app/provider/")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "reading providers config")
	}

	ingressUrl := viper.GetString(fmt.Sprintf("%s.ingress-url", name))

	return &Knative{region: region, name: name, ingressUrl: ingressUrl}, nil
}

func (knative *Knative) Name() string {
	return knative.name
}

func (knative *Knative) TLSConfig() *tls.Config {
	return nil
}

func (knative *Knative) buildFuncInvokeReq(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{
		Scheme: "http",
		Host:   knative.ingressUrl,
	}

	req, err := http.NewRequest("POST", funcUrl.String(), ioutil.NopCloser(bytes.NewReader(*body)))

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = qParams.Encode()
	req.Host = fmt.Sprintf("%s.default.example.com", funcName)

	for k, multiH := range *headers {
		for _, h := range multiH {
			req.Header.Set(k, h)
		}
	}

	return req, nil
}

func (knative *Knative) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func(uniqueId string) (*http.Request, error) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return knative.buildFuncInvokeReq(function.Name(), qParams, &localHeaders, body)
	}
}

func (knative *Knative) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.ConnectDone
}
