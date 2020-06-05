const rp = require('request-promise-native');
const crypto = require('crypto');
'use strict';

module.exports = class StorageService {

  constructor() {
    this.token = '';
    this.signature = '';
    this.headers = {
      'Accept': 'application/json',
      'X-Agile-Recursive': true
    };
  }

  async authenticate (username, password, expiry) {
    const response = await rp.post('http://llnwstor-l.upload.llnw.net/jsonrpc2', {
      body: JSON.stringify({
      id: 1,
      jsonrpc: '2.0',
      method: 'authenticate',
      params: {
        username,
        password,
        expiry,
      }}),
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' 
      }
    })

    this.token = await JSON.parse(response).result.token
  }

  createHmacSignature(access_key, secret_key) {
    const expiry = new Date().getTime() + 10
    const endpoint = '/post/file'
    const mac = [
      `access_key=${access_key}`,
      `expiry=${expiry}`
    ]
    
    Object.keys(this.headers).forEach((header) => {
      if (header.startsWith('X-Agile-')) {
        mac.push(`${header.slice(8).toLowerCase()}=${this.headers[header]}`);
      }
    });

    mac.sort();
    const payload = `${endpoint}?${mac.join('&')}`;
    const hmac = crypto.createHmac('sha256', secret_key);
    hmac.write(payload);

    this.signature = `${payload}&signature=${hmac.digest('base64')}`;
  }

  async upload (body, path, name) {

    if (this.token) {
      this.headers['X-Agile-Authorization'] = this.token;
    } else {
      this.headers['X-Agile-Signature'] = this.signature;
    }

    await rp.post({
      url: 'http://llnwstor-l.upload.llnw.net/post/file',
      resolveWithFullResponse: true,
      headers: this.headers,
      formData: {
        uploadFile: {
          value: body,
          options: {
            filename: name,
          }
        },
        directory: path,
      }
    });
  }
}