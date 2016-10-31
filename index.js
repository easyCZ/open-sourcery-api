const firebase = require('firebase');
const GithubApi = require('github');
const githubCredentials = require('./github.credentials.js');


firebase.initializeApp({
  databaseURL: 'https://open-sourcery.firebaseio.com/',
  serivceAccount: './service.credentials.json'
})

const db = firebase.database();
const github = new GithubApi({
  host: 'api.github.com',
  protocol: 'https',
  headers: {
    'user-agent': 'OpenSourcery v0.1'
  }
}).authenticate({

})