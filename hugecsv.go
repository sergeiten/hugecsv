package hugecsv

import (
	"log"

	"github.com/pkg/errors"
)

//type Item struct {
//	PersonID     int
//	Sex          int
//	StndY        int
//	AgeGroup     int
//	Sido         int
//	Sgg          int
//	KeySeq       int
//	YkihoID      int
//	RecuFrDt     int
//	DsbjtCd      int
//	DmdTramt     int
//	DmdSbrdnAmt  int
//	MainSick     string
//	SubSick      string
//	YkihoGubunCd int
//	YkihoSido    int
//}

// LogPrint prints error message with stack trace without exited program.
func LogPrint(err error, message string) {
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, message))
	}
}

// LogPrintf prints error message with stach trace.
// Arguments are handled in the manner of fmt.Printf.
func LogPrintf(err error, message string, args ...interface{}) {
	if err != nil {
		log.Printf("%+v", errors.Wrapf(err, message, args...))
	}
}

// LogFatal prints error message with stack trace with exited program.
func LogFatal(err error, message string) {
	if err != nil {
		log.Fatalf("%+v", errors.Wrap(err, message))
	}
}
