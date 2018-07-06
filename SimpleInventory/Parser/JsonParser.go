package Parser
import (
	"encoding/json"
	"fmt"
	"log"
)

type Inventory struct{
	Name string
	Catalog_num string
	Amount float32
	Price float32
	Comment string
	Purchase_date string //list
	Supplier string
	Purchase_by string
	Request_by string
	Unix_time int64
	
}
func (parser *Inventory) logIfError(err error) {
        if err != nil {
                log.Fatal(err)
        }
}
func (parser *Inventory) Tostring(item *Inventory) string{
	str, _ := json.Marshal(item)
	return string(str)
}
func (parser *Inventory) ToJson(data []byte) Inventory{
	var ex Inventory
	err := json.Unmarshal(data, &ex)
    parser.logIfError(err)
    fmt.Printf("Auto Unmarshal: %+v \n", ex)
    return ex
}
