package web

import (
	"fmt"
	"github.com/SimpleInventory/web/controllers"
	"net/http"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home.html", app.HomeHandler)
	//http.HandleFunc("/request.html", app.RequestHandler)
	http.HandleFunc("/newItem.html", app.NewItemHandler)
	http.HandleFunc("/findItem.html", app.FindItemHandler)
	http.HandleFunc("/requestItem.html", app.RequestItemHandler)
	http.HandleFunc("/trans.html", app.TransHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
