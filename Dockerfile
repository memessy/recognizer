FROM golang:buster as builder
WORKDIR /src
COPY . .
RUN apt-get update -qq \
    && apt-get install -y -qq \
    libtesseract-dev \
    libleptonica-dev
RUN go build -o /out/recognizer .

FROM debian:buster-slim as runner

ENV PORT=8000
ENV FILE_KEY=file
ENV FILE_MAX_SIZE=10000000

ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get update -qq \
    && apt-get install -y -qq \
    tesseract-ocr-eng \
    tesseract-ocr-rus

COPY --from=builder /out/recognizer /

ENTRYPOINT ["/recognizer"]