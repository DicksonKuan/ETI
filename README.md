# ETI
Dickson, S10192803

1. Design consideration
- scaleable and elastic
* The microservice is able to talk to multiple database if there is a need for an storage upgrade or if one of the storage fails.
* JSON was used to transfer object data so as different services are compatible to access the microservice
* To allow HTTP rest from the web, the program has CORS enabled.

- Agility
* The database tables are not connected to each other so that each Microservice would have thier own database table. 
* To simulate the database are being setup in various machines, I have reroute the ports to show that the server would need to get from different Database.
* 3 Microservices are being setup to make the easier to develop, maintain and deploy

- Portability
* Instead of using the SQL Workbench, I have chosen to use XAMPP which shows there isn't any vendor lock in database.
- The 3 microservice are able to communicate with each other via API. This is to allow them to validate certain data

- Security
* With security as one of the main concern I have added a password for the user account to ensure whoever is accessing the SQL is a valid user. 
* With the 3 microservices communicating with each other via API, this restrict the amount of data they are able to gather from each database. 

- Format of REST API URL
* The API has the domain, API, Version, Services and Resource to allow user to quickly idenitfy which URL they need to access 


2. Architecture diagram


3. Instructions for setting up and running the microservice
3.1 Setting up database
- Download XAMPP
- Update Portnumber to 3308 in XAMPP Config file
* https://www.phpflow.com/php/how-to-change-xampp-apache-server-port/ - instructions
- Go to PHPMyAdmin
- Set password for all account
* UserAccount > root > change password to 'password'
- Create new database 
- Insert SQL code via 'rideshare.sql' OR 'SQL.txt'
* Please do not insert all the SQL statement at once for SQL.txt
- Click 'GO'

3.2 Run API
- Launch command prompt and run 
* Trip folder, Main.go
* Customer folder, Main.go
* Driver folder, Main.go

3.3 Test API
- Run html files 
* Front-end folder > Trip/passenger/ driver.html

Lib and platform used
1. https://github.com/gorilla/mux
2. XAMPP
3. https://github.com/gorilla/handlers
4. Visual studio code
5. https://github.com/go-sql-driver/mysql