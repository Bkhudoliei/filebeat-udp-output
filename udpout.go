package udpout

import (
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/codec"
	"github.com/elastic/beats/libbeat/publisher"
	"net"
)

func init() {
	outputs.RegisterType("udp", makeUdpout)
}

type udpOutput struct {
	connection    *net.UDPConn
	remoteAddress *net.UDPAddr
	beat          beat.Info
	observer      outputs.Observer
	codec         codec.Codec
}

// makeUdpout instantiates a new file output instance.
func makeUdpout(
	beat beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	// disable bulk support in publisher pipeline
	err := cfg.SetInt("bulk_max_size", -1, -1)
	if err != nil {
		logp.Warn("cfg.SetInt failed with: %v", err)
	}
	uo := &udpOutput{
		beat:     beat,
		observer: observer,
	}
	if err := uo.init(beat, config); err != nil {
		return outputs.Fail(err)
	}

	return outputs.Success(-1, 0, uo)
}

func (out *udpOutput) init(beat beat.Info, c udpoutConfig) error {

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	logp.Info("UDP server address: %v", address)

	var err error

	out.codec, err = codec.CreateEncoder(beat, c.Codec)
	if err != nil {
		return err
	}

	server, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		return err
	}
	out.remoteAddress = server
	out.connection = conn
	logp.Info("Initialized udp output. "+
		"Server address=%v", address)

	return nil
}

// Implement Outputer
func (out *udpOutput) Close() error {
	return out.connection.Close()
}

func (out *udpOutput) Publish(
	batch publisher.Batch,
) error {
	defer batch.ACK()

	st := out.observer
	events := batch.Events()
	st.NewBatch(len(events))

	dropped := 0
	for i := range events {
		event := &events[i]
		serializedEvent, err := out.codec.Encode(out.beat.Beat, &event.Content)
		if err != nil {
			if event.Guaranteed() {
				logp.Critical("Failed to serialize the event: %v", err)
			} else {
				logp.Warn("Failed to serialize the event: %v", err)
			}
			logp.Debug("udp", "Failed event: %v", event)

			dropped++
			continue
		}
		_, err = out.connection.Write([]byte(serializedEvent))
		if err != nil {
			st.WriteError(err)
			if event.Guaranteed() {
				logp.Critical("Writing event to UDP failed with: %v", err)
			} else {
				logp.Warn("Writing event to UDP failed with: %v", err)
			}
			dropped++
			continue
		}

		st.WriteBytes(len(serializedEvent) + 1)
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}

func (out *udpOutput) String() string {
	return "UDP(" + out.remoteAddress.String() + ")"
}
