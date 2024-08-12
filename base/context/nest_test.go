package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	var connHandler = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx.Value(NestCtxKey).(*ContextInfo).Params["LocalAddrContextKey"] = "哈哈"
	}

	var authHandler = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "123" {
			ctx.Value(NestCtxKey).(*ContextInfo).Params["Valid"] = true
		}
	}

	nest := New()
	nest.UseFunc(connHandler)
	nest.UseFunc(authHandler)

	mux := httprouter.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()
		params := ctx.Value(NestCtxKey).(*ContextInfo).Params
		localAddr := params["LocalAddrContextKey"].(string)
		valid, ok := params["Valid"]
		if !ok {
			valid = false
		}

		w.Write([]byte(fmt.Sprintf("hello world. localAddr: %s, valid: %t", localAddr, valid)))
	})

	nest.SetMux(mux)
	nest.Run(":9999")
}
