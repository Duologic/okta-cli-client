package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var ApplicationGroupsCmd = &cobra.Command{
	Use:  "applicationGroups",
	Long: "Manage ApplicationGroupsAPI",
}

func init() {
	rootCmd.AddCommand(ApplicationGroupsCmd)
}

var ListApplicationGroupAssignmentsappId string

func NewListApplicationGroupAssignmentsCmd() *cobra.Command {
	var (
		q, expand string
		opts      listOptions
	)
	cmd := &cobra.Command{
		Use:  "listApplicationGroupAssignments",
		Long: "List all Assigned Groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.ListApplicationGroupAssignments(apiClient.GetConfig().Context, ListApplicationGroupAssignmentsappId)
			if q != "" {
				req = req.Q(q)
			}
			if expand != "" {
				req = req.Expand(expand)
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

	cmd.Flags().StringVarP(&ListApplicationGroupAssignmentsappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")
	cmd.Flags().StringVar(&q, "q", "", "Search assigned groups by name prefix")
	cmd.Flags().StringVar(&expand, "expand", "", "Expand embedded resources (group, metadata)")
	addListFlags(cmd, &opts)

	return cmd
}

func init() {
	ListApplicationGroupAssignmentsCmd := NewListApplicationGroupAssignmentsCmd()
	ApplicationGroupsCmd.AddCommand(ListApplicationGroupAssignmentsCmd)
}

var (
	GetApplicationGroupAssignmentappId string

	GetApplicationGroupAssignmentgroupId string
)

func NewGetApplicationGroupAssignmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "getApplicationGroupAssignment",
		Long: "Retrieve an Assigned Group",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.GetApplicationGroupAssignment(apiClient.GetConfig().Context, GetApplicationGroupAssignmentappId, GetApplicationGroupAssignmentgroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&GetApplicationGroupAssignmentappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&GetApplicationGroupAssignmentgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	GetApplicationGroupAssignmentCmd := NewGetApplicationGroupAssignmentCmd()
	ApplicationGroupsCmd.AddCommand(GetApplicationGroupAssignmentCmd)
}

var (
	AssignGroupToApplicationappId string

	AssignGroupToApplicationgroupId string

	AssignGroupToApplicationdata string
)

func NewAssignGroupToApplicationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "assignGroupToApplication",
		Long: "Assign a Group",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.AssignGroupToApplication(apiClient.GetConfig().Context, AssignGroupToApplicationappId, AssignGroupToApplicationgroupId)

			if AssignGroupToApplicationdata != "" {
				req = req.Data(AssignGroupToApplicationdata)
			}

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&AssignGroupToApplicationappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&AssignGroupToApplicationgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	cmd.Flags().StringVarP(&AssignGroupToApplicationdata, "data", "", "", "")
	cmd.MarkFlagRequired("data")

	return cmd
}

func init() {
	AssignGroupToApplicationCmd := NewAssignGroupToApplicationCmd()
	ApplicationGroupsCmd.AddCommand(AssignGroupToApplicationCmd)
}

var (
	UnassignApplicationFromGroupappId string

	UnassignApplicationFromGroupgroupId string
)

func NewUnassignApplicationFromGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unassignApplicationFromGroup",
		Long: "Unassign a Group",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := apiClient.ApplicationGroupsAPI.UnassignApplicationFromGroup(apiClient.GetConfig().Context, UnassignApplicationFromGroupappId, UnassignApplicationFromGroupgroupId)

			resp, err := req.Execute()
			if err != nil {
				if resp != nil && resp.Body != nil {
					d, err := io.ReadAll(resp.Body)
					if err == nil {
						utils.PrettyPrintByte(d)
					}
				}
				return err
			}
			d, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			utils.PrettyPrintByte(d)
			// cmd.Println(string(d))
			return nil
		},
	}

	cmd.Flags().StringVarP(&UnassignApplicationFromGroupappId, "appId", "", "", "")
	cmd.MarkFlagRequired("appId")

	cmd.Flags().StringVarP(&UnassignApplicationFromGroupgroupId, "groupId", "", "", "")
	cmd.MarkFlagRequired("groupId")

	return cmd
}

func init() {
	UnassignApplicationFromGroupCmd := NewUnassignApplicationFromGroupCmd()
	ApplicationGroupsCmd.AddCommand(UnassignApplicationFromGroupCmd)
}
