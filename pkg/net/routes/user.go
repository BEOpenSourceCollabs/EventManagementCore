package routes

import (
	"encoding/json"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
)

type userRoutes struct {
	userRepository repository.UserRepository
}

func NewUserRoutes(router net.AppRouter, userRepository repository.UserRepository) userRoutes {
	routes := userRoutes{
		/* inject dependencies */
		userRepository: userRepository,
	}

	// mount routes to router.
	router.Post("/api/users/create", http.HandlerFunc(routes.PostCreateUser))
	router.Get("/api/users/{id}", http.HandlerFunc(routes.GetUserById))
	router.Put("/api/users/update/{id}", http.HandlerFunc(routes.PutUpdateUserById))
	router.Delete("/api/users/delete/{id}", http.HandlerFunc(routes.DeleteUserById))

	return routes
}

func (u userRoutes) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	payload := &models.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := u.userRepository.CreateUser(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": payload,
	})
}

func (u userRoutes) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := u.userRepository.GetUserByID(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

func (u userRoutes) PutUpdateUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Load the user first
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Merge in the request payload
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// submit the changes
	if err := u.userRepository.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

func (u userRoutes) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := u.userRepository.DeleteUser(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
