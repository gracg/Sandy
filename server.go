package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var db dbHelper

func ServerMain(cnfg Config) {
	DBpath := cnfg.Server.SqLiteDB
	DBexist, err := exists(DBpath)

	//var db dbHelper

	if DBexist == false {
		db.CreateDatabase(DBpath)
	}

	db.NewHelper(DBpath)
	_, _ = DBexist, err

	router := mux.NewRouter()
	router.HandleFunc("/public", ServerPublicKey)
	router.HandleFunc("/group/{groupname}", ServerGroupKeys)

	http.ListenAndServe("127.0.0.1:"+strconv.Itoa(cnfg.Server.Port), router)
}

func ServerPublicKey(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getKey()))
}

func ServerGroupKeys(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	does, err := db.DoesGroupExist(params["groupname"])

	var res ResultKeys

	if does == false {
		res.Keys = []ResultKey{}
		res.IsError = true
		res.Error = "Group does not exist"

		jsonBytes, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
		return
	}

	if err != nil {
		res.Keys = []ResultKey{}
		res.IsError = true
		res.Error = err.Error()

		jsonBytes, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBytes)
		return
	}

	keys, err := db.GetGroupKeys(params["groupname"])
	if err != nil {
		res.Keys = []ResultKey{}
		res.IsError = true
		res.Error = err.Error()

		jsonBytes, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonBytes)
		return
	}

	res.Keys = []ResultKey{}
	fmt.Println(params["groupname"])
	for x := range keys {
		res.Keys = append(res.Keys, keys[x])
	}

	res.IsError = false

	jsonBytes, err := json.Marshal(&res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return

}
