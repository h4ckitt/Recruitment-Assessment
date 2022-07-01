package router

import (
	"assessment/interface/mux/controller"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRouter(controller *controller.Controller) *mux.Router {
	router := mux.NewRouter()

	pathRouter := router.PathPrefix("/numbersvc").Subrouter()

	pathRouter.HandleFunc("", controller.FetchAllPhoneNumbers)
	pathRouter.HandleFunc("", controller.FilterPhoneNumbersByCountryAndState).Queries("country", "{country}", "state", "{state}", "limit", "{limit}", "page", "{page}").Methods(http.MethodGet)
	pathRouter.HandleFunc("", controller.FetchAllPhoneNumbers).Queries("limit", "{limit}", "page", "{page}").Methods(http.MethodGet)

	return router
}
