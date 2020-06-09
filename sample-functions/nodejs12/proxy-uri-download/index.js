const http = require('follow-redirects').http;

exports.proxy = function(request, context, callback) {
    var uri = process.env['URI'];
    var fnres = {};
    var chunks = [];

    http.get(uri, function(res) {
        console.log('called http get on %s', uri)
        fnres.statusCode = res.statusCode;
        fnres.headers = res.headers;
        res.on('data', function(chunk) {
            chunks.push(Buffer.from(chunk));
        }).on('end', function() {
            var buf = Buffer.concat(chunks);
            fnres.body = buf.toString();
            callback(null, fnres);
        }).on('error', function(err) {
            callback(err, null)
        });
    }).on('error', function(err) {
        console.log('failed to proxy %s : %s', uri, err.message);
    });
    return;
}
