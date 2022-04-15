# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [fabcar.proto](#fabcar.proto)
    - [Car](#examples.fabcar.Car)
    - [CarCreated](#examples.fabcar.CarCreated)
    - [CarDeleted](#examples.fabcar.CarDeleted)
    - [CarDetail](#examples.fabcar.CarDetail)
    - [CarDetailDeleted](#examples.fabcar.CarDetailDeleted)
    - [CarDetailId](#examples.fabcar.CarDetailId)
    - [CarDetails](#examples.fabcar.CarDetails)
    - [CarDetailsUpdated](#examples.fabcar.CarDetailsUpdated)
    - [CarId](#examples.fabcar.CarId)
    - [CarOwner](#examples.fabcar.CarOwner)
    - [CarOwnerDeleted](#examples.fabcar.CarOwnerDeleted)
    - [CarOwnerId](#examples.fabcar.CarOwnerId)
    - [CarOwners](#examples.fabcar.CarOwners)
    - [CarOwnersUpdated](#examples.fabcar.CarOwnersUpdated)
    - [CarUpdated](#examples.fabcar.CarUpdated)
    - [CarView](#examples.fabcar.CarView)
    - [Cars](#examples.fabcar.Cars)
    - [CreateCarRequest](#examples.fabcar.CreateCarRequest)
    - [CreateMakerRequest](#examples.fabcar.CreateMakerRequest)
    - [Maker](#examples.fabcar.Maker)
    - [MakerCreated](#examples.fabcar.MakerCreated)
    - [MakerDeleted](#examples.fabcar.MakerDeleted)
    - [MakerName](#examples.fabcar.MakerName)
    - [Makers](#examples.fabcar.Makers)
    - [SetCarDetail](#examples.fabcar.SetCarDetail)
    - [SetCarOwner](#examples.fabcar.SetCarOwner)
    - [UpdateCarDetailsRequest](#examples.fabcar.UpdateCarDetailsRequest)
    - [UpdateCarOwnersRequest](#examples.fabcar.UpdateCarOwnersRequest)
    - [UpdateCarRequest](#examples.fabcar.UpdateCarRequest)
  
    - [DetailType](#examples.fabcar.DetailType)
  
  
    - [FabCarService](#examples.fabcar.FabCarService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="fabcar.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fabcar.proto



<a name="examples.fabcar.Car"></a>

### Car



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| make | [string](#string) |  |  |
| model | [string](#string) |  |  |
| colour | [string](#string) |  |  |
| number | [uint64](#uint64) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.fabcar.CarCreated"></a>

### CarCreated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| make | [string](#string) |  |  |
| model | [string](#string) |  |  |
| colour | [string](#string) |  |  |
| number | [uint64](#uint64) |  |  |






<a name="examples.fabcar.CarDeleted"></a>

### CarDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| make | [string](#string) |  |  |
| model | [string](#string) |  |  |
| colour | [string](#string) |  |  |
| number | [uint64](#uint64) |  |  |
| owners | [CarOwners](#examples.fabcar.CarOwners) |  |  |
| details | [CarDetails](#examples.fabcar.CarDetails) |  |  |






<a name="examples.fabcar.CarDetail"></a>

### CarDetail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car_id | [string](#string) | repeated |  |
| type | [DetailType](#examples.fabcar.DetailType) |  |  |
| make | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.fabcar.CarDetailDeleted"></a>

### CarDetailDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| detail | [CarDetail](#examples.fabcar.CarDetail) |  |  |






<a name="examples.fabcar.CarDetailId"></a>

### CarDetailId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car_id | [string](#string) | repeated |  |
| type | [DetailType](#examples.fabcar.DetailType) |  |  |






<a name="examples.fabcar.CarDetails"></a>

### CarDetails



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [CarDetail](#examples.fabcar.CarDetail) | repeated |  |






<a name="examples.fabcar.CarDetailsUpdated"></a>

### CarDetailsUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| details | [CarDetails](#examples.fabcar.CarDetails) |  |  |






<a name="examples.fabcar.CarId"></a>

### CarId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |






<a name="examples.fabcar.CarOwner"></a>

### CarOwner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car_id | [string](#string) | repeated |  |
| first_name | [string](#string) |  |  |
| second_name | [string](#string) |  |  |
| vehicle_passport | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.fabcar.CarOwnerDeleted"></a>

### CarOwnerDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner | [CarOwner](#examples.fabcar.CarOwner) |  |  |






<a name="examples.fabcar.CarOwnerId"></a>

### CarOwnerId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car_id | [string](#string) | repeated |  |
| first_name | [string](#string) |  |  |
| second_name | [string](#string) |  |  |






<a name="examples.fabcar.CarOwners"></a>

### CarOwners



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [CarOwner](#examples.fabcar.CarOwner) | repeated |  |






<a name="examples.fabcar.CarOwnersUpdated"></a>

### CarOwnersUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owners | [CarOwners](#examples.fabcar.CarOwners) |  |  |






<a name="examples.fabcar.CarUpdated"></a>

### CarUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| colour | [string](#string) |  |  |






<a name="examples.fabcar.CarView"></a>

### CarView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car | [Car](#examples.fabcar.Car) |  |  |
| owners | [CarOwners](#examples.fabcar.CarOwners) |  |  |
| details | [CarDetails](#examples.fabcar.CarDetails) |  |  |






<a name="examples.fabcar.Cars"></a>

### Cars



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Car](#examples.fabcar.Car) | repeated |  |






<a name="examples.fabcar.CreateCarRequest"></a>

### CreateCarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| make | [string](#string) |  |  |
| model | [string](#string) |  |  |
| colour | [string](#string) |  |  |
| number | [uint64](#uint64) |  |  |
| owners | [SetCarOwner](#examples.fabcar.SetCarOwner) | repeated |  |
| details | [SetCarDetail](#examples.fabcar.SetCarDetail) | repeated |  |






<a name="examples.fabcar.CreateMakerRequest"></a>

### CreateMakerRequest
Entities


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| country | [string](#string) |  |  |
| foundation_year | [uint64](#uint64) |  | in 1886 was founded the oldest automaker - Mercedes-Benz |






<a name="examples.fabcar.Maker"></a>

### Maker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| country | [string](#string) |  |  |
| foundation_year | [uint64](#uint64) |  |  |






<a name="examples.fabcar.MakerCreated"></a>

### MakerCreated
Events


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| country | [string](#string) |  |  |
| foundation_year | [uint64](#uint64) |  |  |






<a name="examples.fabcar.MakerDeleted"></a>

### MakerDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| country | [string](#string) |  |  |
| foundation_year | [uint64](#uint64) |  |  |






<a name="examples.fabcar.MakerName"></a>

### MakerName



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="examples.fabcar.Makers"></a>

### Makers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Maker](#examples.fabcar.Maker) | repeated |  |






<a name="examples.fabcar.SetCarDetail"></a>

### SetCarDetail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [DetailType](#examples.fabcar.DetailType) |  |  |
| make | [string](#string) |  |  |






<a name="examples.fabcar.SetCarOwner"></a>

### SetCarOwner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| first_name | [string](#string) |  |  |
| second_name | [string](#string) |  |  |
| vehicle_passport | [string](#string) |  |  |






<a name="examples.fabcar.UpdateCarDetailsRequest"></a>

### UpdateCarDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car_id | [string](#string) | repeated |  |
| details | [SetCarDetail](#examples.fabcar.SetCarDetail) | repeated |  |






<a name="examples.fabcar.UpdateCarOwnersRequest"></a>

### UpdateCarOwnersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| car_id | [string](#string) | repeated |  |
| owners | [SetCarOwner](#examples.fabcar.SetCarOwner) | repeated |  |






<a name="examples.fabcar.UpdateCarRequest"></a>

### UpdateCarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| color | [string](#string) |  |  |
| owners | [SetCarOwner](#examples.fabcar.SetCarOwner) | repeated |  |
| details | [SetCarDetail](#examples.fabcar.SetCarDetail) | repeated |  |





 


<a name="examples.fabcar.DetailType"></a>

### DetailType
Dictionaries

| Name | Number | Description |
| ---- | ------ | ----------- |
| WHEELS | 0 |  |
| BATTERY | 1 |  |


 

 


<a name="examples.fabcar.FabCarService"></a>

### FabCarService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateMaker | [CreateMakerRequest](#examples.fabcar.CreateMakerRequest) | [Maker](#examples.fabcar.Maker) |  |
| DeleteMaker | [MakerName](#examples.fabcar.MakerName) | [Maker](#examples.fabcar.Maker) |  |
| GetMaker | [MakerName](#examples.fabcar.MakerName) | [Maker](#examples.fabcar.Maker) |  |
| ListMakers | [.google.protobuf.Empty](#google.protobuf.Empty) | [Makers](#examples.fabcar.Makers) |  |
| CreateCar | [CreateCarRequest](#examples.fabcar.CreateCarRequest) | [CarView](#examples.fabcar.CarView) |  |
| UpdateCar | [UpdateCarRequest](#examples.fabcar.UpdateCarRequest) | [CarView](#examples.fabcar.CarView) |  |
| DeleteCar | [CarId](#examples.fabcar.CarId) | [CarView](#examples.fabcar.CarView) |  |
| GetCar | [CarId](#examples.fabcar.CarId) | [Car](#examples.fabcar.Car) |  |
| GetCarView | [CarId](#examples.fabcar.CarId) | [CarView](#examples.fabcar.CarView) |  |
| ListCars | [.google.protobuf.Empty](#google.protobuf.Empty) | [Cars](#examples.fabcar.Cars) |  |
| UpdateCarOwners | [UpdateCarOwnersRequest](#examples.fabcar.UpdateCarOwnersRequest) | [CarOwners](#examples.fabcar.CarOwners) |  |
| DeleteCarOwner | [CarOwnerId](#examples.fabcar.CarOwnerId) | [CarOwner](#examples.fabcar.CarOwner) |  |
| GetCarOwner | [CarOwnerId](#examples.fabcar.CarOwnerId) | [CarOwner](#examples.fabcar.CarOwner) |  |
| ListCarOwners | [CarId](#examples.fabcar.CarId) | [CarOwners](#examples.fabcar.CarOwners) |  |
| UpdateCarDetails | [UpdateCarDetailsRequest](#examples.fabcar.UpdateCarDetailsRequest) | [CarDetails](#examples.fabcar.CarDetails) |  |
| DeleteCarDetail | [CarDetailId](#examples.fabcar.CarDetailId) | [CarDetail](#examples.fabcar.CarDetail) |  |
| GetCarDetail | [CarDetailId](#examples.fabcar.CarDetailId) | [CarDetail](#examples.fabcar.CarDetail) |  |
| ListCarDetails | [CarId](#examples.fabcar.CarId) | [CarDetails](#examples.fabcar.CarDetails) |  |

 



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

