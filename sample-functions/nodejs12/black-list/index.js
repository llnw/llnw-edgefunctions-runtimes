const list = require('./resources/blacklist');
const url = require('url');

exports.handler = async (request, context) => {

  try {
    if (request.headers && request.headers.Referer) {
      const referer = request.headers.Referer;
      const parsedUrl = url.parse(referer);
      const domain = parsedUrl.host.replace('www.', '').toLowerCase();

      if (list[domain]) {

        const reformattedReferer = parsedUrl.href.replace("?", "&")
        const redirect = `http://error.rakuten.co.jp/403.html?type=frd&ref=${reformattedReferer}`

        return {
          statusCode: 302,
          headers: {
            'Location': redirect
          }
        }
      }
    }

    return  {
      statusCode: 200
    }

  } catch (e) {
    return {
      statusCode: 400,
      body: JSON.stringify(e)
    }
  }

}

