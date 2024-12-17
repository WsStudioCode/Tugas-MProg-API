package routes

import (
	"net/http"
	"pencatatan-keuangan-api/controllers"
)

func RegisterAuthRoutes() {
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/register", controllers.Register)
}
