package handlers

import (
	"github.com/gorilla/mux"
)

func MakeHTTPHandlers(router mux.Router)  {
	// User handlers. TODO: change nil to Methods.
	router.HandleFunc("/user/register", nil)
	router.HandleFunc("/user/login", nil)
	router.HandleFunc("/user/", nil) // PUT
	router.HandleFunc("/user/", nil) // DELETE
	router.HandleFunc("/user/confirm/", nil) // GET
}
