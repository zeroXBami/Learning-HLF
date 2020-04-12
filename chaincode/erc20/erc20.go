package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type ERC20Chaincode struct {
}

// Transfer Event
type TransferEvent struct {
	Sender    string `json:"sender`
	Recipient string `json:"recipient`
	Amount    int    `json:"amount"`
}

// Token information
type ERC20Metadata struct {
	Name        string `json:"name`
	Symbol      string `json:"symbol`
	Owner       string `json:"owner`
	TotalSupply uint64 `json:"totalSupply"`
}

func main() {
	err := shim.Start(new(ERC20Chaincode))
	if err != nil {
		fmt.Printf("Error starting ERC20 chaincode: %s", err)
	}
}

// Init function, called when chaincode installed on network
// params - tokenName, symbol, owner, amount
func (cc *ERC20Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Unit called with params", params)
	if len(params) != 4 {
		return shim.Error("incorect number of params")
	}

	tokenName, symbol, owner, amount := params[0], params[1], params[2], params[3]

	amountUint, err := strconv.ParseUint(string(amount), 10, 64)

	if err != nil {
		return shim.Error("Amount must be a number and can not be negative")
	}

	if len(tokenName) == 0 || len(symbol) == 0 || len(owner) == 0 {
		return shim.Error("tokenName, symbol or owner can not be empty")
	}

	erc20 := &ERC20Metadata{Name: tokenName, Symbol: symbol, Owner: owner, TotalSupply: amountUint}

	erc20Bytes, err := json.Marshal(erc20)
	if err != nil {
		return shim.Error(" Failed to Marshal erc20")
	}

	err = stub.PutState(tokenName, erc20Bytes)
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}

	err = stub.PutState(owner, []byte(amount))
	if err != nil {
		return shim.Error(" Failed to putstate, error: " + err.Error())
	}

	return shim.Success(nil)
}

func (cc *ERC20Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + fcn + " with params " + params[0])

	if fcn == "totalSupply" {
		return cc.totalSupply(stub, params)
	} else if fcn == "balanceOf" {
		return cc.balanceOf(stub, params)
	} else if fcn == "transfer" {
		return cc.transfer(stub, params)
	}

	fmt.Println("invoke did not find func: " + fcn) //error
	return shim.Error("Received unknown function invocation" + fcn)
}

func (cc *ERC20Chaincode) totalSupply(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 1 {
		return shim.Error("Incorrect number of params")

	}

	tokenName := params[0]
	erc20 := ERC20Metadata{}
	erc20Bytes, err := stub.GetState(tokenName)

	if err != nil || erc20Bytes == nil {
		return shim.Error("failed to get bytes data, error: " + err.Error())
	}

	err = json.Unmarshal(erc20Bytes, &erc20)
	if err != nil {
		return shim.Error("failed to get bytes data, error: " + err.Error())
	}

	totalSupplyBytes, err := json.Marshal(erc20.TotalSupply)
	if err != nil {
		return shim.Error("failed to Marshal erc20, error: " + err.Error())
	}

	fmt.Println(tokenName + "'s totla supply is " + string(totalSupplyBytes))

	return shim.Success(totalSupplyBytes)
}

func (cc *ERC20Chaincode) balanceOf(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters")
	}

	address := params[0]

	amountBytes, err := stub.GetState(address)
	if err != nil {
		return shim.Error("failed to get states, error :" + err.Error())
	}

	fmt.Println(address + "'s balance is " + string(amountBytes))
	if amountBytes == nil {
		return shim.Success([]byte("0"))
	}
	return shim.Success(amountBytes)
}

func (cc *ERC20Chaincode) transfer(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 3 {
		return shim.Error("Incorrect number of params")
	}
	callerAddress, recipientAddress, transferAmount := params[0], params[1], params[2]
	transferAmountInt, err := strconv.Atoi(transferAmount)
	if err != nil {
		return shim.Error("transfer Amount must be integer")
	}
	if transferAmountInt <= 0 {
		return shim.Error("transfer amount must be positive")
	}
	callerAmount, err := stub.GetState(callerAddress)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}
	callerAmountInt, err := strconv.Atoi(string(callerAmount))
	if err != nil {
		return shim.Error("caller amount must be integer or must be higher than transferAmount: " + err.Error())
	}

	recipientAmount, err := stub.GetState(recipientAddress)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}
	if recipientAmount == nil {
		recipientAmount = []byte("0")
	}

	recipientAmountInt, err := strconv.Atoi(string(recipientAmount))
	if err != nil {
		return shim.Error("recipe amount must be integer: " + err.Error())
	}

	callerResultAmount := callerAmountInt - transferAmountInt
	recipientResultAmount := recipientAmountInt + transferAmountInt
	// save each amounts
	err = stub.PutState(callerAddress, []byte(strconv.Itoa(callerResultAmount)))
	if err != nil {
		return shim.Error("failed to PutState of caller, error: " + err.Error())
	}
	err = stub.PutState(recipientAddress, []byte(strconv.Itoa(recipientResultAmount)))
	if err != nil {
		return shim.Error("failed to PutState of caller, error: " + err.Error())
	}

	transferEvent := TransferEvent{Sender: callerAddress, Recipient: recipientAddress, Amount: transferAmountInt}
	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return shim.Error("failed to Marshal transferEvent, error: " + err.Error())
	}
	err = stub.SetEvent("transferEvent", transferEventBytes)
	if err != nil {
		return shim.Error("failed to setEvent, error: " + err.Error())
	}

	fmt.Println(callerAddress + "send" + transferAmount + "to" + recipientAddress)
	return shim.Success([]byte("transfer Success"))
}

func (cc *ERC20Chaincode) allowance(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) approve(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) transferFrom(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) increaseAllowance(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) decreaseAllowance(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) mint(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) burn(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	return shim.Success(nil)
}
