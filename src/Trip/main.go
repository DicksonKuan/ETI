package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Trip struct {
	TripID          int    `json: "TripID"`
	CustomerID      int    `json: "CustomerID"`
	DriverID        int    `json: "DriverID"`
	PickUpLocation  string `json: "PickUpLocation"`
	DropOffLocation string `json: "DropOffLocation"`
	PickUpTime      string `json: "PickUpTime"`
	DropOffTime     string `json: "DropOffTime"`
	Status          string `json: "Status"`
}

//Database function
func CheckTrip(db *sql.DB, CustomerID int, DriverID int) bool {
	//To check if there is any uncompleted trips
	query := fmt.Sprintf("Select * FROM Trip WHERE CustomerID = %d OR DriverID = %d", CustomerID, DriverID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var currentTrip Trip
	for results.Next() {
		// map this type to the record in the table and update the object with new data
		err = results.Scan(&currentTrip.TripID, &currentTrip.CustomerID, &currentTrip.DriverID, &currentTrip.PickUpLocation, &currentTrip.DropOffLocation, &currentTrip.PickUpTime, &currentTrip.DropOffTime, &currentTrip.Status)
		if err != nil {
			panic(err.Error())
		} else if currentTrip.Status != "Completed" { //To check if trip have finish
			return false
		}
	}
	return true
}

func SearchAvailDriver(db *sql.DB) int {
	query := "SELECT driver.DriverID FROM driver LEFT JOIN trip ON driver.DriverID = trip.DriverID WHERE trip.Status = 'Completed' OR trip.DriverID IS NULL"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var DriverID int
	for results.Next() {
		err = results.Scan(&DriverID)
		if err != nil {
			panic(err.Error())
		}
	}
	return DriverID
}

func GetTrip(db *sql.DB, ID int) Trip {
	query := fmt.Sprintf("Select * FROM Trip WHERE TripID = %d", ID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&trip.TripID, &trip.DriverID, &trip.CustomerID, &trip.PickUpLocation, &trip.DropOffLocation, &trip.PickUpTime, &trip.DropOffTime, &trip.Status)
		if err != nil {
			panic(err.Error())
		}
	}
	return trip
}

func GetTrips(db *sql.DB, CustomerID int) []Trip {
	query := fmt.Sprintf("Select * FROM Trip WHERE CustomerID = %d", CustomerID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trips []Trip
	for results.Next() {
		var trip Trip
		// map this type to the record in the table
		err = results.Scan(&trip.TripID, &trip.DriverID, &trip.CustomerID, &trip.PickUpLocation, &trip.DropOffLocation, &trip.PickUpTime, &trip.DropOffTime, &trip.Status)
		if err != nil {
			panic(err.Error())
		}
		trips = append(trips, trip)
	}
	return trips
}

func CreateTrip(db *sql.DB, trip Trip) bool {
	if !CheckTrip(db, trip.CustomerID, trip.DriverID) {
		return false
	}
	query := fmt.Sprintf(
		"INSERT INTO Trip (DriverID, CustomerID, PickupLocation, DropoffLocation, PickUpTime, DropOffTime, Status) VALUES(%d, %d, '%s', '%s', '%s', ' ', 'Pending')",
		trip.DriverID,
		trip.CustomerID,
		trip.PickUpLocation,
		trip.DropOffLocation,
		trip.PickUpTime)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditTrip(db *sql.DB, trip Trip) bool {
	if trip.TripID == 0 {
		return false
	}
	query := fmt.Sprintf("UPDATE Trip SET DriverID = %d, CustomerID = %d, PickUpLocation = '%s', DropOffLocation='%s', PickUpTime = '%s', DropOffTime = '%s', Status = '%s' WHERE TripID = %d",
		trip.DriverID,
		trip.CustomerID,
		trip.PickUpLocation,
		trip.DropOffLocation,
		trip.PickUpTime,
		trip.DropOffTime,
		trip.Status,
		trip.TripID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

//API
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RideShare Trip API")
}
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func APIRouter(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	if r.Method == "DELETE" {
		println("You can't delete Trip records")
	}
	if r.Header.Get("Content-type") == "application/json" {
		//Database
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
		if err != nil {
			fmt.Println(err)
		}
		if r.Method == "GET" { //GET trip
			var TripInformation Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &TripInformation)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("There was an error encoding the json. err = %s", err)
				} else if TripInformation.CustomerID == 0 && TripInformation.DriverID == 0 { // Check if data is empty
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"There is no trip session"))
					return
				} else {
					json.NewEncoder(w).Encode(GetTrip(db, TripInformation.TripID))
					w.WriteHeader(http.StatusAccepted)
					return
				}
			}
		} else if r.Method == "POST" {
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("There was an error encoding the json. err = %s", err)
				} else {
					if newTrip.CustomerID == 0 || newTrip.DriverID == 0 {
						w.WriteHeader(
							http.StatusUnprocessableEntity)
						w.Write([]byte(
							"422 - Please supply in JSON format"))
						return
					} else {
						if CheckTrip(db, newTrip.CustomerID, newTrip.DriverID) {
							CreateTrip(db, newTrip)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Trip created successfully"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("The current driver or customer has an ongoing trip"))
							return
						}
					}
				}
			}
		} else if r.Method == "PUT" {
			var updatedTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &updatedTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("There was an error encoding the json. err = %s", err)
				}
				if updatedTrip.CustomerID == 0 || updatedTrip.DriverID == 0 || updatedTrip.TripID == 0 {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply customer information"))
					return
				} else {
					if CheckTrip(db, updatedTrip.CustomerID, updatedTrip.DriverID) { //To check with the database if there is any record
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no trip records found"))
					} else {
						//To update user details
						if EditTrip(db, updatedTrip) {
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Trip updated successfully"))
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("Trip unable to update"))
						}
					}
				}
			}
		}
	}
}

func GetAllTrips(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	//Database
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)

	if params["CustomerID"] == "" { // Check if data is empty
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte(
			"There is no trip session"))
		return
	} else {
		CustomerID, err := strconv.Atoi(params["CustomerID"])
		if err != nil {
			fmt.Println(err)
		}
		println(CustomerID)
		json.NewEncoder(w).Encode(GetTrips(db, CustomerID))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func main() {
	// //Database
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var trip Trip
	// //Dummy data
	// trip.TripID = 3
	// trip.DriverID = 3
	// trip.CustomerID = 1
	// trip.PickUpLocation = "123456"
	// trip.DropOffLocation = "123457"
	// trip.PickUpTime = "14:00"
	// trip.DropOffTime = "15:00"
	// trip.Status = "Completed"
	//println(GetTrip(db, trip.TripID).Status) - working
	//println(CreateTrip(db, trip)) - working
	//println(EditTrip(db, trip)) - working

	//API
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/Trip", home)                                           //Test API
	router.HandleFunc("/api/v1/Trip/Router", APIRouter).Methods("GET", "PUT", "POST") //API Manipulation
	router.HandleFunc("/api/v1/Trip/{CustomerID}", GetAllTrips)

	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
