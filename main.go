package main

import (
	"fmt"
	"net/http"
	"time"

	postgres "github.com/cassiusbessa/satback/db"
	controllers "github.com/cassiusbessa/satback/http"
)

func main() {
	for i := 0; i < 30; i++ {
		err := postgres.CreateDb()
		if err != nil {
			fmt.Printf("Erro ao conectar ao banco de dados: %v\n", err)
			time.Sleep(time.Second * 2)
		} else {
			break
		}
	}
	controllers.RegisterHandlers()
	http.ListenAndServe(":8080", nil)
}
