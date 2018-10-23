// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario
 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract example simple Chaincode implementation
type SmartContract struct {
}

/* Define Property structure, with 4 properties.
Structure tags are used by encoding/json library
*/
type Property struct {
	PropertyName string `json:"propertyname"`
	Timestamp    string `json:"timestamp"`
	Location     string `json:"location"`
	Holder       string `json:"holder"`
}

/*
 * The Init method *
 called when the Smart Contract "demo-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "demo-chaincode"
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("Invoke initiated for " + function)

	// Route to the appropriate handler function to interact with the ledger
	if function == "queryProperty" {
		return s.queryProperty(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordProperty" {
		return s.recordProperty(APIstub, args)
	} else if function == "queryAllProperty" {
		return s.queryAllProperty(APIstub)
	} else if function == "changePropertyHolder" {
		return s.changePropertyHolder(APIstub, args)
	} else if function == "getHistoryForProperty" {
		return s.getHistoryForProperty(APIstub, args)
	} else if function == "queryPropertyByHolder" {
		return s.queryPropertyByHolder(APIstub, args)
	} else if function == "delete" {
		return s.delete(APIstub, args)
	}
	fmt.Println("Invoke failed for unknown function " + function)
	return shim.Error("Invalid Smart Contract function name")
}

/*
 * The recordProperty method *
This method takes in five arguments (attributes to be saved in the ledger).
*/
func (s *SmartContract) recordProperty(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	if len(args[0]) <= 0 && args[0] != "" {
		return shim.Error("Property ID must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Property Name must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("Location must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("Timestamp must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("Holder must be a non-empty string")
	}
	timeStamp := args[3]
	// timeStamp = string(int64(time.Now().UnixNano())) // 2038 year problem
	// var property = Property{PropertyName: args[1], Location: args[2], Timestamp: args[3], Holder: args[4]}
	var property = Property{PropertyName: args[1], Location: args[2], Timestamp: timeStamp, Holder: args[4]}
	propertyAsBytes, _ := json.Marshal(property)
	err := APIstub.PutState(args[0], propertyAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record property: %s", args[0]))
	}
	fmt.Println("End recordProperty")
	return shim.Success(nil)
}

/*
 * The queryProperty method *
Used to view the records of one particular Property
It takes one argument -- the key for the Property in question
*/
func (s *SmartContract) queryProperty(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if len(args[0]) <= 0 {
		return shim.Error("Property ID must be a non-empty string")
	}
	propertyAsBytes, _ := APIstub.GetState(args[0])
	if propertyAsBytes == nil {
		return shim.Error(fmt.Sprintf("Could not locate property: %s", args[0]))
	}
	fmt.Println("Property found: %s", args[0])
	return shim.Success(propertyAsBytes)
}

/*
 * The queryAllProperty method *
allows for assessing all the records added to the ledger(all properties)
This method does not take any arguments. Returns JSON string containing results.
*/
func (s *SmartContract) queryAllProperty(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "0"
	endKey := "999"
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	// Buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a coma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write it as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Println("- queryAllProperty:\n%s\n", buffer.String)
	return shim.Success(buffer.Bytes())
}

/*
 * The initLedger method *
Will add test data (5 properties)to the network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// manually insert entries into the ledger
	property := []Property{
		Property{PropertyName: "Taj Mahal", Location: "27.1750, 78.0422", Timestamp: "1485153091", Holder: "Vanguard"},
		Property{PropertyName: "Kremlin", Location: "55.7520, 37.6175", Timestamp: "1504054225", Holder: "Vettel"},
		Property{PropertyName: "One World Trade Centre", Location: "40.7127, 74.0134", Timestamp: "1504057825", Holder: "Hamilton"},
		Property{PropertyName: "Burj Khalifa", Location: "25.1972, 55.2744", Timestamp: "1493517025", Holder: "Schumacher"},
		Property{PropertyName: "Louvre Museum", Location: "48.8606, 2.3376", Timestamp: "1485153091", Holder: "Verstappen"},
	}
	i := 0
	for i < len(property) {
		fmt.Println("i is ", i)
		propertyAsBytes, _ := json.Marshal(property[i])
		APIstub.PutState(strconv.Itoa(i+1), propertyAsBytes)
		fmt.Println("Added ", property[i])
		i = i + 1
	}
	fmt.Println("End initLedger")
	return shim.Success(nil)
}

/*
 * The changePropertyHolder method *
The data in the world state can be updated with who has possession.
This function takes in 2 arguments, propertyID and new holder name.
*/
func (s *SmartContract) changePropertyHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 && args[0] != "" {
		return shim.Error("Property ID must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Holder must be a non-empty string")
	}
	propertyAsBytes, _ := APIstub.GetState(args[0])
	if propertyAsBytes == nil {
		return shim.Error(fmt.Sprintf("Could not locate property: %s", args[0]))
	}
	property := Property{}
	json.Unmarshal(propertyAsBytes, &property)
	newOwner := args[1]
	if property.Holder == newOwner {
		return shim.Error("Current and New Property Owner cannot be same")
	}

	// Check that the specified argument is a valid holder of property
	property.Holder = args[1]
	propertyAsBytes, _ = json.Marshal(property)
	err := APIstub.PutState(args[0], propertyAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change property holder: %s", args[0]))
	}
	fmt.Println("End changePropertyHolder")
	return shim.Success(nil)
}

/*
 * getHistoryForProperty function *
get history for a property
*/
func (s *SmartContract) getHistoryForProperty(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	// propertyAsBytes, _ := APIstub.GetState(args[0])
	// if propertyAsBytes == nil {
	// 	return shim.Error(fmt.Sprintf("Could not locate property: %s", args[0]))
	// }
	// fmt.Println("Property found: %s", args[0])
	resultsIterator, err := APIstub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	// buffer is a JSON array containing historic values for the property
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it is a delete operation on a given key, then we need to set the
		// corresponding value null. Else, we'll write the response.Value as-is
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Println("- getHistoryForProperty returning:\n%s\n", buffer.String)
	return shim.Success(buffer.Bytes())
}

/*
 *queryPropertyByHolder*
 Parametrised Rich Query
 This function queries for properties based on a pssed in holder
 This is an example of a parameterised query where the query logic is baked into the chaincode,
 and accepting a single query parameter (holder)
 Only available on state databases that support rich query (e.g. CouchDB)
*/
func (s *SmartContract) queryPropertyByHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	if len(args[0]) <= 0 {
		return shim.Error("Holder must be a non-empty string")
	}
	// holder := strings.ToLower(args[0])
	holder := args[0]
	// queryString := fmt.Sprintf("{\"selector\":{\"propertyname\":\"property\",\"holder\":\"%s\"}}", holder)
	queryString := fmt.Sprintf("{\"selector\":{\"holder\":\"%s\"}}", holder)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

/*
 *getQueryResultForQueryString*
Executes the passed in query string.
Result set is built and returned as a byte array containing JSON results.
*/
func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n:", queryString)
	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil

}

/*
 *constructQueryResponseFromIterator*
constructs a JSON array containing query results from a given result iterator
*/
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write it as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return &buffer, nil
}

/*
 *remove a property key/value pair from the state
 */
func (s *SmartContract) delete(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	propertyAsBytes, err := APIstub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get the state for " + key + "\"}"
		return shim.Error(jsonResp)
	} else if propertyAsBytes == nil {
		jsonResp = "{\"Error\":\"Property does not exist " + key + "\"}"
		return shim.Error(jsonResp)
	}
	property := Property{}
	json.Unmarshal(propertyAsBytes, &property)
	err = APIstub.DelState(key) //remove the property from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state: " + err.Error())
	}
	//maintain indexes if any
	//delete index entry to state
	return shim.Success(nil)
}

/*
 * main function *
calls the Start function
The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	} else {
		fmt.Println("SamrtContract started successfully")
	}
}
