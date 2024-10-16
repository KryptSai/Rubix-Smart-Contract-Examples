use std::mem;
use std::ffi::CStr;
use std::os::raw::c_void;
extern crate serde;
extern crate serde_json;
#[macro_use] extern crate serde_derive;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct SCTDataReply {
    BlockNo: u32,
    BlockId: String,
    SmartContractData: String,
}
//testComments-73
#[derive(Serialize, Deserialize, Debug,Clone)]
struct Createnft {
    did :        String, 
    Userid:      String, 
    Nftfileinfo: String,
    Nftfile:     String,
}
#[derive(Serialize, Deserialize, Debug,Clone)]
struct Deploynft{
    nft:         String,
    Did :        String,
    Quorumtype:  i32,
}
#[derive(Serialize, Deserialize, Debug,Clone)]
struct Subscribenft{
    nft:         String,
    // Did :        String,
    // Quorumtype:  i32,
}
#[derive(Serialize, Deserialize, Debug,Clone)]
struct Executenft {
    nft:      String,
    executor :   String,
    receiver:    String,
    quorumType:  i32,
    comment:     String,
    nftValue:   f64,
    
}

#[derive(Deserialize,Serialize, Debug,Clone)]
#[serde(untagged)] // Untagged enums are used when the type is inferred from the structure.
enum SmartContractData {
    CreateNft(Createnft),
    DeployNft(Deploynft),
    // SubscribeNft(Subscribenft),
    ExecuteNft(Executenft),
}

/*This alloc function is used to allocate 1024 bytes and returns a pointer.
When this function is called in the Go code we will receive the pointer.
Whatever data we need is being pushed onto this memory location.
*/
#[no_mangle]
pub extern "C" fn alloc() -> *mut c_void {
    let mut buf = Vec::with_capacity(1024);
    let ptr = buf.as_mut_ptr();

    mem::forget(buf);

    ptr
}
extern "C" {
    fn mint(ptr:*mut u8,len: i32);
    fn test(len:usize );
    fn test2(pointer:*mut u8 );
    // fn mint(ptr:*mut u8);
    fn deploy(ptr:*mut u8,len:i32);
    fn subscribe1(ptr:*mut u8,len:i32);
    fn subscribe2(ptr:*mut u8,len:i32);
    fn execute(ptr:*mut u8,len:i32);
    
}
#[no_mangle]
pub unsafe extern "C" fn dealloc(ptr: *mut c_void) {
    let _ = Vec::from_raw_parts(ptr, 0, 1024);
}
//extracts the smartcontract data from the blocks
// fn extract_smartcontract_data(blocks: &[SCTDataReply]) -> Vec<Createnft> {   
//     let mut vec_sc_data: Vec<Createnft> = Vec::new();
//     for block in blocks {
//         let  scdata = &block.SmartContractData;
//         if scdata.is_empty() {
//             continue;
//         }
//         if let Ok(data) = serde_json::from_str::<Createnft>(&block.SmartContractData) {
//            vec_sc_data.push(data);
//         }
//     }
//     vec_sc_data
// }


fn extract_smartcontract_data(blocks: &[SCTDataReply]) -> Vec<SmartContractData> {
    let mut vec_sc_data: Vec<SmartContractData> = Vec::new();
    
    for block in blocks {
        let scdata = &block.SmartContractData;
        
        if scdata.is_empty() {
            continue;
        }
        
        // Try deserializing to the enum that can hold different types (Createnft or Deploynft)
        if let Ok(data) = serde_json::from_str::<SmartContractData>(scdata) {
            vec_sc_data.push(data);
        } 
         
    }
    
    vec_sc_data
}
#[no_mangle]
pub unsafe fn create_nft(ptr: *mut u8) {
  
     let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    //  test("Om Sri Sairam,Good morning".len());
    // Deserialize the JSON data into a vector of SCTDataReply structs
    let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    // test("Om Sri Sairam,Good morning".len());
    let smartcontract_data_vec = extract_smartcontract_data(&blocks);
    // test("Om Sri Sairam".len());
    let vec_len = smartcontract_data_vec.len();
    let createnft_data =  smartcontract_data_vec[0].clone();
    //Serialize the smartcontract data 
    let mut serialized_createnft_data = serde_json::to_string(&createnft_data).unwrap();
    test(vec_len);
    let mut serialized_createnft_data_len = serialized_createnft_data.len();
    // let deploynft_data =  smartcontract_data_vec[vec_len-1].clone();

    // let mut serialized_deploynft_data = serde_json::to_string(&deploynft_data).unwrap();
    // let mut serialized_deploynft_data_len = serialized_deploynft_data.len();
    // test();
    // creating NFT by calling it from the wasm
    mint(serialized_createnft_data.as_mut_ptr(),serialized_createnft_data_len.try_into().unwrap());
    // test();
    // deploy(serialized_deploynft_data.as_mut_ptr(),serialized_deploynft_data_len.try_into().unwrap());

}
#[no_mangle]
pub unsafe fn deploy_nft(ptr: *mut u8) {
    let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    // Deserialize the JSON data into a vector of SCTDataReply structs
    let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    // test();
    let smartcontract_data_vec = extract_smartcontract_data(&blocks);
    // test();
    let vec_len = smartcontract_data_vec.len();
    let deploynft_data =  smartcontract_data_vec[1].clone();
    let mut serialized_deploynft_data = serde_json::to_string(&deploynft_data).unwrap();
    let mut serialized_deploynft_data_len = serialized_deploynft_data.len();
    deploy(serialized_deploynft_data.as_mut_ptr(),serialized_deploynft_data_len.try_into().unwrap());

}
#[no_mangle]
pub unsafe fn subscribe_nft1(ptr: *mut u8) {
    let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    // Deserialize the JSON data into a vector of SCTDataReply structs
    let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    // test();
    let smartcontract_data_vec = extract_smartcontract_data(&blocks);
    // test("sairam".len());
    let vec_len = smartcontract_data_vec.len();
    let subscribenft_data =  smartcontract_data_vec[2].clone();
    let mut serialized_subscribenft_data = serde_json::to_string(&subscribenft_data).unwrap();
    test(vec_len);
    let mut serialized_subscribenft_data_len = serialized_subscribenft_data.len();
    subscribe1(serialized_subscribenft_data.as_mut_ptr(),serialized_subscribenft_data_len.try_into().unwrap());
    
}
#[no_mangle]
pub unsafe fn subscribe_nft2(ptr: *mut u8) {
    let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    // Deserialize the JSON data into a vector of SCTDataReply structs
    let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    // test();
    let smartcontract_data_vec = extract_smartcontract_data(&blocks);
    test("sairam".len());
    let vec_len = smartcontract_data_vec.len();
    let subscribenft_data =  smartcontract_data_vec[2].clone();
    let mut serialized_subscribenft_data = serde_json::to_string(&subscribenft_data).unwrap();
    test(vec_len);
    let mut serialized_subscribenft_data_len = serialized_subscribenft_data.len();
    subscribe2(serialized_subscribenft_data.as_mut_ptr(),serialized_subscribenft_data_len.try_into().unwrap());
    
}
#[no_mangle]
pub unsafe fn execute_nft(ptr: *mut u8) {
    let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    // Deserialize the JSON data into a vector of SCTDataReply structs
    let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    // test();
    let smartcontract_data_vec = extract_smartcontract_data(&blocks);
    
    // test("sairam".len());
    let vec_len = smartcontract_data_vec.len();
    let executenft_data =  smartcontract_data_vec[vec_len-1].clone();
    // test2(executenft_data.as_ptr);
    let mut serialized_executenft_data: String = match serde_json::to_string(&executenft_data) {
        Ok(res) => res,
        Err(err) => format!("'hi': '{err}'")
    };
    
    test(vec_len);
    let mut serialized_executenft_data_len = serialized_executenft_data.len();
    execute(serialized_executenft_data.as_mut_ptr(),serialized_executenft_data_len.try_into().unwrap());
    
}

fn do_some_panic(_x: u8) {
    panic!("cessna 172 ")
}


