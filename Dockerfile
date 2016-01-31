# Publish Archlinux packages built by maze-build
#
#     docker build --rm=true -t mikkeloscar/maze-publish .

FROM scratch
MAINTAINER Mikkel Oscar Lyderik <mikkeloscar@gmail.com>

# Add binary
ADD maze-publish /

ENTRYPOINT ["/maze-publish"]
