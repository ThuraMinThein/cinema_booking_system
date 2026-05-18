PROTO_DIR      := proto
PROTO_SRC      := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT         := .
SERVICES       := gateway bookings users seats
NAMESPACE      := cinema-booking
DOCKERHUB_USER ?= your-dockerhub-username
IMAGE_TAG      ?= latest

# ── Proto generation ───────────────────────────────────────────────────────────
.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

# ── Docker (local builds) ──────────────────────────────────────────────────────
.PHONY: docker-build
docker-build:
	@for svc in $(SERVICES); do \
		echo "Building cinema-$$svc:$(IMAGE_TAG)..."; \
		docker build -f $$svc/Dockerfile -t $(DOCKERHUB_USER)/cinema-$$svc:$(IMAGE_TAG) .; \
	done

.PHONY: docker-push
docker-push:
	@for svc in $(SERVICES); do \
		echo "Pushing cinema-$$svc:$(IMAGE_TAG)..."; \
		docker push $(DOCKERHUB_USER)/cinema-$$svc:$(IMAGE_TAG); \
	done

# ── Docker Compose (local dev) ─────────────────────────────────────────────────
.PHONY: compose-up
compose-up:
	docker compose up --build -d

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: compose-logs
compose-logs:
	docker compose logs -f

# ── Kubernetes ─────────────────────────────────────────────────────────────────
.PHONY: k8s-apply
k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmaps/ -n $(NAMESPACE)
	kubectl apply -f k8s/secrets/   -n $(NAMESPACE)
	kubectl apply -f k8s/deployments/ -n $(NAMESPACE)
	kubectl apply -f k8s/services/  -n $(NAMESPACE)
	kubectl apply -f k8s/hpa/       -n $(NAMESPACE)
	kubectl apply -f k8s/ingress/   -n $(NAMESPACE)

.PHONY: k8s-delete
k8s-delete:
	kubectl delete namespace $(NAMESPACE)

.PHONY: k8s-status
k8s-status:
	kubectl get pods,svc,ingress -n $(NAMESPACE)

.PHONY: k8s-rollout
k8s-rollout:
	@for dep in $(SERVICES); do \
		kubectl rollout status deployment/$$dep -n $(NAMESPACE); \
	done