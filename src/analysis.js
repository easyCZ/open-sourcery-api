// const services = require('./services.js');
const github = require('./github.js');
const repository = require('./models/repository.js');
const R = require('ramda');
const fs = require('fs');

// const db = services.firebaseDB;

const scry = (x) => console.log(x) || x;

const getReposByLabel = (pages = 20) => {
    const labelsForRepos = getLabelsForMostPopularRepos(pages);
    const labelsAndIds = R.map(repoWithLabels => {
        let labels = repoWithLabels.labels;
        let id = repoWithLabels.id;
        return {labels, id};
    }, labelsForRepos)
    return R.groupBy()
}


const getLabelsForMostPopularRepos = (pages = 10) => {
    return github.getMostStarred(1)
        .then(R.map(repository.map))
        .then(repos => {
            let reposWithLabels = repos.map(repo => github.getAllLabels(repo)
                .then(R.map(l => l.name))
                .then(labels => Object.assign({}, { labels }, repo))
                .catch(err => console.error(err))
            );
            return Promise.all(reposWithLabels);
        })
        .catch(err => console.error(err))
}

const dumpLabelsForMostPopularRepos = (pages = 10, filename = 'labels.json') => {
    return getLabelsForMostPopularRepos(pages)
        .then(reposWithLabels => new Promise((resolve, reject) => {
            fs.writeFile('tags.json', JSON.stringify(reposWithLabels), err => {
                if (err) reject(err);
                resolve();
            });
        }));

}

const dumpReposByLabel = (pages = 20, filename = 'repos_by_label.json') => {

    return github.getReposByLabel(pages)
        .then(reposByLabel => new Promise((resolve, reject) => {
            fs.writeFile(filename, JSON.stringify(reposByLabel), err => {
                return err ? reject(err) : resolve();
        })
        }))

}


module.exports = {
    getLabelsForMostPopularRepos,
    dumpLabelsForMostPopularRepos,
    dumpReposByLabel
}