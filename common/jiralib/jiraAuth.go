package jiralib

import "net/http"

type Authentication interface {
	PreAuth()
	GetHeader() Header
}

type Auth struct {
	Header   Header
	AuthType string
}

type Header struct {
	Name  string
	Value string
}

func (header Header) SetHeader(req *http.Request) {
	req.Header.Set(header.Name, header.Value)
	req.Header.Set("Content-Type", "application/json")
}
