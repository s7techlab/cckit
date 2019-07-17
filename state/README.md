# Working with Hyperledger Fabric chaincode state with CCKit

Chaincode is a domain specific program which relates to specific business process. It programmatically accesses 
two distinct pieces of the ledger – a blockchain, which immutably records the history of all transactions, and a world state
that holds a cache of the current value of these states. The job of a smart contract developer is to take an existing business 
process that might govern financial prices or delivery conditions, and express it as a smart contract in a programming language

Smart contracts primarily put, get and delete states in the world state, and can also query the state change history.
Chaincode “shim” APIs implements [ChaincodeStubInterface](https://godoc.org/github.com/hyperledger/fabric/core/chaincode/shim#ChaincodeStubInterface) 
which contain methods for access and modify the ledger, and to make invocations between chaincodes. Main methods are:

* `GetState(key string) ([]byte, error)` performs a query to retrieve information about the current state of a business object

* `PutState(key string, value []byte) error` creates a new business object or modifies an existing one in the ledger world state

* `DelState(key string) error` removes of a business object from the current state of the ledger, but not its history

* `GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)` 
   queries the state in the ledger based on given partial composite key
   
* `GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)` returns a history of key values across time.


All this methods use string key as record identifier and slice of bytes as state value. Most of examples uses JSON documents as 
chaincode state value. Hyperledger Fabric supports both LevelDB as CouchDB to serve as state database, holding the latest state of each object.
LevelDB is the default key-value state database embedded in every peer. CouchDB is an optional alternative external 
state database with more features - it supports rich queries against JSON documents in chaincode state, whereas LevelDB only supports
 queries against keys.
 

## Querying and updating state with ChaincodeStubInterface methods

As shown in many examples, assets can be represented as complex structures - Golang structs. 
The chaincode itself can store data as a string in a key/value pair setup. Thus, we need to marshal struct to JSON string 
before putting into chaincode state and unmarshal after getting from state.

With `ChaincodeStubInterface` methods these operations looks like 
[this](https://github.com/IBM/build-blockchain-insurance-app/blob/master/web/chaincode/src/bcins/invoke_insurance.go):

```go
    ct := ContractType{}
    
    err := json.Unmarshal([]byte(args[0]), &req)
    if err != nil {
        return shim.Error(err.Error())
    }
    
    key, err := stub.CreateCompositeKey(prefixContractType, []string{req.UUID})
    if err != nil {
        return shim.Error(err.Error())
    }
    
    valAsBytes, err := stub.GetState(key)
    if err != nil {
        return shim.Error(err.Error())
    }
    if len(valAsBytes) == 0 {
        return shim.Error("Contract Type could not be found")
    }
    err = json.Unmarshal(valAsBytes, &ct)
    if err != nil {
        return shim.Error(err.Error())
    }
    
    ct.Active = req.Active
    
    valAsBytes, err = json.Marshal(ct)
    if err != nil {
        return shim.Error(err.Error())
    }
    
    err = stub.PutState(key, valAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }
    
    return shim.Success(nil)
```

In the example above smart contract code explicitly performs many auxiliary actions:

* Creating composite key 
* Unmarshaling data after receiving it from state
* Marshaling data before placing it to state

## Modelling chaincode state with CCKit

### State methods wrapper

CCKit contains [wrapper](state.go) on `ChaincodeStubInterface` methods to working with chaincode state. This methods
simplifies chaincode key creation and data transformation during working with chaincode state.

```go
type State interface {
    // Get returns value from state, converted to target type
    // entry can be Key (string or []string) or type implementing Keyer interface
    Get(entry interface{}, target ...interface{}) (result interface{}, err error)
    
    // Get returns value from state, converted to int
    // entry can be Key (string or []string) or type implementing Keyer interface
    GetInt(entry interface{}, defaultValue int) (result int, err error)
    
    // GetHistory returns slice of history records for entry, with values converted to target type
    // entry can be Key (string or []string) or type implementing Keyer interface
    GetHistory(entry interface{}, target interface{}) (result HistoryEntryList, err error)
    
    // Exists returns entry existence in state 
    // entry can be Key (string or []string) or type implementing Keyer interface
    Exists(entry interface{}) (exists bool, err error)
    
    // Put returns result of putting entry to state
    // entry can be Key (string or []string) or type implementing Keyer interface
    // if entry is implements Keyer interface and it's struct or type implementing
    // ToByter interface value can be omitted
    Put(entry interface{}, value ...interface{}) (err error)
    
    // Insert returns result of inserting entry to state
    // If same key exists in state error wil be returned
    // entry can be Key (string or []string) or type implementing Keyer interface
    // if entry is implements Keyer interface and it's struct or type implementing
    // ToByter interface value can be omitted
    Insert(entry interface{}, value ...interface{}) (err error)
    
    // List returns slice of target type
    // namespace can be part of key (string or []string) or entity with defined mapping
    List(namespace interface{}, target ...interface{}) (result []interface{}, err error)
    
    // Delete returns result of deleting entry from state
    // entry can be Key (string or []string) or type implementing Keyer interface
    Delete(entry interface{}) (err error)

	...
}
``` 

### Converting from/to bytes while operating with chaincode state

State wrapper allows to automatically marshal golang type to/from slice of bytes. This type can be:

* Any type implementing [ToByter and FromByter](../convert/convert.go) interface

```go
type (
	// FromByter interface supports FromBytes func for converting from slice of bytes to target type
	FromByter interface {
		FromBytes([]byte) (interface{}, error)
	}

	// ToByter interface supports ToBytes func for converting to slice of bytes from source type
	ToByter interface {
		ToBytes() ([]byte, error)
	}
)
```

* Golang struct or one of supported types ( `int`, `string`, `[]string`)
* [Protobuf](https://developers.google.com/protocol-buffers/docs/gotutorial) golang struct


Golang structs [automatically](../convert) marshals/ unmarshals using [json.Marshal](https://golang.org/pkg/encoding/json/#Marshal) and
and [json.Umarshal](https://golang.org/pkg/encoding/json/#Unmarshal) methods. 
[proto.Marshal](https://godoc.org/github.com/golang/protobuf/proto#Marshal) and 
[proto.Unmarshal](https://godoc.org/github.com/golang/protobuf/proto#Unmarshal) is used to convert protobuf.

### Creating state keys

In the chaincode data model we often need to store many instances of one type on the ledger, such as multiple commercial papers,
letters of credit, and so on.  In this case, the keys of those instances will be typically constructed from a combination of attributes—
for example:

> `CommercialPaper` + {Issuer} + {PaperId}

yielding series of chaincode state entries keys [ `CommercialPaperIssuer1Id1`, `CommercialPaperIssuer2Id2`, ...]


The logic of creation primary key of an instance can be customized in the code, or API functions can be provided in SHIM to
construct a composite key (in other words, a unique key) of an instance based on
a combination of several attributes. Composite keys can then be used as a normal string key to
record and retrieve values using the PutState() and GetState() functions.

The following snippet shows a list of functions that create and work with composite keys:

```go
// The function creates a key by combining the attributes into a single string.
// The arguments must be valid utf8 strings and must not contain U+0000 (nil byte) and U+10FFFF charactres.
func CreateCompositeKey(objectType string, attributes []string) (string, error)

// The function splits the compositeKey into attributes from which the key was formed.
// This function is useful for extracting attributes from keys returned by range queries.
func SplitCompositeKey(compositeKey string) (string, []string, error)
````

When putting or getting data to/from chaincode state you must provide key. CCKit have 3 options for dealing with entries key:
 
* Key can be passed explicit:

```go
c.State().Put ( `my-key`, &myStructInstance)
```

* Key type can implement [Keyer](state.go) interface 

```go
    Key []string

    // Keyer interface for entity containing logic of its key creation
    Keyer interface {
        Key() (Key, error)
    }
````

`Key` type - is essentially slice of string, this slice will be automatically converted to string using `shim.CreateCompositeKey` method.

and in chaincode you need to provide only type instance
```go
c.State().Put (&myStructInstance)
```

* Type can have associate mapping

Mapping defines rules for namespace, primary and other key creation. Mapping mainly used with `protobuf` state schema.

### Range queries

As well as retrieving assets with a unique key, SHIM offers API functions the opportunity to retrieve sets of assets based on a range criteria. 
Moreover, composite keys can be modeled to enable queries against multiple components of the key.

The range functions return an iterator (StateQueryIteratorInterface) over a set of keys matching the query criteria. The returned keys are in lexical order. 
Additionally, when a composite key has multiple attributes, the range query function, `GetStateByPartialCompositeKey()`, can be used to search for keys matching a
subset of the attributes.

For example, the key of a `CommercialPaper` composed of `Issuer` and `PaperId` attributes can be searched for entries only from one Issuer.
  
## Protobuf state example

This example uses [Commercial paper scenario](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/scenario.html) and
implements same functionality as [Node.JS](https://github.com/hyperledger/fabric-samples/tree/release-1.4/commercial-paper/organization/digibank/contract) 
chaincode sample from official documentation. Example code located [here](../examples/cpaper).


Protobuf schema advantages:

1. Schema abstraction layer

Encoding the semantics of your business objects once,  in proto format, is enough to help ensure that the signal doesn’t get lost between applications, 
and that the boundaries you create enforce your business rules.

2. Extensions - validators etc

Protobuf v3 does not support validating required parameters, but there are third party projects for proto validation, for example
https://github.com/mwitkow/go-proto-validators. It allows to encode, at the schema level, the shape of your data structure, and the validation rules.

3. Easy Language Interoperability
Because Protocol Buffers are implemented in a variety of languages, they make interoperability between polyglot applications in your architecture that much simpler. 
If you’re introducing a new service using Java or Node.Js SDK  you simply have to hand the proto file to the code generator written in the target language and you
have guarantees about the safety and interoperability between those architectures. 



### Defining model

Protobuf (short for Protocol buffers) is a way of encoding structured data in an efficient and extensible format. 
With protocol buffers, you write a .proto description of the data structure you wish to store. 
From that, the protocol buffer compiler creates a golang struct  (or ant) that implements automatic encoding and parsing 
of the protocol buffer data with an efficient binary format. The generated class provides getters and setters 
for the fields that make up a protocol buffer and takes care of the details of reading and writing the protocol
 buffer as a unit.

In `Commercial Paper` example first we define messages, that will be stored in chaincode state or as events:

* `CommercialPaper` will be stored in chaincode state
* `CommercialPaperId` defines unique id part of commercial paper message
* `IssueCommercialPaper` payload for `issue` transaction and event triggered when new commercial paper issued
* `BuyCommercialPaper` payload for `buy` transaction and event triggered when commercial paper change owner
* `RedeemCommercialPaper` payload for `redeem` transaction and event triggered when commercial paper redeemed

```proto
syntax = "proto3";
package schema;

import "google/protobuf/timestamp.proto";

message CommercialPaper {

    enum State {
        ISSUED = 0;
        TRADING = 1;
        REDEEMED = 2;
    }

    string issuer = 1;
    string paper_number = 2;
    string owner = 3;
    google.protobuf.Timestamp issue_date = 4;
    google.protobuf.Timestamp maturity_date = 5;
    int32 face_value = 6;
    State state = 7;
}

// CommercialPaperId identifier part
message CommercialPaperId {
    string issuer = 1;
    string paper_number = 2;
}

// IssueCommercialPaper event
message IssueCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    google.protobuf.Timestamp issue_date = 3;
    google.protobuf.Timestamp maturity_date = 4;
    int32 face_value = 5;
}

// BuyCommercialPaper event
message BuyCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    string current_owner = 3;
    string new_owner = 4;
    int32 price = 5;
    google.protobuf.Timestamp purchase_date = 6;
}

// RedeemCommercialPaper event
message RedeemCommercialPaper {
    string issuer = 1;
    string paper_number = 2;
    string redeeming_owner = 3;
    google.protobuf.Timestamp redeem_date = 4;
}
```

### Defining protobuf to chaincode state mapping 

Protocol buffers to chaincode  mapper can be used to store schema instances in chaincode state. Every schema type (protobuf or struct) can have mapping rules:

* Primary key creation logic
* Namespace logic
* Secondary key creation logic

For example, in definition below, we defined thant `schema.CommercialPaper` mapped to chaincode state with key attributes
from `schema.CommercialPaperId` message (`Issuer`, `PaperNumber`). Also we define event

```go
var (
    // State mappings
   	StateMappings = m.StateMappings{}.
   		//key namespace will be <`CommercialPaper`, Issuer, PaperNumber>
   		Add(&schema.CommercialPaper{}, m.PKeySchema(&schema.CommercialPaperId{}))
   
   	// EventMappings
   	EventMappings = m.EventMappings{}.
   		// event name will be `IssueCommercialPaper`,  payload - same as issue payload
   		Add(&schema.IssueCommercialPaper{}).
   		Add(&schema.BuyCommercialPaper{}).
   		Add(&schema.RedeemCommercialPaper{})
)
```

### Chaincode

```go
package cpaper

import (
	"fmt"
	
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	m "github.com/s7techlab/cckit/state/mapping"
)


func NewCC() *router.Chaincode {

	r := router.New(`commercial_paper`)

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	// store in chaincode state information about chaincode first instantiator
	r.Init(owner.InvokeSetFromCreator)

	// method for debug chaincode state
	debug.AddHandlers(r, `debug`, owner.Only)

	r.
		// read methods
		Query(`list`, cpaperList).

		// Get method has 2 params - commercial paper primary key components
		Query(`get`, cpaperGet, defparam.Proto(&schema.CommercialPaperId{})).

		// txn methods
		Invoke(`issue`, cpaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
		Invoke(`buy`, cpaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
		Invoke(`redeem`, cpaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
		Invoke(`delete`, cpaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

	return router.NewChaincode(r)
}


func cpaperList(c router.Context) (interface{}, error) {
	// commercial paper key is composite key <`CommercialPaper`>, {Issuer}, {PaperNumber} >
	// where `CommercialPaper` - namespace of this type
	// list method retrieves entries from chaincode state 
	// using GetStateByPartialCompositeKey method, then unmarshal received from state bytes via proto.Ummarshal method
	// and creates slice of *schema.CommercialPaper
	return c.State().List(&schema.CommercialPaper{})
}

func cpaperIssue(c router.Context) (interface{}, error) {
	var (
		issue  = c.Param().(*schema.IssueCommercialPaper) //default parameter
		cpaper = &schema.CommercialPaper{
			Issuer:       issue.Issuer,
			PaperNumber:  issue.PaperNumber,
			Owner:        issue.Issuer,
			IssueDate:    issue.IssueDate,
			MaturityDate: issue.MaturityDate,
			FaceValue:    issue.FaceValue,
			State:        schema.CommercialPaper_ISSUED, // initial state
		}
		err error
	)

	if err = c.Event().Set(issue); err != nil {
		return nil, err
	}

	return cpaper, c.State().Insert(cpaper)
}

func cpaperBuy(c router.Context) (interface{}, error) {

	var (
		cpaper *schema.CommercialPaper

		// but tx payload
		buy = c.Param().(*schema.BuyCommercialPaper)

		// current commercial paper state
		cp, err = c.State().Get(&schema.CommercialPaper{
			Issuer:      buy.Issuer,
			PaperNumber: buy.PaperNumber}, &schema.CommercialPaper{})
	)

	if err != nil {
		return nil, errors.Wrap(err, `not found`)
	}
	cpaper = cp.(*schema.CommercialPaper)

	// Validate current owner
	if cpaper.Owner != buy.CurrentOwner {
		return nil, fmt.Errorf(`paper %s %s is not owned by %s`, cpaper.Issuer, cpaper.PaperNumber, buy.CurrentOwner)
	}

	// First buy moves state from ISSUED to TRADING
	if cpaper.State == schema.CommercialPaper_ISSUED {
		cpaper.State = schema.CommercialPaper_TRADING
	}

	// Check paper is not already REDEEMED
	if cpaper.State == schema.CommercialPaper_TRADING {
		cpaper.Owner = buy.NewOwner
	} else {
		return nil, fmt.Errorf(`paper %s %s is not trading.current state = %s`, cpaper.Issuer, cpaper.PaperNumber, cpaper.State)
	}

    if err = c.Event().Set(buy); err != nil {
		return nil, err
	}
	return cpaper, c.State().Put(cpaper)
}

func cpaperRedeem(c router.Context) (interface{}, error) {
	
	// implement me

	return nil, nil
}

func cpaperGet(c router.Context) (interface{}, error) {
	return c.State().Get(c.Param().(*schema.CommercialPaperId))
}

func cpaperDelete(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(c.Param().(*schema.CommercialPaperId))
}
```

### Tests
We can [test](mapping/mapping_test.go) all chaincode use case scenarios using [MockStub](../testing)

```go
package mapping_test

import (
	"testing"

	"github.com/hyperledger/fabric/protos/peer"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/examples/cpaper/testdata"
	"github.com/s7techlab/cckit/state"

	"github.com/s7techlab/cckit/examples/cpaper"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	actors   identity.Actors
	cPaperCC *testcc.MockStub
	err      error
)
var _ = Describe(`Mapping`, func() {

	BeforeSuite(func() {
		actors, err = identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
			`owner`: `s7techlab.pem`,
		}, examplecert.Content)

		Expect(err).To(BeNil())

		//Create commercial papers chaincode mock - protobuf based schema
		cPaperCC = testcc.NewMockStub(`cpapers`, cpaper.NewCC())
		cPaperCC.From(actors[`owner`]).Init()

	})

	Describe(`Protobuf based schema`, func() {
		It("Allow to add data to chaincode state", func(done Done) {

			events := cPaperCC.EventSubscription()
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[0]))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `IssueCommercialPaper`,
				Payload:   testcc.MustProtoMarshal(&testdata.CPapers[0]),
			}))

			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[1]))
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[2]))

			close(done)
		}, 0.2)

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(cPaperCC.Invoke(`issue`, &testdata.CPapers[0]))
		})

		It("Allow to get entry list", func() {
			cpapers := expectcc.PayloadIs(cPaperCC.Query(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)
			Expect(len(cpapers)).To(Equal(3))
			Expect(cpapers[0].Issuer).To(Equal(testdata.CPapers[0].Issuer))
			Expect(cpapers[0].PaperNumber).To(Equal(testdata.CPapers[0].PaperNumber))
		})

		It("Allow to get entry raw protobuf", func() {
			cp := testdata.CPapers[0]
			cpaperProtoFromCC := cPaperCC.Query(`get`, &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}).Payload

			stateCpaper := &schema.CommercialPaper{
				Issuer:       cp.Issuer,
				PaperNumber:  cp.PaperNumber,
				Owner:        cp.Issuer,
				IssueDate:    cp.IssueDate,
				MaturityDate: cp.MaturityDate,
				FaceValue:    cp.FaceValue,
				State:        schema.CommercialPaper_ISSUED, // initial state
			}
			cPaperProto, _ := proto.Marshal(stateCpaper)
			Expect(cpaperProtoFromCC).To(Equal(cPaperProto))
		})

		It("Allow update data in chaincode state", func() {
			cp := testdata.CPapers[0]
			expectcc.ResponseOk(cPaperCC.Invoke(`buy`, &schema.BuyCommercialPaper{
				Issuer:       cp.Issuer,
				PaperNumber:  cp.PaperNumber,
				CurrentOwner: cp.Issuer,
				NewOwner:     `some-new-owner`,
				Price:        cp.FaceValue - 10,
				PurchaseDate: ptypes.TimestampNow(),
			}))

			cpaperFromCC := expectcc.PayloadIs(
				cPaperCC.Query(`get`, &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}),
				&schema.CommercialPaper{}).(*schema.CommercialPaper)

			// state is updated
			Expect(cpaperFromCC.State).To(Equal(schema.CommercialPaper_TRADING))
			Expect(cpaperFromCC.Owner).To(Equal(`some-new-owner`))
		})

		It("Allow to delete entry", func() {

			cp := testdata.CPapers[0]
			toDelete := &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}

			expectcc.ResponseOk(cPaperCC.Invoke(`delete`, toDelete))
			cpapers := expectcc.PayloadIs(cPaperCC.Invoke(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)

			Expect(len(cpapers)).To(Equal(2))
			expectcc.ResponseError(cPaperCC.Invoke(`get`, toDelete), state.ErrKeyNotFound)
		})
	})

})
```