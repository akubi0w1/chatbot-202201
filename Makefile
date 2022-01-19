GCP_PROJECT=
GCP_REGION=asia-northeast1

UNAME_S=$(shell uname -s)

###############################
# run test
# usege: go-test
go-test:
	go run -v --cover ./...

###############################
# build
# usege: go-build PLATFORM=gcp
go-build:
ifeq ($(UNAME_S),Linux)
	GOOS=linux GOARCH=amd64 go build -o bin/$(PLATFORM)/app cmd/$(PLATFORM)/main.go
endif
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(PLATFORM)/app cmd/$(PLATFORM)/main.go
endif

##############################
# push-image
# usage: push-image PLATFORM=gcp
push-image:
ifeq ($(PLATFORM),gcp)
	GOOS=linux GOARCH=amd64 go build -o ./bin/gcp/chatbot cmd/gcp/main.go
	gcloud builds submit --config build/gcp/cloudbuild.yaml . --substitutions TAG_NAME=$(TAG)
endif

##############################
# deploy to api
# usage: deploy TAG=0.0.1 PLATFORM=gcp
deploy:
	GOOS=linux GOARCH=amd64 go build -o bin/$(PLATFORM)/app cmd/$(PLATFORM)/main.go
ifeq ($(PLATFORM),gcp)
	gcloud run deploy akubi-post-chatbot \
	--image=asia.gcr.io/$(GCP_PROJECT)/chatbot:$(TAG) \
	--platform=managed \
	--region=$(GCP_REGION) \
	--project=$(GCP_PROJECT)
endif
ifeq ($(PLATFORM),aws)
	rm -rf ./functions/$(PLATFORM)
	mkdir -p ./functions/$(PLATFORM)
	zip -j ./functions/$(PLATFORM)/chatbot.zip ./bin/$(PLATFORM)/chatbot
	aws lambda update-function-code \
		--profile sandbox \$
		--function-name akubiPostChatbot \
		--zip-file fileb://functions/$(PLATFORM)/chatbot.zip
endif

################################
# terraform init
# usage: terraform-init PLATFORM=gcp SERVICE=common
terraform-init:
ifeq ($(SERVICE),common)
	cd terraform/$(PLATFORM)/$(SERVICE) && terraform init
else
	cd terraform/$(PLATFORM)/$(SERVICE) && terraform init
endif

################################
# terraform plan
# usage: terraform-plan PLATFORM=gcp PROFILE=production
terraform-plan:
	cd terraform/$(PLATFORM) && terraform workspace select $(PROFILE)
	@echo ----------------------------------------------
	@echo Execute workspace: `cd terraform/$(PLATFORM) && terraform workspace show`
	cd terraform/$(PLATFORM) && terraform plan -var-file=vars/$(PROFILE).tfvars

################################
# terraform apply
# usage: terraform-apply PLATFORM=gcp PROFILE=production
terraform-apply:
	cd terraform/$(PLATFORM) && terraform workspace select $(PROFILE)
	@echo ----------------------------------------------
	@echo Execute workspace: `cd terraform/$(PLATFORM) && terraform workspace show`
	cd terraform/$(PLATFORM) && terraform apply -var-file=vars/$(PROFILE).tfvars

################################
# terraform destroy
# usage: terraform-destroy PLATFORM=gcp PROFILE=production
terraform-destroy:
	cd terraform/$(PLATFORM) && terraform workspace select $(PROFILE)
	@echo ----------------------------------------------
	@echo Execute workspace: `cd terraform/$(PLATFORM) && terraform workspace show`
	cd terraform/$(PLATFORM) && terraform destroy -var-file=vars/$(PROFILE).tfvars

################################
# terraform workspace list
# usave: terraform-workspace PLATFORM=gcp
terraform-workspace:
	cd terraform/$(PLATFORM) && terraform workspace list

################################
# terraform workspace new
# usage: terraform-workspace-new PLATFORM=gcp PROFILE=production
terraform-workspace-new:
	cd terraform/$(PLATFORM) && terraform workspace new $(PROFILE)