package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/SimonTanner/hello-world-redis-app/redis"
	"github.com/gorilla/mux"
)

const (
	addr    = "127.0.0.1:6379"
	expTime = time.Second * 10
)

type Api struct {
	Router      *mux.Router
	RedisClient redis.Client
}

func NewApi(redCli redis.Client) Api {
	router := mux.NewRouter()
	router.HandleFunc("/message/{key}", GetMessage).Methods("GET")
	router.HandleFunc("/", HomePage).Methods("GET", "POST")

	api := Api{
		Router:      router,
		RedisClient: redCli,
	}

	return api
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address:    addr,
		ExpireTime: expTime,
	})

	vars := mux.Vars(r)
	fmt.Println(vars)

	key := vars["key"]

	val, err := redisClient.Get(r.Context(), key)
	if err != nil {
		log.Printf("Error getting message: %e", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		jsvals, _ := json.Marshal(map[string]string{key: val})
		w.Write(jsvals)
	}
}

var testTemplate *template.Template

func HomePage(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address:    addr,
		ExpireTime: expTime,
	})

	testTemplate, err := template.ParseFiles("./api/hello_world.html")
	if err != nil {
		log.Fatalf("error parsing file: %e", err)
	}

	log.Print(r.Method)
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println(r.Form["message"])
		fmt.Println(r.Form["key"])
		redisClient.Set(r.Context(), r.Form["key"][0], r.Form["message"][0])
	}

	log.Print("Getting all messages from Redis")
	vals, err := redisClient.GetAll(r.Context())
	data := struct {
		Messages []string
	}{
		Messages: vals,
	}

	w.Header().Set("Content-Type", "text/html")
	err = testTemplate.Execute(w, data)
	if err != nil {
		log.Fatalf("error rendering template: %e", err)
	}
}
