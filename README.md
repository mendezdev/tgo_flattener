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
- **URL** ```POST /flats```
  - **INFO**: This will accept a JSON array with nested arrays of simple values like `string`, `int`, `float` and even `null`
  - **RESPONSE**: 
    - **404**: if you send an object value inside the array
    - **500**: this is work in progress and the algorithm should be improved
    - **200**: returns an JSON object with the flatted array and max depth of it
      - **BODY EXAMPLE**: 
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
      - **RESPONSE EXAMPLE**:
      ```
      {
        "max_depth": 1,
        "flatted_data": [3,"0_lvl","1_lvl",1,2]
      }
      ```
- **URL** ```GET /flats```
  - **RESPONSE**:
    - **500**: this is work in progress and the algorithm should be improved
    - **200**: returns a JSON array with the last 100 items processed with the ID, the time from when this was processed, the flatted and unflatted array
      - **RESPONSE EXAMPLE**:
      ```
      [
        {
          "id": "60b5a1727c09e9d6a3cefec4",
          "processed_at": "2021-06-01T02:54:42.088Z",
          "unflatted": [
              2,
              3,
              "0_lvl",
              [
                  "1_lvl"
              ],
              1
          ],
          "flatted": ["0_lvl","1_lvl",1,2,3]
        }
      ]
      ``` 
