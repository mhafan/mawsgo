LFUNC_NAME=mawsgo-zip-demo
LFUNC_AWS_IAMROLE=

zip:
	GOOS=linux CGO_ENABLED=0 go build -o main
	zip main.zip main

create:
	aws lambda create-function --function-name $(LFUNC_NAME) --runtime go.1x --handler main --role $(LFUNC_AWS_IAMROLE)

update: zip
	aws lambda update-function-code --function-name $(LFUNC_NAME) --zip-file fileb://main.zip
		