package main

import (
	imagestore "iapoc_elephanttrunkarch/ImageStore/go/ImageStore"
	"os"
	"flag"
	"github.com/golang/glog"
)

func checkErr(resp string, err error) {
	if err != nil {
		glog.Errorf("Error: %v", err)
	} else {
		if resp != "" {
			glog.Infof("Response: %s", resp)
		}
	}
}

func readFile(filename string) []byte {

	file, err := os.Open(filename)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		glog.Errorf("Error: %v", err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	return buffer
}

func writeFile(filename string, message string) {
	f, err := os.Create(filename)
	if err != nil {
		glog.Errorf("Error: %v", err)
	}
	defer f.Close()
	n3, _ := f.WriteString(message)
	glog.Infof("wrote %d bytes\n", n3)
	f.Sync()
}

func main() {

	var inputFile string
	var outputFile string
	flag.StringVar(&inputFile, "input_file", "", "input file path to write to ImageStore")
	flag.StringVar(&outputFile, "output_file", "", "output file that gets" +
				   "created from ImageStore read")

	flag.Parse()

	if len(os.Args) < 2 {
		glog.Errorf("Usage: go run DataAgent/da_grpc/test/clientTest.go " +
			"-input_file=<input_file_path> [-output_file=<output_file_path>]")
		os.Exit(-1)
	}

	flag.Lookup("alsologtostderr").Value.Set("true")
	defer glog.Flush()


	imagestore, err := imagestore.NewImageStore()
	if err != nil {
		glog.Errorf("Failed to instantiate ImageStore. Error: %s", err)
	} else {
		var err error
		var data string
		var keyname string

		data, err = imagestore.Read("inmem")
		checkErr(data, err)

		imagestore.SetStorageType("inmemory")
		keyname, err = imagestore.Store([]byte("vivek"))
		checkErr(keyname, err)

		data, err = imagestore.Read(keyname)
		checkErr("Read success", err)

		err = imagestore.Remove(keyname)
		checkErr("Keyname " + keyname + " removed successfully", err)

		//Reading Files
		inputData := readFile(inputFile)
		keyname, err = imagestore.Store(inputData)
		checkErr(keyname, err)

		data, err = imagestore.Read(keyname)
		checkErr("Read success", err)

		//Writing files
		writeFile(outputFile, data)
	}
}