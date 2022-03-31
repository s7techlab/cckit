

<!-- generated from fabcar.proto -->
# 


<a name="examples.fabcar.FabCarService"></a>

## Методы сервиса FabCarService


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

 <!-- end services -->

## Структура данных


<a name="examples.fabcar.Car"></a>

### Car



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| make | [string](#string) |  |  |
| model | [string](#string) |  |  |
| colour | [string](#string) |  |  |
| number | [uint64](#uint64) |  |  |
| owners_quantity | [uint64](#uint64) |  |  |
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
| owners_quantity | [uint64](#uint64) |  |  |






<a name="examples.fabcar.CarDeleted"></a>

### CarDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| make | [string](#string) |  |  |
| model | [string](#string) |  |  |
| colour | [string](#string) |  |  |
| number | [uint64](#uint64) |  |  |
| owners_quantity | [uint64](#uint64) |  |  |
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
| number | [uint64](#uint64) |  |  |
| owners_quantity | [uint64](#uint64) |  |  |






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
| number | [uint64](#uint64) |  |  |
| owners | [SetCarOwner](#examples.fabcar.SetCarOwner) | repeated |  |
| details | [SetCarDetail](#examples.fabcar.SetCarDetail) | repeated |  |





 <!-- end messages -->



## Словари


<a name="examples.fabcar.DetailType"></a>

### DetailType
Dictionaries

| Name | Number | Description |
| ---- | ------ | ----------- |
| WHEELS | 0 |  |
| BATTERY | 1 |  |



<!-- end enums -->

 <!-- end HasExtensions -->



