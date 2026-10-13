package fmap

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ptrTLS struct {
	Cert string
	Key  *string
}

type ptrServer struct {
	Host    string
	Timeout time.Duration
	TLS     *ptrTLS
}

type ptrConfig struct {
	Name   string
	Server ptrServer
	Backup *ptrServer
	Double **ptrServer
}

func TestStructPtr_Paths(t *testing.T) {
	storage, err := Get[ptrConfig]()
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"Name",
		"Server",
		"Server.Host",
		"Server.Timeout",
		"Server.TLS",
		"Server.TLS.Cert",
		"Server.TLS.Key",
		"Backup",
		"Backup.Host",
		"Backup.Timeout",
		"Backup.TLS",
		"Backup.TLS.Cert",
		"Backup.TLS.Key",
		// only single pointer to the struct is traversed
		"Double",
	}, storage.GetAllPaths())

	field := storage.MustFind("Backup.TLS.Cert")
	assert.Equal(t, "Cert", field.GetName())
	assert.Equal(t, "Backup.TLS", field.GetParent().GetStructPath())
	assert.Equal(t, "Backup", field.GetParent().GetParent().GetStructPath())
	assert.Nil(t, field.GetParent().GetParent().GetParent())
	// offset from the start of the struct stored by the pointer
	assert.Equal(t, storage.MustFind("Server.TLS.Key").GetOffset(), storage.MustFind("Backup.TLS.Key").GetOffset())
	assert.Equal(t, storage.MustFind("Backup.TLS.Key").GetOffset(), field.GetOffset()+16)
}

func TestStructPtr_GetNil(t *testing.T) {
	conf := &ptrConfig{}
	storage, err := GetFrom(conf)
	assert.NoError(t, err)

	for _, path := range []string{"Backup.Host", "Backup.Timeout", "Backup.TLS", "Backup.TLS.Cert", "Server.TLS.Key"} {
		field := storage.MustFind(path)
		assert.Nil(t, field.Get(conf), path)
		value, ok := field.GetDereferenced(conf)
		assert.Nil(t, value, path)
		assert.False(t, ok, path)
	}

	// typed nil pointer
	ptr, ok := storage.MustFind("Backup.TLS.Cert").GetPtr(conf).(*string)
	assert.True(t, ok)
	assert.Nil(t, ptr)
	keyPtr, ok := storage.MustFind("Server.TLS.Key").GetPtr(conf).(**string)
	assert.True(t, ok)
	assert.Nil(t, keyPtr)

	// the object is not changed by reading
	assert.Equal(t, ptrConfig{}, *conf)

	// pointer field itself is available
	assert.Equal(t, (*ptrServer)(nil), storage.MustFind("Backup").Get(conf))
	assert.Equal(t, &conf.Backup, storage.MustFind("Backup").GetPtr(conf))
}

func TestStructPtr_Get(t *testing.T) {
	key := "key"
	conf := &ptrConfig{
		Server: ptrServer{TLS: &ptrTLS{Cert: "server cert"}},
		Backup: &ptrServer{Host: "backup", Timeout: time.Second, TLS: &ptrTLS{Cert: "backup cert", Key: &key}},
	}
	storage, err := GetFrom(conf)
	assert.NoError(t, err)

	assert.Equal(t, "server cert", storage.MustFind("Server.TLS.Cert").Get(conf))
	assert.Equal(t, (*string)(nil), storage.MustFind("Server.TLS.Key").Get(conf))
	assert.Equal(t, "backup", storage.MustFind("Backup.Host").Get(conf))
	assert.Equal(t, int64(time.Second), storage.MustFind("Backup.Timeout").Get(conf))
	assert.Equal(t, conf.Backup.TLS, storage.MustFind("Backup.TLS").Get(conf))
	assert.Equal(t, "backup cert", storage.MustFind("Backup.TLS.Cert").Get(conf))
	assert.Equal(t, &key, storage.MustFind("Backup.TLS.Key").Get(conf))

	value, ok := storage.MustFind("Backup.TLS.Key").GetDereferenced(conf)
	assert.True(t, ok)
	assert.Equal(t, "key", value)

	assert.Same(t, &conf.Backup.TLS.Cert, storage.MustFind("Backup.TLS.Cert").GetPtr(conf))
	assert.Same(t, &conf.Backup.TLS, storage.MustFind("Backup.TLS").GetPtr(conf))

	// works with the struct passed by value too
	assert.Equal(t, "backup cert", storage.MustFind("Backup.TLS.Cert").Get(conf))
}

func TestStructPtr_Set(t *testing.T) {
	conf := &ptrConfig{Name: "name"}
	storage, err := GetFrom(conf)
	assert.NoError(t, err)

	// nil pointers in the path are replaced with new structs
	storage.MustFind("Backup.TLS.Cert").Set(conf, "cert")
	if assert.NotNil(t, conf.Backup) && assert.NotNil(t, conf.Backup.TLS) {
		assert.Equal(t, "cert", conf.Backup.TLS.Cert)
	}

	// existing structs are not replaced
	backup, tls := conf.Backup, conf.Backup.TLS
	key := "key"
	storage.MustFind("Backup.TLS.Key").Set(conf, &key)
	storage.MustFind("Backup.Host").Set(conf, "backup")
	storage.MustFind("Backup.Timeout").Set(conf, time.Minute)
	assert.Same(t, backup, conf.Backup)
	assert.Same(t, tls, conf.Backup.TLS)

	storage.MustFind("Server.TLS.Cert").Set(conf, "server cert")

	assert.Equal(t, ptrConfig{
		Name:   "name",
		Server: ptrServer{TLS: &ptrTLS{Cert: "server cert"}},
		Backup: &ptrServer{Host: "backup", Timeout: time.Minute, TLS: &ptrTLS{Cert: "cert", Key: &key}},
	}, *conf)

	// the whole struct can be replaced by the pointer field
	replaced := &ptrTLS{Cert: "replaced"}
	storage.MustFind("Backup.TLS").Set(conf, replaced)
	assert.Same(t, replaced, conf.Backup.TLS)
	assert.Equal(t, "replaced", storage.MustFind("Backup.TLS.Cert").Get(conf))

	// and reset
	storage.MustFind("Backup").Set(conf, (*ptrServer)(nil))
	assert.Nil(t, conf.Backup)
	assert.Nil(t, storage.MustFind("Backup.TLS.Cert").Get(conf))
}

func TestStructPtr_GetFieldByPtr(t *testing.T) {
	conf := &ptrConfig{
		Server: ptrServer{TLS: &ptrTLS{}},
		Backup: &ptrServer{TLS: &ptrTLS{}},
	}
	storage, err := GetFrom(conf)
	assert.NoError(t, err)

	tests := map[string]any{
		"Name":            &conf.Name,
		"Server":          &conf.Server,
		"Server.Host":     &conf.Server.Host,
		"Server.TLS":      &conf.Server.TLS,
		"Server.TLS.Cert": &conf.Server.TLS.Cert,
		"Server.TLS.Key":  &conf.Server.TLS.Key,
		"Backup":          &conf.Backup,
		"Backup.Host":     &conf.Backup.Host,
		"Backup.Timeout":  &conf.Backup.Timeout,
		"Backup.TLS":      &conf.Backup.TLS,
		"Backup.TLS.Cert": &conf.Backup.TLS.Cert,
		"Backup.TLS.Key":  &conf.Backup.TLS.Key,
	}
	for path, ptr := range tests {
		field, err := storage.GetFieldByPtr(conf, ptr)
		if assert.NoError(t, err, path) {
			assert.Equal(t, path, field.GetStructPath())
		}
	}

	// fields of not created structs can not be found
	empty := &ptrConfig{}
	_, err = storage.GetFieldByPtr(empty, &conf.Backup.TLS.Cert)
	assert.Error(t, err)
	_, err = storage.GetFieldByPtr(empty, (*string)(nil))
	assert.Error(t, err)
}

type ptrNode struct {
	Value int
	Next  *ptrNode
}

type ptrLoopA struct {
	Name string
	B    *ptrLoopB
}

type ptrLoopB struct {
	// by value, but the type is already in the path
	A     ptrLoopA
	Other *ptrLoopB
	Count int
}

func TestStructPtr_RecursiveTypes(t *testing.T) {
	storage, err := Get[ptrNode]()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Value", "Next"}, storage.GetAllPaths())

	node := &ptrNode{}
	storage.MustFind("Next").Set(node, &ptrNode{Value: 2})
	assert.Equal(t, 2, node.Next.Value)

	storage, err = Get[ptrLoopA]()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Name", "B", "B.A", "B.Other", "B.Count"}, storage.GetAllPaths())

	loop := &ptrLoopA{}
	storage.MustFind("B.Count").Set(loop, 3)
	storage.MustFind("B.A").Set(loop, ptrLoopA{Name: "nested"})
	assert.Equal(t, &ptrLoopB{A: ptrLoopA{Name: "nested"}, Count: 3}, loop.B)
	assert.Equal(t, ptrLoopA{Name: "nested"}, storage.MustFind("B.A").Get(loop))

	// the same type in different branches is not a recursion
	type branches struct {
		Left  *ptrTLS
		Right *ptrTLS
	}
	storage, err = Get[branches]()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Left", "Left.Cert", "Left.Key", "Right", "Right.Cert", "Right.Key"}, storage.GetAllPaths())
}

func TestStructPtr_Embedded(t *testing.T) {
	type Base struct {
		ID int
	}
	type Model struct {
		*Base
		Title string
	}
	storage, err := Get[Model]()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Base", "Base.ID", "Title"}, storage.GetAllPaths())

	model := &Model{}
	assert.Nil(t, storage.MustFind("Base.ID").Get(model))
	storage.MustFind("Base.ID").Set(model, 7)
	assert.Equal(t, 7, model.ID)
	assert.Equal(t, 7, storage.MustFind("Base.ID").Get(model))
}

func TestStructPtr_ConcurrentRead(t *testing.T) {
	conf := &ptrConfig{Backup: &ptrServer{Host: "backup", TLS: &ptrTLS{Cert: "cert"}}}
	storage, err := GetFrom(conf)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				assert.Equal(t, "cert", storage.MustFind("Backup.TLS.Cert").Get(conf))
				assert.Nil(t, storage.MustFind("Server.TLS.Cert").Get(conf))
				field, err := storage.GetFieldByPtr(conf, &conf.Backup.Host)
				assert.NoError(t, err)
				assert.Equal(t, "Backup.Host", field.GetStructPath())
			}
		}()
	}
	wg.Wait()
}
