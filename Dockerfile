FROM alpine:latest

# Declare build arguments
ARG MY_BUILD_ARG
ARG ANOTHER_ARGS

# Echo build arguments to console during build
RUN echo "MY_BUILD_ARG=${MY_BUILD_ARG}" && \
    echo "ANOTHER_ARGS=${ANOTHER_ARGS}"

# Optional: Set environment variables based on build arguments
ENV MY_ENV_VAR=${MY_BUILD_ARG}
ENV ANOTHER_ENV_VAR=${ANOTHER_ARG}

CMD ["sh"]
