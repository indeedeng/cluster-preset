# Contributing

All contributions are welcome to this project!

The general process for contributing code is as follows.

1. Open an issue describing the desired feature, bug fix, or improvment you'd like to see before making any changes.
1. Assign the ticket to yourself.
1. Fork the repository and make the necessary changes on a branch identified by the issue (e.g. `gh-##`). This ensures the branch is linked to the corresponding issue.
1. Open a [pull request](https://help.github.com/en/articles/about-pull-requests) for the maintainers of the repository to review.
1. Be sure to add the maintainers to the pull request as reviewers to ensure a speedy review.
1. Upon reviewing the changes, your change may be:
    * merged, requiring no further changes on your part
    * annotated with a set of changes that must be made before the change will be merged.

## Unit Tests

All changes submitted must include updated unit tests (excluding documentation and Kubernetes configuration).
In order for your change to be accepted, the [TravisCI](https://travis-ci.com/indeedeng/cluster-preset) pipeline must complete successfully.

## Code of Conduct

ClusterPreset is governed by the [Contributor Covenant v1.4.1](CODE_OF_CONDUCT.md).
