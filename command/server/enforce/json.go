package enforce

import (
	"net/http"
	"io/ioutil"
	"io"
	"strings"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

// according to https://www.iana.org/assignments/media-types/application/json
// application/json does not have an optional parameter such as charset
// i.e. charset=UTF-8 should not be appended to the content type
//	application/json; charset=UTF-8 should not be used or accepted
const (
	HEADER_JSON_CONTENT_TYPE = "application/json"
)

func EnforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()
		logger := logging.Logger( reqCtx )

		logger.Debug("enforcing request content-type",
			zap.String("enforce_type", "json"),
		)

		// check the request's header file to determine if the content-type was specified
		// application/json should be specified as the content type
		contentType := r.Header.Get( "Content-Type" )
		if contentType == "" {
			logger.Error( "content-type not specified",
				zap.String( "status", "error" ),
			)

			w.WriteHeader( http.StatusBadRequest )
			return
		} else if strings.ToLower(contentType) != HEADER_JSON_CONTENT_TYPE {
			logger.Error( "unsupported media type",
				zap.String("supplied_type", contentType ),
				zap.String("expected_type", HEADER_JSON_CONTENT_TYPE ),
				zap.String( "status", "error" ),
			)

			w.WriteHeader( http.StatusBadRequest )
			return
		}

		// Check the request's message body size, return error if no content was supplied
		if r.ContentLength == 0 {
			logger.Error( "unable to read request body",
				zap.Int64("content_len", r.ContentLength ),
				zap.String( "status", "error" ),
			)

			w.WriteHeader( http.StatusBadRequest )
			return
		}

		// only the first 512 bytes are used to detect the content type
		body, err := ioutil.ReadAll( io.LimitReader(r.Body, 512 ))
		if err != nil {
			logger.Error( "unable to process message body",
				zap.String( "status", "error" ),
				zap.String( "error", err.Error() ),
			)
		}

		bodyType := http.DetectContentType(body)
		if bodyType != HEADER_JSON_CONTENT_TYPE {
			logger.Error( "unsupported media type",
				zap.String("detected_type", bodyType ),
				zap.String("expected_type", HEADER_JSON_CONTENT_TYPE ),
				zap.String( "status", "error" ),
			)

			w.WriteHeader( http.StatusUnsupportedMediaType )
			return
		}

		next.ServeHTTP(w, r)
	})
}

