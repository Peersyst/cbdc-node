FROM golang:1.23.8 AS base
USER root
RUN apt update && \
    apt-get install -y \
        build-essential \
        ca-certificates
WORKDIR /app
COPY . .
RUN make install


FROM base AS build
ARG VERSION=0.0.0
RUN make build


FROM base AS integration
RUN touch /test.lock

FROM golang:1.23.8 AS release
WORKDIR /
COPY --from=integration /test.lock /test.lock
COPY --from=build /app/bin/cbdcd /usr/bin/cbdcd
ENTRYPOINT ["/bin/sh", "-ec"]
CMD ["cbdcd"]