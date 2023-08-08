package main

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
)

type Config struct {
    BodyKey string `json:"bodyKey,omitempty"`
    HeaderKey string `json:"headerKey,omitempty"`
}

type attachHeaderFromBodyMiddleware struct {
    next   http.Handler
    config *Config
}

func New(next http.Handler, config *Config) (http.Handler, error) {
    return &attachHeaderFromBodyMiddleware{next: next, config: config}, nil
}

func (m *attachHeaderFromBodyMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
    if err != nil {
        // If there's an error reading the body, just pass the request to the next handler
        m.next.ServeHTTP(rw, req)
        return
    }

	var bodyMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &bodyMap)
	if err == nil {
		if value, ok := bodyMap[m.config.BodyKey].(string); ok {
			req.Header.Set(m.config.HeaderKey, value)
		}
	}
	
	req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

    
    m.next.ServeHTTP(rw, req)
}
