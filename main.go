package main

import (
	"net/http"

	postgres "github.com/cassiusbessa/satback/db"
	controllers "github.com/cassiusbessa/satback/http"
)

func main() {
	postgres.CreateDb()
	controllers.RegisterHandlers()
	http.ListenAndServe(":8080", nil)
}
