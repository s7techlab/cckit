# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [owner/chaincode_owner.proto](#owner/chaincode_owner.proto)
    - [ChaincodeOwner](#extensions.owner.ChaincodeOwner)
    - [ChaincodeOwnerCreated](#extensions.owner.ChaincodeOwnerCreated)
    - [ChaincodeOwnerDeleted](#extensions.owner.ChaincodeOwnerDeleted)
    - [ChaincodeOwnerUpdated](#extensions.owner.ChaincodeOwnerUpdated)
    - [ChaincodeOwners](#extensions.owner.ChaincodeOwners)
    - [CreateOwnerRequest](#extensions.owner.CreateOwnerRequest)
    - [OwnerId](#extensions.owner.OwnerId)
    - [UpdateOwnerRequest](#extensions.owner.UpdateOwnerRequest)
  
    - [ChaincodeOwnerService](#extensions.owner.ChaincodeOwnerService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="owner/chaincode_owner.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## owner/chaincode_owner.proto



<a name="extensions.owner.ChaincodeOwner"></a>

### ChaincodeOwner
State: information stored in chaincode state about chaincode owner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| subject | [string](#string) |  | certificate subject |
| issuer | [string](#string) |  | certificate issuer |
| expires_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | cert valid not after |
| cert | [bytes](#bytes) |  | Certificate |
| updated_by_msp_id | [string](#string) |  | Creator identity info |
| updated_by_cert | [bytes](#bytes) |  | Certificate |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Updated at |






<a name="extensions.owner.ChaincodeOwnerCreated"></a>

### ChaincodeOwnerCreated
Event: new chaincode owner registered


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| subject | [string](#string) |  | certificate subject |
| issuer | [string](#string) |  | certificate issuer |
| expires_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | cert valid not after |






<a name="extensions.owner.ChaincodeOwnerDeleted"></a>

### ChaincodeOwnerDeleted
Event: chaincode owner deleted`


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| subject | [string](#string) |  | certificate subject |






<a name="extensions.owner.ChaincodeOwnerUpdated"></a>

### ChaincodeOwnerUpdated
Event: new chaincode owner registered


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| subject | [string](#string) |  | certificate subject |
| expires_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | cert valid not after |






<a name="extensions.owner.ChaincodeOwners"></a>

### ChaincodeOwners
List: Chaincode owners


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | repeated |  |






<a name="extensions.owner.CreateOwnerRequest"></a>

### CreateOwnerRequest
Request: register owner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| cert | [bytes](#bytes) |  | Certificate |






<a name="extensions.owner.OwnerId"></a>

### OwnerId
Id: owner identifier


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| subject | [string](#string) |  | Certificate subject |






<a name="extensions.owner.UpdateOwnerRequest"></a>

### UpdateOwnerRequest
Request: update owner certificate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| msp_id | [string](#string) |  | Msp Id |
| cert | [bytes](#bytes) |  | Current certificate |





 

 

 


<a name="extensions.owner.ChaincodeOwnerService"></a>

### ChaincodeOwnerService
ChaincodeOwnerService allows to store information about chaincode &#34;owners&#34; in chaincode state

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetOwnerByTxCreator | [.google.protobuf.Empty](#google.protobuf.Empty) | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | Checks tx creator is owner |
| ListOwners | [.google.protobuf.Empty](#google.protobuf.Empty) | [ChaincodeOwners](#extensions.owner.ChaincodeOwners) | Get owners list |
| GetOwner | [OwnerId](#extensions.owner.OwnerId) | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | Get owner by msp_id and certificate subject |
| CreateOwner | [CreateOwnerRequest](#extensions.owner.CreateOwnerRequest) | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | Register new chaincode owner, method can be call by current owner or if no owner exists If chaincode owner with same MspID, certificate subject and issuer exists - throws error |
| CreateOwnerTxCreator | [.google.protobuf.Empty](#google.protobuf.Empty) | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | Register tx creator as chaincode owner |
| UpdateOwner | [UpdateOwnerRequest](#extensions.owner.UpdateOwnerRequest) | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | Update chaincode owner. Msp id and certificate subject must be equal to current owner certificate |
| DeleteOwner | [OwnerId](#extensions.owner.OwnerId) | [ChaincodeOwner](#extensions.owner.ChaincodeOwner) | Delete owner |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

