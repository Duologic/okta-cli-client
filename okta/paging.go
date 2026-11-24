package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/okta/okta-cli-client/sdk"
	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

// listOptions holds the common pagination flags for list commands.
type listOptions struct {
	after string
	limit int32
	all   bool
}

func addListFlags(cmd *cobra.Command, opts *listOptions) {
	cmd.Flags().StringVar(&opts.after, "after", "", "Pagination cursor returned by a previous request")
	cmd.Flags().Int32Var(&opts.limit, "limit", 0, "Number of results per page (API default when 0)")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Follow pagination links and return all results")
}

// printList prints the response of a list request. With all set, it follows
// the Link rel="next" headers and prints the merged JSON array. Without it,
// the cursor for the next page (if any) is written to stderr.
func printList(cmd *cobra.Command, resp *http.Response, reqErr error, all bool) error {
	if reqErr != nil {
		if resp != nil && resp.Body != nil {
			if d, err := io.ReadAll(resp.Body); err == nil {
				utils.PrettyPrintByte(d)
			}
		}
		return reqErr
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	if !all {
		if err := utils.PrettyPrintByte(body); err != nil {
			return err
		}
		if next := nextLink(resp); next != "" {
			cmd.PrintErrf("next page: --after %s\n", cursorFromLink(next))
		}
		return nil
	}

	var items []json.RawMessage
	seen := map[string]bool{}
	for {
		var page []json.RawMessage
		if err := json.Unmarshal(body, &page); err != nil {
			return fmt.Errorf("decoding page: %w", err)
		}
		items = append(items, page...)

		// The System Log always returns a next link for polling; an empty
		// page means there is nothing more to fetch.
		next := nextLink(resp)
		if next == "" || len(page) == 0 {
			break
		}
		if seen[next] {
			return fmt.Errorf("repeated pagination link: %s", next)
		}
		seen[next] = true

		resp, body, err = getRaw(next, nil)
		if err != nil {
			return err
		}
	}

	if items == nil {
		items = []json.RawMessage{}
	}
	d, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return utils.PrettyPrintByte(d)
}

// getRaw performs an authenticated GET against an Okta API path (or absolute
// URL), for endpoints that the generated SDK does not cover.
func getRaw(path string, query url.Values) (*http.Response, []byte, error) {
	ctx := apiClient.GetConfig().Context
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		base, err := apiClient.GetConfig().ServerURLWithContext(ctx, "ApplicationAPIService.ListApplications")
		if err != nil {
			return nil, nil, err
		}
		path = base + path
	}
	if query == nil {
		query = url.Values{}
	}
	headers := map[string]string{"Accept": "application/json"}
	req, err := apiClient.PrepareRequest(ctx, path, http.MethodGet, nil, headers, query, url.Values{}, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := apiClient.Do(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 300 {
		utils.PrettyPrintByte(body)
		return nil, nil, fmt.Errorf("%s", resp.Status)
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return resp, body, nil
}

func httpResponse(resp *sdk.APIResponse) *http.Response {
	if resp == nil {
		return nil
	}
	return resp.Response
}

func nextLink(resp *http.Response) string {
	for _, header := range resp.Header.Values("Link") {
		for _, link := range strings.Split(header, ",") {
			if !strings.Contains(link, `rel="next"`) {
				continue
			}
			start, end := strings.Index(link, "<"), strings.Index(link, ">")
			if start >= 0 && end > start {
				return link[start+1 : end]
			}
		}
	}
	return ""
}

func cursorFromLink(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	return u.Query().Get("after")
}
