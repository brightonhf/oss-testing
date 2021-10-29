it:
	echo hello
# GOOS=linux GOARCH=amd64 go build -o app .
# kubectl cp ./app scm/chatbot-service-api-6bdfff7d5d-7cbmq:/
# kubectl cp ./app scm/owl-search-service-rpc-7db959db55-lhl97:/