const jwt = require("jsonwebtoken");
const fetch = require('node-fetch');

exports.handler = async (payload, context) => {
  try {
    let isValid = await isValidJwt(payload);
    
    if (!isValid) {
      return {
        "statusCode": 401,
        "body": "JWT Authorization Failed"
      }

    } else {
        return parseResponse(payload)
    }

  } catch(err) {
    console.log("Error caught by function try/catch");
    return {
      "statusCode": 401,
      "body": err.message
    }  
  } 
}

// get jwt token and check if the signature is valid
async function isValidJwt(request) {
  const encodedToken = getJwt(request);
  if (encodedToken === null) {
    return false
  }

  return isValidJwtSignature(encodedToken)
}

// get jwt from request payload
function getJwt(payload) {
  const authHeader = payload.headers.Authorization;
  // check for bearer prefix and remove it
  if (!authHeader || authHeader.substring(0, 6).toUpperCase() !== 'BEARER') {
    return null
  }
  return authHeader.substring(7).trim()
}

// Get secret from env vars and verify jwt signature
async function isValidJwtSignature(token) {
  
  let secret = process.env['JWT_SECRET'];
    
  try {
    var decoded = jwt.verify(token, secret)
    return true
  }
  catch(exception) {
    console.log(exception)
    return false
  }
}

// make and parse http request into ep response
async function parseResponse(payload){
  // construct url to fetch from host and path
  let url = "http://" + payload.host + payload.path;

  // fetch http response from url
  let response = await fetch(url)

  // parse repsonse
  let epBody = await response.json()
  var singleValueHeaders = {}
  var multiValueHeaders = {}

  // parse headers
  for (var [key, val] of response.headers.entries()) {
    // all single value headers will be returned as strings
    if (typeof val == "string") {
      singleValueHeaders[key] = val
    }
    // multi value headers will be returned as lists
    else {
      multiValueHeaders[key] = val
    }
  }

  return {
    "statusCode": response.status,
    "body": epBody,
    "headers": singleValueHeaders,
    "multiValueHeaders": multiValueHeaders
  }
}