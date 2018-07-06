package controllers

import (
	"net/http"
	"fmt"
	"github.com/SimpleInventory/Parser"
	"strconv"
	"time"
	"math/rand"
	"strings"
)

func (app *Application) NewItemHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
		Message string
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
		Message: "",
	}
	//add new item to ledger, if same item exist, the amount will be updated
	//transaction will be recored seperately
	if r.FormValue("submitted") == "true" {
		name := r.FormValue("name")
		cat := r.FormValue("cat_num")
		amount := r.FormValue("amount")
		amount_f,err:= strconv.ParseFloat(amount,32)
		price := r.FormValue("price")
		price_f, err := strconv.ParseFloat(price,32)
		if err != nil {
			fmt.Errorf("Error casting string to float")
		}
		staff:=r.FormValue("staff")
		supplier:=r.FormValue("supplier")
		comment:=r.FormValue("comment")
		now:=time.Now()
		t:=strings.Split(now.String()," ")
		day:=t[0]
		unix_time:=now.Unix()


		//convert to Json string
		j1 := &Parser.Inventory {
			Name : string(name),
			Catalog_num: string(cat),
			Amount: float32(amount_f),
			Price: float32(price_f),
			Purchase_date: string(day),
			Comment: comment,
			Supplier: supplier,
			Purchase_by: staff,
			Unix_time: unix_time,
		}

		sample := j1.Tostring(j1)
		s := string(sample)
		if name!="" && cat!="" && amount !="" && supplier!="" && staff!=""{
			key := "new:"+string(name)+":"+string(cat)
			key=strings.ToLower(key)
			//check if key exists, if yes, information will be updated otherwise be created
			result_state :=checkState(app,key,w,r)
			if result_state !="" {
				//update the amount information and add to ledger
				var key_value []string
				key_value = append(key_value,key)
				fmt.Println("key exists, appending to the key in the ledger")
				jsonstr:=[]byte(result_state)
				item:=j1.ToJson(jsonstr)
				item.Amount=item.Amount+float32(amount_f)
				item_str:=j1.Tostring(&item)
				key_value = append(key_value,item_str)


				txid, err := app.Fabric.Invokekey(key_value)
				if err != nil {
					http.Error(w, "Unable to invoke key in the blockchain", 500)
				}
				data.TransactionId = txid
				data.Success = true
				data.Response = true
				data.Message="ledger has been updated but purchase tranaction has not been recoreded "
			}else{
				var key_value []string
				key_value = append(key_value,key)
				key_value = append(key_value,s)
				txid, err := app.Fabric.Invokekey(key_value)
				if err != nil {
					http.Error(w, "Unable to invoke key in the blockchain", 500)
				}
				data.TransactionId = txid
				data.Success = true
				data.Response = true
				data.Message="ledger has been updated but purchase tranaction has not been recoreded "
			}
			//purchase tranaction/////
			if data.Success && data.Response{
				n1:=rand.Intn(100)
				n2:=rand.Intn(100)
				key="buy:"+string(name)+":"+string(cat)+":"+string(n1)+string(n2)
				key=strings.ToLower(key)
				var key_value []string
				key_value = append(key_value,key)
				key_value = append(key_value,s)
				_, err := app.Fabric.Invokekey(key_value)
				if err != nil {
					http.Error(w, "Unable to invoke key in the blockchain", 500)
				}
				data.Message="ledger has been updated and purchase tranaction has been recoreded "
			}
		}else{
			fmt.Println("Please enter product name, catalog number, amount, staff and supplier ")
			data.Message="Please enter product name, catalog number, amount, staff and supplier information"
		}
	}
	renderTemplate(w, r, "newItem.html", data)
}
//check if key is already in the ledger
func checkState (app *Application, key string,w http.ResponseWriter, r *http.Request) string {
	Value, err := app.Fabric.QueryKey(string(key))
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}
	return Value
}

