/*
 * Copyright (C) 2016 Tim Mathews <tim@signalk.org>
 *
 * This file is part of Argo.
 *
 * Argo is free software: you can redistribute it and/or modify it under the
 * terms of the GNU General Public License as published by the Free Software
 * Foundation, either version 3 of the License, or (at your option) any later
 * version.
 *
 * Argo is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program. If not, see <http://www.gnu.org/licenses/>.
 */

package actisense

import (
	"fmt"
	"math"
)

const (
	// N2K commands
	N2kMsgRecv = 0x93
	N2kMsgSend = 0x94

	// NGT commands
	ACmdRecv = 0xA0
	ACmdSend = 0xA1
)

// The following commands are used with ACmdRecv and ACmdSend

const (
	ACmdReInitMainApp byte = iota
	ACmdCommitToEEPROM
	ACmdCommitToFlash
)
const (
	ACmdHardwareInfo byte = iota + 0x10
	ACmdOperatingMode
	ACmdPortBaudCfg
	ACmdPortPCodeCfg
	ACmdPortDupDelete
	ACmdTotalTime
	ACmdHardwareBaud
)
const (
	ACmdSupportedPGNList byte = iota + 0x40
	ACmdProductInfoN2K
	ACmdCANConfig
	ACmdCANInfoField1
	ACmdCANInfoField2
	ACmdCANInfoField3
	ACmdRxPGNEnable
	ACmdTxPGNEnable
	ACmdRxPGNEnableList
	ACmdTxPGNEnableList
	ACmdDeletePGNEnableList
	ACmdActivatePGNEnableLists
	ACmdDefaultPGNEnableList
	ACmdParamsPGNEnableLists
	ACmdRxPGNEnableListF2
	ACmdTxPGNEnableListF2
)
const (
	ACmdStartupStatus byte = iota + 0xF0
	ACmdErrorReport
	ACmdSystemStatus
)

const ACmdNegativeAcknowledge byte = 0xF4

type ActisenseMode uint16

const (
	OpModeFilter ActisenseMode = iota + 1
	OpModeRxAll
)

const (
	RxPGNList byte = iota
	TxPGNList
)

func (p *ActisensePort) Reboot() (int, error) {
	return p.write(0x01, 0x11)
}

func (p *ActisensePort) ReInitMainApp() (int, error) {
	return p.write(ACmdSend, ACmdReInitMainApp)
}

func (p *ActisensePort) CommitToEEPROM() (int, error) {
	return p.write(ACmdSend, ACmdCommitToEEPROM)
}

func (p *ActisensePort) CommitToFlash() (int, error) {
	return p.write(ACmdSend, ACmdCommitToFlash)
}

func (p *ActisensePort) GetHardwareInfo() (int, error) {
	return p.write(ACmdSend, ACmdHardwareInfo)
}

func (p *ActisensePort) GetOperatingMode() (int, error) {
	return p.write(ACmdSend, ACmdOperatingMode)
}

func (p *ActisensePort) SetOperatingMode(mode ActisenseMode) (int, error) {
	return p.write(ACmdSend, ACmdOperatingMode, byte(mode), byte(mode>>8))
}

func (p *ActisensePort) GetPortBaudCodes() (int, error) {
	return p.write(ACmdSend, ACmdPortBaudCfg)
}

func (p *ActisensePort) SetPortBaudCodes() (int, error) {
	return p.write(ACmdSend, ACmdPortBaudCfg)
}

func (p *ActisensePort) GetPortPCodes() (int, error) {
	return p.write(ACmdSend, ACmdPortPCodeCfg)
}

func (p *ActisensePort) SetPortPCodes(pcodes ...byte) (int, error) {
	d := []byte{ACmdPortPCodeCfg, byte(len(pcodes) / 2)}
	d = append(d, pcodes...)
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) GetPortDupDelete() (int, error) {
	return p.write(ACmdSend, ACmdPortDupDelete)
}

func (p *ActisensePort) SetPortDupDelete() (int, error) {
	return p.write(ACmdSend, ACmdPortDupDelete)
}

func (p *ActisensePort) GetTotalTime() (int, error) {
	return p.write(ACmdSend, ACmdTotalTime)
}

func (p *ActisensePort) SetTotalTime() (int, error) {
	return p.write(ACmdSend, ACmdTotalTime)
}

func (p *ActisensePort) GetHardwareBaudCodes() (int, error) {
	return p.write(ACmdSend, ACmdHardwareBaud)
}

func (p *ActisensePort) SetHardwareBaudCodes() (int, error) {
	return p.write(ACmdSend, ACmdHardwareBaud)
}

func (p *ActisensePort) GetSupportedPGNList() (int, error) {
	return p.write(ACmdSend, ACmdSupportedPGNList)
}

func (p *ActisensePort) GetProductInfoN2K() (int, error) {
	return p.write(ACmdSend, ACmdProductInfoN2K)
}

func (p *ActisensePort) GetCANConfig() (int, error) {
	return p.write(ACmdSend, ACmdCANConfig)
}

func (p *ActisensePort) SetCANConfig() (int, error) {
	return p.write(ACmdSend, ACmdCANConfig)
}

func (p *ActisensePort) GetCANInfoField1() (int, error) {
	return p.write(ACmdSend, ACmdCANInfoField1)
}

func (p *ActisensePort) SetCANInfoField1(info string) (int, error) {
	d := []byte{ACmdCANInfoField1}
	d = append(d, []byte(truncateOrPad(info, 32))...)
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) GetCANInfoField2() (int, error) {
	return p.write(ACmdSend, ACmdCANInfoField2)
}

func (p *ActisensePort) SetCANInfoField2(info string) (int, error) {
	d := []byte{ACmdCANInfoField2}
	d = append(d, []byte(truncateOrPad(info, 32))...)
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) GetCANInfoField3() (int, error) {
	return p.write(ACmdSend, ACmdCANInfoField3)
}

func (p *ActisensePort) GetRxPGN(pgn int) (int, error) {
	d := []byte{ACmdRxPGNEnable}
	d = append(d, intToBytes(pgn)...)
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) SetRxPGN(pgn, enable int) (int, error) {
	d := []byte{ACmdRxPGNEnable}
	d = append(d, intToBytes(pgn)...)
	d = append(d, byte(enable))
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) SetRxPGNEx(pgn, enable, mask int) (int, error) {
	return -1, nil
}

func (p *ActisensePort) GetTxPGN(pgn int) (int, error) {
	d := []byte{ACmdTxPGNEnable}
	d = append(d, intToBytes(pgn)...)
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) SetTxPGN(pgn, enable int) (int, error) {
	d := []byte{ACmdTxPGNEnable}
	d = append(d, intToBytes(pgn)...)
	d = append(d, byte(enable))
	return p.write(ACmdSend, d...)
}

func (p *ActisensePort) SetTxPGNEx(pgn, enable, mask int) (int, error) {
	return -1, nil
}

func (p *ActisensePort) GetRxPGNList() (int, error) {
	return p.write(ACmdSend, ACmdRxPGNEnableList)
}

func (p *ActisensePort) GetTxPGNList() (int, error) {
	return p.write(ACmdSend, ACmdTxPGNEnableList)
}

func (p *ActisensePort) ClearRxPGNList() (int, error) {
	return p.write(ACmdSend, ACmdDeletePGNEnableList, RxPGNList)
}

func (p *ActisensePort) ClearTxPGNList() (int, error) {
	return p.write(ACmdSend, ACmdDeletePGNEnableList, TxPGNList)
}

func (p *ActisensePort) ActivatePGNEnableLists() (int, error) {
	return p.write(ACmdSend, ACmdActivatePGNEnableLists)
}

func (p *ActisensePort) SetDefaultPGNEnableList(id byte) (int, error) {
	return p.write(ACmdSend, ACmdDefaultPGNEnableList, id)
}

func (p *ActisensePort) GetParamsPGNEnableLists() (int, error) {
	return p.write(ACmdSend, ACmdParamsPGNEnableLists)
}

func intMin(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func intMax(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func truncateOrPad(s string, l int) (o string) {
	o = s[0:intMin(len(s), l)]

	if len(o) < l {
		f := fmt.Sprintf("%%-%vv", l)
		o = fmt.Sprintf(f, o)
	}

	return
}

func intToBytes(i int) []byte {
	return []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
}
