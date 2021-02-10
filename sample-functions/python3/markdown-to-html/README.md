EdgeFunctions Functions Layout
==============================
EdgeFunctions functions must export a handler as the entry point for the application. This function uses the standard naming convention of `index.handler`.

Function Description
====================
**Markdown to HTML Converter**

This function converts one or more [markdown files](https://en.wikipedia.org/wiki/Markdown) to HTML using the [Python-Markdown](https://pypi.org/project/Markdown/) library.

Before creating the ZIP achive for this function, you'll need to install the following two libraries in the same folder as your function:
* `$ pip install --target <path_to_function_folder> markdown`
* `$ pip install --target <path_to_function_folder> requests`

Parameters
=====================
This function requires two parameters, `path` and `doc`, which can be provided in either the Environment Variables or in the query string. The optional parameters `ext`, `format` and `encode` can also be provided by either approach.

* `path` = the fully-qualified path to the converted file(s)
* `doc` = a comma-separated list of input filenames (without extensions); documents are concatenated in the response in the order listed
* `ext` (optional) = the file extension to append to the filename(s) specified in the doc parameter; default = `md`
* `format` (optional) = the content format to return; default = `html`
  * `html` returns the markdown converted to HTML
  * `markdown` returns the original markdown wrapped in HTML "pre" tags
  * `json` returns the results in a JSON object with the form: {"results": [{"doc": "{filename}", "html": "{output}", "markdown": "{input}"}]}
* `encode` = whether to encode response HTML when `format=json`; default = `false`
  * `false` returns unencoded HTML in JSON responses
  * `true` returns encoded HTML in JSON responses; needed when HTML includes the `"` character

Environment Variables
=====================
Environment variables set the default behavior of the function. For example:

    [
        {
            "name":"path",
            "value":"https://support.limelight.com/public/demo/files/ef/markdown"
        },
        {
            "name":"doc",
            "value":"markdown-demo"
        }
    ]

Query Parameters
================
If provided, query parameters overwrite the parameters of the same name as set in the Environment Variables.

Returns
=======

**Success**: The result is returned in the response body.

**Failure**: If a required parameter is not provided, the function returns an error description in the body.

Sample Use
==========
This function must be  invoked via a public URL that has first been mapped to a special origin: the `/epInvoke` endpoint. The endpoint takes the form:

`https://apis.llnw.com/ef-api/v1/{shortname}/functions/{function}/epInvoke`

The mapping is accomplished by creating a new Limelight CDN configuration (a *Service Instance*) associated only with this specific function. The Service Instance must be created and successfully deployed before the function will run.

In addition, the Service Instance must be created using a special *Service Profile* (a template of predefined configuration settings) created specifically for your use with EdgeFunctions. This Service Profile is automatically generated when you successfully create (upload) your first function, and is named using the form `{shortname}-EdgeFunctions`.

For more on creating EdgeFunctions and their associated Service Instances, see [Creating a “Hello World” EdgeFunction via API](https://developers.limelight.com/community?id=community_blog&sys_id=7557090e1b0e6090b93d43b3cd4bcb4f).


