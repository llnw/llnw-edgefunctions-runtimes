import json
import time
import os
from functools import reduce
import hmac
from hashlib import md5

def handler(request, context):
    try:
        # Generate random seed
        seed = int(round(time.time() * 1000))
        MAX_INT = int("0xffffffff", 16)

        # Grab path from queries
        path = ''
        if 'path' in request['queries']:
            path = request['queries']['path']

        # Grab env_var name from queries or set to LLNW_HOSTS
        env_var = 'LLNW_HOSTS'
        if 'envvar' in request['queries']:
            env_var = request['queries']['envvar']

        # Grab tag from queries or set to seed
        tag = str((48271 * seed) % 2147483647)
        if 'tag' in request['queries']:
            tag = request['queries']['tag']

        # Check environment vars present
        if env_var not in os.environ:
            raise Exception("No environment vars found")

        hosts = json.loads(os.environ[env_var])

        if not isinstance(hosts, list):
            raise Exception("No environment vars found")

        # Parse through hosts weights and distribute them
        hosts = distribute_weighting(hosts)

        # Create random percentage based on host
        h = hmac.new(tag.encode('utf-8'), env_var.encode('utf-8'), md5)
        md5_hash = h.hexdigest()[:8]
        hash_int = int("0x" + md5_hash, 16)
        random = hash_int / MAX_INT

        filtered = list(filter(lambda host: host['weight'] > random, hosts))

        redirect = filtered[0]['name'] + "/" + path

        return {
            'statusCode': 302,
            'headers': {
                'Location': redirect
            },
            'body': ''
        }
    except Exception as e:
        body = str(e)
        if hasattr(e, 'message'):
            body = e.message
        return {
            'statusCode': 404,
            'headers': {},
            'body': json.dumps(str(body))
        }

def get_weight(weight, default_weight):
    return weight if weight != float('Inf') and weight != -float('Inf') else default_weight

def distribute_weighting(hosts):
    default_weight = 0.5
    cumulative_weight = 0

    total_weight = reduce( (lambda a, b: a + b), list( map(lambda host: get_weight(host['weight'], default_weight), hosts) ) )

    # normalize and sort host weights 
    hosts = list( map(lambda host: { 'weight': get_weight(host['weight'], default_weight)/ total_weight, 'name': host['name']}, hosts) )
    hosts.sort(key=lambda host: host['weight'])

    # compute and set cumulative weights
    for host in hosts:
        host['weight'] += cumulative_weight
        cumulative_weight = host['weight']

    return hosts