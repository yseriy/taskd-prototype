package connector

import (
	"taskd/internal/taskd"
	"taskd/internal/transport"
)

type connector struct {
	configConverter configConverter
}

func (connector *connector) Input() transport.InputConnector {
	repeater := connector.repeater()
	return newInputConnector(
		repeater,
		connector.configConverter.mainQueue(),
		connector.configConverter.consumeParams(),
	)
}

func (connector *connector) Output() transport.OutputConnector {
	repeater := connector.repeater()
	return newOutputConnector(
		repeater,
		connector.configConverter.destinationQueueTemplate(),
		connector.configConverter.publishingParams(),
	)
}

func (connector *connector) repeater() dialer {
	dialer := newDialer(connector.configConverter.url())
	return newRepeater(dialer, connector.configConverter.timeout())
}

func New(rabbitConfig *taskd.RabbitConfig) transport.Connector {
	return &connector{configConverter{inputConfig: rabbitConfig}}
}

type configConverter struct {
	inputConfig *taskd.RabbitConfig
}

func (converter *configConverter) url() string {
	return converter.inputConfig.Url
}

func (converter *configConverter) timeout() int {
	return converter.inputConfig.Timeout
}

func (converter *configConverter) mainQueue() mainQueue {
	return mainQueue{
		Name:       converter.inputConfig.MainQueue.Name,
		Durable:    converter.inputConfig.MainQueue.Durable,
		AutoDelete: converter.inputConfig.MainQueue.AutoDelete,
		Exclusive:  converter.inputConfig.MainQueue.Exclusive,
		NoWait:     converter.inputConfig.MainQueue.NoWait,
		Args:       converter.inputConfig.MainQueue.Args,
	}
}

func (converter *configConverter) consumeParams() consumeParams {
	return consumeParams{
		Name:      converter.inputConfig.ConsumeParams.Name,
		AutoAck:   converter.inputConfig.ConsumeParams.AutoAck,
		Exclusive: converter.inputConfig.ConsumeParams.Exclusive,
		NoLocal:   converter.inputConfig.ConsumeParams.NoLocal,
		NoWait:    converter.inputConfig.ConsumeParams.NoWait,
		Args:      converter.inputConfig.ConsumeParams.Args,
	}
}

func (converter *configConverter) destinationQueueTemplate() destinationQueueTemplate {
	return destinationQueueTemplate{
		Durable:    converter.inputConfig.DestinationQueueTemplate.Durable,
		AutoDelete: converter.inputConfig.DestinationQueueTemplate.AutoDelete,
		Exclusive:  converter.inputConfig.DestinationQueueTemplate.Exclusive,
		NoWait:     converter.inputConfig.DestinationQueueTemplate.NoWait,
		Args:       converter.inputConfig.DestinationQueueTemplate.Args,
	}
}

func (converter *configConverter) publishingParams() publishingParams {
	return publishingParams{
		Mandatory: converter.inputConfig.PublishingParams.Mandatory,
		Immediate: converter.inputConfig.PublishingParams.Immediate,
	}
}
