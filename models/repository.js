const R = require('ramda')


const map = (repo) => ({
  id: repo.id,
  owner: repo.owner.login,
  repo: repo.name,
  full_name: repo.full_name,
  description: repo.description,
  stargazers_count: repo.stargazers_count,
  forsk: repo.forks
})


module.exports = {
  map
}