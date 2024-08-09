package common

import (
	"cloud.google.com/go/firestore"
	"github.com/takara2314/bsam-server/pkg/environment"
	"github.com/takara2314/bsam-server/pkg/racehub"
)

var (
	FirestoreClient *firestore.Client
	Env             *environment.Variables
	Hubs            map[string]*racehub.Hub = make(map[string]*racehub.Hub)
)
