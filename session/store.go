package session

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Store struct {
	name       string
	id         string
	handler    Handler
	attributes map[string]interface{}
}

func NewStore(name string, handler Handler) *Store {
	s := &Store{
		name:       name,
		handler:    handler,
		attributes: make(map[string]interface{}),
	}
	return s
}

func (s *Store) GetHandler() Handler {
	return s.handler
}

func (s *Store) GetId() string {
	return s.id
}

func (s *Store) SetId(id string) {
	if len(id) < 1 {
		id = generateSessionId()
	}

	s.id = id
}

func (s *Store) GetName() string {
	return s.name
}

func (s *Store) Start() {
	data := s.handler.Read(s.GetId())

	decodeData, _ := base64.StdEncoding.DecodeString(data)

	udata := make(map[string]interface{})

	json.Unmarshal(decodeData, &udata)

	for k, v := range udata {
		s.attributes[k] = v
	}
}

func (s *Store) Get(name string, value ...interface{}) interface{} {
	if v, ok := s.attributes[name]; ok {
		return v
	}
	if len(value) > 0 {
		return value[0]
	}
	return nil
}

func (s *Store) Set(name string, value interface{}) {
	s.attributes[name] = value
}

func (s *Store) All() map[string]interface{} {
	return s.attributes
}

func (s *Store) Remove(name string) interface{} {
	value := s.Get(name)
	delete(s.attributes, name)
	return value
}

func (s *Store) Forget(names ...string) {
	for _, name := range names {
		delete(s.attributes, name)
	}
}

func (s *Store) Clear() {
	s.attributes = nil
}

func (s *Store) Save() {
	data, _ := json.Marshal(s.attributes)

	encodeData := base64.StdEncoding.EncodeToString(data)
	s.handler.Write(s.GetId(), encodeData)
}

func generateSessionId() string {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	b := make([]byte, 48)
	io.ReadFull(rand.Reader, b)
	id = id + base64.URLEncoding.EncodeToString(b)

	h := sha1.New()
	h.Write([]byte(id))

	return fmt.Sprintf("%x", h.Sum(nil))
}
