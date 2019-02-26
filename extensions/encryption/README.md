# Hyperledger Fabric chaincode kit (CCKit) - Encryption extension

Allows to encrypt all data in ledger. Based on [example](https://github.com/hyperledger/fabric/tree/master/examples/chaincode/go/enccc_example)


## Using ECDH for establishing secret key 

ECDH (Elliptic curve Diffie-Hellman): both parties can establish a secret value by sending only the public key 
of their ephemeral or static key pair to the other party. If the key pair of one of the parties is trusted by the other
party then that key pair may also be used for authentication. 

 
## Flow of a CHAINCODE transaction 

### 1. Sending proposal from client to peer 

The [proposal](https://github.com/hyperledger/fabric/blob/master/protos/peer/proposal.proto)  is basically a request to 
do something on a [chaincode](https://github.com/hyperledger/fabric/blob/master/protos/peer/chaincode.proto), 
that will result on some action - some change in the state of a chaincode and/or some data to be committed to the ledger.   

```
SignedProposal
|\_ Signature                                    (signature on the Proposal message by the creator specified in the header)
 \_ Proposal
    |\_ Header                                  
    |\_ ChaincodeProposalPayload             
    |   |\_ ChaincodeInvocationSpec
    |   |    \_ ChaincodeSpec
    |   |        |\_ Chaincode_id 
    |   |         \_ Input 
    |   |             \_ Args
    |    \_ TransiendMap
     \_ ChaincodeAction                          (the actions for this proposal - optional for a proposal)
```     
    
Proposal contains:
    
     * `Signature` is a signature of the client over the Proposal.
     * `ChannelId` - name of the fabric channel that is the target of this transaction 
     * `ChaincodeName` and `ChaincodeVersion` - the name and version of the chaincode that is being invoked by this proposal.
     * `TxId` -  the transaction identifier, computed as the hash of SignatureHeader.
     * `Creator` - is the serialized identity of the client.
     * `Nonce` - an array of random bytes.
     * `Args` - the set of arguments for this chaincode invocation.
    
     
> For encrypting data can be used contents of `Trasient Map` field - data (e.g. cryptographic material) that might be used to implement 
some form of application-level confidentiality. 

> The contents of this field are supposed to always be omitted from the transaction and
excluded from the ledger.



### 2. Peer sends proposal response back to client

The proposal response contains an endorser's response to a client's proposal. A proposal response contains a success/error code, 
a response payload and a signature (also referred to as endorsement) over the response payload.

The response payload contains a hash of the proposal (to securely link this response to the corresponding proposal) 
and an opaque extension field that depends on the type specified in the header of the corresponding proposal. A
proposal response contains the following messages:

```
ProposalResponse
|\_ Endorsement                                  (the endorser's signature over the whole response payload)
 \_ ProposalResponsePayload                      (the payload of the proposal response)
```

### 3. Client assembles endorsements into a transaction
 
A transaction message assembles one or more proposals and corresponding responses into a message to be sent to orderers. 
After ordering, (batches of) transactions are delivered to committing peers for validation and final delivery into the ledger. 
A transaction contains one or more actions. Each of them contains a header (same as that of the proposal that requested it) 
and an opaque payload that depends on the type specified in the header.

```
SignedTransaction
|\_ Signature                                    (signature on the Transaction message by the creator specified in the header)
 \_ Transaction
     \_ TransactionAction (1...n)
        |\_ Header (1)                           (the header of the proposal that requested this action)
         \_ Payload (1)                          (the payload for this action)
```


