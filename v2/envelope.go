package v2

import (
	"time"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
)

func createLogEnvelope(appID, message, sourceType, sourceInstance string, logType loggregator_v2.Log_Type) *loggregator_v2.Envelope {
	env := &loggregator_v2.Envelope{
		Timestamp:  time.Now().UnixNano(),
		SourceId:   appID,
		InstanceId: sourceInstance,
		Message: &loggregator_v2.Envelope_Log{
			Log: &loggregator_v2.Log{
				Payload: []byte(message),
				Type:    logType,
			},
		},
		Tags: map[string]*loggregator_v2.Value{
			"source_type":     newTextValue(sourceType),
			"source_instance": newTextValue(sourceInstance),
		},
	}
	return env
}

func newTextValue(t string) *loggregator_v2.Value {
	return &loggregator_v2.Value{Data: &loggregator_v2.Value_Text{Text: t}}
}

func newGaugeValue(f float64) *loggregator_v2.GaugeValue {
	return &loggregator_v2.GaugeValue{Value: f}
}

func newGaugeValueFromUInt64(i uint64) *loggregator_v2.GaugeValue {
	return newGaugeValue(float64(i))
}
