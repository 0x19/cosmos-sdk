package codec

import (
	"encoding/json"

	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// MarshalYAML marshals toPrint using JSONCodec to leverage specialized MarshalJSON methods
// (usually related to serialize data with protobuf or amin depending on a configuration).
// This involves additional roundtrip through JSON.
func MarshalYAML(cdc JSONCodec, toPrint proto.Message) ([]byte, error) {
	// We are OK with the performance hit of the additional JSON roundtip. MarshalYAML is not
	// used in any critical parts of the system.
	// TODO temporarily disabled yaml to reduce dependencies

	bz, err := cdc.MarshalJSON(toPrint)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := json.Unmarshal(bz, &data); err != nil {
		return nil, err
	}
	return yaml.Marshal(data)
}
