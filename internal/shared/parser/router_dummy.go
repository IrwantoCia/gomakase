package parser

import (
	"fmt"
)

func Routes() {
	var router = NewRouter()
	_ = "bar"

	router.GET("/", func() { fmt.Println("Hello, World!") })

}

type Router interface {
	GET(path string, handler func())
}

type router struct {
}

func NewRouter() Router {
	return &router{}
}

func (r *router) GET(path string, handler func()) {
	handler()
}
