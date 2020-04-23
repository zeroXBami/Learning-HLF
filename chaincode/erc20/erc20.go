package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
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

const SymbolKey = `tokenSymbol`
const NameKey = `tokenName`
const TotalSupplyKey = `totalSupply`
const PublisherKey = `publisher`

func main() {
	err := shim.Start(new(ERC20Chaincode))
	if err != nil {
		fmt.Printf("Error starting ERC20 chaincode: %s", err)
	}
}

// Init function, called when chaincode installed on network
// params - tokenName, symbol, publisher
func (cc *ERC20Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Unit called with params", params)
	if len(params) != 4 {
		return shim.Error("Incorect number of params")
	}

	tokenName, symbol, publisher, initSupply := params[0], params[1], params[2], params[3]

	if len(tokenName) == 0 || len(symbol) == 0 || len(publisher) == 0 {
		return shim.Error("tokenName, symbol or owner can not be empty")
	}

	err := stub.PutState(NameKey, []byte(tokenName))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}

	err = stub.PutState(SymbolKey, []byte(symbol))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}

	err = stub.PutState(PublisherKey, []byte(publisher))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}
	// Mint amount of token (initSupply) to Intage at deploy time
	err = stub.PutState("admin@intage.example.com", []byte(initSupply))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}
	err = stub.PutState(TotalSupplyKey, []byte(initSupply))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}

	return shim.Success(nil)
}

// Invoke chaincode function
func (cc *ERC20Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + fcn + " with params " + params[0])

	if fcn == "balanceOf" {
		return cc.balanceOf(stub, params)
	} else if fcn == "getTokenInfor" {
		return cc.getTokenInfor(stub)
	} else if fcn == "mint" {
		return cc.mint(stub, params)
	} else if fcn == "transfer" {
		return cc.transfer(stub, params)
	}
	return shim.Error("Received unknown function invocation" + fcn)
}

// Get token information as tokenName, symbol, publisher, totalSupply
func (cc *ERC20Chaincode) getTokenInfor(stub shim.ChaincodeStubInterface) pb.Response {
	tokenName, err := stub.GetState(NameKey)
	if err != nil {
		return shim.Error("Can not get state token name " + err.Error())
	}

	tokenSymbol, err := stub.GetState(SymbolKey)
	if err != nil {
		return shim.Error("Can not get state token symbol " + err.Error())
	}

	tokenPublisher, err := stub.GetState(PublisherKey)
	if err != nil {
		return shim.Error("Can not get state token publisher " + err.Error())
	}

	totalTokenSupply, err := stub.GetState(TotalSupplyKey)
	if err != nil {
		return shim.Error("Can not get state total supply " + err.Error())
	}

	tokenInfor := string(tokenName) + ", " + string(tokenSymbol) + ", " + string(tokenPublisher) + ", " + string(totalTokenSupply)
	return shim.Success([]byte(tokenInfor))
}

// Get balance of specific address
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

// transfer token
func (cc *ERC20Chaincode) transfer(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 3 {
		return shim.Error("Incorrect number of params")
	}
	from, to, amount := params[0], params[1], params[2]
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return shim.Error("transfer amount must be integer" + err.Error())
	}
	if amountInt <= 0 {
		return shim.Error("transfer amount must be positive")
	}
	balanceOfSender, err := stub.GetState(from)
	if err != nil {
		return shim.Error("Failed to get balance of caller, error: " + err.Error())
	}
	balanceOfSenderInt, err := strconv.Atoi(string(balanceOfSender))
	if err != nil || (balanceOfSenderInt-amountInt) < 0 {
		return shim.Error("balance of sender must be integer or must be higher than transferAmount: " + err.Error())
	}

	balanceOfReceiver, err := stub.GetState(to)
	if err != nil {
		return shim.Error("Failed to get balance of receiver, error: " + err.Error())
	}
	if balanceOfReceiver == nil {
		balanceOfReceiver = []byte("0")
	}

	balanceOfReceiverInt, err := strconv.Atoi(string(balanceOfReceiver))
	if err != nil {
		return shim.Error("recipe amount must be integer: " + err.Error())
	}

	senderResultAmountInt := balanceOfSenderInt - amountInt
	receiverResultAmountInt := balanceOfReceiverInt + amountInt
	// save each amounts
	err = stub.PutState(from, []byte(strconv.Itoa(senderResultAmountInt)))
	if err != nil {
		return shim.Error("Failed to PutState of caller, error: " + err.Error())
	}
	err = stub.PutState(to, []byte(strconv.Itoa(receiverResultAmountInt)))
	if err != nil {
		return shim.Error("Failed to PutState of receiver, error: " + err.Error())
	}

	transferEvent := TransferEvent{Sender: from, Recipient: to, Amount: amountInt}
	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return shim.Error("Failed to Marshal transferEvent, error: " + err.Error())
	}
	err = stub.SetEvent("transferEvent", transferEventBytes)
	if err != nil {
		return shim.Error("Failed to emit transferEvent, error: " + err.Error())
	}

	return shim.Success([]byte("Transfer success" + amount))
}

// Intage mint token to other org
func (cc *ERC20Chaincode) mint(stub shim.ChaincodeStubInterface, params []string) pb.Response {

	if len(params) != 2 {
		return shim.Error("Incorect number of params")
	}
	x509, _ := cid.GetX509Certificate(stub)
	senderName := x509.Subject.CommonName
	mspid, _ := cid.GetMSPID(stub)
	if senderName != "admin" || mspid != "IntageMSP" {
		return shim.Error("Only Intage admin can mint token")
	}
	to, amount := params[0], params[1]
	amountInt, err := strconv.Atoi(string(amount))
	if err != nil || amountInt <= 0 {
		return shim.Error("Mint amount must be higher than zero" + err.Error())
	}
	currentBalance, err := stub.GetState(to)
	if err != nil {
		return shim.Error("Can not get balance of" + to + err.Error())
	}
	if currentBalance == nil {
		currentBalance = []byte("0")
	}

	currentBalanceInt, err := strconv.Atoi(string(currentBalance))
	if err != nil {
		return shim.Error("Balance must be integer" + err.Error())
	}
	newBalanceInt := currentBalanceInt + amountInt
	err = stub.PutState(to, []byte(strconv.Itoa(newBalanceInt)))
	if err != nil {
		return shim.Error("Failed to put new state balance")
	}
	currentTotalSupply, err := stub.GetState(TotalSupplyKey)
	newTotalSupply, err := strconv.Atoi(string(currentTotalSupply))
	totalSupply := newTotalSupply + amountInt
	err = stub.PutState(TotalSupplyKey, []byte(strconv.Itoa(totalSupply)))

	transferEvent := TransferEvent{Sender: "erc20", Recipient: to, Amount: amountInt}
	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return shim.Error("Failed to Marshal transferEvent, error: " + err.Error())
	}
	err = stub.SetEvent("transferEvent", transferEventBytes)
	if err != nil {
		return shim.Error("Failed to emit transferEvent, error: " + err.Error())
	}

	return shim.Success([]byte("Success mint token to" + to))
}
