FROM sc4h/alpine-s6overlay:v2-3.15

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ENV \
  PATH="/app/cvm:${PATH}" \
  APP_CONFIG="/config/config.yml" \
  APP_LOG="/config/activity.log" \
  APP_VERBOSITY="0"

# Binary
COPY ["dist/cvm_${TARGETOS}_${TARGETARCH}${TARGETVARIANT:+_7}/cvm", "/app/cvm/cvm"]

# Add root files
COPY ["docker/run", "/etc/services.d/cvm/run"]

# Volume
VOLUME ["/config"]

# Port
EXPOSE 8685