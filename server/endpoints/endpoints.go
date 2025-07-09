// Package endpoints holds logic of api endpoint
package endpoints

import "github.com/labstack/echo/v4"

type endpoint struct {
	method      string
	path        string
	handler     echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

var endpoints []endpoint

func registerEndpoint(e endpoint) {
	endpoints = append(endpoints, e)
}

func PopulateEndpoints(e *echo.Echo) {
	for _, i := range endpoints {
		e.Add(i.method, i.path, i.handler, i.middlewares...)
	}
}
