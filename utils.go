package docker

import (
	"fmt"
	"strings"
	"os"
	"path/filepath"
	"syscall"
)

// Compare two Config struct. Do not compare the "Image" nor "Hostname" fields
// If OpenStdin is set, then it differs
func CompareConfig(a, b *Config) bool {
	if a == nil || b == nil ||
		a.OpenStdin || b.OpenStdin {
		return false
	}
	if a.AttachStdout != b.AttachStdout ||
		a.AttachStderr != b.AttachStderr ||
		a.User != b.User ||
		a.Memory != b.Memory ||
		a.MemorySwap != b.MemorySwap ||
		a.CpuShares != b.CpuShares ||
		a.OpenStdin != b.OpenStdin ||
		a.Tty != b.Tty ||
		a.VolumesFrom != b.VolumesFrom {
		return false
	}
	if len(a.Cmd) != len(b.Cmd) ||
		len(a.Dns) != len(b.Dns) ||
		len(a.Env) != len(b.Env) ||
		len(a.PortSpecs) != len(b.PortSpecs) ||
		len(a.Entrypoint) != len(b.Entrypoint) ||
		len(a.Volumes) != len(b.Volumes) {
		return false
	}

	for i := 0; i < len(a.Cmd); i++ {
		if a.Cmd[i] != b.Cmd[i] {
			return false
		}
	}
	for i := 0; i < len(a.Dns); i++ {
		if a.Dns[i] != b.Dns[i] {
			return false
		}
	}
	for i := 0; i < len(a.Env); i++ {
		if a.Env[i] != b.Env[i] {
			return false
		}
	}
	for i := 0; i < len(a.PortSpecs); i++ {
		if a.PortSpecs[i] != b.PortSpecs[i] {
			return false
		}
	}
	for i := 0; i < len(a.Entrypoint); i++ {
		if a.Entrypoint[i] != b.Entrypoint[i] {
			return false
		}
	}
	for key := range a.Volumes {
		if _, exists := b.Volumes[key]; !exists {
			return false
		}
	}
	return true
}

func MergeConfig(userConf, imageConf *Config) error {
	if userConf.User == "" {
		userConf.User = imageConf.User
	}
	if userConf.Memory == 0 {
		userConf.Memory = imageConf.Memory
	}
	if userConf.MemorySwap == 0 {
		userConf.MemorySwap = imageConf.MemorySwap
	}
	if userConf.CpuShares == 0 {
		userConf.CpuShares = imageConf.CpuShares
	}
	if userConf.PortSpecs == nil || len(userConf.PortSpecs) == 0 {
		userConf.PortSpecs = imageConf.PortSpecs
	} else {
		for _, imagePortSpec := range imageConf.PortSpecs {
			found := false
			imageNat, err := parseNat(imagePortSpec)
			if err != nil {
				return err
			}
			for _, userPortSpec := range userConf.PortSpecs {
				userNat, err := parseNat(userPortSpec)
				if err != nil {
					return err
				}
				if imageNat.Proto == userNat.Proto && imageNat.Backend == userNat.Backend {
					found = true
				}
			}
			if !found {
				userConf.PortSpecs = append(userConf.PortSpecs, imagePortSpec)
			}
		}
	}
	if !userConf.Tty {
		userConf.Tty = imageConf.Tty
	}
	if !userConf.OpenStdin {
		userConf.OpenStdin = imageConf.OpenStdin
	}
	if !userConf.StdinOnce {
		userConf.StdinOnce = imageConf.StdinOnce
	}
	if userConf.Env == nil || len(userConf.Env) == 0 {
		userConf.Env = imageConf.Env
	} else {
		for _, imageEnv := range imageConf.Env {
			found := false
			imageEnvKey := strings.Split(imageEnv, "=")[0]
			for _, userEnv := range userConf.Env {
				userEnvKey := strings.Split(userEnv, "=")[0]
				if imageEnvKey == userEnvKey {
					found = true
				}
			}
			if !found {
				userConf.Env = append(userConf.Env, imageEnv)
			}
		}
	}
	if userConf.Cmd == nil || len(userConf.Cmd) == 0 {
		userConf.Cmd = imageConf.Cmd
	}
	if userConf.Dns == nil || len(userConf.Dns) == 0 {
		userConf.Dns = imageConf.Dns
	} else {
		//duplicates aren't an issue here
		userConf.Dns = append(userConf.Dns, imageConf.Dns...)
	}
	if userConf.Entrypoint == nil || len(userConf.Entrypoint) == 0 {
		userConf.Entrypoint = imageConf.Entrypoint
	}
	if userConf.WorkingDir == "" {
		userConf.WorkingDir = imageConf.WorkingDir
	}
	if userConf.VolumesFrom == "" {
		userConf.VolumesFrom = imageConf.VolumesFrom
	}
	if userConf.Volumes == nil || len(userConf.Volumes) == 0 {
		userConf.Volumes = imageConf.Volumes
	} else {
		for k, v := range imageConf.Volumes {
			userConf.Volumes[k] = v
		}
	}
	return nil
}

func parseLxcConfOpts(opts ListOpts) ([]KeyValuePair, error) {
	out := make([]KeyValuePair, len(opts))
	for i, o := range opts {
		k, v, err := parseLxcOpt(o)
		if err != nil {
			return nil, err
		}
		out[i] = KeyValuePair{Key: k, Value: v}
	}
	return out, nil
}

func parseLxcOpt(opt string) (string, string, error) {
	parts := strings.SplitN(opt, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Unable to parse lxc conf option: %s", opt)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func FileTreeSize(root string) uint64 {
	var totalSize int64 = 0
	filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {
		totalSize += fileInfo.Size()
		return nil
	})
	return uint64(totalSize)
}

func FileSystemAvailSize(path string) (uint64, error) {
	var st syscall.Statfs_t
	err := syscall.Statfs(path, &st)
	if err != nil {
		return 0, err
	}
	return uint64(st.Bsize) * st.Bavail, nil
}
