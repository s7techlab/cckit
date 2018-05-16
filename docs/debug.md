# Hyperledger Fabric chaincode kit (CCKit)

## Chaincode debugging


Use logger in chaincode

`c.Logger().Debug( "debug message")`



Change logging level while performing tests

`CORE_CHAINCODE_LOGGING_LEVEL=debug go test`
