define copy_template
	@if [ -f $(2) ]; then \
		echo "\"$(2)\" already exists, skipping"; \
	else \
		cp $(1) $(2) && echo "\"$(2)\" created"; \
	fi
endef

setup-config:
	@echo "Setting up config files from templates..."

	$(call copy_template,docs/setup-examples/env-example,.env)

	@mkdir -p config
	$(call copy_template,docs/setup-examples/config-exmaple.yaml,config/config.yaml)