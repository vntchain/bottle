FROM ubuntu:xenial

ENV PATH=/usr/lib/go-1.10/bin:$PATH

RUN \
    apt-get update && apt-get upgrade -q -y && \
    apt-get install -y --no-install-recommends golang-1.10 xz-utils git make gcc wget libc-dev ca-certificates && \
    (make -p ~/go/src/github.com/vntchain && cd ~/go/src/github.com/vntchain) && \
    git clone --depth 1 -b docker https://github.com/vntchain/bottle && \
    go get -u github.com/kardianos/govendor && \
    (cd bottle && TERM=xterm-256color make bottle) && \
    cp bottle/build/bin/bottle /bottle && \
    apt-get remove -y golang-1.10 git xz-utils make gcc wget libc-dev && apt autoremove -y && apt-get clean && \
    rm -rf /bottle


ENTRYPOINT ["/bottle"]