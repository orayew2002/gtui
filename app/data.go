package main

type movie struct {
	Title    string
	Storage  string
	Duration string
	Size     string
	Status   string
	Active   bool
	Tracks   string
}

var movieRows = []movie{
	{"Garagum gijeleri", "disk-01", "1:52:10", "2.1 GB", "ready", true, "tk ru en"},
	{"Köneürgenç", "disk-01", "1:34:02", "1.7 GB", "ready", true, "tk ru"},
	{"Ýaz ýagşy", "disk-02", "2:04:55", "2.8 GB", "encoding", false, "tk"},
	{"Altyn asyr", "disk-02", "1:12:40", "1.2 GB", "ready", true, "tk ru en"},
	{"Deňiz kenary", "disk-03", "1:46:18", "2.0 GB", "failed", false, "tk ru"},
	{"Bagtyýar nesil", "disk-01", "1:28:33", "1.5 GB", "ready", true, "tk"},
}

var serialRows = [][]string{
	{"Gumly obada", "3", "48", "18 GB", "ready", "on"},
	{"Şäher ertekisi", "2", "24", "9 GB", "ready", "on"},
	{"Gizlin jaň", "1", "12", "4 GB", "encoding", "off"},
	{"Maşgala", "5", "92", "41 GB", "ready", "on"},
}

var bookRows = [][]string{
	{"Ruhnama", "A. Myradow", "24", "480 MB", "ready", "on"},
	{"Gorkut ata", "G. Seýidow", "18", "320 MB", "ready", "on"},
	{"Şükür bagşy", "M. Annaýew", "12", "210 MB", "encoding", "off"},
}

var songRows = [][]string{
	{"Aýjemal", "Ş. Söýünow", "Ýaz", "3:42", "18 402", "on"},
	{"Garaşaryn", "M. Durdyýew", "Ýol", "4:05", "12 118", "on"},
	{"Bagtym", "A. Nuryýewa", "Säher", "3:18", "9 806", "on"},
	{"Ýürek sesi", "Ş. Söýünow", "Ýaz", "3:56", "7 445", "off"},
}

var storageRows = [][]string{
	{"disk-01", "/mnt/media01", "2 TB", "1.8 TB", "active"},
	{"disk-02", "/mnt/media02", "2 TB", "1.2 TB", "active"},
	{"disk-03", "/mnt/media03", "1 TB", "0.7 TB", "active"},
	{"usb pool", "/mnt/usb", "512 GB", "310 GB", "pending"},
}

var logLines = []struct{ Time, Text string }{
	{"09:41:02", "release usb-014 finished, 182 items"},
	{"09:38:55", "movie Garagum gijeleri encoded in 4m12s"},
	{"09:31:10", "storage disk-02 scanned, 12 new files"},
	{"09:28:44", "station mary depot sync delayed, warn"},
	{"09:21:03", "admin signed in from 10.0.0.14"},
	{"09:18:30", "encode failed: Deňiz kenary track 2"},
	{"09:12:58", "banner Kitap güni disabled"},
	{"09:04:21", "book Şükür bagşy queued"},
}

var storages = []string{"disk-01", "disk-02", "disk-03", "usb pool"}
