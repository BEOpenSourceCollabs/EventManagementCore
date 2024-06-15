package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
)

type eventRoutes struct {
	eventRepository repository.EventRepository
}

func NewEventRoutes(router net.AppRouter, eventRepository repository.EventRepository) eventRoutes {
	routes := eventRoutes{
		eventRepository: eventRepository,
	}

	// mount routes
	router.Post("/api/events/create", http.HandlerFunc(routes.PostCreateEvent))
	router.Get("/api/events/{id}", http.HandlerFunc(routes.GetEventById))
	router.Get("/api/organizers/{org_id}/events", http.HandlerFunc(routes.GetEventsByOrganizerId))
	router.Put("/api/events/update/{id}", http.HandlerFunc(routes.PutUpdateEventById))
	router.Delete("/api/events/delete/{id}", http.HandlerFunc(routes.DeleteEventById))

	return routes
}

func (e eventRoutes) PostCreateEvent(w http.ResponseWriter, r *http.Request) {
	payload := &models.EventModel{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := e.eventRepository.CreateEvent(payload); err != nil {
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

func (e eventRoutes) GetEventById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	event, err := e.eventRepository.GetEventByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"event": event,
	})
}

func (e eventRoutes) GetEventsByOrganizerId(w http.ResponseWriter, r *http.Request) {
	org_id := r.PathValue("org_id")

	events, err := e.eventRepository.GetEventsByOrganizer(org_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
	})
}

func (e eventRoutes) PutUpdateEventById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// get the event
	event, err := e.eventRepository.GetEventByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Merge in the request payload
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// update the event
	if err := e.eventRepository.UpdateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"event": event,
	})
}

func (e eventRoutes) DeleteEventById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := e.eventRepository.DeleteEvent(id); err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
