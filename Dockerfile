FROM golang:1.22
WORKDIR /app
COPY . ./
RUN make
EXPOSE 9999
CMD ["/app/bin/shortener"]
