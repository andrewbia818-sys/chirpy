//package main

//import (
//	"log"
//	"net/http"
//)

//func main() {
//	const filepathroot = "."//	const port = "8080"

// Create a new http.ServeMux to route requests
//	sm := http.NewServeMux()
// Serve the root directory which defaults to index.html
//	sm.Handle("/", http.FileServer(http.Dir("filepathroot")))

// Create a new http.Server struct.
//	s := &http.Server{
//		Addr:    ":" + port,
//		Handler: sm,
//	}

// Use the ListenAndServe method to start the server
//	if err := s.ListenAndServe(); err != nil {
//		log.Fatalf("Could not start server: %s\n", err.Error())
//	}

//	log.Printf("Serving files from %s on port: %s\n", filepathroot, port)
//	log.Fatal(s.ListenAndServe())

//}

package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	sm := http.NewServeMux()

	// Serve the assets directory at /assets/
	sm.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// Serve the root directory which defaults to index.html
	// sm.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	log.Printf("Serving on port %s\n", port)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
