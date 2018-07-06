package controllers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/SimpleInventory/Parser"
	"strconv"
	"strings"
)



	type Data struct {
		Name string
		Cat string
		Price string
		Amount string
		Success       bool
		Response      bool
		Message string
		Partial_key bool
		Comment string
		List_data []Parser.Inventory
	}

func (app *Application) FindItemHandler(w http.ResponseWriter, r *http.Request) {
	data := &Data{
		Name: "",
		Cat: "",
		Price: "",
		Amount: "",
		Message: "",
		Comment:"",
		Partial_key: false,
		Success: false,
		Response: false,
		List_data: []Parser.Inventory{},
	}

	//for composite keys
	//var stub shim.ChaincodeStubInterface
	if r.FormValue("submitted") == "true" {
		name := r.FormValue("name")
		cat := r.FormValue("cat_num")
		var info Parser.Inventory
		//name must be entered
		if name ==""{
			data.Message="Please enter the item name"
			data.Success=false
			data.Response=false
		}
		//full key, stock amount, buy and used transaction will all be queried
		if name!="" && cat!=""{
			key := "new:"+string(name)+":"+string(cat)
			Value, err := app.Fabric.QueryKey(string(key))
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Unable to query the blockchain!!!", 500)
			}
			if Value!=""{
				jsonstr:=[]byte(Value)
				item:=info.ToJson(jsonstr)
				data.Name=item.Name
				data.Cat=item.Catalog_num
				data.Price=FloatToString(float64(item.Price))
				data.Amount=FloatToString(float64(item.Amount))
				data.Comment=item.Comment
				data.Success=true
				data.Response=true

				fmt.Println(item.Name)
				fmt.Println(item.Catalog_num)
				fmt.Println(item.Amount)
			}
		}
		// partial key, stock information related to possible keys will be displayed
		var list_item []string
		if name!="" && cat==""{
			key := "new:"+string(name)
			Value, err := app.Fabric.QueryKey_partial(key)
			if err != nil {
				http.Error(w, "Unable to query the blockchain!!!", 500)
			}
			if len(Value)!=0{
				list_item=strings.Split(Value,"|")
			}
			//fmt.Println(list_item)
			data.Partial_key=true
		}
		var ex Parser.Inventory
		for _, item := range list_item{
			fmt.Println(item)
			err := json.Unmarshal([]byte(item), &ex)
			if err != nil {
				fmt.Println(err)
				jsonResp := "{\"Error\":\"Failed to decode JSON of: " + string(item) + "\"}"
				fmt.Println(jsonResp)
			}
			data.List_data=append(data.List_data,ex)
		}

	}
	renderTemplate(w, r, "findItem.html", data)
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}
