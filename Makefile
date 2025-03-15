.PHONY: frontend backend all

frontend:
	cd frontend && npm ci && npm run deploy

backend:
	cd backend && make deploy

all:
	@$(MAKE) frontend & $(MAKE) backend & wait