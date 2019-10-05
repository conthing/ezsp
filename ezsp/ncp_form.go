package ezsp

import (
	"fmt"
	"math/rand"
	"time"

	"encoding/binary"

	"github.com/conthing/utils/common"
)

// NcpFormNetwork radioChannel=0xff时自动根据能量扫描选择channel
func NcpFormNetwork(radioChannel byte) (err error) {
	var channelMask uint32
	if radioChannel == 0xff {
		channelMask = EMBER_RECOMMENDED_802_15_4_CHANNELS_MASK
	} else if radioChannel >= EMBER_MIN_802_15_4_CHANNEL_NUMBER && radioChannel <= EMBER_MAX_802_15_4_CHANNEL_NUMBER {
		channelMask = 1 << radioChannel
	} else {
		return fmt.Errorf("unsupported channel %d", radioChannel)
	}
	err = ncpTrustCenterInit()
	if err != nil {
		return
	}
	return ncpStartScan(channelMask)
}

func ncpTrustCenterInit() (err error) {
	emberInitialSecurityState := EmberInitialSecurityState{}
	emberInitialSecurityState.bitmask |= EMBER_TRUST_CENTER_GLOBAL_LINK_KEY
	emberInitialSecurityState.bitmask |= EMBER_HAVE_PRECONFIGURED_KEY
	emberInitialSecurityState.bitmask |= EMBER_HAVE_NETWORK_KEY
	emberInitialSecurityState.bitmask |= EMBER_NO_FRAME_COUNTER_RESET
	emberInitialSecurityState.bitmask |= EMBER_REQUIRE_ENCRYPTED_KEY
	emberInitialSecurityState.bitmask |= EMBER_DISTRIBUTED_TRUST_CENTER_MODE
	copy(emberInitialSecurityState.preconfiguredKey[:], "ZigBeeAlliance09")
	//生成随机networkKey
	rand.Seed(time.Now().Unix())
	binary.LittleEndian.PutUint64(emberInitialSecurityState.networkKey[:], rand.Uint64())
	binary.LittleEndian.PutUint64(emberInitialSecurityState.networkKey[8:], rand.Uint64())

	err = EzspSetInitialSecurityState(&emberInitialSecurityState)
	if err != nil {
		return fmt.Errorf("EzspSetInitialSecurityState failed: %v", err)
	}

	extended := EMBER_JOINER_GLOBAL_LINK_KEY
	err = EzspSetValue_EXTENDED_SECURITY_BITMASK(extended)
	if err != nil {
		return fmt.Errorf("EzspSetValue_EXTENDED_SECURITY_BITMASK failed: %v", err)
	}

	err = EzspSetPolicy(EZSP_TC_KEY_REQUEST_POLICY, EZSP_DENY_TC_KEY_REQUESTS)
	if err != nil {
		return fmt.Errorf("EzspSetPolicy EZSP_DENY_TC_KEY_REQUESTS failed: %v", err)
	}

	err = EzspSetPolicy(EZSP_APP_KEY_REQUEST_POLICY, EZSP_ALLOW_APP_KEY_REQUESTS)
	if err != nil {
		return fmt.Errorf("EzspSetPolicy EZSP_ALLOW_APP_KEY_REQUESTS failed: %v", err)
	}

	err = EzspSetPolicy(EZSP_TRUST_CENTER_POLICY, EZSP_ALLOW_PRECONFIGURED_KEY_JOINS)
	if err != nil {
		return fmt.Errorf("EzspSetPolicy EZSP_ALLOW_PRECONFIGURED_KEY_JOINS failed: %v", err)
	}

	return
}

const (
	FORM_AND_JOIN_NOT_SCANNING  = byte(0)
	FORM_AND_JOIN_NEXT_NETWORK  = byte(1)
	FORM_AND_JOIN_ENERGY_SCAN   = byte(2)
	FORM_AND_JOIN_PAN_ID_SCAN   = byte(3)
	FORM_AND_JOIN_JOINABLE_SCAN = byte(4)

	// The minimum significant difference between energy scan results.
	// Results that differ by less than this are treated as identical.
	ENERGY_SCAN_FUZZ = byte(25)

	NUM_PAN_ID_CANDIDATES = 16

	// ZigBee specifies that active scans have a duration of 3 (138 msec).
	// See documentation for emberStartScan in include/network-formation.h
	// for more info on duration values.
	ACTIVE_SCAN_DURATION = byte(3)

	// 507 ms duration
	ENERGY_SCAN_DURATION = byte(5)
)

var formAndJoinScanType = FORM_AND_JOIN_NOT_SCANNING
var networkCount byte
var channelEnergies [EMBER_NUM_802_15_4_CHANNELS]byte
var panIdCandidates [NUM_PAN_ID_CANDIDATES]uint16
var channelCache byte

func ncpStartScan(channelMask uint32) (err error) {
	if isScanning() {
		return fmt.Errorf("already in scan")
	}
	formAndJoinScanType = FORM_AND_JOIN_ENERGY_SCAN
	networkCount = 0
	for i := range channelEnergies {
		channelEnergies[i] = byte(0xff)
	}
	err = startScan(FORM_AND_JOIN_ENERGY_SCAN, channelMask, ENERGY_SCAN_DURATION)
	return
}

func EzspEnergyScanResultHandler(channel byte, maxRssiValue byte) {
	if isScanning() {
		common.Log.Debug("SCAN: found energy ", maxRssiValue, " dBm on channel ", channel)
		channelEnergies[channel-EMBER_MIN_802_15_4_CHANNEL_NUMBER] = maxRssiValue
	}
}

func EzspScanCompleteHandler(channel byte, emberStatus byte) {
	common.Log.Debug("ezspScanCompleteHandler channel ", channel, ", status ", emberStatus, ", formAndJoinScanType ", formAndJoinScanType)
	if !isScanning() {
		common.Log.Error("not in scaning")
		return
	}

	if FORM_AND_JOIN_ENERGY_SCAN != formAndJoinScanType {
		// This scan is an Active Scan.
		// Active Scans potentially report transmit failures through this callback.
		if EMBER_SUCCESS != emberStatus {
			// The Active Scan is still in progress.  This callback is informing us
			// about a failure to transmit the beacon request on this channel.
			// If necessary we could save this failing channel number and start
			// another Active Scan on this channel later (after this current scan is
			// complete).
			common.Log.Error("ezspScanCompleteHandler status error")
			return
		}
	}

	switch formAndJoinScanType {
	case FORM_AND_JOIN_ENERGY_SCAN:
		energyScanComplete()
	case FORM_AND_JOIN_PAN_ID_SCAN:
		panIdScanComplete()
	default:
		common.Log.Error("unknown scan completed ", formAndJoinScanType)
	}
}

func EzspNetworkFoundHandler(networkFound *EmberZigbeeNetwork, lqi byte, rssi int8) {
	common.Log.Debug("SCAN: found ", networkFound, ", lqi ", lqi, ", rssi: ", rssi)

	switch formAndJoinScanType {

	case FORM_AND_JOIN_PAN_ID_SCAN:
		for i := 0; i < NUM_PAN_ID_CANDIDATES; i++ {
			if panIdCandidates[i] == networkFound.PanId {
				panIdCandidates[i] = uint16(0xFFFF)
			}
		}

	default:
		common.Log.Error("unknown scan  ", formAndJoinScanType)
	}
}

func isScanning() bool {
	return formAndJoinScanType >= FORM_AND_JOIN_ENERGY_SCAN
}

// Pick a channel from among those with the lowest energy and then look for
// a PAN ID not in use on that channel.
//
// The energy scans are not particularly accurate, especially as we don't run
// them for very long, so we add in some slop to the measurements and then pick
// a random channel from the least noisy ones.  This avoids having several
// coordinators pick the same slightly quieter channel.
func energyScanComplete() {
	cutoff := byte(0xFF)
	candidateCount := byte(0)
	var channelIndex byte
	var i int

	// cutoff = min energy + ENERGY_SCAN_FUZZ
	for i = 0; i < EMBER_NUM_802_15_4_CHANNELS; i++ {
		if channelEnergies[i] < cutoff-ENERGY_SCAN_FUZZ {
			cutoff = channelEnergies[i] + ENERGY_SCAN_FUZZ
		}
	}

	// There must be at least one channel,
	// so there will be at least one candidate.
	// 能量低于cutoff的频道比较适合创建新的网络
	for i = 0; i < EMBER_NUM_802_15_4_CHANNELS; i++ {
		if channelEnergies[i] < cutoff {
			candidateCount++
		}
	}

	// If for some reason we never got any energy scan results
	// then our candidateCount will be 0.  We want to avoid that case and
	// bail out (since we will do a divide by 0 below)
	if candidateCount == 0 {
		formAndJoinScanType = FORM_AND_JOIN_NOT_SCANNING
		common.Log.Error("never got any energy scan results")
		return
	}

	// 在这些candidate中随机取第channelIndex个
	channelIndex = byte(rand.Uint32()) % candidateCount

	for i = 0; i < EMBER_NUM_802_15_4_CHANNELS; i++ {
		if channelEnergies[i] < cutoff {
			if channelIndex == 0 {
				channelCache = byte(EMBER_MIN_802_15_4_CHANNEL_NUMBER + i)
				break
			}
			channelIndex--
		}
	}

	common.Log.Debug("select channel ", channelCache, ", start PANID scan")
	startPanIdScan()
}

// Form a network using one of the unused PAN IDs.  If we got unlucky we
// pick some more and try again.
func panIdScanComplete() {

	for i := 0; i < NUM_PAN_ID_CANDIDATES; i++ {
		if panIdCandidates[i] != 0xFFFF {
			unusedPanIdFoundHandler(panIdCandidates[i], channelCache)
			formAndJoinScanType = FORM_AND_JOIN_NOT_SCANNING
			return
		}
	}

	// XXX: Do we care this could keep happening forever?
	// In practice there couldn't be as many PAN IDs heard that
	// conflict with ALL our randomly selected set of candidate PANs.
	// But in theory we could get the same random set of numbers
	// (more likely due to a bug) and we could hear the same set of
	// PAN IDs that conflict with our random set.

	startPanIdScan() // Start over with new candidates.
}

func startPanIdScan() {

	// PAN IDs can be 0..0xFFFE.  We pick some trial candidates and then do a scan
	// to find one that is not in use.
	for i := 0; i < NUM_PAN_ID_CANDIDATES; {
		panId := uint16(rand.Uint32())
		if panId != 0xFFFF {
			panIdCandidates[i] = panId
			i++
		}
	}

	formAndJoinScanType = FORM_AND_JOIN_PAN_ID_SCAN
	startScan(EZSP_ACTIVE_SCAN, uint32(1)<<channelCache, ACTIVE_SCAN_DURATION)
}

func startScan(scanType byte, channelMask uint32, duration byte) (err error) {
	err = EzspStartScan(scanType, channelMask, duration)
	if err != nil {
		formAndJoinScanType = FORM_AND_JOIN_NOT_SCANNING
	}
	return
}

func unusedPanIdFoundHandler(panId uint16, channel byte) {
	networkParams := EmberNetworkParameters{}
	networkParams.RadioChannel = channel
	networkParams.PanId = panId
	networkParams.RadioTxPower = -1

	common.Log.Debug("unusedPanIdFoundHandler")
	err := EzspFormNetwork(&networkParams)
	if err != nil {
		common.Log.Errorf("ezsp form error: %v", err)
	}
}
