include .env
export

bot:
	go run cmd/bot/main.go \
		-token ${TOKEN} \
		-chat-id ${CHAT_ID}
