.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./hello-world/hello-world
	
build:
	GOOS=linux GOARCH=amd64 go build -o hello-world/bin/hello-world ./hello-world

test:
	go test -v ./...

local-api:
	sam local start-api --host 0.0.0.0 --env-vars ./localenv.json

deploy:
	@sam package -t ./template.yaml --s3-bucket lambda-meraki-oh.app.ducbase.com --output-template-file deploy.yaml
	@sam deploy --template-file ./deploy.yaml --stack-name hello --capabilities CAPABILITY_IAM	