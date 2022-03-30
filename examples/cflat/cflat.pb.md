

<!-- generated from cflat.proto -->
# 


<a name="examples.cflat.CFlatService"></a>

## Методы сервиса CFlatService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateFlat | [CreateFlatRequest](#examples.cflat.CreateFlatRequest) | [FlatView](#examples.cflat.FlatView) |  |
| UpdateFlat | [UpdateFlatRequest](#examples.cflat.UpdateFlatRequest) | [FlatView](#examples.cflat.FlatView) |  |
| DeleteFlat | [FlatId](#examples.cflat.FlatId) | [FlatView](#examples.cflat.FlatView) |  |
| GetFlat | [FlatId](#examples.cflat.FlatId) | [Flat](#examples.cflat.Flat) |  |
| GetFlatView | [FlatId](#examples.cflat.FlatId) | [FlatView](#examples.cflat.FlatView) |  |
| ListFlats | [.google.protobuf.Empty](#google.protobuf.Empty) | [Flats](#examples.cflat.Flats) |  |
| UpdateFlatResident | [UpdateFlatResidentRequest](#examples.cflat.UpdateFlatResidentRequest) | [FlatResident](#examples.cflat.FlatResident) |  |
| DeleteFlatResident | [FlatResidentId](#examples.cflat.FlatResidentId) | [FlatResident](#examples.cflat.FlatResident) |  |
| GetFlatResident | [FlatResidentId](#examples.cflat.FlatResidentId) | [FlatResident](#examples.cflat.FlatResident) |  |
| ListFlatResidents | [FlatId](#examples.cflat.FlatId) | [FlatResidents](#examples.cflat.FlatResidents) |  |
| UpdateFlatRoom | [UpdateFlatRoomRequest](#examples.cflat.UpdateFlatRoomRequest) | [FlatRoom](#examples.cflat.FlatRoom) |  |
| DeleteFlatRoom | [FlatRoomId](#examples.cflat.FlatRoomId) | [FlatRoom](#examples.cflat.FlatRoom) |  |
| GetFlatRoom | [FlatRoomId](#examples.cflat.FlatRoomId) | [FlatRoom](#examples.cflat.FlatRoom) |  |
| ListFlatRooms | [FlatId](#examples.cflat.FlatId) | [FlatRooms](#examples.cflat.FlatRooms) |  |

 <!-- end services -->

## Структура данных


<a name="examples.cflat.CreateFlatRequest"></a>

### CreateFlatRequest
Entities


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| country | [string](#string) |  |  |
| region | [string](#string) |  |  |
| city | [string](#string) |  |  |
| street | [string](#string) |  |  |
| house_num | [uint64](#uint64) |  |  |
| flat_num | [uint64](#uint64) |  |  |
| type | [PropertyType](#examples.cflat.PropertyType) |  |  |
| residents | [SetFlatResident](#examples.cflat.SetFlatResident) | repeated |  |
| rooms | [SetFlatRoom](#examples.cflat.SetFlatRoom) | repeated |  |






<a name="examples.cflat.Flat"></a>

### Flat



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| country | [string](#string) |  |  |
| region | [string](#string) |  |  |
| city | [string](#string) |  |  |
| street | [string](#string) |  |  |
| house_num | [uint64](#uint64) |  |  |
| flat_num | [uint64](#uint64) |  |  |
| type | [PropertyType](#examples.cflat.PropertyType) |  |  |
| residents_quantity | [uint64](#uint64) |  |  |
| area | [uint64](#uint64) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.cflat.FlatCreated"></a>

### FlatCreated
Events


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| country | [string](#string) |  |  |
| region | [string](#string) |  |  |
| city | [string](#string) |  |  |
| street | [string](#string) |  |  |
| house_num | [uint64](#uint64) |  |  |
| flat_num | [uint64](#uint64) |  |  |
| type | [PropertyType](#examples.cflat.PropertyType) |  |  |
| residents_quantity | [uint64](#uint64) |  |  |
| area | [uint64](#uint64) |  |  |






<a name="examples.cflat.FlatDeleted"></a>

### FlatDeleted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| country | [string](#string) |  |  |
| region | [string](#string) |  |  |
| city | [string](#string) |  |  |
| street | [string](#string) |  |  |
| house_num | [uint64](#uint64) |  |  |
| flat_num | [uint64](#uint64) |  |  |
| type | [PropertyType](#examples.cflat.PropertyType) |  |  |
| residents_quantity | [uint64](#uint64) |  |  |
| area | [uint64](#uint64) |  |  |
| residents | [FlatResidents](#examples.cflat.FlatResidents) |  |  |
| rooms | [FlatRooms](#examples.cflat.FlatRooms) |  |  |






<a name="examples.cflat.FlatId"></a>

### FlatId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |






<a name="examples.cflat.FlatResident"></a>

### FlatResident



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| first_name | [string](#string) |  | repeated string resident_id = 2 [(validator.field) = {repeated_count_min: 1}]; |
| second_name | [string](#string) |  |  |
| type | [AccommodationType](#examples.cflat.AccommodationType) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.cflat.FlatResidentId"></a>

### FlatResidentId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| resident_id | [string](#string) | repeated |  |






<a name="examples.cflat.FlatResidents"></a>

### FlatResidents



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| residents | [FlatResident](#examples.cflat.FlatResident) | repeated |  |






<a name="examples.cflat.FlatResidentsDeleted"></a>

### FlatResidentsDeleted







<a name="examples.cflat.FlatResidentsUpdated"></a>

### FlatResidentsUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| residents | [FlatResidents](#examples.cflat.FlatResidents) |  |  |






<a name="examples.cflat.FlatRoom"></a>

### FlatRoom



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| type | [RoomType](#examples.cflat.RoomType) |  | repeated string room_id = 2 [(validator.field) = {repeated_count_min: 1}]; |
| area | [uint64](#uint64) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="examples.cflat.FlatRoomId"></a>

### FlatRoomId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| room_id | [string](#string) | repeated |  |






<a name="examples.cflat.FlatRooms"></a>

### FlatRooms



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rooms | [FlatRoom](#examples.cflat.FlatRoom) | repeated |  |






<a name="examples.cflat.FlatRoomsDeleted"></a>

### FlatRoomsDeleted







<a name="examples.cflat.FlatRoomsUpdated"></a>

### FlatRoomsUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| rooms | [FlatRooms](#examples.cflat.FlatRooms) |  |  |






<a name="examples.cflat.FlatUpdated"></a>

### FlatUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| type | [PropertyType](#examples.cflat.PropertyType) |  |  |
| residents_quantity | [uint64](#uint64) |  |  |
| area | [uint64](#uint64) |  |  |






<a name="examples.cflat.FlatView"></a>

### FlatView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat | [Flat](#examples.cflat.Flat) |  |  |
| residents | [FlatResidents](#examples.cflat.FlatResidents) |  |  |
| rooms | [FlatRooms](#examples.cflat.FlatRooms) |  |  |






<a name="examples.cflat.Flats"></a>

### Flats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flats | [Flat](#examples.cflat.Flat) | repeated |  |






<a name="examples.cflat.SetFlatResident"></a>

### SetFlatResident



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| first_name | [string](#string) |  |  |
| second_name | [string](#string) |  |  |
| type | [AccommodationType](#examples.cflat.AccommodationType) |  |  |






<a name="examples.cflat.SetFlatRoom"></a>

### SetFlatRoom



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [RoomType](#examples.cflat.RoomType) |  |  |
| area | [uint64](#uint64) |  |  |






<a name="examples.cflat.UpdateFlatRequest"></a>

### UpdateFlatRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) | repeated |  |
| type | [PropertyType](#examples.cflat.PropertyType) |  |  |
| residents | [SetFlatResident](#examples.cflat.SetFlatResident) | repeated |  |
| rooms | [SetFlatRoom](#examples.cflat.SetFlatRoom) | repeated |  |






<a name="examples.cflat.UpdateFlatResidentRequest"></a>

### UpdateFlatResidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| residents | [SetFlatResident](#examples.cflat.SetFlatResident) | repeated |  |






<a name="examples.cflat.UpdateFlatRoomRequest"></a>

### UpdateFlatRoomRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flat_id | [string](#string) | repeated |  |
| rooms | [SetFlatRoom](#examples.cflat.SetFlatRoom) | repeated |  |





 <!-- end messages -->



## Словари


<a name="examples.cflat.AccommodationType"></a>

### AccommodationType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNREGISTERED | 0 |  |
| REGISTERED | 1 |  |
| RENT | 2 |  |



<a name="examples.cflat.PropertyType"></a>

### PropertyType
Dictionaries

| Name | Number | Description |
| ---- | ------ | ----------- |
| PRIVATE | 0 |  |
| STATE | 1 |  |



<a name="examples.cflat.RoomType"></a>

### RoomType


| Name | Number | Description |
| ---- | ------ | ----------- |
| KITCHEN | 0 |  |
| LIVING | 1 |  |
| DINING | 2 |  |
| BATHROOM | 3 |  |
| BEDROOM | 4 |  |
| HALLWAY | 5 |  |
| NURSERY | 6 |  |



<!-- end enums -->

 <!-- end HasExtensions -->



