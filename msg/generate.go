package msg
//go:generate mockery --quiet --name ".*(Message|MessageReceiver|Producer|Consumer|Event|Reply|MessageSubscriber|CommandMessagePublisher|EntityEventMessagePublisher|EventMessagePublisher|ReplyMessagePublisher|MessagePublisher)$" --dir . --output ./msgmocks/ --case underscore --outpkg=msgmocks
