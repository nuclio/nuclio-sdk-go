package nuclio

type MemoryEvent struct {
	AbstractEvent
	Method      string
	ContentType string
	Body        []byte
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
