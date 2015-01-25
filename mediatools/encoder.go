package mediatools

import (
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
)

func (m *MediaInfo) EncodeVideoFile(fileLoc string, fileName string) (err error) {
	file := fileLoc + fileName
	m.log.Infoln(file)

	var outFile string

	outFile, err = m.encodeMP4(fileLoc, fileName)
	if err != nil {
		// ako je mp4 puko kod encoding, ne mozemo nastaviti jer svi ostali zele taj file kao source
		return err
	}
	outFile, err = m.encodeOGG(fileLoc, fileName)
	outFile, err = m.encodeWEBM(fileLoc, fileName)

	m.log.Debugf("%#v", outFile)
	return nil
}

// HandBrakeCLI -i _test/master_1080.mp4 -o _test/out/master_1080_2.mp4 -e x264 -q 22 -r 15 -B 64 -X 480 -O -x level=4.0:ref=9:bframes=16:b-adapt=2:direct=auto:analyse=all:8x8dct=0:me=tesa:merange=24:subme=11:trellis=2:fast-pskip=0:vbv-bufsize=25000:vbv-maxrate=20000:rc-lookahead=60

func (m *MediaInfo) encodeMP4(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	fileNameOut = m.returnBaseFilename(fileName) + ".mp4"
	fileDestination := fileLoc + "out/" + fileNameOut

	maxWidth := "480"
	extraParams := `level=4.0:ref=9:bframes=16:b-adapt=2:direct=auto:analyse=all:8x8dct=0:me=tesa:merange=24:subme=11:trellis=2:fast-pskip=0:vbv-bufsize=25000:vbv-maxrate=20000:rc-lookahead=60`

	out, err := sh.Command("HandBrakeCLI", "-i", fileSource, "-o", fileDestination, "-e", "x264", "-q", "22", "-r", "15", "-B", "64", "-X", maxWidth, "-O", "-x", extraParams).Output()
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	//m.log.Debug(out)
	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)

	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileNameOut, nil

}

// ffmpeg2theora Master_1080.mp4 --two pass --videobitrate 900 -x 1280 -y 720

func (m *MediaInfo) encodeOGG(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	fileNameOut = m.returnBaseFilename(fileName) + ".ogg"
	fileDestination := fileLoc + "out/" + fileNameOut

	maxWidth := "1280"
	maxHeight := "720"

	out, err := sh.Command("ffmpeg2theora", fileSource, "-o", fileDestination, "--two pass", "--videobitrate", "900", "-x", maxWidth, "-y", maxHeight).Output()
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	//m.log.Debug(out)
	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)
	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileDestination, nil
}

/*
ffmpeg -i _test/master_1080.mp4 -pass 1 -passlogfile hattrick.webm -keyint_min 0 -g 250 -skip_threshold 0 -vcodec libvpx -b 600k -s 1280x720 -aspect 16:9 -an -y hattrick.webm
Output file is empty, nothing was encoded (check -ss / -t / -frames parameters if used)

ffmpeg -i _test/master_1080.mp4 -pass 2 -passlogfile hattrick.webm -keyint_min 0 -g 250 -skip_threshold 0 -vcodec libvpx -b 600k -s 1280x720 -aspect 16:9 -acodec libvorbis -y hattrick.webm
*/

func (m *MediaInfo) encodeWEBM(fileLoc string, fileName string) (fileNameOut string, err error) {
	fileSource := fileLoc + fileName
	fileNameOut = m.returnBaseFilename(fileName) + ".webm"
	fileDestination := fileLoc + "out/" + fileNameOut

	// out, err := sh.
	// 	Command("ffmpeg", "-i", fileSource, "-pass", "1", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-an", "-y", fileDestination).
	// 	Command("ffmpeg", "-i", fileSource, "-pass", "2", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-acodec", "libvorbis", "-y", fileDestination).
	// 	Run()
	out, err := sh.
		Command("ffmpeg", "-i", fileSource, "-pass", "1", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-an", "-y", fileDestination).
		Output()

	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	out, err = sh.
		Command("ffmpeg", "-i", fileSource, "-pass", "2", "-passlogfile", fileDestination, "-keyint_min", "0", "-g", "250", "-skip_threshold", "0", "-vcodec", "libvpx", "-b", "600k", "-s", "1280x720", "-aspect", "16:9", "-acodec", "libvorbis", "-y", fileDestination).
		Output()

	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return "", err
	}

	//m.log.Debug(out)
	m.log.Debugf("output:(%s), err(%v)\n", string(out), err)
	ok, err := m.checkIfFileExists(fileDestination)
	if err != nil {
		return "", err
	}
	if ok == false {
		return "", errors.New(fmt.Sprintf("File '%s' does not exists. Encoding failed?", fileDestination))
	}

	return fileDestination, nil
}