// SPDX-FileCopyrightText: 2024 Hunter Matthews
// SPDX-License-Identifier: LGPL-2.1-only

package subcmd

import "github.com/alecthomas/kong"

// Registration describes a command that can be added to the root CLI.
type Registration struct {
	Name    string
	Help    string
	Command any
	Tags    []string
}

var registrations []Registration

// Register adds a command registration to the root CLI.
func Register(registration Registration) {
	registrations = append(registrations, registration)
}

// KongOptions returns dynamic Kong command options for all registrations.
func KongOptions() []kong.Option {
	options := make([]kong.Option, 0, len(registrations))
	for _, registration := range registrations {
		options = append(options, kong.DynamicCommand(
			registration.Name,
			registration.Help,
			"", // Kong command group; phase one uses a flat command list.
			registration.Command,
			registration.Tags...,
		))
	}

	return options
}
