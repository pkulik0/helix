package helix

import (
	"net/http"
	"testing"
)

func TestGetVideos(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode   int
		options      *Options
		VideosParams *VideosParams
		respBody     string
	}{
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&VideosParams{IDs: []string{"225470646"}, GameID: "21779"},
			`{"error":"Bad Request","status":400,"message":"Must provide only one of the following query params: user ID, game ID,or one or more video IDs."}`,
		},
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&VideosParams{GameID: "21779", Period: "month", Type: "highlight", Language: "en", Sort: "views", First: 1},
			`{"data":[{"id":"224404190","user_id":"30080751","user_name":"Dyrus","title":"jhin mains LUL","description":"LUL","created_at":"2018-01-31T22:35:55Z","published_at":"2018-01-31T22:35:55Z","url":"https://www.twitch.tv/videos/224404190","thumbnail_url":"https://static-cdn.jtvnw.net/s3_vods/427483724d153cb8c673_dyrus_27413838016_782035045//thumb/thumb224404190-%{width}x%{height}.jpg","viewable":"public","view_count":4924,"language":"en","type":"highlight","duration":"50s"}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7Ik9mZnNldCI6MX19"}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetVideos(testCase.VideosParams)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be %d, got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusBadRequest {
			if resp.Error != "Bad Request" {
				t.Errorf("expected error to be %s, got %s", "Bad Request", resp.Error)
			}

			if resp.ErrorStatus != http.StatusBadRequest {
				t.Errorf("expected error status to be %d, got %d", http.StatusBadRequest, resp.ErrorStatus)
			}

			expectedErrMsg := "Must provide only one of the following query params: user ID, game ID,or one or more video IDs."
			if resp.ErrorMessage != expectedErrMsg {
				t.Errorf("expected error message to be %s, got %s", expectedErrMsg, resp.ErrorMessage)
			}

			continue
		}

		if len(resp.Data.Videos) != 1 {
			t.Errorf("expected single video, got %d results", len(resp.Data.Videos))
		}

		if resp.Data.Videos[0].Type != testCase.VideosParams.Type {
			t.Errorf("expected video type to be %s, got %s", testCase.VideosParams.Type, resp.Data.Videos[0].Type)
		}

		if resp.Data.Videos[0].Language != testCase.VideosParams.Language {
			t.Errorf("expected video language to be %s, got %s", testCase.VideosParams.Language, resp.Data.Videos[0].Language)
		}
	}
}
