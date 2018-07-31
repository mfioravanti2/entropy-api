package scoringdb

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/model/request"
)

type ReqAttribute struct {
	ReqAttributeID uint	`gorm:"primary_key"`
	ReqRecordId	uint

	Mnemonic 	string
	Format 		string
	Tag 		string
	Nationality string
}

type ReqAttributes []ReqAttribute

type ReqRecord struct {
	ReqRecordId  uint		`gorm:"primary_key"`

	RequestId 	string
	Time 		time.Time
	Locale 		string
	People 		int
	Attributes	[]ReqAttribute	`gorm:"ForeignKey:ReqRecordId"`
}

// Convert a request object into a request summary
func NewReqRecord( req *request.Request, reqId string, t time.Time ) ( *ReqRecord, error ) {
	var r ReqRecord

	r.RequestId = reqId
	r.Time = t.UTC()
	r.Locale = req.Locale
	r.People = len(req.People)

	for _, p := range req.People {
		for _, a := range p.Attributes {
			t, e := NewReqAttribute( p.Nationality, a )
			if e == nil {
				r.Attributes = append( r.Attributes, t )
			}
		}
	}

	return &r, nil
}

// Convert a request attribute into a request attribute summary
func NewReqAttribute( nationality string, r request.Attribute ) ( ReqAttribute, error ) {
	var a ReqAttribute

	a.Nationality = nationality
	a.Mnemonic = r.Mnemonic
	a.Format = r.Format
	a.Tag = r.Tag

	return a, nil
}

// Save a request summary to the data store
func (ds *DataStore) SaveRequest( ctx context.Context, r *ReqRecord ) error {
	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "SaveRequest", "scoringdb" )
	} else {
		ctx = logging.WithFuncId( ctx, "SaveRequest", "scoringdb" )
	}

	logger := logging.Logger( ctx )

	logger.Info("logging scoring request",
	)

	if !ds.Active {
		return nil
	}

	err := ds.g.Create(r).Error
	ds.LastUse = time.Now().UTC()

	return err
}

// Determine if the data store is configured to store request summaries
func (ds *DataStore) readyRequest() bool {
	if !ds.Active {
		return true
	}

	if ds != nil && ds.g != nil {
		if ds.g.HasTable(&ReqRecord{}) {
			if ds.g.HasTable(&ReqAttribute{}) {
				return true
			}
		}
	}

	ctx := logging.WithFuncId( context.Background(), "readyRequest", "scoringdb" )

	logger := logging.Logger( ctx )
	logger.Info("checking dataStore status",
		zap.String( "recordType", "request" ),
		zap.Bool( "ready", false ),
	)

	return false
}
