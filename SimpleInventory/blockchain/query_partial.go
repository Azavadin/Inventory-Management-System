package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"strings"
)
// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryKey_partial(key string) (string, error) {

	// Prepare arguments
	var args []string
	key=strings.ToLower(key)
	args = append(args, "invoke")
	args = append(args, "partial")
	args = append(args, key)

	response, err := setup.client.Query(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	
	if err != nil {
		fmt.Println("ERROR at partial query")
		fmt.Println(err)
		//return "", fmt.Errorf("failed to query partial composite key: %v", err)
	}
	//fmt.Println(string(response.Payload))
	return string(response.Payload), nil
}
