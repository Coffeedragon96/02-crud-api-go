package main

// Import required libraries
import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	// Using gorilla mux for routing
	"github.com/gorilla/mux"
)

// Creating a simple movie object for this app's purpose
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Creating a nested object called Director for saving further information
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Since we are not using any database, creating a simple slice to store all the values
var movies []Movie

// Function to get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {

	// Setting header for response
	w.Header().Set("Content-Type", "application/json")

	// Encoding all values present in slice movies
	json.NewEncoder(w).Encode(movies)
}

// Function to delete movies
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieving all the parameters coming in the request
	params := mux.Vars(r)

	// Iterating over movies slice
	for index, item := range movies {

		// Matching incoming id with existing ids in slice
		if item.ID == params["id"] {

			// Dropping out the matching id in the slice and overwriting a new one
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
}

// Function to get individual movie based on id
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {

		// Returning movie with matching id
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// Function to create movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// creating new instance of the struct object
	var movie Movie

	// Decoding incoming request
	_ = json.NewDecoder(r.Body).Decode(&movie)

	// Generating a new random movie id
	movie.ID = strconv.Itoa(rand.Intn(1000000))

	// Appending new generated movie object in the slice
	movies = append(movies, movie)

	// Returning the newly added object instance
	json.NewEncoder(w).Encode(movie)
}

// Function to update existing movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Delete the existing details of matching id
	for index, items := range movies {
		if items.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)

			// Create new object with same id and updated details and append it to the slice
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}
}

// Main function
func main() {

	// Initialising router instance
	r := mux.NewRouter()

	// Creating temporary variables in slice
	movies = append(movies, Movie{ID: "1", Isbn: "123445", Title: "Movie 1", Director: &Director{Firstname: "Steven", Lastname: "Spielberg"}})
	movies = append(movies, Movie{ID: "2", Isbn: "987655", Title: "Movie 2", Director: &Director{Firstname: "Anurag", Lastname: "Kashyap"}})

	// Creating routes | CRUD Operations
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Declaring server running on port
	fmt.Printf("Starting server at port 8000\n")

	// Initialising router on port 8000 using gorilla mux router and logging in case of error
	log.Fatal(http.ListenAndServe(":8000", r))
}
