FROM alpine:3.19
WORKDIR /app

# Create a non-root user for OpenShift
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy the pre-built binary
COPY main .

# Expose backend port
EXPOSE 8080

# Run the backend
CMD ["./main"]