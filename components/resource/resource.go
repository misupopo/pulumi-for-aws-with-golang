package resource

import (
	"encoding/json"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"io/ioutil"
)

func Setup(ctx *pulumi.Context) error {
	config, _ := ioutil.ReadFile("./config.json")

	var r Region

	json.Unmarshal(config, &r)

	region := newRegion(r)

	deployment := new(Deployment)

	newVpc, err := deployment.createNewVpc(ctx, region)

	if err != nil {
		ctx.Export("createNewVpc error", pulumi.Printf("%v", err))
		return err
	}

	newSubnets, err := deployment.createNewSubnet(ctx, region, newVpc)

	if err != nil {
		ctx.Export("createNewSubnet error", pulumi.Printf("%v", err))
		return err
	}

	newNetworkInterface, err := deployment.createNetworkInterface(ctx, region, newSubnets)

	if err != nil {
		ctx.Export("createNetworkInterface error", pulumi.Printf("%v", err))
		return err
	}

	newSecurityGroup, err := deployment.createNewSecurityGroup(ctx, region, newVpc)

	if err != nil {
		ctx.Export("createNewSecurityGroup error", pulumi.Printf("%v", err))
		return err
	}

	_, err = deployment.createNewInstance(ctx, region, newNetworkInterface, newSecurityGroup)

	if err != nil {
		ctx.Export("createNewInstance error", pulumi.Printf("%v", err))
		return err
	}

	//ctx.Export("newInstance", pulumi.Printf("%v", newInstance))

	return nil
}

