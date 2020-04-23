package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type DataMarket struct {
}

// Share Event
type ShareEvent struct {
	From string `json:"from`
	To   string `json:"to`
	Id   string `json:"id"`
}

// Upload Event
type UploadEvent struct {
	Id         string `json:"id`
	Owner      string `json:"owner`
	UploadTime int64  `json:"uploadTme`
}

type PrivateData struct {
	Id       string `json:"id`
	Owner    string `json:"owner`
	DataDesc string `json:"desc`
}

const MarketOwner = `marketOwner`

func main() {
	err := shim.Start(new(DataMarket))
	if err != nil {
		fmt.Printf("Error starting DataMarket chaincode: %s", err)
	}
}

// Init function, called when chaincode installed, upgraded on network
func (cc *DataMarket) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := stub.PutState(MarketOwner, []byte("Intage"))
	if err != nil {
		return shim.Error("Chaincode init failed due to error: " + err.Error())
	}
	return shim.Success([]byte("Chaincode init succcessfully"))
}

func (cc *DataMarket) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fcn, params := stub.GetFunctionAndParameters()

	if fcn == "uploadPrivateData" {
		return cc.uploadPrivateData(stub, params)
	} else if fcn == "compareAndPutPrivateData" {
		return cc.compareAndPutPrivateData(stub, params)
	} else if fcn == "getPrivateData" {
		return cc.getPrivateData(stub, params)
	}
	return shim.Error("Received unknown function invocation" + fcn)
}

func (cc *DataMarket) uploadPrivateData(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		return shim.Error("Failed to upload data with error: " + err.Error())
	}
	MSP := string(mspid)
	txId := stub.GetTxID()
	if err != nil {
		return shim.Error("Failed to upload data with error: " + err.Error())
	}
	if len(params) != 2 {
		return shim.Error("Incorrect number of params" + mspid)
	}
	dataOwner := params[0]
	dataDesc := params[1]

	newPrivateData := PrivateData{Id: txId, Owner: dataOwner, DataDesc: dataDesc}
	newPrivateDataBytes, err := json.Marshal(newPrivateData)
	err = stub.PutPrivateData("_implicit_org_"+MSP, txId, newPrivateDataBytes)
	if err != nil {
		return shim.Error("Can not put private data" + err.Error())
	}

	txTimeStamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Can not get tx timestamp with error: " + err.Error())
	}
	uploadEvent := UploadEvent{Id: txId, Owner: dataOwner, UploadTime: txTimeStamp.GetSeconds()}
	uploadEventBytes, err := json.Marshal(uploadEvent)
	if err != nil {
		return shim.Error("Failed to Marshal uploadEvent, error: " + err.Error())
	}
	err = stub.SetEvent("uploadEvent", uploadEventBytes)
	if err != nil {
		return shim.Error("Failed to emit uploadEvent, error: " + err.Error())
	}

	return shim.Success(newPrivateDataBytes)
}

func (cc *DataMarket) compareAndPutPrivateData(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 5 {
		return shim.Error("Incorrect number of params")
	}
	// Get params and compare hash to check valid data
	mspidOfDataOwner := params[0]
	dataID := params[1]
	dataString := params[2]
	onChainDataHash, err := stub.GetPrivateDataHash("_implicit_org_"+mspidOfDataOwner, dataID)
	if err != nil {
		return shim.Error("Can not get hash of data with id: " + dataID + "by error:" + err.Error())
	}
	if !isValidData(onChainDataHash, []byte(dataString)) {
		return shim.Error("Data does not match with hash")
	}

	// Get MSPID and put private data to org if data is valid
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("Can not get MSPID of sender with error: " + err.Error())
	}
	err = stub.PutPrivateData("_implicit_org_"+mspid, dataID, []byte(dataString))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get param and emit share data event
	from := params[3]
	to := params[4]
	shareEvent := ShareEvent{From: from, To: to, Id: dataID}
	shareEventBytes, err := json.Marshal(shareEvent)
	if err != nil {
		return shim.Error("Failed to Marshal shareEvent, error: " + err.Error())
	}
	err = stub.SetEvent("shareEvent", shareEventBytes)
	if err != nil {
		return shim.Error("Failed to emit shareEvent, error: " + err.Error())
	}
	// Invoke ERC20 chaincode to transfer fee, 1 share = 1 token
	return invokeERC20(stub, from, to, 1)
}

// transfer(from, to, amount)
func invokeERC20(stub shim.ChaincodeStubInterface, from string, to string, amount int) pb.Response {
	var args [][]byte

	args = append(args, []byte("transfer"))
	args = append(args, []byte(from))
	args = append(args, []byte(to))
	args = append(args, []byte(strconv.Itoa(amount)))
	response := stub.InvokeChaincode("erc20", args, "mychannel")
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Invoke failed, error: %s", response.Payload)
		return shim.Error(errStr)
	}
	return shim.Success([]byte("Success put private data and transfer fee"))
}

func isValidData(onChainDatahash []byte, dataBytes []byte) bool {
	// [32]byte of hash
	dataHash := sha256.Sum256(dataBytes)
	// convert to byte
	hashByte := dataHash[:]
	if !bytes.Equal(hashByte, onChainDatahash) {
		return false
	}
	return true
}

func (cc *DataMarket) getPrivateData(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 2 {
		return shim.Error("Incorrect number of params")
	}
	dataID := params[0]
	mspid := params[1]
	privateData, err := stub.GetPrivateData("_implicit_org_"+mspid, dataID)
	if err != nil {
		return shim.Error("Can not get private data with error: " + err.Error())
	}
	return shim.Success(privateData)
}
