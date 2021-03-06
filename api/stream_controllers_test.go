package api_test

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ello/streams/model"
	"github.com/m4rw3r/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamController", func() {
	var id uuid.UUID

	BeforeEach(func() {
		id, _ = uuid.V4()
	})

	Context("when adding content via PUT /streams", func() {

		It("should return a status 201 when passed a correct body", func() {
			item1ID, _ := uuid.V4()
			item2ID, _ := uuid.V4()
			items := []model.StreamItem{{
				StreamID:  id.String(),
				Timestamp: time.Now(),
				Type:      0,
				ID:        item1ID.String(),
			}, {
				StreamID:  id.String(),
				Timestamp: time.Now(),
				Type:      1,
				ID:        item2ID.String(),
			}}
			itemsJSON, _ := json.Marshal(items)
			Request("PUT", "/streams", string(itemsJSON))
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))

			//verify the items passed into the service are the same
			checkAll(streamService.lastItemsOnAdd, items)
		})

		It("should return a status 201 when passed a correct body string", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				},
				{
					"id":"c8f17401-62d0-444c-a5d6-639b01f6070f",
					"ts":"2015-11-16T11:59:29.313068877-07:00",
					"type":1,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		It("should return a status 422 when passed an invalid date (non ISO8601)", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid type", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":a,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when validation error is in later element", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				},
				{
					"id":"c8f17401-62d0-444c-a5d6-639b01f6070f",
					"ts":"2015-11-16T11:59:29.313068877-07:00",
					"type":a,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid body/query", func() {
			Request("PUT", "/streams", "hi")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
	})

	Context("when removing content via DELETE /streams", func() {

		It("should return a status 200 when passed a correct body", func() {
			item1ID, _ := uuid.V4()
			item2ID, _ := uuid.V4()
			items := []model.StreamItem{{
				StreamID:  id.String(),
				Timestamp: time.Now(),
				Type:      0,
				ID:        item1ID.String(),
			}, {
				StreamID:  id.String(),
				Timestamp: time.Now(),
				Type:      1,
				ID:        item2ID.String(),
			}}
			itemsJSON, _ := json.Marshal(items)

			// Create First
			Request("PUT", "/streams", string(itemsJSON))
			Expect(response.Code).To(Equal(http.StatusCreated))

			// Now delete
			Request("DELETE", "/streams", string(itemsJSON))
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))

			//verify the items passed into the service are the same
			checkAll(streamService.lastItemsOnRemove, items)
		})
	})

	Context("when retrieving a stream via /stream/:id", func() {

		It("should return a status 201 when accessed with a valid ID", func() {
			Request("GET", "/stream/"+id.String(), "")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should use the pagination slug/limit when calling the stream service", func() {
			Request("GET", "/stream/12345?from=CBA321&limit=15", "")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(streamService.lastLimit).To(Equal(15))
			Expect(streamService.lastFromSlug).To(Equal("CBA321"))
		})
	})
	Context("when retrieving streams via /streams/coalesce", func() {

		It("should return a status 200 with a valid query string", func() {
			q := model.StreamQuery{
				Streams: []string{id.String()},
			}
			json, _ := json.Marshal(q)
			Request("POST", "/streams/coalesce", string(json))
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 200 with a valid query string", func() {
			q := `{"streams":["10e30ca7-b64d-4510-aaff-775fad0f62ed","6da0fb88-f8f5-40d3-a42c-97147a41011d"]}`
			Request("POST", "/streams/coalesce", q)
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 422 when passed an invalid query", func() {
			Request("POST", "/streams/coalesce", "")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should use the paginations slug/limit when calling the stream service", func() {
			Request("POST", "/streams/coalesce?from=CBA321&limit=15", `{"streams":["10e30ca7-b64d-4510-aaff-775fad0f62ed","6da0fb88-f8f5-40d3-a42c-97147a41011d"]}`)
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(streamService.lastLimit).To(Equal(15))
			Expect(streamService.lastFromSlug).To(Equal("CBA321"))
		})
	})
})
