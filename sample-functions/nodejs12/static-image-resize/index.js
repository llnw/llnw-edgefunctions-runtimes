const resize = require('./lib/resize');
const storageService = require('./lib/storage-service')
const path = require('path');
const parsePath = require("parse-path");


exports.handler = function(request, context, callback) {
  // queries
  // url  = source image url (required)
  // path = path in storage for resized image (required)
  // join = join storage path to url path (if present true)
  // x    = width: 100
  // y    = height: 100
  // fmt  = options: 'png'
  // fit  = options: 'cover', 'contain', 'fill', 'inside', 'outside', 'entropy', 'attention' default cover
  // pos  = options: 'centre' 'top' 'right' 'left' 'bottom' 'left,top' default 'centre'
  // hex  = background color: 'ffffff,0.8'

  
  contentTypes = {
    'jpg': 'image/jpeg',
    'jpeg': 'image/jpeg',
    'webp': 'image/webp',
    'tiff': 'image/tiff',
    'png': 'image/png'
  }

  if (!entryCheck(request, callback)) return;

  const url = request.queries.url;
  const parsedUrl = path.parse(parsePath(url).pathname);  
  const format = contentTypes.hasOwnProperty(request.queries) ? request.queries.fmt.toLowerCase() : 'png';
  const joinPath = request.queries.join ? true : false;
  let storagePath = request.queries.path;
  storagePath = joinPath ? path.join(storagePath, parsedUrl.dir) : storagePath;
  const fileName = `${parsedUrl.name}.${format}`
  const width = request.queries.x ? parseInt(request.queries.x) : null;
  const height = request.queries.y ? parseInt(request.queries.y) : null;
  
  resize(url, width, height, format, request.queries).toBuffer().then((outputBuffer) => {

    uploadFile(outputBuffer, storagePath, fileName, contentTypes[format]).then((response) => {
      callback(null, response);
    })
    
  }).catch((err) => {
    callback(null, {
      statusCode: 404,
      headers: {},
      body: err,
    });
    return;
  });
}

const entryCheck = (request, callback) => {
  let failedReason = '';
  if (!request.queries.hasOwnProperty('url')){
    failedReason = "Please provide an image url";
  } else if (!request.queries.hasOwnProperty('path')) {
    failedReason = "Please provide a storage path";
  } else if (!(process.env.HMAC_SECRET && process.env.HMAC_ACCESS) && !(process.env.LLNW_USERNAME && process.env.LLNW_PASSWORD)) {
    failedReason = "Please provide an authentication mechanism";
  }

  if (failedReason) {
    callback(null, {
      statusCode: 404,
      headers: {},
      body: failedReason,
    });
    return false;
  }

  return true;
}

const uploadFile = async (file, storagePath, fileName, contentType) => {

  try {
    const storageApi = new storageService();
    if (!process.env.HMAC_SECRET) {
      await storageApi.authenticate(process.env.LLNW_USERNAME, process.env.LLNW_PASSWORD, 5000)
    } else {
      storageApi.createHmacSignature(process.env.HMAC_ACCESS, process.env.HMAC_SECRET);
    }    

    await storageApi.upload(file, storagePath, fileName)

    return {
      statusCode: 200,
      headers: { 
        'Content-type': contentType
      },
      body: file,
    }
  } catch(e) {

    return {
      statusCode: 404,
      headers: {},
      body: JSON.stringify(e)
    }
  }
}
