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

	return router
}

func CorsHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
