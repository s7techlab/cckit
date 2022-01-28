# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [chaincode.proto](#chaincode.proto)
    - [BlockLimit](#cckit.gateway.BlockLimit)
    - [ChaincodeEvent](#cckit.gateway.ChaincodeEvent)
    - [ChaincodeEvents](#cckit.gateway.ChaincodeEvents)
    - [ChaincodeEventsRequest](#cckit.gateway.ChaincodeEventsRequest)
    - [ChaincodeEventsStreamRequest](#cckit.gateway.ChaincodeEventsStreamRequest)
    - [ChaincodeExec](#cckit.gateway.ChaincodeExec)
    - [ChaincodeInput](#cckit.gateway.ChaincodeInput)
    - [ChaincodeInput.TransientEntry](#cckit.gateway.ChaincodeInput.TransientEntry)
    - [ChaincodeInstanceEventsRequest](#cckit.gateway.ChaincodeInstanceEventsRequest)
    - [ChaincodeInstanceEventsStreamRequest](#cckit.gateway.ChaincodeInstanceEventsStreamRequest)
    - [ChaincodeInstanceExec](#cckit.gateway.ChaincodeInstanceExec)
    - [ChaincodeInstanceInput](#cckit.gateway.ChaincodeInstanceInput)
    - [ChaincodeInstanceInput.TransientEntry](#cckit.gateway.ChaincodeInstanceInput.TransientEntry)
    - [ChaincodeLocator](#cckit.gateway.ChaincodeLocator)
    - [RawJson](#cckit.gateway.RawJson)
  
    - [InvocationType](#cckit.gateway.InvocationType)
  
    - [ChaincodeEventsService](#cckit.gateway.ChaincodeEventsService)
    - [ChaincodeInstanceEventsService](#cckit.gateway.ChaincodeInstanceEventsService)
    - [ChaincodeInstanceService](#cckit.gateway.ChaincodeInstanceService)
    - [ChaincodeService](#cckit.gateway.ChaincodeService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="chaincode.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## chaincode.proto
Gateway to network/chaincode
Two types of gateways: 1. Gateway to all chaincodes in Network 2. Gateway to some concrete chaincode instance in some channel


<a name="cckit.gateway.BlockLimit"></a>

### BlockLimit
Block limit number for event stream subscription or event list
Values can be negative


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| num | [int64](#int64) |  | Block number |






<a name="cckit.gateway.ChaincodeEvent"></a>

### ChaincodeEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event | [protos.ChaincodeEvent](#protos.ChaincodeEvent) |  |  |
| block | [uint64](#uint64) |  |  |
| tx_timestamp | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| payload | [RawJson](#cckit.gateway.RawJson) |  |  |






<a name="cckit.gateway.ChaincodeEvents"></a>

### ChaincodeEvents



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chaincode | [ChaincodeLocator](#cckit.gateway.ChaincodeLocator) |  |  |
| from_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| to_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| items | [ChaincodeEvent](#cckit.gateway.ChaincodeEvent) | repeated |  |






<a name="cckit.gateway.ChaincodeEventsRequest"></a>

### ChaincodeEventsRequest
Chaincode events list request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chaincode | [ChaincodeLocator](#cckit.gateway.ChaincodeLocator) |  |  |
| from_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| to_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| event_name | [string](#string) | repeated |  |
| limit | [uint32](#uint32) |  |  |






<a name="cckit.gateway.ChaincodeEventsStreamRequest"></a>

### ChaincodeEventsStreamRequest
Chaincode events stream request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chaincode | [ChaincodeLocator](#cckit.gateway.ChaincodeLocator) |  |  |
| from_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| to_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| event_name | [string](#string) | repeated |  |






<a name="cckit.gateway.ChaincodeExec"></a>

### ChaincodeExec
Chaincode execution specification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [InvocationType](#cckit.gateway.InvocationType) |  |  |
| input | [ChaincodeInput](#cckit.gateway.ChaincodeInput) |  |  |






<a name="cckit.gateway.ChaincodeInput"></a>

### ChaincodeInput
Chaincode invocation input


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chaincode | [ChaincodeLocator](#cckit.gateway.ChaincodeLocator) |  |  |
| args | [bytes](#bytes) | repeated | Input contains the arguments for invocation. |
| transient | [ChaincodeInput.TransientEntry](#cckit.gateway.ChaincodeInput.TransientEntry) | repeated | TransientMap contains data (e.g. cryptographic material) that might be used to implement some form of application-level confidentiality. The contents of this field are supposed to always be omitted from the transaction and excluded from the ledger. |






<a name="cckit.gateway.ChaincodeInput.TransientEntry"></a>

### ChaincodeInput.TransientEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bytes](#bytes) |  |  |






<a name="cckit.gateway.ChaincodeInstanceEventsRequest"></a>

### ChaincodeInstanceEventsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| to_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| event_name | [string](#string) | repeated |  |
| limit | [uint32](#uint32) |  |  |






<a name="cckit.gateway.ChaincodeInstanceEventsStreamRequest"></a>

### ChaincodeInstanceEventsStreamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| to_block | [BlockLimit](#cckit.gateway.BlockLimit) |  |  |
| event_name | [string](#string) | repeated |  |






<a name="cckit.gateway.ChaincodeInstanceExec"></a>

### ChaincodeInstanceExec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [InvocationType](#cckit.gateway.InvocationType) |  |  |
| input | [ChaincodeInstanceInput](#cckit.gateway.ChaincodeInstanceInput) |  |  |






<a name="cckit.gateway.ChaincodeInstanceInput"></a>

### ChaincodeInstanceInput
Chaincode instance chaincode input spec


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| args | [bytes](#bytes) | repeated | Input contains the arguments for invocation. |
| transient | [ChaincodeInstanceInput.TransientEntry](#cckit.gateway.ChaincodeInstanceInput.TransientEntry) | repeated | TransientMap contains data (e.g. cryptographic material) that might be used to implement some form of application-level confidentiality. The contents of this field are supposed to always be omitted from the transaction and excluded from the ledger. |






<a name="cckit.gateway.ChaincodeInstanceInput.TransientEntry"></a>

### ChaincodeInstanceInput.TransientEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bytes](#bytes) |  |  |






<a name="cckit.gateway.ChaincodeLocator"></a>

### ChaincodeLocator
Chaincode locator - channel name and chaincode name


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chaincode | [string](#string) |  | Chaincode name |
| channel | [string](#string) |  | Channel name |






<a name="cckit.gateway.RawJson"></a>

### RawJson



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [bytes](#bytes) |  |  |





 


<a name="cckit.gateway.InvocationType"></a>

### InvocationType
Chaincode invocation type

| Name | Number | Description |
| ---- | ------ | ----------- |
| INVOCATION_TYPE_QUERY | 0 | Simulation |
| INVOCATION_TYPE_INVOKE | 1 | Simulation and applying to ledger |


 

 


<a name="cckit.gateway.ChaincodeEventsService"></a>

### ChaincodeEventsService
Chaincode events subscription service

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EventsStream | [ChaincodeEventsStreamRequest](#cckit.gateway.ChaincodeEventsStreamRequest) | [ChaincodeEvent](#cckit.gateway.ChaincodeEvent) stream | Chaincode events stream |
| Events | [ChaincodeEventsRequest](#cckit.gateway.ChaincodeEventsRequest) | [ChaincodeEvents](#cckit.gateway.ChaincodeEvents) | Chaincode events |


<a name="cckit.gateway.ChaincodeInstanceEventsService"></a>

### ChaincodeInstanceEventsService
Chaincode instance events subscription service

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EventsStream | [ChaincodeInstanceEventsStreamRequest](#cckit.gateway.ChaincodeInstanceEventsStreamRequest) | [ChaincodeEvent](#cckit.gateway.ChaincodeEvent) stream | Chaincode events stream |
| Events | [ChaincodeInstanceEventsRequest](#cckit.gateway.ChaincodeInstanceEventsRequest) | [ChaincodeEvents](#cckit.gateway.ChaincodeEvents) | Chaincode events s |


<a name="cckit.gateway.ChaincodeInstanceService"></a>

### ChaincodeInstanceService
Chaincode instance communication service. Channel/chaincode already fixed.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exec | [ChaincodeInstanceExec](#cckit.gateway.ChaincodeInstanceExec) | [.protos.Response](#protos.Response) | Exec: Query or Invoke |
| Query | [ChaincodeInstanceInput](#cckit.gateway.ChaincodeInstanceInput) | [.protos.Response](#protos.Response) | Query chaincode on home peer. Do NOT send to orderer. |
| Invoke | [ChaincodeInstanceInput](#cckit.gateway.ChaincodeInstanceInput) | [.protos.Response](#protos.Response) | Invoke chaincode on peers, according to endorsement policy and the SEND to orderer |
| EventsStream | [ChaincodeInstanceEventsStreamRequest](#cckit.gateway.ChaincodeInstanceEventsStreamRequest) | [.protos.ChaincodeEvent](#protos.ChaincodeEvent) stream | Chaincode events stream |
| Events | [ChaincodeInstanceEventsRequest](#cckit.gateway.ChaincodeInstanceEventsRequest) | [ChaincodeEvents](#cckit.gateway.ChaincodeEvents) | Chaincode events |


<a name="cckit.gateway.ChaincodeService"></a>

### ChaincodeService
Chaincode communication service. Allows to locate channel/chaincode.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exec | [ChaincodeExec](#cckit.gateway.ChaincodeExec) | [.protos.Response](#protos.Response) | Exec: Query or Invoke |
| Query | [ChaincodeInput](#cckit.gateway.ChaincodeInput) | [.protos.Response](#protos.Response) | Query chaincode on home peer. Do NOT send to orderer. |
| Invoke | [ChaincodeInput](#cckit.gateway.ChaincodeInput) | [.protos.Response](#protos.Response) | Invoke chaincode on peers, according to endorsement policy and the SEND to orderer |
| EventsStream | [ChaincodeEventsStreamRequest](#cckit.gateway.ChaincodeEventsStreamRequest) | [ChaincodeEvent](#cckit.gateway.ChaincodeEvent) stream | Chaincode events stream |
| Events | [ChaincodeEventsRequest](#cckit.gateway.ChaincodeEventsRequest) | [ChaincodeEvents](#cckit.gateway.ChaincodeEvents) | Chaincode events |

 



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

