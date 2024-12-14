package main

import (
	"fmt"
	"gorestserviceagain/api"
	"gorestserviceagain/postgresdb"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg, err := configFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// db, err := sqldb.NewSqlite(cfg.DSN)
	db, err := postgresdb.NewPostgres(cfg.DSN)
	if err != nil {
		log.Fatal(nil)
	}
	db.Migrate()

	r := chi.NewRouter()

	r.Route("/orders", api.OrdersController{Repo: db}.RegisterRoutes)
	r.Route("/dishes", api.DishesController{Repo: db}.RegisterRoutes)
	api.DiscountDetailController{Repo: db}.RegisterRoutes(r)
	fmt.Println("Staring serve on", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
