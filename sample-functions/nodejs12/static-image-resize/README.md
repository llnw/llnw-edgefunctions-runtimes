# EdgeFunctions Function Layout

An EdgeFunctions function must export a handler as the entry point for the application.  This function uses the standard naming convention of index.handler.

# Function Description

Image resize function.  

Returns and stores resized image in Origin Storage by leveraging Sharp https://sharp.dimens.io/en/stable/.

This module supports reading JPEG, PNG, WebP, TIFF, GIF and SVG images.

Output images can be in JPEG, PNG, WebP and TIFF formats as well as uncompressed raw pixel data.

## Environment Vars

Required: One of the follow authentication methods:

### Token

  * LLNW_USERNAME: string - origin storage username
  * LLNW_PASSWORD: string - origin storage password

### HMAC

  * HMAC_ACCESS: origin storage access key
  * HMAC_SECRET: origin storage secret key

## Request Parameters

required request parameters:

* url  = original image url - url should NOT contain request parameters as they may cause conflicts with function
* path = path in storage where resized image will be placed

optional request parameters:

* join = bool: if set this will join the image url path to the specified storage path url
  * ie: ?url=http://foo.com/bar/baz.png&path=/images will create /images/bar/baz.png in storage account provided
* fmt  = format: string - options: 'png', 'jpeg', 'webp' and 'tiff' default is png
* fit  = fit: string - options: 'cover', 'contain', 'fill', 'inside', 'outside', 'entropy', 'attention' default cover
* pos  = position: string - options: 'centre' 'top' 'right' 'left' 'bottom' 'left,top' etc default 'centre'
* hex  = hex: string - background color: 'ffffff,0.8'


## Returns

The function returns a binary buffer and status code 200 when successful.

The function returns status code 404 on failure.


## Sample Use

Create or update a function with the zip archive staticResize.zip and of the environment variable authentication methods, e.g:

    "environmentVariables": [
        {
            "Name": "LLNW_USERNAME",
            "Value": "username"
        },		
        {
            "Name": "LLNW_PASSWORD",
            "Value": "password"
        }
    ]


Invoke the function with desired request parameters, e.g:

Sample Request: `{SHORT_NAME}/functions/{FUNCTION_NAME}/epInvoke/?url=http://global.mt.lldns.net/llnwstor/faas/images/faas-image-resize-sample.png&pos=bottom`
