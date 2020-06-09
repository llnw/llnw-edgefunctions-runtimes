const _ = require('lodash/core');
const crypto = require('crypto');

let seed = new Date().getTime()

exports.handler = function(request, context, callback) {
  // queries
  // path = required
  // tag  = optional string used to create deterministc hash (used for sticky requests)
  // envvar = optional environment variable name otherwise will default of llnw_hosts
  
  try {  
    seed = `${(48271 * seed) % 2147483647}`
    const path = request.queries.path ? request.queries.path : '';
    const envVar = request.queries.envvar ? request.queries.envvar : 'LLNW_HOSTS';
    let tag = request.queries.tag ? request.queries.tag : seed;    

    if (!process.env[envVar]) throw Error('No environment vars found');

    let hosts = JSON.parse(process.env[envVar]);

    if(!_.isArray(hosts)) throw Error('No matching env vars found');
    
    hosts = distributeWeighting(hosts);
    
    if (tag === undefined) {
      tag = seed
    }
    
    let hmac = crypto.createHmac('md5', tag);
    hmac.write(envVar);
    const md5Hash =hmac.digest('hex').substr(0,8);   
    const hashAsInt = parseInt("0x" + md5Hash, 16);
    const maxInt = parseInt("0xffffffff", 16);
    const random = hashAsInt / maxInt;
    
    const filtered = hosts.filter((t) => {
      return t.weight > random;
    });
    
    const redirect = `${filtered[0].name}/${path}`;
    
    callback(null, {
      statusCode: 302,
      headers: {
        'Location': redirect
      },
      body: '',
    });

    return;
  } catch (e) {
    callback(null, {
      statusCode: 404,
      headers: {},
      body: JSON.stringify(e),
    });

    return;
  }
}

function distributeWeighting (hosts) {
  const defaultWeight = 0.5;
  let cumulativeWeight = 0

  var totalWeight = _.map(hosts, (host) => {
    return _.isFinite(host.weight) ? host.weight : defaultWeight;
  }).reduce((a, b) => {
    return a + b
  });
  
  hosts = _.map(hosts, (host) => {
    if (host.name && typeof host.name === 'string') {
      return {
        weight: (_.isFinite(host.weight) ? host.weight : defaultWeight) / totalWeight,
        name: host.name
      }
    } else {
      throw Error('Poorly formated hosts object');
    }
  }).sort((a,b) => {
    return a.weight > b.weight;
  });

  hosts.forEach((host) => {
    host.weight += cumulativeWeight;
    cumulativeWeight = host.weight;
  });

  return hosts;
}
