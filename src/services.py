import os
import pyrebase
from github import Github as GithubClient


class Github(GithubClient):

    def __init__(self):
        """
        OpenSourcery custom github access layer
        """
        super().__init__(
            client_id=os.environ['GITHUB_KEY'],
            client_secret=os.environ['GITHUB_SECRET'],
            per_page=100,
            user_agent='OpenSourcery v0.1'
        )

    def repos_by_stars(self):
        """
        Find all repositories ordered by stars
        """
        return self.search_repositories('stars:>1', 'stars')

    def labels_by_stars(self):
        """
        yields a generator of (label, repo) tuples
        """
        for repo in self.repos_by_stars():
            yield ((label, repo) for label in repo.get_labels())


class Firebase(object):

    def __init__(self):
        config = {}
        self.firebase = pyrebase.initialize_app(config)
        self.db = self.firebase.database()

    def get_repository(self, id):
        repo = self.db.child('repositories').get()
        return repo

    def get_or_create_repo(self, repo):
        pass
