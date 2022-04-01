# Fabcar S7Techlab Chaincode

Fabcar S7Techlab is modification of Hyperledger fabric-samples [fabcar chaincode](https://github.com/hyperledger/fabric-samples/blob/main/chaincode/fabcar/go/fabcar.go)

Fabcar Hyperledger Fabric Chaincode (FHF) short description:
1) FHF is made without code generating
2) At FHF you can create cat at once with method 'CreateCar'. Payload example:
```json
{
  "car_number": "CAR1",
  "make": "Toyota",
  "model": "Prius",
  "colour": "blue",
  "owner": "Tomoko"
}
```

2) Then you can get car with method 'QueryCar'. Payload example:
```json
{
  "car_number": "CAR1"
}
```

3) Or get all cars with method 'QueryAllCars' without payload

4) Last method is 'ChangeCarOwner' to change car owner. Payload example:
```json
{
  "car_number": "CAR1",
  "new_owner": "Brad"
}
```



Fabcar S7Techlab Chaincode (FS7) has four entities: Maker, Car, Owner and Detail.
FS7 has some difference from FHF:
1) FS7 gateway was generated with proto 
2) You can not create car at once, because before car's maker have to be created and put at BC state.
   Use method 'CreateMaker', payload example:
```json
{
  "name": "Toyota",
  "country": "Japan",
  "foundation_year,omitempty": "1937" // it must be more than 1886, because this year was founded the oldest automaker - Mercedes-Benz
}
```

2) You can get (method 'GetMaker') or delete (method 'DeleteMaker') maker by its name. For example:
```json
{
  "name": "Toyota"
}
```

3) And get all cars with method 'ListMakers' without payload

4) Now your car can be created with 'CreateCar' method, for example:
```json
{
  "make": "Toyota", // if maker is not created programm will return error
  "model": "Prius",
  "colour": "blue",
  "number": 111111,
  "owners": [
    {
      "first_name": "Tomoko",
      "second_name": "Uemura",
      "vehicle_passport": "bbb222"
    }
  ],
  "details": [
    {
      "type": WHEELS,
      "make": "Michelin"
    },
    {
      "type": BATTERY,
      "make": "BYD"
    }
  ]
}
```

The response is:
```json
{
  "car": {
    "id": ["Toyota", "Prius", "111111"],
    "make": "",
    "model": "Prius",
    "colour": "blue",
    "number": 111111,
    "owners_quantity": 1
  },
  "owners": {
    "items": [
      {
        "car_id": ["Toyota", "Prius", "111111"],
        "first_name": "Tomoko",
        "second_name": "Uemura",
        "vehicle_passport": "bbb222"
      }
    ]
  },
  "details": {
    "items": [
      {
        "car_id": ["Toyota", "Prius", "111111"],
        "type": WHEELS,
        "make": "Michelin"
      },
      {
        "car_id": ["Toyota", "Prius", "111111"],
        "type": BATTERY,
        "make": "BYD"
      }
    ]
  }
}
```

5) Car updating makes with 'UpdateCar', payload:
```json
{
  "id": ["Toyota", "Prius", "111111"],
  "color": "red",
  "owners": [
    {
      "car_id": ["Toyota", "Prius", "111111"],
      "first_name": "Tomoko",
      "second_name": "Uemura",
      "vehicle_passport": "ccc333"
    },
    {
      "car_id": ["Toyota", "Prius", "111111"],
      "first_name": "Michel",
      "second_name": "Uemura",
      "vehicle_passport": "ddd444"
    }
  ],
  "details": [
    {
      "car_id": ["Toyota", "Prius", "111111"],
      "type": BATTERY,
      "make": "Panasonic"
    }
  ]
}
```

The response is:
```json
{
  "car": {
    "id": ["Toyota", "Prius", "111111"],
    "make": "",
    "model": "Prius",
    "colour": "red", // it was 'blue'
    "number": 111111,
    "owners_quantity": 2 // become more
  },
  "owners": {
    "items": [ // become more
      {
        "car_id": ["Toyota", "Prius", "111111"],
        "first_name": "Tomoko",
        "second_name": "Uemura",
        "vehicle_passport": "ccc333" // it was 'bbb222'
      },
      { // it was added
        "car_id": ["Toyota", "Prius", "111111"],
        "first_name": "Michel",
        "second_name": "Tailor",
        "vehicle_passport": "ddd444"
      }
    ]
  },
  "details": {
    "items": [
      {
        "car_id": ["Toyota", "Prius", "111111"],
        "type": WHEELS,
        "make": "Michelin"
      },
      {
        "car_id": ["Toyota", "Prius", "111111"],
        "type": BATTERY,
        "make": "Panasonic" // it was 'BYD'
      }
    ]
  }
}
```

6) Also, you can delete car with 'DeleteCar', the response is like from 'UpdateCar' method (point 5)

7) If you would like to get car, use 'GetCar' to get it without owners and details or 'GetCarView' with them. Request:
```json
{
  "id": ["Toyota", "Prius", "111111"]
}
```

8) To get all cars use 'ListCars' without payload

9) Also, car owner can be updated without car changing ('UpdateCarOwners' method):
```json
{
  "car_id": ["Toyota", "Prius", "111111"],
  "owners": [
    {
      "first_name": "Tomoko",
      "second_name": "Uemura",
      "vehicle_passport": "eee555"
    },
    {
      "first_name": "Adriana",
      "second_name": "Grande",
      "vehicle_passport": "fff666"
    }
  ]
}
```

The response:
```json
{
  "items": [
    {
      "car_id": ["Toyota", "Prius", "111111"],
      "first_name": "Tomoko",
      "second_name": "Uemura",
      "vehicle_passport": "eee555" // it was 'ccc333'
    },
    { // without changes
      "car_id": ["Toyota", "Prius", "111111"],
      "first_name": "Michel",
      "second_name": "Uemura",
      "vehicle_passport": "ddd444"
    },
    { // was added
      "car_id": ["Toyota", "Prius", "111111"],
      "first_name": "Adriana",
      "second_name": "Grande",
      "vehicle_passport": "fff666"
    }
  ]
}
```

10) To delete ('DeleteCarOwner') or get ('GetCarOwner') car owner use the same payload:
```json
{
  "car_id": ["Toyota", "Prius", "111111"],
  "first_name": "Tomoko",
  "second_name": "Uemura"
}
```

11) Also, you can update car detail ('UpdateCarDetails') without car changes, for example:
```json
{
  "car_id": ["Toyota", "Prius", "111111"],
  "details": [
     {
        "type": WHEELS, 
        "make": "Michelin"
     },
     {
        "type": BATTERY,
        "make": "Contemporary Amperex Technology"
     }
  ]
}
```

The response:
```json
{
   "items": [
      { // without changes
         "car_id": ["Toyota", "Prius", "111111"],
         "type": WHEELS,
         "make": "Michelin"
      },
      {
         "car_id": ["Toyota", "Prius", "111111"],
         "type": BATTERY,
         "make": "Contemporary Amperex Technology" // it was 'Panasonic'
      }
   ]
}
```

12) To delete ('DeleteCarDetail') or get ('GetCarDetail') car owner use the same payload:
```json
{
  "car_id": ["Toyota", "Prius", "111111"],
  "type": BATTERY
}
``` 

13) And, of course, you can get list of car owners (method 'ListCarOwners') or list of car details (method 'ListCarDetails')
    by car id. Use payload:
```json
{
   "id": ["Toyota", "Prius", "111111"] 
}
```
