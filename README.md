# Serverless-API

# This is the ReadMe file for the serverless rest API challenge:
There is a package called **functions** which contains two folders, **get** and **post**. 

Each of them has its main function, test functions, and cover files. To get detailed information on the test coverage, you can check out the **cover.html**.

---
## Deploy:
Install Serverless Framework and configure your AWS credentials on your machine. Now clone this repo into your Go Path Then cd into the project folder.
run command below to construct the services.

```bash
serverless deploy
```
----
## Run:
Once deployed and substituting your API-URL, you can use Postman to interact with the resulting API, whose results can be confirmed in the DynamoDB console.You can have two kinds of requests, post and get. 

### POST:
Use your API <API URL> in the Postman endpoint bar and Select POST from the list of request types. click on the Body tab and click on raw and select format type as JSON. Now you can insert you request body and press Send.

for example:\
 {\
  "id": "id22",\
  "deviceModel": "d.m451",\
  "name": "Sensor",\
  "note": "Testing a sensor.",\
  "serial": "D0010245"\
 }\
and you can see the result.

### GET
Use your API <API URL> in the Postman endpoint bar and select GET from the list of request types. Click on the Params tab, enter key (id) and value (id22) and press Send. You can see all information that you posted about the device with id22.

Note: If you are going to see the result of my implemention, there is my API URL: https://wjqtudixzd.execute-api.us-east-1.amazonaws.com/api/devices

---
## Unit Testing:
You Can run tests by using **go test** command in the get or the post directories.

```bash
go test -v --cover
```
Note: For unit testing please comment line 82 and uncomment lines 85 to 87 in the get function.

