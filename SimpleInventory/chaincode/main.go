package main

import (
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/SimpleInventory/Parser"
	//"encoding/json"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleInventoryChaincode implementation of Chaincode
type SimpleInventoryChaincode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
func (t *SimpleInventoryChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//fmt.Println("########### SimpleInventoryChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call!!")
	}

	// initial the chain with ledger the key/value hello/world 
	err := stub.PutState("test", []byte("hello world"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	fmt.Println("###########Chaincode has been instantiated###########")
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *SimpleInventoryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//fmt.Println("########### SimpleInventoryChaincode Invoke ###########")
	fmt.Println("###########Chaincode is invoked ###########")
	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Check whether it is an invoke request
	if function != "invoke" {
		// In order to manage multiple type of request, we will check the first argument.
		// Here we have one possible argument: query (every query request will read in the ledger without modification)
		return shim.Error("Unknown function call~~~")
	}
	
	if args[0] == "query" {
		return t.query(stub, args)
	}
	// The update argument will add key/value in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}
	if args[0]  == "delete" {
		return t.delete(stub, args)
		//return t.delete(stub, args)
	}
	if args[0] == "partial" {
		return t.query_partial(stub, args)
	}
	
	//return shim.Error("Unknown function call")
	

	/*

	if args[0] == "querypartial" {
		return t.query_partial(stub, args)
	}

	// The del argument will delete key/value in ledger
	if args[0] == "del" {
		return t.delete(stub, args)
	}
	*/

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument!")
}


// query ledger with key
func (t *SimpleInventoryChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Println("########### SimpleInventoryChaincode query key ###########")
	fmt.Println("########### Query chain with key ###########")
	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	key:=t.generatecompositekey(stub,args[1])
	
	// Get the state of the value matching the key hello in the ledger
	state, err := stub.GetState(key)
	fmt.Println(string(state))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get state: %s", args[1]))
	}

	// Return this value in response
	return shim.Success(state)
}
// query ledger with partial composite key
func (t *SimpleInventoryChaincode) query_partial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Println("########### SimpleInventoryChaincode query partial ###########")
	fmt.Println("########### Query chain with Partial key ###########")
	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	//key:=t.generatecompositekey(stub,args[1])
	keys:=strings.Split(args[1],":")
	// Get the state of the value matching the key hello in the ledger
	statesIterator,err:=stub.GetStateByPartialCompositeKey("type~name~cat", keys)
	if err != nil {
		fmt.Println(err.Error())
		return shim.Error(fmt.Sprintf("Failed to get state by partial key: %s", args[1]))
	}
	defer statesIterator.Close()
	state_list:=""
	//var item Parser.Inventory
	var i int
	for i = 0; statesIterator.HasNext(); i++ {
		response, err := statesIterator.Next()
		if err != nil {
			fmt.Println(err.Error())
			//return shim.Error(err.Error())
		}
		state, err := stub.GetState(response.Key)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + response.Key + "\"}"
			return shim.Error(jsonResp)
		} else if state == nil {
			jsonResp := "{\"Error\":\"Marble does not exist: " + response.Key  + "\"}"
			return shim.Error(jsonResp)
		}
		fmt.Println(string(state))
		//state_list=append(state_list,state)
		if len(state_list)<1{
			state_list=string(state)
		}else{
			state_list=state_list+"|"+string(state)
		}
		/*
		err = json.Unmarshal([]byte(state), &item)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to decode JSON of: " + response.Key + "\"}"
			return shim.Error(jsonResp)
		}*/
		
	}
	// Return this value in response
	return shim.Success([]byte(state_list))
}
//del key/value from ledger
func (t *SimpleInventoryChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### SimpleInventoryChaincode delete key ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}
	key:=t.generatecompositekey(stub,args[1])
	// Get the state of the value matching the key hello in the ledger
	err := stub.DelState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to delete asset: %s", args[1]))
	}

	// Return this value in response
	return shim.Success(nil)
}



// invoke
// Every functions that read and write in the ledger will be here
func (t *SimpleInventoryChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//fmt.Println("########### SimpleInventoryChaincode invoke ###########")
	fmt.Println("########### Adding item to the chain ###########")
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
	if  len(args) == 3 {
		key:=t.generatecompositekey(stub,args[1])
		// Write the new value in the ledger
		err := stub.PutState(key, []byte(args[2]))
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to update state: %s", args[1]))
		}

		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
		err = stub.SetEvent("eventInvoke", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return this value in response
		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}
//generate composite keys
func (t *SimpleInventoryChaincode) generatecompositekey(stub shim.ChaincodeStubInterface, str string) string {
	keys:=strings.Split(str,":")
	compKey := "type~name~cat"
	key, err := stub.CreateCompositeKey(compKey, keys)
	if err != nil {
		fmt.Errorf("Error creating composite key")
		return err.Error()
	}
	return key
}



func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(SimpleInventoryChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
