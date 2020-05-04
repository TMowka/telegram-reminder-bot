FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/tmowka/telegram-reminder-bot

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

ENV TOKEN=$TOKEN
ENV CHAT_IDS=$CHAT_IDS

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the package
RUN go build -o bot ./cmd/bot

# Run executable
CMD ./bot -token ${TOKEN} -chat-id ${CHAT_IDS}
