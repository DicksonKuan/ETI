# ETI
Dickson, S10192803

1. Design consideration
- The database tables are not connected to each other so that each Microservice would have thier own database table. 
- The 3 microservice are able to communicate with each other via API. This is to allow them to validate certain data


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