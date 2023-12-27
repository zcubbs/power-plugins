FROM --platform=$TARGETPLATFORM alpine:latest

ARG TARGETPLATFORM
ARG TARGETARCH

# Copy *.so files from the previous stage
COPY bin/$TARGETARCH/*.so /plugins/

CMD ["sh"]
