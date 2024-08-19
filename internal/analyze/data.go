package analyze

import "regexp"

var (
	languages = map[string]string{
		"ara":    "ARA",
		"eng":    "ENG",
		"jpn":    "JPN",
		"ger":    "GER",
		"ita":    "ITA",
		"rus":    "RUS",
		"por-br": "POR-BR",
		"por":    "POR",
		"spa-la": "SPA-LA",
		"spa":    "SPA",
		"fre":    "FRE",
		"pol":    "POL",
		"dut":    "DUT",
		"nob":    "NOB",
		"fin":    "FIN",
		"tur":    "TUR",
		"swe":    "SWE",
		"gre":    "GRE",
		"heb":    "HEB",
		"rum":    "RUM",
		"kor":    "KOR",
		"dan":    "DAN",
		"chi":    "CHI",
		"vie":    "VIE",
		"ukr":    "UKR",
		"hun":    "HUN",
		"ces":    "CES",
		"slo":    "SLO",
		"ind":    "IND",
		"tha":    "THA",
		"may":    "MAY",
		"hrv":    "HRV",
		"fil":    "FIL",
		"hin":    "HIN",
	}

	quality = map[string]string{
		"2160": "2160",
		"4k":   "2160",
		"uhd":  "2160",
		"1440": "1440",
		"2k":   "1440",
		"qhd":  "1440",
		"fhd+": "1440",
		"1920": "1080",
		"1080": "1080",
		"fhd":  "1080",
		"1280": "720",
		"720":  "720",
		"hd":   "720",
		"576":  "576",
		"480":  "480",
		"sd":   "480",
		"360":  "360",
		"240":  "240",
		"144":  "144",
	}

	audioCodecs = map[string]string{
		"aac":     "AAC",
		"mp3":     "MP3",
		"truehd":  "TrueHD",
		"true-hd": "TrueHD",
		"opus":    "Opus",
		"flac":    "FLAC",
		"ac3":     "AC-3",
		"eac3":    "E-AC-3",
		"vorbis":  "Vorbis",
		"pcm":     "PCM",
		"alac":    "ALAC",
		"dts":     "DTS",
		"wma":     "WMA",
		"amr":     "AMR",
		"gsm":     "GSM",
	}
	videoCodecs = map[string]string{
		"264":   "H.264",
		"265":   "H.265",
		"vp9":   "VP9",
		"av1":   "AV1",
		"mpeg2": "MPEG-2",
		"mpeg4": "MPEG-4",
		"vp8":   "VP8",
		"hevc":  "H.265",
		"avc":   "H.264",
		"divx":  "DivX",
		"xvid":  "Xvid",
	}

	extensions = []string{
		".mp4",  // MPEG-4 Part 14
		".mkv",  // Matroska Video
		".avi",  // Audio Video Interleave
		".mov",  // QuickTime File Format
		".flv",  // Flash Video
		".webm", // WebM
		".m4v",  // MPEG-4 Video
		".mpg",  // MPEG Video
		".mpeg", // MPEG Video
		".3gp",  // 3GPP
		".3g2",  // 3GPP2
		".vob",  // DVD Video Object
		".f4v",  // Flash MP4 Video
		".divx", // DivX
		".xvid", // Xvid
	}

	numExp         = regexp.MustCompile(`(\d+)`)
	nyaaExp        = regexp.MustCompile(`view\/([^\/]+)\/`)
	hashExp        = regexp.MustCompile(`btih:([a-fA-F0-9]+)&`)
	bracketsExp    = regexp.MustCompile(`\[(.*?)\]`)
	parenthesesExp = regexp.MustCompile(`\((.*?)\)`)
	epstrExp       = regexp.MustCompile(`^[0-9\-_~.]+$`)
	intsExp        = regexp.MustCompile(`([0-9]+(?:\.[0-9]*)?)~([0-9]+(?:\.[0-9]*)?)|([0-9]+(?:\.[0-9]*)?)`)
	seExp          = regexp.MustCompile(`S(\d+)E(\d+\.?\d*)`)
)
