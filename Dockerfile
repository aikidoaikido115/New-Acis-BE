FROM golang:1.24
WORKDIR /app
ARG TARGETARCH
RUN set -eux; \
	atlas_arch="${TARGETARCH:-amd64}"; \
	curl -fsSL "https://release.ariga.io/atlas/atlas-linux-${atlas_arch}-latest" -o /usr/local/bin/atlas; \
	chmod +x /usr/local/bin/atlas

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["sh", "-c", "atlas migrate apply --env dev && ./main"]