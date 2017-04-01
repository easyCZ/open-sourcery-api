from src.services import GithubClient


class GithubClientTest(object):

    def test_client_constructor(self):
        client = GithubClient()
        assert client != None

