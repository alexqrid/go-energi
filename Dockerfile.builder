# Ubuntu 20.04 is our base image for building
FROM ubuntu:20.04

# set up timezone
ENV TZ=GMT
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# update software
RUN apt -y --fix-missing update
RUN apt -y full-upgrade
RUN apt -y autoremove
RUN apt -y clean

# install docker
RUN apt -y update
RUN apt -y install curl gnupg lsb-release software-properties-common git build-essential wget awscli sudo

# golang variables
ARG golang_version="1.19.4"
ARG golang_hostarch="linux-amd64"
ARG golang_filename="go${golang_version}.${golang_hostarch}.tar.gz"
ARG golang_url="https://golang.org/dl/${golang_filename}"
ARG golang_sha256="c9c08f783325c4cf840a94333159cc937f05f75d36a8b307951d5bd959cf2ab8"

# install golang
RUN wget -nv ${golang_url}
RUN echo "${golang_sha256} ${golang_filename}" > "${golang_filename}.sha256"
RUN sha256sum -c ${golang_filename}.sha256
RUN tar -C /usr/local -xzf ${golang_filename}
RUN rm -rf ${golang_filename}*
ENV PATH="${PATH}:/usr/local/go/bin"
ENV GOROOT="/usr/local/go"

# install go-junit-report
RUN go install -v github.com/RyanLucchese/go-junit-report@latest
ENV PATH="${PATH}:/root/go/bin"

RUN mkdir -p "/energi"
WORKDIR "/energi"
ADD Makefile.release Makefile.release
RUN make -f Makefile.release release-tools
ENV GOPATH="/energi"
ENV GOBIN="/energi/bin"
ENV GO111MODULE="on"
ENV GOFLAGS="-v"
