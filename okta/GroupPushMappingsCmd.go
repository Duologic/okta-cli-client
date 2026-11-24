package okta

import (
	"net/url"
	"strconv"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

// The group push mappings API is not part of the generated SDK yet, so these
// commands call it directly.
// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/GroupPushMapping/

var GroupPushMappingsCmd = &cobra.Command{
	Use:  "groupPushMappings",
	Long: "Manage GroupPushMappingAPI (read-only)",
}

func init() {
	rootCmd.AddCommand(GroupPushMappingsCmd)
}

func NewListGroupPushMappingsCmd() *cobra.Command {
	var (
		appID         string
		sourceGroupID string
		status        string
		lastUpdated   string
		opts          listOptions
	)
	cmd := &cobra.Command{
		Use:  "lists",
		Long: "List all group push mappings for an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := url.Values{}
			if sourceGroupID != "" {
				query.Set("sourceGroupId", sourceGroupID)
			}
			if status != "" {
				query.Set("status", status)
			}
			if lastUpdated != "" {
				query.Set("lastUpdated", lastUpdated)
			}
			if opts.after != "" {
				query.Set("after", opts.after)
			}
			if opts.limit > 0 {
				query.Set("limit", strconv.Itoa(int(opts.limit)))
			}
			resp, _, err := getRaw("/api/v1/apps/"+url.PathEscape(appID)+"/group-push/mappings", query)
			return printList(cmd, resp, err, opts.all)
		},
	}

	cmd.Flags().StringVar(&appID, "appId", "", "")
	cmd.MarkFlagRequired("appId")
	cmd.Flags().StringVar(&sourceGroupID, "sourceGroupId", "", "Only return the mapping for this Okta source group")
	cmd.Flags().StringVar(&status, "status", "", "Filter by mapping status (ACTIVE, INACTIVE, ERROR)")
	cmd.Flags().StringVar(&lastUpdated, "lastUpdated", "", "Only return mappings updated after this ISO 8601 timestamp")
	addListFlags(cmd, &opts)

	return cmd
}

func init() {
	GroupPushMappingsCmd.AddCommand(NewListGroupPushMappingsCmd())
}

func NewGetGroupPushMappingCmd() *cobra.Command {
	var appID, mappingID string
	cmd := &cobra.Command{
		Use:  "get",
		Long: "Retrieve a group push mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, body, err := getRaw("/api/v1/apps/"+url.PathEscape(appID)+"/group-push/mappings/"+url.PathEscape(mappingID), nil)
			if err != nil {
				return err
			}
			return utils.PrettyPrintByte(body)
		},
	}

	cmd.Flags().StringVar(&appID, "appId", "", "")
	cmd.MarkFlagRequired("appId")
	cmd.Flags().StringVar(&mappingID, "mappingId", "", "")
	cmd.MarkFlagRequired("mappingId")

	return cmd
}

func init() {
	GroupPushMappingsCmd.AddCommand(NewGetGroupPushMappingCmd())
}
