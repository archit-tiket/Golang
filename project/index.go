package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "12345678"
	DB_NAME     = "buses"
)

type BUS struct {
	BusID   string `json:"busid"`
	BusName string `json:"busname"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []BUS  `json:"data"`
	Message string `json:"message"`
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db //was written DB
}

func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all buses
	router.HandleFunc("/buses/", GetBuses).Methods("GET")

	// Create a movie
	router.HandleFunc("/buses/", CreateBus).Methods("POST")

	// Delete a specific movie by the movieID
	router.HandleFunc("/buses/{busid}", DeleteBus).Methods("DELETE")

	// Delete all buses
	router.HandleFunc("/buses/", Deletebuses).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all buses

// response and request handlers
func GetBuses(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting buses...")

	
	rows, err := db.Query("SELECT * FROM buses")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var buses []BUS

	// Foreach bus
	for rows.Next() {
		var id int
		var busID string
		var busName string

		err = rows.Scan(&id, &busID, &busName)

		// check errors
		checkErr(err)

		buses = append(buses, BUS{BusID: busID, BusName: busName})
	}

	var response = JsonResponse{Type: "success", Data: buses}

	json.NewEncoder(w).Encode(response)
}

// Create a bus

// response and request handlers
func CreateBus(w http.ResponseWriter, r *http.Request) {
	busID := r.FormValue("busid")
	busName := r.FormValue("busname")

	var response = JsonResponse{}

	if busID == "" || busName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting bus into DB")

		fmt.Println("Inserting new bus with ID: " + busID + " and name: " + busName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO buses(busID, busName) VALUES($1, $2) returning id;", busID, busName).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The bus has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a bus

// response and request handlers
func DeleteBus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	movieID := params["busid"]

	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing busID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting bus from DB")

		_, err := db.Exec("DELETE FROM buses where busID = $1", movieID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The bus has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all buses

// response and request handlers
func Deletebuses(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all buses...")

	_, err := db.Exec("DELETE FROM buses")

	// check errors
	checkErr(err)

	printMessage("All buses have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All buses have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
