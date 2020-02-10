# This Makefile describes the operations for working with NODE in terraform context.
# https://learn.hashicorp.com/terraform/development/running-terraform-in-automation


# Get the terraform bin path.
TERRAFORM := $(shell which terraform)
ifeq ($(TERRAFORM),)
	# Case when bainary is not installed. We have target below for installing stable version.
	TERRAFORM := $(HOME)/bin/terraform
	TERRAFORM_VERSION := 0.12.15
	TERRAFORM_RELEASE_LINK := https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_linux_amd64.zip
endif


# Set golang global env variables.
# https://www.terraform.io/docs/commands/environment-variables.html
export TF_INPUT=0
export TF_IN_AUTOMATION=1


# Find source files.
TERRAFORM_SOURCE_FILES := $(wildcard *.tf)


$(TERRAFORM):
	#### Node( '$(NODE)' ).Call( '$@' )
	mkdir -p $(HOME)/tmp $(HOME)/bin
	wget $(TERRAFORM_RELEASE_LINK) -qO $(HOME)/tmp/terraform.zip
	apt-get -qq update
	apt-get -qq install -y unzip
	unzip $(HOME)/tmp/terraform.zip -d $(HOME)/bin
	rm -rf $(HOME)/tmp/terraform.zip
	$(TERRAFORM) --version


# Target to initialize the terraform working directory.
.terraform/tflock: $(TERRAFORM) $(TERRAFORM_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(TERRAFORM) init
	@touch .terraform/tflock


# Part of the 'lint' global Makefile interface.
.PHONY: terraform-fmt
terraform-fmt: $(TERRAFORM) $(TERRAFORM_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(TERRAFORM) fmt


# Target to create a plan and save it to the local file.
tfplan: $(TERRAFORM) .terraform/tflock $(TERRAFORM_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(TERRAFORM) validate
	$(TERRAFORM) plan -out=tfplan


# Target to create a plan for destroying and save it to the local file.
tfplan_destroy: $(TERRAFORM) .terraform/tflock $(TERRAFORM_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(TERRAFORM) validate
	$(TERRAFORM) plan -out=tfplan_destroy -destroy


# Target to apply the plan stored in the local file (approved plan!).
# Part of the 'publish' global Makefile interface.
.PHONY: terraform-apply
terraform-apply: $(TERRAFORM) tfplan
	#### Node( '$(NODE)' ).Call( '$@' )
	$(TERRAFORM) apply tfplan
	@rm -f tfplan # Plan is stale now, need to rebuild.


# Target to destroy the plan stored in the local file (approved plan!).
.PHONY: terraform-destroy
terraform-destroy: $(TERRAFORM) tfplan_destroy
	#### Node( '$(NODE)' ).Call( '$@' )
	$(TERRAFORM) apply tfplan_destroy
	@rm -f tfplan_destroy # Plan is stale now, need to rebuild.


# terraform-publish is the main entrypoint for CI/CD scripts.
# Stage 1 - generate plan (automatic step).
# Stage 2 - apply generated plan (preferred to be a manual step).
# If plan file exists - publish have dependency to just apply them.
# Otherwise dependency is the target for plan generation.
.PHONY: terraform-publish
ifeq ($(wildcard tfplan),) # Check file exists.
terraform-publish: tfplan
else
terraform-publish: terraform-apply
endif
