NAME=slack-bot
FUNCTION_NAME=bot
VERSION=latest
DATE=`date +"%Y%m%d_%H%M%S"`
TEST_JSON='{"path": "one"}'

clean:
	rm -rf dist

updateLambda: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/main main.go
	cd dist && zip main.zip main
	aws lambda update-function-code --function-name ${FUNCTION_NAME} --zip-file fileb://${pwd}dist/main.zip

invoke:
	aws lambda invoke \
		--function-name "${FUNCTION_NAME}" \
		--log-type "Tail" \
		--payload $(TEST_JSON) \
		output/$(DATE).log \
		| jq -r '.LogResult' | base64 -D
