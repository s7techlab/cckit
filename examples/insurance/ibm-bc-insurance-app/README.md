# Tests for IBM Blockchain insurance app with CCKit

This example shows how to apply concepts of test-driven development to writing chaincode in Golang for Hyperledger Fabric 
using shim MockStub to unit-test your chaincode without having to deploy it in a Blockchain network

 
##  IBM blockchain insurance app  

Blockchain insurance application example https://developer.ibm.com/code/patterns/build-a-blockchain-insurance-app/  
shows how through distributed ledger and smart contracts blockchain can transform insurance processes:
 
>Blockchain presents a huge opportunity for the insurance industry. It offers the chance to innovate around the way
data is exchanged, claims are processed, and fraud is prevented. Blockchain can bring together developers from tech 
companies, regulators, and insurance companies to create a valuable new insurance management asset. 
 
![Architecture](images/arch-blockchain-insurance2.png)


### Source code 

https://github.com/IBM/build-blockchain-insurance-app

Chaincode from  https://github.com/IBM/build-blockchain-insurance-app/tree/master/web/chaincode/src/bcins is copied to  the current
directory for sample creation  test


## Why test with Mockstub

Creating blockchain network and deploying the chaincode(s) to it is quite cumbersome and slow process, especially if your code
is constantly changing during developing process. Mockstub allows you to get the test results almost immediately.