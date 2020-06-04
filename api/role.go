package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"webix/usermanager/data"

	"github.com/go-chi/chi"
)

type RoleUpdate struct {
	data.Role

	Rules []int
}

func rolesAPI(r chi.Router, dao *data.DAO) {
	r.Get("/roles", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Roles.GetAll()

		if err != nil {
			errorResult(w, err)
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/roles", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		t := data.Role{}
		err := decoder.Decode(&t)
		if err != nil {
			errorResult(w, err)
			return
		}

		// ensure that this is a new object
		t.ID = 0
		err = dao.Roles.Save(&t)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, Response{t.ID})
	})

	r.Delete("/roles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		err := dao.Roles.Delete(id)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, Response{id})
	})

	r.Put("/roles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		// get existing role
		obj, err := dao.Roles.GetOne(id)
		if err != nil {
			errorResult(w, err)
			return
		}

		// get update info
		decoder := json.NewDecoder(r.Body)
		t := RoleUpdate{}
		err = decoder.Decode(&t)
		if err != nil {
			errorResult(w, err)
			return
		}

		if obj.Name != t.Name {
			dao.Roles.LogRoleChanged(obj, t.Name)
		}

		// update and save role object
		obj.Name = t.Name
		obj.Color = t.Color
		obj.Details = t.Details

		err = dao.Roles.Save(obj)
		if err != nil {
			errorResult(w, err)
			return
		}

		// update list of assigned roles
		rules, _ := dao.Roles.GetRules(id)
		added, removed := intersect(rules, t.Rules)

		for _, r := range added {
			dao.Roles.AddRule(obj, r)
		}

		for _, r := range removed {
			dao.Roles.RemoveRule(obj, r)
		}

		format.JSON(w, 200, Response{id})
	})
}
