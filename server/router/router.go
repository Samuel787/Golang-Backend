package router

import (
	"../middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()

	// router.Use(middleware.AuthorizeUser)

	router.HandleFunc("/api/user", middleware.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/getUser/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deleteUser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/updateUser", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/addFollower", middleware.AddFollowerToUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/addFollowing", middleware.AddFollowingToUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/nearByFollowing", middleware.GetNearByFollowing).Methods("PUT", "OPTIONS")

	return router
}
