package router

import (
	"fmt"

	"net/http"
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

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        // log.Println(r.RequestURI)
		fmt.Println("This is the auth middleware")
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}