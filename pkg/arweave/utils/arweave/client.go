package arweave

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/config"
	"github.com/go-resty/resty/v2"
	"github.com/teivah/onecontext"
)

type Client struct {
	*BaseClient
}

func NewClient(ctx context.Context, config *config.Config) (self *Client) {
	self = new(Client)
	self.BaseClient = newBaseClient(ctx, config)
	return
}

// https://docs.arweave.org/developers/server/http-api#network-info
func (self *Client) GetNetworkInfo(ctx context.Context) (out *NetworkInfo,
	err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetResult(&NetworkInfo{}).
		ForceContentType("application/json").
		Get("/info")
	if err != nil {
		return
	}

	out, ok := resp.Result().(*NetworkInfo)
	if !ok {
		err = ErrFailedToParse
		return
	}

	return
}

// https://docs.arweave.org/developers/server/http-api#peer-list
func (self *Client) GetPeerList(ctx context.Context) (out []string, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetResult([]string{}).
		ForceContentType("application/json").
		Get("/peers")
	if err != nil {
		return
	}

	peers, ok := resp.Result().(*[]string)
	if !ok {
		err = ErrFailedToParse
		return
	}

	return *peers, nil
}

func (self *Client) CheckPeerConnection(ctx context.Context,
	peer string) (out *NetworkInfo, duration time.Duration, err error) {
	// Disable retrying request with different peer
	ctx = context.WithValue(ctx, ContextDisablePeers, true)
	ctx = context.WithValue(ctx, ContextForcePeer, peer)

	// Set timeout
	ctx, cancel := context.WithTimeout(ctx,
		self.config.Arweave.CheckPeerTimeout)
	defer cancel()

	self.mtx.RLock()
	ctx, _ = onecontext.Merge(self.ctx, ctx)
	self.mtx.RUnlock()

	resp, err := self.client.R().
		SetContext(ctx).
		ForceContentType("application/json").
		SetResult(&NetworkInfo{}).
		Get("/info")
	if err != nil {
		return
	}

	out, ok := resp.Result().(*NetworkInfo)
	if !ok {
		err = ErrFailedToParse
		return
	}

	duration = resp.Time()

	return
}

// https://docs.arweave.org/developers/server/http-api#get-block-by-height
func (self *Client) GetBlockByHeight(ctx context.Context,
	height int64) (out *Block, resp *resty.Response, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err = req.
		SetResult(&Block{}).
		SetPathParam("height", strconv.FormatInt(height, 10)).
		ForceContentType("application/json").
		Get("/block/height/{height}")
	if err != nil {
		return
	}

	out, ok := resp.Result().(*Block)
	if !ok {
		err = ErrFailedToParse
		return
	}

	// self.log.WithField("block", string(resp.Body())).Info("Block")

	return
}

// https://docs.arweave.org/developers/server/http-api#blocks
func (self *Client) GetBlockByHash(ctx context.Context,
	hash string) (out *Block, resp *resty.Response, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err = req.
		SetResult(&Block{}).
		SetPathParam("hash", hash).
		ForceContentType("application/json").
		Get("/block/hash/{hash}")
	if err != nil {
		return
	}

	out, ok := resp.Result().(*Block)
	if !ok {
		err = ErrFailedToParse
		return
	}

	return
}

// https://docs.arweave.org/developers/server/http-api#get-transaction-by-id
func (self *Client) GetTransactionById(ctx context.Context,
	id string) (out *Transaction, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetResult(&Transaction{}).
		SetError(&Error{}).
		ForceContentType("application/json").
		SetPathParam("id", id).
		Get("/tx/{id}")
	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			err = ErrNotFound
			return
		}

		msg, ok := resp.Error().(*Error)
		if !ok {
			err = ErrBadResponse
			return
		}

		if msg.Error != "" {
			err = errors.New(msg.Error)
			return
		}

		if resp.StatusCode() == http.StatusGone {
			err = ErrOverspend
			return
		}

		err = errors.New(string(http.StatusText(resp.StatusCode())))
		return
	}

	if err != nil {
		if resp != nil && string(resp.Body()) == "Pending" {
			err = ErrPending
		}
		return
	}

	out, ok := resp.Result().(*Transaction)
	if !ok {
		if string(resp.Body()) == "Pending" {
			err = ErrPending
		} else {
			err = ErrFailedToParse
		}

		return
	}

	return
}

// https://docs.arweave.org/developers/server/http-api#get-transaction-offset-and-size
func (self *Client) GetTransactionOffsetInfo(ctx context.Context,
	id string) (out *OffsetInfo, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetResult(&OffsetInfo{}).
		ForceContentType("application/json").
		SetPathParam("id", id).
		Get("/tx/{id}/offset")
	if err != nil {
		return
	}

	out, ok := resp.Result().(*OffsetInfo)
	if !ok {
		err = ErrFailedToParse
		return
	}

	if !out.Offset.Valid || !out.Size.Valid {
		err = ErrBadResponse
		return
	}

	return
}

func (self *Client) getChunk(ctx context.Context,
	offset big.Int) (out *ChunkData, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetPathParam("offset", offset.String()).
		SetResult(&ChunkData{}).
		ForceContentType("application/json").
		Get("/chunk/{offset}")
	if err != nil {
		return
	}

	out, ok := resp.Result().(*ChunkData)
	if !ok {
		err = ErrFailedToParse
		return
	}

	return
}

func (self *Client) GetChunks(ctx context.Context,
	tx *Transaction) (out bytes.Buffer, err error) {
	// Download chunks
	info, err := self.GetTransactionOffsetInfo(ctx, tx.ID.Base64())
	if err != nil {
		return
	}

	zero := big.NewInt(0)
	finish := &info.Offset.Int
	size := &info.Size.Int
	offset := finish.Sub(finish, size)
	offset = finish.Add(offset, big.NewInt(1))

	self.log.WithField("size", size.String()).Trace("Downloading")

	for size.Cmp(zero) > 0 {
		var chunk *ChunkData
		chunk, err = self.getChunk(ctx, *offset)
		if err != nil {
			return
		}

		// Chunk field is already decoded from base64
		out.Write(chunk.Chunk.Bytes())

		// Are there more chunks?
		chunkSize := big.NewInt(int64(len(chunk.Chunk.Bytes())))
		size = size.Sub(size, chunkSize)
		offset = offset.Add(offset, chunkSize)

		// self.log.WithField("chunk_size", chunkSize).WithField("new_size", size).Trace("One chunk")
	}

	if out.Len() != int(tx.DataSize.Int64()) {
		err = ErrDataSizeMismatch
		return
	}

	return
}

// https://docs.arweave.org/developers/server/http-api#get-transaction-field
func (self *Client) GetTransactionDataById(ctx context.Context,
	tx *Transaction) (out bytes.Buffer, err error) {
	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetPathParam("id", tx.ID.Base64()).
		Get("/tx/{id}/data")
	if err != nil {
		return
	}

	// Data is base64 url encoded
	data, err := base64.RawURLEncoding.DecodeString(string(resp.Body()))
	if err != nil {
		return
	}

	out.Write(data)
	if out.Len() > 0 {
		// Do the checks and return
		if out.Len() != int(tx.DataSize.Int64()) {
			err = ErrDataSizeMismatch
			return
		}

		return
	}

	out, err = self.GetChunks(ctx, tx)
	if err == nil {
		return
	}

	var err2 error
	out, err2 = self.GetCachedTransactionDataById(ctx, tx)
	if err2 == nil {
		return
	}

	return out, errors.Join(err, err2)
}

func (self *Client) GetCachedTransactionDataById(ctx context.Context,
	tx *Transaction) (out bytes.Buffer, err error) {
	ctx = context.WithValue(ctx, ContextDisablePeers, true)

	req, cancel := self.Request(ctx)
	defer cancel()

	resp, err := req.
		SetPathParam("id", tx.ID.Base64()).
		Get("{id}")
	if err != nil {
		return
	}

	out.Write(resp.Body())

	if out.Len() != int(tx.DataSize.Int64()) {
		err = ErrDataSizeMismatch
		return
	}

	return
}
