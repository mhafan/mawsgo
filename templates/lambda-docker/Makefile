LOCREPO=lama
LAMA_NAME=mawsgo-docker-demo
AWSREPO=856366117410.dkr.ecr.eu-central-1.amazonaws.com/mawsgo-docker-demo
AWSREPO_L=$(AWSREPO):latest

build:
	docker build -t $(LOCREPO) .
	
push:
	aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin $(AWSREPO)
	docker tag $(LOCREPO) $(AWSREPO_L)
	docker push $(AWSREPO_L)
	
update:
	aws lambda update-function-code --function-name $(LAMA_NAME) --image-uri $(AWSREPO_L)