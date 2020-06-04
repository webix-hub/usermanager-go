package api

import (
	"net/http"
	"sort"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"

	"webix/usermanager/data"
)

var format = render.New()

func Routes(r chi.Router, dao *data.DAO, path, server string) {
	usersAPI(r, dao, path, server)
	rulesAPI(r, dao)
	metaAPI(r, dao)
	rolesAPI(r, dao)
	credentialsAPI(r, dao)
	logsAPI(r, dao)
}

type Response struct {
	ID int `json:"id"`
}

func errorResult(w http.ResponseWriter, err error) {
	format.Text(w, 500, err.Error())
}

func intersect(was, now []int) ([]int, []int) {
	sort.Ints(was)
	sort.Ints(now)

	added := make([]int, 0, len(now))
	removed := make([]int, 0, len(was))

	p1 := 0
	p2 := 0
	for p1 < len(was) && p2 < len(now) {
		if was[p1] == now[p2] {
			p1++
			p2++
		} else if was[p1] > now[p2] {
			added = append(added, now[p2])
			p2++
		} else {
			removed = append(removed, was[p1])
			p1++
		}
	}

	if p1 != len(was) {
		removed = append(removed, was[p1:]...)
	}
	if p2 != len(now) {
		added = append(added, now[p2:]...)
	}

	return added, removed
}
