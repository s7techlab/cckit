# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [cpaper_extended/schema/payload.proto](#cpaper_extended/schema/payload.proto)
    - [BuyCommercialPaper](#examples.cpaper_extended.schema.BuyCommercialPaper)
    - [IssueCommercialPaper](#examples.cpaper_extended.schema.IssueCommercialPaper)
    - [RedeemCommercialPaper](#examples.cpaper_extended.schema.RedeemCommercialPaper)
  
- [cpaper_extended/schema/state.proto](#cpaper_extended/schema/state.proto)
    - [CommercialPaper](#examples.cpaper_extended.schema.CommercialPaper)
    - [CommercialPaperId](#examples.cpaper_extended.schema.CommercialPaperId)
    - [CommercialPaperList](#examples.cpaper_extended.schema.CommercialPaperList)
  
    - [CommercialPaper.State](#examples.cpaper_extended.schema.CommercialPaper.State)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cpaper_extended/schema/payload.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cpaper_extended/schema/payload.proto



<a name="examples.cpaper_extended.schema.BuyCommercialPaper"></a>

### BuyCommercialPaper
BuyCommercialPaper event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |
| current_owner | [string](#string) |  |  |
| new_owner | [string](#string) |  |  |
| price | [int32](#int32) |  |  |
| purchase_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.cpaper_extended.schema.IssueCommercialPaper"></a>

### IssueCommercialPaper
IssueCommercialPaper event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |
| issue_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| maturity_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| face_value | [int32](#int32) |  |  |
| external_id | [string](#string) |  | external_id - another unique constraint |






<a name="examples.cpaper_extended.schema.RedeemCommercialPaper"></a>

### RedeemCommercialPaper
RedeemCommercialPaper event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |
| redeeming_owner | [string](#string) |  |  |
| redeem_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 

 

 

 



<a name="cpaper_extended/schema/state.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cpaper_extended/schema/state.proto



<a name="examples.cpaper_extended.schema.CommercialPaper"></a>

### CommercialPaper
Commercial Paper state entry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  | Issuer and Paper number comprises composite primary key of Commercial paper entry |
| paper_number | [string](#string) |  |  |
| owner | [string](#string) |  |  |
| issue_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| maturity_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| face_value | [int32](#int32) |  |  |
| state | [CommercialPaper.State](#examples.cpaper_extended.schema.CommercialPaper.State) |  |  |
| external_id | [string](#string) |  | Additional unique field for entry |






<a name="examples.cpaper_extended.schema.CommercialPaperId"></a>

### CommercialPaperId
CommercialPaperId identifier part


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |






<a name="examples.cpaper_extended.schema.CommercialPaperList"></a>

### CommercialPaperList
Container for returning multiple entities


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [CommercialPaper](#examples.cpaper_extended.schema.CommercialPaper) | repeated |  |





 


<a name="examples.cpaper_extended.schema.CommercialPaper.State"></a>

### CommercialPaper.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_ISSUED | 0 |  |
| STATE_TRADING | 1 |  |
| STATE_REDEEMED | 2 |  |


 

 

 



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

