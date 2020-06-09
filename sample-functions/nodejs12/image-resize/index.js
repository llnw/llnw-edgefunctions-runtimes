const resize = require('./lib/resize');

exports.handler = function(payload, context, callback) {
  // queries
  // url = original image url (required)
  // x   = width: 100
  // y   = height: 100
  // fmt = options: 'png'
  // fit = options: 'cover', 'contain', 'fill', 'inside', 'outside', 'entropy', 'attention' default cover
  // pos = options: 'centre' 'top' 'right' 'left' 'bottom' 'left,top' default 'centre'
  // hex = background color: 'ffffff,0.8'
  
  contentTypes = {
    'jpg': 'image/jpeg',
    'jpeg': 'image/jpeg',
    'webp': 'image/webp',
    'tiff': 'image/tiff',
    'png': 'image/png'
  }

  if (payload.queries.url) {
    const url = payload.queries.url;
    const width = payload.queries.x ? parseInt(payload.queries.x) : null;
    const height = payload.queries.y ? parseInt(payload.queries.y) : null;
    const format = contentTypes.hasOwnProperty(payload.queries) ? payload.queries.fmt.toLowerCase() : 'png';
    resize(url, width, height, format, payload.queries).toBuffer().then((outputBuffer) => {

      callback(null, {
        statusCode: 200,
        headers: { 
          'Content-type': contentTypes[format]
        },
        body: outputBuffer.toString('base64'),
        shouldBase64Decode: true,
      }); 
      return;
    }).catch((err) => {
      callback(null, {
        statusCode: 404,
        headers: {},
        body: err,
      });
      return;
    });
  } else {
    callback(null, {
      statusCode: 404,
      headers: {},
      body: "Please provide an image url",
    });
    return;
  }
}
