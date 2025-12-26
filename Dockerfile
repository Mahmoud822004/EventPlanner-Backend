FROM alpine:3.19
WORKDIR /app

# Copy the pre-built binary first (as root)
COPY main .

# Make it executable
RUN chmod +x main

# Create a non-root user for OpenShift
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose backend port
EXPOSE 8080

# Run the backend
CMD ["./main"]