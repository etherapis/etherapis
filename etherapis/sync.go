// Contains the sync status reporting.

package main

import (
	"time"

	"github.com/gophergala2016/etherapis/etherapis/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
	"github.com/gophergala2016/etherapis/etherapis/geth"
)

// monitorSync runs in an inifite loop, periodically checking if the attached
// node is synchronizing or not, and writing some log entries.
func monitorSync(api *geth.API) {
	start := time.Now()

	for {
		// Wait a bit before checking the status, first is irrelevant anyway
		time.Sleep(time.Second)

		// Retrieve the sync status and go back to sleep if not syncing
		status, err := api.Syncing()
		if err != nil {
			log15.Error("Failed to retrieve sync status", "error", err)
			continue
		}
		if status == nil {
			start = time.Now()
			continue
		}
		// Yay, something's actually happening, display a log
		if status.CurrentBlock > status.StartingBlock {
			var (
				totalBlocks  = status.HighestBlock - status.StartingBlock
				pulledBlocks = status.CurrentBlock - status.StartingBlock
				estimate     = time.Since(start) * time.Duration(totalBlocks) / time.Duration(pulledBlocks)
			)
			log15.Info("Synchronizing with the network...", "at", status.CurrentBlock, "total", status.HighestBlock, "eta", estimate)
		}
	}
}

// waitSync blocks execution until the chain is synchronized with the network,
// which usually entails waiting until syncing reports false, with the added
// criteria that we actually have a decently current block (< 3 mins).
func waitSync(freshness time.Duration, api *geth.API) {
	notified := false

	for {
		// Wait until sync reports false
		for {
			// Fetch the sync status
			status, err := api.Syncing()
			if err != nil {
				log15.Error("Failed to retrieve sync status", "error", err)
			} else if status == nil {
				break
			}
			// Sleep if syncing
			time.Sleep(250 * time.Millisecond)
		}
		// We're supposedly in sync, double check the block timestamp
		head, err := api.BlockNumber()
		if err != nil {
			log15.Error("Failed to retrieve head block number", "error", err)
			break
		}
		mined, err := api.GetBlockTime(head)
		if err != nil {
			log15.Error("Failed to retrieve head block timestamp", "number", head, "error", err)
			break
		}
		if time.Since(mined) <= freshness {
			log15.Info("In sync with the network", "block", head, "freshness", time.Since(mined))
			return
		}
		// Seems we're not in sync, wait a bit and retry
		if !notified {
			log15.Info("You seem out of sync, updating...", "freshness", time.Since(mined), "allowed", freshness)
			notified = true
		}
		time.Sleep(250 * time.Millisecond)
	}
}
