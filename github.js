const GithubApi = require('github');
const githubCredentials = require('./github.credentials.js');


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


const getMostStarred = () => github.search.repos({ q: 'stars:>1' })
  .then(results => results.items)


module.exports = {
  getMostStarred,
  getAllLabels
}