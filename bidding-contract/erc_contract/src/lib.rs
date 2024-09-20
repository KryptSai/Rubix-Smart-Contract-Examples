use std::mem;
use std::ffi::{CString, CStr};
use std::os::raw::c_void;
extern crate serde;
extern crate serde_json;
//Test comments 09aabcdfd
#[macro_use] extern crate serde_derive;

#[derive(Serialize, Deserialize, Debug)]
struct SCTDataReply {
    BlockNo: u32,
    BlockId: String,
    SmartContractData: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct SmartContractData {
    did: String,
    bid: f64,
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
    // fn rbt_transfer();
    // fn create_did(ptr: *mut u8);   // fn rbt_transfer();
    // fn create_did(ptr: *mut u8);
    fn mint(ptr:*mut u8);
}
#[no_mangle]
pub unsafe extern "C" fn dealloc(ptr: *mut c_void) {
    let _ = Vec::from_raw_parts(ptr, 0, 1024);
}

#[no_mangle]
pub unsafe fn bid(ptr: *mut u8) {
    // Assume get_blocks() returns a valid JSON string pointer
    // For testing, we'll use the hardcoded JSON data directly
    //  let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    // Deserialize the JSON data into a vector of SCTDataReply structs
    // let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    // if 1==1{rbt_transfer()}
    let mut empty_vec = Vec::new();
    mint(empty_vec.as_mut_ptr());
    // Find the block with the highest bid
    // match find_highest_bid_did(&blocks) {
    //     Some((block_no, max_bid)) => println!("The block with the highest bid is BlockNo {} with a bid of {}", block_no, max_bid),
    //     None => println!("No valid bids found."),
    // }

}


