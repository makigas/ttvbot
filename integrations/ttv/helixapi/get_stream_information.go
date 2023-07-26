package helixapi

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

type GetStreamInformationRequest struct {
	Id string
}

type GetStreamInformationResponse struct {
	Title    string
	Username string
}

func (ha *HelixApi) GetStreamInformation(req *GetStreamInformationRequest) (*GetStreamInformationResponse, error) {
	return ha.getStreamInformationRetryable(req.Id, true)
}

func (ha *HelixApi) getStreamInformationRetryable(id string, retry bool) (*GetStreamInformationResponse, error) {
	resp, err := ha.client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{id},
	})
	if err != nil {
		return nil, err
	}

	// Handle status codes.
	if resp.StatusCode == 401 {
		if err := ha.refreshApplicationToken(); err != nil {
			return nil, err
		}
		if retry {
			return ha.getStreamInformationRetryable(id, false)
		}
		return nil, ErrAuthentication
	}
	if resp.StatusCode != 200 {
		return nil, ErrUnexpectedResponse
	}
	if len(resp.Data.Streams) != 1 {
		return nil, errors.New("is not live")
	}

	// Get the stream information.
	return &GetStreamInformationResponse{
		Title:    resp.Data.Streams[0].Title,
		Username: resp.Data.Streams[0].UserLogin,
	}, nil
}
