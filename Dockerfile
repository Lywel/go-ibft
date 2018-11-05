FROM golang:latest as builder
WORKDIR /go-modules

# create ssh directory
RUN mkdir ~/.ssh
RUN touch ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

# allow private repo pull
ARG github_oauth
RUN git config --global url."https://$github_oauth:x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY . ./
RUN GOOS=linux go build -a -installsuffix cgo -ldflags "-linkmode external -extldflags -static"

FROM scratch
EXPOSE 3000
COPY --from=builder /go-modules/app .
ENTRYPOINT ['/app']

