package controllers

import (

	"net/http"
	"fmt"
	"github.com/SimpleInventory/Parser"
	"encoding/json"
	"strings"
	"time"
	"strconv"

)

// pull all buy and request transactions by given item name and catalog number
func (app *Application) TransHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct{
		Name string
		Cat string
		Amount string
		Request string
		Message string
		TransactionId string
		Success       bool
		Response      bool
		Totalcost float32
		Totalunit float32
		Totalunit_used float32
		List_buy []Parser.Inventory
		List_used []Parser.Inventory

	}{
		Name: "",
		Cat: "",
		Amount: "",
		Request:"",
		Message: "",
		TransactionId: "",
		Totalcost:0,
		Totalunit:0,
		Totalunit_used:0,
		Success:false,
		Response:false,
		List_buy: []Parser.Inventory{},
		List_used: []Parser.Inventory{},
	}
	days:= r.FormValue("days")
	if r.FormValue("submitted") == "true" && days =="" {
		name := r.FormValue("name")
		cat := r.FormValue("cat_num")
		key := "new:"+string(name)+":"+string(cat)
		//check if key exists, if yes, information will be updated otherwise error message will be returned 
		result_state :=checkState(app,key,w,r)
		data.Name=name
		data.Cat=cat
		if name!="" && cat!=""{
			if result_state !="" {
				name := r.FormValue("name")
				cat := r.FormValue("cat_num")
				key_buy := "buy:"+string(name)+":"+string(cat)
				key_used := "used:"+string(name)+":"+string(cat)
				data.List_buy=querytrans(app,key_buy,w,r)
				data.List_used=querytrans(app,key_used,w,r)
				Totalcost:=totalcost(data.List_buy)
				Totalunit_used:=totalcost(data.List_used)
				data.Totalcost=Totalcost[0]
				data.Totalunit=Totalcost[1]
				data.Totalunit_used=Totalunit_used[1]
				data.Success=true
				data.Response=true
				if len(data.List_buy)==0{
					data.Success=false
					data.Response=false
					data.Message="No purchase transaction found"
				}
				if len(data.List_used)==0{
					data.Success=false
					data.Response=false
					data.Message="No request transaction found"
				}
			}else{
				data.Success=false
				data.Response=false
				data.Message="No record found"
			}
		}else{
			data.Success=false
			data.Response=false
			data.Message="Please enter the name and catalog number"
		}
	}else if r.FormValue("submitted") == "true" && days !="" {
		name := r.FormValue("name")
		cat := r.FormValue("cat_num")
		span,_:= strconv.Atoi(days)
		key := "new:"+string(name)+":"+string(cat)
		//check if key exists, if yes, information will be updated otherwise error message will be returned 
		result_state :=checkState(app,key,w,r)
		data.Name=name
		data.Cat=cat
		if name!="" && cat!=""{
			if result_state !="" {
				name := r.FormValue("name")
				cat := r.FormValue("cat_num")
				key_buy := "buy:"+string(name)+":"+string(cat)
				key_used := "used:"+string(name)+":"+string(cat)
				data.List_buy=querytrans_time(app,key_buy,w,r,span)
				data.List_used=querytrans_time(app,key_used,w,r,span)
				Totalcost:=totalcost(data.List_buy)
				Totalunit_used:=totalcost(data.List_used)
				data.Totalcost=Totalcost[0]
				data.Totalunit=Totalcost[1]
				data.Totalunit_used=Totalunit_used[1]
				data.Success=true
				data.Response=true
				if len(data.List_buy)==0{
					data.Success=false
					data.Response=false
					data.Message="No purchase transaction found"
				}
				if len(data.List_used)==0{
					data.Success=false
					data.Response=false
					data.Message="No request transaction found"
				}
			}else{
				data.Success=false
				data.Response=false
				data.Message="No record found"
			}
		}else{
			data.Success=false
			data.Response=false
			data.Message="Please enter the name and catalog number"
		}
	}
	renderTemplate(w, r, "trans.html", data)
}
func querytrans(app *Application, key string,w http.ResponseWriter, r *http.Request) []Parser.Inventory {
	var result []Parser.Inventory
	var list_item []string
	Value, err := app.Fabric.QueryKey_partial(key)
	if err != nil {
		http.Error(w, "Unable to query the blockchain!!!", 500)
	}
	fmt.Println("Query trans")
	fmt.Println(Value)
	if len(Value)!=0{
		list_item=strings.Split(Value,"|")
	}
	var ex Parser.Inventory
	if len(list_item)>=1{
		for _, item := range list_item{
			err := json.Unmarshal([]byte(item), &ex)
			if err != nil {
				fmt.Println(err)
				jsonResp := "{\"Error\":\"Failed to decode JSON of: " + string(item) + "\"}"
				fmt.Println(jsonResp)
			}
			fmt.Println(item)
			result=append(result,ex)
		}
	}
	return result
}


func querytrans_time(app *Application, key string,w http.ResponseWriter, r *http.Request,span int) []Parser.Inventory {
	t:=time.Now()
	sec:=int64(span)*24*60*60
	cutoff:=t.Unix()-sec
	var result []Parser.Inventory
	var list_item []string
	Value, err := app.Fabric.QueryKey_partial(key)
	if err != nil {
		http.Error(w, "Unable to query the blockchain!!!", 500)
	}
	fmt.Println("Query trans by time span")
	fmt.Println(Value)
	if len(Value)!=0{
		list_item=strings.Split(Value,"|")
	}
	var ex Parser.Inventory
	if len(list_item)>=1{
		for _, item := range list_item{
			err := json.Unmarshal([]byte(item), &ex)
			if err != nil {
				fmt.Println(err)
				jsonResp := "{\"Error\":\"Failed to decode JSON of: " + string(item) + "\"}"
				fmt.Println(jsonResp)
			}
			if ex.Unix_time >= cutoff{
				result=append(result,ex)
			}
		}
	}
	return result
}




func totalcost(data []Parser.Inventory) []float32 {
	total:=float32(0)
	cost:=float32(0)
	units:=float32(0)
	var res []float32
	for _,item:= range data{
		cost=item.Amount*item.Price
		total=total+cost
		units=units+item.Amount
	}
	res=append(res,total)
	res=append(res,units)
	return res
}
