package rmlui

/*
#define intgo swig_intgo

#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
typedef long long intgo;
typedef unsigned long long uintgo;
typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;

typedef _gostring_ swig_type_1;
typedef _gostring_ swig_type_2;
typedef _gostring_ swig_type_3;
typedef _gostring_ swig_type_4;
typedef _gostring_ swig_type_221;
typedef _gostring_ swig_type_230;
typedef _gostring_ swig_type_240;
extern uintptr_t _wrap_new_Vector2Int__SWIG_1_rmlui(swig_intgo arg1);
extern uintptr_t _wrap_new_Vector2Int__SWIG_2_rmlui(swig_intgo arg1, swig_intgo arg2);
extern uintptr_t _wrap_Context_LoadDocument__SWIG_0_rmlui(uintptr_t arg1, swig_type_4 arg2);
extern uintptr_t _wrap_new_Context_rmlui(swig_type_1 arg1);

extern void _wrap_ElementDocument_ProcessHeader_rmlui(uintptr_t arg1, uintptr_t arg2);
extern uintptr_t _wrap_ElementDocument_GetContext_rmlui(uintptr_t arg1);
extern void _wrap_ElementDocument_SetTitle_rmlui(uintptr_t arg1, swig_type_221 arg2);
extern swig_type_230 _wrap_ElementDocument_GetTitle_rmlui(uintptr_t arg1);
extern swig_type_240 _wrap_ElementDocument_GetSourceURL_rmlui(uintptr_t arg1);
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

#undef intgo
*/
import "C"
import "unsafe"

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
	C._wrap_ElementDocument_SetTitle_rmlui(C.uintptr_t(_swig_i_0), *(*C.swig_type_221)(unsafe.Pointer(&_swig_i_1)))
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
