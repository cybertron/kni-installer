// Package baremetal collects bare metal specific configuration.
package baremetal

import (
	"encoding/json"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift-metalkube/kni-installer/pkg/types/baremetal"
	baremetaldefaults "github.com/openshift-metalkube/kni-installer/pkg/types/baremetal/defaults"
	"github.com/openshift-metalkube/kni-installer/pkg/validate"
)

// Platform collects bare metal specific configuration.
func Platform() (*baremetal.Platform, error) {
	var libvirtURI, ironicURI, nodesJSON, apiVIP string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Libvirt Connection URI",
				Help:    "The libvirt connection URI to be used.",
				Default: baremetaldefaults.LibvirtURI,
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, &libvirtURI)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Ironic Connection URI",
				Help:    "The ironic connection URI to be used.",
				Default: baremetaldefaults.IronicURI,
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, &libvirtURI)
	if err != nil {
		return nil, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Master node definition JSON",
				Help:    "JSON data containing information about the baremetal nodes for use by Ironic.",
			},
		},
	}, &nodesJSON)
	if err != nil {
		return nil, err
	}

	var nodes map[string]interface{}
	if err = json.Unmarshal([]byte(nodesJSON), &nodes); err != nil {
		return nil, err
	}

	// Is there some place we can populate this based on a DNS query instead of
	// adding a survey question? Seems like there should be.
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "API VIP",
				Help:    "The IP pointed to by the API DNS entry.",
			},
			Validate: survey.ComposeValidators(survey.Required, ipValidator),
		},
	}, &apiVIP)
	if err != nil {
		return nil, err
	}

	return &baremetal.Platform{
		LibvirtURI: libvirtURI,
		IronicURI: ironicURI,
		Nodes: nodes,
		ApiVIP: apiVIP,
	}, nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}

func ipValidator(ans interface{}) error {
	// TODO(bnemec): Implement me
	return nil
}
