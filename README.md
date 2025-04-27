
# Calculator

This is a simple calculator which computes expressions with integers. You can see it deployed at [https://calculator.etboye.dk](https://calculator.etboye.dk).

### Features

 - You can compute complicated expressions with integers, the usual four binary operators and brackets. Parsing of the input is done using [participle](https://github.com/alecthomas/participle).
 - You can see historic calculations for sessions. The backend api has a paginated endpoint following the [zalando guidelines](https://opensource.zalando.com/restful-api-guidelines/#pagination).


# Backend API

The backend is written in golang using gin. See the included postman collection.

The backend exposes two endpoints.
## Compute
```
POST /sessions/:sessionid/compute
Exapmple body: 
{
	"input": "1+4"
}
Example response: 
{
    "historyRow": {
        "calculationId": 64,
        "calculation": {
            "input": "1+4",
            "result": {
                "num": "5",
                "denom": "1",
                "estimate": "5.00000"
            },
            "errorId": null
        }
    },
    "error": null
}
```
The innermost `errorId`field describes errors you would expect to see from a calculator that can compute expressions - for example DIVISION_BY_ZERO_ERROR, PARSING_ERROR, LEXING_ERROR. The outermost `error` field is used when the response is not 2xx - for example if the `sessionId` does not match the regex `/^[a-zA-Z0-9\-]*$`

The `result`object uses strings because we compute with big integers which might be too large for javascripts number-object.

## History
Returns a list of at most 5 historic computations under the `sessionId`. We use cursor-based pagination as per the [zalando guidelines](https://opensource.zalando.com/restful-api-guidelines/#pagination).
```
GET /sessions/:sessionId/history?cursor=:cursor
Example GET /sessions/init-db-data-session/history
Example Response:
{
    "self": "/sessions/init-db-data-session/history?cursor=13",
    "first": "/sessions/init-db-data-session/history?cursor=13",
    "prev": null,
    "next": "/sessions/init-db-data-session/history?cursor=8",
    "last": "/sessions/init-db-data-session/history?cursor=3",
    "items": [
        {
            "calculationId": 13,
            "calculation": {
                "input": "12",
                "result": {
                    "num": "12",
                    "denom": "1",
                    "estimate": "12"
                },
                "errorId": null
            }
        },
        {
            "calculationId": 12,
            "calculation": {
                "input": "11",
                "result": {
                    "num": "11",
                    "denom": "1",
                    "estimate": "11"
                },
                "errorId": null
            }
        },
        {
            "calculationId": 11,
            "calculation": {
                "input": "10",
                "result": {
                    "num": "10",
                    "denom": "1",
                    "estimate": "10"
                },
                "errorId": null
            }
        },
        {
            "calculationId": 10,
            "calculation": {
                "input": "9",
                "result": {
                    "num": "9",
                    "denom": "1",
                    "estimate": "9"
                },
                "errorId": null
            }
        },
        {
            "calculationId": 9,
            "calculation": {
                "input": "8",
                "result": {
                    "num": "8",
                    "denom": "1",
                    "estimate": "8"
                },
                "errorId": null
            }
        }
    ],
    "error": null
}
```

## Parameters
### Session id
An id defining the "calculation session". Length has to be at least 1 and at most 100, and it has to match the regex `/^[a-zA-Z0-9\-]*$`.

The frontend generates an uuidv4 and stores it in localStorage upon initialization.

### Cursor
An id defining the id of the most recent historic calculation to be returned in the paginated list of the history endpoint. If the cursor query param is not present, the list of the most recent historic calculations will be returned.

  

# Testing locally

### Prerequisites for running the app

This app uses

- Docker compose for easy deployment of the different services.

- Dozzle for container observability

- For the real deployment on digital ocean, I want dozzle behind basic authentication. For testing, you can run

`cp nginx/.htpasswd.unsafe-admin-admin nginx/.htpasswd` to create a simple user with username 'admin' and password 'admin'.

You can test the whole application by using docker-compose by running

`POSTGRES_USER=user POSTGRES_PASSWORD=secret APP_ENV=test ./run.sh`

  

This exposes

 - The whole application as it will be deployed on `http://localhost:8089`
  - Note that requests `http://localhost:8089/api/:restofurl` is proxied to the backend as `/:restofurl`
 - Dozzle on `http://localhost:8089/dozzle`
 - The backend on `http://localhost:8080`

  
  

# Deployment on digital ocean

This application is running at [https://calculator.etboye.dk](https://calculator.etboye.dk) and deployed using a droplet on digital ocean. A few notes:

- The droplet has nginx installed and uses a certificate from LetsEncrypt. Everything from nginx is proxied to `localhost:8089` where the nginx from the docker-compose app is exposed.

- The droplet is behind a firewall configured directly from digital ocean.

  

# Next up

I would like to
 - Write more unit tests
 - Refactor the pagination code
 - Put an index on the history table?
 - Refactor errors in response bodies using the SimpleHttpResponse generic struct
 - Do input and output validation to handle errors when any value is outside the column sizes defined in the database
 - Find all the TODOS in the code and do them


