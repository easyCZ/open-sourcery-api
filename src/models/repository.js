const R = require('ramda')


const fullNameToId = R.replace(/(\.|\#|\$|\[|\])/g, '_')

const map = (repo) => ({
  id: fullNameToId(repo.full_name),
  owner: repo.owner.login,
  repo: repo.name,
  full_name: repo.full_name,
  description: repo.description,
  stargazers_count: repo.stargazers_count,
  forsk: repo.forks
})




module.exports = {
  map,
}