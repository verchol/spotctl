package ocean

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spotinst/spotinst-cli/internal/spotinst"
	"github.com/spotinst/spotinst-cli/internal/utils/flags"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	spotinstsdk "github.com/spotinst/spotinst-sdk-go/spotinst"
)

type (
	CmdCreateLaunchSpecKubernetes struct {
		cmd  *cobra.Command
		opts CmdCreateLaunchSpecKubernetesOptions
	}

	CmdCreateLaunchSpecKubernetesOptions struct {
		*CmdCreateLaunchSpecOptions

		Name             string
		OceanID          string
		ImageID          string
		UserData         string
		SecurityGroupIDs []string
	}
)

func NewCmdCreateLaunchSpecKubernetes(opts *CmdCreateLaunchSpecOptions) *cobra.Command {
	return newCmdCreateLaunchSpecKubernetes(opts).cmd
}

func newCmdCreateLaunchSpecKubernetes(opts *CmdCreateLaunchSpecOptions) *CmdCreateLaunchSpecKubernetes {
	var cmd CmdCreateLaunchSpecKubernetes

	cmd.cmd = &cobra.Command{
		Use:           "kubernetes",
		Short:         "Create a new Kubernetes launchspec",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(*cobra.Command, []string) error {
			return cmd.Run(context.Background())
		},
	}

	cmd.opts.Init(cmd.cmd.Flags(), opts)

	return &cmd
}

func (x *CmdCreateLaunchSpecKubernetes) Run(ctx context.Context) error {
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

func (x *CmdCreateLaunchSpecKubernetes) survey(ctx context.Context) error {
	if x.opts.Noninteractive {
		return nil
	}

	return nil
}

func (x *CmdCreateLaunchSpecKubernetes) log(ctx context.Context) error {
	flags.Log(x.cmd)
	return nil
}

func (x *CmdCreateLaunchSpecKubernetes) validate(ctx context.Context) error {
	return x.opts.Validate()
}

func (x *CmdCreateLaunchSpecKubernetes) run(ctx context.Context) error {
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

	spec, err := oceanClient.CreateLaunchSpec(ctx, x.buildLaunchSpecFromOpts())
	if err != nil {
		return err
	}

	fmt.Fprintln(x.opts.Out, spec.ID)
	return nil
}

func (x *CmdCreateLaunchSpecKubernetes) buildLaunchSpecFromOpts() *spotinst.OceanLaunchSpec {
	var spec interface{}

	switch x.opts.CloudProvider {
	case spotinst.CloudProviderAWS:
		spec = x.buildLaunchSpecFromOptsAWS()
	}

	return &spotinst.OceanLaunchSpec{Obj: spec}
}

func (x *CmdCreateLaunchSpecKubernetes) buildLaunchSpecFromOptsAWS() *aws.LaunchSpec {
	spec := new(aws.LaunchSpec)

	if x.opts.Name != "" {
		spec.SetName(spotinstsdk.String(x.opts.Name))
	}

	if x.opts.OceanID != "" {
		spec.SetOceanId(spotinstsdk.String(x.opts.OceanID))
	}

	if x.opts.ImageID != "" {
		spec.SetImageId(spotinstsdk.String(x.opts.ImageID))
	}

	if x.opts.UserData != "" {
		if _, err := base64.StdEncoding.DecodeString(x.opts.UserData); err != nil {
			x.opts.UserData = base64.StdEncoding.EncodeToString([]byte(x.opts.UserData))
		}

		spec.SetUserData(spotinstsdk.String(x.opts.UserData))
	}

	if len(x.opts.SecurityGroupIDs) > 0 {
		spec.SetSecurityGroupIDs(x.opts.SecurityGroupIDs)
	}

	return spec
}

func (x *CmdCreateLaunchSpecKubernetesOptions) Init(flags *pflag.FlagSet, opts *CmdCreateLaunchSpecOptions) {
	x.initDefaults(opts)
	x.initFlags(flags)
}

func (x *CmdCreateLaunchSpecKubernetesOptions) initDefaults(opts *CmdCreateLaunchSpecOptions) {
	x.CmdCreateLaunchSpecOptions = opts
}

func (x *CmdCreateLaunchSpecKubernetesOptions) initFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&x.Name,
		"name",
		x.Name,
		"name of the launch spec")

	flags.StringVar(
		&x.OceanID,
		"ocean-id",
		x.OceanID,
		"id of the cluster")

	flags.StringVar(
		&x.ImageID,
		"image-id",
		x.ImageID,
		"id of the image")

	flags.StringVar(
		&x.UserData,
		"user-data",
		x.UserData,
		"user data to provide when launching a node (plain-text or base64-encoded)")
}

func (x *CmdCreateLaunchSpecKubernetesOptions) Validate() error {
	return x.CmdCreateLaunchSpecOptions.Validate()
}
