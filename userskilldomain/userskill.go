package userskilldomain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DarioMaestri/01_skill-recap/jsonres"
	"github.com/DarioMaestri/01_skill-recap/skilldomain"
	"github.com/DarioMaestri/01_skill-recap/userdomain"
)

type UserSkill struct {
	UserId  int `json:"user_id"`
	SkillId int `json:"skill_id"`
}

type UserSkills []UserSkill

const Table string = "user_skill"

var db *sql.DB

//
// DATABASE
//

func FindAll(db *sql.DB) (UserSkills, error) {
	var userSkills UserSkills

	result, err := db.Query(fmt.Sprintf("SELECT * FROM %s", Table))
	if err != nil {
		log.Println("Error loading userSkills from database: ", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var userSkill UserSkill
		err = result.Scan(&userSkill.UserId, &userSkill.SkillId)
		if err != nil {
			log.Println("Error while parsing data from database: ", err)
			return nil, err
		}
		userSkills = append(userSkills, userSkill)
	}

	return userSkills, nil
}

func FindAllByUserId(db *sql.DB, id string) (skilldomain.Skills, error) {
	var skills skilldomain.Skills

	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, name, version FROM %s JOIN %s ON user_id=id WHERE user_id=?", Table, skilldomain.Table))
	if err != nil {
		log.Println("Error while preparing statement: ", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query(id)
	if err != nil {
		log.Println("Error while retrieving by user_id: ", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var skill skilldomain.Skill
		err = result.Scan(&skill.Id, &skill.Name, &skill.Version)
		if err != nil {
			log.Println("Error while parsing data from database: ", err)
			return nil, err
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

func FindAllBySkillId(db *sql.DB, id string) (userdomain.Users, error) {
	var users userdomain.Users

	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, username FROM %s JOIN %s ON user_id=id WHERE skill_id=?", Table, userdomain.Table))
	if err != nil {
		log.Println("Error while preparing statement: ", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query(id)
	if err != nil {
		log.Println("Error while retrieving by user_id: ", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var user userdomain.User
		err = result.Scan(&user.Id, &user.Username)
		if err != nil {
			log.Println("Error while parsing data from database: ", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func Insert(db *sql.DB, userskillp *UserSkill) error {

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (user_id, skill_id) VALUES(?, ?)", Table))
	if err != nil {
		log.Println("Error unable to prepare statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userskillp.UserId, userskillp.SkillId)
	if err != nil {
		log.Println("Error while inserting new userSkill: ", err)
		return err
	}

	return nil
}

func Delete(db *sql.DB, userId string, skillId string) error {

	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND skill_id=?", Table))
	if err != nil {
		log.Println("Error unable to prepare statement: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, skillId)
	if err != nil {
		log.Println("Error while executing statement: ", err)
		return err
	}

	return nil

}

//
// ENTRYPOINT
//

func Alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func GetUserskills(w http.ResponseWriter, r *http.Request) {

	userSkills, err := FindAll(db)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error invoking findAll: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, &userSkills)
	}

}

func GetUsersBySkill(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["skillId"]

	users, err := FindAllBySkillId(db, id)
	if users == nil {
		err = sql.ErrNoRows
	}
	switch err {
	case sql.ErrNoRows:
		jsonres.RespondWithError(w, http.StatusNotFound, "No rows returned by id: "+id)
	case nil:
		jsonres.RespondWithJSON(w, http.StatusOK, &users)
	default:
		log.Println("Error while fetching a single row from db: ", err)
		jsonres.RespondWithError(w, http.StatusNotFound, "Error fetching a single row from db")
	}

}

func GetSkillsByUser(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["userId"]

	skills, err := FindAllByUserId(db, id)
	if skills == nil {
		err = sql.ErrNoRows
	}
	switch err {
	case sql.ErrNoRows:
		jsonres.RespondWithError(w, http.StatusNotFound, "No rows returned by id: "+id)
	case nil:
		jsonres.RespondWithJSON(w, http.StatusOK, &skills)
	default:
		log.Println("Error while fetching a single row from db: ", err)
		jsonres.RespondWithError(w, http.StatusNotFound, "Error fetching a single row from db")
	}

}

func InsertUserskill(w http.ResponseWriter, r *http.Request) {
	var userSkill UserSkill

	_ = json.NewDecoder(r.Body).Decode(&userSkill)

	fmt.Println(userSkill)

	err := Insert(db, &userSkill)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while inserting new userSkill: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "OK"})
	}

}

func DeleteUserskill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	skillId := vars["skillId"]

	err := Delete(db, userId, skillId)
	if err != nil {
		jsonres.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting new userSkill: %s", err))
	} else {
		jsonres.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "OK"})
	}

}

func HandleUserSkill(r *mux.Router, sqlDB *sql.DB) {

	db = sqlDB

	r.HandleFunc("/userskill/alive", Alive).Methods("GET")
	r.HandleFunc("/userskill", GetUserskills).Methods("GET")
	r.HandleFunc("/userskill/skills/{userId}", GetSkillsByUser).Methods("GET")
	r.HandleFunc("/userskill/users/{skillId}", GetUsersBySkill).Methods("GET")
	r.HandleFunc("/userskill", InsertUserskill).Methods("POST")
	r.HandleFunc("/userskill/{userId}/{skillId}", DeleteUserskill).Methods("DELETE")
}
