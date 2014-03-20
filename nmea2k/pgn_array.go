package nmea2k

import (
  "encoding/json"
  "encoding/xml"
  "math"
)

const ACTISENSE_BEM = 0x40000 // Actisense specific fake PGNs

const RadianToDegree = 180.0 / math.Pi

const RES_LAT_LONG_PRECISION = 1e7
const RES_LAT_LONG = 1.0e-7
const RES_LAT_LONG_64 = 1.0e-16

const LEN_VARIABLE = 0
const RES_NOTUSED = 0
const RES_DEGREES = 1e-4 * RadianToDegree
const RES_ROTATION = 1e-3 / 32.0 * RadianToDegree
const RES_ASCII = -1
const RES_LATITUDE = -2
const RES_LONGITUDE = -3
const RES_DATE = -4
const RES_TIME = -5
const RES_TEMPERATURE = -6
const RES_6BITASCII = -7
const RES_INTEGER = -8
const RES_LOOKUP = -9
const RES_BINARY = -10
const RES_MANUFACTURER = -11
const RES_STRING = -12
const RES_FLOAT = -13
const RES_PRESSURE = -14
const RES_STRINGLZ = -15
const MAX_RES_LOOKUP = 15

type Field struct {
	Name        string
	Size        uint32
	Resolution  float64
	Signed      bool
	Units       interface{}
	Description string
	Offset      int32
}

type Pgn struct {
	Description     string
	Pgn             uint32
	IsKnown         bool    // Are we pretty sure we've got all fields with correct definitions?
	Size            uint32  // (Minimal) size of this PGN. Helps to determine fast/single frame and initial malloc
	RepeatingFields uint32  // How many fields at the end repeat until the PGN is exhausted?
	FieldList       []Field // Fields in the PGN
}

type PgnArray []Pgn

type PgnLookup map[int]string

func (inVal PgnLookup) MarshalJSON() ([]byte, error) {
  outVal := make(map[string]string)

  for k, v := range inVal {
    outVal[strconv.Itoa(k)] = v
  }

  return json.Marshal(outVal)
}

var lookupAisAccuracy = PgnLookup{
	0: "Low",
	1: "High",
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
	0:   "Reserved for 2000 Use",
	10:  "System tools",
	20:  "Safety systems",
	25:  "Internetwork device",
	30:  "Electrical Distribution",
	35:  "Electrical Generation",
	40:  "Steering and Control surfaces",
	50:  "Propulsion",
	60:  "Navigation",
	70:  "Communication",
	75:  "Sensor Communication Interface",
	80:  "Instrumentation/general systems",
	85:  "External Environment",
	90:  "Internal Environment",
	100: "Deck, cargo and fishing equipment systems",
	120: "Display",
	125: "Entertainment",
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

var lookupAirmarFormatCode = map[uint32]string{
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

var lookupAirmarTempInstance = PgnLookup{
	0: "Device Sensor",
	1: "Onboard Water Sensor",
	2: "Optional Water Sensor",
}

var lookupAirmarSpeedFilter = PgnLookup{
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

// http://www.nmea.org/Assets/20121212%20nmea%202000%20registration%20list.pdf
var lookupCompanyCode = PgnLookup{
	174:  "Volvo Penta",
	199:  "Actia Corporation",
	273:  "Actisense",
	215:  "Aetna Engineering/Fireboy-Xintex",
	135:  "Airmar",
	459:  "Alltek",
	381:  "B&G",
	185:  "Beede Electrical",
	295:  "BEP",
	396:  "Beyond Measure",
	148:  "Blue Water Data",
	163:  "Evinrude/Bombardier",
	384:  "Camano Light",
	394:  "Capi 2",
	176:  "Carling",
	165:  "CPac",
	286:  "Coelmo",
	404:  "ComNav",
	440:  "Cummins",
	329:  "Dief",
	437:  "Digital Yacht",
	201:  "Disenos Y Technologia",
	211:  "DNA Group",
	426:  "Egersund Marine",
	373:  "Electronic Design",
	427:  "Em-Trak",
	224:  "EMMI Network",
	304:  "Empirbus",
	243:  "eRide",
	1863: "Faria Instruments",
	356:  "Fischer Panda",
	192:  "Floscan",
	1855: "Furuno",
	419:  "Fusion",
	78:   "FW Murphy",
	229:  "Garmin",
	385:  "Geonav",
	378:  "Glendinning",
	475:  "GME / Standard",
	272:  "Groco",
	283:  "Hamilton Jet",
	88:   "Hemisphere GPS",
	257:  "Honda",
	467:  "Hummingbird",
	315:  "ICOM",
	1853: "JRC",
	1859: "Kvasar",
	85:   "Kohler",
	345:  "Korea Maritime University",
	499:  "LCJ Capteurs",
	1858: "Litton",
	400:  "Livorsi",
	140:  "Lowrance",
	137:  "Maretron",
	307:  "MBW",
	355:  "Mastervolt",
	144:  "Mercury",
	1860: "MMP",
	198:  "Mystic Valley Comms",
	147:  "Nautibus",
	275:  "Navico",
	1852: "Navionics",
	503:  "Naviop",
	193:  "Nobeltec",
	517:  "Noland",
	374:  "Northern Lights",
	1854: "Northstar",
	305:  "Novatel",
	161:  "Offshore Systems",
	328:  "Qwerty",
	451:  "Parker Hannifin",
	1851: "Raymarine",
	370:  "Rolls Royce",
	235:  "SailorMade/Tetra",
	460:  "San Giorgio",
	471:  "Sea Cross",
	285:  "Sea Recovery",
	1862: "Yamaha",
	1857: "Simrad",
	470:  "Sitex",
	306:  "Sleipner",
	1850: "Teleflex",
	351:  "Thrane and Thrane",
	431:  "Tohatsu",
	518:  "Transas",
	1856: "Trimble",
	422:  "True Heading",
	80:   "Twin Disc",
	1861: "Vector Cantech",
	466:  "Veethree",
	421:  "Vertex",
	504:  "Vesper",
	358:  "Victron",
	493:  "Watcheye",
	154:  "Westerbeke",
	168:  "Xantrex",
	233:  "Yacht Monitoring Solutions",
	172:  "Yanmar",
	228:  "ZF",
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

var PgnList = PgnArray{
	{"Unknown PGN", 0, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, nil, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"ISO Acknowledgement", 59392, true, 8, 0, []Field{
		{"Control", 8, RES_LOOKUP, false, lookupIsoAckResults, "", 0},
		{"Group Function", 8, 1, false, nil, "", 0},
		{"Reserved", 24, RES_BINARY, false, nil, "Reserved", 0},
		{"PGN", 24, RES_INTEGER, false, nil, "Parameter Group Number of requested information", 0}},
	},

	{"ISO Request", 59904, true, 3, 0, []Field{
		{"PGN", 24, RES_INTEGER, false, nil, "", 0}},
	},

	{"ISO Address Claim", 60928, true, 8, 0, []Field{
		{"Unique Number", 21, RES_BINARY, false, nil, "ISO Identity Number", 0},
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, nil, "", 0},
		{"Device Instance Lower", 4, 1, false, nil, "ISO ECU Instance", 0},
		{"Device Instance Upper", 4, 1, false, nil, "ISO Function Instance", 0},
		{"Device Function", 8, 1, false, nil, "ISO Function", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "", 0},
		{"Device Class", 7, RES_LOOKUP, false, lookupDeviceClass, "", 0},
		{"System Instance", 4, 1, false, nil, "ISO Device Class Instance", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 3, 1, false, nil, "ISO Self Configurable", 0}},
	},

	{"ISO: Manu. Proprietary single-frame addressed", 61184, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, nil, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, 1, false, nil, "", 0},
		{"Payload", 48, RES_BINARY, false, nil, "", 0}},
	},

	// Maretron ACM 100 manual documents PGN 65001-65030

	{"Bus #1 Phase C Basic AC Quantities", 65001, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0}},
	},

	{"Bus #1 Phase B Basic AC Quantities", 65002, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0}},
	},

	{"Bus #1 Phase A Basic AC Quantities", 65003, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0}},
	},

	{"Bus #1 Average Basic AC Quantities", 65004, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0}},
	},

	{"Utility Total AC Energy", 65005, false, 8, 0, []Field{
		{"Total Energy Export", 32, 1, false, "kWh", "", 0},
		{"Total Energy Import", 32, 1, false, "kWh", "", 0}},
	},

	{"Utility Phase C AC Reactive Power", 65006, false, 8, 0, []Field{
		{"Reactive Power", 16, 1, false, "VAr", "", 0},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Utility Phase C AC Power", 65007, false, 8, 0, []Field{
		{"Real Power", 32, 1, true, "W", "", -2000000000},
		{"Apparent Power", 32, 1, true, "VA", "", -2000000000}},
	},

	{"Utility Phase C Basic AC Quantities", 65008, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Utility Phase B AC Reactive Power", 65009, false, 8, 0, []Field{
		{"Reactive Power", 16, 1, false, "VAr", "", 0},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Utility Phase B AC Power", 65010, false, 8, 0, []Field{
		{"Real Power", 32, 1, true, "W", "", -2000000000},
		{"Apparent Power", 32, 1, true, "VA", "", -2000000000}},
	},

	{"Utility Phase B Basic AC Quantities", 65011, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Utility Phase A AC Reactive Power", 65012, false, 8, 0, []Field{
		{"Reactive Power", 32, 1, true, "VAr", "", -2000000000},
		{"Power Factor", 16, 1.0 / 16384, true, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Utility Phase A AC Power", 65013, false, 8, 0, []Field{
		{"Real Power", 32, 1, true, "W", "", -2000000000},
		{"Apparent Power", 32, 1, true, "VA", "", -2000000000}},
	},

	{"Utility Phase A Basic AC Quantities", 65014, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Utility Total AC Reactive Power", 65015, false, 8, 0, []Field{
		{"Reactive Power", 32, 1, true, "VAr", "", -2000000000},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Utility Total AC Power", 65016, false, 8, 0, []Field{
		{"Real Power", 32, 1, true, "W", "", -2000000000},
		{"Apparent Power", 32, 1, true, "VA", "", -2000000000}},
	},

	{"Utility Average Basic AC Quantities", 65017, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Generator Total AC Energy", 65018, false, 8, 0, []Field{
		{"Total Energy Export", 32, 1, false, "kWh", "", 0},
		{"Total Energy Import", 32, 1, false, "kWh", "", 0}},
	},

	{"Generator Phase C AC Reactive Power", 65019, false, 8, 0, []Field{
		{"Reactive Power", 16, 1, false, "VAr", "", -2000000000},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Generator Phase C AC Power", 65020, false, 8, 0, []Field{
		{"Real Power", 16, 1, false, "W", "", -2000000000},
		{"Apparent Power", 16, 1, false, "VA", "", -2000000000}},
	},

	{"Generator Phase C Basic AC Quantities", 65021, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Generator Phase B AC Reactive Power", 65022, false, 8, 0, []Field{
		{"Reactive Power", 16, 1, false, "VAr", "", -2000000000},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Generator Phase B AC Power", 65023, false, 8, 0, []Field{
		{"Real Power", 16, 1, false, "W", "", -2000000000},
		{"Apparent Power", 16, 1, false, "VA", "", -2000000000}},
	},

	{"Generator Phase B Basic AC Quantities", 65024, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Generator Phase A AC Reactive Power", 65025, false, 8, 0, []Field{
		{"Reactive Power", 16, 1, false, "VAr", "", -2000000000},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Generator Phase A AC Power", 65026, false, 8, 0, []Field{
		{"Real Power", 16, 1, false, "W", "", -2000000000},
		{"Apparent Power", 16, 1, false, "VA", "", -2000000000}},
	},

	{"Generator Phase A Basic AC Quantities", 65027, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	{"Generator Total AC Reactive Power", 65028, false, 8, 0, []Field{
		{"Reactive Power", 16, 1, false, "VAr", "", -2000000000},
		{"Power Factor", 16, 1.0 / 16384, false, nil, "", 0},
		{"Power Factor Lagging", 2, RES_LOOKUP, false, lookupPowerFactor, "", 0}},
	},

	{"Generator Total AC Power", 65029, false, 8, 0, []Field{
		{"Real Power", 16, 1, false, "W", "", -2000000000},
		{"Apparent Power", 16, 1, false, "VA", "", -2000000000}},
	},

	{"Generator Average Basic AC Quantities", 65030, false, 8, 0, []Field{
		{"Line-Line AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"Line-Neutral AC RMS Voltage", 16, 1, false, "V", "", 0},
		{"AC Frequency", 16, 1.0 / 128, false, "Hz", "", 0},
		{"AC RMS Current", 16, 1, false, "A", "", 0}},
	},

	// ISO 11783 defined this message to provide a mechanism for assigning a
	// network address to a node. The NAME information in the data portion of the
	// message must match the name information of the node whose network address
	// is to be set.
	{"ISO Commanded Address", 65240, false, 8, 0, []Field{
		{"Unique Number", 21, RES_BINARY, false, nil, "ISO Identity Number", 0},
		{"Manufacturer Code", 11, 1, false, nil, "", 0},
		{"Device Instance Lower", 4, 1, false, nil, "ISO ECU Instance", 0},
		{"Device Instance Upper", 4, 1, false, nil, "ISO Function Instance", 0},
		{"Device Function", 8, 1, false, nil, "ISO Function", 0},
		{"Reserved", 1, 1, false, nil, "", 0},
		{"Device Class", 7, RES_LOOKUP, false, lookupDeviceClass, "", 0},
		{"System Instance", 4, 1, false, nil, "ISO Device Class Instance", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 3, 1, false, nil, "ISO Self Configurable", 0},
		{"New Source Address", 8, 1, false, nil, "", 0}},
	},

	// ISO 11783: 65,280 to 65,535 (0xFF00 to 0xFFFF): Propietary PDU-2 messages
	{"ISO: Manu. Proprietary single-frame non-addressed", 65280, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, nil, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, 1, false, nil, "", 0},
		{"Payload", 48, RES_BINARY, false, nil, "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Boot State Acknowledgment", 65285, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Boot State", 4, RES_LOOKUP, false, lookupAirmarBootState, "", 0}},
	},

	{"Lowrance: Temperature", 65285, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=140", "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Temperature Instance", 4, 1, false, nil, "", 0},
		{"Temperature Source", 4, 1, false, nil, "", 0},
		{"Actual Temperature", 16, RES_TEMPERATURE, false, "K", "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Boot State Request", 65286, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Access Level", 65287, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Format Code", 3, RES_LOOKUP, false, lookupAirmarFormatCode, "", 0},
		{"Access Level", 3, RES_LOOKUP, false, lookupAirmarAccessLevel, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Access Seed/Key", 32, RES_INTEGER, false, nil, "When transmitted, it provides a seed for an unlock operation. It is used to provide the key during PGN 126208.", 0}},
	},

	{"Simnet: Configure Temperature Sensor", 65287, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Trim Tab Sensor Calibration", 65289, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Paddle Wheel Speed Configuration", 65290, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Clear Fluid Level Warnings", 65292, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: LGC-2000 Configuration", 65293, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Reprogram Status", 65325, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Autopilot Mode", 65341, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Depth Quality Factor", 65408, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"SID", 8, 1, false, nil, "", 0},
		{"Depth Quality Factor", 4, RES_LOOKUP, false, ",0=No Depth Lock", "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Device Information", 65410, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"SID", 8, 1, false, nil, "", 0},
		{"Internal Device Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Supply Voltage", 16, 0.01, false, "V", "", 0},
		{"Reserved", 8, 1, false, nil, "", 0}},
	},

	{"Simnet: Autopilot Mode", 65480, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	// http://www.maretron.com/support/manuals/DST100UM_1.2.pdf
	{"NMEA - Request group function", 126208, true, 8, 2, []Field{
		{"Function Code", 8, RES_INTEGER, false, "=0", "Request", 0},
		{"PGN", 24, RES_INTEGER, false, nil, "Requested PGN", 0},
		{"Transmission interval", 32, 1, false, nil, "", 0},
		{"Transmission interval offset", 16, 1, false, nil, "", 0},
		{"# of Requested Parameters", 8, 1, false, nil, "How many parameter pairs will follow", 0},
		{"Parameter Index", 8, RES_INTEGER, false, nil, "Parameter index", 0},
		{"Parameter Value", LEN_VARIABLE, RES_INTEGER, false, nil, "Parameter value, variable length", 0}},
	},

	{"NMEA - Command group function", 126208, true, 8, 2, []Field{
		{"Function Code", 8, RES_INTEGER, false, "=1", "Command", 0},
		{"PGN", 24, RES_INTEGER, false, nil, "Commanded or requested PGN", 0},
		{"Priority", 4, 1, false, lookupPriorityLevel, "", 0},
		{"Reserved", 4, 1, false, nil, "", 0},
		{"# of Commanded Parameters", 8, 1, false, nil, "How many parameter pairs will follow", 0},
		{"Parameter Index", 8, RES_INTEGER, false, nil, "Parameter index", 0},
		{"Parameter Value", LEN_VARIABLE, RES_INTEGER, false, nil, "Parameter value, variable length", 0}},
	},

	{"NMEA - Acknowledge group function", 126208, true, 8, 1, []Field{
		{"Function Code", 8, RES_INTEGER, false, "=2", "Acknowledge", 0},
		{"PGN", 24, RES_INTEGER, false, nil, "Commanded or requested PGN", 0},
		{"PGN error code", 4, 1, false, nil, "", 0},
		{"Transmission interval/Priority error code", 4, 1, false, nil, "", 0},
		{"# of Commanded Parameters", 8, 1, false, nil, "", 0},
		{"Parameter Error", 8, RES_INTEGER, false, nil, "", 0}},
	},

	/////////////////////////// RESPONSE TO REQUEST PGNS ////////////////////////
	{"Maretron: Slave Response", 126270, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, nil, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Product Code", 16, 1, false, nil, "0x1b2: SSC200", 0},
		{"Software Code", 16, 1, false, nil, "", 0},
		{"Command", 8, 1, false, nil, "0x50=Deviation calibration result", 0},
		{"Status", 8, 1, false, nil, "", 0}},
	},

	{"PGN List (Transmit and Receive)", 126464, false, 8, 1, []Field{
		{"Function Code", 8, RES_LOOKUP, false, lookupFunctionCode, "Transmit or receive PGN Group Function Code", 0},
		{"PGN", 24, RES_INTEGER, false, nil, "", 0}},
	},

	{"Manufacturer Propietary: Addressable Multi-Frame", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, nil, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Payload", LEN_VARIABLE, RES_BINARY, false, nil, "", 0}},
	},

	{"Airmar: Addressable Multi-Frame", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=0", "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/PB200UserManual.pdf
	{"Airmar: Attitude Offset", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=32", "Attitude Offsets", 0},
		{"Azimuth offset", 16, RES_DEGREES, true, "deg", "Positive: sensor rotated to port, negative: sensor rotated to starboard", 0},
		{"Pitch offset", 16, RES_DEGREES, true, "deg", "Positive: sensor tilted to bow, negative: sensor tilted to stern", 0},
		{"Roll offset", 16, RES_DEGREES, true, "deg", "Positive: sensor tilted to port, negative: sensor tilted to starboard", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/PB200UserManual.pdf
	{"Airmar: Calibrate Compass", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=33", "Calibrate Compass", 0},
		{"Calibrate Function", 8, RES_LOOKUP, false, lookupAirmarCalibrateFunction, "", 0},
		{"Calibration Status", 8, RES_LOOKUP, false, lookupAirmarCalibrationStatus, "", 0},
		{"Verify Score", 8, RES_INTEGER, false, nil, "TBD", 0},
		{"X-axis gain value", 16, 0.01, false, nil, "default 100, range 50 to 500", 0},
		{"Y-axis gain value", 16, 0.01, false, nil, "default 100, range 50 to 500", 0},
		{"Z-axis gain value", 16, 0.01, false, nil, "default 100, range 50 to 500", 0},
		{"X-axis linear offset", 16, 0.01, true, "Tesla", "default 0, range -320.00 to 320.00", 0},
		{"Y-axis linear offset", 16, 0.01, true, "Tesla", "default 0, range -320.00 to 320.00", 0},
		{"Z-axis linear offset", 16, 0.01, true, "Tesla", "default 0, range -320.00 to 320.00", 0},
		{"X-axis angular offset", 16, 0.1, false, "deg", "default 0, range 0 to 3600", 0},
		{"Pitch and Roll damping", 16, 0.05, false, "s", "default 30, range 0 to 200", 0},
		{"Compass/Rate gyro damping", 16, 0.05, true, "s", "default -30, range -2400 to 2400, negative indicates rate gyro is to be used in compass calculations", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/PB200UserManual.pdf
	{"Airmar: True Wind Options", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=34", "True Wind Options", 0},
		{"COG substition for HDG", 2, RES_LOOKUP, false, lookupAirmarCogSubstitute, "Allow use of COG when HDG not available?", 0},
		{"Calibration Status", 8, RES_LOOKUP, false, lookupAirmarCalibrationStatus, "", 0},
		{"Verify Score", 8, RES_INTEGER, false, nil, "TBD", 0},
		{"X-axis gain value", 16, 0.01, false, nil, "default 100, range 50 to 500", 0},
		{"Y-axis gain value", 16, 0.01, false, nil, "default 100, range 50 to 500", 0},
		{"Z-axis gain value", 16, 0.01, false, nil, "default 100, range 50 to 500", 0},
		{"X-axis linear offset", 16, 0.01, true, "Tesla", "default 0, range -320.00 to 320.00", 0},
		{"Y-axis linear offset", 16, 0.01, true, "Tesla", "default 0, range -320.00 to 320.00", 0},
		{"Z-axis linear offset", 16, 0.01, true, "Tesla", "default 0, range -320.00 to 320.00", 0},
		{"X-axis angular offset", 16, 0.1, false, "deg", "default 0, range 0 to 3600", 0},
		{"Pitch and Roll damping", 16, 0.05, false, "s", "default 30, range 0 to 200", 0},
		{"Compass/Rate gyro damping", 16, 0.05, true, "s", "default -30, range -2400 to 2400, negative indicates rate gyro is to be used in compass calculations", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Simulate Mode", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=35", "Simulate Mode", 0},
		{"Simulate Mode", 2, RES_LOOKUP, false, lookupOffOn, "", 0},
		{"Reserved", 22, RES_BINARY, false, nil, "Reserved", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Calibrate Depth", 126720, true, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=40", "Calibrate Depth", 0},
		{"Speed of Sound Mode", 16, 0.1, false, "m/s", "actual allowed range is 1350.0 to 1650.0 m/s", 0},
		{"Reserved", 8, RES_BINARY, false, nil, "Reserved", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Calibrate Speed", 126720, true, 8, 2, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=41", "Calibrate Speed", 0},
		{"Number of pairs of data points", 8, RES_INTEGER, false, nil, "actual range is 0 to 25. 254=restore default speed curve", 0},
		{"Input frequency", 16, 0.1, false, "Hz", "", 0},
		{"Output speed", 16, 0.01, false, "m/s", "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Calibrate Temperature", 126720, true, 8, 2, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, "=42", "Calibrate Temperature", 0},
		{"Temperature Instance", 2, RES_LOOKUP, false, lookupAirmarTempInstance, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "Reserved", 0},
		{"Temperature offset", 16, 0.1, false, "Hz", "", 0},
		{"Temperature offset", 16, 0.001, true, "K", "actual range is -9.999 to +9.999 K", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: Speed Filter", 126720, true, 8, 2, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, "=43", "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, nil, "Speed Filter", 0},
		{"Filter type", 4, RES_LOOKUP, false, lookupAirmarSpeedFilter, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Sample interval", 16, 0.01, false, "s", "", 0},
		{"Filter duration", 16, 0.01, false, "s", "", 0}},
	},

	// http://www.airmartechnology.com/uploads/installguide/DST200UserlManual.pdf
	{"Airmar: NMEA 2000 options", 126720, true, 8, 2, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, "=46", "Marine Industry", 0},
		{"Proprietary ID", 8, RES_INTEGER, false, nil, "NMEA 2000 options", 0},
		{"Transmission Interval", 2, RES_LOOKUP, false, lookupAirmarTransmissionInterval, "", 0},
		{"Reserved", 22, RES_BINARY, false, nil, "Reserved", 0}},
	},

	// http://www.maretron.com/support/manuals/GPS100UM_1.2.pdf
	{"System Time", 126992, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Source", 4, RES_LOOKUP, false, lookupSystemTime, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0}},
	},

	{"Product Information", 126996, false, 0x86, 0, []Field{
		{"NMEA 2000 Version", 16, 1, false, nil, "", 0},
		{"Product Code", 16, 1, false, nil, "", 0},
		{"Model ID", 256, RES_ASCII, false, nil, "", 0},
		{"Software Version Code", 320, RES_ASCII, false, nil, "", 0},
		{"Model Version", 192, RES_ASCII, false, nil, "", 0},
		{"Model Serial Code", 256, RES_ASCII, false, nil, "", 0},
		{"Certification Level", 8, 1, false, nil, "", 0},
		{"Load Equivalency", 8, 1, false, nil, "", 0}},
	},

	{"Configuration Information", 126998, false, 0x2a, 0, []Field{
		{"Station ID", 16, 1, false, nil, "", 0},
		{"Station Name", 16, 1, false, nil, "", 0},
		{"A", 16, 1, false, nil, "", 0},
		{"Manufacturer", 288, RES_ASCII, false, nil, "", 0},
		{"Installation Description #1", 16, 1, false, nil, "", 0},
		{"Installation Description #2", 16, 1, false, nil, "", 0}},
	},

	////////////////////////// PERIODIC DATA PGNs //////////////////////////////
	// http://www.nmea.org/Assets/july%202010%20nmea2000_v1-301_app_b_pgn_field_list.pdf
	// http://www.maretron.com/support/manuals/USB100UM_1.2.pdf
	// http://www8.garmin.com/manuals/GPSMAP4008_NMEA2000NetworkFundamentals.pdf

	{"Heading/Track control", 127237, false, 0x15, 0, []Field{
		{"Rudder Limit Exceeded", 2, 1, false, nil, "", 0},
		{"Off-Heading Limit Exceeded", 2, 1, false, nil, "", 0},
		{"Off-Track Limit Exceeded", 2, 1, false, nil, "", 0},
		{"Override", 2, 1, false, nil, "", 0},
		{"Steering Mode", 4, 1, false, nil, "", 0},
		{"Turn Mode", 4, 1, false, nil, "", 0},
		{"Heading Reference", 3, 1, false, nil, "", 0},
		{"Reserved", 3, 1, false, nil, "", 0},
		{"Commanded Rudder Direction", 2, 1, false, nil, "", 0},
		{"Commanded Rudder Angle", 16, RES_DEGREES, true, "deg", "", 0},
		{"Heading-To-Steer (Course)", 16, RES_DEGREES, false, "deg", "", 0},
		{"Track", 16, RES_DEGREES, false, "deg", "", 0},
		{"Rudder Limit", 16, RES_DEGREES, false, "deg", "", 0},
		{"Off-Heading Limit", 16, RES_DEGREES, false, "deg", "", 0},
		{"Radius of Turn Order", 16, RES_DEGREES, true, "deg", "", 0},
		{"Rate of Turn Order", 16, RES_ROTATION, true, "deg/s", "", 0},
		{"Off-Track Limit", 16, 1, true, "m", "", 0},
		{"Vessel Heading", 16, RES_DEGREES, false, "deg", "", 0}},
	},

	// http://www.maretron.com/support/manuals/RAA100UM_1.0.pdf
	// Haven't actually seen this value yet, lengths are guesses
	{"Rudder", 127245, false, 8, 0, []Field{
		{"Instance", 8, 1, false, nil, "", 0},
		{"Direction Order", 2, 1, false, nil, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "Reserved", 0},
		{"Angle Order", 16, RES_DEGREES, true, "deg", "", 0},
		{"Position", 16, RES_DEGREES, true, "deg", "", 0}},
	},

	// NMEA + Simrad AT10
	// http://www.maretron.com/support/manuals/SSC200UM_1.7.pdf
	// molly_rose_E80start.kees
	{"Vessel Heading", 127250, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Heading", 16, RES_DEGREES, false, "deg", "", 0},
		{"Deviation", 16, RES_DEGREES, true, "deg", "", 0},
		{"Variation", 16, RES_DEGREES, true, "deg", "", 0},
		{"Reference", 2, RES_LOOKUP, false, lookupDirectionReference, "", 0}},
	},

	// http://www.maretron.com/support/manuals/SSC200UM_1.7.pdf
	// Lengths observed from Simrad RC42
	{"Rate of Turn", 127251, true, 5, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Rate", 32, RES_ROTATION * 0.0001, true, "deg/s", "", 0}},
	},

	{"Attitude", 127257, true, 7, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Yaw", 16, RES_ROTATION, true, "deg/s", "", 0},
		{"Pitch", 16, RES_ROTATION, true, "deg/s", "", 0},
		{"Roll", 16, RES_ROTATION, true, "deg/s", "", 0}},
	},

	// NMEA + Simrad AT10
	// http://www.maretron.com/support/manuals/GPS100UM_1.2.pdf
	{"Magnetic Variation", 127258, true, 6, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Source", 4, RES_LOOKUP, false, lookupMagneticVariation, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Age of service", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Variation", 16, RES_DEGREES, true, "deg", "", 0}},
	},

	// Engine group PGNs all derived PGN Numbers from
	// http://www.maretron.com/products/pdf/J2K100-Data_Sheet.pdf
	// http://www.floscan.com/html/blue/NMEA2000.php

	{"Engine Parameters, Rapid Update", 127488, true, 8, 0, []Field{
		{"Engine Instance", 8, RES_LOOKUP, false, lookupEngineInstance, "", 0},
		{"Engine Speed", 16, RES_INTEGER, false, "rpm", "", 0},
		{"Engine Boost Pressure", 16, RES_PRESSURE, false, "hPa", "", 0},
		{"Engine Tilt/Trim", 8, 1, true, nil, "", 0}},
	},

	{"Engine Parameters, Dynamic", 127489, true, 26, 0, []Field{
		{"Engine Instance", 8, RES_LOOKUP, false, lookupEngineInstance, "", 0},
		{"Oil pressure", 16, RES_PRESSURE, false, "hPa", "", 0},
		{"Oil temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Alternator Potential", 16, 0.01, false, "V", "", 0},
		{"Fuel Rate", 16, 0.1, true, "L/h", "", 0},
		{"Total Engine hours", 32, 1.0, false, "s", "", 0},
		{"Coolant Pressure", 16, RES_PRESSURE, false, "hPa", "", 0},
		{"Fuel Pressure", 16, 1, false, nil, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Discrete Status 1", 16, RES_INTEGER, false, nil, "", 0},
		{"Discrete Status 2", 16, RES_INTEGER, false, nil, "", 0},
		{"Percent Engine Load", 8, RES_INTEGER, true, "%", "", 0},
		{"Percent Engine Torque", 8, RES_INTEGER, true, "%", "", 0}},
	},

	{"Transmission Parameters, Dynamic", 127493, true, 7, 0, []Field{
		{"Engine Instance", 2, RES_LOOKUP, false, lookupEngineInstance, "", 0},
		{"Transmission Gear", 2, RES_LOOKUP, false, lookupGearStatus, "", 0},
		{"Reserved", 4, 1, false, nil, "", 0},
		{"Oil pressure", 16, RES_PRESSURE, false, "hPa", "", 0},
		{"Oil temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Discrete Status 1", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"Trip Parameters, Vessel", 127496, true, 10, 0, []Field{
		{"Time to Empty", 32, 0.001, false, "s", "", 0},
		{"Distance to Empty", 32, 0.01, false, "m", "", 0},
		{"Estimated Fuel Remaining", 16, 1, false, "L", "", 0},
		{"Trip Run Time", 32, 0.001, false, "s", "", 0}},
	},

	{"Trip Parameters, Engine", 127497, true, 9, 0, []Field{
		{"Engine Instance", 8, RES_LOOKUP, false, lookupEngineInstance, "", 0},
		{"Trip Fuel Used", 16, 1, false, "L", "", 0},
		{"Fuel Rate, Average", 16, 0.1, true, "L/h", "", 0},
		{"Fuel Rate, Economy", 16, 0.1, true, "L/h", "", 0},
		{"Instantaneous Fuel Economy", 16, 0.1, true, "L/h", "", 0}},
	},

	{"Engine Parameters, Static", 127498, true, 8, 0, []Field{
		{"Engine Instance", 8, RES_LOOKUP, false, lookupEngineInstance, "", 0},
		{"Rated Engine Speed", 16, 1, false, nil, "", 0},
		{"VIN", 8, 1, false, nil, "", 0},
		{"Software ID", 16, 1, false, nil, "", 0}},
	},

	{"Binary Switch Bank Status", 127501, false, 8, 1, []Field{
		{"Indicator Bank Instance", 8, 1, false, nil, "", 0},
		{"Indicator", 2, RES_LOOKUP, false, lookupOffOn, "", 0}},
	},

	{"Switch Bank Control", 127502, false, 8, 1, []Field{
		{"Switch Bank Instance", 8, 1, false, nil, "", 0},
		{"Switch", 2, RES_LOOKUP, false, lookupOffOn, "", 0}},
	},

	// http://www.nmea.org/Assets/nmea-2000-corrigendum-1-2010-1.pdf
	{"AC Input Status", 127503, true, 8, 10, []Field{
		{"AC Instance", 8, 1, false, nil, "", 0},
		{"Number of Lines", 8, 1, false, nil, "", 0},
		{"Line", 2, RES_LOOKUP, false, lookupLine, "", 0},
		{"Acceptability", 2, RES_LOOKUP, false, lookupAcceptability, "", 0},
		{"Reserved", 4, 1, false, nil, "", 0},
		{"Voltage", 16, 0.01, false, "V", "", 0},
		{"Current", 16, 0.1, false, "A", "", 0},
		{"Frequency", 16, 0.01, false, "Hz", "", 0},
		{"Breaker Size", 16, 0.1, false, "A", "", 0},
		{"Real Power", 32, RES_INTEGER, false, "W", "", 0},
		{"Reactive Power", 32, RES_INTEGER, false, "VAr", "", 0},
		{"Power Factor", 8, 0.01, false, nil, "", 0}},
	},

	// http://www.nmea.org/Assets/nmea-2000-corrigendum-1-2010-1.pdf
	{"AC Output Status", 127504, true, 8, 10, []Field{
		{"AC Instance", 8, 1, false, nil, "", 0},
		{"Number of Lines", 8, 1, false, nil, "", 0},
		{"Line", 2, RES_LOOKUP, false, lookupLine, "", 0},
		{"Waveform", 3, RES_LOOKUP, false, lookupWaveform, "", 0},
		{"Reserved", 3, 1, false, nil, "", 0},
		{"Voltage", 16, 0.01, false, "V", "", 0},
		{"Current", 16, 0.1, false, "A", "", 0},
		{"Frequency", 16, 0.01, false, "Hz", "", 0},
		{"Breaker Size", 16, 0.1, false, "A", "", 0},
		{"Real Power", 32, RES_INTEGER, false, "W", "", 0},
		{"Reactive Power", 32, RES_INTEGER, false, "VAr", "", 0},
		{"Power Factor", 8, 0.01, false, nil, "", 0}},
	},

	// http://www.maretron.com/support/manuals/TLA100UM_1.2.pdf
	// Observed from EP65R
	{"Fluid Level", 127505, true, 7, 0, []Field{
		{"Instance", 4, 1, false, nil, "", 0},
		{"Type", 4, RES_LOOKUP, false, lookupTankType, "", 0},
		{"Level", 16, 100.0 / 25000, false, "%", "", 0},
		{"Capacity", 32, 0.1, false, "L", "", 0}},
	},

	{"DC Detailed Status", 127506, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"DC Instance", 8, 1, false, nil, "", 0},
		{"DC Type", 8, 1, false, nil, "", 0},
		{"State of Charge", 8, 1, false, nil, "", 0},
		{"State of Health", 8, 1, false, nil, "", 0},
		{"Time Remaining", 16, 1, false, nil, "", 0},
		{"Ripple Voltage", 16, 0.01, false, "V", "", 0}},
	},

	{"Charger Status", 127507, false, 8, 0, []Field{
		{"Charger Instance", 8, 1, false, nil, "", 0},
		{"Battery Instance", 8, 1, false, nil, "", 0},
		{"Operating State", 8, 1, false, nil, "", 0},
		{"Charge Mode", 8, 1, false, nil, "", 0},
		{"Charger Enable/Disable", 2, 1, false, nil, "", 0},
		{"Equalization Pending", 2, 1, false, nil, "", 0},
		{"Reserved", 4, 1, false, nil, "", 0},
		{"Equalization Time Remaining", 16, 1, false, nil, "", 0}},
	},

	{"Battery Status", 127508, true, 8, 0, []Field{
		{"Battery Instance", 8, 1, false, nil, "", 0},
		{"Voltage", 16, 0.01, true, "V", "", 0},
		{"Current", 16, 0.1, true, "A", "", 0},
		{"Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"SID", 8, 1, false, nil, "", 0}},
	},

	{"Inverter Status", 127509, false, 8, 0, []Field{
		{"Inverter Instance", 8, 1, false, nil, "", 0},
		{"AC Instance", 8, 1, false, nil, "", 0},
		{"DC Instance", 8, 1, false, nil, "", 0},
		{"Operating State", 4, RES_LOOKUP, false, lookupStandbyOn, "", 0},
		{"Inverter", 2, RES_LOOKUP, false, lookupStandbyOn, "", 0}},
	},

	{"Charger Configuration Status", 127510, false, 8, 0, []Field{
		{"Charger Instance", 8, 1, false, nil, "", 0},
		{"Battery Instance", 8, 1, false, nil, "", 0},
		{"Charger Enable/Disable", 2, 1, false, nil, "", 0},
		{"Reserved", 6, 1, false, nil, "", 0},
		{"Charge Current Limit", 16, 0.1, false, "A", "", 0},
		{"Charging Algorithm", 8, 1, false, nil, "", 0},
		{"Charger Mode", 8, 1, false, nil, "", 0},
		{"Estimated Temperature", 16, RES_TEMPERATURE, false, "K", "When no sensor present", 0},
		{"Equalize One Time Enable/Disable", 4, 1, false, nil, "", 0},
		{"Over Charge Enable/Disable", 4, 1, false, nil, "", 0},
		{"Equalize Time", 16, 1, false, nil, "", 0}},
	},

	{"Inverter Configuration Status", 127511, false, 8, 0, []Field{
		{"Inverter Instance", 8, 1, false, nil, "", 0},
		{"AC Instance", 8, 1, false, nil, "", 0},
		{"DC Instance", 8, 1, false, nil, "", 0},
		{"Inverter Enable/Disable", 2, 1, false, nil, "", 0},
		{"Inverter Mode", 8, 1, false, nil, "", 0},
		{"Load Sense Enable/Disable", 8, 1, false, nil, "", 0},
		{"Load Sense Power Threshold", 8, 1, false, nil, "", 0},
		{"Load Sense Interval", 8, 1, false, nil, "", 0}},
	},

	{"AGS Configuration Status", 127512, false, 8, 0, []Field{
		{"AGS Instance", 8, 1, false, nil, "", 0},
		{"Generator Instance", 8, 1, false, nil, "", 0},
		{"AGS Mode", 8, 1, false, nil, "", 0}},
	},

	{"Battery Configuration Status", 127513, false, 8, 0, []Field{
		{"Battery Instance", 8, 1, false, nil, "", 0},
		{"Battery Type", 8, 1, false, nil, "", 0},
		{"Supports Equalization", 2, 1, false, nil, "", 0},
		{"Reserved", 6, 1, false, nil, "", 0},
		{"Nominal Voltage", 16, 0.01, false, "V", "", 0},
		{"Chemistry", 8, 1, false, nil, "", 0},
		{"Capacity", 16, 1, false, nil, "", 0},
		{"Temperature Coefficient", 16, 1, false, nil, "", 0},
		{"Peukert Exponent", 16, 1, false, nil, "", 0},
		{"Charge Efficiency Factor", 16, 1, false, nil, "", 0}},
	},

	{"AGS Status", 127514, false, 8, 0, []Field{
		{"AGS Instance", 8, 1, false, nil, "", 0},
		{"Generator Instance", 8, 1, false, nil, "", 0},
		{"AGS Operating State", 8, 1, false, nil, "", 0},
		{"Generator State", 8, 1, false, nil, "", 0},
		{"Generator On Reason", 8, 1, false, nil, "", 0},
		{"Generator Off Reason", 8, 1, false, nil, "", 0}},
	},

	// http://www.maretron.com/support/manuals/DST100UM_1.2.pdf
	{"Speed", 128259, true, 6, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Speed Water Referenced", 16, 0.01, false, "m/s", "", 0},
		{"Speed Ground Referenced", 16, 0.01, false, "m/s", "", 0},
		{"Speed Water Referenced Type", 4, RES_LOOKUP, false, nil, "", 0}},
	},

	// http://www.maretron.com/support/manuals/DST100UM_1.2.pdf
	{"Water Depth", 128267, true, 5, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Depth", 32, 0.01, false, "m", "Depth below transducer", 0},
		{"Offset", 16, 0.001, true, "m", "Distance between transducer and surface (positive) or keel (negative)", 0}},
	},

	// http://www.nmea.org/Assets/nmea-2000-digital-interface-white-paper.pdf
	{"Distance Log", 128275, true, 14, 0, []Field{
		{"Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Log", 32, 1, false, "m", "Total cumulative distance", 0},
		{"Trip Log", 32, 1, false, "m", "Distance since last reset", 0}},
	},

	{"Tracked Target Data", 128520, true, 27, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Target ID #", 8, 1, false, nil, "Number of route, waypoint, event, mark, etc.", 0},
		{"Track Status", 2, RES_LOOKUP, false, lookupTrackStatus, "", 0},
		{"Reported Target", 1, RES_LOOKUP, false, lookupYesNo, "", 0},
		{"Target Acquisition", 1, RES_LOOKUP, false, lookupTargetAcquisition, "", 0},
		{"Bearing Reference", 2, RES_LOOKUP, false, lookupDirectionReference, "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Bearing", 16, RES_DEGREES, false, "deg", "", 0},
		{"Distance", 32, 0.001, false, "m", "", 0},
		{"Course", 16, RES_DEGREES, false, "deg", "", 0},
		{"Speed", 16, 0.01, false, "m/s", "", 0},
		{"CPA", 32, 0.01, false, "m", "", 0},
		{"TCPA", 32, 0.001, true, "s", "negative = time elapsed since event, positive = time to go", 0},
		{"UTC of Fix", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Name", 2040, RES_ASCII, false, nil, "", 0}},
	},

	{"Position, Rapid Update", 129025, true, 8, 0, []Field{
		{"Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Longitude", 32, RES_LONGITUDE, true, "deg", "", 0}},
	},

	// http://www.maretron.com/support/manuals/GPS100UM_1.2.pdf
	{"COG & SOG, Rapid Update", 129026, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"COG Reference", 2, RES_LOOKUP, false, lookupDirectionReference, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "Reserved", 0},
		{"COG", 16, RES_DEGREES, false, "deg", "", 0},
		{"SOG", 16, 0.01, false, "m/s", "", 0},
		{"Reserved", 16, RES_BINARY, false, nil, "Reserved", 0}},
	},

	{"Position Delta, Rapid Update", 129027, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Time Delta", 16, 1, false, nil, "", 0},
		{"Latitude Delta", 16, 1, true, nil, "", 0},
		{"Longitude Delta", 16, 1, true, nil, "", 0}},
	},

	{"Altitude Delta, Rapid Update", 129028, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Time Delta", 16, 1, true, nil, "", 0},
		{"GNSS Quality", 2, 1, false, nil, "", 0},
		{"Direction", 2, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Course Over Ground", 32, RES_DEGREES, false, "deg", "", 0},
		{"Altitude Delta", 16, 1, true, nil, "", 0}},
	},

	// http://www.maretron.com/support/manuals/GPS100UM_1.2.pdf
	{"GNSS Position Data", 129029, true, 51, 3, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Latitude", 64, RES_LATITUDE, true, "deg", "", 0},
		{"Longitude", 64, RES_LONGITUDE, true, "deg", "", 0},
		{"Altitude", 64, 1e-6, true, "m", "", 0},
		{"GNSS type", 4, RES_LOOKUP, false, lookupGns, "", 0},
		{"Method", 4, RES_LOOKUP, false, lookupGnsMethod, "", 0},
		{"Integrity", 2, RES_LOOKUP, false, lookupGnsIntegrity, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "Reserved", 0},
		{"Number of SVs", 8, 1, false, nil, "Number of satellites used in solution", 0},
		{"HDOP", 16, 0.01, true, nil, "Horizontal dilution of precision", 0},
		{"PDOP", 16, 0.01, true, nil, "Probable dilution of precision", 0},
		{"Geoidal Separation", 16, 0.01, false, "m", "Geoidal Separation", 0},
		{"Reference Stations", 8, 1, false, nil, "Number of reference stations", 0},
		{"Reference Station Type", 4, RES_LOOKUP, false, lookupGns, "", 0},
		{"Reference Station ID", 12, 1, false, nil, "", 0},
		{"Age of DGNSS Corrections", 16, 0.01, false, "s", "", 0}},
	},

	{"Time & Date", 129033, true, 8, 0, []Field{
		{"Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Time", 32, RES_TIME, false, "seconds", "Seconds since midnight", 0},
		{"Local Offset", 16, RES_INTEGER, true, "minutes", "Minutes", 0}},
	},

	{"AIS Class A Position Report", 129038, true, 27, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Longitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Latitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Position Accuracy", 1, RES_LOOKUP, false, lookupPositionAccuracy, "", 0},
		{"RAIM", 1, RES_LOOKUP, false, lookupAisRAIM, "", 0},
		{"Time Stamp", 6, RES_LOOKUP, false, lookupTimeStamp, "0-59 = UTC second when the report was generated", 0},
		{"COG", 16, RES_DEGREES, false, "deg", "", 0},
		{"SOG", 16, 0.01, false, "m/s", "", 0},
		{"Communication State", 19, RES_BINARY, false, nil, "Information used by the TDMA slot allocation algorithm and synchronization information", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Heading", 16, RES_DEGREES, false, "deg", "True heading", 0},
		{"Rate of Turn", 16, RES_ROTATION, false, "deg/s", "", 0},
		{"Nav Status", 8, RES_LOOKUP, false, lookupNavStatus, "", 0},
		{"Reserved for Regional Applications", 8, 1, false, nil, "", 0},
		{"Spare", 8, 1, false, nil, "", 0}},
	},

	{"AIS Class B Position Report", 129039, true, 0x1a, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Longitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Latitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Position Accuracy", 1, RES_LOOKUP, false, lookupPositionAccuracy, "", 0},
		{"RAIM", 1, RES_LOOKUP, false, lookupAisRAIM, "", 0},
		{"Time Stamp", 6, RES_LOOKUP, false, lookupTimeStamp, "0-59 = UTC second when the report was generated", 0},
		{"COG", 16, RES_DEGREES, false, "deg", "", 0},
		{"SOG", 16, 0.01, false, "m/s", "", 0},
		{"Communication State", 19, RES_BINARY, false, nil, "Information used by the TDMA slot allocation algorithm and synchronization information", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Heading", 16, RES_DEGREES, false, "deg", "True heading", 0},
		{"Regional Application", 8, 1, false, nil, "", 0},
		{"Regional Application", 2, 1, false, nil, "", 0},
		{"Unit type", 1, RES_LOOKUP, false, lookupAisUnitType, "", 0},
		{"Integrated Display", 1, RES_LOOKUP, false, lookupYesNo, "Whether the unit can show messages 12 and 14", 0},
		{"DSC", 1, RES_LOOKUP, false, lookupYesNo, "", 0},
		{"Band", 1, RES_LOOKUP, false, lookupAisBand, "", 0},
		{"Can handle Msg 22", 1, RES_LOOKUP, false, lookupYesNo, "Whether device supports message 22", 0},
		{"AIS mode", 1, RES_LOOKUP, false, lookupAisMode, "", 0},
		{"AIS communication state", 1, RES_LOOKUP, false, lookupAisCommState, "", 0}},
	},

	{"AIS Class B Extended Position Report", 129040, true, 33, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Longitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Latitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Position Accuracy", 1, RES_LOOKUP, false, lookupPositionAccuracy, "", 0},
		{"AIS RAIM flag", 1, RES_LOOKUP, false, lookupAisRAIM, "", 0},
		{"Time Stamp", 6, RES_LOOKUP, false, lookupTimeStamp, "0-59 = UTC second when the report was generated", 0},
		{"COG", 16, RES_DEGREES, false, "deg", "", 0},
		{"SOG", 16, 0.01, false, "m/s", "", 0},
		{"Regional Application", 8, 1, false, nil, "", 0},
		{"Regional Application", 4, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "reserved", 0},
		{"Type of ship", 8, RES_LOOKUP, false, lookupShipType, "", 0},
		{"True Heading", 16, RES_DEGREES, false, "deg", "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0},
		{"GNSS type", 4, RES_LOOKUP, false, lookupGnsAis, "", 0},
		{"Length", 16, 0.1, false, "m", "", 0},
		{"Beam", 16, 0.1, false, "m", "", 0},
		{"Position reference from Starboard", 16, 0.1, false, "m", "", 0},
		{"Position reference from Bow", 16, 0.1, false, "m", "", 0},
		{"Name", 160, RES_ASCII, false, nil, "0=unavailable", 0},
		{"DTE", 1, RES_LOOKUP, false, lookupAisDTE, "", 0},
		{"AIS mode", 1, 1, false, lookupAisMode, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0}},
	},

	{"Datum", 129044, true, 24, 0, []Field{
		{"Local Datum", 32, RES_ASCII, false, nil, "defined in IHO Publication S-60, Appendices B and C. " +
			"First three chars are datum ID as per IHO tables." +
			"Fourth char is local datum subdivision code.", 0},
		{"Delta Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Delta Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Delta Altitude", 32, 1e-6, true, "m", "", 0},
		{"Reference Datum", 32, RES_ASCII, false, nil, "defined in IHO Publication S-60, Appendices B and C. " +
			"First three chars are datum ID as per IHO tables. " +
			"Fourth char is local datum subdivision code.", 0}},
	},

	{"User Datum", 129045, true, 37, 0, []Field{
		{"Delta X", 32, 0.01, true, "m", "Delta shift in X axis from WGS 84", 0},
		{"Delta Y", 32, 0.01, true, "m", "Delta shift in Y axis from WGS 84", 0},
		{"Delta Z", 32, 0.01, true, "m", "Delta shift in Z axis from WGS 84", 0},
		{"Rotation in X", 32, RES_FLOAT, true, nil, "Rotational shift in X axis from WGS 84. Rotations presented use the geodetic sign convention.  When looking along the positive axis towards the origin, counter-clockwise rotations are positive.", 0},
		{"Rotation in Y", 32, RES_FLOAT, true, nil, "Rotational shift in Y axis from WGS 84. Rotations presented use the geodetic sign convention.  When looking along the positive axis towards the origin, counter-clockwise rotations are positive.", 0},
		{"Rotation in Z", 32, RES_FLOAT, true, nil, "Rotational shift in Z axis from WGS 84. Rotations presented use the geodetic sign convention.  When looking along the positive axis towards the origin, counter-clockwise rotations are positive.", 0},
		{"Scale", 32, RES_FLOAT, true, "ppm", "Scale factor expressed in parts-per-million", 0},
		{"Ellipsoid Semi-major Axis", 32, 0.01, true, "m", "Semi-major axis (a) of the User Datum ellipsoid", 0},
		{"Ellipsoid Flattening Inverse", 32, RES_FLOAT, true, nil, "Flattening (1/f) of the User Datum ellipsoid", 0},
		{"Datum Name", 32, RES_ASCII, false, nil, "4 character code from IHO Publication S-60,Appendices B and C." +
			"First three chars are datum ID as per IHO tables." +
			"Fourth char is local datum subdivision code.", 0}},
	},

	{"Cross Track Error", 129283, false, 6, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"XTE mode", 4, RES_LOOKUP, false, lookupResidualMode, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Navigation Terminated", 2, RES_LOOKUP, false, lookupYesNo, "", 0},
		{"XTE", 32, 0.01, true, "m", "", 0}},
	},

	{"Navigation Data", 129284, true, 0x22, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Distance to Waypoint", 32, 0.01, false, "m", "", 0},
		{"Course/Bearing reference", 2, RES_LOOKUP, false, lookupDirectionReference, "", 0},
		{"Perpendicular Crossed", 2, RES_LOOKUP, false, lookupYesNo, "", 0},
		{"Arrival Circle Entered", 2, RES_LOOKUP, false, lookupYesNo, "", 0},
		{"Calculation Type", 2, RES_LOOKUP, false, lookupNavCalculation, "", 0},
		{"ETA Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"ETA Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Bearing, Origin to Destination Waypoint", 16, RES_DEGREES, false, "deg", "", 0},
		{"Bearing, Position to Destination Waypoint", 16, RES_DEGREES, false, "deg", "", 0},
		{"Origin Waypoint Number", 32, 1, false, nil, "", 0},
		{"Destination Waypoint Number", 32, 1, false, nil, "", 0},
		{"Destination Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Destination Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Waypoint Closing Velocity", 16, 0.01, true, "m/s", "", 0}},
	},

	{"Navigation - Route/WP Information", 129285, true, 8, 4, []Field{
		{"Start RPS#", 16, 1, false, nil, "", 0},
		{"nItems", 16, 1, false, nil, "", 0},
		{"Database ID", 16, 1, false, nil, "", 0},
		{"Route ID", 16, 1, false, nil, "", 0},
		{"Navigation direction in route", 2, 1, false, nil, "", 0},
		{"Supplementary Route/WP data available", 2, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Route Name", 2040, RES_STRING, false, nil, "", 0},
		{"Reserved", 8, RES_BINARY, false, nil, "Reserved", 0},
		{"WP ID", 16, 1, false, nil, "", 0},
		{"WP Name", 2040, RES_STRING, false, nil, "", 0},
		{"WP Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"WP Longitude", 32, RES_LONGITUDE, true, "deg", "", 0}},
	},

	{"Set & Drift, Rapid Update", 129291, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Set Reference", 2, RES_LOOKUP, false, lookupDirectionReference, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "Reserved", 0},
		{"Set", 16, RES_DEGREES, false, "deg", "", 0},
		{"Drift", 16, 0.01, false, "m/s", "", 0}},
	},

	{"Navigation - Route / Time to+from Mark", 129301, true, 10, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Time to mark", 32, 0.001, true, "s", "negative = elapsed since event, positive = time to go", 0},
		{"Mark Type", 4, RES_LOOKUP, false, lookupNavMarkType, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Mark ID", 32, 1, false, nil, "", 0}},
	},

	{"Bearing and Distance between two Marks", 129302, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Bearing Reference", 4, RES_LOOKUP, false, nil, "", 0},
		{"Calculation Type", 2, RES_LOOKUP, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "Reserved", 0},
		{"Bearing, Origin to Destination", 16, RES_DEGREES, false, "deg", "", 0},
		{"Distance", 32, 0.01, false, "m", "", 0},
		{"Origin Mark Type", 4, RES_LOOKUP, false, nil, "", 0},
		{"Destination Mark Type", 4, RES_LOOKUP, false, nil, "", 0},
		{"Origin Mark ID", 32, 1, false, nil, "", 0},
		{"Destination Mark ID", 32, 1, false, nil, "", 0}},
	},

	// http://www.maretron.com/support/manuals/GPS100UM_1.2.pdf
	// Haven't seen this yet (no way to send PGN 059904 yet) so lengths unknown
	{"GNSS Control Status", 129538, false, 8, 0, []Field{
		{"SV Elevation Mask", 16, 1, false, nil, "Will not use SV below this elevation", 0},
		{"PDOP Mask", 16, 0.01, false, nil, "Will not report position above this PDOP", 0},
		{"PDOP Switch", 16, 0.01, false, nil, "Will report 2D position above this PDOP", 0},
		{"SNR Mask", 16, 0.01, false, nil, "Will not use SV below this SNR", 0},
		{"GNSS Mode (desired)", 3, RES_LOOKUP, false, lookupGnssMode, "", 0},
		{"DGNSS Mode (desired)", 3, RES_LOOKUP, false, lookupDGnssMode, "", 0},
		{"Position/Velocity Filter", 2, 1, false, nil, "", 0},
		{"Max Correction Age", 16, 1, false, nil, "", 0},
		{"Antenna Altitude for 2D Mode", 16, 0.01, false, "m", "", 0},
		{"Use Antenna Altitude for 2D Mode", 2, RES_LOOKUP, false, lookupGnssAntenna, "", 0}},
	},

	// http://www.maretron.com/support/manuals/GPS100UM_1.2.pdf
	{"GNSS DOPs", 129539, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Desired Mode", 3, RES_LOOKUP, false, lookupGnssMode, "", 0},
		{"Actual Mode", 3, RES_LOOKUP, false, lookupGnssMode, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "Reserved", 0},
		{"HDOP", 16, 0.01, true, nil, "Horizontal dilution of precision", 0},
		{"VDOP", 16, 0.01, true, nil, "Vertical dilution of precision", 0},
		{"TDOP", 16, 0.01, true, nil, "Time dilution of precision", 0}},
	},

	{"GNSS Sats in View", 129540, true, 0xff, 7, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Mode", 2, RES_LOOKUP, false, lookupGnssSatMode, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "Reserved", 0},
		{"Sats in View", 8, 1, false, nil, "", 0},
		{"PRN", 8, 1, false, nil, "", 0},
		{"Elevation", 16, RES_DEGREES, true, "deg", "", 0},
		{"Azimuth", 16, RES_DEGREES, true, "deg", "", 0},
		{"SNR", 16, 0.01, false, "dB", "", 0},
		{"Range residuals", 32, 1, true, nil, "", 0},
		{"Status", 4, RES_LOOKUP, false, lookupGnssSatStatus, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0}},
	},

	{"GPS Almanac Data", 129541, false, 8, 0, []Field{
		{"PRN", 8, 1, false, nil, "", 0},
		{"GPS Week number", 8, 1, false, nil, "", 0},
		{"SV Health Bits", 8, 1, false, nil, "", 0},
		{"Eccentricity", 8, 1, false, nil, "", 0},
		{"Almanac Reference Time", 8, 1, false, nil, "", 0},
		{"Inclination Angle", 8, 1, false, nil, "", 0},
		{"Right of Right Ascension", 8, 1, false, nil, "", 0},
		{"Root of Semi-major Axis", 8, 1, false, nil, "", 0},
		{"Argument of Perigee", 8, 1, false, nil, "", 0},
		{"Longitude of Ascension Node", 8, 1, false, nil, "", 0},
		{"Mean Anomaly", 8, 1, false, nil, "", 0},
		{"Clock Parameter 1", 8, 1, false, nil, "", 0},
		{"Clock Parameter 2", 8, 1, false, nil, "", 0}},
	},

	{"GNSS Pseudorange Noise Statistics", 129542, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"RMS of Position Uncertainty", 16, 1, false, nil, "", 0},
		{"STD of Major axis", 8, 1, false, nil, "", 0},
		{"STD of Minor axis", 8, 1, false, nil, "", 0},
		{"Orientation of Major axis", 8, 1, false, nil, "", 0},
		{"STD of Lat Error", 8, 1, false, nil, "", 0},
		{"STD of Lon Error", 8, 1, false, nil, "", 0},
		{"STD of Alt Error", 8, 1, false, nil, "", 0}},
	},

	{"GNSS RAIM Output", 129545, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Integrity flag", 4, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "Reserved", 0},
		{"Latitude expected error", 8, 1, false, nil, "", 0},
		{"Longitude expected error", 8, 1, false, nil, "", 0},
		{"Altitude expected error", 8, 1, false, nil, "", 0},
		{"SV ID of most likely failed sat", 8, 1, false, nil, "", 0},
		{"Probability of missed detection", 8, 1, false, nil, "", 0},
		{"Estimate of pseudorange bias", 8, 1, false, nil, "", 0},
		{"Std Deviation of bias", 8, 1, false, nil, "", 0}},
	},

	{"GNSS RAIM Settings", 129546, false, 8, 0, []Field{
		{"Radial Position Error Maximum Threshold", 8, 1, false, nil, "", 0},
		{"Probability of False Alarm", 8, 1, false, nil, "", 0},
		{"Probability of Missed Detection", 8, 1, false, nil, "", 0},
		{"Pseudorange Residual Filtering Time Constant", 8, 1, false, nil, "", 0}},
	},

	{"GNSS Pseudorange Error Statistics", 129547, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"RMS Std Dev of Range Inputs", 16, 1, false, nil, "", 0},
		{"Std Dev of Major error ellipse", 8, 1, false, nil, "", 0},
		{"Std Dev of Minor error ellipse", 8, 1, false, nil, "", 0},
		{"Orientation of error ellipse", 8, 1, false, nil, "", 0},
		{"Std Dev Lat Error", 8, 1, false, nil, "", 0},
		{"Std Dev Lon Error", 8, 1, false, nil, "", 0},
		{"Std Dev Alt Error", 8, 1, false, nil, "", 0}},
	},

	{"DGNSS Corrections", 129549, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Reference Station ID", 16, 1, false, nil, "", 0},
		{"Reference Station Type", 16, 1, false, nil, "", 0},
		{"Time of corrections", 8, 1, false, nil, "", 0},
		{"Station Health", 8, 1, false, nil, "", 0},
		{"Reserved Bits", 8, 1, false, nil, "", 0},
		{"Satellite ID", 8, 1, false, nil, "", 0},
		{"PRC", 8, 1, false, nil, "", 0},
		{"RRC", 8, 1, false, nil, "", 0},
		{"UDRE", 8, 1, false, nil, "", 0},
		{"IOD", 8, 1, false, nil, "", 0}},
	},

	{"GNSS Differential Correction Receiver Interface", 129550, false, 8, 0, []Field{
		{"Channel", 8, 1, false, nil, "", 0},
		{"Frequency", 8, 1, false, nil, "", 0},
		{"Serial Interface Bit Rate", 8, 1, false, nil, "", 0},
		{"Serial Interface Detection Mode", 8, 1, false, nil, "", 0},
		{"Differential Source", 8, 1, false, nil, "", 0},
		{"Differential Operation Mode", 8, 1, false, nil, "", 0}},
	},

	{"GNSS Differential Correction Receiver Signal", 129551, false, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Channel", 8, 1, false, nil, "", 0},
		{"Signal Strength", 8, 1, false, nil, "", 0},
		{"Signal SNR", 8, 1, false, nil, "", 0},
		{"Frequency", 8, 1, false, nil, "", 0},
		{"Station Type", 8, 1, false, nil, "", 0},
		{"Station ID", 8, 1, false, nil, "", 0},
		{"Differential Signal Bit Rate", 8, 1, false, nil, "", 0},
		{"Differential Signal Detection Mode", 8, 1, false, nil, "", 0},
		{"Used as Correction Source", 8, 1, false, nil, "", 0},
		{"Reserved", 8, 1, false, nil, "Reserved", 0},
		{"Differential Source", 8, 1, false, nil, "", 0},
		{"Time since Last Sat Differential Sync", 8, 1, false, nil, "", 0},
		{"Satellite Service ID No.", 8, 1, false, nil, "", 0}},
	},

	{"GLONASS Almanac Data", 129556, false, 8, 0, []Field{
		{"PRN", 8, 1, false, nil, "", 0},
		{"NA", 8, 1, false, nil, "", 0},
		{"CnA", 8, 1, false, nil, "", 0},
		{"HnA", 8, 1, false, nil, "", 0},
		{"(epsilon)nA", 8, 1, false, nil, "", 0},
		{"(deltaTnA)DOT", 8, 1, false, nil, "", 0},
		{"(omega)nA", 8, 1, false, nil, "", 0},
		{"(delta)TnA", 8, 1, false, nil, "", 0},
		{"tnA", 8, 1, false, nil, "", 0},
		{"(lambda)nA", 8, 1, false, nil, "", 0},
		{"(delta)inA", 8, 1, false, nil, "", 0},
		{"tcA", 8, 1, false, nil, "", 0},
		{"tnA", 8, 1, false, nil, "", 0}},
	},

	{"AIS DGNSS Broadcast Binary Message", 129792, false, 8, 0, []Field{
		{"Message ID", 8, 1, false, nil, "", 0},
		{"Repeat Indicator", 8, 1, false, nil, "", 0},
		{"Source ID", 8, 1, false, nil, "", 0},
		{"NMEA 2000 Reserved", 8, 1, false, nil, "", 0},
		{"AIS Tranceiver Information", 8, 1, false, nil, "", 0},
		{"Spare", 8, 1, false, nil, "", 0},
		{"Longitude", 8, 1, false, nil, "", 0},
		{"Latitude", 8, 1, false, nil, "", 0},
		{"NMEA 2000 Reserved", 8, 1, false, nil, "", 0},
		{"Spare", 8, 1, false, nil, "", 0},
		{"Number of Bits in Binary Data Field", 8, 1, false, nil, "", 0},
		{"Binary Data", 64, RES_BINARY, false, nil, "", 0}},
	},

	{"AIS UTC and Date Report", 129793, false, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Position Accuracy", 1, RES_LOOKUP, false, lookupAisAccuracy, "", 0},
		{"RAIM", 1, RES_LOOKUP, false, lookupAisRAIM, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "NMEA reserved to align next data on byte boundary", 0},
		{"Position Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Communication State", 19, RES_BINARY, false, nil, "Information used by the TDMA slot allocation algorithm and synchronization information", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Position Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "NMEA reserved to align next data on byte boundary", 0},
		{"GNSS type", 4, RES_LOOKUP, false, lookupGnsAis, "", 0},
		{"Spare", 8, RES_BINARY, false, nil, "", 0}},
	},

	// http://www.navcen.uscg.gov/enav/ais/AIS_messages.htm
	{"AIS Class A Static and Voyage Related Data", 129794, true, 0x18, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"IMO number", 32, RES_INTEGER, false, nil, "0=unavailable", 0},
		{"Callsign", 56, RES_ASCII, false, nil, "0=unavailable", 0},
		{"Name", 160, RES_ASCII, false, nil, "0=unavailable", 0},
		{"Type of ship", 8, RES_LOOKUP, false, lookupShipType, "", 0},
		{"Length", 16, 0.1, false, "m", "", 0},
		{"Beam", 16, 0.1, false, "m", "", 0},
		{"Position reference from Starboard", 16, 0.1, false, "m", "", 0},
		{"Position reference from Bow", 16, 0.1, false, "m", "", 0},
		{"ETA Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"ETA Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Draft", 16, 0.01, false, "m", "", 0},
		{"Destination", 160, RES_ASCII, false, nil, "0=unavailable", 0},
		{"AIS version indicator", 2, RES_LOOKUP, false, lookupAisVersion, "", 0},
		{"GNSS type", 4, RES_LOOKUP, false, lookupGnsAis, "", 0},
		{"DTE", 1, RES_LOOKUP, false, lookupAisDTE, "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0}},
	},

	{"AIS Addressed Binary Message", 129795, true, 13, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Sequence Number", 2, 1, false, nil, "", 0},
		{"Destination ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "reserved", 0},
		{"Retransmit flag", 1, 1, false, nil, "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"Number of Bits in Binary Data Field", 16, RES_INTEGER, false, nil, "", 0},
		{"Binary Data", 64, RES_BINARY, false, nil, "", 0}},
	},

	{"AIS Acknowledge", 129796, true, 12, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 32, 1, false, "MMSI", "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Destination ID #1", 32, 1, false, nil, "", 0},
		{"Sequence Number for ID 1", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "reserved", 0},
		{"Sequence Number for ID n", 2, RES_BINARY, false, nil, "reserved", 0}},
	},

	{"AIS Binary Broadcast Message", 129797, true, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 32, 1, false, nil, "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Number of Bits in Binary Data Field", 16, 1, false, nil, "", 0},
		{"Binary Data", 2040, RES_BINARY, false, nil, "", 0}},
	},

	{"AIS SAR Aircraft Position Report", 129798, false, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Position Accuracy", 1, RES_LOOKUP, false, lookupPositionAccuracy, "", 0},
		{"RAIM", 1, RES_LOOKUP, false, lookupAisRAIM, "", 0},
		{"Time Stamp", 6, RES_LOOKUP, false, lookupTimeStamp, "0-59 = UTC second when the report was generated", 0},
		{"COG", 16, RES_DEGREES, false, "deg", "", 0},
		{"SOG", 16, 0.1, false, "m/s", "", 0},
		{"Communication State", 19, RES_BINARY, false, nil, "Information used by the TDMA slot allocation algorithm and synchronization information", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Altitude", 64, 1e-6, true, "m", "", 0},
		{"Reserved for Regional Applications", 8, 1, false, nil, "", 0},
		{"DTE", 1, RES_LOOKUP, false, lookupAisDTE, "", 0},
		{"Reserved", 7, RES_BINARY, false, nil, "reserved", 0}},
	},

	{"Radio Frequency/Mode/Power", 129799, false, 9, 0, []Field{
		{"Rx Frequency", 32, 10, false, "Hz", "", 0},
		{"Tx Frequency", 32, 10, false, "Hz", "", 0},
		{"Radio Channel", 8, 1, false, nil, "", 0},
		{"Tx Power", 8, 1, false, nil, "", 0},
		{"Mode", 8, 1, false, nil, "", 0},
		{"Channel Bandwidth", 8, 1, false, nil, "", 0}},
	},

	{"AIS UTC/Date Inquiry", 129800, false, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 30, 1, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 3, RES_BINARY, false, nil, "reserved", 0},
		{"Destination ID", 30, 1, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0}},
	},

	{"AIS Addressed Safety Related Message", 129801, true, 12, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 32, 1, false, "MMSI", "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Sequence Number", 2, 1, false, nil, "", 0},
		{"Destination ID", 32, 1, false, "MMSI", "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "reserved", 0},
		{"Retransmit flag", 1, 1, false, nil, "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"Safety Related Text", 2040, RES_ASCII, false, nil, "", 0}},
	},

	{"AIS Safety Related Broadcast Message", 129802, false, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 30, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 3, RES_BINARY, false, nil, "reserved", 0},
		{"Safety Related Text", 288, RES_ASCII, false, nil, "", 0}},
	},

	{"AIS Interrogation", 129803, false, 8, 8, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 30, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 3, RES_BINARY, false, nil, "reserved", 0},
		{"Destination ID", 30, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Message ID A", 8, RES_INTEGER, false, nil, "", 0},
		{"Slot Offset A", 14, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Message ID B", 8, RES_INTEGER, false, nil, "", 0},
		{"Slot Offset B", 14, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0}},
	},

	{"AIS Assignment Mode Command", 129804, true, 23, 3, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Reserved", 1, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Destination ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Offset", 16, RES_INTEGER, false, nil, "", 0},
		{"Increment", 16, RES_INTEGER, false, nil, "", 0}},
	},

	{"AIS Data Link Management Message", 129805, false, 8, 4, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 30, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 3, RES_BINARY, false, nil, "reserved", 0},
		{"Offset", 10, RES_INTEGER, false, nil, "", 0},
		{"Number of Slots", 8, RES_INTEGER, false, nil, "", 0},
		{"Timeout", 8, RES_INTEGER, false, nil, "", 0},
		{"Increment", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"AIS Channel Management", 129806, false, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 30, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"AIS Transceiver information", 5, RES_LOOKUP, false, lookupAisTransceiver, "", 0},
		{"Reserved", 3, RES_BINARY, false, nil, "reserved", 0},
		{"Channel A", 7, 1, false, nil, "", 0},
		{"Channel B", 7, 1, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Power", 8, 1, false, nil, "reserved", 0},
		{"Tx/Rx Mode", 8, RES_INTEGER, false, nil, "", 0},
		{"North East Longitude Corner 1", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"North East Latitude Corner 1", 32, RES_LATITUDE, true, "deg", "", 0},
		{"South West Longitude Corner 1", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"South West Latitude Corner 2", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "reserved", 0},
		{"Addressed or Broadcast Message Indicator", 2, 1, false, nil, "", 0},
		{"Channel A Bandwidth", 7, RES_INTEGER, false, nil, "", 0},
		{"Channel B Bandwidth", 7, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Transitional Zone Size", 8, 1, false, nil, "", 0}},
	},

	{"AIS Class B Group Assignment", 129807, false, 8, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat Indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"Source ID", 30, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Tx/Rx Mode", 2, RES_INTEGER, false, nil, "", 0},
		{"Reserved", 6, RES_BINARY, false, nil, "reserved", 0},
		{"North East Longitude Corner 1", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"North East Latitude Corner 1", 32, RES_LATITUDE, true, "deg", "", 0},
		{"South West Longitude Corner 1", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"South West Latitude Corner 2", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Station Type", 8, 1, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Ship and Cargo Filter", 6, 1, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Reporting Interval", 16, 1, false, nil, "", 0},
		{"Quiet Time", 16, 1, false, nil, "", 0}},
	},

	{"DSC Call Information", 129808, false, 8, 2, []Field{
		{"DSC Format Symbol", 8, 1, false, nil, "", 0},
		{"DSC Category Symbol", 8, 1, false, nil, "", 0},
		{"DSC Message Address", 8, 1, false, nil, "", 0},
		{"Nature of Distress or 1st Telecommand", 8, 1, false, nil, "", 0},
		{"Subsequent Communication Mode or 2nd Telecommand", 8, 1, false, nil, "", 0},
		{"Proposed Rx Frequency/Channel", 8, 1, false, nil, "", 0},
		{"Proposed Tx Frequency/Channel", 8, 1, false, nil, "", 0},
		{"Telephone Number", 8, 1, false, nil, "", 0},
		{"Latitude of Vessel Reported", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Longitude of Vessel Reported", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Time of Position", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"User ID of Ship In Distress", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"DSC EOS Symbol", 8, 1, false, nil, "", 0},
		{"Expansion Enabled", 8, 1, false, nil, "", 0},
		{"Calling Rx Frequency/Channel", 8, 1, false, nil, "", 0},
		{"Calling Tx Frequency/Channel", 8, 1, false, nil, "", 0},
		{"Time of Receipt", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Date of Receipt", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"DSC Equipment Assigned Message ID", 8, 1, false, nil, "", 0},
		{"DSC Expansion Field Symbol", 8, 1, false, nil, "", 0},
		{"DSC Expansion Field Data", 8, 1, false, nil, "", 0}},
	},

	{"AIS Class B static data (msg 24 Part A)", 129809, false, 20 + 4 + 1, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Name", 160, RES_ASCII, false, nil, "", 0}},
	},

	{"AIS Class B static data (msg 24 Part B)", 129810, false, 0x25 - 4, 0, []Field{
		{"Message ID", 6, 1, false, nil, "", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Type of ship", 8, RES_LOOKUP, false, lookupShipType, "", 0},
		{"Vendor ID", 56, RES_ASCII, false, nil, "", 0},
		{"Callsign", 56, RES_ASCII, false, nil, "0=unavailable", 0},
		{"Length", 16, 0.1, false, "m", "", 0},
		{"Beam", 16, 0.1, false, "m", "", 0},
		{"Position reference from Starboard", 16, 0.1, false, "m", "", 0},
		{"Position reference from Bow", 16, 0.1, false, "m", "", 0},
		{"Mothership User ID", 32, RES_INTEGER, false, "MMSI", "MMSI of mother ship sent by daughter vessels", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Spare", 6, RES_INTEGER, false, nil, "0=unavailable", 0}},
	},

	{"Route and WP Service - Database List", 130064, false, 8, 9, []Field{
		{"Start Database ID", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of Databases Available", 8, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Database Name", 64, RES_ASCII, false, nil, "", 0},
		{"Database Timestamp", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Database Datestamp", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"WP Position Resolution", 6, 1, false, nil, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "reserved", 0},
		{"Number of Routes in Database", 16, 1, false, nil, "", 0},
		{"Number of WPs in Database", 16, 1, false, nil, "", 0},
		{"Number of Bytes in Database", 16, 1, false, nil, "", 0}},
	},

	{"Route and WP Service - Route List", 130065, false, 8, 6, []Field{
		{"Start Route ID", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of Routes in Database", 8, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"Route Name", 64, RES_ASCII, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "reserved", 0},
		{"WP Identification Method", 2, 1, false, nil, "", 0},
		{"Route Status", 2, 1, false, nil, "", 0}},
	},

	{"Route and WP Service - Route/WP-List Attributes", 130066, false, 8, 0, []Field{
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"Route/WP-List Name", 64, RES_ASCII, false, nil, "", 0},
		{"Route/WP-List Timestamp", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Route/WP-List Datestamp", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Change at Last Timestamp", 8, 1, false, nil, "", 0},
		{"Number of WPs in the Route/WP-List", 16, 1, false, nil, "", 0},
		{"Critical supplementary parameters", 8, 1, false, nil, "", 0},
		{"Navigation Method", 2, 1, false, nil, "", 0},
		{"WP Identification Method", 2, 1, false, nil, "", 0},
		{"Route Status", 2, 1, false, nil, "", 0},
		{"XTE Limit for the Route", 16, 1, false, nil, "", 0}},
	},

	{"Route and WP Service - Route - WP Name & Position", 130067, false, 8, 4, []Field{
		{"Start RPS#", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of WPs in the Route/WP-List", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"WP ID", 8, 1, false, nil, "", 0},
		{"WP Name", 64, RES_ASCII, false, nil, "", 0},
		{"WP Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"WP Longitude", 32, RES_LONGITUDE, true, "deg", "", 0}},
	},

	{"Route and WP Service - Route - WP Name", 130068, false, 8, 2, []Field{
		{"Start RPS#", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of WPs in the Route/WP-List", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"WP ID", 8, 1, false, nil, "", 0},
		{"WP Name", 64, RES_ASCII, false, nil, "", 0}},
	},

	{"Route and WP Service - XTE Limit & Navigation Method", 130069, false, 8, 6, []Field{
		{"Start RPS#", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of WPs with a specific XTE Limit or Nav. Method", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"RPS#", 8, 1, false, nil, "", 0},
		{"XTE limit in the leg after WP", 16, 1, false, nil, "", 0},
		{"Nav. Method in the leg after WP", 4, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0}},
	},

	{"Route and WP Service - WP Comment", 130070, false, 8, 2, []Field{
		{"Start ID", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of WPs with Comments", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"WP ID / RPS#", 8, 1, false, nil, "", 0},
		{"Comment", 64, RES_ASCII, false, nil, "", 0}},
	},

	{"Route and WP Service - Route Comment", 130071, false, 8, 2, []Field{
		{"Start Route ID", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of Routes with Comments", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"Comment", 64, RES_ASCII, false, nil, "", 0}},
	},

	{"Route and WP Service - Database Comment", 130072, false, 8, 2, []Field{
		{"Start Database ID", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of Databases with Comments", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Comment", 64, RES_ASCII, false, nil, "", 0}},
	},

	{"Route and WP Service - Radius of Turn", 130073, false, 8, 2, []Field{
		{"Start RPS#", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of WPs with a specific Radius of Turn", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Route ID", 8, 1, false, nil, "", 0},
		{"RPS#", 8, 1, false, nil, "", 0},
		{"Radius of Turn", 16, 1, false, nil, "", 0}},
	},

	{"Route and WP Service - WP List - WP Name & Position", 130074, false, 8, 4, []Field{
		{"Start WP ID", 8, 1, false, nil, "", 0},
		{"nItems", 8, 1, false, nil, "", 0},
		{"Number of valid WPs in the WP-List", 16, 1, false, nil, "", 0},
		{"Database ID", 8, 1, false, nil, "", 0},
		{"Reserved", 8, 1, false, nil, "reserved", 0},
		{"WP ID", 8, 1, false, nil, "", 0},
		{"WP Name", 64, RES_ASCII, false, nil, "", 0},
		{"WP Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"WP Longitude", 32, RES_LONGITUDE, true, "deg", "", 0}},
	},

	// http://askjackrabbit.typepad.com/ask_jack_rabbit/page/7
	{"Wind Data", 130306, true, 6, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Wind Speed", 16, 0.01, false, "m/s", "", 0},
		{"Wind Angle", 16, RES_DEGREES, false, "deg", "", 0},
		{"Reference", 3, RES_LOOKUP, false, lookupWindReference, "", 0}},
	},

	// Water temperature, Transducer Measurement
	{"Environmental Parameters", 130310, true, 7, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Water Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Outside Ambient Air Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Atmospheric Pressure", 16, RES_PRESSURE, false, "hPa", "", 0}},
	},

	{"Environmental Parameters", 130311, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Temperature Instance", 6, RES_LOOKUP, false, lookupTemperatureSource, "", 0},
		{"Humidity Instance", 2, RES_LOOKUP, false, lookupHumidityInstance, "", 0},
		{"Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Humidity", 16, 0.004, true, "%", "", 0},
		{"Atmospheric Pressure", 16, RES_PRESSURE, false, "hPa", "", 0}},
	},

	{"Temperature", 130312, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Temperature Instance", 8, 1, false, nil, "", 0},
		{"Temperature Source", 8, RES_LOOKUP, false, lookupTemperatureSource, "", 0},
		{"Actual Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Set Temperature", 16, RES_TEMPERATURE, false, "K", "", 0}},
	},

	{"Humidity", 130313, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Humidity Instance", 8, 1, false, nil, "", 0},
		{"Humidity Source", 8, 1, false, nil, "", 0},
		{"Actual Humidity", 16, 100.0 / 25000, true, "%", "", 0},
		{"Set Humidity", 16, 100.0 / 25000, true, "%", "", 0}},
	},

	// Based off the definition for 130315. Appears to be correct
	{"Actual Pressure", 130314, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Pressure Instance", 8, 1, false, nil, "", 0},
		{"Pressure Source", 8, RES_LOOKUP, false, lookupPressureSource, "", 0},
		{"Pressure", 32, 0.1, true, "Pa", "", 0}},
	},

	// Source: http://standards.nmea.org/NSNA/corrigenda/nmea-2000/nmea-2000-corrigendum-3-2009.pdf
	{"Set Pressure", 130315, true, 8, 0, []Field{
		{"SID", 8, 1, false, nil, "", 0},
		{"Pressure Instance", 8, 1, false, nil, "", 0},
		{"Pressure Source", 8, RES_LOOKUP, false, lookupPressureSource, "", 0},
		{"Pressure", 32, 0.1, true, "Pa", "", 0}},
	},

	{"Tide Station Data", 130320, true, 20, 0, []Field{
		{"Mode", 4, RES_LOOKUP, false, lookupResidualMode, "", 0},
		{"Tide Tendency", 2, RES_LOOKUP, false, lookupTideTendency, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "", 0},
		{"Measurement Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Measurement Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Station Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Station Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Tide Level", 16, 0.001, true, "m", "Relative to MLLW", 0},
		{"Tide Level standard deviation", 16, 0.01, false, "m", "", 0},
		{"Station ID", 16, RES_STRING, false, nil, "", 0},
		{"Station Name", 16, RES_STRING, false, nil, "", 0}},
	},

	{"Salinity Station Data", 130321, true, 22, 0, []Field{
		{"Mode", 4, RES_LOOKUP, false, lookupResidualMode, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0},
		{"Measurement Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Measurement Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Station Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Station Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Salinity", 32, RES_FLOAT, true, "ppt", "The average Salinity of ocean water is about 35 grams of salts per kilogram of sea water (g/kg), usually written as 35 ppt which is read as 35 parts per thousand.", 0},
		{"Water Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Station ID", 16, RES_STRING, false, nil, "", 0},
		{"Station Name", 16, RES_STRING, false, nil, "", 0}},
	},

	{"Current Station Data", 130322, false, 8, 0, []Field{
		{"Mode", 4, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0},
		{"Measurement Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Measurement Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Station Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Station Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Measurement Depth", 32, 0.01, false, "m", "Depth below transducer", 0},
		{"Current speed", 16, 0.01, false, "m/s", "", 0},
		{"Current flow direction", 16, RES_DEGREES, false, "deg", "", 0},
		{"Water Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Station ID", 16, RES_STRING, false, nil, "", 0},
		{"Station Name", 16, RES_STRING, false, nil, "", 0}},
	},

	{"Meteorological Station Data", 130323, false, 0x1e, 0, []Field{
		{"Mode", 4, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0},
		{"Measurement Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Measurement Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Station Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Station Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Wind Speed", 16, 0.01, false, "m/s", "", 0},
		{"Wind Direction", 16, RES_DEGREES, false, "deg", "", 0},
		{"Wind Reference", 3, RES_LOOKUP, false, lookupWindReference, "", 0},
		{"Reserved", 5, RES_BINARY, false, "", "reserved", 0},
		{"Wind Gusts", 16, 0.01, false, "m/s", "", 0},
		{"Atmospheric Pressure", 16, RES_PRESSURE, false, "hPa", "", 0},
		{"Ambient Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Station ID", 16, RES_STRING, false, nil, "", 0},
		{"Station Name", 16, RES_STRING, false, nil, "", 0}},
	},

	{"Moored Buoy Station Data", 130324, false, 8, 0, []Field{
		{"Mode", 4, 1, false, nil, "", 0},
		{"Reserved", 4, RES_BINARY, false, nil, "", 0},
		{"Measurement Date", 16, RES_DATE, false, "days", "Days since January 1, 1970", 0},
		{"Measurement Time", 32, RES_TIME, false, "s", "Seconds since midnight", 0},
		{"Station Latitude", 32, RES_LATITUDE, true, "deg", "", 0},
		{"Station Longitude", 32, RES_LONGITUDE, true, "deg", "", 0},
		{"Wind Speed", 16, 0.01, false, "m/s", "", 0},
		{"Wind Direction", 16, RES_DEGREES, false, "deg", "", 0},
		{"Wind Reference", 3, RES_LOOKUP, false, lookupWindReference, "", 0},
		{"Reserved", 5, RES_BINARY, false, "", "reserved", 0},
		{"Wind Gusts", 16, 0.01, false, "m/s", "", 0},
		{"Wave Height", 16, 1, false, nil, "", 0},
		{"Dominant Wave Period", 16, 1, false, nil, "", 0},
		{"Atmospheric Pressure", 16, RES_PRESSURE, false, "hPa", "", 0},
		{"Pressure Tendency Rate", 16, 1, false, "", "", 0},
		{"Air Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Water Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Station ID", 64, RES_ASCII, false, nil, "", 0}},
	},

	{"Small Craft Status", 130576, true, 2, 0, []Field{
		{"Port trim tab", 8, 1, true, nil, "", 0},
		{"Starboard trim tab", 8, 1, true, nil, "", 0}},
	},

	{"Direction Data", 130577, true, 14, 0, []Field{
		{"Data Mode", 4, RES_LOOKUP, false, lookupResidualMode, "", 0},
		{"COG Reference", 2, RES_LOOKUP, false, lookupDirectionReference, "", 0},
		{"Reserved", 2, RES_BINARY, false, nil, "Reserved", 0},
		{"SID", 8, 1, false, nil, "", 0},
		// So far, 2 bytes. Very sure of this given molly rose data
		{"COG", 16, RES_DEGREES, false, "deg", "", 0},
		{"SOG", 16, 0.01, false, "m/s", "", 0},
		{"Heading", 16, RES_DEGREES, false, "deg", "", 0},
		{"Speed through Water", 16, 0.01, false, "m/s", "", 0},
		{"Set", 16, RES_DEGREES, false, "deg", "", 0},
		{"Drift", 16, 0.01, false, "m/s", "", 0}},
	},

	{"Vessel Speed Components", 130578, true, 12, 0, []Field{
		{"Longitudinal Speed, Water-referenced", 16, 0.001, true, "m/s", "", 0},
		{"Transverse Speed, Water-referenced", 16, 0.001, true, "m/s", "", 0},
		{"Longitudinal Speed, Ground-referenced", 16, 0.001, true, "m/s", "", 0},
		{"Transverse Speed, Ground-referenced", 16, 0.001, true, "m/s", "", 0},
		{"Stern Speed, Water-referenced", 16, 0.001, true, "m/s", "", 0},
		{"Stern Speed, Ground-referenced", 16, 0.001, true, "m/s", "", 0}},
	},

	{"SonicHub: Init #2", 130816, false, 9, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=1", "Init #2", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"A", 16, RES_INTEGER, false, nil, "", 0},
		{"B", 16, RES_INTEGER, false, nil, "", 0}},
	},

	{"SonicHub: AM Radio", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, lookupCompanyCode, "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, "=275", "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=4", "AM Radio", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 8, RES_LOOKUP, false, lookupSonicHubTuning, "", 0},
		{"Frequency", 32, 0.001, false, "kHz", "", 0},
		{"Noise level", 2, 1, false, nil, "", 0},  // Not sure about this
		{"Signal level", 4, 1, false, nil, "", 0}, // ... and this, doesn't make complete sense compared to display
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Zone Info", 130816, false, 6, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=5", "Zone info", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Zone", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"SonicHub: Source", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=6", "Source", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Source", 8, RES_LOOKUP, false, lookupSonicHubSource, "", 0}},
	},

	{"SonicHub: Source List", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=8", "Source list", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Source ID", 8, RES_INTEGER, false, nil, "", 0},
		{"A", 8, RES_INTEGER, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Mute Control", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=9", "Control", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 8, RES_LOOKUP, false, lookupSonicHubMute, "", 0}},
	},

	{"SonicHub: FM Radio", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=12", "FM Radio", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 8, RES_LOOKUP, false, lookupSonicHubTuning, "", 0},
		{"Frequency", 32, 0.001, false, "kHz", "", 0},
		{"Noise level", 2, 1, false, nil, "", 0},  // Not sure about this
		{"Signal level", 4, 1, false, nil, "", 0}, // ... and this, doesn't make complete sense compared to display
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Playlist", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=13", "Playlist", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 8, RES_LOOKUP, false, lookupSonicHubPlaylist, "", 0},
		{"A", 8, RES_INTEGER, false, nil, "", 0},
		{"Current Track", 32, RES_INTEGER, false, nil, "", 0},
		{"Tracks", 32, RES_INTEGER, false, nil, "", 0},
		{"Length", 32, 0.001, false, nil, "Seconds", 0},
		{"Position in track", 32, 0.001, false, nil, "Seconds", 0}},
	},

	{"SonicHub: Track", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=14", "Track", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 32, RES_INTEGER, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Artist", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=15", "Artist", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 32, RES_INTEGER, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Album", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=16", "Album", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 32, RES_INTEGER, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Menu Item", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=19", "Menu Item", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Item", 32, RES_INTEGER, false, nil, "", 0},
		{"C", 8, 1, false, nil, "", 0},
		{"D", 8, 1, false, nil, "", 0},
		{"E", 8, 1, false, nil, "", 0},
		{"Text", 256, RES_STRINGLZ, false, nil, "", 0}},
	},

	{"SonicHub: Zones", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=20", "Zones", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Zones", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"SonicHub: Max Volume", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=23", "Max Volume", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Zone", 8, RES_LOOKUP, false, lookupSonicHubZone, "", 0},
		{"Level", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"SonicHub: Volume", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=24", "Volume", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Zone", 8, RES_LOOKUP, false, lookupSonicHubZone, "", 0},
		{"Level", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"SonicHub: Init #1", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=25", "Init #1", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0}},
	},

	{"SonicHub: Position", 130816, true, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=48", "Position", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"Position", 32, 0.001, false, nil, "Seconds", 0}},
	},

	{"SonicHub: Init #3", 130816, false, 9, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=50", "Init #3", 0},
		{"Control", 8, RES_LOOKUP, false, lookupSonicHubControl, "", 0},
		{"A", 8, RES_INTEGER, false, nil, "", 0},
		{"B", 8, RES_INTEGER, false, nil, "", 0}},
	},

	{"Simrad: Text Message", 130816, false, 0x40, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "", 0},
		{"Reserved", 2, 1, false, nil, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, nil, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=50", "Init #3", 0},
		{"A", 8, 1, false, nil, "", 0},
		{"B", 8, 1, false, nil, "", 0},
		{"C", 8, 1, false, nil, "", 0},
		{"SID", 8, 1, false, nil, "", 0},
		{"Prio", 8, 1, false, nil, "", 0},
		{"Text", 256, RES_ASCII, false, nil, "", 0}},
	},

	{"Navico: Product Information", 130817, false, 0x0e, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=275", "Navico", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Product Code", 16, RES_INTEGER, false, 0, "", 0},
		{"Model", 256, RES_ASCII, false, 0, "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"Firmware version", 80, RES_ASCII, false, 0, "", 0},
		{"Firmware date", 256, RES_ASCII, false, 0, "", 0},
		{"Firmware time", 256, RES_ASCII, false, 0, "", 0}},
	},

	{"Simnet: Reprogram Data", 130818, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Version", 16, RES_INTEGER, false, 0, "", 0},
		{"Sequence", 16, RES_INTEGER, false, 0, "", 0},
		{"Data", 2040, RES_BINARY, false, 0, "", 0}},
	},

	{"Simnet: Request Reprogram", 130819, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	/* Fusion */
	{"Fusion: Source Name", 130820, false, 13, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=2", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0},
		{"Source", 40, RES_STRINGLZ, false, 0, "", 0}},
	},

	{"Fusion: Track", 130820, false, 0x20, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=5", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 40, 1, false, 0, "", 0},
		{"Track", 80, RES_STRINGLZ, false, 0, "", 0}},
	},

	{"Fusion: Artist", 130820, false, 0x20, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=6", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 40, 1, false, 0, "", 0},
		{"Artist", 80, RES_STRINGLZ, false, 0, "", 0}},
	},

	{"Fusion: Album", 130820, false, 0x20, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=7", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 40, 1, false, 0, "", 0},
		{"Album", 80, RES_STRINGLZ, false, 0, "", 0}},
	},

	{"Fusion: Play Progress", 130820, false, 9, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=9", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Progress", 24, 0.001, false, "s", "", 0}},
	},

	{"Fusion: AM/FM Station", 130820, false, 0x0A, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=11", "", 0},
		{"A", 24, 1, false, 0, "", 0},
		{"Frequency", 32, 1, false, "Hz", "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Track", 80, RES_STRINGLZ, false, 0, "", 0}},
	},

	{"Fusion: VHF", 130820, false, 9, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=12", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Channel", 8, 1, false, 0, "", 0},
		{"D", 24, 1, false, 0, "", 0}},
	},

	{"Fusion: Squelch", 130820, false, 6, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=13", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Squelch", 8, 1, false, 0, "", 0}},
	},

	{"Fusion: Scan", 130820, false, 6, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=14", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Scan", 8, RES_LOOKUP, false, ",0=Off,1=Scan", "", 0}},
	},

	{"Fusion: Menu Item", 130820, false, 23, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=17", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Line", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0},
		{"F", 8, 1, false, 0, "", 0},
		{"G", 8, 1, false, 0, "", 0},
		{"H", 8, 1, false, 0, "", 0},
		{"I", 8, 1, false, 0, "", 0},
		{"Text", 40, RES_STRINGLZ, false, 0, "", 0}},
	},

	{"Fusion: Replay", 130820, false, 23, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=19", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"Mode", 8, RES_LOOKUP, false, ",9=USB Repeat,10=USB Shuffle,12=iPod Repeat,13=iPod Shuffle", "", 0},
		{"C", 24, 1, false, 0, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0},
		{"Status", 8, RES_LOOKUP, false, ",0=Off,1=One/Track,2=All/Album", "", 0},
		{"H", 8, 1, false, 0, "", 0},
		{"I", 8, 1, false, 0, "", 0},
		{"J", 8, 1, false, 0, "", 0}},
	},

	{"Fusion: Time", 130820, false, 23, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=20", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"Command ID", 8, 1, false, "=59", "", 0},
		{"C", 24, 1, false, 0, "", 0},
		{"Minutes", 8, 1, false, 0, "", 0},
		{"Hours", 8, 1, false, 0, "", 0},
		{"Format", 1, RES_LOOKUP, false, ",0=12h,1=24h", "", 0},
		{"Alarm", 1, RES_LOOKUP, false, ",0=Off,1=On", "", 0},
		{"I", 6, 1, false, 0, "", 0},
		{"J", 8, 1, false, 0, "", 0}},
	},

	{"Fusion: Mute", 130820, false, 5, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=23", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"Mute", 8, RES_LOOKUP, false, ",1=Muted,2=Not Muted", "", 0}},
	},

	// Range: 0 to +24
	{"Fusion: Sub Volume", 130820, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=26", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"Zone 1", 8, 1, false, "vol", "", 0},
		{"Zone 2", 8, 1, false, "vol", "", 0},
		{"Zone 3", 8, 1, false, "vol", "", 0},
		{"Zone 4", 8, 1, false, "vol", "", 0}},
	},

	// Range: -15 to +15
	{"Fusion: Tone", 130820, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=27", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"Bass", 8, 1, true, "vol", "", 0},
		{"Mid", 8, 1, true, "vol", "", 0},
		{"Treble", 8, 1, true, "vol", "", 0}},
	},

	{"Fusion: Volume", 130820, false, 0x0A, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=29", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"Zone 1", 8, 1, false, "vol", "", 0},
		{"Zone 2", 8, 1, false, "vol", "", 0},
		{"Zone 3", 8, 1, false, "vol", "", 0},
		{"Zone 4", 8, 1, false, "vol", "", 0}},
	},

	{"Fusion: Transport", 130820, false, 5, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=419", "Fusion", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 8, 1, false, "=32", "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"Transport", 8, RES_LOOKUP, false, ",1=Paused", "", 0}},
	},

	// M/V Dirona
	{"Furuno: Unknown", 130820, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1855", "Furuno", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0}},
	},

	{"Simnet: Reprogram Status", 130820, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Reserved", 8, 1, false, 0, "", 0},
		{"Status", 8, 1, false, 0, "", 0},
		{"Reserved", 24, 1, false, 0, "", 0}},
	},

	// M/V Dirona
	{"Furuno: Unknown", 130821, false, 0x0c, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1855", "Furuno", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"SID", 8, 1, false, 0, "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0},
		{"F", 8, 1, false, 0, "", 0},
		{"G", 8, 1, false, 0, "", 0},
		{"H", 8, 1, false, 0, "", 0},
		{"I", 8, 1, false, 0, "", 0}},
	},

	// Uwe Lovas has seen this from EP-70R
	{"Lowrance: unknown", 130827, false, 10, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=140", "Lowrance", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"A", 8, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 16, 1, false, 0, "", 0},
		{"F", 16, 1, false, 0, "", 0}},
	},

	{"Simnet: Set Serial Number", 130828, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Suzuki: Engine and Storage Device Config", 130831, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, 0, "", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Fuel Used - High Resolution", 130832, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Engine and Tank Configuration", 130834, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Set Engine and Tank Configuration", 130835, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	// Seen when HDS8 configures EP65R
	{"Simnet: Fluid Level Sensor Configuration", 130836, false, 0x0e, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, 0, "", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"Device", 8, RES_INTEGER, false, 0, "", 0},
		{"Instance", 8, 1, false, 0, "", 0},
		{"F", 1 * 4, 1, false, 0, "", 0},
		{"Tank type", 1 * 4, RES_LOOKUP, false, lookupTankType, "", 0},
		{"Capacity", 32, 0.1, false, 0, "", 0},
		{"G", 8, 1, false, 0, "", 0},
		{"H", 16, 1, true, 0, "", 0},
		{"I", 8, 1, true, 0, "", 0}},
	},

	{"Simnet: Fuel Flow Turbine Configuration", 130837, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Fluid Level Warning", 130838, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Pressure Sensor Configuration", 130839, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Data User Group Configuration", 130840, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	// Where did this come from ?
	// 130842: { "Simnet: DSC Message", 130842, false, 0x08, 0,[]Field{
	//  { "Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad",0},
	//  { "Reserved", 2, 1, false, 0, "", 0},
	//  { "Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	// },

	{"Simnet: AIS Class B static data (msg 24 Part A)", 130842, false, 0x1d, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, lookupCompanyCode, "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 6, 1, false, "=0", "Msg 24 Part A", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Name", 160, RES_ASCII, false, 0, "", 0}},
	},

	{"Simnet: AIS Class B static data (msg 24 Part B)", 130842, false, 0x25, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, lookupCompanyCode, "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 6, 1, false, "=1", "Msg 24 Part B", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"E", 8, 1, false, 0, "", 0},
		{"User ID", 32, RES_INTEGER, false, "MMSI", "", 0},
		{"Type of ship", 8, RES_LOOKUP, false, lookupShipType, "", 0},
		{"Vendor ID", 56, RES_ASCII, false, 0, "", 0},
		{"Callsign", 56, RES_ASCII, false, 0, "0=unavailable", 0},
		{"Length", 16, 0.1, false, "m", "", 0},
		{"Beam", 16, 0.1, false, "m", "", 0},
		{"Position reference from Starboard", 16, 0.1, false, "m", "", 0},
		{"Position reference from Bow", 16, 0.1, false, "m", "", 0},
		{"Mothership User ID", 32, RES_INTEGER, false, "MMSI", "sent by daughter vessels", 0},
		{"", 2, 1, false, 0, "", 0},
		{"Spare", 6, RES_INTEGER, false, 0, "0=unavailable", 0}},
	},

	{"Simnet: Sonar Status, Frequency and DSP Voltage", 130843, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0}},
	},

	{"Simnet: Parameter Handle", 130845, false, 0x0e, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 6, 1, false, 0, "", 0},
		{"Repeat indicator", 2, RES_LOOKUP, false, lookupRepeatIndicator, "", 0},
		{"D", 8, 1, false, 0, "", 0},
		{"Group", 8, 1, false, 0, "", 0},
		{"F", 8, 1, false, 0, "", 0},
		{"G", 8, 1, false, 0, "", 0},
		{"H", 8, 1, false, 0, "", 0},
		{"I", 8, 1, false, 0, "", 0},
		{"J", 8, 1, false, 0, "", 0},
		{"Backlight", 8, RES_LOOKUP, false, lookupSimnetBacklightLevel, "", 0},
		{"L", 16, 1, false, 0, "", 0}},
	},

	{"Simnet: Event Command: Alarm?", 130850, false, 12, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"A", 16, 1, false, 0, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=1", "Alarm command", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"Alarm", 16, RES_LOOKUP, false, ",57=Raise,56=Clear", "", 0},
		{"Message ID", 16, RES_INTEGER, false, 0, "", 0},
		{"F", 8, 1, false, 0, "", 0},
		{"G", 8, 1, false, 0, "", 0}},
	},

	{"Simnet: Event Command: AP command", 130850, false, 12, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=2", "AP command", 0},
		{"B", 16, 1, false, 0, "", 0},
		{"Controlling Device", 8, 1, false, 0, "", 0},
		{"Event", 16, RES_LOOKUP, false, lookupSimnetApEvents, "", 0},
		{"Direction", 8, RES_LOOKUP, false, lookupSimnetDirection, "", 0},
		{"Angle", 16, RES_DEGREES, false, "deg", "", 0},
		{"G", 8, 1, false, 0, "", 0}},
	},

	//  1: { "Simnet: Event Command: Unknown", 130850, false, 12, 0, []Field{
	//    { "Manufacturer Code", 11, RES_MANUFACTURER, false, lookupCompanyCode, "Simrad", 0},
	//    { "Reserved", 2, 1, false, 0, "", 0},
	//    { "Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
	//    { "A", 16, 1, false, 0, "", 0},
	//    { "Proprietary ID", 8, RES_LOOKUP, false, "=1", "Alarm command", 0},
	//    { "B", 8, 1, false, 0, "", 0},
	//    { "C", 16, 1, false, 0, "", 0},
	//    { "D", 16, 1, false, 0, "", 0},
	//    { "E", 16, 1, false, 0, "", 0}},
	//  },

	{"Simnet: Event Reply: AP command", 130851, false, 12, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Proprietary ID", 8, RES_LOOKUP, false, "=2", "AP command", 0},
		{"B", 16, 1, false, 0, "", 0},
		{"Controlling Device", 8, 1, false, 0, "", 0},
		{"Event", 16, RES_LOOKUP, false, lookupSimnetApEvents, "", 0},
		{"Direction", 8, RES_LOOKUP, false, lookupSimnetDirection, "", 0},
		{"Angle", 16, RES_DEGREES, false, "deg", "", 0},
		{"G", 8, 1, false, 0, "", 0}},
	},

	{"Simnet: Alarm Message", 130856, false, 0x08, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=1857", "Simrad", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Message ID", 16, 1, false, 0, "", 0},
		{"B", 8, 1, false, 0, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"Text", 2040, RES_ASCII, false, 0, "", 0}},
	},

	{"Airmar: Additional Weather Data", 130880, false, 9, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "Airmar", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"Apparent Windchill Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"True Windchill Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Dewpoint", 16, RES_TEMPERATURE, false, "K", "", 0}},
	},

	{"Airmar: Heater Control", 130881, false, 9, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"C", 8, 1, false, 0, "", 0},
		{"Plate Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Air Temperature", 16, RES_TEMPERATURE, false, "K", "", 0},
		{"Dewpoint", 16, RES_TEMPERATURE, false, "K", "", 0}},
	},

	{"Airmar: POST", 130944, false, 8, 0, []Field{
		{"Manufacturer Code", 11, RES_MANUFACTURER, false, "=135", "", 0},
		{"Reserved", 2, 1, false, 0, "", 0},
		{"Industry Code", 3, RES_LOOKUP, false, lookupIndustryCode, "", 0},
		{"Control", 4, RES_LOOKUP, false, lookupAirmarControl, "", 0},
		{"Reserved", 7, 1, false, 0, "", 0},
		{"Number of ID/test result pairs to follow", 8, RES_INTEGER, false, 0, "", 0},
		{"Test ID", 8, RES_LOOKUP, false, lookupAirmarTestId, "See Airmar docs for table of IDs and failure codes; these lookup values are for DST200", 0},
		{"Test result", 8, RES_LOOKUP, false, nil, "Values other than 0 are failure codes. See Airmar docs for description.", 0}},
	},

	{"Actisense: Operating mode", ACTISENSE_BEM + 0x11, false, 0x0e, 0, []Field{
		{"SID", 8, 1, false, 0, "", 0},
		{"Model ID", 16, RES_INTEGER, false, 0, "", 0},
		{"Serial ID", 32, RES_INTEGER, false, 0, "", 0},
		{"Error ID", 32, RES_INTEGER, false, 0, "", 0},
		{"Operating Mode", 16, 1, false, 0, "", 0}},
	},

	{"Actisense: Startup status", ACTISENSE_BEM + 0xf0, false, 0x0f, 0, []Field{
		{"SID", 8, 1, false, 0, "", 0},
		{"Model ID", 16, RES_INTEGER, false, 0, "", 0},
		{"Serial ID", 32, RES_INTEGER, false, 0, "", 0},
		{"Error ID", 32, RES_INTEGER, false, 0, "", 0},
		{"Firmware version", 16, 0.001, false, 0, "", 0},
		{"Reset status", 8, 1, false, 0, "", 0},
		{"A", 8, 1, false, 0, "", 0}},
	},

	{"Actisense: System status", ACTISENSE_BEM + 0xf2, false, 0x22, 0, []Field{
		{"SID", 8, 1, false, 0, "", 0},
		{"Model ID", 16, RES_INTEGER, false, 0, "", 0},
		{"Serial ID", 32, RES_INTEGER, false, 0, "", 0},
		{"Error ID", 32, RES_INTEGER, false, 0, "", 0},
		{"Indi channel count", 8, 1, false, 0, "", 0},
		{"Ch1 Rx Bandwidth", 8, 1, false, 0, "", 0},
		{"Ch1 Rx Load", 8, 1, false, 0, "", 0},
		{"Ch1 Rx Filtered", 8, 1, false, 0, "", 0},
		{"Ch1 Rx Dropped", 8, 1, false, 0, "", 0},
		{"Ch1 Tx Bandwidth", 8, 1, false, 0, "", 0},
		{"Ch1 Tx Load", 8, 1, false, 0, "", 0},
		{"Ch2 Rx Bandwidth", 8, 1, false, 0, "", 0},
		{"Ch2 Rx Load", 8, 1, false, 0, "", 0},
		{"Ch2 Rx Filtered", 8, 1, false, 0, "", 0},
		{"Ch2 Rx Dropped", 8, 1, false, 0, "", 0},
		{"Ch2 Tx Bandwidth", 8, 1, false, 0, "", 0},
		{"Ch2 Tx Load", 8, 1, false, 0, "", 0},
		{"Uni channel count", 8, 1, false, 0, "", 0},
		{"Ch1 Bandwidth", 8, 1, false, 0, "", 0},
		{"Ch1 Deleted", 8, 1, false, 0, "", 0},
		{"Ch1 BufferLoading", 8, 1, false, 0, "", 0},
		{"Ch1 PointerLoading", 8, 1, false, 0, "", 0},
		{"Ch2 Bandwidth", 8, 1, false, 0, "", 0},
		{"Ch2 Deleted", 8, 1, false, 0, "", 0},
		{"Ch2 BufferLoading", 8, 1, false, 0, "", 0},
		{"Ch2 PointerLoading", 8, 1, false, 0, "", 0}},
	},

	{"Actisense: ?", ACTISENSE_BEM + 0xf4, false, 17, 0, []Field{
		{"SID", 8, 1, false, 0, "", 0},
		{"Model ID", 16, RES_INTEGER, false, 0, "", 0},
		{"Serial ID", 32, RES_INTEGER, false, 0, "", 0}},
	},
}

func (pp PgnArray) First(id uint32) (int, Pgn) {
	for i, pgn := range pp {
		if id == pgn.Pgn {
			return i, pgn
		}
	}

	return 0, pp[0]
}

func (pp PgnArray) Last(id uint32) (int, Pgn) {
	for i := len(pp) - 1; i >= 0; i-- {
		if id == pp[i].Pgn {
			return i, pp[i]
		}
	}

	return 0, pp[0]
}
