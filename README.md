# FetchRewards Receipt Processor
### Receipt Processor Interview Assessment
 
* POST /receipts/process : Saves a receipts data internally for future access\
Example request body:
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
* GET /receipts/{id}/points : Provides the numerical points value of the requested receipt


### How to run Application ###
1. Clone repository
2. Either use the included Dockerfile to build and image and run that in a contianer or run the go application directly by executing the following commands:
```
go get github.com/google/uuid
go run main.go
```
3. Now the Receipt Processor server should be running and listening for requests

### POST request using /receipts/process ###
1. Using a tool such as Postman, create a POST request to url http://localhost:8080/receipts/process
2. Fill in the body of the POST request with a JSON object such as the example above contining the receipt contents
3. Execute the POST request and observe the response which will contain the unique id of that receipt now store in the server


### GET request using /receipts/{id}/points ###
1. Using a tool such as Postman, create a GET request to url http://localhost:8080/receipts/{id}/points
2. Replace the {id} in the url with the unique id that was returned from the above POST command which processed the receipt
3. Execute the GET request and observe the reponse which will contain the numerical points value of the requested receipt