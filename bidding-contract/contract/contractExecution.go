package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	// setup "github.com/rubixchain/rubixgoplatform/setup"

	wasm "github.com/bytecodealliance/wasmtime-go"
	"github.com/joho/godotenv"
)

type ContractExecution struct {
	wasmPath        string
	stateFile       string
	initialised     bool
	pointerPosition int
	instance        *wasm.Instance
	store           *wasm.Store
	memory          *wasm.Memory

	data []byte
}

type Action struct {
	Function string        `json:"function"`
	Args     []interface{} `json:"args"`
}

type generateTokendata struct {
	Did         string
	WasmPath    string
	SchemaPath  string
	RawCodePath string
	Port        string
}

/*Different functions which we have written here are called from wasm:
  1. "alloc"
  2. "apply_state"
  3. "get_state"
  So there will be a corresponding function with the same name in the rust code too. If the function name
  is different then the system will thorugh an error.
  The initial idea is to make all the mandatory functions into a package which can be easily imported and
  utilised while maintaining the standard for execution.*/

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

	// linker.FuncWrap("env", "rbt_transfer", c.initiateRbtTransfer)
	// linker.FuncWrap("env", "create_did", c.initiatecreateDid2)
	linker.FuncWrap("env", "mint", c.mintToken)
	module, err := wasm.NewModule(c.store.Engine, wasmBytes)
	if err != nil {
		fmt.Println("failed to compile new wasm module,err:", err)
		return nil, err
	}

	instance, err := linker.Instantiate(c.store, module)
	// instance, err := wasm.NewInstance(c.store, module, nil)
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
	//c.apply_state()

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

// func (c *ContractExecution) readAtCurrentPointer() string {
// 	if !c.initialised {
// 		panic("Contract not initialised")
// 	}

// 	pointer := c.pointerPosition
// 	fmt.Println("Pointer position in readAtCurrentPointer function", pointer)
// 	view := c.memory.UnsafeData(c.store)[pointer:]
// 	length := 0
// 	for _, byte := range view {
// 		if byte == 0 {
// 			break
// 		}
// 		length++
// 	}
// 	fmt.Println("length in readAtCurrentPointer function:", length)
// 	str := string(view[:length])
// 	c.pointerPosition += length + 1
// 	return str
// }

func (c *ContractExecution) ReadStateFile() string {
	if !c.initialised {
		panic("Contract not initialised")
	}

	file, err := os.ReadFile(c.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}

		panic(err)
	}

	return string(file)
}

// func (c *ContractExecution) apply_state() {
// 	if !c.initialised {
// 		panic("Contract not initialised")
// 	}

// 	state := c.ReadStateFile()
// 	if state != "" {
// 		pointer := c.write(state)
// 		c.instance.GetExport(c.store, "apply_state").Func().Call(c.store, pointer)
// 	}
// }

func (c *ContractExecution) ProcessActions(actions []Action, jsonStr string) {
	if !c.initialised {
		panic("Contract not initialised")
	}

	fmt.Println("The given json string ", jsonStr)
	for _, action := range actions {
		// map on action.args and store to pointers
		pointers := make([]interface{}, len(action.Args))
		// for i, arg := range action.Args {
		// 	pointers[i] = c.write(arg.(string))
		// }
		pointers[0] = c.write(jsonStr)
		fmt.Println("Pointers in ProcessActions function is:", pointers)
		functionRef := c.instance.GetExport(c.store, action.Function)
		fmt.Println(functionRef)
		fmt.Println("Function", action.Function)
		functionRef.Func().Call(c.store, pointers...)
	}

}

func (c *ContractExecution) initiateMintToken(pointer int32) {
	var data generateTokendata
	data = generateTokendata{
		Did:         "bafybmicsu4nhoifx2pwjghogfxt63ihwsj3lbicuiqekobxhqyxsq7srju",
		WasmPath:    "./bidding_contract/target/wasm32-unknown-unknown/debug/bidding_contract.wasm",
		SchemaPath:  "./data/state/bidding_contract.json",
		RawCodePath: "./bidding_contract/src/lib.rs",
		Port:        "20011",
	}
	fmt.Println("mintToken data", data)
	marshalData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error in marshaling JSON:", err)
		return
	}
	fmt.Println("marshal data", marshalData)
	copy(c.memory.UnsafeData(c.store)[pointer:pointer+int32(len(marshalData))], marshalData)
	datalen := len(marshalData)
	fmt.Println("data length is:", datalen)
	view := c.memory.UnsafeData(c.store)[pointer:]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length of the data which has been read in mintToken func:", length)
	str := string(view[:length])
	fmt.Println("data in mintToken function is:", str)
	// contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)

}
func GenerateToken(data generateTokendata) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the form fields
	_ = writer.WriteField("did", data.Did)

	// Add the binaryCodePath field
	file, _ := os.Open(data.WasmPath)
	defer file.Close()
	binaryPart, _ := writer.CreateFormFile("binaryCodePath", data.WasmPath)
	_, _ = io.Copy(binaryPart, file)

	// Add the rawCodePath field
	rawFile, _ := os.Open(data.RawCodePath)
	defer rawFile.Close()
	rawPart, _ := writer.CreateFormFile("rawCodePath", data.RawCodePath)
	_, _ = io.Copy(rawPart, rawFile)

	// Add the schemaFilePath field
	schemaFile, _ := os.Open(data.SchemaPath)
	defer schemaFile.Close()
	schemaPart, _ := writer.CreateFormFile("schemaFilePath", data.SchemaPath)
	_, _ = io.Copy(schemaPart, schemaFile)

	// Close the writer
	writer.Close()

	// Create the HTTP request
	url := fmt.Sprintf("http://localhost:%s/api/generate-smart-contract", data.Port)
	// http: //localhost:20002/api/createnft
	req, _ := http.NewRequest("POST", url, &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute Contract :", string(data2))

	// Process the response as needed
	fmt.Println("Response status code:", resp.StatusCode)
}
func (c *ContractExecution) mintToken(pointer int32) {
	c.initiateMintToken(pointer)
	// copy(c.data, c.memory.UnsafeData(c.store)[pointer:pointer+int32(c.datalen)])
	view := c.memory.UnsafeData(c.store)[pointer:]
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
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in mintToken:", err3)
	}
	GenerateToken(response)

}
