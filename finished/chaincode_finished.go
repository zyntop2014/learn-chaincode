/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Execting 6")
	}

	var usersArray []string

	var userone User
	userone.Name = args[0]
	userone.Password = args[1]
	balance, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}

	userone.Balance = balance

	b, err := json.Marshal(userone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userone")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	userone.Name = args[3]
	userone.Password = args[4]
	balance, err = strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding at 3 place")
	}

	userone.Balance = balance

	b, err = json.Marshal(userone)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for userone")
	}

	err = stub.PutState(args[3], b)
	if err != nil {
		return nil, err
	}

	usersArray = append(usersArray, args[0])
	usersArray = append(usersArray, args[3])

	b, err = json.Marshal(usersArray)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for usertwo")
	}

	err = stub.PutState("users", b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "list_users" {
		return t.listUsers(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) listUsers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("users")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for users}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
