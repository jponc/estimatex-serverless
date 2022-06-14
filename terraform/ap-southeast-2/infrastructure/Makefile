# Create workspaces per env
init-workspaces:
	terraform workspace new dev
	terraform workspace new staging
	terraform workspace new prod

# Initialise per env

init-dev:
	terraform workspace select dev
	terraform init

init-staging:
	terraform workspace select staging
	terraform init

init-prod:
	terraform workspace select prod
	terraform init

# Plan per env
plan-dev:
	terraform workspace select dev
	terraform plan -var-file="dev.tfvars"

plan-staging:
	terraform workspace select staging
	terraform plan -var-file="staging.tfvars"

plan-prod:
	terraform workspace select prod
	terraform plan -var-file="prod.tfvars"

# Apply per env
apply-dev:
	terraform workspace select dev
	terraform apply -var-file="dev.tfvars"

apply-staging:
	terraform workspace select staging
	terraform apply -var-file="staging.tfvars"

apply-prod:
	terraform workspace select prod
	terraform apply -var-file="prod.tfvars"
