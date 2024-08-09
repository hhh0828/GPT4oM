# syntax=docker/dockerfile:1

FROM gpt4om:latestreal
FROM httpd

# Set destination for COPY


#/
# Download Go modules
COPY ./Static/index.html usr/local/apach2/htdocs


# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
#COPY *.go ./

# Build
#RUN CGO_ENABLED=0 GOOS=linux go build -o /pracforgpt

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
#EXPOSE 8080

# Run
#CMD ["/pracforgpt"]