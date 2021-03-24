package resource

type Region struct {
	ResourceName     string            `json:"ResourceName"`
	Location         string            `json:"Location"`
	Vpc              *Vpc              `json:"vpc"`
	Subnet           *Subnet           `json:"subnet"`
	NetworkInterface *NetworkInterface `json:"networkInterface"`
	SecurityGroup    *SecurityGroup    `json:"SecurityGroup"`
	Instance         *Instance         `json:"instance"`
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
		SecurityGroup:    region.SecurityGroup,
		Instance:         region.Instance,
	}
}
