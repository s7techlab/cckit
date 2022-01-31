# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [cpaper_asservice/cpaper.proto](#cpaper_asservice/cpaper.proto)
    - [BuyCommercialPaper](#examples.cpaper_asservice.BuyCommercialPaper)
    - [CommercialPaper](#examples.cpaper_asservice.CommercialPaper)
    - [CommercialPaperId](#examples.cpaper_asservice.CommercialPaperId)
    - [CommercialPaperList](#examples.cpaper_asservice.CommercialPaperList)
    - [ExternalId](#examples.cpaper_asservice.ExternalId)
    - [IssueCommercialPaper](#examples.cpaper_asservice.IssueCommercialPaper)
    - [RedeemCommercialPaper](#examples.cpaper_asservice.RedeemCommercialPaper)
  
    - [CommercialPaper.State](#examples.cpaper_asservice.CommercialPaper.State)
  
    - [CPaperService](#examples.cpaper_asservice.CPaperService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cpaper_asservice/cpaper.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cpaper_asservice/cpaper.proto



<a name="examples.cpaper_asservice.BuyCommercialPaper"></a>

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






<a name="examples.cpaper_asservice.CommercialPaper"></a>

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
| state | [CommercialPaper.State](#examples.cpaper_asservice.CommercialPaper.State) |  |  |
| external_id | [string](#string) |  | Additional unique field for entry |






<a name="examples.cpaper_asservice.CommercialPaperId"></a>

### CommercialPaperId
CommercialPaperId identifier part


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |






<a name="examples.cpaper_asservice.CommercialPaperList"></a>

### CommercialPaperList
Container for returning multiple entities


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | repeated |  |






<a name="examples.cpaper_asservice.ExternalId"></a>

### ExternalId
ExternalId


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="examples.cpaper_asservice.IssueCommercialPaper"></a>

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






<a name="examples.cpaper_asservice.RedeemCommercialPaper"></a>

### RedeemCommercialPaper
RedeemCommercialPaper event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| issuer | [string](#string) |  |  |
| paper_number | [string](#string) |  |  |
| redeeming_owner | [string](#string) |  |  |
| redeem_date | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 


<a name="examples.cpaper_asservice.CommercialPaper.State"></a>

### CommercialPaper.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_ISSUED | 0 |  |
| STATE_TRADING | 1 |  |
| STATE_REDEEMED | 2 |  |


 

 


<a name="examples.cpaper_asservice.CPaperService"></a>

### CPaperService
Commercial paper chaincode-as-service

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| List | [.google.protobuf.Empty](#google.protobuf.Empty) | [CommercialPaperList](#examples.cpaper_asservice.CommercialPaperList) | List method returns all registered commercial papers |
| Get | [CommercialPaperId](#examples.cpaper_asservice.CommercialPaperId) | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | Get method returns commercial paper data by id |
| GetByExternalId | [ExternalId](#examples.cpaper_asservice.ExternalId) | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | GetByExternalId |
| Issue | [IssueCommercialPaper](#examples.cpaper_asservice.IssueCommercialPaper) | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | Issue commercial paper |
| Buy | [BuyCommercialPaper](#examples.cpaper_asservice.BuyCommercialPaper) | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | Buy commercial paper |
| Redeem | [RedeemCommercialPaper](#examples.cpaper_asservice.RedeemCommercialPaper) | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | Redeem commercial paper |
| Delete | [CommercialPaperId](#examples.cpaper_asservice.CommercialPaperId) | [CommercialPaper](#examples.cpaper_asservice.CommercialPaper) | Delete commercial paper |

 



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

