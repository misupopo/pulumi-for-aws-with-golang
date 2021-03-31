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
) ([]*lb.TargetGroup, error) {
	// 文字数制限32文字以内に収めないといけないためtgに省略
	targetGroup1, err := lb.NewTargetGroup(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-tg-80"),
		&lb.TargetGroupArgs{
			Port:     pulumi.Int(80),
			Protocol: pulumi.String("HTTP"),
			VpcId:    newVpc.ID(),
		})

	if err != nil {
		return nil, err
	}

	targetGroup2, err := lb.NewTargetGroup(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-tg-31001"),
		&lb.TargetGroupArgs{
			Port:     pulumi.Int(31001),
			Protocol: pulumi.String("HTTP"),
			VpcId:    newVpc.ID(),
		})

	if err != nil {
		return nil, err
	}

	var targetGroup []*lb.TargetGroup
	targetGroup = append([]*lb.TargetGroup{}, targetGroup1, targetGroup2)

	return targetGroup, err
}

func (d *Deployment) createNewListener(
	ctx *pulumi.Context,
	region *Region,
	newLoadBalancer *alb.LoadBalancer,
	newTargetGroup []*lb.TargetGroup,
) ([]*lb.Listener, error) {
	listener1, err := lb.NewListener(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-listener-port-80"),
		&lb.ListenerArgs{
			LoadBalancerArn: newLoadBalancer.Arn,
			Port:            pulumi.Int(80),
			Protocol:        pulumi.String("HTTP"),
			DefaultActions: lb.ListenerDefaultActionArray{
				&lb.ListenerDefaultActionArgs{
					Type:           pulumi.String("forward"),
					TargetGroupArn: newTargetGroup[0].Arn,
				},
			},
		})

	if err != nil {
		return nil, err
	}

	listener2, err := lb.NewListener(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-listener-port-31001"),
		&lb.ListenerArgs{
			LoadBalancerArn: newLoadBalancer.Arn,
			Port:            pulumi.Int(31001),
			Protocol:        pulumi.String("HTTP"),
			DefaultActions: lb.ListenerDefaultActionArray{
				&lb.ListenerDefaultActionArgs{
					Type:           pulumi.String("forward"),
					TargetGroupArn: newTargetGroup[1].Arn,
				},
			},
		})

	if err != nil {
		return nil, err
	}

	var listener []*lb.Listener
	listener = append([]*lb.Listener{}, listener1, listener2)

	return listener, nil
}

func (d *Deployment) createNewListenerRule(
	ctx *pulumi.Context,
	region *Region,
	newListener []*lb.Listener,
	newTargetGroup []*lb.TargetGroup,
) (*lb.ListenerRule, error) {
	listenerRule, err := lb.NewListenerRule(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-listener-roule"),
		&lb.ListenerRuleArgs{
			ListenerArn: newListener[1].Arn,
			Priority:    pulumi.Int(99),
			Actions: lb.ListenerRuleActionArray{
				&lb.ListenerRuleActionArgs{
					Type:           pulumi.String("forward"),
					TargetGroupArn: newTargetGroup[1].Arn,
				},
			},
			Conditions: lb.ListenerRuleConditionArray{
				&lb.ListenerRuleConditionArgs{
					PathPattern: &lb.ListenerRuleConditionPathPatternArgs{
						Values: pulumi.StringArray{
							pulumi.String("/api/"),
						},
					},
				},
			},
		})

	if err != nil {
		return nil, err
	}

	return listenerRule, err
}

func (d *Deployment) createNewTargetGroupAttachment(
	ctx *pulumi.Context,
	region *Region,
	newInstance *ec2.Instance,
	newTargetGroup []*lb.TargetGroup,
) ([]*lb.TargetGroupAttachment, error) {
	targetGroupAttachment1, err := lb.NewTargetGroupAttachment(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-tg-attachment1"),
		&lb.TargetGroupAttachmentArgs{
			TargetGroupArn: newTargetGroup[0].Arn,
			TargetId:       newInstance.ID(),
			Port:           pulumi.Int(80),
		})

	if err != nil {
		return nil, err
	}

	targetGroupAttachment2, err := lb.NewTargetGroupAttachment(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-tg-attachment2"),
		&lb.TargetGroupAttachmentArgs{
			TargetGroupArn: newTargetGroup[1].Arn,
			TargetId:       newInstance.ID(),
			Port:           pulumi.Int(31001),
		})

	if err != nil {
		return nil, err
	}

	var targetGroupAttachment []*lb.TargetGroupAttachment
	targetGroupAttachment = append([]*lb.TargetGroupAttachment{}, targetGroupAttachment1, targetGroupAttachment2)

	return targetGroupAttachment, nil
}
