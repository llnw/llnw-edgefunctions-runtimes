# EdgeFunctions Function Layout

An EdgeFunctions function must export a handler as the entry point for the application.  This function uses the standard naming convention of  `index.handler`.

# Function Description

Blacklist Referer redirect function

If the function is invoked with a Referer header, the Referer value will be checked for a matching hostname in the list of hostnames defined in resources/blacklist.js.
If a match is found, the function will respond with a 302 redirect, and if a match is not found it will return with an empty body and a 200 status code.


## Returns

The function returns a 302 redirect on a referer match.

The function returns a 200 status code if a match is not found or a referer header is not provided.

And it returns a 404 on failure.

## Sample Use

Create or update a function with the zip archive blackList.zip.

Then Invoke the function with epInvoke and with the Referer header:

/{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke

* The function will then take the hostname in the Referer header and check for a matching hostname in resources/blacklist.js
