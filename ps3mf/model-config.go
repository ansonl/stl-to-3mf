package ps3mf

import (
	"encoding/xml"
	"strconv"
)

type IdPair struct {
	FirstId int
	LastId int
}

type ModelConfig struct {
	XMLName xml.Name `xml:"config"`
	Objects []ModelConfigObject `xml:"object"`
}

type ModelConfigObject struct {
	XMLName xml.Name `xml:"object"`
	Id string `xml:"id,attr"`
	InstancesCount string `xml:"instances_count,attr"`
	Metadata []ModelConfigMeta `xml:"metadata"`
	Volume []ModelConfigVolume `xml:"volume"`
}

type ModelConfigVolume struct {
	XMLName xml.Name `xml:"volume"`
	FirstId string `xml:"firstid,attr"`
	LastId string `xml:"lastid,attr"`
	Metadata []ModelConfigMeta `xml:"metadata"`
}

type ModelConfigMeta struct {
	XMLName xml.Name `xml:"metadata"`
	Type string `xml:"type,attr"`
	Key string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

func GetModelConfigMeta(typ, key, value string) ModelConfigMeta {
	return ModelConfigMeta{
		Type: typ,
		Key: key,
		Value: value,
	}
}

func boolToIntString(b bool) string {
	if b { return "1" }
	return "0"
}

func (b *Bundle) GetModelConfig(m *ModelXML, idPairs []IdPair) ModelConfig {
	config := ModelConfig{
		Objects: make([]ModelConfigObject, 0, len(b.Model.Resources.Objects)),
	}
	for idx := range m.Resources {
		id := m.Resources[idx].Id
		objectConfig := ModelConfigObject{
			Id: id,
			InstancesCount: "1",
			Metadata: []ModelConfigMeta{
				GetModelConfigMeta("object", "name", "model"),
				GetModelConfigMeta("object", "extruder", b.Extruders[idx]),
				GetModelConfigMeta("object", "wipe_into_infill", boolToIntString(b.WipeIntoInfill[idx])),
				GetModelConfigMeta("object", "wipe_into_objects", boolToIntString(b.WipeIntoModel[idx])),
			},
			Volume: make([]ModelConfigVolume, len(idPairs)),
		}
		for volumeIndex, idPair := range idPairs {
			objectConfig.Volume[volumeIndex] = ModelConfigVolume{
				XMLName:  xml.Name{},
				FirstId:  strconv.Itoa(idPair.FirstId),
				LastId:   strconv.Itoa(idPair.LastId),
				Metadata: []ModelConfigMeta{
					GetModelConfigMeta("volume", "name", b.Names[volumeIndex]),
					GetModelConfigMeta("volume", "volume_type", "ModelPart"),
					// use identity matrix since vertices are already transformed
					GetModelConfigMeta("volume", "matrix", "1 0 0 0 0 1 0 0 0 0 1 0 0 0 0 1"),
					GetModelConfigMeta("volume", "source_object_id", strconv.Itoa(volumeIndex)),
					GetModelConfigMeta("volume", "source_volume_id", "0"),
				},
			}
		}
		config.Objects = append(config.Objects, objectConfig)
	}
	return config
}
