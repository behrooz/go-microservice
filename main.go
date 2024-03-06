package main

func main() {
	a := App{}
	//initPriceTable(a.DB)
	a.Initialize(
		"root",    //os.Getenv("APP_DB_USERNAME"),
		"123",     //os.Getenv("APP_DB_PASSWORD"),
		"product", //os.Getenv("APP_DB_NAME"))
	)

	a.initializeRoutes()
	a.Run(":8010")

}
