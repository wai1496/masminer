package antminer

import (
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

// GetStats returns MinerStats
func (c *Client) GetStats() (stats MinerStats, err error) {
	var wg errgroup.Group
	var mu sync.Mutex

	wg.Go(func() error {
		ret, err := getMinerStatsSummary(c.ssh)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		stats.Summary = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getMinerStatsPools(c.ssh)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		stats.Pools = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getMinerStatsDevs(c.ssh)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		stats.Devs = ret
		return nil
	})
	return stats, wg.Wait()
}

func getMinerStatsSummary(client *ssh.Client) (summary MinerStatsSummary, err error) {
	ret, err := outputRemoteShell(client, minerAPISummaryCMD)
	if err != nil {
		return summary, err
	}
	return parseSummaryFromCGMinerSummary(ret)
}

func parseSummaryFromCGMinerSummary(in []byte) (MinerStatsSummary, error) {
	var summary MinerStatsSummary
	lprops := parseCGMinerStats(in)
	props := lprops[1]

	summary.Elapsed = props["Elapsed"]
	summary.GHS5s = props["GHS 5s"]
	summary.GHSAvarage = props["GHS av"]
	summary.Foundblocks = props["Found Blocks"]
	summary.Getworks = props["Getworks"]
	summary.Accepted = props["Accepted"]
	summary.Rejected = props["Rejected"]
	summary.HardwareErrors = props["Hardware Errors"]
	summary.Utility = props["Utility"]
	summary.Discarded = props["Discarded"]
	summary.Stale = props["Stale"]
	summary.Localwork = props["Local Work"]
	summary.WorkUtility = props["Work Utility"]
	summary.DifficultyAccepted = props["Difficulty Accepted"]
	summary.DifficultyRejected = props["Difficulty Rejected"]
	summary.DifficultyStale = props["Difficulty Stale"]
	summary.Bestshare = props["Best Share"]
	return summary, nil
}

func getMinerStatsPools(client *ssh.Client) (pools []MinerStatsPool, err error) {
	ret, err := outputRemoteShell(client, minerAPIPoolsCMD)
	if err != nil {
		return pools, err
	}
	return parsePoolsFromCGMinerPools(ret)
}

func getMinerStatsDevs(client *ssh.Client) (dev MinerStatsDevs, err error) {
	ret, err := outputRemoteShell(client, minerAPIStatsCMD)
	if err != nil {
		return dev, err
	}
	return parseDevsFromCGMinerStats(ret)
}
