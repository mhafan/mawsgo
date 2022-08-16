lib:
	go build

# ID vyvojarskeho uctu v AWS
AWS_ID=
# ID-nazev lambda, do ktere se tohle bude nahravat
AWS_LAMBDA=
# arn AWS role, ktera tu lambdu bude spoustet
AWS_ROLE=
AWS_REGION=eu-central-1
# ECR private registry
LOC_IMG=cosi
AWS_ECR=$(AWS_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
AWS_ECR_IMG=$(AWS_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(LOC_IMG)

push:
	git add *
	git commit -m 'Makefile push'
	git push

# Demo pro aplikace
zip:	main.go
	GOOS=linux CGO_ENABLED=0 go build -o mainzip main.go
	zip mainzip.zip mainzip

# registrace nove lambdy
# build: make zip
createlambda: zip	
	aws lambda create-function --function-name $(AWS_LAMBDA) --runtime go1.x --handle mainzip --zip-file fileb://mainzip.zip --role $(AWS_ROLE)

# update kodu
updatelambda: zip	
	aws lambda update-function-code --function-name $(AWS_LAMBDA) --zip-file fileb://mainzip.zip

# build docker image
buildimage: 
	docker build -t $(LOC_IMG) .
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ECR)
	docker tag $(LOC_IMG) $(AWS_ECR_IMG)
	docker push $(AWS_ECR_IMG)