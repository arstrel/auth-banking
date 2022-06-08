package app

import (
	"github.com/arstrel/rest-banking/auth/domain"
	"github.com/gorilla/mux"
)

func sanityCheck() {

}

func Start() {
	sanityCheck()
	router := mux.NewRouter()

	authRepository := domain.NewAuthRepository()
}
