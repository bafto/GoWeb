package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var randSource rand.Source = rand.NewSource(time.Now().UnixNano())
var random *rand.Rand = rand.New(randSource)

func ContainsFood(sli []Food, comp Food) bool {
	for _, v := range sli {
		if v.Equals(comp) {
			return true
		}
	}
	return false
}

//send a response with a json body constructed from data over w
func returnJsonFromStruct(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

//handles POST and DELETE methods on /api/editFoods to work on the foods.json file
func EditFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" { //validate that the request contains json
		log.Println("Request does not contain json")
		http.Error(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	//get the food from the request body
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errLog.Println("Error unmarshaling requests json body: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPost: //add food to the list
		food.ID = GenerateFoodID() //generate the food ID
		if food.Label == nil {
			food.Label = make([]string, 0)
		}
		if err != nil {
			log.Println("Error retreiving label list: " + err.Error())
			errLog.Println("Error retreiving label list: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = AddFoodToList(food) //and add the food to the list
		if err != nil {
			log.Println("Error editing Food List: " + err.Error())
			errLog.Println("Error editing Food List: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJsonFromStruct(w, food, http.StatusOK)
	case http.MethodDelete: //delete food from the list
		err = DeleteFoodFromList(food)
		if err != nil {
			log.Println("Error deleting from Food List: " + err.Error())
			errLog.Println("Error deleting from Food List: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("success"))
	default:
		log.Println("Method type \"" + r.Method + "\" not handled")
		http.Error(w, "false Method type", http.StatusBadRequest)
	}
}

//handles POST requests on /api/changeFood to save changes on a single food like adding Labels
func ChangeFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" { //check if the request contains json
		log.Println("Request does not contain json")
		http.Error(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food) //retreive the food
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errLog.Println("Error unmarshaling requests json body: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	changed, err := food.ValidateLabel()
	if err != nil {
		log.Println("Error validating food Label: " + err.Error())
		errLog.Println("Error validating food Label: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if changed {
		returnJsonFromStruct(w, food, http.StatusNotModified)
		return
	}
	switch r.Method {
	case http.MethodPatch:
		err = ChangeFoodInList(food)
		if err != nil {
			log.Println("Error changing food list: " + err.Error())
			errLog.Println("Error changing food list: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		returnJsonFromStruct(w, food, http.StatusOK)
	default:
		log.Println("Method type \"" + r.Method + "\" not handled")
		http.Error(w, "false Method type", http.StatusBadRequest)
	}
}

//responds with a json array containing foods.json
func GetFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("received request on /api/getFood of type " + r.Method + " that should've been GET")
		http.Error(w, "Method should have been GET", http.StatusMethodNotAllowed)
		return
	}
	Foods, err := GetWholeFoodList()
	if err != nil {
		log.Println("Error retreiving whole food list: " + err.Error())
		errLog.Println("Error retreiving whole food list: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	returnJsonFromStruct(w, Foods, http.StatusOK)
}

func GetFoodConstrainedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("received request on /api/getFood of type " + r.Method + " that should've been POST")
		http.Error(w, "Method should have been POST", http.StatusMethodNotAllowed)
		return
	}

	type Data struct {
		Label []string
		Names []string
		Count int
	}

	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errLog.Println("Error unmarshaling requests json body: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Foods, err := GetEveryFoodWithLabel(data.Label)
	if err != nil {
		log.Println("Error getting foods with label: " + err.Error())
		errLog.Println("Error getting foods with label: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var retFoods []Food = make([]Food, 0)
	if data.Count <= len(Foods) {
		for i := 0; i < data.Count; i++ {
			food := Foods[random.Intn(len(Foods))]
			if !ContainsFood(retFoods, food) && !ContainsString(data.Names, food.Name) {
				retFoods = append(retFoods, food)
			} else {
				i--
				continue
			}
		}
	} else {
		retFoods = Foods
	}
	returnJsonFromStruct(w, retFoods, http.StatusOK)
}

func GetLabelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("received request on /api/getFood of type " + r.Method + " that should've been GET")
		http.Error(w, "Method should have been GET", http.StatusMethodNotAllowed)
		return
	}
	label, err := GetWholeLabelList()
	if err != nil {
		log.Println("error receiving label list: " + err.Error())
		errLog.Println("error receiving label list: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	returnJsonFromStruct(w, label, http.StatusOK)
}

func EditLabelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" { //validate that the request contains json
		log.Println("Request does not contain json")
		http.Error(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	//get the label from the request body
	var label string
	err := json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		log.Println("Error unmarshaling requests json body: " + err.Error())
		errLog.Println("Error unmarshaling requests json body: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPost: //add label to the list
		err = AddLabelToList(label)
		if err != nil {
			log.Println("Error adding label to list: " + err.Error())
			errLog.Println("Error adding label to list: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(label))
	case http.MethodDelete: //delete label from the list
		err = DeleteLabelFromList(label)
		if err != nil {
			log.Println("Error deleting from labelList: " + err.Error())
			errLog.Println("Error deleting from labelList: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(label))
	default:
		log.Println("Method type \"" + r.Method + "\" not handled")
		http.Error(w, "false Method type", http.StatusBadRequest)
	}
}
