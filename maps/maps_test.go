package maps

import (
	"bytes"
	"io"
	"io/fs"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/google/pprof/profile"
	"github.com/stretchr/testify/require"
)

type fakefile struct {
	content io.Reader
}

func (f *fakefile) Stat() (fs.FileInfo, error) { return nil, nil }
func (f *fakefile) Read(b []byte) (int, error) { return f.content.Read(b) }
func (f *fakefile) Close() error               { return nil }

type fakefs struct {
	data map[string][]byte
}

func (f *fakefs) Open(name string) (fs.File, error) {
	return &fakefile{content: bytes.NewBuffer(f.data[name])}, nil
}

func testCache() *PidMappingFileCache {
	return &PidMappingFileCache{
		fs: &fakefs{map[string][]byte{
			"/proc/2043862/maps": []byte(`
00400000-00464000 r-xp 00000000 fd:01 2106801                            /main
00464000-004d4000 r--p 00064000 fd:01 2106801                            /main
004d4000-004d9000 rw-p 000d4000 fd:01 2106801                            /main
004d9000-0050b000 rw-p 00000000 00:00 0
c000000000-c004000000 rw-p 00000000 00:00 0
7f47d6714000-7f47d8a85000 rw-p 00000000 00:00 0
7f47d8a85000-7f47e8c05000 ---p 00000000 00:00 0
7f47e8c05000-7f47e8c06000 rw-p 00000000 00:00 0
7f47e8c06000-7f47faab5000 ---p 00000000 00:00 0
7f47faab5000-7f47faab6000 rw-p 00000000 00:00 0
7f47faab6000-7f47fce8b000 ---p 00000000 00:00 0
7f47fce8b000-7f47fce8c000 rw-p 00000000 00:00 0
7f47fce8c000-7f47fd305000 ---p 00000000 00:00 0
7f47fd305000-7f47fd306000 rw-p 00000000 00:00 0
7f47fd306000-7f47fd385000 ---p 00000000 00:00 0
7f47fd385000-7f47fd3e5000 rw-p 00000000 00:00 0
7ffc30d8b000-7ffc30dac000 rw-p 00000000 00:00 0                          [stack]
7ffc30dce000-7ffc30dd1000 r--p 00000000 00:00 0                          [vvar]
7ffc30dd1000-7ffc30dd3000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 r-xp 00000000 00:00 0                  [vsyscall]
			`),
		}},
		logger:     log.NewNopLogger(),
		cache:      map[uint32][]*profile.Mapping{},
		pidMapHash: map[uint32]uint64{},
	}
}

func TestPidMappingFileCache(t *testing.T) {
	c := testCache()
	mapping, err := c.MappingForPid(2043862)
	require.NoError(t, err)
	require.Equal(t, 3, len(mapping))
}

func TestMapping(t *testing.T) {
	m := &Mapping{
		fileCache:   testCache(),
		pidMappings: map[uint32][]*profile.Mapping{},
		pids:        []uint32{},
	}
	mapping, err := m.PidAddrMapping(2043862, 0x45e427)
	require.NoError(t, err)
	require.NotNil(t, mapping)

	resultMappings, _ := m.AllMappings()
	require.Equal(t, 3, len(resultMappings))
}
