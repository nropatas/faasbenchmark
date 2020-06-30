package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/nropatas/httpbench"
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
)

func init() {
	Tests.Register(Test{Id: "LargeResponse", Fn: largeResponse, RequiredStack: "largeresponse", Description: "Benchmark the response time of a function returning a large (4mb) response body. Invoke the function once at a time for a minute."})
}

func largeResponse(test *config.Test) {
	headers := http.Header{}
	body := []byte{}
	queryParams := url.Values{}
	httpConfig := config.Http{
		QueryParams:      &queryParams,
		Headers:          &headers,
		Body:             &body,
		TestType:         httpbench.ConcurrentRequestsSynced.String(),
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		ConcurrencyLimit: 1,
		RequestDelay:     time.Millisecond,
		Duration:         time.Minute,
	}
	httpConfig.QueryParams = &queryParams
	httpConfig.Headers = &headers
	httpConfig.Body = &body

	if test.Config.HasCustomHttpConfig() {
		err := json.Unmarshal(*test.Config.CustomHttpConfig, &httpConfig)
		if err != nil {
			fmt.Println("Failed to read custom HTTP config:", err.Error())
		}
	}

	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(&httpConfig, function)

		if err != nil {
			continue
		}

		newReq := hfConf.Test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
		tlsConfig := hfConf.Test.Config.Provider.TLSConfig()
		wg := &sync.WaitGroup{}
		trace := httpbench.New(newReq, hfConf.HttpConfig.Hook, tlsConfig)
		wg.Add(1)
		go func() {
			defer wg.Done()
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpResult)
		}()
		requestsResult := trace.ConcurrentRequestsSynced(httpConfig.ConcurrencyLimit, httpConfig.RequestDelay, httpConfig.Duration)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
