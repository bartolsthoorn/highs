// This file provides support for constructing and solving models using HiGHS's
// "full" (low-level) API.

package highs

import (
	"fmt"
	"runtime"
	"unsafe"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <interfaces/highs_c_api.h>
import "C"

// A RawModel represents a HiGHS low-level model.
type RawModel struct {
	obj unsafe.Pointer
}

// NewRawModel allocates and returns an empty raw model.
func NewRawModel() *RawModel {
	model := &RawModel{}
	model.obj = C.Highs_create()
	runtime.SetFinalizer(model, func(m *RawModel) {
		C.Highs_destroy(m.obj)
	})
	return model
}

// SetBoolOption assigns a Boolean value to a named option.
func (m *RawModel) SetBoolOption(opt string, v bool) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	var val C.HighsInt
	if v {
		val = 1
	}

	// Set the option.
	switch C.Highs_setBoolOptionValue(m.obj, str, val) {
	case C.kHighsStatusOk:
		return nil
	case C.kHighsStatusError:
		return fmt.Errorf("SetBoolOption error")
	case C.kHighsStatusWarning:
		return fmt.Errorf("SetBoolOption warning")
	default:
		return fmt.Errorf("SetBoolOption unknown status")
	}
}

// SetIntOption assigns an integer value to a named option.
func (m *RawModel) SetIntOption(opt string, v int) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	val := C.HighsInt(v)

	// Set the option.
	switch C.Highs_setIntOptionValue(m.obj, str, val) {
	case C.kHighsStatusOk:
		return nil
	case C.kHighsStatusError:
		return fmt.Errorf("SetIntOption error")
	case C.kHighsStatusWarning:
		return fmt.Errorf("SetIntOption warning")
	default:
		return fmt.Errorf("SetIntOption unknown status")
	}
}

// SetFloat64Option assigns a floating-point value to a named option.
func (m *RawModel) SetFloat64Option(opt string, v float64) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	val := C.double(v)

	// Set the option.
	switch C.Highs_setDoubleOptionValue(m.obj, str, val) {
	case C.kHighsStatusOk:
		return nil
	case C.kHighsStatusError:
		return fmt.Errorf("SetFloat64Option error")
	case C.kHighsStatusWarning:
		return fmt.Errorf("SetFloat64Option warning")
	default:
		return fmt.Errorf("SetFloat64Option unknown status")
	}
}

// SetStringOption assigns a string value to a named option.
func (m *RawModel) SetStringOption(opt string, v string) error {
	// Convert arguments from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))
	val := C.CString(v)
	defer C.free(unsafe.Pointer(val))

	// Set the option.
	switch C.Highs_setStringOptionValue(m.obj, str, val) {
	case C.kHighsStatusOk:
		return nil
	case C.kHighsStatusError:
		return fmt.Errorf("SetStringOption error")
	case C.kHighsStatusWarning:
		return fmt.Errorf("SetStringOption warning")
	default:
		return fmt.Errorf("SetStringOption unknown status")
	}
}

// GetBoolOption returns the Boolean value of a named option.
func (m *RawModel) GetBoolOption(opt string) (bool, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.HighsInt
	switch C.Highs_getBoolOptionValue(m.obj, str, &val) {
	case C.kHighsStatusOk:
		var v bool
		if val != 0 {
			v = true
		}
		return v, nil
	case C.kHighsStatusError:
		return false, fmt.Errorf("GetBoolOption error")
	case C.kHighsStatusWarning:
		return false, fmt.Errorf("GetBoolOption warning")
	default:
		return false, fmt.Errorf("GetBoolOption unknown status")
	}
}

// GetIntOption returns the Integer value of a named option.
func (m *RawModel) GetIntOption(opt string) (int, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.HighsInt
	switch C.Highs_getIntOptionValue(m.obj, str, &val) {
	case C.kHighsStatusOk:
		return int(val), nil
	case C.kHighsStatusError:
		return 0, fmt.Errorf("GetIntOption error")
	case C.kHighsStatusWarning:
		return 0, fmt.Errorf("GetIntOption warning")
	default:
		return 0, fmt.Errorf("GetIntOption unknown status")
	}
}

// GetFloat64Option returns the floating-point value of a named option.
func (m *RawModel) GetFloat64Option(opt string) (float64, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.double
	switch C.Highs_getDoubleOptionValue(m.obj, str, &val) {
	case C.kHighsStatusOk:
		return float64(val), nil
	case C.kHighsStatusError:
		return 0.0, fmt.Errorf("GetFloat64Option error")
	case C.kHighsStatusWarning:
		return 0.0, fmt.Errorf("GetFloat64Option warning")
	default:
		return 0.0, fmt.Errorf("GetFloat64Option unknown status")
	}
}

// GetStringOption returns the string value of a named option.  Do not invoke
// this method in security-sensitive applications because it runs a risk of
// buffer overflow.
func (m *RawModel) GetStringOption(opt string) (string, error) {
	// Convert the option argument from Go to C.
	str := C.CString(opt)
	defer C.free(unsafe.Pointer(str))

	// The value could potentially be of any size.  Allocate "enough"
	// memory and hope for the best.
	val := (*C.char)(C.calloc(65536, 1))
	defer C.free(unsafe.Pointer(val))

	// Get the value.
	switch C.Highs_getStringOptionValue(m.obj, str, val) {
	case C.kHighsStatusOk:
		return C.GoString(val), nil
	case C.kHighsStatusError:
		return "", fmt.Errorf("GetStringOption error")
	case C.kHighsStatusWarning:
		return "", fmt.Errorf("GetStringOption warning")
	default:
		return "", fmt.Errorf("GetStringOption unknown status")
	}
}

// A RawSolution encapsulates all the values returned by various HiGHS solvers
// and provides methods to retrieve additional information.
type RawSolution struct {
	obj          unsafe.Pointer // Underlying highs opaque data type
	Status       ModelStatus    // Status of the LP solve
	ColumnPrimal []float64      // Primal column solution
	RowPrimal    []float64      // Primal row solution
	ColumnDual   []float64      // Dual column solution
	RowDual      []float64      // Dual row solution
	ColumnBasis  []BasisStatus // Basis status of each column
	RowBasis     []BasisStatus // Basis status of each row
	Objective    float64        // Objective value
}

// GetIntInfo returns the integer value of a named piece of information.
func (s *RawSolution) GetIntInfo(info string) (int, error) {
	// Convert the info argument from Go to C.
	str := C.CString(info)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.HighsInt
	switch C.Highs_getIntInfoValue(s.obj, str, &val) {
	case C.kHighsStatusOk:
		return int(val), nil
	case C.kHighsStatusError:
		return 0, fmt.Errorf("GetIntInfo error")
	case C.kHighsStatusWarning:
		return 0, fmt.Errorf("GetIntInfo warning")
	default:
		return 0, fmt.Errorf("GetIntInfo unknown status")
	}
}

// GetInt64Info returns the 64-bit integer value of a named piece of
// information.
func (s *RawSolution) GetInt64Info(info string) (int64, error) {
	// Convert the info argument from Go to C.
	str := C.CString(info)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.int64_t
	switch C.Highs_getInt64InfoValue(s.obj, str, &val) {
	case C.kHighsStatusOk:
		return int64(val), nil
	case C.kHighsStatusError:
		return 0, fmt.Errorf("GetInt64Info error")
	case C.kHighsStatusWarning:
		return 0, fmt.Errorf("GetInt64Info warning")
	default:
		return 0, fmt.Errorf("GetInt64Info unknown status")
	}
}

// GetFloat64Info returns the floating-point value of a named piece of
// information.
func (s *RawSolution) GetFloat64Info(info string) (float64, error) {
	// Convert the info argument from Go to C.
	str := C.CString(info)
	defer C.free(unsafe.Pointer(str))

	// Get the value.
	var val C.double
	switch C.Highs_getDoubleInfoValue(s.obj, str, &val) {
	case C.kHighsStatusOk:
		return float64(val), nil
	case C.kHighsStatusError:
		return 0.0, fmt.Errorf("GetFloat64Info error")
	case C.kHighsStatusWarning:
		return 0.0, fmt.Errorf("GetFloat64Info warning")
	default:
		return 0.0, fmt.Errorf("GetFloat64Info unknown status")
	}
}

// Solve solves a model.
func (m *RawModel) Solve() (*RawSolution, error) {
	// Solve the model.  We assume the user has already set up all the
	// required parameters.
	soln := &RawSolution{}
	success := C.Highs_run(m.obj)
	switch success {
	case C.kHighsStatusOk:
		// Success
	case C.kHighsStatusWarning:
		return soln, fmt.Errorf("model failed with a warning")
	case C.kHighsStatusError:
		return soln, fmt.Errorf("model failed with an error")
	default:
		return soln, fmt.Errorf("model failed with an unknown status")
	}

	// Extract various aspects of the solution as Go data.
	soln.Status = convertHighsModelStatus(C.Highs_getModelStatus(m.obj))
	nc := int(C.Highs_getNumCol(m.obj))
	nr := int(C.Highs_getNumRow(m.obj))
	colValue := make([]C.double, nc)
	colDual := make([]C.double, nc)
	rowValue := make([]C.double, nr)
	rowDual := make([]C.double, nr)
	switch C.Highs_getSolution(m.obj, &colValue[0], &colDual[0],
		&rowValue[0], &rowDual[0]) {
	case C.kHighsStatusOk:
		// Success
	case C.kHighsStatusWarning:
		return &RawSolution{}, fmt.Errorf("Highs_getSolution failed with a warning")
	case C.kHighsStatusError:
		return &RawSolution{}, fmt.Errorf("Highs_getSolution failed with an error")
	default:
		return &RawSolution{}, fmt.Errorf("Highs_getSolution failed with an unknown status")
	}
	soln.ColumnPrimal = convertSlice[float64, C.double](colValue)
	soln.RowPrimal = convertSlice[float64, C.double](rowValue)
	soln.ColumnDual = convertSlice[float64, C.double](colDual)
	soln.RowDual = convertSlice[float64, C.double](rowDual)
	return soln, nil
}
