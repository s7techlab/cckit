# Hyperledger Fabric chaincode kit (CCKit)

## Chaincode testing tools

### Mockstub

Deploying a Hyperledger Fabric blockchain network, chancide installing and initializing is quite complicated to set up and
a long procedure. 

The time to re-install / upgrade the code of a smart contract can be reduced by using 
[chaincode dev mode](https://hyperledger-fabric.readthedocs.io/en/latest/peer-chaincode-devmode.html). Normally chaincodes 
are started and maintained by peer. In “dev” mode, chaincode is built and started by the user. 
This mode is useful during chaincode development phase for rapid code/build/run/debug cycle turnaround. However, the process 
of updating the code will still be slow.

The [shim](https://github.com/hyperledger/fabric/tree/master/core/chaincode/shim) package contains a 
[MockStub](https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/mockstub.go) implementation 
that wraps calls to a chaincode, simulating its behavior in the HLF peer environment. MockStub don't need to start 
multiple docker containers with peer, world state, chaincodes
and allows to get test results almost immediately. MockStub essentially replaces the SDK and peer enviroment, allowing 
to test calls to the chaincode functions

![mockstub](../docs/img/mockstub-hlf-peer.png)