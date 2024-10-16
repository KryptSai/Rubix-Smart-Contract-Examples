package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"mime/multipart"
	"net/http"
	"os"

	wasm "github.com/bytecodealliance/wasmtime-go"
	"github.com/joho/godotenv"
)

var Tokenchaindatapointer int

type ContractExecution struct {
	wasmPath        string
	stateFile       string
	initialised     bool
	pointerPosition int
	instance        *wasm.Instance
	store           *wasm.Store
	memory          *wasm.Memory

	data []byte
	// Tokenchaindatapointer int
}

type Action struct {
	Function string        `json:"function"`
	Args     []interface{} `json:"args"`
}

type generateTokendata struct {
	Did         string `json:"did"`
	Userid      string `json:"UserId"`
	NftFileinfo string `json:"NFTFileInfo"`
	NftFile     string `json:"NFTFile"`
}

// type deployTokendata struct {
// 	Did         string `json:"did"`
// 	Userid      string `json:"UserId"`
// 	NftFileinfo string `json:"NFTFileInfo"`
// 	NftFile     string `json:"NFTFile"`
// }

type deployTokendata struct {
	Nft        string `json:"nft"`
	Did        string `json:"Did"`
	QuorumType int32  `json:"Quorumtype"`
	// port       string
}
type subscribeNFTTokendata struct {
	Nft string `json:"nft"`
	// Port string `json:"port"`
}
type executeNFTTokendata struct {
	Nft        string  `json:"nft"`
	Executor   string  `json:"executor"`
	Receiver   string  `json:"receiver"`
	QuorumType int32   `json:"quorumType"`
	Comment    string  `json:"comment"`
	NFTValue   float64 `json:"nftValue"`
}

func NewContractExecution(contractId string, port string) (*ContractExecution, error) {
	fmt.Println("Port", port)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	fmt.Println("Contract ID", contractId)
	path := os.Getenv(port) + "SmartContract/"
	c := &ContractExecution{
		wasmPath:  fmt.Sprintf(path+"%s/binaryCodeFile.wasm", contractId),
		stateFile: fmt.Sprintf(path+"%s/SchemaCodeFile.json", contractId),
	}
	fmt.Println("Path is ", path)
	fmt.Println("ContractExecution:", c)
	wasmBytes, err := os.ReadFile(c.wasmPath)
	if err != nil {
		return nil, err
	}

	engine := wasm.NewEngine()
	linker := wasm.NewLinker(engine)
	linker.DefineWasi()

	c.store = wasm.NewStore(engine)
	if c.store == nil {
		fmt.Println("not able to create a new store")
	}
	fmt.Println("c.store", c.store.Engine, c.store)
	//
	linker.FuncWrap("env", "mint", c.mintNFTToken)
	linker.FuncWrap("env", "deploy", c.deployNFTToken)
	linker.FuncWrap("env", "subscribe1", c.subscribeNFTToken1)
	linker.FuncWrap("env", "subscribe2", c.subscribeNFTToken2)
	linker.FuncWrap("env", "execute", c.executeNFTToken)
	linker.FuncWrap("env", "test", Test)
	linker.FuncWrap("env", "test2", c.Test2)
	module, err := wasm.NewModule(c.store.Engine, wasmBytes)
	if err != nil {
		fmt.Println("failed to compile new wasm module,err:", err)
		return nil, err
	}
	instance, err := linker.Instantiate(c.store, module)
	if err != nil {
		fmt.Println("failed to instantiate wasm module,err:", err)
		return nil, err
	}
	allocFn := instance.GetExport(c.store, "alloc").Func()
	address, err := allocFn.Call(c.store)
	if err != nil {
		return nil, err
	}

	c.pointerPosition = int(address.(int32))

	c.instance = instance
	c.memory = instance.GetExport(c.store, "memory").Memory()
	c.initialised = true
	fmt.Println("Pointer:", c.pointerPosition)
	fmt.Println("initialisation status", c.initialised)

	return c, nil
}

func (c *ContractExecution) write(str string) int {
	if !c.initialised {
		panic("Contract not initialised")
	}
	ptr := c.pointerPosition
	fmt.Print("length of the string is:", len(str))
	fmt.Print("\n Writing to memory: ")
	fmt.Println(str)

	fmt.Print("Pointer position: ")
	fmt.Println(ptr)

	copy(
		c.memory.UnsafeData(c.store)[ptr:],
		[]byte(str),
	)

	c.pointerPosition += len(str) + 1
	fmt.Println("Latest pointer position", c.pointerPosition)
	return ptr
}
func Test(len int32) {
	fmt.Println("Length of the Vector is:", len)
}
func (c *ContractExecution) Test2(pointer int32) {
	fmt.Println("Test2 function called")
	//reading the data from the wasm memory
	view := c.memory.UnsafeData(c.store)[pointer:]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in Test2 func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in Test2 function is:", str)
	fmt.Println("data in Test2 is: ", string(c.data))

}
func (c *ContractExecution) ProcessActions(actions []Action, jsonStr string) {
	if !c.initialised {
		panic("Contract not initialised")
	}

	fmt.Println("The given json string ", jsonStr)
	for _, action := range actions {
		// map on action.args and store to pointers
		pointers := make([]interface{}, len(action.Args))

		pointers[0] = c.write(jsonStr)
		Tokenchaindatapointer = pointers[0].(int)
		fmt.Println("Pointers in ProcessActions function is:", Tokenchaindatapointer)
		fmt.Println("Pointers in ProcessActions function is:", pointers)
		functionRef := c.instance.GetExport(c.store, action.Function)
		fmt.Println(functionRef)
		fmt.Println("Function", action.Function)
		_, err := functionRef.Func().Call(c.store, pointers...)
		if err != nil {
			fmt.Printf("Error has occured in process action for the function %v, err: %v\n", action.Function, err)
			return
		}
	}

}
func (c *ContractExecution) DeployFuncCaller(action Action) {
	fmt.Println("Pointers in DeployFuncCaller function is:", Tokenchaindatapointer)
	functionRef := c.instance.GetExport(c.store, action.Function)
	fmt.Println(functionRef)
	fmt.Println("Function", action.Function)
	functionRef.Func().Call(c.store, Tokenchaindatapointer)

}
func (c *ContractExecution) SubscribeFuncCaller(action Action) {
	fmt.Println("Pointers in SubscribeFuncCaller function is:", Tokenchaindatapointer)
	functionRef := c.instance.GetExport(c.store, action.Function)
	fmt.Println(functionRef)
	fmt.Println("Function", action.Function)
	functionRef.Func().Call(c.store, Tokenchaindatapointer)

}
func (c *ContractExecution) ExecuteNFTFuncCaller(action Action) {
	fmt.Println("Pointers in ExecuteNFTFuncCaller function is:", Tokenchaindatapointer)
	functionRef := c.instance.GetExport(c.store, action.Function)
	fmt.Println(functionRef)
	fmt.Println("Function", action.Function)
	s, err := functionRef.Func().Call(c.store, Tokenchaindatapointer)
	if err != nil || s == nil {
		fmt.Printf("Error has occured in while calling ExecuteNFT function action for the function %v, err: %v\n", action.Function, err)
		return
	}

}

func generateToken(data generateTokendata, port string) {
	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)
	fmt.Println("Printing the data in generateToken function", data)

	// Add form fields (simple text fields)
	writer.WriteField("did", data.Did)
	writer.WriteField("UserId", data.Userid)

	// Add the NFTFile to the form
	fmt.Println("NFT file name is:", data.NftFile)
	nftFile, err := os.Open(data.NftFile)
	if err != nil {
		fmt.Println("Error opening NFT file:", err)
		return
	}
	defer nftFile.Close()

	// Add the NFTFile part to the form
	nftFormFile, err := writer.CreateFormFile("NFTFile", data.NftFile)
	if err != nil {
		fmt.Println("Error creating NFT form file:", err)
		return
	}

	_, err = io.Copy(nftFormFile, nftFile)
	if err != nil {
		fmt.Println("Error copying NFT file content:", err)
		return
	}

	// Add the NFTFileInfo to the form
	fmt.Println("NFTFileInfo file name is:", data.NftFileinfo)
	nftFileInfo, err := os.Open(data.NftFileinfo)
	if err != nil {
		fmt.Println("Error opening NFTFileInfo file:", err)
		return
	}
	defer nftFileInfo.Close()

	// Add the NFTFileInfo part to the form
	nftInfoFormFile, err := writer.CreateFormFile("NFTFileInfo", data.NftFileinfo)
	if err != nil {
		fmt.Println("Error creating NFTFileInfo form file:", err)
		return
	}

	_, err = io.Copy(nftInfoFormFile, nftFileInfo)
	if err != nil {
		fmt.Println("Error copying NFTFileInfo content:", err)
		return
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart writer:", err)
		return
	}

	// Create the request URL
	url := fmt.Sprintf("http://localhost:%s/api/createnft", port)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set the Content-Type header to multipart/form-data with the correct boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request in generateToken fun:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Read and print the response body
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(data2))

	defer resp.Body.Close()
}
func (c *ContractExecution) mintNFTToken(pointer int32, len int32) {
	fmt.Println("mintToken function called")
	//reading the data from the wasm memory
	view := c.memory.UnsafeData(c.store)[pointer : pointer+len]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in mintToken func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in mintToken function is:", str)
	fmt.Println("data in mintToken is: ", string(c.data))
	var response generateTokendata
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in mintToken:", err3)
	}
	port := "20024"
	generateToken(response, port)

}

func (c *ContractExecution) deployNFTToken(pointer int32, len int32) {
	fmt.Println("deployToken function called")
	//reading the data from the wasm memory
	view := c.memory.UnsafeData(c.store)[pointer : pointer+len]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in deployToken func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in deployToken function is:", str)
	fmt.Println("data in deployToken is: ", string(c.data))
	var response deployTokendata
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in mintToken:", err3)
	}
	port := "20024"
	InitiateDeployNft(response, port)

}
func (c *ContractExecution) subscribeNFTToken1(pointer int32, len int32) {
	fmt.Println("subscribeNFTToken1 function called")
	//reading the data from the wasm memory
	view := c.memory.UnsafeData(c.store)[pointer : pointer+len]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in subscribeNFTToken func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in subscribeNFTToken function is:", str)
	fmt.Println("data in subscribeNFTToken is: ", string(c.data))
	var response subscribeNFTTokendata
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in subscribeNFTToken1:", err3)
	}
	port := "20024"
	InitiateSubscribeNFT(response, port)
}
func (c *ContractExecution) subscribeNFTToken2(pointer int32, len int32) {
	fmt.Println("subscribeNFTToken2 function called")
	//reading the data from the wasm memory
	view := c.memory.UnsafeData(c.store)[pointer : pointer+len]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in subscribeNFTToken func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in subscribeNFTToken function is:", str)
	fmt.Println("data in subscribeNFTToken is: ", string(c.data))
	var response subscribeNFTTokendata
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in subscribeNFTToken2:", err3)
	}
	port := "20024"
	InitiateSubscribeNFT(response, port)
}

func InitiateDeployNft(data deployTokendata, port string) {
	// data := map[string]interface{}{
	// 	"NFT":        NFT,
	// 	"DID":        Did,
	// 	"QuorumType": QuorumType,
	// }
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error in marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/deploy-nft", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in DeployNft :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	SignatureResponse(id, port)

	defer resp.Body.Close()

}
func InitiateSubscribeNFT(data subscribeNFTTokendata, port string) {
	// data := map[string]interface{}{
	// 	"NFT": NFT,
	// }
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/subscribe-nft", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in subscribe NFT :", string(data2))

	defer resp.Body.Close()

}

// func ExecuteNFTFunc() {
// 	NFT := "QmfUxjmQ9Kuykot2KrMfALYP6c6VErEGMwP33DrYZ1nJaY"
// 	Executor := "bafybmidcbhlerxfkrgfcjzi6fd442efcjx6lnbi5lx2p3l3o6a5qzjclfi"
// 	Receiver := "bafybmihkhzcczetx43gzuraoemydxntloct6qb4jkix6xo26fv5jdefq3a"
// 	QuorumType := 2
// 	Comment := "Test execute nft"
// 	NFTValue := 1.0
// 	port := "20024"
// 	InitiateExecutenft(NFT, Executor, Receiver, QuorumType, Comment, NFTValue, port)

// }
func (c *ContractExecution) executeNFTToken(pointer int32, len int32) {
	fmt.Println("executeNFTToken function called")
	fmt.Println("length provided in executeNFTToken function is", len)
	//reading the data from the wasm memory
	view := c.memory.UnsafeData(c.store)[pointer : pointer+len]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in executeNFTToken func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in executeNFTToken function is:", str)
	fmt.Println("data in executeNFTToken is: ", string(c.data))
	var response executeNFTTokendata
	//Unmarshaling the data which has been read from the wasm memory
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in executeNFTToken function is:", err3)
	}
	port := "20024"
	InitiateExecutenft(response, port)

}

// func SubscribeNFTFunc() {
// 	NFT := "QmfUxjmQ9Kuykot2KrMfALYP6c6VErEGMwP33DrYZ1nJaY"

// 	port := "20023"
// 	InitiateSubscribeNFT()

// }
// func SubscribeNFTFunc2() {
// 	NFT := "QmfUxjmQ9Kuykot2KrMfALYP6c6VErEGMwP33DrYZ1nJaY"

// 	port := "20024"
// 	InitiateSubscribeNFT(NFT, port)

// }

func InitiateExecutenft(data executeNFTTokendata, port string) {
	// func InitiateExecutenft(NFT string, Executor string, Receiver string, QuorumType int, Comment string, NFTValue float64, port string) {
	// data := map[string]interface{}{
	// 	"NFT":        NFT,
	// 	"receiver":   Receiver,
	// 	"comment":    Comment,
	// 	"executor":   Executor,
	// 	"quorumType": QuorumType,
	// 	"NFTValue":   NFTValue,
	// }
	fmt.Println("printing the data in InitiateExecutenft function is:", data)
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/execute-nft", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status in InitiateExecutenft:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in InitiateExecutenft :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	SignatureResponse(id, port)

	defer resp.Body.Close()

}
