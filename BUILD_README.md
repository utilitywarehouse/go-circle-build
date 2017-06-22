# Minimalist, boilerplate build templates that almost just work

The `circle.yml`, `Makefile`, and `Dockerfile` files in this repo are designed
to be dropped as-is into a new Go project. If all goes as planned, two simple
steps will result in a working CI setup that:

  - produces a tested Go binary that can be run in an Alpine Linux Docker
    container, and

  - produces and publishes an Alpine Linux-based Docker image that runs that
    binary

If you are feeling ambitious, you can also uncomment the last task in the
`deployment.master.commands` section of `circle.yml` and then you will also
have a setup that auto-deploys to the Kubernetes `dev` cluster.

## Usage

  1. Drop those three files into your project root. Then edit `circle.yml` and
     replace `<service-name>` with the name of your app. (This should almost
     certainly be the same as the name of your project's GitHub repository,
     which should be the same as the basename of the directory containing your
     app's source code.)

  2. Add a Circle CI project for your app and define at least the following
     three environment variables for it:

       - `UW_DOCKER_PASSWORD` (password for Docker registry)
       - `GH_USERNAME` and `GH_PASSWORD` (credentials for GitHub repo)

     If your app doesn't belong to the _telecom_ domain, you will also
     need to define the `UW_DOCKER_USERNAME` and `NAMESPACE` environment
     variables. And if you're setting up auto-deploy, you'll also need to
     define `KUBERNETES_TOKEN` (which is the system user secret in your k8s
     namespace).

## How it works

  - The included Makefile contains an `all` target that, when invoked in your
    project's root folder, will install your app's dependencies, lint your
    code, run your tests, and if all goes well, `go build` a new binary for
    your app. As usual with Go apps, the binary artifact will be named after
    the directory containing your source code.

  - The included `circle.yml` file contains a custom `test` task that will
    launch a [go-alpine][1] container, sync the project's base directory
    (`$PWD`) into the container, install the minimal set of required apk's in
    the container, and finally invoke `make all` inside the container.

  - A `circle.yml` deployment task will then build a Docker image including the
    resulting binary and publish it to the `registry.uw.systems` Docker
    registry.

## Notes

- The GitHub credentials will be shared into the container (the make command
  copies them to `~/.netrc`), allowing private repository access from inside
  the container. This is needed for `go get`.
- The `circle.yml` deployment task for auto-deployment to k8s assumes that your
  k8s deployment will be named `$SERVICE`; if this is not true you'll need to
  modify `circle.yml` (or change your k8s deployment name).
- The same `NAMESPACE` is used to publish images to the Docker registry and
  deploy to k8s. If your app uses different namespaces you'll need to
  tweak this.

## Linting

You can exclude directories from linting (useful for things like generated
code, etc.) by adding a space-separated list of directory names to the
`LINT_EXCLUDE` environment variable. The target in the Makefile will
concatenate this into a pipe delimited string and pass it to the `-e` flag for
[gometalinter](https://github.com/alecthomas/gometalinter).

## Running locally

You can simply copy/paste the `test.override` task from `circle.yml` onto the
command-line if you would like to install/test/build your app in a container on
your local machine. (You may have to set some environment variable to get it to
work.) E.g.:

    $ export SERVICE=$(basename $PWD)
    $ docker run --rm -e GH_USERNAME -e GH_PASSWORD -e SERVICE -e LINT_EXCLUDE \
      -v $PWD:/go/src/github.com/utilitywarehouse/$SERVICE golang:1.8-alpine \
      sh -c 'apk update && apk add make git gcc musl-dev &&
      cd /go/src/github.com/utilitywarehouse/$SERVICE && make all'

This will produce a statically linked executable in your current working
directory, much the same way as `go build` would, with the difference that this
will be statically linked to the `musl` libraries.

You can also experiment with the different make targets individually; most of
them work on the host system. E.g.

    $ make test

Please be aware that this may overwrite your `$HOME/.netrc` file (please check
the dependencies in the Makefile if you are bothered about this; the thing to
avoid is the `bootstrap_github_credential` target).

## Todo

Make linting configurable (which linters to run)

[1]: https://github.com/docker-library/golang/blob/132cd70768e3bc269902e4c7b579203f66dc9f64/.8/alpine/Dockerfile
