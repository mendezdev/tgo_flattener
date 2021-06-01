# tgo_flattener

This is a simple server that receive a multilevel arrays and returns a flat array and the max depth of it. 


## Requirements
- Clone this repo
- This beta version use [MongoDB 4.2.3+](https://docs.mongodb.com/manual/administration/install-community/) for running local and for tests. You need to install it a run on default port :27017.
- Install [go 1.16.3+](https://golang.org/doc/install).
- **IMPORTANT**: if you don't have ```go mod``` enabled, see this [article](https://lets-go.alexedwards.net/sample/02.02-project-setup-and-enabling-modules.html)

## How to run the tests
- Put MongoDB to run on port :27017
- Open the terminal, go to the root folder of this app and execute ```go test ./...```

## How to run the app
- Put MongoDB to run on port :27017
- Open the terminal, go to the root folder of this app and execute ```go run main.go```. This will run on port ```:8080```

## ENDPOINTS
- URL ```POST /flats```
  - BODY: This will accept an array with nested arrays of simple values like `string`, `int`, `float` and even `null`
  - RESPONSE: 
    - 404: if you send an object value inside the array
    - 500: this is work in progress and the algorithm should be improved
    - 200: if all is OK you will get the flatted array and max depth of it
      - body example: 
      ```
      [
        "0_lvl",
        [
          "1_lvl"
        ],
        1,
        2,
        3
      ]
      ```
      - response example:
      ```
      {
        "max_depth": 1,
        "flatted_data": [3,"0_lvl","1_lvl",1,2]
      }
      ```
