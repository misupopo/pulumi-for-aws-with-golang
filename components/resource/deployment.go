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
}

type Vpc struct {
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
	}
}

func (d *Deployment) createNewVpc(
	ctx *pulumi.Context,
	region *Region,
	) (*ec2.Vpc, error) {
	newVpc, err := ec2.NewVpc(ctx,
		fmt.Sprintf("%s%s", region.ResourceName, "-vpc"),
		&ec2.VpcArgs{
		CidrBlock: pulumi.String("172.16.0.0/16"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("tf-example"),
		},
	})

	if err != nil {
		return nil, err
	}

	return newVpc, nil
}
