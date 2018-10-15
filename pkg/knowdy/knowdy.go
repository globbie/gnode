package knowdy

// #cgo CFLAGS: -I${SRCDIR}/knowdy/include
// #cgo CFLAGS: -I${SRCDIR}/knowdy/libs/gsl-parser/include
// #cgo LDFLAGS: ${SRCDIR}/knowdy/build/lib/libknowdy_static.a
// #cgo LDFLAGS: ${SRCDIR}/knowdy/build/lib/libglb-lib_static.a
// #cgo LDFLAGS: ${SRCDIR}/knowdy/build/lib/libgsl-parser_static.a
// #include <knd_shard.h>
// static void kndShard_del__(struct kndShard *shard)
// {
//     if (shard) {
//         kndShard_del(shard);
//     }
// }
import "C"
import (
	"errors"
)

type Shard struct {
	shard *C.struct_kndShard
}

func New(conf string) (*Shard, error) {
	var shard *C.struct_kndShard = nil

	errCode := C.kndShard_new((**C.struct_kndShard)(&shard), C.CString(conf), C.size_t(len(conf)))
	if errCode != C.int(0) {
		return nil, errors.New("could not create shard struct")
	}
	ret := &Shard{
		shard: shard,
	}
	return ret, nil
}

func (s *Shard) Del() error {
	C.kndShard_del__(s.shard)
	return nil
}

func (s *Shard) RunTask(task string) (string, error) {
	var output *C.char
	var outputLen C.size_t

	errCode := C.kndShard_run_task(s.shard, C.CString(task), C.size_t(len(task)), (**C.char)(&output), &outputLen)
	if errCode != C.int(0) {
		return "", errors.New("could not create shard struct")
	}
	return C.GoStringN(output, C.int(outputLen)), nil
}

