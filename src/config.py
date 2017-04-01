
FIREBASE_API_KEY = 'AIzaSyCFP9yQ8dcnhOFtHTFAiS5PrAWdamn_2D4'
FIREBASE_AUTH_DOMAIN = 'open-sourcery.firebase.com'
FIREBASE_DATABASE_URL = 'https://open-sourcery.firebaseio.com'
FIREBASE_STORAGE_BUCKET = 'open-sourcery.firebase.com'


def get_firebase_config(service_creds):
    return {
        'apiKey': FIREBASE_API_KEY,
        'authDomain': FIREBASE_AUTH_DOMAIN,
        'databaseURL': FIREBASE_DATABASE_URL,
        'serviceAccount': service_creds,
        'storageBucket': FIREBASE_STORAGE_BUCKET
    }
