package main

import (
	"context"
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

var errLog *log.Logger

func loadTemplates() {
	templates = template.Must(template.ParseFiles("html/index.html", "html/labelList.html", "html/foodPlanner.html"))
}

//serving the root page (every request on "/" that is not handled specifically)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /")
	loadTemplates()
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error executing the index template: " + err.Error())
		errLog.Println("Error executing the index template: " + err.Error())
	}
}

//serving the labelList page
func labelListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /labelList")
	loadTemplates()
	err := templates.ExecuteTemplate(w, "labelList.html", nil)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error executing the labelList template: " + err.Error())
		errLog.Println("Error executing the labelList template: " + err.Error())
	}
}

//serving the foodPlanner page
func foodPlannerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a request on /foodPlanner")
	loadTemplates()
	err := templates.ExecuteTemplate(w, "foodPlanner.html", nil)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error executing the foodPlanner template: " + err.Error())
		errLog.Println("Error executing the foodPlanner template: " + err.Error())
	}
}

//validating foods.json and creating it if necessary
func setupFoodFile() {
	if _, err := os.Stat("resources/foods.json"); err == nil {
		_, err = GetWholeFile()
		if err != nil {
			log.Fatal("foods.json is in a bad state, check the json format: " + err.Error())
		}
	} else if os.IsNotExist(err) {
		log.Println(err)
		err = os.Mkdir("resources", 0755)
		if err != nil {
			log.Fatal("failed to create the resources directory")
		}
		os.Create("resources/foods.json")
		ioutil.WriteFile("resources/foods.json", []byte(`{"LabelList":[],"IngredientList":[],"FoodLIst":[]}`), 0644)
		_, err = GetWholeFile()
		if err != nil {
			log.Fatal("foods.json is in a bad state even after creating: " + err.Error())
		}
	} else {
		log.Fatal("foods.json is in an unknown state: " + err.Error())
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	errLog = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	errLog.Println("Setupt log for this session")

	//load the templates (in production only here)
	loadTemplates()

	setupFoodFile()

	//create the handler and server
	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":8080", Handler: serverHandler}

	staticHandler := http.FileServer(http.Dir("static")) //static handler to serve files from the static directory
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	//setup the html page handler functions
	serverHandler.HandleFunc("/", indexHandler)
	serverHandler.HandleFunc("/labelList", labelListHandler)
	serverHandler.HandleFunc("/foodPlanner", foodPlannerHandler)
	//setup the handler functions defined in apiHandler.go
	serverHandler.HandleFunc("/api/editFood", EditFoodHandler)
	serverHandler.HandleFunc("/api/getFood", GetFoodHandler)
	serverHandler.HandleFunc("/api/changeFood", ChangeFoodHandler)
	serverHandler.HandleFunc("/api/getLabel", GetLabelHandler)
	serverHandler.HandleFunc("/api/editLabel", EditLabelHandler)
	serverHandler.HandleFunc("/api/getFoodConstrained", GetFoodConstrainedHandler)

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
		errLog.Println(err.Error())
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
			errLog.Println(err.Error())
		} else {
			switch inp {
			case "quit":
				log.Println("Attempting to shutdown server")
				err := server.Shutdown(context.Background())
				if err != nil {
					errLog.Println("Error while trying to shutdown server: " + err.Error())
					log.Fatal("Error while trying to shutdown server: " + err.Error())
				}
				log.Println("Server was shutdown")
				loop = false
			default:
				fmt.Println("cmd not supported")
			}
		}
	}
}
