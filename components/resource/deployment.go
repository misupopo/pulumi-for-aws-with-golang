package resource

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type Region struct {
	ResourceName string `json:"ResourceName"`
	Location     string `json:"Location"`
	Vpc          *Vpc   `json:"vpc"`
	Subnet       *Vpc   `json:"subnet"`
}

type Vpc struct {
	Cidr string `json:"cidr"`
	Tag  string `json:"tag"`
}

type Subnet struct {
	Cidr string `json:"cidr"`
	Tag  string `json:"tag"`
}

type Deployment struct {
}

func newRegion(region Region) *Region{
	return &Region{
		ResourceName: region.ResourceName,
		Location:     region.Location,
		Vpc:          region.Vpc,
		Subnet:       region.Subnet,
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
	newSubnet, err := ec2.NewSubnet(ctx, "mySubnet", &ec2.SubnetArgs{
		VpcId:            newVpc.ID(),
		CidrBlock:        pulumi.String(region.Subnet.Cidr),
		AvailabilityZone: pulumi.String(region.Location),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(region.Subnet.Tag),
		},
	})

	if err != nil {
		return nil, err
	}

	return newSubnet, nil
}
