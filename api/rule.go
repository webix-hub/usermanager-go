package api

import (
	"net/http"
	"webix/usermanager/data"
	"webix/usermanager/rules"

	"github.com/go-chi/chi"
)

func rulesAPI(r chi.Router, dao *data.DAO) {
	r.Get("/rules", func(w http.ResponseWriter, r *http.Request) {
		format.JSON(w, 200, rules.GetDetails())
	})
}
