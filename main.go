package main

import (
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func main() {
	// Load the database we wish to enrich.
	// writer, err := mmdbwriter.Load("GeoLite2-Country.mmdb", mmdbwriter.Options{})
	writer, err := mmdbwriter.Load("GeoLite2-City.mmdb", mmdbwriter.Options{})
	if err != nil {
		log.Fatal(err)
	}

	// Define and insert the new data.
	_, sreNet, err := net.ParseCIDR("56.0.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	sreData := mmdbtype.Map{
		"AcmeCorp.DeptName": mmdbtype.String("SRE"),
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
			mmdbtype.String("production"),
		},
	}
	if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
		log.Fatal(err)
	}

	_, devNet, err := net.ParseCIDR("56.1.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	devData := mmdbtype.Map{
		"AcmeCorp.DeptName": mmdbtype.String("Development"),
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
		},
	}
	if err := writer.InsertFunc(devNet, inserter.TopLevelMergeWith(devData)); err != nil {
		log.Fatal(err)
	}

	_, mgmtNet, err := net.ParseCIDR("56.2.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	mgmtData := mmdbtype.Map{
		"AcmeCorp.DeptName": mmdbtype.String("Management"),
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
		},
	}
	if err := writer.InsertFunc(mgmtNet, inserter.TopLevelMergeWith(mgmtData)); err != nil {
		log.Fatal(err)
	}

	_, city, err := net.ParseCIDR("110.87.70.227/32")
	if err != nil {
		log.Fatal(err)
	}

	//"country": {
	//	"geoname_id": 1814991,
	//	"iso_code": "CN",
	//	"names": {
	//		"en": "China",
	//		"ja": "中国",
	//		"zh-CN": "中国"
	//	}
	//},
	cityData := mmdbtype.Map{
		"country": mmdbtype.Map{
			"geoname_id": mmdbtype.Uint32(1814991),
			"iso_code":   mmdbtype.String("CN"),
			"names": mmdbtype.Map{
				"en":    mmdbtype.String("China"),
				"zh-CN": mmdbtype.String("中国"),
			},
		},
	}
	if err := writer.InsertFunc(city, inserter.TopLevelMergeWith(cityData)); err != nil {
		log.Fatal(err)
	}

	// Write the newly enriched DB to the filesystem.
	fh, err := os.Create("GeoLite2-City-with-Department-Data.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
