const request = require('request');
const sharp = require('sharp');
const hexRgb = require('hex-rgb');

function _formatRgb(hex) {
  hex = hex.split(',');
  const rgb = hexRgb(hex[0]);

  return {
    r: rgb.red,
    g: rgb.green,
    b: rgb.blue,
    alpha: hex[1] ? parseFloat(hex[1]) : 1
  };
}

function _formatRequestArgs(queryStringParams) {
  // fit = options: 'cover', 'contain', 'fill', 'inside', 'outside', 'entropy', 'attention' default cover
  // pos = options: 'centre' 'top' 'right' 'left' 'bottom' 'left,top' default 'centre'
  // hex = background color: 'ffffff,0.8'
  
  let position = queryStringParams.pos ? queryStringParams.pos.replace(',', ' ') : null;

  if (queryStringParams.fit === 'entropy') {
    queryStringParams.fit = 'cover';
    position = sharp.strategy.entropy;    
  } else if (queryStringParams.fit === 'attention') {
      queryStringParams.fit = 'cover';
      position = sharp.strategy.attention
  }


  return {
    fit: queryStringParams.fit ? queryStringParams.fit : null,
    background: queryStringParams.hex ? _formatRgb(queryStringParams.hex) : null,
    position, 
  }
}

module.exports = function resize(url, width, height, format, queryStringParams) {
  const readStream = request({url: url}, (err, resp, bodyBuffer) => {
    return bodyBuffer;
  });

  sharp.cache( { memory: 128 } )

  let transform = sharp().ensureAlpha();

  const options = _formatRequestArgs(queryStringParams);

  if (options.background) {
    transform = transform.flatten(options)
  }

  if (format) {
    transform = transform.toFormat(format);
  }

  if (width || height) {
    transform = transform.resize(width, height, options);
  }

  return readStream.pipe(transform);

}
