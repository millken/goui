package rmlui

/*
#define intgo swig_intgo

#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
typedef long long intgo;
typedef unsigned long long uintgo;

typedef void *swig_voidp;

typedef struct { char *p; intgo n; } _gostring__;
typedef struct { void* array; intgo len; intgo cap; } _goslice__;


typedef _gostring__ swig_type_107;
typedef _gostring__ swig_type_108;
typedef _gostring__ swig_type_109;

typedef _gostring__ swig_type_100;
typedef _gostring__ swig_type_101;
typedef _gostring__ swig_type_102;
typedef _gostring__ swig_type_103;
typedef _gostring__ swig_type_104;
typedef _gostring__ swig_type_105;
typedef _gostring__ swig_type_106;
extern uintptr_t _wrap_new_Vector2Int__SWIG_1_rmlui(swig_intgo arg1);
extern uintptr_t _wrap_new_Vector2Int__SWIG_2_rmlui(swig_intgo arg1, swig_intgo arg2);
extern uintptr_t _wrap_new_Context_rmlui(swig_type_106 arg1);

extern void _wrap_ElementDocument_ProcessHeader_rmlui(uintptr_t arg1, uintptr_t arg2);
extern uintptr_t _wrap_ElementDocument_GetContext_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_SetTitle_rmlui(uintptr_t arg1, swig_type_107 arg2);
extern swig_type_108 _wrap_ElementDocument_GetTitle_rmlui(uintptr_t arg1);
extern swig_type_109 _wrap_ElementDocument_GetSourceURL_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_PullToFront_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_PushToBack_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_Show__SWIG_0_rmlui(uintptr_t arg1, swig_intgo arg2, swig_intgo arg3);
extern void _wrap_ElementDocument_Show__SWIG_1_rmlui(uintptr_t arg1, swig_intgo arg2);
extern void _wrap_ElementDocument_Show__SWIG_2_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_Hide_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_Close_rmlui(uintptr_t arg1);
extern void _wrap_delete_ElementDocument_rmlui(uintptr_t arg1);

extern swig_intgo _wrap_ModalFlag_None_rmlui(void);
extern swig_intgo _wrap_ModalFlag_Modal_rmlui(void);
extern swig_intgo _wrap_ModalFlag_Keep_rmlui(void);

extern swig_intgo _wrap_FocusFlag_None_rmlui(void);
extern swig_intgo _wrap_FocusFlag_Document_rmlui(void);
extern swig_intgo _wrap_FocusFlag_Keep_rmlui(void);
extern swig_intgo _wrap_FocusFlag_Auto_rmlui(void);

extern void _wrap_delete_SystemInterface_rmlui(uintptr_t arg1);
extern double _wrap_SystemInterface_GetElapsedTime_rmlui(uintptr_t arg1);
extern swig_intgo _wrap_SystemInterface_TranslateString_rmlui(uintptr_t arg1, swig_voidp arg2, swig_type_100 arg3);
extern void _wrap_SystemInterface_JoinPath_rmlui(uintptr_t arg1, swig_voidp arg2, swig_type_101 arg3, swig_type_102 arg4);
extern _Bool _wrap_SystemInterface_LogMessage_rmlui(uintptr_t arg1, uintptr_t arg2, swig_type_103 arg3);
extern void _wrap_SystemInterface_SetMouseCursor_rmlui(uintptr_t arg1, swig_type_104 arg2);
extern void _wrap_SystemInterface_SetClipboardText_rmlui(uintptr_t arg1, swig_type_105 arg2);
extern void _wrap_SystemInterface_GetClipboardText_rmlui(uintptr_t arg1, swig_voidp arg2);
extern void _wrap_SystemInterface_ActivateKeyboard_rmlui(uintptr_t arg1, uintptr_t arg2, float arg3);
extern void _wrap_SystemInterface_DeactivateKeyboard_rmlui(uintptr_t arg1);

#undef intgo
*/
import "C"
import (
	"sync"
	"unsafe"
)

func NewVector2Int_1(arg1 int) (_swig_ret Rml_Vector2_Sl_int_Sg_) {
	var swig_r Rml_Vector2_Sl_int_Sg_
	_swig_i_0 := arg1
	swig_r = SwigcptrRml_Vector2_Sl_int_Sg_(C._wrap_new_Vector2Int__SWIG_1_rmlui(C.swig_intgo(_swig_i_0)))
	return swig_r
}

func NewVector2Int_2(arg1 int, arg2 int) (_swig_ret Rml_Vector2_Sl_int_Sg_) {
	var swig_r Rml_Vector2_Sl_int_Sg_
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	swig_r = SwigcptrRml_Vector2_Sl_int_Sg_(C._wrap_new_Vector2Int__SWIG_2_rmlui(C.swig_intgo(_swig_i_0), C.swig_intgo(_swig_i_1)))
	return swig_r
}

type RmlModalFlag int

func _swig_getModalFlag_None() (_swig_ret RmlModalFlag) {
	var swig_r RmlModalFlag
	swig_r = (RmlModalFlag)(C._wrap_ModalFlag_None_rmlui())
	return swig_r
}

var ModalFlag_None RmlModalFlag = _swig_getModalFlag_None()

func _swig_getModalFlag_Modal() (_swig_ret RmlModalFlag) {
	var swig_r RmlModalFlag
	swig_r = (RmlModalFlag)(C._wrap_ModalFlag_Modal_rmlui())
	return swig_r
}

var ModalFlag_Modal RmlModalFlag = _swig_getModalFlag_Modal()

func _swig_getModalFlag_Keep() (_swig_ret RmlModalFlag) {
	var swig_r RmlModalFlag
	swig_r = (RmlModalFlag)(C._wrap_ModalFlag_Keep_rmlui())
	return swig_r
}

var ModalFlag_Keep RmlModalFlag = _swig_getModalFlag_Keep()

type RmlFocusFlag int

func _swig_getFocusFlag_None() (_swig_ret RmlFocusFlag) {
	var swig_r RmlFocusFlag
	swig_r = (RmlFocusFlag)(C._wrap_FocusFlag_None_rmlui())
	return swig_r
}

var FocusFlag_None RmlFocusFlag = _swig_getFocusFlag_None()

func _swig_getFocusFlag_Document() (_swig_ret RmlFocusFlag) {
	var swig_r RmlFocusFlag
	swig_r = (RmlFocusFlag)(C._wrap_FocusFlag_Document_rmlui())
	return swig_r
}

var FocusFlag_Document RmlFocusFlag = _swig_getFocusFlag_Document()

func _swig_getFocusFlag_Keep() (_swig_ret RmlFocusFlag) {
	var swig_r RmlFocusFlag
	swig_r = (RmlFocusFlag)(C._wrap_FocusFlag_Keep_rmlui())
	return swig_r
}

var FocusFlag_Keep RmlFocusFlag = _swig_getFocusFlag_Keep()

func _swig_getFocusFlag_Auto() (_swig_ret RmlFocusFlag) {
	var swig_r RmlFocusFlag
	swig_r = (RmlFocusFlag)(C._wrap_FocusFlag_Auto_rmlui())
	return swig_r
}

var FocusFlag_Auto RmlFocusFlag = _swig_getFocusFlag_Auto()

type RmlcptrElementDocument uintptr

func (p RmlcptrElementDocument) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p RmlcptrElementDocument) SwigIsElementDocument() {
}

func (arg1 RmlcptrElementDocument) GetContext() (_swig_ret Context) {
	var swig_r Context
	_swig_i_0 := arg1
	swig_r = (Context)(SwigcptrContext(C._wrap_ElementDocument_GetContext_rmlui(C.uintptr_t(_swig_i_0))))
	return swig_r
}

func (arg1 RmlcptrElementDocument) SetTitle(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_ElementDocument_SetTitle_rmlui(C.uintptr_t(_swig_i_0), *(*C.swig_type_107)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 RmlcptrElementDocument) GetTitle() (_swig_ret string) {
	var swig_r string
	_swig_i_0 := arg1
	swig_r_p := C._wrap_ElementDocument_GetTitle_rmlui(C.uintptr_t(_swig_i_0))
	swig_r = *(*string)(unsafe.Pointer(&swig_r_p))
	var swig_r_1 string
	swig_r_1 = swigCopyString(swig_r)
	return swig_r_1
}

func (arg1 RmlcptrElementDocument) GetSourceURL() (_swig_ret string) {
	var swig_r string
	_swig_i_0 := arg1
	swig_r_p := C._wrap_ElementDocument_GetSourceURL_rmlui(C.uintptr_t(_swig_i_0))
	swig_r = *(*string)(unsafe.Pointer(&swig_r_p))
	var swig_r_1 string
	swig_r_1 = swigCopyString(swig_r)
	return swig_r_1
}

func (arg1 RmlcptrElementDocument) PullToFront() {
	_swig_i_0 := arg1
	C._wrap_ElementDocument_PullToFront_rmlui(C.uintptr_t(_swig_i_0))
}

func (arg1 RmlcptrElementDocument) PushToBack() {
	_swig_i_0 := arg1
	C._wrap_ElementDocument_PushToBack_rmlui(C.uintptr_t(_swig_i_0))
}

func (arg1 RmlcptrElementDocument) Show__SWIG_0(arg2 RmlModalFlag, arg3 RmlFocusFlag) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	_swig_i_2 := arg3
	C._wrap_ElementDocument_Show__SWIG_0_rmlui(C.uintptr_t(_swig_i_0), C.swig_intgo(_swig_i_1), C.swig_intgo(_swig_i_2))
}

func (arg1 RmlcptrElementDocument) Show__SWIG_1(arg2 RmlModalFlag) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_ElementDocument_Show__SWIG_1_rmlui(C.uintptr_t(_swig_i_0), C.swig_intgo(_swig_i_1))
}

func (arg1 RmlcptrElementDocument) Show__SWIG_2() {
	_swig_i_0 := arg1
	C._wrap_ElementDocument_Show__SWIG_2_rmlui(C.uintptr_t(_swig_i_0))
}

func (p RmlcptrElementDocument) Show(a ...interface{}) {
	argc := len(a)
	if argc == 0 {
		p.Show__SWIG_2()
		return
	}
	if argc == 1 {
		p.Show__SWIG_1(a[0].(RmlModalFlag))
		return
	}
	if argc == 2 {
		p.Show__SWIG_0(a[0].(RmlModalFlag), a[1].(RmlFocusFlag))
		return
	}
	panic("No match for overloaded function call")
}

func (arg1 RmlcptrElementDocument) Hide() {
	_swig_i_0 := arg1
	C._wrap_ElementDocument_Hide_rmlui(C.uintptr_t(_swig_i_0))
}

func (arg1 RmlcptrElementDocument) Close() {
	_swig_i_0 := arg1
	C._wrap_ElementDocument_Close_rmlui(C.uintptr_t(_swig_i_0))
}

func DeleteElementDocument0(arg1 Rml_ElementDocument) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_ElementDocument_rmlui(C.uintptr_t(_swig_i_0))
}

type Rml_ElementDocument interface {
	Swigcptr() uintptr
	SwigIsElementDocument()
	GetContext() (_swig_ret Context)
	SetTitle(arg2 string)
	GetTitle() (_swig_ret string)
	GetSourceURL() (_swig_ret string)
	PullToFront()
	PushToBack()
	Show(a ...interface{})
	Hide()
	Close()
}

type SwigcptrSystemInterface uintptr

func (p SwigcptrSystemInterface) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrSystemInterface) SwigIsSystemInterface() {
}

func DeleteSystemInterface(arg1 SystemInterface) {
	_swig_i_0 := getSwigcptr(arg1)
	C._wrap_delete_SystemInterface_rmlui(C.uintptr_t(_swig_i_0))
}

func (arg1 SwigcptrSystemInterface) GetElapsedTime() (_swig_ret float64) {
	var swig_r float64
	_swig_i_0 := arg1
	swig_r = (float64)(C._wrap_SystemInterface_GetElapsedTime_rmlui(C.uintptr_t(_swig_i_0)))
	return swig_r
}

func (arg1 SwigcptrSystemInterface) TranslateString(arg2 *string, arg3 string) (_swig_ret int) {
	var swig_r int
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	_swig_i_2 := arg3
	swig_r = (int)(C._wrap_SystemInterface_TranslateString_rmlui(C.uintptr_t(_swig_i_0), C.swig_voidp(_swig_i_1), *(*C.swig_type_100)(unsafe.Pointer(&_swig_i_2))))
	if Swig_escape_always_false {
		Swig_escape_val = arg3
	}
	return swig_r
}

func (arg1 SwigcptrSystemInterface) JoinPath(arg2 *string, arg3 string, arg4 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	_swig_i_2 := arg3
	_swig_i_3 := arg4
	C._wrap_SystemInterface_JoinPath_rmlui(C.uintptr_t(_swig_i_0), C.swig_voidp(_swig_i_1), *(*C.swig_type_101)(unsafe.Pointer(&_swig_i_2)), *(*C.swig_type_102)(unsafe.Pointer(&_swig_i_3)))
	if Swig_escape_always_false {
		Swig_escape_val = arg3
	}
	if Swig_escape_always_false {
		Swig_escape_val = arg4
	}
}

func (arg1 SwigcptrSystemInterface) LogMessage(arg2 Log_Type, arg3 string) (_swig_ret bool) {
	var swig_r bool
	_swig_i_0 := arg1
	_swig_i_1 := getSwigcptr(arg2)
	_swig_i_2 := arg3
	swig_r = (bool)(C._wrap_SystemInterface_LogMessage_rmlui(C.uintptr_t(_swig_i_0), C.uintptr_t(_swig_i_1), *(*C.swig_type_103)(unsafe.Pointer(&_swig_i_2))))
	if Swig_escape_always_false {
		Swig_escape_val = arg3
	}
	return swig_r
}

func (arg1 SwigcptrSystemInterface) SetMouseCursor(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_SystemInterface_SetMouseCursor_rmlui(C.uintptr_t(_swig_i_0), *(*C.swig_type_104)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrSystemInterface) SetClipboardText(arg2 string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_SystemInterface_SetClipboardText_rmlui(C.uintptr_t(_swig_i_0), *(*C.swig_type_105)(unsafe.Pointer(&_swig_i_1)))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrSystemInterface) GetClipboardText(arg2 *string) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	C._wrap_SystemInterface_GetClipboardText_rmlui(C.uintptr_t(_swig_i_0), C.swig_voidp(_swig_i_1))
}

func (arg1 SwigcptrSystemInterface) ActivateKeyboard(arg2 Rml_Vector2_Sl_float_Sg_, arg3 float32) {
	_swig_i_0 := arg1
	_swig_i_1 := getSwigcptr(arg2)
	_swig_i_2 := arg3
	C._wrap_SystemInterface_ActivateKeyboard_rmlui(C.uintptr_t(_swig_i_0), C.uintptr_t(_swig_i_1), C.float(_swig_i_2))
}

func (arg1 SwigcptrSystemInterface) DeactivateKeyboard() {
	_swig_i_0 := arg1
	C._wrap_SystemInterface_DeactivateKeyboard_rmlui(C.uintptr_t(_swig_i_0))
}

type SystemInterface interface {
	Swigcptr() uintptr
	SwigIsSystemInterface()
	GetElapsedTime() (_swig_ret float64)
	TranslateString(arg2 *string, arg3 string) (_swig_ret int)
	JoinPath(arg2 *string, arg3 string, arg4 string)
	LogMessage(arg2 Log_Type, arg3 string) (_swig_ret bool)
	SetMouseCursor(arg2 string)
	SetClipboardText(arg2 string)
	GetClipboardText(arg2 *string)
	ActivateKeyboard(arg2 Rml_Vector2_Sl_float_Sg_, arg3 float32)
	DeactivateKeyboard()
}
type _ sync.Mutex

//export cgo_panic__rmlui
func cgo_panic__rmlui(p *byte) {
	s := (*[1024]byte)(unsafe.Pointer(p))[:]
	for i, b := range s {
		if b == 0 {
			panic(string(s[:i]))
		}
	}
	panic(string(s))
}

type SwigcptrLog_Type uintptr
type Log_Type interface {
	Swigcptr() uintptr
}

func (p SwigcptrLog_Type) Swigcptr() uintptr {
	return uintptr(p)
}
