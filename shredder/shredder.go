package shredder

import (
	"archive/zip"
	"errors"
	"fmt"
	"hotel_billing/billing"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const (
	PATH_JOB     = "path_job_hotel_billing"
	PATH_SUMMARY = "path_summary_file"
	NUMBER_TASK  = "number_tasks_process"
	OS           = "OS"
)

func Print() {
	fmt.Println("Hi shredder!")
}

func loadVars() map[string]string {
	mapVars := make(map[string]string)
	pathSummary, ok := os.LookupEnv(PATH_SUMMARY)
	if ok {
		mapVars[PATH_SUMMARY] = pathSummary
	} else {
		path, err := os.Getwd()
		if err == nil {
			path = path + string(os.PathSeparator) + "SUMMARY"
			mapVars[PATH_SUMMARY] = path
			os.Setenv(PATH_SUMMARY, path)
		}

	}
	pathJob, ok := os.LookupEnv(PATH_JOB)
	if ok {
		mapVars[PATH_JOB] = pathJob
	} else {
		path, err := os.Getwd()
		if err == nil {
			path = path + string(os.PathSeparator) + "INVOCES"
			mapVars[PATH_JOB] = path
			os.Setenv(PATH_JOB, path)
		}

	}
	return mapVars
}
func existsPath(path string) bool {
	_, err := os.OpenFile(path, os.O_RDWR, 0644)
	if errors.Is(err, os.ErrNotExist) {
		// handle the case where the file doesn't exist
		return false
	}
	return true
}
func DistributeShred() error {
	fmt.Println("Here in DistributeShred!")
	mVars := loadVars()
	pathJob := mVars[PATH_JOB]
	if !existsPath(pathJob) {
		return fmt.Errorf("error, work path does not exist : %s", pathJob)
	} else {
		now := time.Now()
		year := now.Year()
		month := now.Month()
		pathAbsJob := fmt.Sprintf("%s%s%d%s%d", pathJob, string(os.PathSeparator), year, string(os.PathSeparator), month)
		pathAbsSummary := fmt.Sprintf("%s%s%d%s%d", mVars[PATH_SUMMARY], string(os.PathSeparator), year, string(os.PathSeparator), month)
		if !existsPath(pathAbsSummary) {
			if err := os.MkdirAll(pathAbsSummary, os.ModePerm); err != nil {
				return err
			}
		}
		fmt.Printf("path_abs_job:%s, path_abs_summary:%s\n", pathAbsJob, pathAbsSummary)

		shred(pathAbsJob, pathAbsSummary)
	}

	return nil
}
func joinSummaryFiles(idx int, nameSuammary <-chan string, wg *sync.WaitGroup, wZip *zip.Writer) {
	defer wg.Done()
	for sFile := range nameSuammary {
		f1, err := os.Open(sFile)

		w1, err := wZip.Create(sFile)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(w1, f1); err != nil {
			panic(err)
		}
		f1.Close()
		e := os.Remove(sFile)
		if e != nil {
			log.Fatal(e)
		}
	}
}
func shred(pathJob string, pathSummary string) {
	fmt.Printf("pathSummary:%s, pathJob:%s", pathSummary, pathJob)
	fZip, err := os.Create(pathSummary + string(os.PathSeparator) + "Summary.zip")
	defer fZip.Close()
	if err != nil {
		log.Fatal(err)
	}
	wZip := zip.NewWriter(fZip)
	var wg sync.WaitGroup
	files, err := ioutil.ReadDir(pathJob)
	if err != nil {
		log.Fatal(err)
	}

	summary := make(chan string)

	for i := 0; i < len(files); i++ {
		wg.Add(1)
		go joinSummaryFiles(i, summary, &wg, wZip)
	}

	for x, file := range files {
		nameFile := fmt.Sprintf("%s%s%s", pathJob, string(os.PathSeparator), file.Name())
		billing.ProcessBilling(summary, nameFile, pathSummary, x)
	}
	close(summary)

	fmt.Println("\nWaiting for goroutines to finish...")
	wg.Wait()
	wZip.Close()
	fmt.Println("Done!")
}
