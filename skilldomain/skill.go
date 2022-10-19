package skilldomain

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

type Skill struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Skills []Skill

const Table string = "skill"

var db *sql.DB

//
// DATABASE
//

func FindAll(db *sql.DB) (Skills, error) {
	var skills Skills

	result, err := db.Query(fmt.Sprintf("SELECT * FROM %s", Table))
	if err != nil {
		log.Println("Error loading skills from database: ", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var skill Skill
		err = result.Scan(&skill.Id, &skill.Name, &skill.Version)
		if err != nil {
			log.Println("Error while parsing data from database: ", err)
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func FindById(db *sql.DB, id string) (Skill, error) {
	var skill Skill

	err := db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id=?", Table), id).Scan(&skill.Id, &skill.Name, &skill.Version)
	if err != nil {
		log.Println("Error loading skill from database: ", err)
		return skill, err
	}

	return skill, nil
}

func Insert(db *sql.DB, skillp *Skill) (int64, error) {

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (name, version) VALUES(?, ?)", Table))
	if err != nil {
		log.Println("Error unable to prepare statement:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(skillp.Name, skillp.Version)
	if err != nil {
		log.Println("Error while inserting new skill: ", err)
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

func Update(db *sql.DB, skillp *Skill, id string) (int64, error) {

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET name = ?, version = ? WHERE id = ?", Table))
	if err != nil {
		log.Println("Error unable to prepare statement: ", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(skillp.Name, skillp.Version, id)
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

func GetSkills(w http.ResponseWriter, r *http.Request) {

	skills, err := FindAll(db)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error invoking findAll: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, &skills)
	}

}

func GetSkill(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	skill, err := FindById(db, id)
	switch err {
	case sql.ErrNoRows:
		jsonres.RespondWithError(w, http.StatusNotFound, "No rows returned by id: "+id)
	case nil:
		jsonres.RespondWithJSON(w, http.StatusOK, &skill)
	default:
		log.Println("Error while fetching a single row from db: ", err)
		jsonres.RespondWithError(w, http.StatusNotFound, "Error fetching a single row from db")
	}

}

func InsertSkill(w http.ResponseWriter, r *http.Request) {
	var skill Skill

	_ = json.NewDecoder(r.Body).Decode(&skill)

	idr, err := Insert(db, &skill)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while inserting new skill: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, idr)
	}

}

func DeleteSkill(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	idr, err := Delete(db, id)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting skill id %s: %s", id, err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, idr)
	}

}

func UpdateSkill(w http.ResponseWriter, r *http.Request) {
	var skill Skill

	id := mux.Vars(r)["id"]

	_ = json.NewDecoder(r.Body).Decode(&skill)

	idr, err := Update(db, &skill, id)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating skill id %s: %s", id, err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, idr)
	}

}

func HandleSkill(r *mux.Router, sqlDB *sql.DB) {

	db = sqlDB

	r.HandleFunc("/skills/alive", Alive).Methods("GET")
	r.HandleFunc("/skills", GetSkills).Methods("GET")
	r.HandleFunc("/skills/{id}", GetSkill).Methods("GET")
	r.HandleFunc("/skills", InsertSkill).Methods("POST")
	r.HandleFunc("/skills/{id}", DeleteSkill).Methods("DELETE")
	r.HandleFunc("/skills/{id}", UpdateSkill).Methods("PUT")
}
