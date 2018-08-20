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

	"github.com/nuclio/nuclio/pkg/errors"

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

func (p *Platform) CallFunction(functionName string, event FunctionCallEvent) (*Response, error) {
	request := p.createRequest(functionName, event)

	response := fasthttp.AcquireResponse()

	err := p.client.Do(request, response)
	fasthttp.ReleaseRequest(request)
	if err != nil {
		fasthttp.ReleaseResponse(response)
		return nil, errors.Wrap(err, "Failed to call function")
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

func (p *Platform) createRequest(functionName string, event FunctionCallEvent) *fasthttp.Request {
	request := fasthttp.AcquireRequest()
	request.URI().SetScheme("http")
	request.URI().SetHost(p.getFunctionHost(functionName))
	request.URI().SetPath(event.GetPath())
	request.SetBody(event.GetBody())
	request.Header.SetContentType(event.GetContentType())
	request.Header.SetMethod(event.GetMethod())
	return request
}

func (p *Platform) createResponse(response *fasthttp.Response) *Response {
	result := &Response{}
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

// Special interface for CallFunction in Platform.
// It allows user not to implement all functions in Event interface for using Platform.CallFunction.
// FunctionCallEvent is sub-interface of Event and includes only necessary functions.
type FunctionCallEvent interface {
	GetMethod() string
	GetContentType() string
	GetBody() []byte
	GetPath() string
}

// Basic implementation of FunctionCallEvent interface with few defaults.
type BasicFunctionCallEvent struct {
	Method      string
	ContentType string
	Body        []byte
	Path        string
}

func (be BasicFunctionCallEvent) GetMethod() string {
	if be.Method == "" {
		return "GET"
	}
	return be.Method
}

func (be BasicFunctionCallEvent) GetContentType() string {
	if be.ContentType == "" {
		return "text/plain"
	}
	return be.ContentType
}

func (be BasicFunctionCallEvent) GetBody() []byte {
	return be.Body
}

func (be BasicFunctionCallEvent) GetPath() string {
	return be.Path
}
