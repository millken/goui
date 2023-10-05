package rmlui

/*
#define intgo swig_intgo

#include <stddef.h>
#include <stdint.h>

typedef long long intgo;
typedef unsigned long long uintgo;
typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;

typedef _gostring_ swig_type_1;
typedef _gostring_ swig_type_2;
typedef _gostring_ swig_type_3;
typedef _gostring_ swig_type_4;

extern uintptr_t _wrap_new_Vector2Int__SWIG_1_rmlui(swig_intgo arg1);
extern uintptr_t _wrap_new_Vector2Int__SWIG_2_rmlui(swig_intgo arg1, swig_intgo arg2);
extern uintptr_t _wrap_Context_LoadDocument__SWIG_0_rmlui(uintptr_t arg1, swig_type_4 arg2);
extern uintptr_t _wrap_new_Context_rmlui(swig_type_1 arg1);


#undef intgo
*/
import "C"

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
