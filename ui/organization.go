package ui

import (
	"mime"
	"net/http"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func ServeOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	id := bson.ObjectIdHex(idStr)
	org, err := data.GetOraganization(id)
	catch(r, err)
	if org == nil {
		ServeNotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeOrganization, struct {
		Context *Context
	}{
		Context: ctx,
	})
}

func init() {
	Router.NewRoute().Methods("GET").Path("/organizations/{id}").HandlerFunc(ServeOrganization)
}
