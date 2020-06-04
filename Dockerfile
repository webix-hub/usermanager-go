FROM centurylink/ca-certs
WORKDIR /app
COPY ./usermanager /app
COPY ./demodata /app/demodata

CMD ["/app/usermanager"]