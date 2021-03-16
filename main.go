package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Test event",
		Description: "Just a test :P",
	},
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter data with the event and title description.")
	}

	json.Unmarshal(reqBody, &newEvent)

	if len(newEvent.ID) == 0 {
		var nextID int = 0
		for _, oneEvent := range events {
			i, err := strconv.Atoi(oneEvent.ID)
			if err != nil {
				fmt.Fprintf(w, "Invalid ID (must be integer).")
			} else {
				if i >= nextID {
					nextID = i
				}
			}
		}
		newEvent.ID = strconv.Itoa(nextID + 1)
	}
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, oneEvent := range events {
		if oneEvent.ID == eventID {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(oneEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid request - enter new title and description to update event.")
	}

	json.Unmarshal(reqBody, &updatedEvent)

	for i, oneEvent := range events {
		if oneEvent.ID == eventID {
			if len(updatedEvent.Title) > 0 {
				oneEvent.Title = updatedEvent.Title
			}
			if len(updatedEvent.Description) > 0 {
				oneEvent.Description = updatedEvent.Description
			}
			events[i] = oneEvent
			json.NewEncoder(w).Encode(oneEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	for i, oneEvent := range events {
		if eventID == oneEvent.ID {
			events = append(events[:i], events[i+1:]...)
			json.NewEncoder(w).Encode(oneEvent)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/create", createEvent).Methods("POST")
	router.HandleFunc("/events/update/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/delete/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
