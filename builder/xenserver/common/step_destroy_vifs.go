package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	xsclient "github.com/xenserver/go-xenserver-client"
)

type StepDestroyVIFs struct{}

func (self *StepDestroyVIFs) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("commonconfig").(CommonConfig)
	if !config.DestroyVIFs {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packer.Ui)
	client := state.Get("client").(xsclient.XenAPIClient)

	ui.Say("Step: Destroy VIFs")

	uuid := state.Get("instance_uuid").(string)
	instance, err := client.GetVMByUuid(uuid)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to get VM from UUID '%s': %s", uuid, err.Error()))
		return multistep.ActionHalt
	}

	vifs, err := instance.GetVIFs()
	if err != nil {
		ui.Error(fmt.Sprintf("Error getting VIFs: %s", err.Error()))
		return multistep.ActionHalt
	}

	for _, vif := range vifs {
		err = vif.Destroy()
		if err != nil {
			ui.Error(fmt.Sprintf("Error destroying VIF: %s", err.Error()))
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (self *StepDestroyVIFs) Cleanup(state multistep.StateBag) {}
