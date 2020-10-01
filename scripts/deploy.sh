~/go/bin/build-lambda-zip.exe -output main.zip pack_linux
aws lambda update-function-code --function-name habitica-rebalanced-ep --zip-file fileb://main.zip
sleep 3s