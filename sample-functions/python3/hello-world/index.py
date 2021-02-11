def handler(payload, context):
    return {
        'statusCode': 200,
        'body': 'Hello, world!',
        "headers": {
            "Content-Type": "text/plain",
            "X-LLNW-Faas-Collect-Stdio": "1"
        }
    }