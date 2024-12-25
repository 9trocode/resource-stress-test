# Stage 1: Build the Next.js application
FROM node:18-alpine AS build

# ARG and ENV for environment variables
ARG TEST_ARG
ARG PORT



ENV TEST_ARG=${TEST_ARG}
ENV PORT=${PORT}


# Print environment variables for debugging (build phase)
RUN echo "Decoded Build-time Environment Variables:" && \
    echo "TEST_ARG=$(echo $TEST_ARG | base64 -d)" && \
     echo "PORT=$(echo $PORT | base64 -d)"

CMD ["sleep", "1000"]
