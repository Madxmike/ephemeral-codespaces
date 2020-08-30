package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tierzer0/gateway/auth"
	"github.com/tierzer0/gateway/environment"
	"github.com/tierzer0/gateway/redis"
)

var (
	Port              = flag.String("port", "8080", "the port to serve on")
	RedisAddress      = flag.String("redis", "", "the address of redis")
	FirebaseProjectID = flag.String("firebase", "", "the firebase project id")
)

func main() {
	flag.Parse()
	if *RedisAddress == "" {
		log.Fatal("redis address not provided")
	}

	if *FirebaseProjectID == "" {
		log.Fatal("firebase project id not provided")
	}

	redisConn, err := redigo.Dial("tcp", *RedisAddress)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not establish redis connection"))
	}
	publisher := redis.Publisher{
		Conn: &redisConn,
	}

	authenticator := auth.Authenticator{
		ProjectID: *FirebaseProjectID,
	}
	authenticator.RetrievePublicKeys()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/environment", environment.Routes(&authenticator, publisher))

	go func() {
		log.Println(http.ListenAndServe(":"+*Port, r))
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

type RedisPublisher struct {
}
