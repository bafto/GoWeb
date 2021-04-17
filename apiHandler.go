package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//write an json error to w
func errorJson(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(`{"error":"` + msg + `"}`)
}

//send a response with json body over w
func returnJson(w http.ResponseWriter, JSON string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JSON)
}

//send a response with a json body constructed from data over w
func returnJsonFromStruct(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

//handles POST and DELETE methods on /api/editFoods to work on the foods.json file
func EditFoodHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /api/editFood")
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Request does not contain json")
		errorJson(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPost:
		log.Println("request on /api/editFood was of type POST")
		food = *NewFood(food)
		if food.Label == nil {
			food.Label = make(map[string]bool)
		}
		err = AddFoodToList(food)
		if err != nil {
			log.Println("Error editing Food List: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJsonFromStruct(w, food, http.StatusOK)
	case http.MethodDelete:
		log.Println("request on /api/editFood was of type DELETE")
		err = DeleteFoodFromList(food)
		if err != nil {
			log.Println("Error deleting from Food List: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJson(w, `{"success":true}`, http.StatusOK)
	default:
		log.Println("Method type \"" + r.Method + "\" not handled")
		errorJson(w, "false Method type", http.StatusBadRequest)
	}
}

//handles POST requests on /api/changeFood to save changes on a single food like adding Labels
func ChangeFoodHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /api/changeFood")
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Request does not contain json")
		errorJson(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPost:
		log.Println("request on /api/changeFood was of type POST")
		err = ChangeFoodInList(food)
		if err != nil {
			log.Println("Error changing food list: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJsonFromStruct(w, food, http.StatusOK)
	default:
		log.Println("Method type \"" + r.Method + "\" not handled")
		errorJson(w, "false Method type", http.StatusBadRequest)
	}
}

//responds with a json array containing foods.json
func GetFoodHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /api/getFood")
	if r.Method != http.MethodGet {
		log.Println("received request on /api/getFood of type " + r.Method + " that should've been GET")
		errorJson(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}
	if b, err := CheckFileExist("static/foods.json"); err != nil {
		log.Println("File static/foods.json may or may not exist: " + err.Error())
	} else if !*b {
		os.Create("static/foods.json")
		ioutil.WriteFile("static/foods.json", []byte("[]"), 0644)
	}
	file, err := ioutil.ReadFile("static/foods.json")
	if err != nil {
		log.Println("Error reading static/foods.json: " + err.Error())
		errorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var Foods []Food
	err = json.Unmarshal(file, &Foods)
	if err != nil {
		log.Println("Error unmarshaling static/foods.json: " + err.Error())
		errorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnJsonFromStruct(w, Foods, http.StatusOK)
}
