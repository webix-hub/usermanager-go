package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"webix/usermanager/data"

	"github.com/go-chi/chi"
)

func logsAPI(r chi.Router, dao *data.DAO) {

	r.Get("/logs/by-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		data, err := dao.Logs.GetBy(id, 0, []data.LogType{data.LogUserAdded, data.LogUserChanged, data.LogUserDeleted, data.LogRoleAdded, data.LogRoleChanged, data.LogRoleDeleted})
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, data)
	})

	r.Get("/logs/login/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		data, err := dao.Logs.GetBy(id, 0, []data.LogType{data.LogLogin})
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, data)
	})

	r.Get("/logs/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		data, err := dao.Logs.GetBy(0, id, []data.LogType{data.LogUserAdded, data.LogUserChanged, data.LogUserDeleted})
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, data)
	})

	r.Get("/logs/role/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		data, err := dao.Logs.GetBy(0, id, []data.LogType{data.LogRoleAdded, data.LogRoleChanged, data.LogRoleDeleted})
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, data)
	})

	r.Get("/logs/all", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Logs.GetBy(0, 0, []data.LogType{data.LogUserAdded, data.LogUserChanged, data.LogUserDeleted, data.LogRoleAdded, data.LogRoleChanged, data.LogRoleDeleted})
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, data)
	})

	r.Get("/logs/meta", func(w http.ResponseWriter, r *http.Request) {
		meta := data.GetDescription()

		format.JSON(w, 200, meta)
	})

	// only for demo
	r.Post("/logs/add", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		log := data.Log{}
		err := decoder.Decode(&log)
		if err != nil {
			errorResult(w, err)
			return
		}

		err = dao.GetDB().Save(&log).Error
		if err != nil {
			errorResult(w, err)
			return
		}
		format.JSON(w, 200, &log)
	})
}
