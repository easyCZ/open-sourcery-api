const firebase = require('firebase');
const github = require('./github.js');

console.log('gh', github)


firebase.initializeApp({
  databaseURL: 'https://open-sourcery.firebaseio.com/',
  serivceAccount: './service.credentials.json'
})

const db = firebase.database();


// github.getMostStarred()
//   .then()

github.getAllLabels({
  owner: 'FreeCodeCamp',
  repo: 'FreeCodeCamp'
})
.then(res => console.log(res.length));
