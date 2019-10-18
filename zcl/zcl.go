package zcl

import (
	"encoding/binary"
	"errors"

	"github.com/conthing/utils/common"
)

const (
	CmdReadAttrib                  byte = 0x00
	CmdReadAttribResponse          byte = 0x01
	CmdWriteAttrib                 byte = 0x02
	CmdWriteAttribUndivided        byte = 0x03
	CmdWriteAttribResponse         byte = 0x04
	CmdWriteAttribNoResponse       byte = 0x05
	CmdConfigureReporting          byte = 0x06
	CmdConfigureReportingResponse  byte = 0x07
	CmdReadReportingConfig         byte = 0x08
	CmdReadReportingConfigResponse byte = 0x09
	CmdReportAttrib                byte = 0x0a
	CmdDefaultResponse             byte = 0x0b
)

const (
	AttrZCLVersion         uint16 = 0x0000
	AttrApplicationVersion uint16 = 0x0001
	AttrStackVersion       uint16 = 0x0002
	AttrHWVersion          uint16 = 0x0003
	AttrManufacturerName   uint16 = 0x0004
	AttrModelIdentifier    uint16 = 0x0005
	AttrDateCode           uint16 = 0x0006
	AttrPowerSource        uint16 = 0x0007
)

const (
	UnsupportClusterCommand byte = 0x81
	UnsupportGeneralCommand byte = 0x82
	Success                 byte = 0x00
	Failure                 byte = 0x01
)

const (
	TypeNull           byte = 0x00 ///< Null data type
	Type8Bit           byte = 0x08 ///< 8-bit value data type
	Type16Bit          byte = 0x09 ///< 16-bit value data type
	Type24Bit          byte = 0x0a ///< 24-bit value data type
	Type32Bit          byte = 0x0b ///< 32-bit value data type
	Type40Bit          byte = 0x0c ///< 40-bit value data type
	Type48Bit          byte = 0x0d ///< 48-bit value data type
	Type56Bit          byte = 0x0e ///< 56-bit value data type
	Type64Bit          byte = 0x0f ///< 64-bit value data type
	TypeBool           byte = 0x10 ///< Boolean data type
	Type8BitMap        byte = 0x18 ///< 8-bit bitmap data type
	Type16BitMap       byte = 0x19 ///< 16-bit bitmap data type
	Type24BitMap       byte = 0x1a ///< 24-bit bitmap data type
	Type32BitMap       byte = 0x1b ///< 32-bit bitmap data type
	Type40BitMap       byte = 0x1c ///< 40-bit bitmap data type
	Type48BitMap       byte = 0x1d ///< 48-bit bitmap data type
	Type56BitMap       byte = 0x1e ///< 56-bit bitmap data type
	Type64BitMap       byte = 0x1f ///< 64-bit bitmap data type
	TypeU8             byte = 0x20 ///< Unsigned 8-bit value data type
	TypeU16            byte = 0x21 ///< Unsigned 16-bit value data type
	TypeU24            byte = 0x22 ///< Unsigned 16-bit value data type
	TypeU32            byte = 0x23 ///< Unsigned 32-bit value data type
	TypeU40            byte = 0x24 ///< Unsigned 40-bit value data type
	TypeU48            byte = 0x25 ///< Unsigned 48-bit value data type
	TypeU56            byte = 0x26 ///< Unsigned 56-bit value data type
	TypeU64            byte = 0x27 ///< Unsigned 64-bit value data type
	TypeS8             byte = 0x28 ///< signed 8-bit value data type
	TypeS16            byte = 0x29 ///< signed 16-bit value data type
	TypeS24            byte = 0x2a ///< signed 16-bit value data type
	TypeS32            byte = 0x2b ///< signed 32-bit value data type
	TypeS40            byte = 0x2c ///< signed 40-bit value data type
	TypeS48            byte = 0x2d ///< signed 48-bit value data type
	TypeS56            byte = 0x2e ///< signed 56-bit value data type
	TypeS64            byte = 0x2f ///< signed 64-bit value data type
	TypeEnum8          byte = 0x30 ///
	TypeEnum16         byte = 0x31 ///
	TypeByteArray      byte = 0x41 ///< Byte array data type
	TypeCharString     byte = 0x42 ///< Charactery string (array) data type
	TypeLongByteArray  byte = 0x43 ///< Long Byte array data type
	TypeLongCharString byte = 0x44 ///< Long Charactery string (array) data type
	TypeTimeOfDay      byte = 0xe0 ///
	TypeDate           byte = 0xe1 ///
	TypeUtcTime        byte = 0xe2 ///
	TypeClustId        byte = 0xe8 ///
	TypeAttribId       byte = 0xe9 ///
	TypeIeeeAddr       byte = 0xf0 ///< IEEE address (U64) type
	TypeSecKey         byte = 0xf1 ///
	TypeInvalid        byte = 0xff ///< Invalid data type
)

const (
	ClustBasic      uint16 = 0x0000 ///< Basic cluster ID
	ClustPwrCfg     uint16 = 0x0001 ///< Power configuration cluster ID
	ClustOnOff      uint16 = 0x0006 ///< On/Off cluster ID
	ClustOnOffSwcfg uint16 = 0x0007 ///< On/Off cluster ID
	ClustLevel      uint16 = 0x0008 ///< Level Control cluster ID
)

const (
	ClustOnOffCmdOff    byte = 0x00
	ClustOnOffCmdOn     byte = 0x01
	ClustOnOffCmdToggle byte = 0x02
)

const (
	ClustBasicCmdResetToDefault byte = 0x00
)

const (
	ClustLevelCmdMoveToLevel          byte = 0x00
	ClustLevelCmdMove                 byte = 0x01
	ClustLevelCmdStep                 byte = 0x02
	ClustLevelCmdStop                 byte = 0x03
	ClustLevelCmdMoveToLevelWithOnoff byte = 0x04
	ClustLevelCmdMoveWithOnoff        byte = 0x05
	ClustLevelCmdStepWithOnoff        byte = 0x06
	ClustLevelCmdStop2                byte = 0x07
)

type ZclContext struct {
	LocalEdp  byte //todo 改成local 和 remote
	RemoteEdp byte
	Context   interface{}

	GlobalHandle ZclGlobalHandle

	OnoffClusterHandle ZclOnoffClusterHandle
	BasicClusterHandle ZclBasicClusterHandle
	LevelClusterHandle ZclLevelClusterHandle
}

var SendSequence byte

type ZclGlobalHandle interface {
	AttribReportedHandle(*ZclContext, uint16, []*StAttrib) error
	UnsupportClusterCommandHandle(*ZclContext, uint16, bool, bool, byte, byte, interface{}) ([]byte, error)
}

type ZclBasicClusterHandle interface {
	CommandResetTODefaultHandle(*ZclContext)
}

type ZclOnoffClusterHandle interface {
	CommandOnHandle(*ZclContext)
	CommandOffHandle(*ZclContext)
	CommandToggleHandle(*ZclContext)
}

type ZclLevelClusterHandle interface {
	CommandMoveToLevelHandle(*ZclContext, uint8, uint16)
	CommandMoveHandle(*ZclContext, uint8, uint8)
	CommandStepHandle(*ZclContext, uint8, uint8, uint16)
	CommandStopHandle(*ZclContext)
	CommandStepWithOnoffHandle(*ZclContext, uint8, uint8, uint16)
	CommandMoveWithOnoffHandle(*ZclContext, uint8, uint8)
	CommandMoveToLevelWithOnoffHandle(*ZclContext, uint8, uint16)
}

type StAttrib struct {
	AttributeIdentifier uint16
	AttributeDataType   byte
	AttributeData       interface{}
}

//todo error的详细定义
var ErrUnsupportProfile = errors.New("ErrUnsupportProfile")
var ErrManufCodeNotSupport = errors.New("ErrManufCodeNotSupport")
var ErrFrameTypeNotSupport = errors.New("ErrFrameTypeNotSupport")
var ErrOtherClusterNotSupport = errors.New("ErrOtherClusterNotSupport")
var ErrUnsupportCommand = errors.New("ErrUnsupportCommand")
var ErrorDataTypeNotSupport = errors.New("ErrorDataTypeNotSupport")
var ErrFailToAnalysis = errors.New("ErrFailToAnalysis")
var ErrUnsupportDirection = errors.New("ErrUnsupportDirection")
var ErrUnsupportClusterCommand = errors.New("ErrUnsupportClusterCommand")
var ErrUnsupportGeneralCommand = errors.New("ErrUnsupportGeneralCommand")

func zclPackFrameCtrl(clusterSpecific bool, fromServer bool, disableDefaultResponse bool, manufSpecific bool, manufCode uint16) (data []byte) {
	//frameType为00时 发送command:reportattr,readattr,writeattr,defaultresponse
	//frameType为01时 发送具体的cluster
	//manufSpecific为0
	//direction 回复时和来的时候相反
	//disableDefaultResponse 回复时为true
	data = make([]byte, 1)

	if clusterSpecific {
		data[0] = 1
	}
	if fromServer {
		data[0] |= 8
	}
	if disableDefaultResponse {
		data[0] |= 16
	}
	if manufSpecific {
		data[0] |= 4
		manufCode1 := byte(manufCode >> 8)
		manufCode0 := byte(manufCode & 0xff)
		data = append(data, manufCode0, manufCode1)
	}
	return
}

func zclPackFrame(clusterSpecific bool, fromServer bool, disableDefaultResponse bool, sequenceNumber byte, commandIdentifier byte, payload []byte) (data []byte) {
	data = zclPackFrameCtrl(clusterSpecific, fromServer, disableDefaultResponse, false, 0)
	data = append(data, sequenceNumber, commandIdentifier)
	if payload != nil {
		data = append(data, payload...)
	}
	return
}

//ZclPackOnoffClusterCommandOff pack a OnOff cluster OFF command
func ZclPackOnoffClusterCommandOff() (data []byte) {
	data = zclPackFrame(true, false, false, SendSequence, 0, nil)
	SendSequence++
	return
}

//ZclPackOnoffClusterCommandOn pack a OnOff cluster ON command
func ZclPackOnoffClusterCommandOn() (data []byte) {
	data = zclPackFrame(true, false, false, SendSequence, 1, nil)
	SendSequence++
	return
}

//ZclPackOnoffClusterCommandToggle pack a OnOff cluster TOGGLE command
func ZclPackOnoffClusterCommandToggle() (data []byte) {
	data = zclPackFrame(true, false, false, SendSequence, 2, nil)
	SendSequence++
	return
}

//ZclPackBasicClusterCommandResetToFactory pack a Basic cluster Reset to factory command
func ZclPackBasicClusterCommandResetToFactory() (data []byte) {
	data = zclPackFrame(true, false, false, SendSequence, 2, nil)
	SendSequence++
	return
}

//ZclPackLevelCommandMoveToLevel pack a Level cluster Move to level command
func ZclPackLevelCommandMoveToLevel(level uint8, time uint16) (data []byte) {
	payload := make([]byte, 3)
	payload[0] = level
	payload[1] = byte(time & 0xff)
	payload[2] = byte((time >> 8) & 0xff)
	data = zclPackFrame(true, false, false, SendSequence, 0, payload)
	SendSequence++
	return
}

//ZclPackLevelCommandMove pack a Level cluster Move command
func ZclPackLevelCommandMove(mode uint8, rate uint8) (data []byte) {
	payload := make([]byte, 2)
	payload[0] = mode
	payload[1] = rate
	data = zclPackFrame(true, false, false, SendSequence, 1, payload)
	SendSequence++
	return
}

//ZclPackLevelCommandStep pack a Level cluster Step command
func ZclPackLevelCommandStep(mode uint8, size uint8, time uint16) (data []byte) {
	payload := make([]byte, 4)
	payload[0] = mode
	payload[1] = size
	payload[2] = byte(time & 0xff)
	payload[3] = byte((time >> 8) & 0xff)
	data = zclPackFrame(true, false, false, SendSequence, 2, payload)
	SendSequence++
	return
}

//ZclPackLevelCommandStop pack a Level cluster Stop command
func ZclPackLevelCommandStop() (data []byte) {
	data = zclPackFrame(true, false, false, SendSequence, 3, nil) //3或者7,怎么处理
	SendSequence++
	return
}

//ZclPackLevelCommandMoveToLevelWithOnOff pack a Level cluster Move to level with onoff command
func ZclPackLevelCommandMoveToLevelWithOnOff(level uint8, time uint16) (data []byte) {
	payload := make([]byte, 3)
	payload[0] = level
	payload[1] = byte(time & 0xff)
	payload[2] = byte((time >> 8) & 0xff)
	data = zclPackFrame(true, false, false, SendSequence, 4, payload)
	SendSequence++
	return
}

//ZclPackLevelCommandMoveWithOnOff pack a Level cluster Move with onoff command
func ZclPackLevelCommandMoveWithOnOff(mode uint8, rate uint8) (data []byte) {
	payload := make([]byte, 2)
	payload[0] = mode
	payload[1] = rate
	data = zclPackFrame(true, false, false, SendSequence, 5, payload)
	SendSequence++
	return
}

//ZclPackLevelCommandStepWithOnOff pack a Level cluster Step with onoff command
func ZclPackLevelCommandStepWithOnOff(mode uint8, size uint8, time uint16) (data []byte) {
	payload := make([]byte, 4)
	payload[0] = mode
	payload[1] = size
	payload[2] = byte(time & 0xff)
	data = zclPackFrame(true, false, false, SendSequence, 6, nil)
	SendSequence++
	return
}

//ZclPackReadAttr means ReadAttribute
func ZclPackReadAttr(attrId []byte) (data []byte) {
	//todo zjn 这两个函数写的都有问题
	payload := make([]byte, 1)
	payload = append(payload, attrId...)
	data = zclPackFrame(false, false, false, SendSequence, 0, payload)
	SendSequence++
	return
}

//ZclPackReadAttr means WriteAttribute
func ZclPackWriteAttr(attrId []byte, attrDataType byte, attrData []byte) (data []byte) {
	payload := make([]byte, 1)
	//payload = append(payload,attrId...,attrDataType,attrData...)
	data = zclPackFrame(false, false, false, SendSequence, 2, payload)
	SendSequence++
	return
}

//Parse means parse the frametype data
func (z *ZclContext) Parse(profile uint16, cluster uint16, data []byte) (resp []byte, err error) {
	//if profile != 0x104 {
	//	common.Log.Errorf("ErrUnsupportProfile")
	//	return nil, ErrUnsupportProfile
	//}

	frameCtrl := data[0]
	frameType := frameCtrl & 0x03
	manufSpecific := (frameCtrl & 0x04) != 0
	direction := (frameCtrl & 0x08) != 0
	disableDefaultResponse := (frameCtrl & 0x10) != 0

	if manufSpecific {
		common.Log.Errorf("ErrManufCodeNotSupport")
		return nil, ErrManufCodeNotSupport
	}

	sequenceNumber := data[1]
	commandIdentifier := data[2]

	//common.Log.Debugf("FrameControl:%02x SequenceNumber:%02x CommandIdentifier:%02x Payload: 0x%x\n",
	//	frameCtrl, sequenceNumber, commandIdentifier, data[3:])

	if frameType == 0x00 {
		resp, err = z.ParseGlobalCommand(cluster, direction, disableDefaultResponse, sequenceNumber, commandIdentifier, data[3:])
	} else if frameType == 0x01 {
		resp, err = z.ParseClusterCommand(cluster, direction, disableDefaultResponse, sequenceNumber, commandIdentifier, data[3:])
	} else {
		common.Log.Errorf("ErrFrameTypeNotSupport")
		err = ErrFrameTypeNotSupport
	}

	if resp == nil && commandIdentifier != 0x0b {
		var payload []byte
		payload = make([]byte, 2)
		payload[0] = commandIdentifier
		if err != nil || disableDefaultResponse == false {

			if err == ErrUnsupportClusterCommand {
				payload[1] = UnsupportClusterCommand
			} else if err == ErrUnsupportGeneralCommand {
				payload[1] = UnsupportGeneralCommand
			} else if err != nil {
				payload[1] = Failure
			} else {
				payload[1] = Success
			}
			resp = zclPackFrame(false, false, true, sequenceNumber, CmdDefaultResponse, payload)
		}
	}

	return
}

func getSendDirection(direction bool) byte {
	var sendDirection byte
	if direction {
		sendDirection = 0x00
	} else {
		sendDirection = 0x08
	}
	return sendDirection
}

func (z *ZclContext) attribsReported(cluster uint16, list []*StAttrib) {
	if z.GlobalHandle != nil {
		err := z.GlobalHandle.AttribReportedHandle(z, cluster, list)
		if err != nil {
			common.Log.Errorf("Attrib report parse error %v", err)
		}
	}
}

func (z *ZclContext) getReportAttribs(payload []byte) ([]*StAttrib, error) {
	var stAttrib []*StAttrib

	for i := 0; i+3 <= len(payload); {
		attributeData, attribLen, err := z.getPayload(payload[i+2], payload[i+3:])
		if err != nil {
			common.Log.Errorf("zcl get payload err %v", err)
		}
		stAttribNew := &StAttrib{AttributeIdentifier: binary.LittleEndian.Uint16(payload[i : i+2]), AttributeDataType: payload[i+2], AttributeData: attributeData}
		stAttrib = append(stAttrib, stAttribNew)
		i += (attribLen + 3)
	}
	return stAttrib, nil
}

func (z *ZclContext) getPayload(datatype byte, payload []byte) (val interface{}, length int, err error) {
	var dataPayload []byte
	switch datatype {
	case Type8Bit, Type8BitMap, TypeU8, TypeEnum8:
		length = 1
		if len(payload) < length {
			return nil, length, err
		}
		val = payload[0]
	case TypeBool:
		length = 1
		if len(payload) < length {
			return nil, length, err
		}
		val = (payload[0] != 0)
	case TypeS8:
		length = 1
		if len(payload) < length {
			return nil, length, err
		}
		val = int8(payload[0])
	case Type16Bit, Type16BitMap, TypeU16, TypeEnum16, TypeClustId, TypeAttribId:
		length = 2
		if len(payload) < length {
			return nil, length, err
		}
		val = binary.LittleEndian.Uint16(payload[0:2])
	case TypeS16:
		length = 2
		if len(payload) < length {
			return nil, length, err
		}
		val = int16(binary.BigEndian.Uint16(payload[0:2]))
	case Type32Bit, Type32BitMap, TypeU32, TypeTimeOfDay, TypeDate, TypeUtcTime:
		length = 4
		if len(payload) < length {
			return nil, length, err
		}
		val = binary.LittleEndian.Uint32(payload[0:4])
	case TypeS32:
		length = 4
		if len(payload) < length {
			return nil, length, err
		}
		val = int32(binary.BigEndian.Uint32(payload[0:4]))
	case TypeByteArray:
		payloadlen := payload[0]
		length = 1 + int(payloadlen)
		if len(payload) < length {
			return nil, length, err
		}
		val = payload[1 : 1+payloadlen]
	case TypeCharString:
		payloadlen := payload[0]
		length = 1 + int(payloadlen)
		if len(payload) < length {
			return nil, length, err
		}
		dataPayload = payload[1 : 1+payloadlen]
		val = string(dataPayload)
	case Type64Bit, Type64BitMap, TypeU64, TypeIeeeAddr:
		length = 8
		if len(payload) < length {
			return nil, length, err
		}
		val = binary.LittleEndian.Uint64(payload[0:8])
	case TypeS64:
		length = 8
		if len(payload) < length {
			return nil, length, err
		}
		val = int64(binary.BigEndian.Uint64(payload[0:8]))
	case TypeSecKey:
		length = 16
		if len(payload) < length {
			return nil, length, err
		}
		val = payload[0:16]
	case TypeLongByteArray:
		payloadlength := binary.LittleEndian.Uint16(payload[0:2])
		length = 2 + int(payloadlength)
		if len(payload) < length {
			return nil, length, err
		}
		val = payload[2 : 2+payloadlength]
	case TypeLongCharString:
		payloadlength := binary.LittleEndian.Uint16(payload[0:2])
		length = 2 + int(payloadlength)
		if len(payload) < length {
			return nil, length, err
		}
		val = string(payload[2 : 2+payloadlength])
	case Type24Bit, Type24BitMap, TypeU24:
		length = 3
		if len(payload) < length {
			return nil, length, err
		}
		u24 := uint32(payload[0]) | uint32(payload[1])<<8 | uint32(payload[2])<<16
		val = u24
	case TypeS24:
		length = 3
		if len(payload) < length {
			return nil, length, err
		}
		s24 := uint32(payload[0]) | uint32(payload[1])<<8 | uint32(payload[2])<<16
		if s24 >= 0x800000 {
			s24 |= 0xff000000
		}
		val = int32(s24)
	case Type40Bit, Type40BitMap, TypeU40:
		length = 5
		if len(payload) < length {
			return nil, length, err
		}
		u40 := uint64(payload[0]) | uint64(payload[1])<<8 | uint64(payload[2])<<16 | uint64(payload[3])<<24 | uint64(payload[4])<<32
		val = u40
	case TypeS40:
		length = 5
		if len(payload) < length {
			return nil, length, err
		}
		s40 := uint64(payload[0]) | uint64(payload[1])<<8 | uint64(payload[2])<<16 | uint64(payload[3])<<24 | uint64(payload[4])<<32
		if s40 >= 0x8000000000 {
			s40 |= 0xffffff0000000000
		}
		val = int64(s40)
	case Type48Bit, Type48BitMap, TypeU48:
		length = 6
		if len(payload) < length {
			return nil, length, err
		}
		u48 := uint64(payload[0]) | uint64(payload[1])<<8 | uint64(payload[2])<<16 | uint64(payload[3])<<24 | uint64(payload[4])<<32 | uint64(payload[5])<<40
		val = u48
	case TypeS48:
		length = 6
		if len(payload) < length {
			return nil, length, err
		}
		s48 := uint64(payload[0]) | uint64(payload[1])<<8 | uint64(payload[2])<<16 | uint64(payload[3])<<24 | uint64(payload[4])<<32 | uint64(payload[5])<<40
		if s48 >= 0x800000000000 {
			s48 |= 0xffff000000000000
		}
		val = int64(s48)
	case Type56Bit, Type56BitMap, TypeU56:
		length = 7
		if len(payload) < length {
			return nil, length, err
		}
		u56 := uint64(payload[0]) | uint64(payload[1])<<8 | uint64(payload[2])<<16 | uint64(payload[3])<<24 | uint64(payload[4])<<32 | uint64(payload[5])<<40 | uint64(payload[6])<<48
		val = u56
	case TypeS56:
		length = 7
		if len(payload) < length {
			return nil, length, err
		}
		s56 := uint64(payload[0]) | uint64(payload[1])<<8 | uint64(payload[2])<<16 | uint64(payload[3])<<24 | uint64(payload[4])<<32 | uint64(payload[5])<<40 | uint64(payload[6])<<48
		if s56 >= 0x80000000000000 {
			s56 |= 0xff00000000000000
		}
		val = int64(s56)
	case TypeNull, TypeInvalid:
		length = 0
		val = nil
	default:
		common.Log.Errorf("ErrorDataTypeNotSupport")
		return nil, 0, ErrorDataTypeNotSupport
	}
	return
}

//ParseGlobalCommand means choose every different global command
func (z *ZclContext) ParseGlobalCommand(cluster uint16, direction bool, disableDefaultResponse bool, sequenceNumber byte,
	commandIdentifier byte, data []byte) ([]byte, error) {
	if direction {
		switch commandIdentifier {
		case CmdReadAttribResponse:
			common.Log.Errorf("Read Attribute not support")
		case CmdWriteAttribResponse:
			common.Log.Errorf("Write Attribute not support")
		case CmdReportAttrib:
			//common.Log.Debug("Attributes Reported")
			attribs, err := z.getReportAttribs(data)
			z.attribsReported(cluster, attribs)
			if err != nil {
				return nil, ErrFailToAnalysis
			}
			//z.reportAttrib(cluster, data[3:])
		case CmdDefaultResponse:

		default:
			common.Log.Errorf("ErrUnsupportGeneralCommand")
			return nil, ErrUnsupportGeneralCommand
		}
	} else {
		common.Log.Errorf("ErrUnsupportGeneralCommand")
		return nil, ErrUnsupportGeneralCommand
	}

	return nil, nil
}

//ParseClusterCommand means choose every different cluster
func (z *ZclContext) ParseClusterCommand(cluster uint16, direction bool, disableDefaultResponse bool, sequenceNumber byte,
	commandIdentifier byte, data []byte) (resp []byte, err error) {
	switch cluster {
	case ClustBasic:
		resp, err = z.BasicClusterCommandHandle(cluster, direction, disableDefaultResponse, sequenceNumber, commandIdentifier, data)
	case ClustOnOff:
		resp, err = z.OnoffClusterCommandHandle(cluster, direction, disableDefaultResponse, sequenceNumber, commandIdentifier, data)
	case ClustLevel:
		resp, err = z.LevelClusterCommandHandle(cluster, direction, disableDefaultResponse, sequenceNumber, commandIdentifier, data)
	default:
		resp, err = nil, ErrUnsupportClusterCommand
	}

	if err == ErrUnsupportClusterCommand {
		if z.GlobalHandle != nil {
			resp, err = z.GlobalHandle.UnsupportClusterCommandHandle(z, cluster, direction, disableDefaultResponse, sequenceNumber, commandIdentifier, data)
		}
	}

	return
}

//OnoffClusterCommandHandle means choose OnoffCluster
func (z *ZclContext) OnoffClusterCommandHandle(cluster uint16, direction bool, disableDefaultResponse bool, sequenceNumber byte,
	commandIdentifier byte, data []byte) ([]byte, error) {
	if direction {
		switch commandIdentifier {
		case ClustOnOffCmdOff:
			if z.OnoffClusterHandle != nil {
				z.OnoffClusterHandle.CommandOffHandle(z)
			}
		case ClustOnOffCmdOn:
			if z.OnoffClusterHandle != nil {
				z.OnoffClusterHandle.CommandOnHandle(z)
			}
		case ClustOnOffCmdToggle:
			if z.OnoffClusterHandle != nil {
				z.OnoffClusterHandle.CommandToggleHandle(z)
			}
		default:
			return nil, ErrUnsupportCommand
		}
	} else {
		return nil, ErrUnsupportDirection
	}

	return nil, nil
}

//BasicClusterCommandHandle means choose BasicCluster
func (z *ZclContext) BasicClusterCommandHandle(cluster uint16, direction bool, disableDefaultResponse bool, sequenceNumber byte,
	commandIdentifier byte, data []byte) ([]byte, error) {
	if direction {
		switch commandIdentifier {
		case ClustBasicCmdResetToDefault:
			if z.BasicClusterHandle != nil {
				z.BasicClusterHandle.CommandResetTODefaultHandle(z)
			}
		default:
			return nil, ErrUnsupportCommand
		}
	} else {
		return nil, ErrUnsupportDirection
	}

	return nil, nil
}

//LevelClusterCommandle means choose levelCluster
func (z *ZclContext) LevelClusterCommandHandle(cluster uint16, direction bool, disableDefaultResponse bool, sequenceNumber byte,
	commandIdentifier byte, data []byte) ([]byte, error) {
	var level uint8
	var time uint16
	var mode uint8
	var size uint8
	var rate uint8
	if direction {
		switch commandIdentifier {
		case ClustLevelCmdMoveToLevel:
			if z.LevelClusterHandle != nil {
				level = data[0]
				time = binary.LittleEndian.Uint16(data[1:3])
				z.LevelClusterHandle.CommandMoveToLevelHandle(z, level, time)
			}
		case ClustLevelCmdMove:
			if z.LevelClusterHandle != nil {
				mode = data[0]
				rate = data[1]
				z.LevelClusterHandle.CommandMoveHandle(z, mode, rate)
			}
		case ClustLevelCmdStep:
			if z.LevelClusterHandle != nil {
				mode = data[0]
				size = data[1]
				time = binary.LittleEndian.Uint16(data[2:4])
				z.LevelClusterHandle.CommandStepHandle(z, mode, size, time)
			}
		case ClustLevelCmdStop, ClustLevelCmdStop2:
			if z.LevelClusterHandle != nil {
				z.LevelClusterHandle.CommandStopHandle(z)
			}
		case ClustLevelCmdMoveToLevelWithOnoff:
			if z.LevelClusterHandle != nil {
				level = data[0]
				time = binary.LittleEndian.Uint16(data[1:3])
				z.LevelClusterHandle.CommandMoveToLevelWithOnoffHandle(z, level, time)
			}
		case ClustLevelCmdMoveWithOnoff:
			if z.LevelClusterHandle != nil {
				mode = data[0]
				rate = data[1]
				z.LevelClusterHandle.CommandMoveWithOnoffHandle(z, mode, rate)
			}
		case ClustLevelCmdStepWithOnoff:
			if z.LevelClusterHandle != nil {
				mode = data[0]
				size = data[1]
				time = binary.LittleEndian.Uint16(data[2:4])
				z.LevelClusterHandle.CommandStepWithOnoffHandle(z, mode, size, time)
			}
		default:
			return nil, ErrUnsupportCommand
		}
	} else {
		return nil, ErrUnsupportDirection
	}

	return nil, nil
}
