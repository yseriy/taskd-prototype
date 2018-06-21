package amqp

import (
	"taskd/internal/transport/amqp/driver"
	"taskd/internal/transport/amqp/driver/rabbitmq"
	"time"
)

type channelManager struct {
	connector  driver.AMQPConnector
	connection driver.AMQPConnection
	channel    driver.AMQPChannel
}

func newManager(url string, timeout time.Duration) (*channelManager, error) {
	manager := &channelManager{connector: &dialer{connector: rabbitmq.NewConnector(url), timeout: timeout}}
	if err := manager.setupConnection(); err != nil {
		return nil, err
	}
	return manager, nil
}

func (m *channelManager) Channel() (driver.AMQPChannel, error) {
	if err := m.setupChannel(); err != nil {
		return nil, err
	}
	return m.channel, nil
}

func (m *channelManager) setupChannel() error {
	m.closeChannel()
	for {
		channel, err := m.connection.Channel()
		if err != nil {
			if err := m.setupConnection(); err != nil {
				return err
			} else {
				continue
			}
		}
		m.channel = channel
		return nil
	}
}

func (m *channelManager) closeChannel() error {
	if m.channel != nil {
		m.channel.Close()
	}
	return nil
}

func (m *channelManager) setupConnection() error {
	m.Close()
	connection, err := m.connector.Dial()
	if err != nil {
		return err
	}
	m.connection = connection
	return nil
}

func (m *channelManager) Close() error {
	if m.connection != nil {
		return m.connection.Close()
	}
	return nil
}
