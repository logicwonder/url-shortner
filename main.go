package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	h "github.com/logicwonder/url-shortner/http"
	mr "github.com/logicwonder/url-shortner/repository/mongo"
	rr "github.com/logicwonder/url-shortner/repository/redis"
	"github.com/logicwonder/url-shortner/shortner"
)

func main() {
	repo := chooseRepo()

	service := shortner.NewRedirectService(repo)
	handler := h.NewHandler(service)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortner.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		redisTimeout, _ := strconv.Atoi(os.Getenv("REDIS_TIMEOUT"))
		repo, err := rr.NewRedisRepository(redisURL, redisTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil

}
