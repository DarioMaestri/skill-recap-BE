package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/DarioMaestri/01_skill-recap/skilldomain"
	"github.com/DarioMaestri/01_skill-recap/userdomain"
	"github.com/DarioMaestri/01_skill-recap/userskilldomain"
)

var db *sql.DB
var err error

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSkill struct {
	UserId  int `json:"userId"`
	SkillId int `json:"skillId"`
}

func main() {

	fmt.Println("Drivers: ", sql.Drivers())

	fmt.Println("Open DB")
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/sys")
	if err != nil {
		log.Fatal("Error couldn't connect to db:", err)
	}

	defer db.Close()
	fmt.Println("DB OK")

	fmt.Println("Start server")
	r := mux.NewRouter()

	skilldomain.HandleSkill(r, db)
	userdomain.HandleUser(r, db)
	userskilldomain.HandleUserSkill(r, db)

	fmt.Println("Server OK")

	fmt.Println("GO running")

	http.ListenAndServe(":8080", r)
}
