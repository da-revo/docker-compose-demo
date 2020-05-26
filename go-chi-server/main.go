package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mgoSession *mgo.Session

var db *DB

// M is short for bson.M
type M bson.M

// var secretKey = []byte(os.Getenv("SECRET"))
var secretKey = []byte("A_SECRET")

// DB stores the database session, and collections information.
type DB struct {
	food *mgo.Collection
}

// NewDB is used to create DB instance
func NewDB(
	food *mgo.Collection,
) *DB {
	return &DB{
		food: food,
	}
}

// Food holds the data for food
type Food struct {
	ID   *bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string         `json:"name,omitempty" bson:"name,omitempty"`
}

func main() {

	fmt.Println("Starting Server at", time.Now())

	mgoSession, err := mgo.Dial("database-mongodb:27017")
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()

	food := mgoSession.DB("test1").C("food")
	db = NewDB(
		food,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world, you have reached the go-chi-server with mongodb !!! \n\n /add \n /show"))
	})

	r.Get("/add", func(w http.ResponseWriter, r *http.Request) {
		newItem := Food{
			Name: "pasta",
		}
		err := db.food.Insert(newItem)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("added pasta"))
		}
	})

	r.Get("/show", func(w http.ResponseWriter, r *http.Request) {

		var items []Food
		err = db.food.Find(M{}).All(&items)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			response, _ := json.Marshal(items)
			w.Write(response)
		}
	})

	http.ListenAndServe(":8080", r)
}
