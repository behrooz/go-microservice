# What is this repo

In this repositiry i'm trying to implement a CRUD go project with mysql database and gorm ORM

# How to run this application

In main.go file we have this code 
```
a.Initialize(
    os.Getenv("APP_DB_USERNAME"),
    os.Getenv("APP_DB_PASSWORD"),
    os.Getenv("APP_DB_NAME"))
)

```
As you can see here we have some varibales to initilize and we have .env file to set varibales in os env

```
source .env
```
then run 
```
go run main
```