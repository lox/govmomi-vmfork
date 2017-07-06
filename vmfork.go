package vmfork

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"

	methods "github.com/lox/govmomi-vmfork/internal/methods"
	types "github.com/lox/govmomi-vmfork/internal/types"
	vstypes "github.com/vmware/govmomi/vim25/types"
)

type VirtualMachine struct {
	Name string
	*object.VirtualMachine
	client *vim25.Client
}

type CreateChildSpec struct {
	Name       string
	Script     string
	Persistent bool
}

func (vm *VirtualMachine) Fork(ctx context.Context, spec CreateChildSpec) error {
	if on, _ := vm.IsPoweredOn(ctx); !on {
		task, err := vm.PowerOn(ctx)
		if err != nil {
			return err
		}
		debugf("waiting for PowerOn(%v)", task)
		if err = task.Wait(ctx); err != nil {
			return err
		}
	}

	q, err := vm.IsQuiescedForkParent(ctx)
	if err != nil {
		return err
	}

	if !q {
		debugf("Parent isn't quiesced")

		err = vm.EnableForkParent(ctx)
		if err != nil {
			return nil
		}

		auth := &vstypes.NamePasswordAuthentication{
			Username: "vmkite",
			Password: "vmkite",
		}

		debugf("Starting %s on %s", spec.Script, vm.Name)
		err = vm.startProgram(ctx, auth, vstypes.GuestProgramSpec{
			ProgramPath: spec.Script,
		})
		if err != nil {
			return nil
		}

		err := vm.AwaitQueiscence(ctx)
		if err != nil {
			return fmt.Errorf("Failed to quiesce parent: %v", err)
		}
	}

	if err = vm.CreateForkChild(ctx, spec); err != nil {
		return err
	}

	f, err := DefaultFinder(ctx, vm.client)
	if err != nil {
		return err
	}

	debugf("finder.VirtualMachine(%v)", spec.Name)
	child, err := f.VirtualMachine(ctx, spec.Name)
	if err != nil {
		return err
	}

	task, err := child.PowerOn(ctx)
	if err != nil {
		return err
	}
	debugf("waiting for PowerOn(%v)", task)
	if err = task.Wait(ctx); err != nil {
		return err
	}

	return nil
}

func (vm *VirtualMachine) IsPoweredOn(ctx context.Context) (bool, error) {
	debugf("vm.PowerState(%v)", vm.Name)
	state, err := vm.PowerState(ctx)
	if err != nil {
		return false, err
	}
	debugf("PowerState = %s", state)
	return state == "poweredOn", nil
}

func (vm *VirtualMachine) AwaitQueiscence(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tick := time.Tick(time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick:
			q, err := vm.IsQuiescedForkParent(ctx)
			if err != nil {
				return err
			}
			if q {
				debugf("Parent is Queisced")
				return nil
			}
		}
	}
}

func (vm *VirtualMachine) IsQuiescedForkParent(ctx context.Context) (bool, error) {
	var o mo.VirtualMachine

	debugf("vm.Properties(%v, %q)", vm.Name, "summary.runtime.quiescedForkParent")
	if err := vm.Properties(ctx, vm.Reference(), []string{"summary.runtime.quiescedForkParent"}, &o); err != nil {
		return false, err
	}

	if o.Summary.Runtime.QuiescedForkParent == nil {
		return false, nil
	}

	return *o.Summary.Runtime.QuiescedForkParent, nil
}

func (vm *VirtualMachine) EnableForkParent(ctx context.Context) error {
	req := types.EnableForkParent_Task{
		This: vm.Reference(),
	}

	debugf("vm.EnableForkParent(%s)", vm.Name)
	res, err := methods.NewEnableForkParent_Task(ctx, vm.client, &req)
	if err != nil {
		return err
	}

	task := object.NewTask(vm.client, res.Returnval)
	debugf("waiting for EnableForkParent %v", task)
	if err := task.Wait(ctx); err != nil {
		return err
	}

	return nil
}

func (vm *VirtualMachine) CreateForkChild(ctx context.Context, spec CreateChildSpec) error {
	req := types.CreateForkChild_Task{
		This: vm.Reference(),
		Name: spec.Name,
		Spec: types.VirtualMachineCreateChildSpec{
			Persistent: &spec.Persistent,
		},
	}

	debugf("vm.CreateForkChild(%s)", req.Name)
	res, err := methods.NewCreateForkChild_Task(ctx, vm.client, &req)
	if err != nil {
		return err
	}

	task := object.NewTask(vm.client, res.Returnval)
	debugf("waiting for CreateForkChild %v", task)
	if err := task.Wait(ctx); err != nil {
		return err
	}

	return nil
}

func (vm *VirtualMachine) startProgram(ctx context.Context, auth vstypes.BaseGuestAuthentication, spec vstypes.GuestProgramSpec) error {
	o := guest.NewOperationsManager(vm.client, vm.Reference())
	procs, err := o.ProcessManager(ctx)

	pid, err := procs.StartProgram(ctx, auth, &spec)
	if err != nil {
		return err
	}

	fmt.Printf("Spawned pid %d\n", pid)
	return nil
}

func debugf(format string, data ...interface{}) {
	log.Printf("[vmfork] "+format, data...)
}
