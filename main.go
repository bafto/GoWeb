package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup //to handle waiting on all goroutines

var serverHandler *http.ServeMux //the handler that handles requests
var server http.Server           //the server itself which uses the handler serverHandler

var templates *template.Template //html templates for all files that need to be served

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

func loadTemplates() {
	templates = template.Must(template.ParseFiles("index.html"))
}

//serving the root page (every request on "/" that is not handled specifically)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /")
	loadTemplates()
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		log.Fatal("Error executing the index template: " + err.Error())
	}
}

//handles POST and DELETE methods on /api/editFoods to work on the foods.json file
func editFoodHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /api/editFood")
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Request does not contain json")
		errorJson(w, "Request Header is not application/json", http.StatusBadRequest)
		return
	}
	var food Food
	err := json.NewDecoder(r.Body).Decode(&food)
	log.Println(food.ID)
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
			food.Label = make([]string, 0)
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
	}
	log.Println(food)
}

//responds with a json array containing foods.json
func getFoodHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	//load the templates (in production only here)
	loadTemplates()

	//create the handler and server
	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":8080", Handler: serverHandler}

	staticHandler := http.FileServer(http.Dir("static")) //static handler to serve files from the static directory
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	//setup the handler functions of the serverHandler
	serverHandler.HandleFunc("/", indexHandler)
	serverHandler.HandleFunc("/api/editFood", editFoodHandler)
	serverHandler.HandleFunc("/api/getFood", getFoodHandler)

	//start the goroutine that handles some commands mainly for debugging but also to shutdown the server
	log.Println("Starting cmd goroutine")
	wg.Add(1)
	go func() {
		defer wg.Done() //tell the waiter group that we are finished at the end
		cmdInterface()
		log.Println("cmd goroutine finished")
	}()

	//startup the server to listen for requests
	fmt.Println("server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	} else if err == http.ErrServerClosed {
		log.Println("Server not listening anymore")
	}

	//wait for the cmd goroutine to finish (maybe more routines later)
	wg.Wait()
}

//the cmd interface should only be used during development and not while stuff is getting loged to the console, as the terminal can get quite messy pretty quickly otherwise
func cmdInterface() {
	for loop := true; loop; {
		var inp string
		_, err := fmt.Scanln(&inp)
		if err != nil {
			log.Println(err.Error())
		} else {
			switch inp {
			case "quit":
				log.Println("Attempting to shutdown server")
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal("Error while trying to shutdown server: " + err.Error())
				}
				log.Println("Server was shutdown")
				loop = false
			case "update-templates":
				loadTemplates() //only use to debug, as changing the templates variable is not properly synchronized between goroutines
			case "add":
				fmt.Print("Food Name to add: ")
				var inp2 string
				_, err := fmt.Scanln(&inp2)
				if err != nil {
					log.Println(err.Error())
				} else {
					err = AddFoodToList(Food{Name: inp2})
					if err != nil {
						log.Println("Error editing Food list: " + err.Error())
					}
				}
			}
		}
	}
}
