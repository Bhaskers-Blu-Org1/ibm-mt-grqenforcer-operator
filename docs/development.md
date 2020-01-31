# Development Guide

## Prerequisite

- git
- go version v1.12+
- Linting Tools

    | linting tool | version |
    | ------------ | ------- |
    | [hadolint](https://github.com/hadolint/hadolint#install) | [v1.17.2](https://github.com/hadolint/hadolint/releases/tag/v1.17.2) |
    | [shellcheck](https://github.com/koalaman/shellcheck#installing) | [v0.7.0](https://github.com/koalaman/shellcheck/releases/tag/v0.7.0) |
    | [yamllint](https://github.com/adrienverge/yamllint#installation) | [v1.17.0](https://github.com/adrienverge/yamllint/releases/tag/v1.17.0)
    | [helm client](https://helm.sh/docs/using_helm/#install-helm) | [v2.10.0](https://github.com/helm/helm/releases/tag/v2.10.0) |
    | [golangci-lint](https://github.com/golangci/golangci-lint#install) | [v1.18.0](https://github.com/golangci/golangci-lint/releases/tag/v1.18.0) |
    | [autopep8](https://github.com/hhatto/autopep8#installation) IBMDEV SKIPPED | [v1.4.4](https://github.com/hhatto/autopep8/releases/tag/v1.4.4) |
    | [mdl](https://github.com/markdownlint/markdownlint#installation) | [v0.5.0](https://github.com/markdownlint/markdownlint/releases/tag/v0.5.0) |
    | [awesome_bot](https://github.com/dkhamsing/awesome_bot#installation) | [1.19.1](https://github.com/dkhamsing/awesome_bot/releases/tag/1.19.1) |
    | [sass-lint](https://github.com/sasstools/sass-lint#install) IBMDEV SKIPPED | [v1.13.1](https://github.com/sasstools/sass-lint/releases/tag/v1.13.1) |
    | [tslint](https://github.com/palantir/tslint#installation--usage) IBMDEV SKIPPED | [v5.18.0](https://github.com/palantir/tslint/releases/tag/5.18.0)
    | [prototool](https://github.com/uber/prototool/blob/dev/docs/install.md) | `7df3b95` |
    | [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) IBMDEV SKIPPED | `3792095` |

## Developer quick start

- Setup `GIT_HOST` to override the setting for your custom path.

```bash
export GIT_HOST=github.com/<YOUR_GITHUB_ID>
```

- Run the `linter` and `test` before building the binary.

```bash
make check
make test
make build
```

- Build and push the docker image for local development.

<!-- IBMDEV Correct build-push-images to build-push-image (no 's') -->
<!-- IBMDEV Add export of DOCKER_USERNAME and DOCKER_PASSWORD -->
<!-- IBMDEV this section seems completely wrong -- Makefile uses IMAGE_REPO and IMAGE_NAME, not REGISTRY and IMG -->
<!-- IBMDEV replace REGISTRY and IMG with IMAGE_REPO and IMAGE_NAME -->

```bash
export IMAGE_NAME=<YOUR_CUSTOMIZED_IMAGE_NAME>
export IMAGE_REPO=<YOUR_CUSTOMIZED_IMAGE_REGISTRY>
export DOCKER_USERNAME=<YOUR_DOCKER_USERNAME>
export DOCKER_PASSWORD=<YOUR_DOCKER_PASSWORD>
make build-push-images
```

> **Note:** You need to login the docker registry before running the command above.
