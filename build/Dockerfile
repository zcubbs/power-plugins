FROM --platform=$TARGETPLATFORM alpine:latest

ARG TARGETPLATFORM
ARG TARGETARCH

# Copy binaries from buildx
# TARGETPLATFORM is set by buildx
COPY bin/$TARGETPLATFORM/ /plugins/

CMD ["sh"]
