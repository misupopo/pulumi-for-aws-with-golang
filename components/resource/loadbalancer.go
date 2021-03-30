package resource

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/alb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type LoadBalancer struct {
	LoadBalancerType string `json:"loadBalancerType"`
}

func (d *Deployment) createNewLoadBalancer(
	ctx *pulumi.Context,
	region *Region,
	newSubnets []*ec2.Subnet,
	newSecurityGroup *ec2.SecurityGroup,
) (*alb.LoadBalancer, error) {
	loadBalancer, err := alb.NewLoadBalancer(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-loadbalancer"),
		&alb.LoadBalancerArgs{
			Name: pulumi.String(fmt.Sprintf("%s%s", region.ResourceName, "-loadbalancer")),
			Internal: pulumi.Bool(false),
			LoadBalancerType: pulumi.String(region.LoadBalancer.LoadBalancerType),
			Subnets: pulumi.StringArray{
				newSubnets[0].ID(),
				newSubnets[1].ID(),
			},
			SecurityGroups: pulumi.StringArray{
				newSecurityGroup.ID(),
			},
		})

	if err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (d *Deployment) createNewTargetGroup(
	ctx *pulumi.Context,
	region *Region,
	newVpc *ec2.Vpc,
) (*lb.TargetGroup, error) {
	targetGroup, err := lb.NewTargetGroup(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-target-group"),
		&lb.TargetGroupArgs{
			Port:     pulumi.Int(80),
			Protocol: pulumi.String("HTTP"),
			VpcId:    newVpc.ID(),
		})

	if err != nil {
		return nil, err
	}

	return targetGroup, err
}

