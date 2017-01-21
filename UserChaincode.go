package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
)

// User Chaincode implementation
type UserChaincode struct {
}

var userIndexTxStr = "_userIndexTxStr"

type UserData struct{
	USER_ID string `json:"USER_ID"`
	USER_NAME string `json:"USER_NAME"`
	BANK string `json:"BANK"`
	KYC_DATE string `json:"KYC_DATE"`
	KYC_EXPIRATION_DATE string `json:"KYC_EXPIRATION_DATE"`
	KYC_DOCUMENT string `json:"KYC_DOCUMENT"`
	}


func (t *UserChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var err error
	// Initialize the chaincode
	
	fmt.Printf("Deployment of Loyalty is completed\n")
	
	var emptyUserTxs []UserData
	jsonAsBytes, _ := json.Marshal(emptyUserTxs)
	err = stub.PutState(userIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	
	return nil, nil
}

// Add KYC Data for User
func (t *UserChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == userIndexTxStr {		
		return t.AddUser(stub, args)
	}
	return nil, nil
}


func (t *UserChaincode)  AddUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var UserDataObj UserData
	var UserDataList []UserData
	var err error

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Need 5 arguments")
	}

	// Initialize the chaincode
	UserDataObj.USER_ID = args[0]
	UserDataObj.USER_NAME = args[1]
	UserDataObj.BANK = args[2]
	UserDataObj.KYC_DATE = args[3]
	UserDataObj.KYC_EXPIRATION_DATE = args[4]
	UserDataObj.KYC_DOCUMENT = args[5]
	
	
	fmt.Printf("Input from user:%s\n", UserDataObj)
	
	userTxsAsBytes, err := stub.GetState(userIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get consumer Transactions")
	}
	json.Unmarshal(userTxsAsBytes, &UserDataList)
	
	UserDataList = append(UserDataList, UserDataObj)
	jsonAsBytes, _ := json.Marshal(UserDataList)
	
	err = stub.PutState(userIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *UserChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	
	var UserId string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting UserId of the person to query")
	}

	UserId = args[0]
	
	resAsBytes, err = t.GetUserDetails(stub, UserId)
	
	fmt.Printf("Query Response:%s\n", resAsBytes)
	
	if err != nil {
		return nil, err
	}
	
	return resAsBytes, nil
}

func (t *UserChaincode)  GetUserDetails(stub shim.ChaincodeStubInterface, UserId string) ([]byte, error) {
	
	var objFound bool
	userTxsAsBytes, err := stub.GetState(userIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get User Transactions")
	}
	var UserTxObjects []UserData
	var UserTxObjects1 []UserData
	json.Unmarshal(userTxsAsBytes, &UserTxObjects)
	length := len(UserTxObjects)
	fmt.Printf("Output from chaincode: %s\n", userTxsAsBytes)
	
	if UserId == "" {
		res, err := json.Marshal(UserTxObjects)
	 	if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
	
	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := UserTxObjects[i]
		if UserId == obj.USER_ID  {
			UserTxObjects1 = append(UserTxObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}
	
	if objFound {
		res, err := json.Marshal(UserTxObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}

func main() {
	err := shim.Start(new(UserChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}