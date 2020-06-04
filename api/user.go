package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"webix/usermanager/data"

	"github.com/go-chi/chi"
)

type UserUpdate struct {
	data.User

	Roles []int
	Rules []int
}

type UploadResponse struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

func usersAPI(r chi.Router, dao *data.DAO, path, server string) {
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		d, err := dao.Users.GetAll()

		if err != nil {
			errorResult(w, err)
		} else {
			format.JSON(w, 200, d)
		}
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		t := UserUpdate{}
		err := decoder.Decode(&t)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}

		// ensure that this is a new object
		u := data.User{
			ID:      0,
			Name:    t.Name,
			Email:   t.Email,
			Details: t.Details,
			Status:  t.Status,
		}

		dao.Users.Save(&u)

		format.JSON(w, 200, Response{u.ID})
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		err := dao.Users.Delete(id)
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, Response{id})
	})

	r.Put("/users/{id}/anonymize", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		// get existing user
		u, err := dao.Users.GetOne(id)
		if err != nil {
			errorResult(w, err)
			return
		}

		// update and save user object
		u.Name = "anonymous"
		u.Email = "anonymous@some.com"
		u.Avatar = ""

		err = dao.Users.Save(u)
		if err != nil {
			errorResult(w, err)
			return
		}
		format.JSON(w, 200, Response{id})
	})

	r.Post("/users/{id}/avatar", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		// the max upload size
		var limit = int64(4 << 20) // default is 4MB
		r.Body = http.MaxBytesReader(w, r.Body, limit)
		r.ParseMultipartForm(limit)

		file, _, err := r.FormFile("upload")
		if err != nil {
			errorResult(w, err)
			return
		}
		defer file.Close()

		u, err := updateAvatar(id, file, path, server, dao)
		if err != nil {
			errorResult(w, err)
		} else {
			format.JSON(w, 200, &UploadResponse{Status: "server", Value: u.Avatar})
		}
	})

	r.Get("/users/{id}/avatar/{name}", func(w http.ResponseWriter, r *http.Request) {
		//id := chi.URLParam(r, "id")
		name := chi.URLParam(r, "name")
		http.ServeFile(w, r, filepath.Join(path, "avatars", name))
	})

	r.Put("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		// get existing user
		u, err := dao.Users.GetOne(id)
		if err != nil {
			errorResult(w, err)
			return
		}

		// get update info
		decoder := json.NewDecoder(r.Body)
		t := UserUpdate{}
		err = decoder.Decode(&t)
		if err != nil {
			errorResult(w, err)
			return
		}

		if u.Email != t.Email {
			dao.Users.LogAddEmail(u, t.Email)
		}

		// update and save user object
		u.Name = t.Name
		u.Email = t.Email
		u.Status = t.Status
		u.Details = t.Details

		err = dao.Users.Save(u)
		if err != nil {
			errorResult(w, err)
			return
		}

		// update list of assigned roles
		roles, _ := dao.Users.GetRoles(id)
		added, removed := intersect(roles, t.Roles)

		for _, r := range added {
			dao.Users.AddRole(u, r)
		}

		for _, r := range removed {
			dao.Users.RemoveRole(u, r)
		}

		// update list of assigned rules
		rules, _ := dao.Users.GetRules(id)
		added, removed = intersect(rules, t.Rules)

		for _, r := range added {
			dao.Users.AddRule(u, r)
		}

		for _, r := range removed {
			dao.Users.RemoveRule(u, r)
		}

		format.JSON(w, 200, Response{id})
	})
}
