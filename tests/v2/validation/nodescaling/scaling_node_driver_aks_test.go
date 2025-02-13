package nodescaling

import (
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/aks"
	"github.com/rancher/rancher/tests/framework/extensions/scalinginput"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AKSNodeScalingTestSuite struct {
	suite.Suite
	client        *rancher.Client
	session       *session.Session
	scalingConfig *scalinginput.Config
}

func (s *AKSNodeScalingTestSuite) TearDownSuite() {
	s.session.Cleanup()
}

func (s *AKSNodeScalingTestSuite) SetupSuite() {
	testSession := session.NewSession()
	s.session = testSession

	s.scalingConfig = new(scalinginput.Config)
	config.LoadConfig(scalinginput.ConfigurationFileKey, s.scalingConfig)

	client, err := rancher.NewClient("", testSession)
	require.NoError(s.T(), err)

	s.client = client
}

func (s *AKSNodeScalingTestSuite) TestScalingAKSNodePools() {
	scaleOneNode := aks.NodePool{
		NodeCount: &oneNode,
	}

	scaleTwoNodes := aks.NodePool{
		NodeCount: &twoNodes,
	}

	tests := []struct {
		name     string
		aksNodes aks.NodePool
		client   *rancher.Client
	}{
		{"Scaling agentpool by 1", scaleOneNode, s.client},
		{"Scaling agentpool by 2", scaleTwoNodes, s.client},
	}

	for _, tt := range tests {
		clusterID, err := clusters.GetClusterIDByName(s.client, s.client.RancherConfig.ClusterName)
		require.NoError(s.T(), err)

		s.Run(tt.name, func() {
			ScalingAKSNodePools(s.T(), s.client, clusterID, &tt.aksNodes)
		})
	}
}

func (s *AKSNodeScalingTestSuite) TestScalingAKSNodePoolsDynamicInput() {
	if s.scalingConfig.AKSNodePool == nil {
		s.T().Skip()
	}

	clusterID, err := clusters.GetClusterIDByName(s.client, s.client.RancherConfig.ClusterName)
	require.NoError(s.T(), err)

	ScalingAKSNodePools(s.T(), s.client, clusterID, s.scalingConfig.AKSNodePool)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAKSNodeScalingTestSuite(t *testing.T) {
	suite.Run(t, new(AKSNodeScalingTestSuite))
}
