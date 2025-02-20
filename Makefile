SERVICES = \
    user-service \
    order-service \
    shared-service

run:
	@echo "Services are running..."
	@bash -c ' \
		for svc in $(SERVICES); do \
			(cd $$svc && go run .) & \
		done; \
		trap "echo Stopping services...; kill $$(jobs -p)" INT; \
		wait'