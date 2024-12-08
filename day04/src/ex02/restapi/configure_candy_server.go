// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import "C"
import (
	"crypto/tls"
	"fmt"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"net/http"
	"unsafe"
)

// #cgo CFLAGS: -I.
// #include "cow.h"
import "C"
import (
	"day04/ex02/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target ../../ex01v2 --name CandyServer --spec ../api/swagger.yaml --principal interface{}

const (
	CECost = 10
	AACost = 15
	NTCost = 17
	DECost = 21
	YRCost = 23
)

func configureFlags(api *operations.CandyServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CandyServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(func(params operations.BuyCandyParams) middleware.Responder {
		//бизнес логика ручки

		errorMessage := ""
		var candyCost int64

		if *params.Order.CandyType == "CE" {
			candyCost = CECost
		} else if *params.Order.CandyType == "AA" {
			candyCost = AACost
		} else if *params.Order.CandyType == "NT" {
			candyCost = NTCost
		} else if *params.Order.CandyType == "DE" {
			candyCost = DECost
		} else if *params.Order.CandyType == "YR" {
			candyCost = YRCost
		} else {
			errorMessage = "incorrect candy type"
		}

		if *params.Order.CandyCount < 0 {
			errorMessage = "invalid candy count"
		}

		if *params.Order.Money < 0 {
			errorMessage = "invalid amount of money value"
		}

		if errorMessage != "" {
			res := operations.NewBuyCandyBadRequest()
			resp := operations.BuyCandyBadRequestBody{errorMessage}
			res.SetPayload(&resp)
			return res
		}

		change := *params.Order.Money - *params.Order.CandyCount*candyCost
		if change < 0 {
			res := operations.NewBuyCandyPaymentRequired()
			resp := operations.BuyCandyPaymentRequiredBody{fmt.Sprintf("«Вам нужно еще %d денег!»", change*-1)}
			res.SetPayload(&resp)
			return res
		} else {
			res := operations.NewBuyCandyCreated()

			thanks := "Thank you!"
			cThanks := C.CString(thanks)
			defer C.free(unsafe.Pointer(cThanks))

			cow := C.ask_cow(cThanks)
			defer C.free(unsafe.Pointer(cow))

			goCow := C.GoString(cow)

			resp := operations.BuyCandyCreatedBody{Thanks: goCow, Change: change}
			res.SetPayload(&resp)
			return res

		}
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
