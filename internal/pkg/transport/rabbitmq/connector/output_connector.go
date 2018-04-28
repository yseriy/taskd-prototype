package connector

type OutputConnector struct {
	connector        connector
	queueTemplate    DestinationQueueTemplate
	publishingParams PublishingParams
}

func newOutputConnector(connector connector, queueTemplate DestinationQueueTemplate,
	publishingParams PublishingParams) *OutputConnector {
	return &OutputConnector{connector: connector, queueTemplate: queueTemplate, publishingParams: publishingParams}
}

func (outputConnector *OutputConnector) Connect() *PublisherHelper {
	channel := outputConnector.connector.connect()
	return newPublisherHelper(channel, outputConnector.queueTemplate, outputConnector.publishingParams)
}
