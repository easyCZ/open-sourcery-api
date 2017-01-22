const GithubApi = require('github');
const githubCredentials = require('../github.credentials.js');


const github = new GithubApi({
  // debug: true,
  host: 'api.github.com',
  protocol: 'https',
  headers: {
    'user-agent': 'OpenSourcery v0.1'
  }
})

github.authenticate(githubCredentials)


const getLabels = ({ owner, repo }) => github.issues
  .getLabels({ owner, repo })
  .then(labels => labels.map(label => label.name))

const iterate = (promise, items = [], map = (i) => i) =>
  promise.then(response => github.hasNextPage(response)
    ? iterate(github.getNextPage(response), items.concat(map(response)))
    : items.concat(map(response))
)

const getAllLabels = ({owner, repo}) =>
  iterate(github.issues.getLabels({ owner, repo, per_page: 100 }));


const getMostStarred = (page = 1, per_page = 100) => github.search
  .repos({ q: 'stars:>1', per_page, page })
  .then(results => results.items)


/*
 * Get repositories index by label
 * Returns {
 *   'accesibility': ['facebook/react', 'google/hammerrow'],
 *   ...
 * }
 */
const getReposByLabel = (page = 1) => {

  getLabelsForMostPopularRepos(page)
    .then(labels => {
        labels = R.map(repo => ({ labels: repo.labels, id: repo.id }), labels)
        labels = R.map(repo => R.map(label => ({ label, id: repo.id}), repo.labels), labels)
        const flatLabels = R.flatten(labels)

        var reduceToLabelsBy = R.reduceBy((acc, repo) => acc.concat(repo.id), []);
        return reduceToLabelsBy(repo => repo.label)(flatLabels);
    });

}

module.exports = {
  github,
  getMostStarred,
  getAllLabels,
  getReposByLabel
}