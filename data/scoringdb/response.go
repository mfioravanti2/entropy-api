package scoringdb

import (
	"context"
	"time"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/model/response"
	"go.uber.org/zap"
)

type RespAttribute struct {
	RespAttributeID uint	`gorm:"primary_key"`
	RespRecordId	uint

	Mnemonic 	string
	Format 		string
	Tag 		string
	Nationality string
	Score		float64
}

type RespAttributes []RespAttribute

type RespRecord struct {
	RespRecordId  uint		`gorm:"primary_key"`

	RequestId 	string
	Time 		time.Time
	Locale 		string
	Score		float64
	People 		int
	Attributes	[]RespAttribute	`gorm:"ForeignKey:RespRecordId"`
}

// Convert a response object into a response summary
func NewRespRecord( resp *response.Response, reqId string, t time.Time ) ( *RespRecord, error ) {
	var r RespRecord

	r.RequestId = reqId
	r.Time = t.UTC()
	r.Locale = resp.Data.Locale
	r.People = len(resp.Data.People)
	r.Score = resp.Data.Score

	for _, p := range resp.Data.People {
		for _, a := range p.Attributes {
			t, e := NewRespAttribute( p.Nationality, a )
			if e == nil {
				r.Attributes = append( r.Attributes, t )
			}
		}
	}

	return &r, nil
}

// Convert a response attribute into a response attribute summary
func NewRespAttribute( nationality string, r response.Attribute ) ( RespAttribute, error ) {
	var a RespAttribute

	a.Nationality = nationality
	a.Mnemonic = r.Mnemonic
	a.Format = r.Format
	a.Tag = r.Tag
	a.Score = r.Score

	return a, nil
}

// Save a response summary to the data store
func (ds *DataStore) SaveResponse( ctx context.Context, r *RespRecord ) error {
	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "SaveResponse", "scoringdb" )
	} else {
		ctx = logging.WithFuncId( ctx, "SaveResponse", "scoringdb" )
	}

	logger := logging.Logger( ctx )

	logger.Info("logging scoring response",
	)

	if !ds.Active {
		return nil
	}

	err := ds.g.Create(r).Error
	ds.LastUse = time.Now().UTC()

	return err
}

// Determine if the data store is configured to store request summaries
func (ds *DataStore) readyResponse() bool {
	if !ds.Active {
		return true
	}

	if ds != nil && ds.g != nil {
		if ds.g.HasTable(&RespRecord{}) {
			if ds.g.HasTable(&RespAttribute{}) {
				return true
			}
		}
	}

	ctx := logging.WithFuncId( context.Background(), "readyResponse", "scoringdb" )

	logger := logging.Logger( ctx )
	logger.Info("checking dataStore status",
		zap.String( "recordType", "request" ),
		zap.Bool( "ready", false ),
	)

	return false
}
