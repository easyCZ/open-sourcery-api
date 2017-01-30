import os
from github import Github as GithubClient


class Github(GithubClient):

    def __init__(self):
        print(os.environ['GITHUB_KEY'], os.environ['GITHUB_SECRET'])
        super().__init__(
            client_id=os.environ['GITHUB_KEY'],
            client_secret=os.environ['GITHUB_SECRET'],
            per_page=100,
            user_agent='OpenSourcery v0.1'
        )

    def repos_by_stars(self):
        return self.search_repositories('stars:>1', 'stars')

    def labels_by_stars(self):
        """
        yields a generator of (label, repo) tuples
        """
        for repo in self.repos_by_stars():
            yield ((label, repo) for label in repo.get_labels())