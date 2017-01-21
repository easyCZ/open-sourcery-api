const R = require('ramda');
const firebaseDB = require('../services').firebaseDB;
const fs = require('fs');



const list = () => firebaseDB
    .ref('repositories')
    .once('value')
    .then(snapshot => snapshot.val())
    .then(repositories => {

        let labels = [];

        Object.keys(repositories).forEach(repositoryOwner => {
            Object.keys(repositories[repositoryOwner]).forEach(repositoryName => {
                labels = labels.concat(repositories[repositoryOwner][repositoryName].labels);
            })
        })

        return labels;
    })


const listWithCounts = () => list().then(R.countBy(R.identity))

listWithCounts()
    .then(tags => {
        const pairs = R.toPairs(tags)
        const sorted = R.sortBy(pair => pair[1], pairs);
        fs.writeFile('tags.json', JSON.stringify(sorted), err => {
            console.log('error', err);
            console.log('done');
            process.exit(0);
        })
    })
    .catch(err => console.error(err))

module.exports = {
    list
}