package service

import "github.com/ello/streams/model"

// StreamService is the interface to the underlying stream storage system.
type StreamService interface {

	//Add will add the content items to the embedded stream id
	Add(items []model.StreamItem) error

	//Remove will remove the content items from the embedded stream id
	Remove(items []model.StreamItem) error

	//Load will pull a coalesced view of the streams in the query
	Load(query model.StreamQuery, limit int, fromSlug string) (*model.StreamQueryResponse, error)
}
