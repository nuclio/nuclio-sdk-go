/*
Copyright 2017 The Nuclio Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nuclio

import (
	"fmt"
	"strconv"

	"github.com/nuclio/logger"
	"github.com/valyala/fasthttp"
)

type Platform struct {
	client    fasthttp.Client
	logger    logger.Logger
	kind      string
	namespace string
}

func NewPlatform(parentLogger logger.Logger, kind string, namespace string) (*Platform, error) {
	return &Platform{
		client:    fasthttp.Client{},
		logger:    parentLogger.GetChild("platform"),
		kind:      kind,
		namespace: namespace,
	}, nil
}

func (p *Platform) CallFunction(functionName string, event Event) (Response, error) {
	var emptyResponse Response

	request := p.createRequest(functionName, event)

	response := fasthttp.AcquireResponse()

	err := p.client.Do(request, response)
	fasthttp.ReleaseRequest(request)
	if err != nil {
		fasthttp.ReleaseResponse(response)
		return emptyResponse, err
	}
	wrappedResponse := p.createResponse(response)
	fasthttp.ReleaseResponse(response)
	return wrappedResponse, nil
}

func (p *Platform) getFunctionHost(name string) string {
	if p.kind == "local" {
		return fmt.Sprintf("%s-%s:8080", p.namespace, name)
	}
	return fmt.Sprintf("%s:8080", name)
}

func (p *Platform) createRequest(functionName string, event Event) *fasthttp.Request {
	request := fasthttp.AcquireRequest()
	request.URI().SetScheme("http")
	request.URI().SetHost(p.getFunctionHost(functionName))
	request.URI().SetPath(event.GetPath())
	request.SetBody(event.GetBody())
	request.Header.SetContentType(event.GetContentType())
	request.Header.SetMethod(event.GetMethod())

	for headerKey, headerValue := range event.GetHeaders() {
		switch typedHeaderValue := headerValue.(type) {
		case string:
			request.Header.Set(headerKey, typedHeaderValue)

		case int:
			request.Header.Set(headerKey, strconv.Itoa(typedHeaderValue))

		case bool:
			request.Header.Set(headerKey, strconv.FormatBool(typedHeaderValue))

		case []byte:
			request.Header.Set(headerKey, string(typedHeaderValue))

		default:
			p.logger.WarnWith("Header value is of an unsupported type. Ignoring it",
				"headerKey",
				headerKey,
				"headerValue",
				headerValue)
		}
	}

	return request
}

func (p *Platform) createResponse(response *fasthttp.Response) Response {
	result := Response{}
	if len(response.Header.ContentType()) == 0 {
		result.ContentType = "text/plain"
	} else {
		result.ContentType = string(response.Header.ContentType())
	}

	result.StatusCode = response.StatusCode()

	result.Headers = make(map[string]interface{}, response.Header.Len())
	response.Header.VisitAll(func(key, value []byte) {
		result.Headers[string(key)] = string(value)
	})

	result.Body = append(result.Body, response.Body()...)

	return result
}
