package ipvanish

import "github.com/Marlinski/icannzplit/util"

// Stop all openvpn and cleanup routes
func Stop() {
	close(stopChan)
	wg.Wait()
	util.Log.Noticef("all routes cleanup")
}
