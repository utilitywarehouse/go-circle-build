# build template
These files enable quick bootstrapping of an automated build via circle, that auto deploys to the kubernetes dev
cluster.  The build is executed in a `docker run` container, that has the checkout folder synced to it. This causes the
compiled binary to end up in the host system, but compiled with the container library dependencies. The container used
to run the build is
[go-alpine](https://github.com/docker-library/golang/blob/132cd70768e3bc269902e4c7b579203f66dc9f64/1.8/alpine/Dockerfile).
The compiled binary is then copied to the runtime container in `docker build` step. It is then published to the
configured docker registry, and finally the deployment in kubernetes gets PATCHed to use this container.

# usage
clone this repository locally, change it's upstream (or purge the `.git` folder) and configure

# configuration

## circle env vars

the github credentials will be shared into the container (the make command copies them to `~/.netrc`), allowing private
repository access from inside the container. this is needed for `go get`
- GH_USERNAME 
- GH_PASSWORD

- SERVICE 
	- service name most likely you would set this to `basename $PWD`
	- important to note: the deployment in kubernetes is expected to have the same name
- UW_DOCKER_PASS - the registry password for the docker repository (telco user)
- KUBERNETES_TOKEN - the kubernetes token (system user secret in the telecom namespace)

## linting

you can exclude directories from linting (useful for things like generated code etc) by adding a space seperated list of
folder names to the LINT_EXCLUDE environment variable. the target in the makefile will concatenate this into a pipe
delimited string and pass it to the `-e` flag for [gometalinter](https://github.com/alecthomas/gometalinter)

## .gitignore
add the produced binary (normally `basename $PWD`) to `.gitignore`

# running locally

	docker run --rm -e GH_USERNAME -e GH_PASSWORD -e SERVICE -e LINT_EXCLUDE \
		-v $PWD:/go/src/github.com/utilitywarehouse/$SERVICE golang:1.8-alpine \
		sh -c 'apk update && apk add make git gcc musl-dev && cd /go/src/github.com/utilitywarehouse/$SERVICE && make build-container'

this will produce a statically linked executable in this directory, much the same way as `go build` would, with the
difference that this will be statically linked to the musl libraries.
you can also experiment with the different make steps individually, most of them work on the host system. please be
aware that there is a possibility that this will overwrite your $HOME/.netrc file, please check the dependencies in the
makefile, the thing to avoid is the `bootstrap_github_credential` target, if you are bothered about this.
currently these commands are safe:
- make
- make test
- make lint

# todo
make linting configurable (which linters to run)
