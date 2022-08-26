# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [token/service/allowance/allowance.proto](#token/service/allowance/allowance.proto)
    - [Allowance](#examples.erc20_service.service.allowance.Allowance)
    - [AllowanceId](#examples.erc20_service.service.allowance.AllowanceId)
    - [AllowanceRequest](#examples.erc20_service.service.allowance.AllowanceRequest)
    - [Allowances](#examples.erc20_service.service.allowance.Allowances)
    - [ApproveRequest](#examples.erc20_service.service.allowance.ApproveRequest)
    - [Approved](#examples.erc20_service.service.allowance.Approved)
    - [TransferFromRequest](#examples.erc20_service.service.allowance.TransferFromRequest)
    - [TransferFromResponse](#examples.erc20_service.service.allowance.TransferFromResponse)
    - [TransferredFrom](#examples.erc20_service.service.allowance.TransferredFrom)
  
  
  
    - [AllowanceService](#examples.erc20_service.service.allowance.AllowanceService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="token/service/allowance/allowance.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## token/service/allowance/allowance.proto



<a name="examples.erc20_service.service.allowance.Allowance"></a>

### Allowance
Allowance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| spender_address | [string](#string) |  |  |
| token | [string](#string) | repeated |  |
| amount | [uint64](#uint64) |  |  |






<a name="examples.erc20_service.service.allowance.AllowanceId"></a>

### AllowanceId
Allowance identifier


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| spender_address | [string](#string) |  |  |
| token | [string](#string) | repeated |  |






<a name="examples.erc20_service.service.allowance.AllowanceRequest"></a>

### AllowanceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| spender_address | [string](#string) |  |  |
| token | [string](#string) | repeated |  |






<a name="examples.erc20_service.service.allowance.Allowances"></a>

### Allowances



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Allowance](#examples.erc20_service.service.allowance.Allowance) | repeated |  |






<a name="examples.erc20_service.service.allowance.ApproveRequest"></a>

### ApproveRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| spender_address | [string](#string) |  |  |
| amount | [uint64](#uint64) |  |  |
| token | [string](#string) | repeated |  |






<a name="examples.erc20_service.service.allowance.Approved"></a>

### Approved
Approved event is emitted when Approve method has been invoked


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| spender_address | [string](#string) |  |  |
| token | [string](#string) | repeated |  |
| amount | [uint64](#uint64) |  |  |






<a name="examples.erc20_service.service.allowance.TransferFromRequest"></a>

### TransferFromRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| recipient_address | [string](#string) |  |  |
| amount | [uint64](#uint64) |  |  |
| token | [string](#string) | repeated |  |






<a name="examples.erc20_service.service.allowance.TransferFromResponse"></a>

### TransferFromResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| recipient_address | [string](#string) |  |  |
| token | [string](#string) | repeated |  |
| amount | [uint64](#uint64) |  |  |






<a name="examples.erc20_service.service.allowance.TransferredFrom"></a>

### TransferredFrom
TransferredFrom event is emitted when TransferFrom method has been invoked


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner_address | [string](#string) |  |  |
| spender_address | [string](#string) |  |  |
| recipient_address | [string](#string) |  |  |
| token | [string](#string) | repeated |  |
| amount | [uint64](#uint64) |  |  |





 

 

 


<a name="examples.erc20_service.service.allowance.AllowanceService"></a>

### AllowanceService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAllowance | [AllowanceRequest](#examples.erc20_service.service.allowance.AllowanceRequest) | [Allowance](#examples.erc20_service.service.allowance.Allowance) | Returns the remaining number of tokens that spender will be allowed to spend on behalf of owner through transfersender. This is zero by default. |
| Approve | [ApproveRequest](#examples.erc20_service.service.allowance.ApproveRequest) | [Allowance](#examples.erc20_service.service.allowance.Allowance) | Sets amount as the allowance of spender over the caller’s tokens. Emits an ApprovalEvent |
| TransferFrom | [TransferFromRequest](#examples.erc20_service.service.allowance.TransferFromRequest) | [TransferFromResponse](#examples.erc20_service.service.allowance.TransferFromResponse) | Moves amount tokens from sender to recipient using the allowance mechanism. Amount is then deducted from the caller’s allowance. Emits TransferEvent |

 



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

