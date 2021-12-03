package main

import (
	"bank/db"
	"bank/db/memdb"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//memdb.NewDatabase()
	http.HandleFunc("/users", userHandler)
	http.ListenAndServe(":8080", nil)
	// /defer memdb.NewDatabase().Db.Close()
}

/*
{
	"Name":"vijay",
	"Email":"vijay@test.com"
}
*/

var store db.Database = memdb.NewDatabase()

func userHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// TODO: crete user
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		u := &db.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		u, err = store.CreateUser(r.Context(), u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err = json.Marshal(u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)

		defer r.Body.Close()
	}
	if r.Method == http.MethodGet {
		users, err := store.ListAllUsers(r.Context())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err := json.Marshal(users)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
		defer r.Body.Close()
	}
	if r.Method == http.MethodDelete {
		Id := r.URL.Query().Get("Id")
		NewID, err1 := strconv.ParseInt(Id, 0, 0)
		if err1 != nil {
			fmt.Println("Error While Parsing Int")
		}
		err := store.DeleteUser(r.Context(), int(NewID))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write([]byte("user is deleted"))
		defer r.Body.Close()
	}
	if r.Method == http.MethodPatch {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		u := &db.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		// Id := r.URL.Query().Get("Id")
		// NewID, err1 := strconv.ParseInt(Id, 0, 0)
		// if err1 != nil {
		// 	fmt.Println("Error While Parsing Int")
		// }
		us, err := store.UserById(r.Context(), u.Id)
		if errors.Is(err, db.ErrNotFound) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		changed := false
		if u.Name != us.Name {
			changed = true
		}
		if u.Email != us.Email {
			changed = true
		}
		if !changed {
			return
		}
		err = store.UpdateUser(r.Context(), u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err = json.Marshal(u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
		defer r.Body.Close()
	}
}
