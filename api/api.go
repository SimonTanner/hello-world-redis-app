package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/SimonTanner/hello-world-redis-app/redis"
	"github.com/gorilla/mux"
)

const (
	addr    = "127.0.0.1:6379"
	expTime = time.Second * 300
)

var testTemplate *template.Template

type Api struct {
	Router *mux.Router
}

func NewApi() Api {
	router := mux.NewRouter()
	router.HandleFunc("/message/{key}", GetMessage).Methods("GET")
	router.HandleFunc("/", HomePage).Methods("GET", "POST")

	api := Api{
		Router: router,
	}

	return api
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address:    addr,
		ExpireTime: expTime,
	})

	vars := mux.Vars(r)
	key := vars["key"]

	val, err := redisClient.Get(r.Context(), key)
	if err != nil {
		log.Printf("Error getting message: %e", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		jsvals, _ := json.Marshal(val)
		w.Write(jsvals)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address:    addr,
		ExpireTime: expTime,
	})

	testTemplate, err := template.ParseFiles("./api/hello_world.html")
	if err != nil {
		log.Fatalf("error parsing file: %e", err)
	}

	if r.Method == "POST" {
		r.ParseForm()
		redisClient.Set(r.Context(), r.Form.Get("key"), redis.Message{Str: r.Form.Get("message")})
	}

	log.Print("Getting all messages from Redis")
	var msgs []redis.Message
	msgs, err = redisClient.GetAll(r.Context())
	if err != nil {
		log.Printf("error rendering template: %e", err)
	}

	data := struct {
		Messages []redis.Message
	}{
		Messages: msgs,
	}

	w.Header().Set("Content-Type", "text/html")
	err = testTemplate.Execute(w, data)
	if err != nil {
		log.Fatalf("error rendering template: %e", err)
	}
}
