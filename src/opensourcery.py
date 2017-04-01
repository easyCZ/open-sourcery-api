import argparse
import os


DESCRIPTION = 'OpenSourcery'
DEST_HELP = 'The destination of output'

INDEX_ACTION = 'index'

ACTION_HELP = 'The type of action to take'
ACTION_CHOICES = [INDEX_ACTION]
ACTION = 'action'

EXPECTED_ENV_VARS = [
    'GITHUB_KEY',
    'GITHUB_SECRET',
    'FIREBASE_KEY'
]


def check_env():
    print(
        'Missing ENV VARIABLES: ',
        [key for key in EXPECTED_ENV_VARS if key not in os.environ])


def main():
    check_env()

    parser = argparse.ArgumentParser(description=DESCRIPTION)
    parser.add_argument(ACTION, choices=ACTION_CHOICES, help=ACTION_HELP)
    args = parser.parse_args()

    if args.action == INDEX_ACTION:
        pass


if __name__ == '__main__':
    main()
