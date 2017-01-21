const firebase = require('firebase');


firebase.initializeApp({
  databaseURL: 'https://open-sourcery.firebaseio.com',
  serviceAccount: './firebase.credentials.json',
})

const firebaseDB = firebase.database();

module.exports = {
    firebaseDB
}