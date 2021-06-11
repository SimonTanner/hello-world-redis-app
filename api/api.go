package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/SimonTanner/hello-world-redis-app/redis"
	"github.com/gorilla/mux"
)

const addr = "127.0.0.1:6379"

type Api struct {
	Router      *mux.Router
	RedisClient redis.Client
}

func NewApi(redCli redis.Client) Api {
	router := mux.NewRouter()
	router.HandleFunc("/message/{key}", GetMessage).Methods("GET")
	router.HandleFunc("/add/{key}/{message}", SetMessage).Methods("GET")
	router.HandleFunc("/", GetMessages).Methods("GET")
	router.HandleFunc("/home", HomePage).Methods("GET", "POST")

	api := Api{
		Router:      router,
		RedisClient: redCli,
	}

	return api
}

func GetMessage(w http.ResponseWriter, r *http.Request) {

	redisClient := redis.NewClient(redis.RedisConf{
		Address: addr,
	})

	vars := mux.Vars(r)
	fmt.Println(vars)

	key := vars["key"]

	val, err := redisClient.Get(r.Context(), key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		jsvals, _ := json.Marshal(map[string]string{key: val})
		w.Write(jsvals)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address: addr,
	})

	vals, err := redisClient.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if len(vals) == 0 {
			fmt.Println("no messages found")
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusOK)
			jsvals, _ := json.Marshal(vals)
			w.Write(jsvals)
		}
	}
}

func SetMessage(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address: addr,
	})

	vars := mux.Vars(r)
	fmt.Println(vars["key"])
	fmt.Println(vars["message"])

	if err := redisClient.Set(r.Context(), vars["key"], vars["message"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

var testTemplate *template.Template

func HomePage(w http.ResponseWriter, r *http.Request) {
	redisClient := redis.NewClient(redis.RedisConf{
		Address: addr,
	})
	testTemplate, err := template.ParseFiles("./api/hello_world.html")
	fmt.Println(r.Method)
	if r.Method == "GET" {
		if err != nil {
			log.Fatalf("error parsing file: %e", err)
		}

		w.Header().Set("Content-Type", "text/html")
		err = testTemplate.Execute(w, nil)
		if err != nil {
			log.Fatalf("error rendering template: %e", err)
		}
	} else {
		r.ParseForm()
		fmt.Println(r.Form["message"])
		redisClient.Set(r.Context(), "message", r.Form["message"][0])
	}
}
