// Command distance computes a great-circle distance BOTH ways — the haversine
// formula (P2) and the stereographic chordal metric (P1) — and prints them side
// by side with the independent arccos reference. This is the worked "both
// methods" example that the papers reference but do not print in full.
//
// Usage:
//
//	go run ./cmd/distance                         # defaults to London -> New York
//	go run ./cmd/distance 51.5074 -0.1278 40.7128 -74.0060
//	go run ./cmd/distance -r 6378.137 <lat1> <lon1> <lat2> <lon2>   # WGS-84 mean radius
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"haversine-check/geo"
)

func main() {
	radius := flag.Float64("r", geo.EarthRadiusKm, "sphere radius (distance is returned in this unit)")
	flag.Parse()

	lat1, lon1, lat2, lon2 := 51.5074, -0.1278, 40.7128, -74.0060 // London -> New York
	label := "London -> New York (default)"
	if args := flag.Args(); len(args) == 4 {
		var err error
		if lat1, err = parse(args[0]); err == nil {
			lon1, err = parse(args[1])
		}
		if err == nil {
			lat2, err = parse(args[2])
		}
		if err == nil {
			lon2, err = parse(args[3])
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "could not parse coordinates:", err)
			os.Exit(1)
		}
		label = "custom coordinates"
	} else if len(flag.Args()) != 0 {
		fmt.Fprintln(os.Stderr, "expected 0 or 4 positional args: lat1 lon1 lat2 lon2")
		os.Exit(1)
	}

	hav := geo.Haversine(lat1, lon1, lat2, lon2, *radius)
	cho := geo.Chordal(lat1, lon1, lat2, lon2, *radius)
	ref := *radius * geo.CentralAngle(lat1, lon1, lat2, lon2)

	fmt.Printf("%s\n", label)
	fmt.Printf("  from (%.4f, %.4f) to (%.4f, %.4f), radius %.4f\n\n", lat1, lon1, lat2, lon2, *radius)
	fmt.Printf("  haversine formula (P2)      : %.6f\n", hav)
	fmt.Printf("  stereographic chordal (P1)  : %.6f\n", cho)
	fmt.Printf("  reference  arccos(u1 . u2)  : %.6f\n", ref)
	fmt.Printf("\n  |haversine - chordal|       : %.3e\n", math.Abs(hav-cho))
	fmt.Printf("  |haversine - reference|     : %.3e\n", math.Abs(hav-ref))
}

func parse(s string) (float64, error) { return strconv.ParseFloat(s, 64) }
