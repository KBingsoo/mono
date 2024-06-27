FROM public.ecr.aws/docker/library/golang:1.22.3-alpine3.18 AS build-env

RUN apk --no-cache add g++ ca-certificates tzdata build-base git

COPY . /project
WORKDIR /project
RUN go install -tags musl

FROM public.ecr.aws/docker/library/alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /workspace
COPY --from=build-env /go/bin/cards /go/bin/cards
COPY --from=build-env /project/.env /workspace/.env

EXPOSE 8080

CMD ["/go/bin/mono", "server"]