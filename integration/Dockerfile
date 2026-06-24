# Integration harness: build the yup-tac binary for Linux, then run it inside a
# Debian image (GNU coreutils) so its flags, inputs, and outputs can be checked
# against the real `tac` reference. Invoked by `make integration`.
FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/yup-tac .

FROM debian:stable-slim
COPY --from=build /out/yup-tac /usr/local/bin/yup-tac
COPY integration/run.sh /run.sh
ENTRYPOINT ["/bin/sh", "/run.sh"]
