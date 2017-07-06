package vmfork

import (
	"context"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
)

var defaultFinder *find.Finder

func DefaultFinder(ctx context.Context, client *vim25.Client) (*find.Finder, error) {
	if defaultFinder == nil {
		debugf("finder.NewFinder()")
		defaultFinder = find.NewFinder(client, true)

		debugf("finder.DefaultDatacenter()")
		dc, err := defaultFinder.DefaultDatacenter(ctx)
		if err != nil {
			return nil, err
		}
		debugf("finder.SetDatacenter(%v)", dc)
		defaultFinder.SetDatacenter(dc)
	}

	return defaultFinder, nil
}

func FindVirtualMachine(ctx context.Context, name string, client *vim25.Client) (*VirtualMachine, error) {
	f, err := DefaultFinder(ctx, client)
	if err != nil {
		return nil, err
	}

	debugf("finder.VirtualMachine(%v)", name)
	vm, err := f.VirtualMachine(ctx, name)
	if err != nil {
		return nil, err
	}

	return &VirtualMachine{name, vm, client}, nil
}
