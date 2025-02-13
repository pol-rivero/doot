package install

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pol-rivero/doot/lib/common/log"
	. "github.com/pol-rivero/doot/lib/types"
	"github.com/pol-rivero/doot/lib/utils/optional"
)

type HostnameFilter struct {
	hostSpecificDirPrefix optional.Optional[string]
	ignorePrefixes        []string
}

func getHostnameFilter(hosts map[string]string) HostnameFilter {
	hostname := getHostname()
	result := HostnameFilter{
		hostSpecificDirPrefix: optional.Empty[string](),
		ignorePrefixes:        make([]string, 0, len(hosts)+1),
	}

	// The doot directory should never be symlinked
	result.ignorePrefixes = append(result.ignorePrefixes, "doot"+string(filepath.Separator))

	for host, dir := range hosts {
		dirPrefix := dir + string(filepath.Separator)
		if host == hostname {
			log.Info("Using host-specific directory: %s", dir)
			result.hostSpecificDirPrefix = optional.Of(dirPrefix)
		} else {
			result.ignorePrefixes = append(result.ignorePrefixes, dirPrefix)
		}
	}
	return result
}

func (hf HostnameFilter) isHostSpecific(path RelativePath) (isHostSpecific bool, prefixLen int) {
	if !hf.hostSpecificDirPrefix.HasValue() {
		return false, 0
	}
	prefix := hf.hostSpecificDirPrefix.Value()
	if !strings.HasPrefix(path.Str(), prefix) {
		return false, 0
	}
	return true, len(prefix)
}

func (hf HostnameFilter) isIgnored(path RelativePath) bool {
	for _, prefix := range hf.ignorePrefixes {
		if strings.HasPrefix(path.Str(), prefix) {
			return true
		}
	}
	return false
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Failed to get hostname: %v", err)
		return ""
	}
	return hostname
}
