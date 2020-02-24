package api

import (
	add2 "github.com/theonlyrob/vercer/webserver/services/predicate/api/add"
	delete2 "github.com/theonlyrob/vercer/webserver/services/predicate/api/delete"
	get2 "github.com/theonlyrob/vercer/webserver/services/predicate/api/get"
	search2 "github.com/theonlyrob/vercer/webserver/services/predicate/api/search"
	update2 "github.com/theonlyrob/vercer/webserver/services/predicate/api/update"
	"net/http"
)

func Register(mux *http.ServeMux) {
	add2.Register(mux)
	delete2.Register(mux)
	get2.Register(mux)
	search2.Register(mux)
	update2.Register(mux)
}
