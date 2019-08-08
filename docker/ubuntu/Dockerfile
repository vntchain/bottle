FROM ubuntu:xenial

ENV PATH=/usr/lib/go-1.10/bin:/root/go/bin:$PATH

ENV GOPATH=/root/go

RUN \
    apt-get update && apt-get upgrade -q -y && \
    apt-get install -y --no-install-recommends libxml2-dev golang-1.10 xz-utils git make gcc wget libc-dev ca-certificates && \
    mkdir -p ~/go/src/github.com/vntchain && \
    (cd ~/go/src/github.com/vntchain && git clone --depth 1 https://github.com/vntchain/bottle -b docker -v ) && \
    (cd ~/go/src/github.com/vntchain/bottle && TERM=xterm-256color make bottle) && \
    cp ~/go/src/github.com/vntchain/bottle/build/bin/bottle /bottle && \
    apt-get remove -y golang-1.10 git xz-utils make gcc wget libc-dev && apt autoremove -y && apt-get clean 

ENTRYPOINT ["/bottle"]
