package controllers

import (

	"net/http"
	"fmt"
	"github.com/SimpleInventory/Parser"
	"strconv"
	"math/rand"
	"time"
	"strings"

)


// handle item request and update the ledger
func (app *Application) RequestItemHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct{
		Name string
		Cat string
		Amount string
		Request string
		Message string
		TransactionId string
		Success       bool
		Response      bool

	}{
		Name: "",
		Cat: "",
		Amount: "",
		Request:"",
		Message: "",
		TransactionId: "",
		Success:false,
		Response:false,
	}
	now:=time.Now()
	t:=strings.Split(now.String()," ")
	day:=t[0]
	if r.FormValue("submitted") == "true" {
		name := r.FormValue("name")
		cat := r.FormValue("cat_num")
		amount:=r.FormValue("amount")
		request:=r.FormValue("request")
		amount_f,_:= strconv.ParseFloat(amount,32)

		key := "new:"+string(name)+":"+string(cat)

		//check if key exists, if yes, information will be updated otherwise error message will be returned 
		result_state :=checkState(app,key,w,r)
		var j1 Parser.Inventory
		if name!="" && cat!="" && request!="" && amount!=""{
			if result_state !="" {
				//update the amount information and add to ledger
				var key_value []string
				key_value = append(key_value,key)
				fmt.Println("key exists, updating the ledger")
				jsonstr:=[]byte(result_state)
				item:=j1.ToJson(jsonstr)
				item.Amount= (item.Amount-float32(amount_f))
				fmt.Println(item.Amount)

				item.Request_by=request
				item_str:=j1.Tostring(&item)
				key_value = append(key_value,item_str)
				txid, err := app.Fabric.Invokekey(key_value)
				data.Message="ledger has been successfully updated"
				if err != nil {
					http.Error(w, "Unable to invoke key in the blockchain", 500)
				}
				data.TransactionId=txid
				data.Success=true
				data.Response=true
				data.Message="ledger has been updated but request tranaction has not been recoreded "
			}else{
				data.Message="No record found"
			}
			//request transaction
			if data.Success && data.Response{
				n1:=rand.Intn(100)
				n2:=rand.Intn(100)
				key="used:"+string(name)+":"+string(cat)+":"+string(n1)+string(n2)
				key=strings.ToLower(key)
				j1.Name=name
				j1.Catalog_num=cat
				j1.Amount=float32(amount_f)
				j1.Request_by=request
				j1.Purchase_date=day
				s:=j1.Tostring(&j1)
				var key_value []string
				key_value = append(key_value,key)
				key_value = append(key_value,s)
				_, err := app.Fabric.Invokekey(key_value)
				if err != nil {
					http.Error(w, "Unable to invoke key in the blockchain", 500)
				}
				data.Message="ledger has been updated and request tranaction has been recoreded "
			}
		}else{
			data.Message="Please enter the name and catalog number"
		}

	}
	renderTemplate(w, r, "requestItem.html", data)
}
