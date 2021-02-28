package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct {
	dispatcher *mux.Router
}

func NewMuxRouter() Router {
	muxDispatcher := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    muxDispatcher.PathPrefix("/static/").Handler(s)
	return &muxRouter{dispatcher: muxDispatcher}
}

func (router *muxRouter) GET(uri string, f func(response http.ResponseWriter, request *http.Request)) {
	router.dispatcher.HandleFunc(uri, f).Methods("GET")
}

func (router *muxRouter) POST(uri string, f func(response http.ResponseWriter, request *http.Request)) {
	router.dispatcher.HandleFunc(uri, f).Methods("POST")
}

func (router *muxRouter) SERVE(port string) {
	http.ListenAndServe(port, router.dispatcher)
}
