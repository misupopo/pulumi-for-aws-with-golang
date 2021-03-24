package resource

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type Vpc struct {
	Cidr string `json:"cidr"`
	Tag  string `json:"tag"`
}

type Subnet struct {
	Cidr             string `json:"cidr"`
	AvailabilityZone string `json:"availabilityZone"`
	Tag              string `json:"tag"`
}

type NetworkInterface struct {
	PrivateIp string `json:"privateIp"`
	Tag       string `json:"tag"`
}

type SecurityGroup struct {
	Ingress []Ingress `json:"ingress"`
	Egress  []Egress  `json:"egress"`
}

type Ingress struct {
	Protocol    string   `json:"protocol"`
	ToPort      int      `json:"toPort"`
	FromPort    int      `json:"fromPort"`
	Description string   `json:"description"`
	CidrBlocks  []string `json:"cidrBlocks"`
}

type Egress struct {
	Protocol    string   `json:"protocol"`
	ToPort      int      `json:"toPort"`
	FromPort    int      `json:"fromPort"`
	Description string   `json:"description"`
	CidrBlocks  []string `json:"cidrBlocks"`
}

type Instance struct {
	AMI          string `json:"ami"`
	InstanceType string `json:"instanceType"`
	Tag          string `json:"tag"`
}

func (d *Deployment) createNewVpc(
	ctx *pulumi.Context,
	region *Region,
) (*ec2.Vpc, error) {
	newVpc, err := ec2.NewVpc(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-vpc"),
		&ec2.VpcArgs{
			CidrBlock: pulumi.String(region.Vpc.Cidr),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(region.Vpc.Tag),
			},
		})

	if err != nil {
		return nil, err
	}

	return newVpc, nil
}

func (d *Deployment) createNewSubnet(
	ctx *pulumi.Context,
	region *Region,
	newVpc *ec2.Vpc,
) (*ec2.Subnet, error) {
	newSubnet, err := ec2.NewSubnet(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-subnet"),
		&ec2.SubnetArgs{
			VpcId:            newVpc.ID(),
			CidrBlock:        pulumi.String(region.Subnet.Cidr),
			AvailabilityZone: pulumi.String(region.Subnet.AvailabilityZone),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(region.Subnet.Tag),
			},
		})

	if err != nil {
		return nil, err
	}

	return newSubnet, nil
}

func (d *Deployment) createNetworkInterface(
	ctx *pulumi.Context,
	region *Region,
	newSubnet *ec2.Subnet,
) (*ec2.NetworkInterface, error) {
	networkInterface, err := ec2.NewNetworkInterface(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-networkInterface"),
		&ec2.NetworkInterfaceArgs{
			SubnetId: newSubnet.ID(),
			PrivateIps: pulumi.StringArray{
				pulumi.String(region.NetworkInterface.PrivateIp),
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(region.NetworkInterface.Tag),
			},
		})

	if err != nil {
		return nil, err
	}

	return networkInterface, nil
}

func (d *Deployment) createNewSecurityGroup(
	ctx *pulumi.Context,
	region *Region,
	newVpc *ec2.Vpc,
) (*ec2.SecurityGroup, error) {
	//var createdEgress []ec2.SecurityGroupIngressArray
	//
	//for _, v := range region.SecurityGroup.Ingress {
	//	createdEgress = append(createdEgress, ec2.SecurityGroupIngressArgs{
	//		Protocol:    pulumi.String(v.Protocol),
	//		ToPort:      pulumi.Int(v.ToPort),
	//		FromPort:    pulumi.Int(v.FromPort),
	//		Description: pulumi.String(v.Description),
	//		CidrBlocks:  pulumi.StringArray{
	//			pulumi.String("0.0.0.0/0"),
	//		},
	//	})
	//	//createdEgress[i] = ec2.SecurityGroupIngressArgs{
	//	//	Protocol:    pulumi.String(v.Protocol),
	//	//	ToPort:      pulumi.Int(v.ToPort),
	//	//	FromPort:    pulumi.Int(v.FromPort),
	//	//	Description: pulumi.String(v.Description),
	//	//}
	//}

	securityGroup, err := ec2.NewSecurityGroup(ctx,
		"ssh-sg",
		&ec2.SecurityGroupArgs{
			Name:        pulumi.String("ssh-sg"),
			VpcId:       newVpc.ID(),
			Description: pulumi.String("Allows SSH traffic to bastion hosts"),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:    pulumi.String(region.SecurityGroup.Ingress[0].Protocol),
					ToPort:      pulumi.Int(region.SecurityGroup.Ingress[0].ToPort),
					FromPort:    pulumi.Int(region.SecurityGroup.Ingress[0].FromPort),
					Description: pulumi.String(region.SecurityGroup.Ingress[0].Description),
					CidrBlocks:  pulumi.StringArray{
						pulumi.String(region.SecurityGroup.Ingress[0].CidrBlocks[0]),
					},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:    pulumi.String(region.SecurityGroup.Egress[0].Protocol),
					ToPort:      pulumi.Int(region.SecurityGroup.Egress[0].ToPort),
					FromPort:    pulumi.Int(region.SecurityGroup.Egress[0].FromPort),
					Description: pulumi.String(region.SecurityGroup.Egress[0].Description),
					CidrBlocks:  pulumi.StringArray{
						pulumi.String(region.SecurityGroup.Egress[0].CidrBlocks[0]),
					},
				},
			},
		})

	if err != nil {
		return nil, err
	}

	return securityGroup, nil
}

func (d *Deployment) createNewInstance(
	ctx *pulumi.Context,
	region *Region,
	newNetworkInterface *ec2.NetworkInterface,
	newSecurityGroup *ec2.SecurityGroup,
) (*ec2.Instance, error) {
	instance, err := ec2.NewInstance(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "Instance"),
		&ec2.InstanceArgs{
			VpcSecurityGroupIds: pulumi.StringArray{newSecurityGroup.Name},
			Ami:                 pulumi.String(region.Instance.AMI),
			InstanceType:        pulumi.String(region.Instance.InstanceType),
			NetworkInterfaces: ec2.InstanceNetworkInterfaceArray{
				&ec2.InstanceNetworkInterfaceArgs{
					NetworkInterfaceId: newNetworkInterface.ID(),
					DeviceIndex:        pulumi.Int(0),
				},
			},
		})

	if err != nil {
		return nil, err
	}

	return instance, nil
}
