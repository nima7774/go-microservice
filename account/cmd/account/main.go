package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nima7774/go-microservice/account"
	"github.com/tinrab/retry"
)

type config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to process env config: %v", err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("failed to connect to postgres, retrying...", err)
		}
		return
	})

	defer r.Close()
	log.Println("Listening on port 8080")
	s := account.NewAccountService(r)
	log.Fatal(account.ListenGRPC(s, 8080))

}
