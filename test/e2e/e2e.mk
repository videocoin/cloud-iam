ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: e2e
e2e:
	$(call log_info,Starting test environment:)
	cd $(ROOT_DIR) && docker-compose up -d --build mysql
	sleep 30
	cd $(ROOT_DIR) && docker-compose up -d --build migrate iamd
	$(call log_info,Test environment started!)
	make e2e-venom

.PHONY: e2e-nobuild
e2e-nobuild:
	$(call log_info,Starting test environment:)
	cd $(ROOT_DIR) && docker-compose up -d --no-build mysql 
	sleep 30
	cd $(ROOT_DIR) && docker-compose up -d --no-build migrate iamd
	$(call log_info,Test environment started!)
	make e2e-venom

.PHONY: e2e-venom
e2e-venom:
	sh $(ROOT_DIR)/e2e-venom.sh

.PHONY: e2e-stop
e2e-stop:
	cd $(ROOT_DIR) && docker-compose down --volumes