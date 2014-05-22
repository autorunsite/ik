package ik

import (
	"log"
	"io"
	"math/rand"
)

type FluentRecord struct {
	Tag       string
	Timestamp uint64
	Data      map[string]interface{}
}

type TinyFluentRecord struct {
	Timestamp uint64
	Data      map[string]interface{}
}

type FluentRecordSet struct {
	Tag       string
	Records   []TinyFluentRecord
}

type Port interface {
	Emit(recordSets []FluentRecordSet) error
}

type Spawnee interface {
	Run() error
	Shutdown() error
}

type PluginInstance interface {
	Spawnee
	Factory() Plugin
}

type Input interface {
	PluginInstance
	Port() Port
}

type Output interface {
	PluginInstance
	Port
}

type MarkupAttributes int

const (
	Red        = 0x00001
	Green      = 0x00002
	Yellow     = 0x00003
	Blue       = 0x00004
	Magenta    = 0x00005
	Cyan       = 0x00006
	White      = 0x00007
	Embolden   = 0x10000
	Underlined = 0x20000
)

type MarkupChunk struct {
	Attrs MarkupAttributes
	Text  string
}

type Markup struct {
	Chunks []MarkupChunk
}

type ScoreValueFetcher interface {
	PlainText(PluginInstance) (string, error)
	Markup(PluginInstance) (Markup, error)
}

type Disposable interface {
	Dispose()
}

type Plugin interface {
	Name() string
	BindScorekeeper(*Scorekeeper)
}

type ScorekeeperTopic struct {
	Plugin Plugin
	Name string
	DisplayName string
	Description string
	Fetcher ScoreValueFetcher
}

type Engine interface {
	Disposable
	Logger() *log.Logger
	RandSource() rand.Source
	Scorekeeper() *Scorekeeper
	DefaultPort() Port
	Spawn(Spawnee) error
	Launch(PluginInstance) error
	SpawneeStatuses() ([]SpawneeStatus, error)
	PluginInstances() []PluginInstance
}

type InputFactory interface {
	Plugin
	New(engine Engine, config *ConfigElement) (Input, error)
}

type InputFactoryRegistry interface {
	RegisterInputFactory(factory InputFactory) error
	LookupInputFactory(name string) InputFactory
}

type OutputFactory interface {
	Plugin
	New(engine Engine, config *ConfigElement) (Output, error)
}

type OutputFactoryRegistry interface {
	RegisterOutputFactory(factory OutputFactory) error
	LookupOutputFactory(name string) OutputFactory
}

type PluginRegistry interface {
	Plugins() []Plugin
}

type Scoreboard interface {
	PluginInstance
}

type ScoreboardFactory interface {
	Plugin
	New(engine Engine, pluginRegistry PluginRegistry, config *ConfigElement) (Scoreboard, error)
}

type JournalChunk interface {
	GetReader() (io.Reader, error)
	GetNextChunk() JournalChunk
	TakeOwnership() bool
	Dispose() error
}

type JournalChunkListener func (JournalChunk) error

type Journal interface {
	Key() string
	Write(data []byte) error
	GetTailChunk() JournalChunk
	AddNewChunkListener(JournalChunkListener)
	AddFlushListener(JournalChunkListener)
	Flush(func (JournalChunk) error) error
	Dispose() error
}

type JournalGroup interface {
	GetJournal(key string) Journal
	GetJournalKeys() []string
	Dispose() error
}

type JournalGroupFactory interface {
	GetJournalGroup() JournalGroup
}

type RecordPacker interface {
	Pack(record FluentRecord) ([]byte, error)
}
