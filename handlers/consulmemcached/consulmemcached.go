package consulmemcached

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/netflix/rend/common"
	"github.com/netflix/rend/handlers"
)

type Handler struct {
	mc *memcache.Client
}

func New(mclient *memcache.Client) handlers.HandlerConst {
	return func() (handlers.Handler, error) {
		fmt.Println(">> Starting Proxy ! <<")
		handler := &Handler{
			mc: mclient,
		}
		return handler, nil
	}
}

func (h *Handler) Set(cmd common.SetRequest) error {
	log.WithFields(log.Fields{
		"key":  cmd.Key,
		"data": cmd.Data,
		"ttl":  strconv.FormatInt(int64(cmd.Exptime), 10),
	}).Info("Set operation")

	h.mc.Set(&memcache.Item{Key: string(cmd.Key), Value: cmd.Data})
	return nil
}

func (h *Handler) Add(cmd common.SetRequest) error {
	return nil
}

func (h *Handler) Replace(cmd common.SetRequest) error {
	return nil
}

func (h *Handler) Append(cmd common.SetRequest) error {
	return nil
}

func (h *Handler) Prepend(cmd common.SetRequest) error {
	return nil
}

func (h *Handler) Get(cmd common.GetRequest) (<-chan common.GetResponse, <-chan error) {
	dataOut := make(chan common.GetResponse, len(cmd.Keys))
	defer close(dataOut)
	errorOut := make(chan error)
	defer close(errorOut)

	log.Debug("Get operation")

	for idx, bk := range cmd.Keys {
		item, err := h.mc.Get(string(bk))

		if err != nil {
			log.WithError(err).Debug("Get fail")
			//if err == common.ErrKeyNotFound {
			dataOut <- common.GetResponse{
				Miss:   true,
				Quiet:  cmd.Quiet[idx],
				Opaque: cmd.Opaques[idx],
				Key:    bk,
			}
			continue
		}

		log.WithFields(log.Fields{
			"key":  item.Key,
			"data": item.Value,
			"ttl":  strconv.FormatInt(int64(item.Expiration), 10),
		}).Info("Get operation")

		dataOut <- common.GetResponse{
			Miss:   false,
			Quiet:  cmd.Quiet[idx],
			Opaque: cmd.Opaques[idx],
			Flags:  item.Flags,
			Key:    bk,
			Data:   item.Value,
		}
	}
	return dataOut, errorOut
}

func (h *Handler) GetE(cmd common.GetRequest) (<-chan common.GetEResponse, <-chan error) {
	return nil, nil
}

func (h *Handler) GAT(cmd common.GATRequest) (common.GetResponse, error) {
	return common.GetResponse{}, nil
}

func (h *Handler) Delete(cmd common.DeleteRequest) error {
	return nil
}

func (h *Handler) Touch(cmd common.TouchRequest) error {
	return nil
}

func (h *Handler) Close() error {
	return nil
}
