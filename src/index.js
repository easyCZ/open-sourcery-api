const services = require('./services.js');
const github = require('./github.js');
const repository = require('./models/repository.js');
const R = require('ramda');

const db = services.firebaseDB;

const scry = (x) => console.log(x) || x;


github.getMostStarred(4)
    .then(R.map(repository.map))
    .then(repos => {
        let reposWithLabels = repos.map(repo => github.getAllLabels(repo)
            .then(R.map(l => l.name))
            .then(labels => Object.assign({}, { labels }, repo))
            .catch(err => console.error(err))
        );
        return Promise.all(reposWithLabels);
    })
    // .then(scry)
    .then(repos => Promise.all(repos.map(repo => db
        .ref('repositories/' + repo.id)
        .set(repo)
        .then(updatedValue => console.log('saved', repo.full_name))
    )))
    .then(() => console.log('Done'))
    .then(() => process.exit(0))
    .catch(err => console.error(err))