# Super minimalist Dockerfile defining an image that, when run, executes
# a single binary identified by the `SERVICE` build argument
FROM alpine:3.5

# When building the image, pass name of app (usually just the basename of the
# directory containing your app) via the `SERVICE` build argument, e.g.
#
#     docker build --build-arg SERVICE=my-app-name`
ARG SERVICE

# We copy the value of the `SERVICE` build argument to an ENV variable; the ENV
# value will be persisted in the built image, allowing (1) the binary to be
# persisted in the image under a descriptive name (IOW, not "app") and (2) users
# to launch the app by executing simply `docker run <image>`, without (3) also
# requiring users of this generic Dockerfile to hard-code the name of their app
# into the file.
ENV APP=${SERVICE}

RUN apk add --no-cache ca-certificates && mkdir /app
COPY ${SERVICE} /app/${SERVICE}
ENTRYPOINT /app/${APP}

