package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	if r.Header.Get("Content-Type") != "application/json" { //validate that the request contains json
		log.Println("Request does not contain json")
		errorJson(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	//get the food from the request body
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPost: //add food to the list
		log.Println("request on /api/editFood was of type POST")
		food.ID = GenerateFoodID() //generate the food ID
		if food.Label == nil {
			food.Label = make(map[string]bool)
		}
		LabelList, err := GetWholeLabelList()
		if err != nil {
			log.Println("Error retreiving label list: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
		}
		for _, v := range LabelList { //add all labels to the food
			food.Label[v] = false
		}
		err = AddFoodToList(food) //and add the food to the list
		if err != nil {
			log.Println("Error editing Food List: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJsonFromStruct(w, food, http.StatusOK)
	case http.MethodDelete: //delete food from the list
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
	if r.Header.Get("Content-Type") != "application/json" { //check if the request contains json
		log.Println("Request does not contain json")
		errorJson(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food) //retreive the food
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
	Foods, err := GetWholeFoodList()
	if err != nil {
		log.Println("Error retreiving whole food list: " + err.Error())
		errorJson(w, err.Error(), http.StatusInternalServerError)
	}
	returnJsonFromStruct(w, Foods, http.StatusOK)
}

func GetLabelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /api/getLabel")
	if r.Method != http.MethodGet {
		log.Println("received request on /api/getFood of type " + r.Method + " that should've been GET")
		errorJson(w, `{"success":false}`, http.StatusMethodNotAllowed)
		return
	}
	label, err := GetWholeLabelList()
	if err != nil {
		log.Println("error receiving label list: " + err.Error())
		errorJson(w, err.Error(), http.StatusInternalServerError)
	}
	returnJsonFromStruct(w, label, http.StatusOK)
}

func EditLabelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /api/editLabel")
	if r.Header.Get("Content-Type") != "application/json" { //validate that the request contains json
		log.Println("Request does not contain json")
		errorJson(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	//get the food from the request body
	var label string
	err := json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPost: //add food to the list
		log.Println("request on /api/editLabel was of type POST")
		err = AddLabelToList(label)
		if err != nil {
			log.Println("Error adding label to list: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
		}
		returnJson(w, label, http.StatusOK)
	case http.MethodDelete: //delete food from the list
		log.Println("request on /api/editLabel was of type DELETE")
		err = DeleteLabelFromList(label)
		if err != nil {
			log.Println("Error deleting from labelList: " + err.Error())
			errorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJson(w, label, http.StatusOK)
	default:
		log.Println("Method type \"" + r.Method + "\" not handled")
		errorJson(w, "false Method type", http.StatusBadRequest)
	}
}
