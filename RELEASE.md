# Release process

This document describes the release process for the Go bindings.

## Description

1. Ensure that the tests on the master are passing. The contributions should be done through PRs with tests valid as security check so it should be the case.

2. Update the CHANGELOG.md with the description of the changes

3. Create a new tag, example:

```sh
git tag v0.0.15
git push --tags
```

4. The CI job will build the artifacts and create a draft release with the artifacts uploaded.

5. Copy the description added in the `CHANGELOG.md` file to the release description.

6. Publish it.

Once published, the artifacts can be downloaded using  the `version`, example: 

`https://github.com/codex-storage/codex-go-bindings/releases/download/v0.0.16/codex-linux-amd64.zip`

It is not recommended to use the `latest` URL because you may face cache issues.

