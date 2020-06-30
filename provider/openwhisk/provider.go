package openwhisk

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

type OpenWhisk struct {
	region     string
	name       string
	ingressUrl string
}

func New() (*OpenWhisk, error) {
	name := "openwhisk"
	region := "eu-central-1"

	viper.SetConfigName("providers")
	viper.SetConfigType("json")
	viper.AddConfigPath("/app/provider/")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "reading providers config")
	}

	ingressUrl := viper.GetString(fmt.Sprintf("%s.ingress-url", name))

	return &OpenWhisk{region: region, name: name, ingressUrl: ingressUrl}, nil
}

func (openwhisk *OpenWhisk) Name() string {
	return openwhisk.name
}

func (openwhisk *OpenWhisk) TLSConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

func (openwhisk *OpenWhisk) buildFuncInvokeReq(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}
	funcUrl.Scheme = "https"
	funcUrl.Host = openwhisk.ingressUrl
	funcUrl.Path = path.Join(funcUrl.Path, "api", "v1", "web", "guest", "default", funcName)

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

func (openwhisk *OpenWhisk) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func(uniqueId string) (*http.Request, error) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return openwhisk.buildFuncInvokeReq(function.Name(), qParams, &localHeaders, body)
	}
}

func (openwhisk *OpenWhisk) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}
