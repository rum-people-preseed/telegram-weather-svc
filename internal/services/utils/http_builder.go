package utils

import "net/url"

type HTTPBuilder struct {
	URL    string
	params []HTTPParam
}

type HTTPParam struct {
	key   string
	value string
}

func NewHTTPBuilder(URL string) *HTTPBuilder {
	return &HTTPBuilder{
		URL:    URL,
		params: make([]HTTPParam, 0),
	}
}

func NewHTTPParam(key string, value string) *HTTPParam {
	return &HTTPParam{
		key:   key,
		value: value,
	}
}

func (b *HTTPBuilder) AddParams(params ...*HTTPParam) *HTTPBuilder {
	for _, param := range params {
		b.params = append(b.params, *param)
	}
	return b
}

func (b *HTTPBuilder) Build() string {
	u, _ := url.Parse(b.URL)

	baseQuery := u.Query()

	for _, param := range b.params {
		baseQuery.Add(param.key, param.value)
	}

	u.RawQuery = baseQuery.Encode()
	return u.String()
}

func BuildURL(URL string, params ...*HTTPParam) string {
	httpBuilder := NewHTTPBuilder(URL)
	httpBuilder.AddParams(params...)
	return httpBuilder.Build()
}
