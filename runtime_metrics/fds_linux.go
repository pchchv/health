package runtime_metrics

import "os"

func getFDUsage() (uint64, error) {
	fds, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return 0, err
	}
	return uint64(len(fds)), nil
}
