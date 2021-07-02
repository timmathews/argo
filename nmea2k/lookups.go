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

package nmea2k

type PgnLookup map[int]string
type PgnSubLookup map[int]PgnLookup

var lookupActisenseCANBaudCode = PgnLookup{
	0: "10,000",
	1: "25,000",
	2: "50,000",
	3: "100,000",
	4: "125,000",
	5: "250,000",
	6: "500,000",
	7: "1,000,000",
}

var lookupActisenseUARTBaudCode = PgnLookup{
	0:  "1,200",
	1:  "2,400",
	2:  "4,800",
	3:  "9,600",
	4:  "19,200",
	5:  "38,400",
	6:  "57,600",
	7:  "115,200",
	8:  "230,400",
	9:  "460,800",
	10: "921,600",
	11: "Undefined",
	12: "Undefined",
	13: "31,250",
	14: "62,500",
	15: "125,000",
	16: "250,000",
	17: "500,000",
	18: "1,000,000",
	19: "1,000,000",
	20: "Undefined",
	21: "Undefined",
}

var lookupActisensePCode = PgnLookup{
	0: "Disabled",
	1: "Enabled",
	2: "Enabled for session",
}

var lookupAlertType = PgnLookup{
	0:  "Reserved",
	1:  "Energy Alarm",
	2:  "Alarm",
	3:  "Reserved",
	4:  "Reserved",
	5:  "Warning",
	6:  "Reserved",
	7:  "Reserved",
	8:  "Caution",
	13: "Reserved",
	14: "Data out of range",
	15: "Data not available",
}

var lookupAlertCategory = PgnLookup{
	0:  "Navigational",
	1:  "Technical",
	13: "Reserved",
	14: "Data out of range",
	15: "Data not available",
}

var lookupSilenceStatus = PgnLookup{
	0: "Not temporary status",
	1: "Temporary status",
}

var lookupAcknowledgeStatus = PgnLookup{
	0: "Not acknowledged",
	1: "Acknowledged",
}

var lookupEscalationStatus = PgnLookup{
	0: "Not escalated",
	1: "Escalated",
}

var lookupSpeedDirection = PgnLookup{
	0: "Ahead",
	1: "Astern",
}

var lookupSpeedReference = PgnLookup{
	0:   "Paddle Wheel",
	1:   "Pitot Tube",
	2:   "Doppler Log",
	3:   "Correlation Log",
	4:   "Electromagnetic",
	253: "Not supported",
}

var lookupSupport = PgnLookup{
	0: "Not supported",
	1: "Supported",
}

var lookupTriggerCondition = PgnLookup{
	0:  "Manual",
	1:  "Auto",
	2:  "Test",
	13: "Reserved",
	14: "Data out of range",
	15: "Data not available",
}

var lookupThresholdStatus = PgnLookup{
	0:  "Normal",
	1:  "Threshold exceeded",
	2:  "Extreme threshold exceeded",
	3:  "Low threshold exceeded",
	4:  "Extreme low threshold exceeded",
	13: "Reserved",
	14: "Data out of range",
	15: "Data not available",
}

var lookupAlertState = PgnLookup{
	0:   "Disabled",
	1:   "Normal",
	2:   "Active",
	3:   "Silenced",
	4:   "Acknowleged",
	5:   "Awaiting acknowlege",
	253: "Reserved",
	254: "Data out of range",
	255: "Data not available",
}

var lookupAisAccuracy = PgnLookup{
	0: "Low",
	1: "High",
}

var lookupAisAtoNType = PgnLookup{
	0:  "Not specified",
	1:  "Reference point",
	2:  "RACON",
	3:  "Fix structure offshore",
	4:  "Reserved",
	5:  "Light, without sectors",
	6:  "Light, with sectors",
	7:  "Leading light, front",
	8:  "Leading light, rear",
	9:  "Beacon, cardinal N",
	10: "Beacon, cardinal E",
	11: "Beacon, cardinal S",
	12: "Beacon, cardinal W",
	13: "Beacon, port hand",
	14: "Beacon, starboard hand",
	15: "Beacon, preferred channel, port hand",
	16: "Beacon, preferred channel, starboard hand",
	17: "Beacon, isolated danger",
	18: "Beacon, safe water",
	19: "Beacon, special mark",
	20: "Cardinal mark N",
	21: "Cardinal mark E",
	22: "Cardinal mark S",
	23: "Cardinal mark W",
	24: "Port hand mark",
	25: "Starboard hand mark",
	26: "Preferred channel, port hand",
	27: "Preferred channel, starboard hand",
	28: "Isolated danger",
	29: "Safe water",
	30: "Special mark",
	31: "Light vessel/LANBY/Rigs",
}

var lookupAisBand = PgnLookup{
	0: "top 525 kHz of marine band",
	1: "entire marine band",
}

var lookupAisCommState = PgnLookup{
	0: "SOTDMA",
	1: "ITDMA",
}

var lookupAisDTE = PgnLookup{
	0: "Available",
	1: "Not available",
}

var lookupAisMode = PgnLookup{
	0: "Autonomous",
	1: "Assigned",
}

var lookupAisRAIM = PgnLookup{
	0: "not in use",
	1: "in use",
}

var lookupAisTransceiver = PgnLookup{
	0: "Channel A VDL reception",
	1: "Channel B VDL reception",
	2: "Channel A VDL transmission",
	3: "Channel B VDL transmission",
	4: "Own information not broadcast",
	5: "Reserved",
}

var lookupAisUnitType = PgnLookup{
	0: "SOTDMA",
	1: "CS",
}

var lookupAisVersion = PgnLookup{
	0: "ITU-R M.1371-1",
	1: "ITU-R M.1371-3",
}

// http://www.nmea.org/Assets/20120726%20nmea%202000%20class%20&%20function%20codes%20v%202.00.pdf
var lookupDeviceClass = PgnLookup{
	0:   "Reserved for NMEA 2000 Use",
	10:  "System Tools",
	20:  "Safety Systems",
	25:  "Inter/Intranetwork Device",
	30:  "Electrical Distribution",
	35:  "Electrical Generation",
	40:  "Steering and Control Surfaces",
	50:  "Propulsion",
	60:  "Navigation",
	70:  "Communication",
	75:  "Sensor Communication Interface",
	80:  "Instrumentation/General Systems (Deprecated)",
	85:  "External Environment",
	90:  "Internal Environment",
	100: "Deck, Cargo and Fishing Equipment",
	120: "Display",
	125: "Entertainment",
}

var lookupDeviceFunction = PgnSubLookup{
	0: PgnLookup{
		0: "Reserved for NMEA 2000 Use",
	},
	10: PgnLookup{
		130: "Diagnostic",
		140: "Bus Traffic Logger",
	},
	20: PgnLookup{
		110: "Alarm Enunciator (Deprecated)",
		130: "Emergency Position Indicating Radio Beacon (EPIRB)",
		135: "Man Overboard",
		140: "Voyage Data Recorder",
		150: "Camera",
	},
	25: PgnLookup{
		130: "PC Gateway",
		131: "NMEA 2000 to Analog Gateway",
		132: "Analog to NMEA 2000 Gateway",
		135: "NMEA 0183 Gateway",
		140: "Router",
		150: "Bridge",
		160: "Repeater",
	},
	30: PgnLookup{
		130: "Binary Event Monitor",
		140: "Load Controller",
		141: "AC/DC Input",
		150: "Function Controller",
	},
	35: PgnLookup{
		140: "Engine",
		141: "DC Generator/Alternator",
		142: "Solar Panel (Solar Array)", // Documentation lists this as 141 also
		143: "Wind Generator (DC)",
		144: "Fuel Cell",
		145: "Network Power Supply",
		151: "AC Generator",
		152: "AC Bus",
		153: "AC Mains (Utility/Shore)",
		160: "Power Converter - Battery Charger",
		161: "Power Converter - Battery Charger+Inverter",
		162: "Power Converter - Inverter",
		163: "Power Converter - DC",
		170: "Battery",
		180: "Engine Gateway",
	},
	40: PgnLookup{
		130: "Follow-up Controller",
		140: "Mode Controller",
		150: "Autopilot",
		155: "Rudder",
		160: "Heading Sensors (Deprecated)",
		170: "Trim (Tabs)/Interceptors",
		180: "Attitude (Pitch, Roll, Yaw) Control",
	},
	50: PgnLookup{
		130: "Engineroom Monitoring (Deprecated)",
		140: "Engine",
		141: "DC Generator/Alternator",
		150: "Engine Controller (Deprecated)",
		151: "AC Generator",
		155: "Motor",
		160: "Engine Gateway",
		165: "Transmission",
		170: "Throttle/Shift Control",
		180: "Actuator (Deprecated)",
		190: "Gauge Interface (Deprecated)",
		200: "Gauge Large (Deprecated)",
		210: "Gauge Small (Deprecated)",
	},
	60: PgnLookup{
		130: "Bottom Depth",
		135: "Bottom Depth/Speed",
		140: "Ownship Attitude",
		145: "Ownship Positon (GNSS)",
		150: "Ownship Position (Loran C)",
		155: "Speed",
		160: "Turn Rate Indicator (Deprecated)",
		170: "Integrated Navigation (Deprecated)",
		175: "Integrated Navigation System",
		190: "Navigation Management",
		195: "Automatic Identification System (AIS)",
		200: "Radar",
		201: "Infrared Imaging",
		205: "ECDIS (Deprecated)",
		210: "ECS (Deprecated)",
		220: "Direction Finder (Deprecated)",
		230: "Voyage Status",
	},
	70: PgnLookup{
		130: "EPIRB (Deprecated)",
		140: "AIS (Deprecated)",
		150: "DSC (Deprecated)",
		160: "Data Receiver/Transmitter",
		170: "Satellite",
		180: "Radio-telephone (MF/HF) (Deprecated)",
		190: "Radiotelephone",
	},
	75: PgnLookup{
		130: "Temperature",
		140: "Pressure",
		150: "Fluid Level",
		160: "Flow",
		170: "Humidity",
	},
	80: PgnLookup{
		130: "Time/Date Systems (Deprecated)",
		140: "VDR (Deprecated)",
		150: "Integrated Instrumentation (Deprecated)",
		160: "General Purpose Displays (Deprecated)",
		170: "General Sensor Box (Deprecated)",
		180: "Weather Instruments (Deprecated)",
		190: "Transducer/General (Deprecated)",
		200: "NMEA 0183 Converter (Deprecated)",
	},
	85: PgnLookup{
		130: "Atmospheric",
		160: "Aquatic",
	},
	90: PgnLookup{
		130: "HVAC",
	},
	100: PgnLookup{
		130: "Scale (Catch)",
	},
	120: PgnLookup{
		130: "Display",
		140: "Alarm Enunciator",
	},
	125: PgnLookup{
		130: "Multimedia Player",
		140: "Multimedia Controller",
	},
}

var lookupIndustryCode = PgnLookup{
	4: "Marine",
}

var lookupRepeatIndicator = PgnLookup{
	0: "Initial",
	1: "First retransmission",
	2: "Second retransmission",
	3: "Final retransmission",
}

var lookupSteeringMode = PgnLookup{
	0: "Main Steering",
	1: "Non-Followup Device",
	2: "Followup Device",
	3: "Heading Control Standaline",
	4: "Heading Control",
	5: "Track Control",
}

var lookupTurnMode = PgnLookup{
	0: "Rudder Limit Controlled",
	1: "Turn Rate Controlled",
	2: "Radius Controlled",
}

var lookupCommandedRudderDirection = PgnLookup{
	0: "No Order",
	1: "Move to starboard",
	2: "Move to port",
}

var lookupEngineInstance = PgnLookup{
	0: "Single Engine or Dual Engine Port",
	1: "Dual Engine Starboard",
}

var lookupGearStatus = PgnLookup{
	0: "Forward",
	1: "Neutral",
	2: "Reverse",
	3: "Unknown",
}

var lookupPositionAccuracy = PgnLookup{
	0: "Low",
	1: "High",
}

// http://www.navcen.uscg.gov/?pageName=AISMessagesAStatic
var lookupShipType = PgnLookup{
	0:  "Unavailable",
	20: "Wing in ground",
	21: "Wing in ground hazard cat A",
	22: "Wing in ground hazard cat B",
	23: "Wing in ground hazard cat C",
	24: "Wing in ground hazard cat D",
	29: "Wing in ground (no other information)",
	30: "Fishing",
	31: "Towing",
	32: "Towing and length exceeds 200m or wider than 25m",
	33: "Engaged in dredging or underwater operations",
	34: "Engaged in diving operations",
	35: "Engaged in military operations",
	36: "Sailing",
	37: "Pleasure",
	40: "High speed craft",
	41: "High speed craft hazard cat A",
	42: "High speed craft hazard cat B",
	43: "High speed craft hazard cat C",
	44: "High speed craft hazard cat D",
	49: "High speed craft (no additional information)",
	50: "Pilot vessel",
	51: "SAR",
	52: "Tug",
	53: "Port tender",
	54: "Anti-pollution",
	55: "Law enforcement",
	56: "Spare",
	57: "Spare #2",
	58: "Medical",
	59: "RR Resolution No.18",
	60: "Passenger ship",
	61: "Passenger ship hazard cat A",
	62: "Passenger ship hazard cat B",
	63: "Passenger ship hazard cat C",
	64: "Passenger ship hazard cat D",
	69: "Passenger ship (no additional information)",
	70: "Cargo ship",
	71: "Cargo ship hazard cat A",
	72: "Cargo ship hazard cat B",
	73: "Cargo ship hazard cat C",
	74: "Cargo ship hazard cat D",
	79: "Cargo ship (no additional information)",
	80: "Tanker",
	81: "Tanker hazard cat A",
	82: "Tanker hazard cat B",
	83: "Tanker hazard cat C",
	84: "Tanker hazard cat D",
	89: "Tanker (no additional information)",
	90: "Other",
	91: "Other hazard cat A",
	92: "Other hazard cat B",
	93: "Other hazard cat C",
	94: "Other hazard cat D",
	99: "Other (no additional information)",
}

var lookupTimeStamp = PgnLookup{
	60: "Not available",
	61: "Manual input mode",
	62: "Dead reckoning mode",
	63: "Positioning system is inoperative",
}

var lookupGns = PgnLookup{
	0: "GPS",
	1: "GLONASS",
	2: "GPS+GLONASS",
	3: "GPS+SBAS/WAAS",
	4: "GPS+SBAS/WAAS+GLONASS",
	5: "Chayka",
	6: "Integrated",
	7: "Surveyed",
	8: "Galileo",
}

var lookupGnsAis = PgnLookup{
	0: "Undefined",
	1: "GPS",
	2: "GLONASS",
	3: "GPS+GLONASS",
	4: "Loran-C",
	5: "Chayka",
	6: "Integrated",
	7: "Surveyed",
	8: "Galileo",
}

var lookupGnsIntegrity = PgnLookup{
	0: "No integrity checking",
	1: "Safe",
	2: "Caution",
}

var lookupGnsMethod = PgnLookup{
	0: "No GNSS",
	1: "GNSS Fix",
	2: "DGNSS Fix",
	3: "Precise GNSS",
	4: "RTK Fixed Integer",
	5: "RTK Float",
	6: "Estimated (DR) Mode",
	7: "Manual Input",
	8: "Simulate Mode",
}

var lookupSystemTime = PgnLookup{
	0: "GPS",
	1: "GLONASS",
	2: "Radio Station",
	3: "Local Cesium clock",
	4: "Local Rubidium clock",
	5: "Local Crystal clock",
}

var lookupMagneticVariation = PgnLookup{
	0: "Manual",
	1: "Automatic Chart",
	2: "Automatic Table",
	3: "Automatic Calculation",
	4: "WMM 2000",
	5: "WMM 2005",
	6: "WMM 2010",
	7: "WMM 2015",
	8: "WMM 2020",
}

var lookupNavCalculation = PgnLookup{
	0: "Great Circle",
	1: "Rhumb Line",
}

var lookupNavMarkType = PgnLookup{
	0: "Collision",
	1: "Turning point",
	2: "Reference",
	3: "Wheelover",
	4: "Waypoint",
}

var lookupResidualMode = PgnLookup{
	0: "Autonomous",
	1: "Differential Enhanced",
	2: "Estimated",
	3: "Simulator",
	4: "Manual",
}

var lookupWindReference = PgnLookup{
	0: "True (ground referenced to North)",
	1: "Magnetic (ground referenced to Magnetic North)",
	2: "Apparent",
	3: "True (boat referenced)",
	4: "True (water referenced)",
}

var lookupYesNo = PgnLookup{
	0: "No",
	1: "Yes",
}

var lookupDirectionReference = PgnLookup{
	0: "True",
	1: "Magnetic",
	2: "Error",
	3: "Null",
}

var lookupNavStatus = PgnLookup{
	0: "Under way using engine",
	1: "At anchor",
	2: "Not under command",
	3: "Restricted manoeuverability",
	4: "Constrained by her draught",
	5: "Moored",
	6: "Aground",
	7: "Engaged in Fishing",
	8: "Under way sailing",
}

var lookupPowerFactor = PgnLookup{
	0: "Leading",
	1: "Lagging",
	2: "Error",
}

var lookupTemperatureSource = PgnLookup{
	0: "Sea Temperature",
	1: "Outside Temperature",
	2: "Inside Temperature",
	3: "Engine Room Temperature",
	4: "Main Cabin Temperature",
	5: "Live Well Temperature",
	6: "Bait Well Temperature",
	7: "Refridgeration Temperature",
	8: "Heating System Temperature",
	9: "Freezer Temperature",
}

var lookupHumidityInstance = PgnLookup{
	0: "Inside",
	1: "Outside",
}

var lookupPressureSource = PgnLookup{
	0: "Atmospheric",
	1: "Water",
	2: "Steam",
	3: "Compressed Air",
	4: "Hydraulic",
}

var lookupTankType = PgnLookup{
	0: "Fuel",
	1: "Water",
	2: "Gray water",
	3: "Live well",
	4: "Oil",
	5: "Black water",
}

var lookupTideTendency = PgnLookup{
	0: "Falling",
	1: "Rising",
}

var lookupIsoAckResults = PgnLookup{
	0: "ACK",
	1: "NAK",
	2: "Access Denied",
	3: "Address Busy",
}

var lookupWaveform = PgnLookup{
	0: "Sine Wave",
	1: "Modified Sine Wave",
	6: "Error",
	7: "Data Not Available",
}

var lookupOffOn = PgnLookup{
	0: "Off",
	1: "On",
}

var lookupStandbyOn = PgnLookup{
	0: "Standby",
	1: "On",
}

var lookupAcceptability = PgnLookup{
	0: "Bad Level",
	1: "Bad Frequency",
	2: "Being Qualified",
	3: "Good",
}

var lookupTrackStatus = PgnLookup{
	0: "Cancelled",
	1: "Acquiring",
	2: "Tracking",
	3: "Lost",
}

var lookupTargetAcquisition = PgnLookup{
	0: "Manual",
	1: "Automatic",
}

var lookupLine = PgnLookup{
	0: "Line 1",
	1: "Line 2",
	2: "Line 3",
	3: "Reserved",
}

var lookupFunctionCode = PgnLookup{
	0: "Transmit PGN List",
	1: "Receive PGN List",
}

var lookupGnssMode = PgnLookup{
	0: "1D",
	1: "2D",
	2: "3D",
	3: "Auto",
	4: "Reserved",
	5: "Reserved",
	6: "Error",
}

var lookupDGnssMode = PgnLookup{
	0: "no SBAS",
	1: "SBAS",
	3: "SBAS",
}

var lookupGnssAntenna = PgnLookup{
	0: "use last 3D height",
	1: "Use antenna altitude",
}

var lookupGnssSatMode = PgnLookup{
	3: "Range residuals used to calculate position",
}

var lookupGnssSatStatus = PgnLookup{
	0: "Not tracked",
	1: "Tracked",
	2: "Used",
	3: "Not tracked+Diff",
	4: "Tracked+Diff",
	5: "Used+Diff",
}

var lookupAirmarBootState = PgnLookup{
	0: "In Startup Monitor",
	1: "Running Bootloader",
	2: "Running Application",
}

var lookupAirmarCalibrateFunction = PgnLookup{
	0: "Normal/cancel calibration",
	1: "Enter calibration mode",
	2: "Reset calibration to 0",
	3: "Verify",
	4: "Reset compass to defaults",
	5: "Reset damping to defaults",
}

var lookupAirmarCalibrationStatus = PgnLookup{
	0: "Queried",
	1: "Passed",
	2: "Failed - timeout",
	3: "Failed - tilt error",
	4: "Failed - other",
	5: "In progress",
}

var lookupAirmarFormatCode = PgnLookup{
	1: "Format Code 1",
}

var lookupAirmarAccessLevel = PgnLookup{
	0: "Locked",
	1: "Unlocked Level 1",
	2: "Unlocked Level 2",
}

var lookupAirmarCogSubstitute = PgnLookup{
	0: "Use HDG only",
	1: "Allow COG to replace HDG",
}

var lookupAirmarControl = PgnLookup{
	0: "Report previous values",
	1: "Generate new values",
}

var lookupAirmarDepthQuality = PgnLookup{
	0: "No depth lock",
}

var lookupAirmarTempInstance = PgnLookup{
	0: "Device Sensor",
	1: "Onboard Water Sensor",
	2: "Optional Water Sensor",
}

var lookupAirmarSpeedFilter = PgnLookup{
	0: "No filter",
	1: "Basic IIR filter",
}

var lookupAirmarTempFilter = PgnLookup{
	0: "No filter",
	1: "Basic IIR filter",
}

var lookupAirmarTransmissionInterval = PgnLookup{
	0: "Measure Interval",
	1: "Requested by user",
}

var lookupAirmarTestId = PgnLookup{
	1: "Format Code",
	2: "Factory EEPROM",
	3: "User EEPROM",
	4: "Water Temp Sensor",
	5: "Sonar Transceiver",
	6: "Speed sensor",
	7: "Internal temperature sensor",
	8: "Battery voltage sensor",
}

var lookupSonicHubControl = PgnLookup{
	0:   "Set",
	128: "Ack",
}

var lookupSonicHubSource = PgnLookup{
	0: "AM",
	1: "FM",
	2: "iPod",
	3: "USB",
	4: "AUX",
	5: "AUX 2",
	6: "Mic",
}

var lookupSonicHubTuning = PgnLookup{
	1: "Seeking Up",
	2: "Tuned",
	3: "Seeking Down",
}

var lookupSonicHubMute = PgnLookup{
	1: "Mute On",
	2: "Mute Off",
}

var lookupSonicHubPlaylist = PgnLookup{
	1: "Report",
	4: "Next Song",
	6: "Previous Song",
}

var lookupSonicHubZone = PgnLookup{
	0: "Zone 1",
	1: "Zone 2",
	2: "Zone 3",
}

// https://www.nmea.org/Assets/20140409%20nmea%202000%20registration%20list.pdf
var lookupCompanyCode = PgnLookup{
	174:  "Volvo Penta",
	199:  "Actia Corporation",
	273:  "Actisense",
	578:  "Advansea",
	215:  "Aetna Engineering/Fireboy-Xintex",
	135:  "Airmar",
	459:  "Alltek Marine Electronics Group",
	274:  "Amphenol LTW Technology",
	502:  "Attwood Marine",
	381:  "B&G",
	185:  "Beede Electrical",
	295:  "BEP",
	396:  "Beyond Measure",
	148:  "Blue Water Data",
	163:  "Evinrude/BRP Bombardier",
	394:  "Capi 2",
	176:  "Carling Technologies",
	165:  "CPac Systems AB",
	286:  "Coelmo Srl Italy",
	404:  "Com Nav",
	440:  "Cummins",
	329:  "Dief",
	437:  "Digital Yacht Ltd",
	201:  "Disenos Y Technologia",
	211:  "DNA Group, Inc",
	426:  "Egersund Marine Electronics AS",
	373:  "Electronic Design",
	427:  "Em-Trak Marine Electronics Ltd",
	224:  "EMMI Network",
	304:  "Empirbus",
	243:  "eRide",
	1863: "Faria Instruments",
	356:  "Fischer Panda",
	192:  "Floscan Instrument Co Inc",
	1855: "Furuno USA",
	419:  "Fusion",
	78:   "FW Murphy",
	229:  "Garmin",
	385:  "Geonav",
	378:  "Glendinning",
	475:  "GME / Standard Communications Pty",
	272:  "Groco",
	283:  "Hamilton Jet",
	88:   "Hemisphere GPS",
	257:  "Honda",
	467:  "Hummingbird Marine Electronics",
	315:  "ICOM",
	1853: "Japan Radio Co",
	1859: "Kvasar AB",
	579:  "KVH",
	85:   "Kohler",
	345:  "Korea Maritime University",
	499:  "LCJ Capteurs",
	1858: "Litton",
	400:  "Livorsi Marine",
	140:  "Lowrance Electronics",
	137:  "Maretron",
	571:  "Marinecraft (South Korea)",
	307:  "MBW Technologies",
	355:  "Mastervolt",
	144:  "Mercury Marine",
	1860: "MMP",
	198:  "Mystic Valley Communications",
	147:  "Nautibus Electronic GmbH",
	275:  "Navico",
	1852: "Navionics",
	503:  "Naviop",
	193:  "Nobeltec",
	517:  "Noland",
	374:  "Northern Lights",
	1854: "Northstar",
	305:  "Novatel",
	478:  "Ocean Sat BV",
	161:  "Offshore Systems UK",
	573:  "Orolia Ltd",
	328:  "Qwerty",
	451:  "Parker Hannifin",
	1851: "Raymarine, Inc",
	370:  "Rolls Royce Marine",
	384:  "Rose Point Navigation Systems", // was Capano Light
	235:  "SailorMade/Tetra",
	580:  "San Jose Technologies",
	460:  "San Giorgio S.E.I.N. srl",
	1862: "Sanshin Industries / Yamaha",
	471:  "Sea Cross Marine AB",
	285:  "Sea Recovery",
	1857: "Simrad",
	470:  "Sitex",
	306:  "Sleipner Motor AS",
	1850: "Teleflex",
	351:  "Thrane and Thrane",
	431:  "Tohatsu Co JP",
	518:  "Transas USA",
	1856: "Trimble",
	422:  "True Heading",
	80:   "Twin Disc",
	1861: "Vector Cantech",
	466:  "Veethree",
	421:  "Vertex Standard Co Ltd",
	504:  "Vesper Marine",
	358:  "Victron",
	493:  "Watcheye",
	154:  "Westerbeke Corp",
	168:  "Xantrex Technology",
	233:  "Yacht Monitoring Solutions",
	172:  "Yanmar",
	228:  "ZF Marine Electronics",
}

var lookupPriorityLevel = PgnLookup{
	8: "Leave priority unchanged",
}

var lookupSimnetBacklightLevel = PgnLookup{
	1:  "Day Mode",
	4:  "Night Mode",
	11: "Level 1",
	22: "Level 2",
	33: "Level 3",
	44: "Level 4",
	55: "Level 5",
	66: "Level 6",
	77: "Level 7",
	88: "Level 8",
	99: "Level 9",
}

var lookupSimnetApEvents = PgnLookup{
	6:  "Standby",
	9:  "Auto mode",
	10: "Nav mode",
	13: "Non Follow Up mode",
	15: "Wind mode",
	26: "Change Course",
}

var lookupSimnetDirection = PgnLookup{
	2: "Port",
	3: "Starboard",
	4: "Left rudder (port)",
	5: "Right rudder (starboard)",
}

var lookupSimnetAlarm = PgnLookup{
	56: "Clear",
	57: "Raise",
}

var lookupFusionMute = PgnLookup{
	1: "Muted",
	2: "Not muted",
}

var lookupFusionTransport = PgnLookup{
	1: "Paused",
}

var lookupFusionTimeFormat = PgnLookup{
	0: "12h",
	1: "24h",
}

var lookupFusionReplayMode = PgnLookup{
	9:  "USB repeat",
	10: "USB shuffle",
	12: "iPod repeat",
	13: "iPod shuffle",
}

var lookupFusionReplayStatus = PgnLookup{
	0: "Off",
	1: "One/Track",
	2: "All/Album",
}
