package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Trip struct {
	TripID          int    `json: "TripID"`
	CustomerID      int    `json: "CustomerID"`
	DriverID        string `json: "DriverID"`
	PickUpLocation  string `json: "PickUpLocation"`
	DropOffLocation string `json: "DropOffLocation"`
	PickUpTime      string `json: "PickUpTime"`
	DropOffTime     string `json: "DropOffTime"`
	Status          string `json: "Status"`
}

type Data struct { //Structure to get neccasry input from front end
	CustomerEmail   string `json: "CustomerEmail"`
	PickUpLocation  string `json: "PickUpLocation"`
	DropOffLocation string `json: "DropOffLocation"`
	PickUpTime      string `json: "PickUpTime"`
	DropOffTime     string `json: "DropOffTime"`
}

type AllTrip struct {
	DriverCarPlate  string `json: "CarPlate"`
	PickUpLocation  string `json: "PickUpLocation"`
	DropOffLocation string `json: "DropOffLocation"`
	PickUpTime      string `json: "PickUpTime"`
	DropOffTime     string `json: "DropOffTime"`
	Status          string `json: "Status"`
}

//Database function
func CheckTrip(db *sql.DB, CustomerID int, DriverID string) bool {
	//To check if there is any uncompleted trips
	query := fmt.Sprintf("Select Status FROM Trip WHERE CustomerID = %d OR DriverID = '%s'", CustomerID, DriverID)
	results, err := db.Query(query)
	if err != nil {
		panic("Error here" + err.Error())
	}
	var currentTrip Trip
	for results.Next() {
		// map this type to the record in the table and update the object with new data
		err = results.Scan(&currentTrip.Status)
		if err != nil {
			panic(err.Error())
		} else if currentTrip.Status != "Completed" { //To check if trip have finish
			return false
		}
	}
	return true
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func CheckDriverAvail(db *sql.DB, IDs []string) string {
	//Check if driver is currently on a trip
	//QueryString := fmt.Sprintf("Select * FROM Trip WHERE Status != 'Completed' AND DriverID = '%s' ", IDs)
	QueryString := "Select DriverID FROM Trip WHERE Status != 'Completed' "

	for i := range IDs {
		if i == 0 {
			QueryString += fmt.Sprintf("AND DriverID = '%s' ", string(IDs[i]))
		} else {
			QueryString += fmt.Sprintf("OR DriverID = '%s' ", string(IDs[i]))
		}

	}
	results, err := db.Query(QueryString)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var DriverID string
		// map this type to the record in the table
		err = results.Scan(&DriverID)
		if err != nil {
			panic(err.Error())
		}
		IDs = remove(IDs, DriverID)
	}
	return IDs[0]
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
	query := fmt.Sprintf("Select * FROM Trip WHERE CustomerID = %d ORDER BY TripID DESC", CustomerID)
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
	query := fmt.Sprintf(
		"INSERT INTO Trip (DriverID, CustomerID, PickupLocation, DropoffLocation, PickUpTime, DropOffTime, Status) VALUES('%s', %d, '%s', '%s', '%s', ' ', 'Pending')",
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
	query := fmt.Sprintf("UPDATE Trip SET DriverID = '%s', CustomerID = %d, PickUpLocation = '%s', DropOffLocation='%s', PickUpTime = '%s', DropOffTime = '%s', Status = '%s' WHERE TripID = %d",
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
func CheckCustomer(Email string) int {
	//Check Customer exsist in the database
	URL := "http://localhost:5000/api/v1/CheckUser/" + Email

	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
		return 0
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			data, err := strconv.Atoi(string(responseData))
			if err != nil {
				println(err)
			}
			return data
		}
	}
	return 0
}
func GetDriverPlate(DriverID string) string {
	//Check Customer exsist in the database
	URL := "http://localhost:4000/api/v1/GetDriverPlate/" + DriverID

	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			return string(responseData)
		}
	}
	return ""
}
func CheckDriver(Email string) string {
	URL := "http://localhost:4000/api/v1/CheckUser/" + Email
	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			return string(responseData)
		}
	}
	return ""
}
func GetAllDriver() []string {
	response, err := http.Get("http://localhost:4000/api/v1/GetAllDriver")
	if err != nil {
		fmt.Print(err.Error())
	}
	if response.StatusCode == http.StatusAccepted {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			IDs := strings.Split(string(responseData), ",")
			replacer := strings.NewReplacer(",", "")
			var newIDs []string
			for i := range IDs {
				newIDs = append(newIDs, replacer.Replace(IDs[i]))
			}
			return newIDs
		}
	}
	return nil
}
func APIRouter(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == "DELETE" {
		println("You can't delete Trip records")
	} else if r.Method == "GET" {
		params := mux.Vars(r)
		if params["TripID"] == " " {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Please provide TRIPID"))
			return
		}
		//GET trip using TripID
		TripID, err := strconv.Atoi(params["TripID"])
		if err != nil {
			fmt.Println(err)
		}
		TripInformation := GetTrip(db, TripID)
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else if TripInformation.TripID == 0 { // Check if data is empty
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("There is no trip session"))
			return
		} else {
			json.NewEncoder(w).Encode(GetTrip(db, TripInformation.TripID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			//POST Trip data with Customer Email and trip info
			var newTripData Data
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				fmt.Println(err)
			} else {
				err := json.Unmarshal(reqBody, &newTripData)
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("There was an error encoding the json."))
					return
				} else {
					if newTripData.CustomerEmail == "" {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply in JSON format"))
						return
					}
					var newTrip Trip
					newTrip.CustomerID = CheckCustomer(newTripData.CustomerEmail) //Check if user is in database
					if newTrip.CustomerID == 0 {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("Customer account not found"))
						return
					}
					//Find Avail driver
					AvailDriverID := CheckDriverAvail(db, GetAllDriver())
					if AvailDriverID != "" {
						newTrip.DriverID = AvailDriverID
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no avail driver"))
						return
					}
					//Add trip details
					newTrip.PickUpLocation = newTripData.PickUpLocation
					newTrip.DropOffLocation = newTripData.DropOffLocation
					newTrip.PickUpTime = newTripData.PickUpTime
					newTrip.DropOffTime = newTripData.DropOffTime

					//Check if customer or driver has uncompleted trips
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
				if updatedTrip.CustomerID == 0 || updatedTrip.DriverID == "" || updatedTrip.TripID == 0 {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply customer information"))
					return
				} else {
					if CheckTrip(db, updatedTrip.CustomerID, updatedTrip.DriverID) { //To check with the database if there is any record
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no trip records found"))
						return
					} else {
						//To update user details
						if EditTrip(db, updatedTrip) { //To edit trip
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Trip updated successfully"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("Trip unable to update"))
							return
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

	if params["Email"] == "" { // Check if data is empty
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("Please provide valid email"))
		return
	} else {
		CustomerID := CheckCustomer(params["Email"])
		if CustomerID != 0 {
			TripData := GetTrips(db, CustomerID)
			var JSONObject []AllTrip
			for _, data := range TripData {
				var TempAllTripData = AllTrip{PickUpLocation: data.PickUpLocation, DropOffLocation: data.DropOffLocation,
					PickUpTime: data.PickUpTime, DropOffTime: data.DropOffTime,
					Status: data.Status, DriverCarPlate: GetDriverPlate(data.DriverID)}
				JSONObject = append(JSONObject, TempAllTripData)
			}

			json.NewEncoder(w).Encode(JSONObject)
			w.WriteHeader(http.StatusAccepted)
		}

		return
	}
}

func main() {
	//API
	router := mux.NewRouter()
	//Web front-end CORS
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/Trip", home)                                                    //Test API
	router.HandleFunc("/api/v1/Trip/Router/{TripID}", APIRouter).Methods("GET", "PUT", "POST") //API Manipulation
	router.HandleFunc("/api/v1/Trip/{Email}", GetAllTrips)

	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router)))
}
