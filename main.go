package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var wg sync.WaitGroup //to handle waiting on all goroutines

var serverHandler *http.ServeMux //the handler that handles requests
var server http.Server           //the server itself which uses the handler serverHandler

var templates *template.Template //html templates for all files that need to be served

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
	serverHandler.HandleFunc("/api/editFood", EditFoodHandler)
	serverHandler.HandleFunc("/api/getFood", GetFoodHandler)
	serverHandler.HandleFunc("/api/changeFood", ChangeFoodHandler)

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
			default:
				fmt.Println("cmd not supported")
			}
		}
	}
}
