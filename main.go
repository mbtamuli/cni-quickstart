// go:build darwin || linux
package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

type PluginConf struct {
	RuntimeConfig *struct{} `json:"runtimeConfig,omitempty"`

	RawPrevResult *map[string]interface{} `json:"prevResult"`
	PrevResult    *current.Result         `json:"-"`

	Debug    bool   `json:"debug"`
	DebugDir string `json:"debugDir"`
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("dummy"))
}

func cmdAdd(args *skel.CmdArgs) error {
	conf, _, err := parseConfig(args.StdinData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	if conf.PrevResult == nil {
		return fmt.Errorf("must be used as a chained plugin")
	}

	result, err := current.NewResultFromResult(conf.PrevResult)
	if err != nil {
		return fmt.Errorf("unable to generate new results: %v", err)
	}

	_, ipNet, err := net.ParseCIDR("1.2.3.4/24")
	if err != nil {
		return fmt.Errorf("unable to parse CIDR: %v", err)
	}

	result.IPs = []*current.IPConfig{
		{
			Address:   *ipNet,
			Gateway:   net.ParseIP("1.2.3.1"),
			Interface: current.Int(1),
		},
	}

	return types.PrintResult(result, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func cmdCheck(args *skel.CmdArgs) error {
	return nil
}

func parseConfig(stdin []byte) (*types.NetConf, *current.Result, error) {
	conf := types.NetConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	// Parse previous result.
	var result *current.Result
	if conf.RawPrevResult != nil {
		var err error
		if err = version.ParsePrevResult(&conf); err != nil {
			return nil, nil, fmt.Errorf("could not parse prevResult: %v", err)
		}

		result, err = current.NewResultFromResult(conf.PrevResult)
		if err != nil {
			return nil, nil, fmt.Errorf("could not convert result to current version: %v", err)
		}
	}

	return &conf, result, nil
}
