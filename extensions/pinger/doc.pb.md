

<!-- generated from chaincode_pinger.proto -->
# 


<a name="extensions.pinger.ChaincodePingerService"></a>

## Методы сервиса ChaincodePingerService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Ping | [.google.protobuf.Empty](#google.protobuf.Empty) | [PingInfo](#extensions.pinger.PingInfo) | ping chaincode |

 <!-- end services -->

## Структура данных


<a name="extensions.pinger.PingInfo"></a>

### PingInfo
stores time and certificate of ping tx creator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| invoker_id | [string](#string) |  |  |
| invoker_cert | [bytes](#bytes) |  |  |
| time | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->


<!-- end enums -->

 <!-- end HasExtensions -->



