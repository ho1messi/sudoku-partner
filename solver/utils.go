package solver

var digitStrs = []string{
	".", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}

var easyPuzzles = [][]string{
	{
		".6.593...9.1...5...3.4...9.1.8.2...44..3.9..12...1.6.9.8...6.2...4...8.7...785.1.",
		"762593148941278536835461792198627354476359281253814679387146925514932867629785413",
	},
	{
		"981..3.4.....7925..7.1.6.83.9.4.75.2..8.1.7..7.36.5.1.31.7.4.9..6923.....5.9..324",
		"981523647634879251275146983196487532548312769723695418312754896469238175857961324",
	},
	{
		"981..3.4.....7925..751.6.83.9.4.75.25.831.7.97236.541831.7.4.9..6923.1.5.579.1324",
		"981523647634879251275146983196487532548312769723695418312754896469238175857961324",
	},
	{
		"981..3.4.6...7925.2751.698319.4.75325.831.7.97236.541831.7.4.9..6923.1.5.579.1324",
		"981523647634879251275146983196487532548312769723695418312754896469238175857961324",
	},
}

var hardPuzzles = [][]string{
	{
		"8..........36......7..9.2...5...7.......4.7.....1.5.3...1....68..85...1..9....4..",
		"",
	},
}

var veryHardPuzzles = [][]string{
	{
		"8..........36......7..9.2...5...7.......457.....1...3...1....68..85...1..9....4..",
		"",
	},
}
