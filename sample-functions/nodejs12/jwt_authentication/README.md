# EdgeFunctions Function Layout

An EdgeFunctions function must export a handler as the entry point for the application.  This function uses the standard naming convention of `index.handler`.

# Function Description

JWT Authentication proxy function.

The function takes in an HTTP request forwarded from EdgePrism. Validates the JWT Authorization header then proxies the request and returns the response.

The function does simple HMAC JWT authorization, and you can set your HMAC secret with the environment variables `JWT_SECRET` inside the function.

The function gets the JWT authorization token from the request header named Authorization.

In order to create your own zip archive, you need to add the libraries jsonwebtoken and node-fetch to the package.  This can done using the following commands:
- $ npm install jsonwebtoken
- $ npm install node-fetch

## Request Paramaters

Use standard EP Invoke request. Will proxy to the `http://{host}:{path}` with the specified HTTP `{method}` set in the EP Invoke request and copy all headers and query params.

Optional Parameters:
- xLLnwHost = host name of proxy request
- xLlnwPath = path name of proxy rqeuest
- xLlnwHttpMethod =  request method to be made (e.g put, post, get, delete)

## Environment Variables

Required:
- JWT_SECRET = HMAC secret used to validate the signature in the JWT authorization token


## Returns

If the JWT is valid, the function will proxy the request and return the response.

If the JWT is invalid or fails, the function returns a status code 401.

## Sample Use

Create or update a function with the provided zip archive and with an environment variable named JWT_SECRET, e.g:

    "environmentVariables": [
        {
            "Name": "JWT_SECRET",
            "Value": "test"
        }
    ]

Invoke the function with epInvoke, the optional query parameters, and by passing the JSON Web Token as a header named 'Authorization':

/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke?xLLnwHost=A.com&xLLnwPath=/example&xLlnwHttpMethod=get

* The signature of the Authorization header is validated using the env var JWT_SECRET

* If it is valid, then the function will make a get request to A.com/example and return the response


/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke?xLLnwHost=B.com&xLLnwPath=/mydata&xLlnwHttpMethod=post

* The signature of the Authorization header is validated using the env var JWT_SECRET

* If it is valid, then the function will make a post request to B.com/mydata and return the response