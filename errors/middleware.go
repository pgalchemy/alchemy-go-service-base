package errors

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

// Middleware provides a gin middleware to collect errors and deliver a response
func Middleware(c *gin.Context) {
	// Wait for handler
	c.Next()

	// Only continue if response not yet written
	if !c.Writer.Written() {
		var err error
		if len(c.Errors) > 0 {
			err = c.Errors[0].Err
		} else {
			err = ResourceNotFoundError("Route not found.")
		}

		switch v := err.(type) {
		case *HTTPError:
			c.JSON(v.StatusCode, gin.H{"success": false, "error": v})
			return
		case pg.Error:
			httperr := getPostgresError(v.Field('C'))
			c.JSON(httperr.StatusCode, gin.H{"success": false, "error": httperr})
			return
		case *json.SyntaxError:
			httperr := InvalidContentError("Invalid data provided.")
			c.JSON(httperr.StatusCode, gin.H{"success": false, "error": httperr})
			return
		default:
			httperr := InvalidContentError("Invalid content.")
			c.JSON(httperr.StatusCode, gin.H{"success": false, "error": httperr})
			return
		}
	}
}

func getPostgresError(code string) *HTTPError {
	switch code {
	case "23502":
		return MissingParameterError("Missing parameter.")
	case "23505":
		return ConflictError("Resource already exists.")
	default:
		return InvalidContentError("Invalid content.")
	}
}
