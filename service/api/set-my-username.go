package api

import (
	api "github.com/zoematr/WASAPhoto/service/api/structs"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/zoematr/WASAPhoto/service/api/reqcontext"
	"github.com/zoematr/WASAPhoto/service/database"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"io"
)

