# SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
# SPDX-License-Identifier: Apache-2.0

################################################################################

FROM golang:1.25 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY internal/ internal/
COPY cmd/ cmd/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o /workspace ./...

################################################################################

# Ref: https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static:nonroot AS scheduler
WORKDIR /
COPY --from=builder /workspace/scheduler .
USER 65532:65532
ENTRYPOINT ["/scheduler"]

################################################################################

# Ref: https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static:nonroot AS controllers
WORKDIR /
COPY --from=builder /workspace/controllers .
USER 65532:65532
ENTRYPOINT ["/controllers"]

################################################################################

# Ref: https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static:nonroot AS admission
WORKDIR /
COPY --from=builder /workspace/admission .
USER 65532:65532
ENTRYPOINT ["/admission"]
