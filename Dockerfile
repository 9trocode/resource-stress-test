FROM alpine:latest

ARG MY_BUILD_ARG
ARG ANOTHER_ARG

RUN echo "MY_BUILD_ARG=${MY_BUILD_ARG}" > /build_args.txt && \
    echo "ANOTHER_ARG=${ANOTHER_ARG}" >> /build_args.txt

CMD ["sh"]
