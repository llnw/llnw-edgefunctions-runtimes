# EdgeFunctions Function Layout

An EdgeFunctions function must export a handler as the entry point for the application.  This function uses the standard naming convention of index.handler.

# Function Description

A-B testing function.

The function distributes requests between 2+ hosts according to a weight function.
Subsequent requests can be made sticky to original host by specifying an optional tag. 

Host options, including host name and weights, should be stored in the function's environment variables.

The weight of each host can be set using any positive integer or decimal value.  
The probability of a host being selected is represented by (weight of host)/(total weight of all hosts).
For example the weights of 0.5 & 0.5, 1 & 1, or 2 & 2 are equivalent, and the weights 0.25 & 0.75 or 1 & 3 are equivalent.

## Required Environment Variables

LLNW_HOSTS = a json serialized set of hosts optionally weighted to pick from, default weight is even distribution
e.g:
{
    "Name":"LLNW_HOSTS",
    "Value": '[{"name":"llnw.com","weight":0.1},{"name":"limelight.com","weight":0.9}]'
}

## Optional Request Parameters

path = url path appended to the selected host name, no path is default
* With host names A.com and B.com and path = c/example, redirects will be to A.com/c/example and B.com/c/example

envvar = alternate name for the required environment variable, default is 'LLNW_HOSTS'
* If envvar = HOSTS, then there must be a environment variable with name HOSTS, e.g: 
    {
        "Name":"HOSTS",
        "Value": '[{"name":"llnw.com","weight":0.5},{"name":"limelight.com","weight":0.5}]'
    }

tag = string value used as a seed for the randomly generated decimal that determines the host selected
* Each invocation with the same tag will return the same host
  

## Returns

The function returns a 302 redirect to a host when successful and a 404 on failure.


## Sample Use

Create or update a function with environment variables defining your host set
e.g:

    "environmentVariables": [
        {
            "Name": "LLNW_HOSTS",
            "Value": '[{"name":"llnw.com","weight":0.5},{"name":"limelight.com","weight":0.5}]'
        },		
        {
            "Name": "EXAMPLE_HOSTS",
            "Value": '[{"name":"llnw_example.com","weight":1},{"name":"limelight_example.com","weight":3]'
        }
    ]


Invoke the function with epInvoke and optional query paramters (path and envvar):

/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke 

* Will invoke the function with the set of hosts [{"name":"llnw.com","weight":0.5},{"name":"limelight.com","weight":0.5}]

* And will return a redirect on llnw.com or limelight.com with equal probability


/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke?path=test

* Will invoke the function with the set of hosts [{"name":"llnw.com","weight":0.5},{"name":"limelight.com","weight":0.5}]

* And will return a redirect on llnw.com/test or limelight.com/test with equal probability


/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke?path=test&envvar=EXAMPLE_HOSTS

* Will invoke the function with the set of hosts [{"name":"llnw_example.com","weight":1},{"name":"limelight_example.com","weight":3]

* And will return a redirect llnw_example.com/test with 25% probability or limelight_example.com/test with 75% probability

/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke?tag=2

* Will invoke the function with the set of hosts [{"name":"llnw.com","weight":0.5},{"name":"limelight.com","weight":0.5}]

* And will return a redirect to the same host on every invoke where tag=2