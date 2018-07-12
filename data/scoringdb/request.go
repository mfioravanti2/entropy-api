package scoringdb

import (
	"context"
	"time"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/model/request"
)

type ReqAttribute struct {
	ReqAttributeID uint	`gorm:"primary_key"`
	RecordID	uint

	Mnemonic 	string
	Format 		string
	Tag 		string
	Nationality string
}

type ReqAttributes []ReqAttribute

type ReqRecord struct {
	ReqRecordId  	uint		`gorm:"primary_key"`

	RequestId 	string
	Time 		time.Time
	Locale 		string
	People 		int
	Attributes	[]ReqAttribute	`gorm:"ForeignKey:ReqRecordId"`
}

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

func NewReqAttribute( nationality string, r request.Attribute ) ( ReqAttribute, error ) {
	var a ReqAttribute

	a.Nationality = nationality
	a.Mnemonic = r.Mnemonic
	a.Format = r.Format
	a.Tag = r.Tag

	return a, nil
}

func (ds *DataStore) SaveRequest( ctx context.Context, r *ReqRecord ) error {
	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "SaveRequest", "scoringdb" )
	} else {
		ctx = logging.WithFuncId( ctx, "SaveRequest", "scoringdb" )
	}

	logger := logging.Logger( ctx )

	logger.Info("logging scoring request",
		//zap.String( "personId", p.PersonID ),
	)

	ds.g.Create(r)

	return nil
}
