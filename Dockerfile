FROM golang:1.16-alpine as builder
ENV GO111MODULE=on
RUN mkdir /app
COPY . /app
RUN cd /app && \
  go build -o activity-manager ./cmd/

FROM golang:1.16-alpine
COPY --from=builder /app /app
WORKDIR /app

#EXPOSE 8081

CMD ["/app/activity-manager"]