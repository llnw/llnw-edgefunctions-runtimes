import json
import os
import requests
import jwt

# env var name of jwt secret
JWTENV = "JWT_SECRET"

def handler(request, context):
    try:
        secret = loadSecret()

        tokenString = getJWTToken(request)

        isValid(tokenString, secret)

        return proxy(request)

    except Exception as e:
        body = str(e)
        if hasattr(e, 'message'):
            body = e.message
        return {
            'statusCode': 401,
            'headers': {},
            'body': json.dumps(str(body))
        }
        
# get secret from function environment variables
def loadSecret():
    if JWTENV in os.environ:
        value = os.environ[JWTENV]
    else:
        raise Exception("environment variable secret not found")

    return value

# get JWT Token from the request queries
def getJWTToken(request):
    if 'Authorization' not in request['headers']:
        raise Exception("authorization token not in headers")    
    value = request['headers']['Authorization']

    # check for and remove bearer prefix 
    if "BEARER" not in value.upper():
        raise Exception("BEARER not in authorization header")
    tokenString = value[7:]
    return tokenString

# decode token string using secret from env vars, checking token signature for validity
def isValid(tokenString, secret):
    decodedToken = jwt.decode(tokenString, secret, algorithms=['HS256'])

# proxy request and return response
def proxy(request):
    response = makeRequest(request)

    return parseResponse(response)

# parse ep request and perform http request 
def makeRequest(request):
    body = None
    if request["body"] != "":
        body = request["body"]

    url = "http://" + request["host"] + request["path"]

    # Merge headers and multivalue headers
    headers = {**request["headers"], **request["multiValueHeaders"]}

    # Merge query and multivalue queries
    params = {**request["queries"], **request["multiValueQueries"]}

    # Make request, request method based on httpMethod in the ep request
    response = requests.request(request["httpMethod"], url, headers=headers, params=params)

    return response


# parse http response into ep response
def parseResponse(response):
    epResponse = {
        "statusCode": response.status_code,
        "body": response.text
    }

    headers = {}
    multiValueHeaders = {}
    
    # copy and seperate multi-value headers and single value headers
    for key, val in response.headers.items():
        if isinstance(val, list):
            multiValueHeaders[key] = val
        else:
            headers[key] = val
            
    # add headers to ep response
    epResponse["headers"] = headers
    epResponse["multiValueHeaders"] = multiValueHeaders
    
    return epResponse

