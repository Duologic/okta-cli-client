package okta

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var SystemLogCmd = &cobra.Command{
	Use:  "systemLog",
	Long: "Manage SystemLogAPI",
}

func init() {
	rootCmd.AddCommand(SystemLogCmd)
}

// parseLogTime accepts either an RFC 3339 timestamp or a duration relative to now (e.g. 72h).
func parseLogTime(v string) (time.Time, error) {
	if d, err := time.ParseDuration(v); err == nil {
		return time.Now().Add(-d), nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time %q: use RFC 3339 or a duration like 72h", v)
	}
	return t, nil
}

func NewListLogEventsCmd() *cobra.Command {
	var (
		since, until, q, filter, sortOrder string
		opts                               listOptions
	)
	cmd := &cobra.Command{
		Use:  "listLogEvents",
		Long: "List all System Log Events",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.SystemLogAPI.ListLogEvents(apiClient.GetConfig().Context)
			if since != "" {
				t, err := parseLogTime(since)
				if err != nil {
					return err
				}
				req = req.Since(t)
			}
			if until != "" {
				t, err := parseLogTime(until)
				if err != nil {
					return err
				}
				req = req.Until(t)
			}
			if q != "" {
				req = req.Q(q)
			}
			if filter != "" {
				req = req.Filter(filter)
			}
			if sortOrder != "" {
				req = req.SortOrder(sortOrder)
			}
			if opts.after != "" {
				req = req.After(opts.after)
			}
			if opts.limit > 0 {
				req = req.Limit(opts.limit)
			}

			resp, err := req.Execute()
			return printList(cmd, httpResponse(resp), err, opts.all)
		},
	}

	cmd.Flags().StringVar(&since, "since", "720h", "Start of the window: RFC 3339 timestamp or duration ago (e.g. 72h)")
	cmd.Flags().StringVar(&until, "until", "", "End of the window: RFC 3339 timestamp or duration ago")
	cmd.Flags().StringVar(&q, "q", "", "Keyword search across event fields")
	cmd.Flags().StringVar(&filter, "filter", "", `Filter expression, e.g. 'eventType sw "application.provision.group_push"'`)
	cmd.Flags().StringVar(&sortOrder, "sortOrder", "DESCENDING", "ASCENDING or DESCENDING")
	addListFlags(cmd, &opts)

	return cmd
}

func init() {
	ListLogEventsCmd := NewListLogEventsCmd()
	SystemLogCmd.AddCommand(ListLogEventsCmd)
}
