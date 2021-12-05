package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Driver struct {
	DriverID         string `json: "DriverID"`
	FirstName        string `json: "FirstName"`
	LastName         string `json: "LastName"`
	MobileNumber     string `json: "MobileNumber"`
	EmailAddress     string `json: "EmailAddress"`
	CarLicenseNumber string `json: "CarLicenseNumber"`
	Password         string `json: "Password"`
}

//Database function
func CheckUser(db *sql.DB, email string) bool {
	query := fmt.Sprintf("Select * FROM Driver WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&driver.DriverID, &driver.FirstName,
			&driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.Password, &driver.CarLicenseNumber)
		if err != nil {
			panic(err.Error())
		} else if driver.EmailAddress == email {
			return true
		}
	}
	return false
}

func GetUser(db *sql.DB, email string, password string) Driver {
	query := fmt.Sprintf("Select * FROM Driver WHERE EmailAddress= '%s' AND Password= '%s'", email, password)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&driver.DriverID, &driver.FirstName,
			&driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.Password, &driver.CarLicenseNumber)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver
}

func GetUserWithEmail(db *sql.DB, email string) string {
	query := fmt.Sprintf("Select * FROM Driver WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&driver.DriverID, &driver.FirstName,
			&driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.Password, &driver.CarLicenseNumber)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver.DriverID
}

func CreateNewUser(db *sql.DB, driver Driver) bool {
	if CheckUser(db, driver.EmailAddress) {
		return false
	}
	query := fmt.Sprintf(
		"INSERT INTO Driver(DriverID, FirstName, LastName, MobileNo, EmailAddress, Password, CarLiscenseNo) VALUES('%s','%s','%s','%s','%s','%s', '%s');",
		driver.DriverID,
		driver.FirstName,
		driver.LastName,
		driver.MobileNumber,
		driver.EmailAddress,
		driver.Password,
		driver.CarLicenseNumber)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditUser(db *sql.DB, driver Driver) bool {
	if driver.DriverID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE Driver SET FirstName = '%s', LastName = '%s', MobileNo= '%s', CarLiscenseNo = '%s', Password = '%s', EmailAddress = '%s' WHERE DriverID = '%s';", driver.FirstName, driver.LastName, driver.MobileNumber, driver.CarLicenseNumber, driver.Password, driver.EmailAddress, driver.DriverID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func SearchAvailDriver(db *sql.DB) string {
	//GET All driver ID
	query := "SELECT DriverID FROM driver"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var DriverIDs string
	for results.Next() {
		var TempID string
		err = results.Scan(&TempID)
		if err != nil {
			panic(err.Error())
		}
		DriverIDs += TempID + ","
	}
	return DriverIDs
}

//API
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

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RideShare Driver API")
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
		println("You can't delete Driver account")
	} else if r.Method == "GET" { //GET User
		params := mux.Vars(r)
		if params["Email"] == " " && params["Password"] == " " { //Check if email and password has value
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Please provide driver's account details"))
			return
		}
		loginInformation := GetUser(db, params["Email"], params["Password"])        //Get user according to email and password
		if loginInformation.EmailAddress != "" || loginInformation.Password != "" { // Check if data is empty
			json.NewEncoder(w).Encode(loginInformation)
			w.WriteHeader(http.StatusAccepted)
			return
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Account is invalid"))
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" { //CREATE new Driver
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newDriver)
				if err != nil {
					fmt.Printf("There was an error encoding the json. err = %s", err)
				} else {
					if newDriver.EmailAddress == "" {
						w.WriteHeader(
							http.StatusUnprocessableEntity)
						w.Write([]byte(
							"422 - Please supply in JSON format"))
						return
					} else {
						if !CheckUser(db, newDriver.EmailAddress) { //Check if user exsists
							CreateNewUser(db, newDriver)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Account has bee created successfully"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("This email address is in use"))
							return
						}
					}
				}
			}
		} else if r.Method == "PUT" { //UPDATE Driver Details
			var DriverInformation Driver
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &DriverInformation)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("There was an error encoding the json. err = %s", err)
				}
				if DriverInformation.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply customer information"))
					return
				} else {
					if !CheckUser(db, DriverInformation.EmailAddress) { //To check with the database if there is any record
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no exsiting account with for " + DriverInformation.EmailAddress))
					} else {
						if EditUser(db, DriverInformation) { //To update user details
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Account updated successfully"))
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("Account unable to update"))
						}
					}
				}
			}
		}
	}
}

func CheckDriver(w http.ResponseWriter, r *http.Request) {
	//Database
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["UserEmail"] == "" {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte(
			"422 - Please supply email information"))
		return
	} else if CheckUser(db, params["UserEmail"]) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(GetUserWithEmail(db, params["UserEmail"])))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func GetAllDriver(w http.ResponseWriter, r *http.Request) {
	//Database
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(SearchAvailDriver(db)))
}

func main() {
	// //Database
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var driver Driver
	// //Dummy data
	// driver.EmailAddress = "lily@gmail.com"
	// driver.FirstName = "YuEn"
	// driver.LastName = "John"
	// driver.MobileNumber = "91234568"
	// driver.Password = "password"
	// driver.CarLicenseNumber = "S101II"
	// driver.ID = 2
	//println(GetUser(db, "lily@gmail.com", "password").FirstName) - working
	//CreateNewUser(db, driver) - Working
	//println(EditUser(db, driver)) - Working

	//API part
	router := mux.NewRouter()
	//Web Front-end CORS
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/GetAllDriver", GetAllDriver)                                                //Get All driver
	router.HandleFunc("/api/v1/CheckUser/{UserEmail}", CheckDriver)                                        //Check Driver exsist
	router.HandleFunc("/api/v1/Driver", home)                                                              //Test API
	router.HandleFunc("/api/v1/Driver/Router/{Email}/{Password}", APIRouter).Methods("GET", "PUT", "POST") //API Manipulation

	fmt.Println("Listening at port 4000")
	log.Fatal(http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(router)))
}
