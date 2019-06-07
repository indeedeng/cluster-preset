FROM debian:stretch-slim

# TODO: download from github releases
COPY cluster-preset /usr/local/bin/cluster-preset

RUN useradd -ms /bin/sh cluster-preset
WORKDIR /home/cluster-preset
USER cluster-preset

ENTRYPOINT [ "cluster-preset" ]
