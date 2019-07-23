# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [schema.proto](#schema.proto)
    - [BuyCommercialPaper](#schema.BuyCommercialPaper)
    - [CommercialPaper](#schema.CommercialPaper)
    - [CommercialPaperId](#schema.CommercialPaperId)
    - [CommercialPaperList](#schema.CommercialPaperList)
    - [ExternalId](#schema.ExternalId)
    - [IssueCommercialPaper](#schema.IssueCommercialPaper)
    - [RedeemCommercialPaper](#schema.RedeemCommercialPaper)
  
    - [CommercialPaper.State](#schema.CommercialPaper.State)
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="schema.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema.proto



<a name="schema.BuyCommercialPaper"></a>

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






<a name="schema.CommercialPaper"></a>

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
| state | [CommercialPaper.State](#schema.CommercialPaper.State) |  |  |
| external_id | [string](#string) |  | Additional unique field for entry |






<a name="schema.CommercialPaperId"></a>

### CommercialPaperId
CommercialPaperId identifier part


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |






<a name="schema.CommercialPaperList"></a>

### CommercialPaperList
Container for returning multiple entities


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [CommercialPaper](#schema.CommercialPaper) | repeated |  |






<a name="schema.ExternalId"></a>

### ExternalId
ExternalId


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="schema.IssueCommercialPaper"></a>

### IssueCommercialPaper
IssueCommercialPaper event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |
| issue_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| maturity_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| face_value | [int32](#int32) |  |  |
| external_id | [string](#string) |  | external_id - once more uniq id of state entry |






<a name="schema.RedeemCommercialPaper"></a>

### RedeemCommercialPaper
RedeemCommercialPaper event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |
| redeeming_owner | [string](#string) |  |  |
| redeem_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 


<a name="schema.CommercialPaper.State"></a>

### CommercialPaper.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| ISSUED | 0 |  |
| TRADING | 1 |  |
| REDEEMED | 2 |  |


 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

