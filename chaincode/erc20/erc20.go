package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/tools/protolator"
)

type ERC20Chaincode struct {
}

// Transfer Event
type TransferEvent struct {
	Sender    string `json:"sender`
	Recipient string `json:"recipient`
	Amount    int    `json:"amount"`
}

type PrivateData struct {
	Id    string `json:"id`
	Owner string `json:"owner`
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
		return shim.Error("incorect number of params")
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

	err = stub.PutState("Intage", []byte(initSupply))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}
	err = stub.PutState(TotalSupplyKey, []byte(initSupply))
	if err != nil {
		return shim.Error("Failed to putstate, error: " + err.Error())
	}

	return shim.Success(nil)
}

func (cc *ERC20Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + fcn + " with params " + params[0])

	if fcn == "totalSupply" {
		return cc.totalSupply(stub)
	} else if fcn == "balanceOf" {
		return cc.balanceOf(stub, params)
	} else if fcn == "uploadData" {
		return cc.uploadData(stub, params)
	} else if fcn == "getTokenInfor" {
		return cc.getTokenInfor(stub)
	} else if fcn == "requestViewData" {
		return cc.requestViewData(stub, params)
	}

	fmt.Println("invoke did not find func: " + fcn) //error
	return shim.Error("Received unknown function invocation" + fcn)
}

func (cc *ERC20Chaincode) uploadData(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 2 {
		return shim.Error("Incorrect number of params")
	}

	id, owner := params[0], params[1]
	if len(id) == 0 || len(owner) == 0 {
		return shim.Error("Incorrect params value")
	}

	isExistID, err := stub.GetState(id)

	if isExistID != nil {
		return shim.Error("Id already existed")
	}
	newPrivateData := PrivateData{id, owner}
	privDataAsBytes, err := json.Marshal(newPrivateData)
	if err != nil {
		return shim.Error("Can not Marshal data" + err.Error())
	}
	err = stub.PutState(newPrivateData.Id, privDataAsBytes)
	if err != nil {
		return shim.Error("Can not put state new privateData" + err.Error())
	}

	return cc.mint(stub, owner, 1)
}

func (cc *ERC20Chaincode) requestViewData(stub shim.ChaincodeStubInterface, params []string) pb.Response {

	if len(params) != 2 {
		return shim.Error(" Incorrect number of params")
	}

	requester, dataID := params[0], params[1]

	priVdata := PrivateData{}
	dataAsBytes, err := stub.GetState(dataID)
	if err != nil {
		return shim.Error("Can not find dataId" + err.Error())
	}
	err = json.Unmarshal(dataAsBytes, &priVdata)
	if err != nil {
		return shim.Error("failed to get bytes data, error: " + err.Error())
	}

	if err != nil {
		return shim.Error("failed to Marshal Private Data, error: " + err.Error())
	}

	return cc.transfer(stub, []string{requester, priVdata.Owner, "1"})

}

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

	fmt.Println("Total supply is " + string(totalTokenSupply))
	tokenInfor := string(tokenName) + ", " + string(tokenSymbol) + ", " + string(tokenPublisher) + ", " + string(totalTokenSupply)
	return shim.Success([]byte(tokenInfor))
}

func (cc *ERC20Chaincode) totalSupply(stub shim.ChaincodeStubInterface) pb.Response {
	totalTokenSupply, err := stub.GetState(TotalSupplyKey)
	if err != nil {
		return shim.Error("Can not get state total supply " + err.Error())
	}

	fmt.Println("Total supply is " + string(totalTokenSupply))
	return shim.Success(totalTokenSupply)
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
		return shim.Error("transfer Amount must be integer" + err.Error())
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
	return shim.Success([]byte("transfer Success" + transferAmount))
}

func (cc *ERC20Chaincode) mint(stub shim.ChaincodeStubInterface, to string, amount int) pb.Response {

	if len(to) == 0 {
		return shim.Error("Invalid type of to")
	}

	currentBalance, err := stub.GetState(to)
	if err != nil {
		stub.PutState(to, []byte(strconv.Itoa(amount)))
		return shim.Success([]byte("Success Mint Token"))
	}

	currentBalanceInt, err := strconv.Atoi(string(currentBalance))
	if err != nil {
		currentBalanceInt = 0
	}
	newBalanceInt := currentBalanceInt + 1
	err = stub.PutState(to, []byte(strconv.Itoa(newBalanceInt)))
	if err != nil {
		return shim.Error("Failed to put new state balance")
	}
	currentTotalSupply, err := stub.GetState(TotalSupplyKey)
	newTotalSupply, err := strconv.Atoi(string(currentTotalSupply))
	totalSupply := newTotalSupply + 1
	err = stub.PutState(TotalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	return shim.Success([]byte("Success upload data and mint Token"))
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

func (cc *ERC20Chaincode) checkMsgSender(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("\nBegin*** getCreator \n")
	creator, err := stub.GetCreator()
	if err != nil {
		fmt.Printf("GetCreator Error")
		return shim.Error(err.Error())
	}

	si := &msp.SerializedIdentity{}
	err2 := proto.Unmarshal(creator, si)
	if err2 != nil {
		fmt.Printf("Proto Unmarshal Error")
		return shim.Error(err2.Error())
	}
	buf := &bytes.Buffer{}
	protolator.DeepMarshalJSON(buf, si)
	fmt.Printf("End*** getCreator \n")
	fmt.Printf(string(buf.Bytes()))

	return shim.Success([]byte(buf.Bytes()))
}
