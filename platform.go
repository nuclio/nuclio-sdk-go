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

type Platform interface {
	CallFunction(name string, event FunctionCallEvent) (Response, error)
}

type FunctionCallEvent interface {
	GetMethod() string
	GetContentType() string
	GetBody() []byte
	GetPath() string
}

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
