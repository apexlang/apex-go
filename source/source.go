package source

type Source struct {
	Body []byte `json:"body,omitempty"`
	Name string `json:"name,omitempty"`
}

func NewSource(name string, body []byte) *Source {
	return &Source{
		Name: name,
		Body: body,
	}
}
