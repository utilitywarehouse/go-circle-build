FROM alpine:3.5
ARG SERVICE
ENV APP=${SERVICE}
RUN apk add --no-cache ca-certificates && mkdir /app

COPY $SERVICE /app/${SERVICE}

ENTRYPOINT /app/${APP}

