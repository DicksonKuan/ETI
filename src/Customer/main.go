package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID           int    `json:"CustomerID"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	MobileNumber string `json:"MobileNumber"`
	EmailAddress string `json:"EmailAddress"`
	Password     string `json:"Password"`
}

//Database function
func CheckUser(db *sql.DB, email string) bool {
	query := fmt.Sprintf("Select * FROM Customer WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
		return false
	}
	var customer Customer
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&customer.ID, &customer.FirstName,
			&customer.LastName, &customer.MobileNumber, &customer.EmailAddress, &customer.Password)
		if err != nil {
			panic(err.Error())
		} else if customer.EmailAddress == email {
			return true
		}
	}
	return false
}

func GetUser(db *sql.DB, email string, password string) Customer {
	query := fmt.Sprintf("Select * FROM Customer WHERE EmailAddress= '%s' AND Password= '%s'", email, password)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var customer Customer
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&customer.ID, &customer.FirstName,
			&customer.LastName, &customer.MobileNumber, &customer.EmailAddress, &customer.Password)
		if err != nil {
			panic(err.Error())
		}
	}
	return customer
}

func CreateNewUser(db *sql.DB, customer Customer) bool {
	query := fmt.Sprintf(
		"INSERT INTO Customer(FirstName, LastName, MobileNumber, EmailAddress, Password) VALUES('%s','%s','%s','%s','%s');",
		customer.FirstName,
		customer.LastName,
		customer.MobileNumber,
		customer.EmailAddress,
		customer.Password)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
		return false
	}
	return true
}
func DeleteUser(db *sql.DB, customer Customer) bool {
	query := fmt.Sprintf("DELETE FROM Trip WHERE CustomerID =%d", customer.ID)
	_, err := db.Query(query)

	if err != nil {
		return false
	} else {
		query := fmt.Sprintf("DELETE FROM Customer WHERE CustomerID =%d", customer.ID)
		_, err := db.Query(query)

		if err != nil {
			panic(err.Error())
			return false
		}
		return true
	}

}

func EditUser(db *sql.DB, customer Customer) bool {
	if customer.ID == 0 {
		return false
	}
	query := fmt.Sprintf("UPDATE Customer SET FirstName = '%s', LastName = '%s', MobileNumber= '%s'	WHERE ID = %d;", customer.FirstName, customer.LastName, customer.MobileNumber, customer.ID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

//API Function
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
		println("DELETE is working")
	}

	if r.Header.Get("Content-type") == "application/json" {
		//Database
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
		if err != nil {
			fmt.Println(err)
		}
		if r.Method == "GET" { //GET User
			var loginInformation Customer
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &loginInformation)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("There was an error encoding the json. err = %s", err)
				} else if loginInformation.EmailAddress != "" || loginInformation.Password != "" { // Check if data is empty
					json.NewEncoder(w).Encode(GetUser(db, loginInformation.EmailAddress, loginInformation.Password))
					w.WriteHeader(http.StatusAccepted)
					return
				} else {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"Account is invalid"))
					return
				}
			}
		}
		if r.Method == "POST" {
			var newCustomer Customer
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newCustomer)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("There was an error encoding the json. err = %s", err)
				} else {
					if newCustomer.EmailAddress == "" {
						w.WriteHeader(
							http.StatusUnprocessableEntity)
						w.Write([]byte(
							"422 - Please supply in JSON format"))
						return
					} else {
						if !CheckUser(db, newCustomer.EmailAddress) {
							CreateNewUser(db, newCustomer)
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Account created successfully"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("This email address is in use"))
							return
						}
					}
				}
			}
		} else if r.Method == "PUT" {
			var customerInformation Customer
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &customerInformation)

				if customerInformation.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply customer information"))
					return
				} else {
					if !CheckUser(db, customerInformation.EmailAddress) { //To check with the database if there is any record
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no exsiting account with for " + customerInformation.EmailAddress))
					} else {
						//To update user details
						EditUser(db, customerInformation)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Account updated successfully"))
					}
				}
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RideShare Passenger API")
}

func main() {
	// // //Database
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3308)/rideshare") //Connecting to database
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var customer Customer
	// //Dummy data
	// customer.EmailAddress = "lily@gmail.com"
	// customer.FirstName = "Pika"
	// customer.LastName = "John"
	// customer.MobileNumber = "91234568"
	// customer.Password = "password"

	// //To get User
	// customer = GetUser(db, "John@np.com", "password") - Working
	// println(string(JsonResult)) - Working
	// To CREATE new user
	//CreateNewUser(db, customer) - Working
	//To UPDATE User
	//println(EditUser(db, customer)) - Working

	//jSON serlised
	// JsonResult, err := json.Marshal(
	// 	Customer{
	// 		ID:           customer.ID,
	// 		FirstName:    customer.FirstName,
	// 		LastName:     customer.LastName,
	// 		EmailAddress: customer.EmailAddress,
	// 		Password:     customer.Password})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// println(string(JsonResult))

	//API part
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/Passenger", home)                                           //Test API
	router.HandleFunc("/api/v1/Passenger/Router", APIRouter).Methods("GET", "PUT", "POST") //API Manipulation

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
