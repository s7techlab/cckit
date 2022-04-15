# Pinger - hyperledger fabric chaincode ping extension

Often there is a need to find out if everything is correct with chaincode. To do it we would like to ping CC
and know some information.

CCKit provides `pinger` extension for implementing ping opportunity in Hyperledger Fabric chaincodes. When chaincode will be pinged,
you get information about invoker, his certificate and time.

Pinger implemented in two version:

1. As chaincode [handlers](pinger.go)
2. As [service](chaincode_pinger.proto), that can be embedded in chaincode, using [chaincode-as-service mode](../../gateway)