from config import get_firebase_config, FIREBASE_API_KEY, FIREBASE_AUTH_DOMAIN, FIREBASE_DATABASE_URL, FIREBASE_STORAGE_BUCKET


class FirebaseConfigTest(object):

    service_creds = '/random/location/creds.json'
    config = get_firebase_config(service_creds)

    def test_service_creds_passed_through(self):
        assert self.config['serviceAccount'] == self.service_creds

    def test_api_key(self):
        assert self.config['apiKey'] == FIREBASE_API_KEY

