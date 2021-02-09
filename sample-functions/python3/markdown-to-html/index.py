import markdown
import requests

import html
import json
import os

def convert(text):
    # return markdown converted to HTML
    return markdown.markdown(text, extensions=['extra'])

def handler(payload, context):
    queries = payload['queries']
    envar = os.environ
    
    # get path (required) from queries or environment variables
    if "path" in queries.keys():
        path = queries["path"].strip()
    elif "path" in os.environ:
        path = envar["path"].strip()
    else:
        return {
            "statusCode": 500,
            "body": "No path specified in either query string or Environment Variables"
        }
    
    if path == "" or path is None:
        return {
            "statusCode": 500,
            "body": "Path name is empty"
        }
    
    # get document name (required) from queries or environment variables
    if "doc" in queries.keys():
        doc = queries["doc"].strip()
    elif "doc" in envar:
        doc = envar["doc"].strip()
    else:
        return {
            "statusCode": 500,
            "body": "No document specified in either query string or Environment Variables"
        }
    
    if doc == "" or doc is None:
        return {
            "statusCode": 500,
            "body": "Document name is empty"
        }
    
    # get document extension (optional) from queries or environment variables
    ext = "md"
    if "ext" in queries.keys():
        ext = queries["ext"].strip()
    elif "ext" in envar:
        ext = envar["ext"].strip()
    
    if ext == "" or ext is None:
        return {
            "statusCode": 500,
            "body": "Extension name is empty"
        }
    
    # get response format (optional) from queries or environment variables
    format = "html"
    if "format" in queries.keys():
        format = queries["format"].strip()
    elif "format" in envar:
        format = envar["format"].strip()
    
    if format == "" or format is None:
        return {
            "statusCode": 500,
            "body": "Format name is empty"
        }
    elif format not in ["markdown", "html", "json"]:
        return {
            "statusCode": 500,
            "body": "Unexpected format '" + format + "'"
        }
    
    # convert the document name(s) to an array
    docs = doc.split(",")
    result = "" if format != "json" else {"results": []}
    type = "text/html" if format != "json" else "application/json"
    
    for doc in docs:
        # get the markdown text
        url = path + "/" + doc + "." + ext
        response = requests.get(url, verify=True)
        
        if response.ok:
            text = response.text
            
            if format == "json":
                # append the result as an object
                md = convert(text)
                md = html.escape(md, quote=True)
                
                result["results"].append({"doc": doc, "html": md, "markdown": text})
            elif format == "html":
                # append the converted markdown (HTML)
                result += convert(text)
            else:
                # append the original markdown wrapped in pre tags for better readability in browsers
                result += "<pre>" + text + "</pre>"
        else:
            # fail - return the error code and body
            return {
                "statusCode": response.status_code,
                "body": html.escape(response.text, quote=True)
            }
        
    # success - return the result
    if format == "json":
        result = json.dumps(result)
    return {
        "statusCode": 200,
        "body": result,
        "headers": {
            "Content-Type": type
        }
    }

# def test():
    # payload = {
       # "httpMethod": "",
       # "remoteAddress": "",
       # "host": "",
       # "path": "",
       # "headers": {},
       # "multiValueHeaders": {},
       # "queries": {
           # "path": "https://support.limelight.com/public/demo/files/ef/markdown",
           # "doc": "demo-markdown-to-html"
       # },
       # "multiValueQueries": {},
       # "body": "",
       # "isBase64Encoded": False
    # }
    
    # result = handler(payload, '')
    # return (result)

# print (test())
