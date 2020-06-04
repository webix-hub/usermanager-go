package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"webix/usermanager/data"
)

func credentialsAPI(r chi.Router, dao *data.DAO) {

	r.Get("/user/{id}/credentials", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, _ := strconv.Atoi(id)

		data, err := dao.Credentials.GetByUser(idInt)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, data)
	})

	r.Put("/user/{id}/credentials", func(w http.ResponseWriter, r *http.Request) {
		// get update info
		decoder := json.NewDecoder(r.Body)
		cred := data.Credential{}
		err := decoder.Decode(&cred)
		if err != nil {
			errorResult(w, err)
			return
		}

		// ignore incoming data, regenerating credentials
		cred.Record = "password" // do not use in production !!!

		// do create credentials
		err = dao.Credentials.Save(&cred)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, Response{ID: cred.ID})
	})

	r.Post("/user/{id}/credentials", func(w http.ResponseWriter, r *http.Request) {
		// get update info
		decoder := json.NewDecoder(r.Body)
		upd := data.Credential{}
		err := decoder.Decode(&upd)
		if err != nil {
			errorResult(w, err)
			return
		}

		// do create credentials
		err = dao.Credentials.Save(&upd)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, Response{ID: upd.ID})
	})

	r.Delete("/user/{id}/credentials/{cid}", func(w http.ResponseWriter, r *http.Request) {
		cid := chi.URLParam(r, "cid")
		cidInt, _ := strconv.Atoi(cid)

		err := dao.Credentials.Delete(cidInt)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, Response{ID : cidInt})
	})
}
