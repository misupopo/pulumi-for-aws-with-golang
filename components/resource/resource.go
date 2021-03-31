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

	newNetworkInterface, err := deployment.createNewNetworkInterface(ctx, region, newSubnets)

	if err != nil {
		ctx.Export("createNewNetworkInterface error", pulumi.Printf("%v", err))
		return err
	}

	nweInternetGateway, err := deployment.createNewInternetGateway(ctx, region, newVpc)

	if err != nil {
		ctx.Export("createNewNetworkInterface error", pulumi.Printf("%v", err))
		return err
	}

	routeTable, err := deployment.createNewRouteTable(ctx, region, newVpc, nweInternetGateway)

	if err != nil {
		ctx.Export("createNewRouteTable error", pulumi.Printf("%v", err))
		return err
	}

	err = deployment.createNewRouteTableAssociation(ctx, region, newSubnets, routeTable)

	if err != nil {
		ctx.Export("createNewRouteTable error", pulumi.Printf("%v", err))
		return err
	}

	newSecurityGroup, err := deployment.createNewSecurityGroup(ctx, region, newVpc)

	if err != nil {
		ctx.Export("createNewSecurityGroup error", pulumi.Printf("%v", err))
		return err
	}

	newInstance, err := deployment.createNewInstance(ctx, region, newNetworkInterface, newSecurityGroup)

	if err != nil {
		ctx.Export("createNewInstance error", pulumi.Printf("%v", err))
		return err
	}

	_, err = deployment.createNewEip(ctx, region, newInstance)

	if err != nil {
		ctx.Export("createNewEip error", pulumi.Printf("%v", err))
		return err
	}

	newLoadBalancer, err := deployment.createNewLoadBalancer(ctx, region, newSubnets, newSecurityGroup)

	if err != nil {
		ctx.Export("createNewLoadBalancer error", pulumi.Printf("%v", err))
		return err
	}

	newTargetGroup, err := deployment.createNewTargetGroup(ctx, region, newVpc)

	if err != nil {
		ctx.Export("createNewTargetGroup error", pulumi.Printf("%v", err))
		return err
	}

	newListener, err := deployment.createNewListener(ctx, region, newLoadBalancer, newTargetGroup)

	if err != nil {
		ctx.Export("createNewListener error", pulumi.Printf("%v", err))
		return err
	}

	_, err = deployment.createNewListenerRule(ctx, region, newListener, newTargetGroup)

	if err != nil {
		ctx.Export("createNewListener error", pulumi.Printf("%v", err))
		return err
	}

	//ctx.Export("newInstance", pulumi.Printf("%v", newInstance))

	return nil
}

