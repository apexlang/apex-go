package errors

import (
	json "encoding/json"
	"os"

	"github.com/tetratelabs/tinymem"
)

func Return(errs ...error) (ptrSize uint64) {
	cerrs := Convert(errs...)
	jsonBytes, _ := json.Marshal(cerrs)
	jsonString := string(jsonBytes)
	ptr, size := tinymem.StringToPtr(jsonString)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

func Write(errs ...error) {
	cerrs := Convert(errs...)
	jsonBytes, _ := json.Marshal(cerrs)
	os.Stderr.Write(jsonBytes)
	os.Exit(1)
}

func Convert(errs ...error) Errors {
	e := make(Errors, len(errs))
	for i, err := range errs {
		if ee, ok := err.(*Error); ok {
			e[i] = ee
		} else {
			e[i] = &Error{
				Message: err.Error(),
			}
		}
	}
	return e
}
