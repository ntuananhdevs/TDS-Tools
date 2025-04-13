package facebook

import (
	"fmt"
	"net/http"
)

type FacebookClient struct {
	Cookie string
}

func NewFacebookClient(cookie string) *FacebookClient {
	return &FacebookClient{Cookie: cookie}
}

func (fb *FacebookClient) Like(postID string) error {
	url := fmt.Sprintf("https://mbasic.facebook.com/ufi/reaction/?ft_ent_identifier=%s", postID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", fb.Cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Like failed")
	}
	return nil
}
