// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-09-01-preview/hcpopenshiftclusters"
)

func TestHcpOpenShiftClusterPlatformExpandFlatten(t *testing.T) {
	input := []HcpOpenShiftClusterPlatformProfile{{
		ManagedResourceGroup:    "managed-rg",
		SubnetId:                "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/cluster",
		VnetIntegrationSubnetId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/integration",
		OutboundType:            string(hcpopenshiftclusters.OutboundTypeLoadBalancer),
		NetworkSecurityGroupId:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/nsg",
		OperatorsAuthentication: []HcpOpenShiftClusterOperatorsAuthentication{{
			UserAssignedIdentities: []HcpOpenShiftClusterUserAssignedIdentities{{
				ControlPlaneOperators: map[string]string{
					"cloud-controller-manager": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/control-plane",
				},
				DataPlaneOperators: map[string]string{
					"image-registry": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/data-plane",
				},
				ServiceManagedIdentity: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/service",
			}},
		}},
	}}

	expanded := expandHcpOpenShiftClusterPlatform(input)
	expanded.IssuerURL = "https://issuer.example.com"
	actual, err := flattenHcpOpenShiftClusterPlatform(&expanded)
	if err != nil {
		t.Fatalf("flattening platform: %+v", err)
	}
	input[0].IssuerUrl = "https://issuer.example.com"

	if !reflect.DeepEqual(actual, input) {
		t.Fatalf("unexpected platform round-trip\nactual: %#v\nexpected: %#v", actual, input)
	}
}

func TestHcpOpenShiftClusterPlatformFlattenCanonicalizesResourceIds(t *testing.T) {
	input := hcpopenshiftclusters.PlatformProfile{
		SubnetId:                "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.NETWORK/VIRTUALNETWORKS/vnet/SUBNETS/cluster",
		VnetIntegrationSubnetId: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.NETWORK/VIRTUALNETWORKS/vnet/SUBNETS/integration",
		NetworkSecurityGroupId:  "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.NETWORK/NETWORKSECURITYGROUPS/nsg",
		OperatorsAuthentication: hcpopenshiftclusters.OperatorsAuthenticationProfile{
			UserAssignedIdentities: hcpopenshiftclusters.UserAssignedIdentitiesProfile{
				ControlPlaneOperators: map[string]string{
					"cloud-controller-manager": "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/control-plane",
				},
				DataPlaneOperators: map[string]string{
					"image-registry": "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/data-plane",
				},
				ServiceManagedIdentity: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/rg/PROVIDERS/MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/service",
			},
		},
	}

	actual, err := flattenHcpOpenShiftClusterPlatform(&input)
	if err != nil {
		t.Fatalf("flattening platform: %+v", err)
	}

	expected := []HcpOpenShiftClusterPlatformProfile{{
		SubnetId:                "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/cluster",
		VnetIntegrationSubnetId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/integration",
		NetworkSecurityGroupId:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/nsg",
		OperatorsAuthentication: []HcpOpenShiftClusterOperatorsAuthentication{{
			UserAssignedIdentities: []HcpOpenShiftClusterUserAssignedIdentities{{
				ControlPlaneOperators: map[string]string{
					"cloud-controller-manager": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/control-plane",
				},
				DataPlaneOperators: map[string]string{
					"image-registry": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/data-plane",
				},
				ServiceManagedIdentity: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/service",
			}},
		}},
	}}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected canonical platform\nactual: %#v\nexpected: %#v", actual, expected)
	}
}

func TestHcpOpenShiftClusterPlatformFlattenRejectsInvalidResourceIds(t *testing.T) {
	validPlatform := func() hcpopenshiftclusters.PlatformProfile {
		return hcpopenshiftclusters.PlatformProfile{
			SubnetId:                "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/cluster",
			VnetIntegrationSubnetId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/integration",
			NetworkSecurityGroupId:  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/nsg",
			OperatorsAuthentication: hcpopenshiftclusters.OperatorsAuthenticationProfile{
				UserAssignedIdentities: hcpopenshiftclusters.UserAssignedIdentitiesProfile{
					ControlPlaneOperators: map[string]string{
						"control-plane": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/control-plane",
					},
					DataPlaneOperators: map[string]string{
						"data-plane": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/data-plane",
					},
					ServiceManagedIdentity: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/service",
				},
			},
		}
	}

	testCases := map[string]struct {
		mutate        func(*hcpopenshiftclusters.PlatformProfile)
		expectedError string
	}{
		"subnet": {
			mutate:        func(input *hcpopenshiftclusters.PlatformProfile) { input.SubnetId = "invalid" },
			expectedError: "`subnet_id`",
		},
		"vnet integration subnet": {
			mutate:        func(input *hcpopenshiftclusters.PlatformProfile) { input.VnetIntegrationSubnetId = "invalid" },
			expectedError: "`vnet_integration_subnet_id`",
		},
		"network security group": {
			mutate:        func(input *hcpopenshiftclusters.PlatformProfile) { input.NetworkSecurityGroupId = "invalid" },
			expectedError: "`network_security_group_id`",
		},
		"control plane operator identity": {
			mutate: func(input *hcpopenshiftclusters.PlatformProfile) {
				input.OperatorsAuthentication.UserAssignedIdentities.ControlPlaneOperators["control-plane"] = "invalid"
			},
			expectedError: "`control_plane_operators`",
		},
		"data plane operator identity": {
			mutate: func(input *hcpopenshiftclusters.PlatformProfile) {
				input.OperatorsAuthentication.UserAssignedIdentities.DataPlaneOperators["data-plane"] = "invalid"
			},
			expectedError: "`data_plane_operators`",
		},
		"service managed identity": {
			mutate: func(input *hcpopenshiftclusters.PlatformProfile) {
				input.OperatorsAuthentication.UserAssignedIdentities.ServiceManagedIdentity = "invalid"
			},
			expectedError: "`service_managed_identity`",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			input := validPlatform()
			testCase.mutate(&input)

			_, err := flattenHcpOpenShiftClusterPlatform(&input)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), testCase.expectedError) {
				t.Fatalf("expected error containing %q, got %q", testCase.expectedError, err)
			}
		})
	}
}

func TestHcpOpenShiftClusterEtcdExpandFlatten(t *testing.T) {
	input := []HcpOpenShiftClusterEtcdProfile{{
		DataEncryption: []HcpOpenShiftClusterEtcdDataEncryptionProfile{{
			KeyManagementMode: string(hcpopenshiftclusters.EtcdDataEncryptionKeyManagementModeTypeCustomerManaged),
			CustomerManaged: []HcpOpenShiftClusterCustomerManagedEncryption{{
				EncryptionType: string(hcpopenshiftclusters.CustomerManagedEncryptionTypeKMS),
				Kms: []HcpOpenShiftClusterKmsEncryption{{
					VaultName:  "vault",
					Visibility: string(hcpopenshiftclusters.KeyVaultVisibilityPrivate),
					ActiveKey: []HcpOpenShiftClusterKmsKey{{
						Name:    "key",
						Version: "version-1",
					}},
				}},
			}},
		}},
	}}

	expanded := expandHcpOpenShiftClusterEtcd(input)
	actual := flattenHcpOpenShiftClusterEtcd(expanded)
	if !reflect.DeepEqual(actual, input) {
		t.Fatalf("unexpected etcd round-trip\nactual: %#v\nexpected: %#v", actual, input)
	}
}

func TestHcpOpenShiftClusterApiExpandFlatten(t *testing.T) {
	input := []HcpOpenShiftClusterApiProfile{{
		Visibility:      string(hcpopenshiftclusters.VisibilityPublic),
		AuthorizedCidrs: []string{"10.0.0.0/24", "192.168.0.0/24"},
	}}

	expanded := expandHcpOpenShiftClusterApi(input)
	expanded.Url = "https://api.example.com"
	actual := flattenHcpOpenShiftClusterApi(expanded)
	input[0].Url = "https://api.example.com"

	if !reflect.DeepEqual(actual, input) {
		t.Fatalf("unexpected API round-trip\nactual: %#v\nexpected: %#v", actual, input)
	}
}
