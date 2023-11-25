package balancer

import (
	"math/rand"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

var (
	MinWeight = 1
	MaxWeight = 15
)

// weightKey is the type used as the key to store WeightAddrInfo in the Attributes
// field of resolver.Address.
type weightKey struct{}

// WeightAddrInfo will be stored inside Address metadata in order to use weighted balancer.
type WeightAddrInfo struct {
	Weight int
}

// SetAddrInfo returns a copy of addr in which the Attributes field is updated
// with addrInfo.
func SetAddrInfo(addr resolver.Address, addrInfo WeightAddrInfo) resolver.Address {
	addr.Attributes = attributes.New(weightKey{}, addrInfo)
	return addr
}

// GetAddrInfo returns the WeightAddrInfo stored in the Attributes fields of addr.
func GetAddrInfo(addr resolver.Address) WeightAddrInfo {
	v := addr.Attributes.Value(weightKey{})
	ai, _ := v.(WeightAddrInfo)
	return ai
}

type weightPickerBuilder struct {
}

// Name is the name of weight balancer.
const Name = "weight"

// NewBuilder creates a new weight balancer builder.
func newWeightPicker() balancer.Builder {
	return base.NewBalancerBuilder(Name, &weightPickerBuilder{}, base.Config{HealthCheck: false})
}

func init() {
	balancer.Register(newWeightPicker())
}

func (*weightPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("weightPicker: newPicker called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for subConn, addr := range info.ReadySCs {
		node := GetAddrInfo(addr.Address)
		if node.Weight <= 0 {
			node.Weight = MinWeight
		} else if node.Weight > MaxWeight {
			node.Weight = MaxWeight
		}
		for i := 0; i < node.Weight; i++ {
			scs = append(scs, subConn)
		}
	}
	return &weightPicker{
		subConns: scs,
	}
}

type weightPicker struct {
	subConns []balancer.SubConn
}

func (p *weightPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	sc := p.subConns[rand.Intn(len(p.subConns))]
	return balancer.PickResult{SubConn: sc}, nil
}
