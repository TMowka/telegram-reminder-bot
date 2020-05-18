FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/tmowka/telegram-reminder-bot

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

ENV BOT_BOT_TOKEN=$BOT_BOT_TOKEN
ENV BOT_CHAT_ID=$BOT_CHAT_ID
ENV BOT_BOT_LOCATION=$BOT_BOT_LOCATION

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the package
RUN go build -o bot ./cmd/bot

# Run executable
CMD ./bot \
    -token ${BOT_BOT_TOKEN} \
    -chat-ids ${BOT_CHAT_ID} \
    -location ${BOT_BOT_LOCATION}
