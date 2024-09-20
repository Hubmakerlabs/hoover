package ao

import "github.com/Hubmakerlabs/hoover/pkg/arweave/goar/types"

const Gateway = "https://arweave.net"

func AppendTag(bundle *types.BundleItem, name, value string) {
	bundle.Tags = append(bundle.Tags, types.Tag{Name: name, Value: value})
}

type EventData struct {
	Content   string     `json:"Content,omitempty"`
	EventTags [][]string `json:"Tags,omitempty"`
}

func NewEventData(content string) *EventData {
	return &EventData{Content: content}
}

func (evData *EventData) Append(values ...string) {
	evData.EventTags = append(evData.EventTags, values)
}
