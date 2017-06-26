const functions = require('firebase-functions');

// Create and Deploy Your First Cloud Functions
// https://firebase.google.com/docs/functions/write-firebase-functions

exports.helloWorld = functions.https.onRequest((request, response) => {

 return fetch('http://github.com')
     .then(r => r.text())
     .then(text => console.log(text))
     .then(r => response.send("Hello from Firebase!"))

});
