/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudconfig

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	aws "github.com/kubermatic/machine-controller/pkg/cloudprovider/provider/aws/types"
	azure "github.com/kubermatic/machine-controller/pkg/cloudprovider/provider/azure/types"
	gce "github.com/kubermatic/machine-controller/pkg/cloudprovider/provider/gce/types"
	kubevirt "github.com/kubermatic/machine-controller/pkg/cloudprovider/provider/kubevirt/types"
	vsphere "github.com/kubermatic/machine-controller/pkg/cloudprovider/provider/vsphere/types"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/provider/cloud/gcp"
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/resources/cloudconfig/openstack"
)

// CloudConfig returns the cloud-config for the supplied data.
func CloudConfig(
	cluster *kubermaticv1.Cluster,
	dc *kubermaticv1.Datacenter,
	credentials resources.Credentials,
) (cloudConfig string, err error) {
	cloud := cluster.Spec.Cloud
	switch {
	case cloud.AWS != nil:
		DualstackIPFamilies := []string{"ipv4", "ipv6"}
		awsCloudConfig := &aws.CloudConfig{
			// Dummy AZ, so that K8S can extract the region from it.
			// https://github.com/kubernetes/kubernetes/blob/v1.15.0/staging/src/k8s.io/legacy-cloud-providers/aws/aws.go#L1199
			// https://github.com/kubernetes/kubernetes/blob/v1.15.0/staging/src/k8s.io/legacy-cloud-providers/aws/aws.go#L1174
			Global: aws.GlobalOpts{
				Zone:                        dc.Spec.AWS.Region + "x",
				VPC:                         cloud.AWS.VPCID,
				KubernetesClusterID:         cluster.Name,
				DisableSecurityGroupIngress: false,
				RouteTableID:                cloud.AWS.RouteTableID,
				DisableStrictZoneCheck:      true,
				RoleARN:                     cloud.AWS.ControlPlaneRoleARN,
			},
		}
		if cluster.IsDualStack() {
			awsCloudConfig.Global.NodeIPFamilies = DualstackIPFamilies
		}

		// in 1.27, the AWS CCM removed the DisableStrictZoneCheck field, so
		// for those setups we must set it to false, which will omit it from
		// the generated cloudConfig
		if cluster.Spec.Version.Semver().Minor() >= 27 {
			awsCloudConfig.Global.DisableStrictZoneCheck = false
		}

		cloudConfig, err = aws.CloudConfigToString(awsCloudConfig)
		if err != nil {
			return cloudConfig, err
		}

	case cloud.Azure != nil:
		azureCloudConfig := &azure.CloudConfig{
			Cloud:                      "AZUREPUBLICCLOUD",
			TenantID:                   credentials.Azure.TenantID,
			SubscriptionID:             credentials.Azure.SubscriptionID,
			AADClientID:                credentials.Azure.ClientID,
			AADClientSecret:            credentials.Azure.ClientSecret,
			ResourceGroup:              cloud.Azure.ResourceGroup,
			Location:                   dc.Spec.Azure.Location,
			VNetName:                   cloud.Azure.VNetName,
			SubnetName:                 cloud.Azure.SubnetName,
			RouteTableName:             cloud.Azure.RouteTableName,
			SecurityGroupName:          cloud.Azure.SecurityGroup,
			PrimaryAvailabilitySetName: cloud.Azure.AvailabilitySet,
			VnetResourceGroup:          cloud.Azure.VNetResourceGroup,
			UseInstanceMetadata:        false,
			LoadBalancerSku:            string(cloud.Azure.LoadBalancerSKU),
		}
		cloudConfig, err = azure.CloudConfigToString(azureCloudConfig)
		if err != nil {
			return cloudConfig, err
		}

	case cloud.Openstack != nil:
		manageSecurityGroups := dc.Spec.Openstack.ManageSecurityGroups
		trustDevicePath := dc.Spec.Openstack.TrustDevicePath
		useOctavia := dc.Spec.Openstack.UseOctavia
		if cluster.Spec.Cloud.Openstack.UseOctavia != nil {
			useOctavia = cluster.Spec.Cloud.Openstack.UseOctavia
		}
		openstackCloudConfig := &openstack.CloudConfig{
			Global: openstack.GlobalOpts{
				AuthURL:                     dc.Spec.Openstack.AuthURL,
				Username:                    credentials.Openstack.Username,
				Password:                    credentials.Openstack.Password,
				DomainName:                  credentials.Openstack.Domain,
				ProjectName:                 credentials.Openstack.Project,
				ProjectID:                   credentials.Openstack.ProjectID,
				Region:                      dc.Spec.Openstack.Region,
				ApplicationCredentialSecret: credentials.Openstack.ApplicationCredentialSecret,
				ApplicationCredentialID:     credentials.Openstack.ApplicationCredentialID,
			},
			BlockStorage: openstack.BlockStorageOpts{
				BSVersion:       "auto",
				TrustDevicePath: trustDevicePath != nil && *trustDevicePath,
				IgnoreVolumeAZ:  dc.Spec.Openstack.IgnoreVolumeAZ,
			},
			LoadBalancer: openstack.LoadBalancerOpts{
				ManageSecurityGroups: manageSecurityGroups == nil || *manageSecurityGroups,
				UseOctavia:           useOctavia,
			},
			Version: cluster.Status.Versions.ControlPlane.String(),
		}

		if cluster.Spec.Cloud.Openstack.EnableIngressHostname != nil {
			openstackCloudConfig.LoadBalancer.EnableIngressHostname = cluster.Spec.Cloud.Openstack.EnableIngressHostname
		}

		if cluster.Spec.Cloud.Openstack.IngressHostnameSuffix != nil {
			openstackCloudConfig.LoadBalancer.IngressHostnameSuffix = cluster.Spec.Cloud.Openstack.IngressHostnameSuffix
		}

		cloudConfig, err = openstack.CloudConfigToString(openstackCloudConfig)
		if err != nil {
			return cloudConfig, err
		}

	case cloud.VSphere != nil:
		vsphereCloudConfig, err := GetVSphereCloudConfig(cluster, dc, credentials)
		if err != nil {
			return cloudConfig, err
		}
		cloudConfig, err = vsphere.CloudConfigToString(vsphereCloudConfig)
		if err != nil {
			return cloudConfig, err
		}

	case cloud.GCP != nil:
		b, err := base64.StdEncoding.DecodeString(credentials.GCP.ServiceAccount)
		if err != nil {
			return "", fmt.Errorf("error decoding service account: %w", err)
		}
		sam := map[string]string{}
		err = json.Unmarshal(b, &sam)
		if err != nil {
			return "", fmt.Errorf("failed unmarshaling service account: %w", err)
		}
		projectID := sam["project_id"]
		if projectID == "" {
			return "", errors.New("empty project_id")
		}

		tag := fmt.Sprintf("kubernetes-cluster-%s", cluster.Name)

		if len(dc.Spec.GCP.ZoneSuffixes) == 0 {
			return "", errors.New("empty zoneSuffixes")
		}

		localZone := dc.Spec.GCP.Region + "-" + dc.Spec.GCP.ZoneSuffixes[0]

		// By default, all GCP clusters are assumed to be the in the same zone. If the control plane
		// and worker nodes are not it the same zone (localZone), the GCP cloud controller fails
		// to find nodes that are not in the localZone: https://github.com/kubermatic/kubermatic/issues/5025
		// to avoid this, we should enable multizone or regional configuration.
		// It's not easily possible to access the MachineDeployment object from here to compare
		// localZone with the user cluster zone. Additionally, ZoneSuffixes are not used
		// to limit available zones for the user. So, we will just enable multizone support as a workaround.

		// FIXME: Compare localZone to MachineDeployment.Zone and set multizone to true
		// when they differ, or if len(dc.Spec.GCP.ZoneSuffixes) > 1
		multizone := true

		if cloud.GCP.Network == "" || cloud.GCP.Network == gcp.DefaultNetwork {
			// NetworkName is used by the gce cloud provider to populate the provider's NetworkURL.
			// This value can be provided in the config as a name or a url. Internally,
			// the gce cloud provider checks it and if it's a name, it will infer the URL from it.
			// However, if the name has a '/', the provider assumes it's a URL and uses it as is.
			// This breaks routes cleanup since the routes are matched against the URL,
			// which would be incorrect in this case.
			// On the provider side, the "global/networks/default" format is the valid
			// one since it's used internally for firewall rules and and network interfaces,
			// so it has to be kept this way.
			// tl;dr: use "default" or a full network URL, not "global/networks/default"
			cloud.GCP.Network = "default"
		}

		gcpCloudConfig := &gce.CloudConfig{
			Global: gce.GlobalOpts{
				ProjectID:      projectID,
				LocalZone:      localZone,
				MultiZone:      multizone,
				Regional:       dc.Spec.GCP.Regional,
				NetworkName:    cloud.GCP.Network,
				SubnetworkName: cloud.GCP.Subnetwork,
				TokenURL:       "nil",
				NodeTags:       []string{tag},
			},
		}
		cloudConfig, err = gcpCloudConfig.AsString()
		if err != nil {
			return cloudConfig, err
		}

	case cloud.Kubevirt != nil:
		cc := kubevirt.CloudConfig{
			Kubeconfig: "/etc/kubernetes/cloud/infra-kubeconfig",
			Namespace:  cluster.Status.NamespaceName,
		}
		return cc.String()
	}

	return cloudConfig, err
}

func GetVSphereCloudConfig(
	cluster *kubermaticv1.Cluster,
	dc *kubermaticv1.Datacenter,
	credentials resources.Credentials,
) (*vsphere.CloudConfig, error) {
	vsphereURL, err := url.Parse(dc.Spec.VSphere.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vsphere endpoint: %w", err)
	}
	port := "443"
	if urlPort := vsphereURL.Port(); urlPort != "" {
		port = urlPort
	}
	datastore := dc.Spec.VSphere.DefaultDatastore
	// if a datastore is provided at cluster level override the default
	// datastore provided at datacenter level.
	if cluster.Spec.Cloud.VSphere.Datastore != "" {
		datastore = cluster.Spec.Cloud.VSphere.Datastore
	} else if cluster.Spec.Cloud.VSphere.DatastoreCluster != "" {
		datastore = cluster.Spec.Cloud.VSphere.DatastoreCluster
	}

	// Originally, we have been setting cluster-id to the vSphere Compute Cluster name
	// (provided via the Datacenter object), however, this is supposed to identify the
	// Kubernetes cluster, therefore it must be unique. This feature flag is enabled by
	// default for new vSphere clusters, while existing vSphere clusters must be
	// migrated manually (preferably by following advice here:
	// https://kb.vmware.com/s/article/84446).
	clusterID := dc.Spec.VSphere.Cluster
	if cluster.Spec.Features[kubermaticv1.ClusterFeatureVsphereCSIClusterID] {
		clusterID = cluster.Name
	}

	cc := &vsphere.CloudConfig{
		Global: vsphere.GlobalOpts{
			User:             credentials.VSphere.Username,
			Password:         credentials.VSphere.Password,
			VCenterIP:        vsphereURL.Hostname(),
			VCenterPort:      port,
			InsecureFlag:     dc.Spec.VSphere.AllowInsecure,
			Datacenter:       dc.Spec.VSphere.Datacenter,
			DefaultDatastore: datastore,
			WorkingDir:       cluster.Name,
			ClusterID:        clusterID,
		},
		Workspace: vsphere.WorkspaceOpts{
			// This is redundant with what the Vsphere cloud provider itself does:
			// https://github.com/kubernetes/kubernetes/blob/9d80e7522ab7fc977e40dd6f3b5b16d8ebfdc435/pkg/cloudprovider/providers/vsphere/vsphere.go#L346
			// We do it here because the fields in the "Global" object
			// are marked as deprecated even thought the code checks
			// if they are set and will make the controller-manager crash
			// if they are not - But maybe that will change at some point
			VCenterIP:        vsphereURL.Hostname(),
			Datacenter:       dc.Spec.VSphere.Datacenter,
			Folder:           cluster.Spec.Cloud.VSphere.Folder,
			DefaultDatastore: datastore,
		},
		Disk: vsphere.DiskOpts{
			SCSIControllerType: "pvscsi",
		},
		VirtualCenter: map[string]*vsphere.VirtualCenterConfig{
			vsphereURL.Hostname(): {
				User:        credentials.VSphere.Username,
				Password:    credentials.VSphere.Password,
				VCenterPort: port,
				Datacenters: dc.Spec.VSphere.Datacenter,
			},
		},
	}
	if cluster.IsDualStack() && resources.ExternalCloudProviderEnabled(cluster) {
		cc.Global.IPFamily = "ipv4,ipv6"
	}
	return cc, nil
}
