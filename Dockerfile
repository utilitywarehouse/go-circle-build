FROM alpine:3.5
ARG SERVICE
RUN apk add --no-cache ca-certificates && mkdir /app
WORKDIR /app

COPY $SERVICE $SERVICE

ENTRYPOINT "/app/$SERVICE"

