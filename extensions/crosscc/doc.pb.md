# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [crosscc/cclocator_setting.proto](#crosscc/cclocator_setting.proto)
    - [PingServiceResponse](#crosscc.PingServiceResponse)
    - [PingServiceResponses](#crosscc.PingServiceResponses)
    - [ServiceLocator](#crosscc.ServiceLocator)
    - [ServiceLocatorId](#crosscc.ServiceLocatorId)
    - [ServiceLocatorSet](#crosscc.ServiceLocatorSet)
    - [ServiceLocatorSetRequest](#crosscc.ServiceLocatorSetRequest)
    - [ServiceLocators](#crosscc.ServiceLocators)
  
    - [SettingService](#crosscc.SettingService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="crosscc/cclocator_setting.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## crosscc/cclocator_setting.proto



<a name="crosscc.PingServiceResponse"></a>

### PingServiceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| locator | [ServiceLocator](#crosscc.ServiceLocator) |  |  |
| error | [string](#string) |  |  |






<a name="crosscc.PingServiceResponses"></a>

### PingServiceResponses



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| responses | [PingServiceResponse](#crosscc.PingServiceResponse) | repeated |  |






<a name="crosscc.ServiceLocator"></a>

### ServiceLocator
State: ervice resolving setting


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [string](#string) |  | service identifier |
| channel | [string](#string) |  | channel id |
| chaincode | [string](#string) |  | chaincode name |






<a name="crosscc.ServiceLocatorId"></a>

### ServiceLocatorId
Id: service resolving setting identifier


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [string](#string) |  | service identifier |






<a name="crosscc.ServiceLocatorSet"></a>

### ServiceLocatorSet
Event: service resolving settings was set


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [string](#string) |  | service identifier |
| channel | [string](#string) |  | channel id |
| chaincode | [string](#string) |  | chaincode name |






<a name="crosscc.ServiceLocatorSetRequest"></a>

### ServiceLocatorSetRequest
Request: set service resolving setting


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [string](#string) |  | service identifier |
| channel | [string](#string) |  | channel id |
| chaincode | [string](#string) |  | chaincode name |






<a name="crosscc.ServiceLocators"></a>

### ServiceLocators
List: service resolving settings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ServiceLocator](#crosscc.ServiceLocator) | repeated |  |





 

 

 


<a name="crosscc.SettingService"></a>

### SettingService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ServiceLocatorSet | [ServiceLocatorSetRequest](#crosscc.ServiceLocatorSetRequest) | [ServiceLocator](#crosscc.ServiceLocator) |  |
| ServiceLocatorGet | [ServiceLocatorId](#crosscc.ServiceLocatorId) | [ServiceLocator](#crosscc.ServiceLocator) |  |
| ListServiceLocators | [.google.protobuf.Empty](#google.protobuf.Empty) | [ServiceLocators](#crosscc.ServiceLocators) |  |
| PingService | [ServiceLocatorId](#crosscc.ServiceLocatorId) | [PingServiceResponse](#crosscc.PingServiceResponse) | Try to query chaincodes from service chaincode settings |
| PingServices | [.google.protobuf.Empty](#google.protobuf.Empty) | [PingServiceResponses](#crosscc.PingServiceResponses) |  |

 



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

