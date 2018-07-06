package controllers

import (
	"net/http"
)

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Hello string
	}{
		Hello: "",
	}
	if r.FormValue("submitted") == "true" {
		key:= r.FormValue("key")
		
		helloValue, err := app.Fabric.QueryKey(string(key))
		if err != nil {
			http.Error(w, "Unable to query the blockchain", 500)
		}
	data.Hello=helloValue
	}

	renderTemplate(w, r, "home.html", data)
}
