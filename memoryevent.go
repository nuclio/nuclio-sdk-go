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
		} else {
			return "POST"
		}
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
