// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	yaml "gopkg.in/yaml.v3"
)

const (
	ConfigFile = "/etc/slurm-bridge/config.yaml"
)

type Config struct {
	SchedulerName     string   `yaml:"schedulerName"`
	SlurmRestApi      string   `yaml:"slurmRestApi"`
	ManagedNamespaces []string `yaml:"managedNamespaces"`
	MCSLabel          string   `yaml:"mcsLabel"`
	Partition         string   `yaml:"partition"`
}

func Unmarshal(in []byte) (*Config, error) {
	out := &Config{}
	if err := yaml.Unmarshal(in, out); err != nil {
		return nil, err
	}
	return out, nil
}

func UnmarshalOrDie(in []byte) *Config {
	cfg, err := Unmarshal(in)
	if err != nil {
		panic(err)
	}
	return cfg
}
