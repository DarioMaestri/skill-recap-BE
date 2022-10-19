package userdomain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DarioMaestri/01_skill-recap/jsonres"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Users []User

const Table string = "user"

var db *sql.DB

//
// DATABASE
//

func FindAll(db *sql.DB) (Users, error) {
	var users Users

	result, err := db.Query(fmt.Sprintf("SELECT * FROM %s", Table))
	if err != nil {
		log.Println("Error loading users from database: ", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var user User
		err = result.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Println("Error while parsing data from database: ", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func FindById(db *sql.DB, id string) (User, error) {
	var user User

	err := db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id=?", Table), id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println("Error loading user from database: ", err)
		return user, err
	}

	return user, nil
}

func FindByUsernamePassword(db *sql.DB, userp *User) (User, error) {
	var user User

	err := db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE username=? AND password=?", Table), userp.Username, userp.Password).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println("Error loading user from database: ", err)
		return user, err
	}

	return user, nil
}

func Insert(db *sql.DB, userp *User) (int64, error) {

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (username, password) VALUES(?, ?)", Table))
	if err != nil {
		log.Println("Error unable to prepare statement:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userp.Username, userp.Password)
	if err != nil {
		log.Println("Error while inserting new user: ", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error while retrieving new id: ", err)
		return 0, err
	}

	return id, nil
}

func Delete(db *sql.DB, id string) (int64, error) {

	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?", Table))
	if err != nil {
		log.Println("Error unable to prepare statement: ", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		log.Println("Error while executing statement: ", err)
		return 0, err
	}

	idr, err := result.RowsAffected()
	if err != nil {
		log.Println("Error while retrieving row affected id: ", err)
		return 0, err
	}

	return idr, nil

}

func Update(db *sql.DB, userp *User, id string) (int64, error) {

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET username = ?, password = ? WHERE id = ?", Table))
	if err != nil {
		log.Println("Error unable to prepare statement: ", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userp.Username, userp.Password, id)
	if err != nil {
		log.Println("Error while executing statement: ", err)
		return 0, err
	}

	idr, err := result.RowsAffected()
	if err != nil {
		log.Println("Error while retrieving row affected id: ", err)
		return 0, err
	}

	return idr, nil

}

//
// ENTRYPOINT
//

func Alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := FindAll(db)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error invoking findAll: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, &users)
	}

}

func GetUser(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	user, err := FindById(db, id)
	switch err {
	case sql.ErrNoRows:
		jsonres.RespondWithError(w, http.StatusNotFound, "No rows returned by id: "+id)
	case nil:
		jsonres.RespondWithJSON(w, http.StatusOK, &user)
	default:
		log.Println("Error while fetching a single row from db: ", err)
		jsonres.RespondWithError(w, http.StatusNotFound, "Error fetching a single row from db")
	}

}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)

	idr, err := Insert(db, &user)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while inserting new user: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, idr)
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	idr, err := Delete(db, id)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting user id %s: %s", id, err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, idr)
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	id := mux.Vars(r)["id"]

	_ = json.NewDecoder(r.Body).Decode(&user)

	idr, err := Update(db, &user, id)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating user id %s: %s", id, err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, idr)
	}

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)

	fmt.Println(user)

	_, err := FindByUsernamePassword(db, &user)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Credentials not correct: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, true)
	}

}

func HandleUser(r *mux.Router, sqlDB *sql.DB) {

	db = sqlDB

	r.HandleFunc("/users/alive", Alive).Methods("GET")
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users", InsertUser).Methods("POST")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/login", LoginUser).Methods("POST", "OPTIONS")
}
