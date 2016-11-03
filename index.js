const firebase = require('firebase');
const github = require('./github.js');
const repository = require('./models/repository.js');
const R = require('ramda');


firebase.initializeApp({
  databaseURL: 'https://open-sourcery.firebaseio.com/',
  serviceAccount: require('./firebase.credentials.json')
})

const db = firebase.database();


github.getMostStarred()
  .then(R.map(repository.map))
  .then(repos => {

    repos.forEach(repo =>{
      github.getAllLabels(repo)
        .then(R.map(l => l.name))
        .then(labels => {
          repo = Object.assign({}, { labels }, repo);
          db.ref('repositories/' + repo.full_name)
            .set(repo)
        })
    })

  })
  .catch(err => console.error(err))