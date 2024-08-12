package main

import (
	"context"
	"net/http"
)

// Context实战 - web框架
type contextKey struct {
	name string
}

type ContextInfo struct {
	Params map[string]any
}

var NestCtxKey = contextKey{
	name: "NestCtxKey",
}

// NewCtx Key需要自定义，不能直接用string，否则在包与包之间使用的过程中，可能会出现key冲突
func NewCtx() context.Context {
	ctx := context.WithValue(context.Background(), NestCtxKey, &ContextInfo{
		Params: make(map[string]any),
	})
	return ctx
}

// Handler 中间件
type Handler interface {
	ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (fn HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fn(ctx, w, r)
}

type Nest struct {
	middleware []Handler
	mux        http.Handler
}

func New() *Nest {
	return &Nest{
		mux: http.DefaultServeMux,
	}
}

func (nest *Nest) Use(handler Handler) {
	nest.middleware = append(nest.middleware, handler)
}

func (nest *Nest) UseFunc(handlerFunc HandlerFunc) {
	nest.Use(handlerFunc)
}

func (nest *Nest) SetMux(mux http.Handler) {
	nest.mux = mux
}

func (nest *Nest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewCtx()

	for _, handler := range nest.middleware {
		handler.ServeHTTP(ctx, w, r)
	}

	nest.mux.ServeHTTP(w, r.WithContext(ctx))
}

func (nest *Nest) Run(addr string) error {
	return http.ListenAndServe(addr, nest)
}
