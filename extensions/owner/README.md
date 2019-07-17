# Owner - access control hyperledger fabric chaincode extension

In many cases during chaincode instantiating we need to define permissions for chaincode functions -
"who is allowed to do this thing", incredibly important in the world of smart contracts, but in many examples access control 
implemented at the application level but not at the blockchain layer. 

The most common and basic form of access control is the concept of `ownership`: there's one account (combination
of MSP  and certificate identifiers) that is the owner and can do administrative tasks on contracts. This 
approach is perfectly reasonable for contracts that only have a single administrative user.

CCKit provides `owner` extension for implementing ownership and access control in Hyperledger Fabric chaincodes.