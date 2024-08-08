FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go get sigs.k8s.io/controller-runtime@v0.15.1
RUN go mod tidy
RUN go build -o k8s-security-policy-auditor .
FROM gcr.io/distroless/base
COPY --from=builder /app/k8s-security-policy-auditor /k8s-security-policy-auditor
ENTRYPOINT ["/k8s-security-policy-auditor"]
