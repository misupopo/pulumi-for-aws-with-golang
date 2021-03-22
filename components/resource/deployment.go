package resource

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type Region struct {
	ResourceName     string            `json:"ResourceName"`
	Location         string            `json:"Location"`
	Vpc              *Vpc              `json:"vpc"`
	Subnet           *Subnet           `json:"subnet"`
	NetworkInterface *NetworkInterface `json:"networkInterface"`
	Instance         *Instance         `json:"instance"`
}

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

type Instance struct {
	AMI          string `json:"ami"`
	InstanceType string `json:"instanceType"`
	Tag          string `json:"tag"`
}

type Deployment struct {
}

func newRegion(region Region) *Region {
	return &Region{
		ResourceName:     region.ResourceName,
		Location:         region.Location,
		Vpc:              region.Vpc,
		Subnet:           region.Subnet,
		NetworkInterface: region.NetworkInterface,
		Instance:         region.Instance,
	}
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

func (d *Deployment) createNewInstance(
	ctx *pulumi.Context,
	region *Region,
	newNetworkInterface *ec2.NetworkInterface,
) (*ec2.Instance, error) {
	instance, err := ec2.NewInstance(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "Instance"),
		&ec2.InstanceArgs{
		Ami:          pulumi.String(region.Instance.AMI),
		InstanceType: pulumi.String(region.Instance.InstanceType),
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


