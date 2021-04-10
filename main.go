package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var templates *template.Template

func loadTemplates() {
	templates = template.Must(template.ParseFiles("index.html"))
}

//serving the root page (every request on "/" that is not handled specifically)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	loadTemplates()
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		log.Fatal("Error executing the index template: " + err.Error())
	}
}

var wg sync.WaitGroup //to handle waiting on all goroutines

var serverHandler *http.ServeMux //the handler that handles requests
var server http.Server           //the server itself which uses the handler serverHandler

func main() {
	//load the templates (in production only here)
	loadTemplates()

	//create the handler and server
	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":8080", Handler: serverHandler}

	staticHandler := http.FileServer(http.Dir("static")) //static handler to serve .js and .css files from the static directory
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	//setup the handler functions of the serverHandler
	serverHandler.HandleFunc("/", indexHandler)

	//start the goroutine that handles some commands mainly for debugging but also to shutdown the server
	log.Println("Starting cmd goroutine")
	go func() {
		wg.Add(1)
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
			}
		}
	}
}
