package msg
//go:generate mockery --quiet --name ".*(Message|Producer|Consumer|Event|Reply|MessageSubscriber|CommandMessagePublisher|EntityEventMessagePublisher|EventMessagePublisher|ReplyMessagePublisher|MessagePublisher)$" --dir . --output ./msgmocks/ --case underscore --outpkg=msgmocks
