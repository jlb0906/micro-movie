FROM alpine
RUN mkdir /app
ADD config-grpc-srv /app
ADD conf/ /app
WORKDIR /app
ENTRYPOINT [ "/app/config-grpc-srv" ]