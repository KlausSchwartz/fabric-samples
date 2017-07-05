/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Trade structure, with its properties.  Structure tags are used by encoding/json library
type Trade struct {
	Date   string   `json:"date"`
	Values []string `json:"values"`
}

/*
 * The Init method is called when the Smart Contract "etrade" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "etrade"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryTrade" {
		return s.queryTrade(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createTrade" {
		return s.createTrade(APIstub, args)
	} else if function == "queryAllTrades" {
		return s.queryAllTrades(APIstub)
	} else if function == "changeTradeValues" {
		return s.changeTradeValues(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryTrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tradeAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(tradeAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	trades := []Trade{
		Trade{Date: "01.01.2017", Values: []string{"3.148", "3.222"}},
		Trade{Date: "02.01.2017", Values: []string{"4.547", "3.333"}},
		Trade{Date: "03.01.2017", Values: []string{"1.478", "3.444"}},
		Trade{Date: "04.01.2017", Values: []string{"8.465", "3.555"}},
		Trade{Date: "05.01.2017", Values: []string{"4.112", "3.666"}},
	}

	i := 0
	for i < len(trades) {
		fmt.Println("i is ", i)
		tradeAsBytes, _ := json.Marshal(trades[i])
		APIstub.PutState("TRADE"+strconv.Itoa(i), tradeAsBytes)
		fmt.Println("Added", trades[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createTrade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var trade = Trade{Date: args[1], Values: []string{args[2]}}

	tradeAsBytes, _ := json.Marshal(trade)
	APIstub.PutState(args[0], tradeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllTrades(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "TRADE0"
	endKey := "TRADE99999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTrades:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeTradeValues(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	tradeAsBytes, _ := APIstub.GetState(args[0])
	trade := Trade{}

	json.Unmarshal(tradeAsBytes, &trade)
	// trade.Values = args[1]
	trade.Values = []string{args[1]}

	tradeAsBytes, _ = json.Marshal(trade)
	APIstub.PutState(args[0], tradeAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
