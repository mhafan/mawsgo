FROM public.ecr.aws/amazonlinux/amazonlinux:latest as BUILD

#
ENV DEBIAN_FRONTEND=noninteractive

# build as root
USER 0

# Go: golang
# C++: gcc-c++ make libcurl-devel
RUN yum update -y && yum -y install golang

#
WORKDIR /src
ADD . .
#
RUN go mod tidy
RUN go build -tags lambda.norpc -ldflags="-s -w" -o main

#####
FROM public.ecr.aws/amazonlinux/amazonlinux:latest

#
COPY --from=BUILD /src/main .

ENTRYPOINT [ "./main" ]
