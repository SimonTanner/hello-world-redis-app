package api

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc("/{key}", GetMessage).Methods("GET")
	router.HandleFunc("/addmessage/{key}/{message}", SetMessage).Methods("GET")
	router.HandleFunc("/", GetMessages).Methods("GET")

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

	fmt.Println("what??")

	if err := redisClient.Set(r.Context(), vars["key"], vars["message"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
