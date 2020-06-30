package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nropatas/httpbench/engine"
	"github.com/nropatas/httpbench/syncedtrace"
	"github.com/nuweba/faasbenchmark/provider/aws"
	"github.com/nuweba/faasbenchmark/provider/azure"
	"github.com/nuweba/faasbenchmark/provider/fission"
	"github.com/nuweba/faasbenchmark/provider/google"
	"github.com/nuweba/faasbenchmark/provider/knative"
	"github.com/nuweba/faasbenchmark/provider/kubeless"
	"github.com/nuweba/faasbenchmark/provider/openfaas"
	"github.com/nuweba/faasbenchmark/provider/openwhisk"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/pkg/errors"
)

type RequestFilter = func(sleepTime time.Duration, tr *engine.TraceResult, funcDuration time.Duration, reused bool) (report.Result, error)

type Filter interface {
	HttpResult(sleepTime time.Duration, tr *engine.TraceResult, funcDuration time.Duration, reused bool) (report.Result, error)
}

type FaasProvider interface {
	Filter
	Name() string
	TLSConfig() *tls.Config
	HttpInvocationTriggerStage() syncedtrace.TraceHookType
	NewStack(stackPath string) (stack.Stack, error)
	NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func(uniqueId string) (*http.Request, error)
}

type Providers int

const (
	AWS Providers = iota
	Google
	Azure
	Knative
	OpenFaaS
	OpenWhisk
	Kubeless
	Fission
	ProvidersCount
)

func (p Providers) String() string {
	return [...]string{
		"aws",
		"google",
		"azure",
		"knative",
		"openfaas",
		"openwhisk",
		"kubeless",
		"fission",
	}[p]
}

func (p Providers) Description() string {
	return [...]string{
		"aws lambda functions",
		"google cloud functions",
		"azure functions",
		"knative",
		"openfaas",
		"openwhisk",
		"kubeless",
		"fission",
	}[p]
}

func NewProvider(providerName string) (FaasProvider, error) {
	var faasProvider FaasProvider
	var err error

	switch strings.ToLower(providerName) {
	case strings.ToLower(AWS.String()):
		faasProvider, err = aws.New()
	case strings.ToLower(Google.String()):
		faasProvider, err = google.New()
	case strings.ToLower(Azure.String()):
		faasProvider, err = azure.New()
	case strings.ToLower(Knative.String()):
		faasProvider, err = knative.New()
	case strings.ToLower(OpenFaaS.String()):
		faasProvider, err = openfaas.New()
	case strings.ToLower(OpenWhisk.String()):
		faasProvider, err = openwhisk.New()
	case strings.ToLower(Kubeless.String()):
		faasProvider, err = kubeless.New()
	case strings.ToLower(Fission.String()):
		faasProvider, err = fission.New()
	default:
		faasProvider, err = nil, errors.New(fmt.Sprintf("provider not supported: %s", providerName))
	}

	if err != nil {
		return nil, err
	}

	return faasProvider, nil
}

func List() []string {
	var providers []string
	for providerId := Providers(0); providerId < ProvidersCount; providerId++ {
		providers = append(providers, providerId.String())
	}

	return providers
}
