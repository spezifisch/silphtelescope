package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spezifisch/silphtelescope/pkg/geodex"
	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

var rootCmd = &cobra.Command{
	Use:   "geodexgen",
	Short: "Generate silpht's geodex",
	Long:  `Get Pokestop and Gym data from multiple sources and fill our geodex.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		tStart := time.Now()

		// setup BookOfQuests parser
		boqFiles, _ := cmd.Flags().GetStringArray("boq")
		usingBOQ := len(boqFiles) > 0
		boqOutput := make(chan *geodex.BOQCell)
		boqCancel := make(chan bool)
		boq, err := geodex.NewBOQDB(boqFiles, boqOutput, boqCancel)
		if err != nil {
			log.WithError(err).Error("got invalid boq files")
			return
		}

		// setup GeoDex
		ddbBasePath, _ := cmd.Flags().GetString("geodex")
		ddb := geodex.NewDiskDB(&ddbBasePath)

		// setup MAD MariaDB connection
		sdbHostname, _ := cmd.Flags().GetString("sql-hostname")
		sdbDatabase, _ := cmd.Flags().GetString("sql-database")
		sdbUsername, _ := cmd.Flags().GetString("sql-username")
		sdbPassword, _ := cmd.Flags().GetString("sql-password")
		sdb, err := geodex.NewSQLDB(sdbHostname, sdbDatabase, sdbUsername, sdbPassword)
		if err != nil {
			log.WithError(err).Error("sqldb connection failed")
			return
		}
		defer sdb.Close()

		version, err := sdb.GetVersion()
		if err != nil {
			log.WithError(err).Error("getting version failed")
			return
		}
		log.Infof("Connected to sqldb %s running %s", sdbHostname, version)

		// setup silpht Tile38 connection
		tdbHostname, _ := cmd.Flags().GetString("t-hostname")
		tdbPassword, _ := cmd.Flags().GetString("t-password")
		tdb, err := geodex.NewTDB(tdbHostname, tdbPassword)
		if err != nil {
			log.WithError(err).Error("tdb connection failed")
			return
		}
		defer tdb.Close()

		// setup done
		timeTrack(tStart, "setup")
		tStart = time.Now()

		// get Pokestops from MAD and insert them into tile38
		ps, err := sdb.NewMADPokestopScanner()
		if err != nil {
			log.WithError(err).Error("selecting pokestops failed")
			return
		}
		defer ps.Close()

		psCount := 0
		for ps.Next() {
			p, err := ps.ScanPokestop()
			if err != nil {
				log.WithError(err).Error("parsing pokestop failed")
				return
			}

			f := p.ToFort()
			tdb.InsertFort(f)
			ddb.MergeFort(f)
			psCount++
		}
		log.Infoln("Pokestops read from MAD:", psCount)
		timeTrack(tStart, "mad pokestop import")
		tStart = time.Now()

		// get Gyms from MAD and insert them into tile38
		gs, err := sdb.NewMADGymScanner()
		if err != nil {
			log.WithError(err).Error("selecting gyms failed")
			return
		}
		defer gs.Close()

		gsCount := 0
		for gs.Next() {
			g, err := gs.ScanGym()
			if err != nil {
				log.WithError(err).Error("parsing gym failed")
				return
			}

			f := g.ToFort()
			tdb.InsertFort(f)
			ddb.MergeFort(f)
			gsCount++
		}
		log.Infoln("Gyms read from MAD:", gsCount)
		timeTrack(tStart, "mad gym import")
		tStart = time.Now()

		// get Gym names from BOQ
		if usingBOQ {
			// our signal when BOQ parser is done
			boqDone := make(chan bool)
			// let BOQ reader parse all files, outputting cells to boqOutput
			go func() {
				boq.Run()
				boqDone <- true
			}()

			boqCellCount := 0
			boqPOICount := 0
			boqGymCount := 0
			namesAdded := 0
			namesKept := 0
			for {
				done := false

				select {
				case cell := <-boqOutput:
					boqCellCount++
					for _, poi := range cell.Stops {
						boqPOICount++
						if poi.IsGym {
							boqGymCount++

							// check data from BOQ
							if len(poi.Location.Coordinates) != 2 {
								log.Error("invalid coordinates:", poi.Location.Coordinates)
								return
							}
							if poi.Name == "" {
								continue
							}

							// get gym GUID from tile db
							gymLocation := pogo.Location{
								Latitude:  poi.Location.Coordinates[1],
								Longitude: poi.Location.Coordinates[0],
							}
							tFort, err := tdb.GetNearestFort(gymLocation, 0.1)
							if err != nil {
								// fort doesn't exist in tile38 db, that's ok
								continue
							}

							// get fort from disk
							dFort, err := ddb.GetFort(*tFort.GUID)
							if err != nil {
								// doesn't exist on disk. that's ok
								continue
							}
							if dFort.Name != nil {
								// already has a name
								namesKept++
								continue
							}

							// set name and save to disk
							dFort.Name = &poi.Name
							if err = ddb.SaveFort(dFort); err != nil {
								log.Errorf("couldn't edit fort %s", *tFort.GUID)
								return
							}
							namesAdded++
						}
					}
				case <-boqDone: // boq.Run() ended
					done = true
				}

				if done {
					break
				}
			}

			log.Infof("processed BOQ data: %d cells containing %d POIs with %d gyms",
				boqCellCount, boqPOICount, boqGymCount)
			log.Infof("added names to %d gyms, got %d gyms which already had a name",
				namesAdded, namesKept)

			timeTrack(tStart, "boq import")
			tStart = time.Now()
		}

		// the thing we're doing this for
		defer timeTrack(tStart, "example lookup")
		printFortName(tdb, ddb, 52.5395, 13.4161)
		printFortName(tdb, ddb, 52.5399, 13.4208)
	},
}

// from: https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("> %s took %s", name, elapsed)
}

func printFortName(tdb *geodex.TDB, ddb *geodex.DiskDB, lat, lon float64) {
	searchCenter := pogo.Location{
		Latitude:  lat,
		Longitude: lon,
	}
	sFort, err := tdb.GetNearestFort(searchCenter, 1000)
	if err != nil {
		log.WithError(err).Error("can't find nearby fort")
		return
	}

	fFort, err := ddb.GetFort(*sFort.GUID)
	if err != nil {
		log.WithError(err).Warn("can't find name for fort")
	}

	log.Printf("Fort nearest to (%f,%f): %s",
		searchCenter.Latitude, searchCenter.Longitude, fFort.ToString())
}

func main() {
	rootCmd.PersistentFlags().String("sql-hostname", "", "SQL DB hostname")
	rootCmd.PersistentFlags().String("sql-database", "rocketdb", "SQL DB database for MAD")
	rootCmd.PersistentFlags().String("sql-username", "rocketdb", "SQL DB user")
	rootCmd.PersistentFlags().String("sql-password", "rocketdb", "SQL DB password")

	rootCmd.PersistentFlags().String("t-hostname", "", "Tile38 DB hostname")
	rootCmd.PersistentFlags().String("t-password", "", "Tile38 DB password")

	rootCmd.PersistentFlags().String("geodex", "geodex-storage", "GeoDex storage path")

	rootCmd.PersistentFlags().StringArrayP("boq", "b", []string{}, "BookOfQuests JSON file(s)")

	rootCmd.MarkPersistentFlagRequired("sql-hostname")
	rootCmd.MarkPersistentFlagRequired("t-hostname")
	rootCmd.MarkPersistentFlagRequired("geodex")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
