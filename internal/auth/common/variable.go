package common

import (
	"cloud.google.com/go/firestore"
	"github.com/takara2314/bsam-server/pkg/environment"
)

var (
	FirestoreClient *firestore.Client
	Env             *environment.Variables
)
