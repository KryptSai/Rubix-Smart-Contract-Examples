package main

import (
	"bidding-contract/contract"
	contractModule "bidding-contract/contract"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var ContractExecutionInstace *contractModule.ContractExecution

// var Tokenchaindatapointer int

type SmartContractDataReply struct {
	BasicResponse
	SCTDataReply []SCTDataReply
}

type BasicResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type SCTDataReply struct {
	BlockNo           uint64
	BlockId           string
	SmartContractData string
}

func GenerateSmartContract() {
	/*
		This did, wasmPath, schemaPath, rawcodePath and Port should be replaced according to your Rubix node configuration and
		the respective paths
	*/
	did := "bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a"
	wasmPath := "./nft_contract/target/wasm32-unknown-unknown/debug/nft_contract.wasm"
	schemaPath := "/home/rubix/Sai-Rubix/Nft-Rubix/nft/nft_contract/target/nft_contract.json"
	rawCodePath := "./nft_contract/src/lib.rs"
	port := "20023"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

// This function is intended to pass the smart contract hash which is retruned while generating smart contract
func smartContractHash() string {
	return "QmYvV7Qx3ujrqT4LkmLtGwdusYYnuraoyK5XScsVzqnCTw"
}

func DeploySmartContract() {
	/*
		port : The port corresponding to the deployer node.
	*/
	comment := "Deploying NFT contract"
	deployerAddress := "bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a"
	quorumType := 2
	rbtAmount := 1
	smartContractToken := smartContractHash()
	port := "20023"
	id := contract.DeploySmartContract(comment, deployerAddress, quorumType, rbtAmount, smartContractToken, port)
	fmt.Println("Contract ID: " + id)
	contract.SignatureResponse(id, port)

}

// func Deploynft() {
// 	did := "bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a"
// 	nft := "QmQDkDP7kutZGQUv4J7ppt7Ksz3de9PHeG91vUT3ok7C7J"
// 	quorumtype := 2
// 	port := "20023"
// 	contract.DeployNft(nft, did, quorumtype, port)

// }

func ExecuteSmartContractForCreateNFTdata() {
	/*

		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on node11"
	executorAddress := "bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi"
	quorumType := 2
	smartContractData := `{"did":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","Userid":"Saibaba","Nftfileinfo":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/metadata.json","Nftfile":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/testimage20.png"}`
	smartContractToken := smartContractHash()
	port := "20024"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}
func ExecuteSmartContractForDeployNFTdata() {
	/*
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on node10"
	executorAddress := "bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi"
	quorumType := 2
	smartContractData := `{"nft":"Qman2hz94gQHw8MkEboRUPLRtocL4RjDRVewhofdUDofE1","Did":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","Quorumtype":2}`
	smartContractToken := smartContractHash()
	port := "20024"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}
func ExecuteSmartContractForSubscribeNFTnode23() {
	/*
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract to SubscribeNFT on Node23"
	executorAddress := "bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a"
	quorumType := 2
	smartContractData := `{"nft":"Qman2hz94gQHw8MkEboRUPLRtocL4RjDRVewhofdUDofE1"}`
	smartContractToken := smartContractHash()
	port := "20023"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}
func ExecuteSmartContractForSubscribeNFTnode24() {
	/*
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract to SubscribeNFT on Node24"
	executorAddress := "bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi"
	quorumType := 2
	smartContractData := `{"nft":"Qman2hz94gQHw8MkEboRUPLRtocL4RjDRVewhofdUDofE1"}`
	smartContractToken := smartContractHash()
	port := "20024"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractForExecuteNFT() {
	/*
		port : The port corresponding to the executor node.
	*/
	comment := "Executing Test Smart Contract on Node12"
	executorAddress := "bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi"
	quorumType := 2
	// smartContractData := `{"nft":"QmSUJ7r5D72ae7kaxbs14wHdyGfbNASaFLbKkHVZjot5o4","executor":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","Quorumtype":2}`
	smartContractData := `{"nft":"Qman2hz94gQHw8MkEboRUPLRtocL4RjDRVewhofdUDofE1","executor":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","receiver":"bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a","quorumType": 2,"comment":"Test execute nft","nftValue":1.0}`
	// smartContractData := `{"executor":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","nft":"QmfUxjmQ9Kuykot2KrMfALYP6c6VErEGMwP33DrYZ1nJaY","quorumType": 2,"comment": "Test execute nft"}`
	// smartContractData := `{"did":"bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi","Userid":"Saibaba","Nftfileinfo":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/metadata.json","Nftfile":"/home/rubix/Sai-Rubix/Nft-Rubix/nft/testimage20.png"}`
	smartContractToken := smartContractHash()
	port := "20024"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}
func SubscribeSmartContractnode23(port string) {
	contractToken := smartContractHash()
	//	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

// This function is responsible for subscribing to a particular smart contract.
func SubscribeSmartContractnode24(port string) {
	contractToken := smartContractHash()
	//	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

// func SubscribeSmartContractnode11(port string) {
// 	contractToken := smartContractHash()
// 	//contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
// 	contract.SubscribeSmartContract(contractToken, port)
// }

// func SubscribeSmartContractnode12(port string) {
// 	contractToken := smartContractHash()
// 	//contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
// 	contract.SubscribeSmartContract(contractToken, port)
// }

// Function to manually set a delay and trigger ContractExecution
func CreatenftFunc(port string, seconds int) {
	contractId := smartContractHash()
	fmt.Printf("Setting a delay of %d seconds before triggering ContractExecution...\n", seconds)
	time.After(time.Duration(seconds))
	contractExec, err := contractModule.NewContractExecution(contractId, port)
	ContractExecutionInstace = contractExec
	smartContractTokenData := contract.GetSmartContractData(port, contractId)
	fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))
	var dataReply SmartContractDataReply

	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data reply in RunSmartContract", dataReply)
	action := contractModule.Action{
		Function: "create_nft",
		Args:     []interface{}{""},
	}
	actions := []contractModule.Action{action}
	fmt.Println("actions in SetDelay function", actions)
	smartContractData := dataReply.SCTDataReply
	fmt.Println("Smart Contract Data :", smartContractData)
	jsonString, err := json.Marshal(smartContractData)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Print the JSON string
	fmt.Println(string(jsonString))
	contractExec.ProcessActions(actions, string(jsonString))
}
func DeployNFTFunc() {
	action := contractModule.Action{
		Function: "deploy_nft",
		Args:     []interface{}{""},
	}

	ContractExecutionInstace.DeployFuncCaller(action)

}
func SubscribeNFTFunc1() {
	action := contractModule.Action{
		Function: "subscribe_nft1",
		Args:     []interface{}{""},
	}
	ContractExecutionInstace.SubscribeFuncCaller(action)

}

func SubscribeNFTFunc2() {
	action := contractModule.Action{
		Function: "subscribe_nft2",
		Args:     []interface{}{""},
	}

	ContractExecutionInstace.SubscribeFuncCaller(action)

}
func ExecuteNFTFunc() {
	action := contractModule.Action{
		Function: "execute_nft",
		Args:     []interface{}{""},
	}

	ContractExecutionInstace.ExecuteNFTFuncCaller(action)

}

func main() {

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enlighten me with the function to be executed ")
		fmt.Println(`
		1. Generate Contract 
		2. Subscribe Contract node9 aka Deployer Node
		3. Subscribe Contract node10 
		4. Deploy Contract
		5. Execute Contract for creating NFT data
		6. Execute Contract for deploy NFT data
		7.Execute Contract for subscribe NFT in node23 
		8.Execute Contract for subscribe NFT in node24
		9.Execute Contract for Execute NFT in node24
		10.createNFT  function
		11.Deploynft function
		12.Subscribing NFT function in 20023
		13.Subscribing NFT function in 20024
		14.ExecuteNFT function`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in TestNode1 aka Deployer Node")
			SubscribeSmartContractnode23("20023")
		case "3":
			fmt.Println("Subscribing Smart Contract in node10")
			SubscribeSmartContractnode24("20024")

		case "4":
			fmt.Println("Deploying Smart Contract in node9")
			DeploySmartContract()
		case "5":
			fmt.Println("Executing Smart Contract to createNFT")
			ExecuteSmartContractForCreateNFTdata()

		case "6":
			fmt.Println("Executing Smart Contract to deployNFT")
			ExecuteSmartContractForDeployNFTdata()
		case "7":
			fmt.Println("Executing Smart Contract for subscribe NFT for node23")
			ExecuteSmartContractForSubscribeNFTnode23()
		case "8":
			fmt.Println("Executing Smart Contract for subscribe NFT for node24")
			ExecuteSmartContractForSubscribeNFTnode24()

		case "9":
			fmt.Println("Executing Smart Contract for execute NFT")
			ExecuteSmartContractForExecuteNFT()

		case "10":
			fmt.Println("Creating an NFT: Calling the mint function")
			CreatenftFunc("20023", 20)
		case "11":
			fmt.Println("Deploying an NFT: Calling the Deploynft function")
			DeployNFTFunc()
		case "12":
			fmt.Println("Subscribing an NFT:SubscribeNFT  for node23")
			SubscribeNFTFunc1()
		case "13":
			fmt.Println("Subscribing an NFT:SubscribeNFT function for node24")
			SubscribeNFTFunc2()
		case "14":
			fmt.Println("Executing an NFT: Calling the ExecuteNFT function")
			ExecuteNFTFunc()

		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
