# EdgeFunctions Function Layout

An EdgeFunctions function must export a handler as the entry point for the application.  This function uses the **non-standard** naming convention of `index.proxy`.

# Function Description

Proxy URI download function.

When invoked the function will do an http get on a URI and pass the body on in the invoke response.
To use it, you configure a URI environment variable, which the function sends a get request to.


## Required Environment Variables

URI = uri which the function proxies a request to, e.g:
{
    "Name":"URI",
    "Value": "http://www.limelightnetworks.com/"
}

## Returns

The function returns the response from the proxied request to the URI.

## Sample Use

Create or update a function with an environment variable defining the value of your URI,
e.g:

```
{
        "name": "proxy",
        "description": "do a proxy download via URI environment variable",
        "functionArchive": "...",
        "environmentVariables": [
            {
                "name": "URI",
                "value": "http://www.limelightnetworks.com/"
            }
        ],
        "handler": "index.proxy",
        "memory": 128,
        "runtime": "nodejs12",
        "timeout": 4000
}
```

Invoke the function:

/{SHORT_NAME}/functions/{FUNCTION_NAME}/invoke 

*  Will proxy a request to http://www.limelightnetworks.com/ and return the response
