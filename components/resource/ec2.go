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

type Subnet []struct {
	Cidr             string `json:"cidr"`
	AvailabilityZone string `json:"availabilityZone"`
	Tag              string `json:"tag"`
}

type NetworkInterface struct {
	PrivateIp string `json:"privateIp"`
	Tag       string `json:"tag"`
}

type SecurityGroup struct {
	Description string `json:"description"`
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
) ([]*ec2.Subnet, error) {
	newSubnet1, err := ec2.NewSubnet(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-subnet1"),
		&ec2.SubnetArgs{
			VpcId:            newVpc.ID(),
			CidrBlock:        pulumi.String((*region.Subnet)[0].Cidr),
			AvailabilityZone: pulumi.String((*region.Subnet)[0].AvailabilityZone),
			Tags: pulumi.StringMap{
				"Name": pulumi.String((*region.Subnet)[0].Tag),
			},
		})

	if err != nil {
		return nil, err
	}

	newSubnet2, err := ec2.NewSubnet(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-subnet2"),
		&ec2.SubnetArgs{
			VpcId:            newVpc.ID(),
			CidrBlock:        pulumi.String((*region.Subnet)[1].Cidr),
			AvailabilityZone: pulumi.String((*region.Subnet)[1].AvailabilityZone),
			Tags: pulumi.StringMap{
				"Name": pulumi.String((*region.Subnet)[1].Tag),
			},
		})

	var subnets []*ec2.Subnet
	subnets = append([]*ec2.Subnet{}, newSubnet1, newSubnet2)

	return subnets, nil
}

func (d *Deployment) createNetworkInterface(
	ctx *pulumi.Context,
	region *Region,
	newSubnet []*ec2.Subnet,
) (*ec2.NetworkInterface, error) {
	networkInterface, err := ec2.NewNetworkInterface(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-networkInterface"),
		&ec2.NetworkInterfaceArgs{
			SubnetId: newSubnet[0].ID(),
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
	var createdIngress ec2.SecurityGroupIngressArray

	for _, v := range region.SecurityGroup.Ingress {
		createdIngress = append(createdIngress, ec2.SecurityGroupIngressArgs{
			Protocol:    pulumi.String(v.Protocol),
			ToPort:      pulumi.Int(v.ToPort),
			FromPort:    pulumi.Int(v.FromPort),
			Description: pulumi.String(v.Description),
			CidrBlocks:  pulumi.StringArray{
				pulumi.String(v.CidrBlocks[0]),
			},
		})
	}

	var createdEgress ec2.SecurityGroupEgressArray

	for _, v := range region.SecurityGroup.Ingress {
		createdEgress = append(createdEgress, ec2.SecurityGroupEgressArgs{
			Protocol:    pulumi.String(v.Protocol),
			ToPort:      pulumi.Int(v.ToPort),
			FromPort:    pulumi.Int(v.FromPort),
			Description: pulumi.String(v.Description),
			CidrBlocks:  pulumi.StringArray{
				pulumi.String(v.CidrBlocks[0]),
			},
		})
	}

	securityGroup, err := ec2.NewSecurityGroup(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-security-group"),
		&ec2.SecurityGroupArgs{
			Name:        pulumi.String(fmt.Sprintf("%s%s", region.ResourceName, "-security-group")),
			VpcId:       newVpc.ID(),
			Description: pulumi.String(region.SecurityGroup.Description),
			Ingress: createdIngress,
			Egress: createdEgress,
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
		fmt.Sprintf("%s%s", region.ResourceName, "-instance"),
		&ec2.InstanceArgs{
			//VpcSecurityGroupIds: pulumi.StringArray{newSecurityGroup.ID()}, これを指定しなくても紐づくのでコメントアウト
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

	_, err = ec2.NewNetworkInterfaceSecurityGroupAttachment(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-security-group-attachment"),
		&ec2.NetworkInterfaceSecurityGroupAttachmentArgs{
			SecurityGroupId:    newSecurityGroup.ID(),
			NetworkInterfaceId: instance.PrimaryNetworkInterfaceId,
		})

	return instance, nil
}

func (d *Deployment) createNewInternetGateway(
	ctx *pulumi.Context,
	region *Region,
	newVpc *ec2.Vpc,
) (*ec2.InternetGateway, error) {
	internetGateway, err := ec2.NewInternetGateway(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-gateway"),
		&ec2.InternetGatewayArgs{
			VpcId: newVpc.ID(),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("%s%s", region.ResourceName, "-gateway")),
			},
		})

	if err != nil {
		return nil, err
	}

	return internetGateway, nil
}

func (d *Deployment) createNewRouteTable(
	ctx *pulumi.Context,
	region *Region,
	newVpc *ec2.Vpc,
	nweInternetGateway *ec2.InternetGateway,
) (*ec2.RouteTable, error) {
	routeTable, err := ec2.NewRouteTable(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-route-table"),
		&ec2.RouteTableArgs{
			VpcId: newVpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: nweInternetGateway.ID(),
				},
			},
		})

	if err != nil {
		return nil, err
	}

	return routeTable, nil
}

func (d *Deployment) createNewRouteTableAssociation(
	ctx *pulumi.Context,
	region *Region,
	newSubnets []*ec2.Subnet,
	newRouteTable *ec2.RouteTable,
) error {
	_, err := ec2.NewRouteTableAssociation(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-route-table-association1"),
		&ec2.RouteTableAssociationArgs{
			SubnetId:     newSubnets[0].ID(),
			RouteTableId: newRouteTable.ID(),
		})

	if err != nil {
		return err
	}

	_, err = ec2.NewRouteTableAssociation(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-route-table-association2"),
		&ec2.RouteTableAssociationArgs{
			SubnetId:     newSubnets[1].ID(),
			RouteTableId: newRouteTable.ID(),
		})

	if err != nil {
		return err
	}

	return nil
}
