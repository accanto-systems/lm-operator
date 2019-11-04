package controller

import (
	"github.com/orgs/accanto-systems/lm-operator/pkg/controller/alm"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, alm.Add)
}
