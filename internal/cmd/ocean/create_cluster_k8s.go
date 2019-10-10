package ocean

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spotinst/spotinst-cli/internal/spotinst"
	"github.com/spotinst/spotinst-cli/internal/utils/flags"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	spotinstsdk "github.com/spotinst/spotinst-sdk-go/spotinst"
)

type (
	CmdCreateClusterKubernetes struct {
		cmd  *cobra.Command
		opts CmdCreateClusterKubernetesOptions
	}

	CmdCreateClusterKubernetesOptions struct {
		*CmdCreateClusterOptions

		// Base
		Name   string
		Region string

		// Strategy
		SpotPercentage           float64
		UtilizeReservedInstances bool
		FallbackToOnDemand       bool

		// Capacity
		MinSize    int
		MaxSize    int
		TargetSize int

		// Compute
		SubnetIDs                []string
		InstanceTypesWhitelist   []string
		InstanceTypesBlacklist   []string
		SecurityGroupIDs         []string
		ImageID                  string
		KeyPair                  string
		UserData                 string
		RootVolumeSize           int
		AssociatePublicIPAddress bool
		EnableMonitoring         bool
		EnableEBSOptimization    bool

		// Auto Scaling
		EnableAutoScaler bool
		EnableAutoConfig bool
		Cooldown         int
	}
)

func NewCmdCreateClusterKubernetes(opts *CmdCreateClusterOptions) *cobra.Command {
	return newCmdCreateClusterKubernetes(opts).cmd
}

func newCmdCreateClusterKubernetes(opts *CmdCreateClusterOptions) *CmdCreateClusterKubernetes {
	var cmd CmdCreateClusterKubernetes

	cmd.cmd = &cobra.Command{
		Use:           "kubernetes",
		Short:         "Create a new Kubernetescluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(*cobra.Command, []string) error {
			return cmd.Run(context.Background())
		},
	}

	cmd.opts.Init(cmd.cmd.Flags(), opts)

	return &cmd
}

func (x *CmdCreateClusterKubernetes) Run(ctx context.Context) error {
	steps := []func(context.Context) error{
		x.survey,
		x.log,
		x.validate,
		x.run,
	}

	for _, step := range steps {
		if err := step(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (x *CmdCreateClusterKubernetes) survey(ctx context.Context) error {
	if x.opts.Noninteractive {
		return nil
	}

	return nil
}

func (x *CmdCreateClusterKubernetes) log(ctx context.Context) error {
	flags.Log(x.cmd)
	return nil
}

func (x *CmdCreateClusterKubernetes) validate(ctx context.Context) error {
	return x.opts.Validate()
}

func (x *CmdCreateClusterKubernetes) run(ctx context.Context) error {
	spotinstClientOpts := []spotinst.ClientOption{
		spotinst.WithCredentialsProfile(x.opts.Profile),
	}

	spotinstClient, err := x.opts.Clients.NewSpotinst(spotinstClientOpts...)
	if err != nil {
		return err
	}

	oceanClient, err := spotinstClient.Services().Ocean(x.opts.CloudProvider, spotinst.OrchestratorKubernetes)
	if err != nil {
		return err
	}

	cluster, err := oceanClient.CreateCluster(ctx, x.buildClusterFromOpts())
	if err != nil {
		return err
	}

	fmt.Fprintln(x.opts.Out, cluster.ID)
	return nil
}

func (x *CmdCreateClusterKubernetes) buildClusterFromOpts() *spotinst.OceanCluster {
	var cluster interface{}

	switch x.opts.CloudProvider {
	case spotinst.CloudProviderAWS:
		cluster = x.buildClusterFromOptsAWS()
	}

	return &spotinst.OceanCluster{Obj: cluster}
}

func (x *CmdCreateClusterKubernetes) buildClusterFromOptsAWS() *aws.Cluster {
	cluster := new(aws.Cluster)

	if x.opts.Name != "" {
		cluster.SetName(spotinstsdk.String(x.opts.Name))
	}

	return cluster
}

func (x *CmdCreateClusterKubernetesOptions) Init(flags *pflag.FlagSet, opts *CmdCreateClusterOptions) {
	x.initDefaults(opts)
	x.initFlags(flags)
}

func (x *CmdCreateClusterKubernetesOptions) initDefaults(opts *CmdCreateClusterOptions) {
	x.CmdCreateClusterOptions = opts
}

func (x *CmdCreateClusterKubernetesOptions) initFlags(flags *pflag.FlagSet) {
	// Base
	{
		flags.StringVar(
			&x.Name,
			"name",
			x.Name,
			"name of the cluster")

		flags.StringVar(
			&x.Region,
			"region",
			x.Region,
			"")
	}

	// Strategy
	{
		flags.Float64Var(
			&x.SpotPercentage,
			"spot-percentage",
			x.SpotPercentage,
			"")

		flags.BoolVar(
			&x.UtilizeReservedInstances,
			"utilize-reserved-instances",
			x.UtilizeReservedInstances,
			"")

		flags.BoolVar(
			&x.FallbackToOnDemand,
			"fallback-ondemand",
			x.FallbackToOnDemand,
			"")
	}

	// Capacity
	{
		flags.IntVar(
			&x.MinSize,
			"min-size",
			x.MinSize,
			"")

		flags.IntVar(
			&x.MaxSize,
			"max-size",
			x.MaxSize,
			"")

		flags.IntVar(
			&x.TargetSize,
			"target-size",
			x.TargetSize,
			"")
	}

	// Compute
	{
		flags.StringSliceVar(
			&x.SubnetIDs,
			"subnet-ids",
			x.SubnetIDs,
			"")

		flags.StringSliceVar(
			&x.InstanceTypesWhitelist,
			"instance-types-whitelist",
			x.InstanceTypesWhitelist,
			"")

		flags.StringSliceVar(
			&x.InstanceTypesBlacklist,
			"instance-types-blacklist",
			x.InstanceTypesBlacklist,
			"")

		flags.StringSliceVar(
			&x.SecurityGroupIDs,
			"security-group-ids",
			x.SecurityGroupIDs,
			"")

		flags.StringVar(
			&x.ImageID,
			"image-id",
			x.ImageID,
			"")

		flags.StringVar(
			&x.KeyPair,
			"key-pair",
			x.KeyPair,
			"")

		flags.StringVar(
			&x.UserData,
			"user-data",
			x.UserData,
			"")

		flags.IntVar(
			&x.RootVolumeSize,
			"root-volume-size",
			x.RootVolumeSize,
			"")

		flags.BoolVar(
			&x.AssociatePublicIPAddress,
			"associate-public-ip-address",
			x.AssociatePublicIPAddress,
			"")

		flags.BoolVar(
			&x.EnableMonitoring,
			"enable-monitoring",
			x.EnableMonitoring,
			"")

		flags.BoolVar(
			&x.EnableEBSOptimization,
			"enable-ebs-optimization",
			x.EnableEBSOptimization,
			"")

	}

	// Auto Scaling
	{
		flags.BoolVar(
			&x.EnableAutoScaler,
			"enable-auto-scaler",
			x.EnableAutoScaler,
			"")

		flags.BoolVar(
			&x.EnableAutoConfig,
			"enable-auto-scaler-autoconfig",
			x.EnableAutoConfig,
			"")

		flags.IntVar(
			&x.Cooldown,
			"cooldown",
			x.Cooldown,
			"")
	}
}

func (x *CmdCreateClusterKubernetesOptions) Validate() error {
	return x.CmdCreateClusterOptions.Validate()
}
