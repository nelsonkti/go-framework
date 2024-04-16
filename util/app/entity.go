package app

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HandlerFunc func(*Context)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         httprouter.Params
}
