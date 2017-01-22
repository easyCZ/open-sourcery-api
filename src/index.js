const argparse = require('command-line-args');


const argumentDefinitions = [
  { name: 'verbose', alias: 'v', type: Boolean },
  { name: 'src', type: String, multiple: true, defaultOption: true },
  { name: 'timeout', alias: 't', type: Number }
]








const temp = () => {
    return github.getMostStarred(4)
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
}
