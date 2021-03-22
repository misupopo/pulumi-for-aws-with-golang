package resource

import (
	"encoding/json"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"io/ioutil"
)

// GetInstancesは
func GetInstances(ctx *pulumi.Context) (*ec2.GetInstancesResult, error) {
	result, err := ec2.GetInstances(ctx, &ec2.GetInstancesArgs{
		InstanceTags: map[string]string{},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func Setup(ctx *pulumi.Context) error {
	config, _ := ioutil.ReadFile("./config.json")

	var r Region

	json.Unmarshal(config, &r)

	region := newRegion(r)

	deployment := new(Deployment)

	newVpc, err := deployment.createNewVpc(ctx, region)

	if err != nil {
		return err
	}

	newSubnet, err := deployment.createNewSubnet(ctx, region, newVpc)

	if err != nil {
		return err
	}

	//ctx.Export("test", pulumi.Printf("%v", r.Vpc))
	//fmt.Printf("%v", r)


	return nil
}

