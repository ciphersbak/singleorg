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
	} else if function == "changePropertyHolder" {
		return s.changePropertyHolder(APIstub, args)
	} else if function == "queryAllProperty" {
		return s.queryAllProperty(APIstub)
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
	var property = Property{PropertyName: args[1], Location: args[2], Timestamp: args[3], Holder: args[4]}
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
Will add test data (4 properties)to our network
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// manually insert entries into the ledger
	property := []Property{
		Property{PropertyName: "21D Avalon Midtown West", Location: "67.0006, -70.5476", Timestamp: "1504054225", Holder: "Vettel"},
		Property{PropertyName: "314 Jacobs Creek", Location: "91.2395, -49.4594", Timestamp: "1504057825", Holder: "Hamilton"},
		Property{PropertyName: "Burj Khalifa", Location: "58.0148, 59.01391", Timestamp: "1493517025", Holder: "Schumacher"},
		Property{PropertyName: "Empire State Tower", Location: "153.0054, 12.6429", Timestamp: "1485153091", Holder: "Verstappen"},
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
	propertyAsBytes, _ := APIstub.GetState(args[0])
	if propertyAsBytes != nil {
		return shim.Error(fmt.Sprintf("Could not locate property: %s", args[0]))
	}
	property := Property{}
	json.Unmarshal(propertyAsBytes, &property)
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
