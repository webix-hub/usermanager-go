package api

import (
	"net/http"
	"webix/usermanager/data"

	"github.com/go-chi/chi"
)

func metaAPI(r chi.Router, dao *data.DAO) {
	r.Get("/meta", func(w http.ResponseWriter, r *http.Request) {
		meta, err := dao.Meta.GetLinks()
		if err != nil {
			errorResult(w, err)
			return
		}

		format.JSON(w, 200, meta)
	})
}
