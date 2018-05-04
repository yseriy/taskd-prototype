package connector

import "taskd/internal/transport"

type outputConnector struct {
	dialer           dialer
	queueTemplate    destinationQueueTemplate
	publishingParams publishingParams
}

func newOutputConnector(dialer dialer, queueTemplate destinationQueueTemplate,
	publishingParams publishingParams) transport.OutputConnector {
	return &outputConnector{dialer: dialer, queueTemplate: queueTemplate, publishingParams: publishingParams}
}

func (outputConnector *outputConnector) Connect() (transport.SendHelper, error) {
	connection, err := outputConnector.dialer.dial()
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return newSendHelper(channel, outputConnector.queueTemplate, outputConnector.publishingParams), nil
}
