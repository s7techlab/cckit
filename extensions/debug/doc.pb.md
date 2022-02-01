# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [debug/debug_state.proto](#debug/debug_state.proto)
    - [CompositeKey](#extensions.debug.CompositeKey)
    - [CompositeKeys](#extensions.debug.CompositeKeys)
    - [Prefix](#extensions.debug.Prefix)
    - [Prefixes](#extensions.debug.Prefixes)
    - [PrefixesMatchCount](#extensions.debug.PrefixesMatchCount)
    - [PrefixesMatchCount.MatchesEntry](#extensions.debug.PrefixesMatchCount.MatchesEntry)
    - [Value](#extensions.debug.Value)
  
    - [DebugStateService](#extensions.debug.DebugStateService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="debug/debug_state.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## debug/debug_state.proto



<a name="extensions.debug.CompositeKey"></a>

### CompositeKey
State key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | repeated |  |






<a name="extensions.debug.CompositeKeys"></a>

### CompositeKeys
State keys


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| keys | [CompositeKey](#extensions.debug.CompositeKey) | repeated |  |






<a name="extensions.debug.Prefix"></a>

### Prefix
State key prefix


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | repeated | parts of key |






<a name="extensions.debug.Prefixes"></a>

### Prefixes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| prefixes | [Prefix](#extensions.debug.Prefix) | repeated |  |






<a name="extensions.debug.PrefixesMatchCount"></a>

### PrefixesMatchCount
State key prefix match count


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| matches | [PrefixesMatchCount.MatchesEntry](#extensions.debug.PrefixesMatchCount.MatchesEntry) | repeated |  |






<a name="extensions.debug.PrefixesMatchCount.MatchesEntry"></a>

### PrefixesMatchCount.MatchesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [uint32](#uint32) |  |  |






<a name="extensions.debug.Value"></a>

### Value
State value


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | repeated |  |
| value | [bytes](#bytes) |  |  |
| json | [string](#string) |  |  |





 

 

 


<a name="extensions.debug.DebugStateService"></a>

### DebugStateService
Debug state service
allows to directly manage chaincode state

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListKeys | [Prefix](#extensions.debug.Prefix) | [CompositeKeys](#extensions.debug.CompositeKeys) | Get keys list, returns all keys or, if prefixes are defined, only prefix matched |
| GetState | [CompositeKey](#extensions.debug.CompositeKey) | [Value](#extensions.debug.Value) | Get state value by key |
| PutState | [Value](#extensions.debug.Value) | [Value](#extensions.debug.Value) | Put state value |
| DeleteState | [CompositeKey](#extensions.debug.CompositeKey) | [Value](#extensions.debug.Value) | Delete state value |
| DeleteStates | [Prefixes](#extensions.debug.Prefixes) | [PrefixesMatchCount](#extensions.debug.PrefixesMatchCount) | Delete all states or, if prefixes are defined, only prefix matched |

 



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

