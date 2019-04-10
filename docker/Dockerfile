FROM ubuntu:xenial

ENV PATH=/usr/lib/go-1.10/bin:$PATH

RUN \
    apt-get update && apt-get upgrade -q -y && \
    apt-get install -y --no-install-recommends golang-1.10 git make gcc wget libc-dev ca-certificates && \
    git clone --depth 1 https://github.com/vntchain/bottle && \
    (cd bottle && TERM=xterm-256color make bottle) && \
    cp bottle/build/bin/bottle /bottle && \
    apt-get remove -y golang-1.10 git make gcc libc-dev && apt autoremove -y && apt-get clean && \
    rm -rf /bottle


ENTRYPOINT ["/bottle"]