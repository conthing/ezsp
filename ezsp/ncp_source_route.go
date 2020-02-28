package ezsp

import (
	"ezsp/ash"

	"github.com/conthing/utils/common"
)

type StSourceRouteTableEntry struct {
	destination uint16
	closerIndex byte // The entry one hop closer to the gateway.
	olderIndex  byte // The entry touched before this one.
}

const (
	// A special index. For destinations that are neighbors of the gateway,
	// closerIndex is set to 0xFF. For the oldest entry, olderIndex is set to
	// 0xFF.
	NULL_INDEX = byte(0xFF)

	/** @brief The size of the source route table on the EZSP host.
	 *
	 * @note This configuration value sets the size of the source route table
	 * on the host, not on the node.
	 * ::EMBER_SOURCE_ROUTE_TABLE_SIZE sets ::EZSP_CONFIG_SOURCE_ROUTE_TABLE_SIZE
	 * if ezsp-utils.c is used, which sets the size of the source route table on
	 * the NCP.
	 */
	EZSP_HOST_SOURCE_ROUTE_TABLE_SIZE = 64
)

var sourceRouteTable [EZSP_HOST_SOURCE_ROUTE_TABLE_SIZE]StSourceRouteTableEntry

// The number of entries in use.
var entryCount = 0

// The index of the most recently added entry.
var newestIndex = NULL_INDEX

var NcpSourceRouteTraceOn bool

func ncpSourceRouteTrace(format string, v ...interface{}) {
	if NcpSourceRouteTraceOn {
		common.Log.Debugf(format, v...)
	}
}

func sourceRouteFindIndex(id uint16) byte {
	for i := 0; i < entryCount; i++ {
		if sourceRouteTable[i].destination == id {
			return byte(i)
		}
	}
	return NULL_INDEX
}

// Create an entry with the given id or update an existing entry. furtherIndex
// is the entry one hop further from the gateway.
func sourceRouteAddEntry(id uint16, furtherIndex byte) byte {
	// See if the id already exists in the table.
	index := sourceRouteFindIndex(id)

	if index == NULL_INDEX {
		if entryCount < EZSP_HOST_SOURCE_ROUTE_TABLE_SIZE {
			// No existing entry. Table is not full. Add new entry.
			index = byte(entryCount)
			entryCount++
		} else {
			// No existing entry. Table is full. Replace oldest entry.
			index = newestIndex
			for sourceRouteTable[index].olderIndex != NULL_INDEX {
				index = sourceRouteTable[index].olderIndex
			}
		}
	}

	// Update the pointers (only) if something has changed.
	if index != newestIndex {
		for i := 0; i < entryCount; i++ {
			if sourceRouteTable[i].olderIndex == index {
				sourceRouteTable[i].olderIndex = sourceRouteTable[index].olderIndex
				break
			}
		}
		sourceRouteTable[index].olderIndex = newestIndex
		newestIndex = index
	}

	// Add the entry.
	sourceRouteTable[index].destination = id
	sourceRouteTable[index].closerIndex = NULL_INDEX

	// The current index is one hop closer to the gateway than furtherIndex.
	if furtherIndex != NULL_INDEX {
		sourceRouteTable[furtherIndex].closerIndex = index
	}

	// Return the current index to save the caller having to look it up.
	return index
}

func EzspIncomingRouteRecordHandler(source uint16, sourceEui uint64, lastHopLqi byte, lastHopRssi int8, relay []uint16) {
	ncpSourceRouteTrace("NCP get source route for 0x%04x, %v", source, relay)
	// The source of the route record is furthest from the gateway. We start there
	// and work closer.
	previous := sourceRouteAddEntry(source, NULL_INDEX)

	// Go through the relay list and add them to the source route table.
	for _, id := range relay {
		// We pass the index of the previous entry to link the route together.
		previous = sourceRouteAddEntry(id, previous)
	}
}

// Note: We assume that the given relayList location is big enough to handle the
// longest source route.
func ncpFindSourceRoute(destination uint16) (exist bool, relayList []uint16) {
	index := sourceRouteFindIndex(destination)

	if index == NULL_INDEX {
		exist = false
		return
	}

	// Fill in the relay list. The first relay in the list is the closest to the
	// destination (furthest from the gateway).
	for sourceRouteTable[index].closerIndex != NULL_INDEX {
		index = sourceRouteTable[index].closerIndex
		relayList = append(relayList, sourceRouteTable[index].destination)
	}
	exist = true
	return
}

func NcpSetSourceRoute(id uint16) (err error) {
	exist, relayList := ncpFindSourceRoute(id)
	if !exist {
		ncpSourceRouteTrace("NCP cannot find source route for 0x%04x, send directly", id)
		return nil //不存在没有错，直接发送
	}
	ncpSourceRouteTrace("NCP set source route for 0x%04x, %v, \nTranceiverStep=%d SendStep=%d", id, relayList, ash.TransceiverStep, SendStep)
	err = EzspSetSourceRoute(id, relayList)
	return
}
