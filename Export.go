package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Chat co_chat
type Chat struct {
	Time               int64
	User, Wid, X, Y, Z int
	Message            string
}

//User co_user
type User struct {
	ID   int
	Time int64
	User string
	UUID string
}

var db *gorm.DB

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.FormValue("key") != "MscWeekily" {
			w.Write([]byte("wrong key"))
			return
		}
		getMsg(w)
	})
	http.ListenAndServe("0.0.0.0:45777", nil)
}

func getMsg(w http.ResponseWriter) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "co_" + defaultTableName
	}
	var err error
	db, err = gorm.Open("sqlite3", "database.db")
	defer db.Close()
	db.SingularTable(true)
	if err != nil {
		log.Fatal(err)
	}

	var data []Chat
	db.Where("time > ?", time.Now().AddDate(0, 0, -7).Unix()).Find(&data)
	for _, v := range data {
		w.Write([]byte(time.Unix(v.Time, 0).Format("[2006-01-02 15:04:05]") + getName(v.User) + ":" + v.Message + "\n"))
	}
}

func getName(id int) string {
	var user User
	db.Where("id = ?", id).First(&user)
	return user.User
}
