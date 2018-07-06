package controllers

import (
	"net/http"
	"fmt"
	"github.com/SimpleInventory/Parser"
	"strconv"
)

func (app *Application) RequestHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
	}
	if r.FormValue("submitted") == "true" {
		helloValue := r.FormValue("hello")
		name := r.FormValue("name")
		cat := r.FormValue("cat_num")
		amount := r.FormValue("amount")
		amount_f,err:= strconv.ParseFloat(amount,32)
		comments := r.FormValue("comments")
		purchase_day := r.FormValue("purday")
		supplier := r.FormValue("supplier")
		purchase_by := r.FormValue("purchaseby")
		if err != nil {
			fmt.Errorf("Error casting string to float")
		}

		//convert to Json string
		j1 := &Parser.Inventory {
			Name : string(name),
			Catalog_num: string(cat),
			Amount: float32(amount_f),
			Comment : string(comments),
			Purchase_date: purchase_day,
			Supplier : string(supplier),
			Purchase_by : string(purchase_by)}
		sample := j1.Tostring(j1)
		s := string(sample)
		key := string(name)+string(cat)
		fmt.Println(key)
		fmt.Println(s)
		//key and value array
		var key_value []string
		key_value = append(key_value,key)
		key_value = append(key_value,s)

		fmt.Println(helloValue)
		//check if key exists, if yes, information will be updated otherwise be created
		result_state :=checkState(app,key,w,r)
		if result_state !="" {
			fmt.Println("key exists, appending to the key in the ledger")
			jsonstr:=[]byte(result_state)
			item:=j1.ToJson(jsonstr)
			fmt.Println(item)

		}else{
			txid, err := app.Fabric.Invokekey(key_value)
			if err != nil {
				http.Error(w, "Unable to invoke key in the blockchain", 500)
			}
			data.TransactionId = txid
			data.Success = true
			data.Response = true
		}
	}
	renderTemplate(w, r, "request.html", data)
}
/*
//check if key is already in the ledger
func checkState (app *Application, key string,w http.ResponseWriter, r *http.Request) string {
	Value, err := app.Fabric.QueryKey(string(key))
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}
	return Value
}*/

