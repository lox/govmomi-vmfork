package types

import (
	vstypes "github.com/vmware/govmomi/vim25/types"
)

// vm.enableForkParent

type EnableForkParentRequestType struct {
	This vstypes.ManagedObjectReference  `xml:"_this"`
	Host *vstypes.ManagedObjectReference `xml:"host,omitempty"`
}

type EnableForkParent_Task EnableForkParentRequestType

type EnableForkParent_TaskResponse struct {
	Returnval vstypes.ManagedObjectReference `xml:"returnval"`
}

// vm.createForkChild
// -----------------------------------------

type VirtualMachineCreateChildSpec struct {
	vstypes.DynamicData

	Persistent *bool `xml:"persistent"`
}

type CreateForkChildRequestType struct {
	This vstypes.ManagedObjectReference `xml:"_this"`
	Name string                         `xml:"name"`
	Spec VirtualMachineCreateChildSpec  `xml:"spec"`
}

type CreateForkChild_Task CreateForkChildRequestType

type CreateForkChild_TaskResponse struct {
	Returnval vstypes.ManagedObjectReference `xml:"returnval"`
}
