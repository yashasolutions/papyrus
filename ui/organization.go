package ui

import (
	"math"
	"mime"
	"net/http"

	"github.com/gophergala2016/papyrus/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gopkg.in/mgo.v2/bson"
)

func ServeOrganizationList(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	acc := ctx.Account

	orgs, err := data.ListOraganizationsOwner(acc.ID, 0, math.MaxInt32)
	catch(r, err)

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeOrganizationList, struct {
		Context       *Context
		Organizations []data.Organization
	}{
		Context:       ctx,
		Organizations: orgs,
	})
}

func ServeOrganizationNew(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeOrganizationNew, struct {
		Context *Context
	}{
		Context: ctx,
	})
}

func HandleOrganizationCreate(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err := r.ParseForm()
	catch(r, err)

	body := struct {
		Name string `schema:"name"`
	}{}

	err = schema.NewDecoder().Decode(&body, r.PostForm)
	catch(r, err)

	switch {
	case body.Name == "":
		RedirectBack(w, r)
		return
	}

	org := data.Organization{
		Name:      body.Name,
		OwnerID:   ctx.Account.ID,
		CreatorID: ctx.Account.ID,
	}
	err = org.Put()
	catch(r, err)

	http.Redirect(w, r, "/organizations/"+org.ID.Hex(), http.StatusSeeOther)
}

func ServeOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	if ctx.Account == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	if org.OwnerID != ctx.Account.ID {
		ServeForbidden(w, r)
		return
	}

	prjs, err := data.ListProjectsOrganization(org.ID, 0, math.MaxInt32)
	catch(r, err)

	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	ServeHTMLTemplate(w, r, tplServeOrganization, struct {
		Context      *Context
		Organization *data.Organization
		Projects     []data.Project
	}{
		Context:      ctx,
		Organization: org,
		Projects:     prjs,
	})
}

func init() {
	Router.NewRoute().Methods("GET").Path("/organizations").HandlerFunc(ServeOrganizationList)
	Router.NewRoute().Methods("GET").Path("/organizations/new").HandlerFunc(ServeOrganizationNew)
	Router.NewRoute().Methods("POST").Path("/organizations/new").HandlerFunc(HandleOrganizationCreate)
	Router.NewRoute().Methods("GET").Path("/organizations/{id}").HandlerFunc(ServeOrganization)
}
