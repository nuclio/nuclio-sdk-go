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

type MemoryEvent struct {
	AbstractEvent
	Method      string
	ContentType string
	Body        []byte
	Headers     map[string]interface{}
	Path        string
}

func (me *MemoryEvent) GetMethod() string {
	if me.Method == "" {
		if len(me.Body) == 0 {
			return "GET"
		}
		return "POST"
	}
	return me.Method
}

func (me *MemoryEvent) GetContentType() string {
	if me.ContentType == "" {
		return "text/plain"
	}
	return me.ContentType
}

func (me *MemoryEvent) GetBody() []byte {
	return me.Body
}

func (me *MemoryEvent) GetPath() string {
	return me.Path
}

func (me *MemoryEvent) GetHeaders() map[string]interface{} {
	return me.Headers
}

func (me *MemoryEvent) GetHeader(key string) interface{} {
	if val, ok := me.Headers[key]; ok {
		return val
	}
	return ""
}
