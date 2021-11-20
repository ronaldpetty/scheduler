FROM ubuntu:20.04 as tlsbuilder
RUN apt update
RUN apt install openssl -y
RUN openssl genrsa -out server.key 2048
RUN openssl ecparam -genkey -name secp384r1 -out server.key
RUN openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650  -subj "/C=US/ST=California/L=San Francisco/O=RX-M/OU=K8s/CN=kube-scheduler"

FROM golang:1.17 as gobuilder
ADD . /code
WORKDIR /code
RUN go build -tags "netgo" -o kube-scheduler *.go

FROM scratch
COPY --from=tlsbuilder server.key /server.key
COPY --from=tlsbuilder server.crt /server.crt

WORKDIR /
ENV PATH=/
COPY --from=gobuilder /code/kube-scheduler /kube-scheduler
ENTRYPOINT ["nothing-to-see-here"]
