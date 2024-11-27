package bluesky

import (
	"context"
	"fmt"
	"strings"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/dgraph-io/badger/v4"
	"github.com/minio/sha256-simd"
)

type Resolver struct {
	*badger.DB
	*identity.BaseDirectory
}

func NewResolver(path string) (r *Resolver, err error) {
	r = &Resolver{BaseDirectory: &identity.BaseDirectory{}}
	// initialize
	opts := badger.DefaultOptions(path)
	if r.DB, err = badger.Open(opts); chk.E(err) {
		return
	}
	return
}

func (r *Resolver) Close() (err error) { return r.DB.Close() }

func (r *Resolver) Find(did string, c context.Context) (pubkey string, err error) {
	// generate key hash for identity key (16 byte truncated hash)
	k := sha256.Sum256([]byte(did))
	key := k[:]
	var val []byte
	// first search for the DID:
	if err = r.DB.View(func(txn *badger.Txn) (err error) {
		var item *badger.Item
		if item, err = txn.Get(key); err != nil {
			return
		}
		if item == nil {
			err = fmt.Errorf("'%s' not found", did)
			return
		}
		val, err = item.ValueCopy(nil)
		if err != nil {
			return
		}
		return
	}); err != nil {
	}
	if len(val) != 0 {
		log.I.F("found pubkey for %s: %s", did, string(val))
		pubkey = string(val)
		err = nil
		return
	}
	var didd *identity.DIDDocument
	didd, err = r.BaseDirectory.ResolveDID(c, syntax.DID(did))
	if err != nil {
		return
	}
	if len(didd.VerificationMethod) < 1 {
		return
	}
	for _, v := range didd.VerificationMethod {
		if strings.HasPrefix(v.ID, did) && v.Type == "Multikey" {
			// this is the matching pubkey to the identity
			if err = r.DB.Update(func(txn *badger.Txn) (err error) {
				log.I.F("storing pubkey for %s: %s", did, v.PublicKeyMultibase)
				return txn.Set(key, []byte(v.PublicKeyMultibase))
			}); err != nil {
				continue
			}
			// we found it
			break
		}
	}
	return
}
