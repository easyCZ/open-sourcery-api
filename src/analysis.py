import collections

from services import Github


gh = Github()


def repos_by_label(repos=1000):
    index = collections.defaultdict(list)
    repo_labels_gen = zip(range(repos), gh.labels_by_stars())

    for (i, labels_gen) in repo_labels_gen:
        for label, repo in labels_gen:
            index[label.name].append(repo.full_name)

    return index
