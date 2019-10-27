package spotinst

import (
	"context"

	"github.com/spotinst/spotinst-cli/internal/log"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
)

type oceanKubernetesAWS struct {
	svc aws.Service
}

func (x *oceanKubernetesAWS) ListClusters(ctx context.Context) ([]*OceanCluster, error) {
	log.Debugf("Listing all Kubernetes clusters")

	output, err := x.svc.ListClusters(ctx, &aws.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	clusters := make([]*OceanCluster, len(output.Clusters))
	for i, cluster := range output.Clusters {
		clusters[i] = &OceanCluster{
			TypeMeta: TypeMeta{
				Kind: typeOf(OceanCluster{}),
			},
			ObjectMeta: ObjectMeta{
				ID:        spotinst.StringValue(cluster.ID),
				Name:      spotinst.StringValue(cluster.Name),
				CreatedAt: spotinst.TimeValue(cluster.CreatedAt),
				UpdatedAt: spotinst.TimeValue(cluster.UpdatedAt),
			},
			Obj: cluster,
		}
	}

	return clusters, nil
}

func (x *oceanKubernetesAWS) ListLaunchSpecs(ctx context.Context) ([]*OceanLaunchSpec, error) {
	log.Debugf("Listing all Kubernetes launch specs")

	output, err := x.svc.ListLaunchSpecs(ctx, &aws.ListLaunchSpecsInput{})
	if err != nil {
		return nil, err
	}

	specs := make([]*OceanLaunchSpec, len(output.LaunchSpecs))
	for i, spec := range output.LaunchSpecs {
		specs[i] = &OceanLaunchSpec{
			TypeMeta: TypeMeta{
				Kind: typeOf(OceanLaunchSpec{}),
			},
			ObjectMeta: ObjectMeta{
				ID:        spotinst.StringValue(spec.ID),
				Name:      spotinst.StringValue(spec.Name),
				CreatedAt: spotinst.TimeValue(spec.CreatedAt),
				UpdatedAt: spotinst.TimeValue(spec.UpdatedAt),
			},
			Obj: spec,
		}
	}

	return specs, nil
}

func (x *oceanKubernetesAWS) GetCluster(ctx context.Context, clusterID string) (*OceanCluster, error) {
	log.Debugf("Getting a Kubernetes cluster by ID: %s", clusterID)

	input := &aws.ReadClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	output, err := x.svc.ReadCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	cluster := &OceanCluster{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.Cluster.ID),
			Name:      spotinst.StringValue(output.Cluster.Name),
			CreatedAt: spotinst.TimeValue(output.Cluster.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.Cluster.UpdatedAt),
		},
		Obj: output.Cluster,
	}

	return cluster, nil
}

func (x *oceanKubernetesAWS) GetLaunchSpec(ctx context.Context, specID string) (*OceanLaunchSpec, error) {
	log.Debugf("Getting a Kubernetes launch spec by ID: %s", specID)

	input := &aws.ReadLaunchSpecInput{
		LaunchSpecID: spotinst.String(specID),
	}

	output, err := x.svc.ReadLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	spec := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanLaunchSpec{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.LaunchSpec.ID),
			Name:      spotinst.StringValue(output.LaunchSpec.Name),
			CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return spec, nil
}

func (x *oceanKubernetesAWS) CreateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error) {
	log.Debugf("Creating a new Kubernetes cluster")

	input := &aws.CreateClusterInput{
		Cluster: cluster.Obj.(*aws.Cluster),
	}

	output, err := x.svc.CreateCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	created := &OceanCluster{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.Cluster.ID),
			Name:      spotinst.StringValue(output.Cluster.Name),
			CreatedAt: spotinst.TimeValue(output.Cluster.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.Cluster.UpdatedAt),
		},
		Obj: output.Cluster,
	}

	return created, nil
}

func (x *oceanKubernetesAWS) CreateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error) {
	log.Debugf("Creating a new Kubernetes launch spec")

	input := &aws.CreateLaunchSpecInput{
		LaunchSpec: spec.Obj.(*aws.LaunchSpec),
	}

	output, err := x.svc.CreateLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	created := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.LaunchSpec.ID),
			Name:      spotinst.StringValue(output.LaunchSpec.Name),
			CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return created, nil
}

func (x *oceanKubernetesAWS) UpdateCluster(ctx context.Context, cluster *OceanCluster) (*OceanCluster, error) {
	log.Debugf("Updating a Kubernetes cluster by ID: %s", cluster.ID)

	input := &aws.UpdateClusterInput{
		Cluster: cluster.Obj.(*aws.Cluster),
	}

	// Remove read-only fields.
	input.Cluster.Region = nil
	input.Cluster.UpdatedAt = nil
	input.Cluster.CreatedAt = nil

	output, err := x.svc.UpdateCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	updated := &OceanCluster{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.Cluster.ID),
			Name:      spotinst.StringValue(output.Cluster.Name),
			CreatedAt: spotinst.TimeValue(output.Cluster.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.Cluster.UpdatedAt),
		},
		Obj: output.Cluster,
	}

	return updated, nil
}

func (x *oceanKubernetesAWS) UpdateLaunchSpec(ctx context.Context, spec *OceanLaunchSpec) (*OceanLaunchSpec, error) {
	log.Debugf("Updating a Kubernetes launch spec by ID: %s", spec.ID)

	input := &aws.UpdateLaunchSpecInput{
		LaunchSpec: spec.Obj.(*aws.LaunchSpec),
	}

	// Remove read-only fields.
	input.LaunchSpec.UpdatedAt = nil
	input.LaunchSpec.CreatedAt = nil

	output, err := x.svc.UpdateLaunchSpec(ctx, input)
	if err != nil {
		return nil, err
	}

	updated := &OceanLaunchSpec{
		TypeMeta: TypeMeta{
			Kind: typeOf(OceanCluster{}),
		},
		ObjectMeta: ObjectMeta{
			ID:        spotinst.StringValue(output.LaunchSpec.ID),
			Name:      spotinst.StringValue(output.LaunchSpec.Name),
			CreatedAt: spotinst.TimeValue(output.LaunchSpec.CreatedAt),
			UpdatedAt: spotinst.TimeValue(output.LaunchSpec.UpdatedAt),
		},
		Obj: output.LaunchSpec,
	}

	return updated, nil
}

func (x *oceanKubernetesAWS) DeleteCluster(ctx context.Context, clusterID string) error {
	log.Debugf("Deleting a Kubernetes cluster by ID: %s", clusterID)

	input := &aws.DeleteClusterInput{
		ClusterID: spotinst.String(clusterID),
	}

	_, err := x.svc.DeleteCluster(ctx, input)
	return err
}

func (x *oceanKubernetesAWS) DeleteLaunchSpec(ctx context.Context, specID string) error {
	log.Debugf("Deleting a Kubernetes launch spec by ID: %s", specID)

	input := &aws.DeleteLaunchSpecInput{
		LaunchSpecID: spotinst.String(specID),
	}

	_, err := x.svc.DeleteLaunchSpec(ctx, input)
	return err
}