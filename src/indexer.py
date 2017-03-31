import yaml
import os


def get_project_configs(directory):
    for config_filename in os.listdir(directory):
        path = os.path.join(directory, config_filename)
        with open(path) as config:
            yield yaml.parse(config)


def index(repository):
   owner = repository['owner']
   repo_name = repository['name']

   # Get repository from firebase

   # If already created, run a if-modified-since request to update fields

   # If not, get full details from github

   # Store updated repository details to firebase

   # For each label, find open issues
   # Store each issue in firebase
