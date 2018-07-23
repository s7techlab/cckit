# Tests for IBM Blockchain insurance app with CCKit

This example shows how to test Hyperledger Fabric chaincode using shim Mockstub

##  IBM blockchain insurance app  

Blockchain insurance application example https://developer.ibm.com/code/patterns/build-a-blockchain-insurance-app/  
shows how through distributed ledger and smart contracts blockchain can transform  insurance processes:

>With its distributed ledger, smart contracts, and non-repudiation capabilities, blockchain is revolutionizing the way 
financial organizations do business, and the insurance industry is no exception. 
This code pattern shows you how to implement a web-based blockchain app using Hyperledger Fabric to facilitate insurance 
sales and claims.


### Source code 

https://github.com/IBM/build-blockchain-insurance-app

Chaincode from  https://github.com/IBM/build-blockchain-insurance-app/tree/master/web/chaincode/src/bcins is copied to  the current
directory for sample creation  test


## Tests

Here we show how to use shim MockStub to unit-test your chaincode without having to deploy it in a Blockchain network