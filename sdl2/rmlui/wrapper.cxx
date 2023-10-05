#include "RmlUi/Core.h"


#include <typeinfo>
#include <stdexcept>
#include <string>
#include <vector>
#include <stdexcept>

#include "wrapper.h"

typedef long long intgo;
typedef unsigned long long uintgo;
typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;



static _gostring_ Swig_AllocateString(const char *p, size_t l) {
  _gostring_ ret;
  ret.p = (char*)malloc(l);
  memcpy(ret.p, p, l);
  ret.n = l;
  return ret;
}
#ifdef __cplusplus
extern "C" {
#endif

Rml::Vector2< int > *_wrap_new_Vector2Int__SWIG_1_rmlui(intgo _swig_go_0) {
  int arg1 ;
  Rml::Vector2< int > *result = 0 ;
  Rml::Vector2< int > *_swig_go_result;
  
  arg1 = (int)_swig_go_0; 
  
  result = (Rml::Vector2< int > *)new Rml::Vector2< int >(arg1);
  *(Rml::Vector2< int > **)&_swig_go_result = (Rml::Vector2< int > *)result; 
  return _swig_go_result;
}


Rml::Vector2< int > *_wrap_new_Vector2Int__SWIG_2_rmlui(intgo _swig_go_0, intgo _swig_go_1) {
  int arg1 ;
  int arg2 ;
  Rml::Vector2< int > *result = 0 ;
  Rml::Vector2< int > *_swig_go_result;
  
  arg1 = (int)_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  result = (Rml::Vector2< int > *)new Rml::Vector2< int >(arg1,arg2);
  *(Rml::Vector2< int > **)&_swig_go_result = (Rml::Vector2< int > *)result; 
  return _swig_go_result;
}

Rml::ElementDocument *_wrap_Context_LoadDocument__SWIG_0_rmlui(Rml::Context *_swig_go_0, _gostring_ _swig_go_1) {
  Rml::Context *arg1 = (Rml::Context *) 0 ;
  Rml::String *arg2 = 0 ;
  Rml::ElementDocument *result = 0 ;
  Rml::ElementDocument *_swig_go_result;
  
  arg1 = *(Rml::Context **)&_swig_go_0; 
  
  Rml::String arg2_str(_swig_go_1.p, _swig_go_1.n);
  arg2 = &arg2_str;
  
  
  result = (Rml::ElementDocument *)(arg1)->LoadDocument((Rml::String const &)*arg2);
  *(Rml::ElementDocument **)&_swig_go_result = (Rml::ElementDocument *)result; 
  return _swig_go_result;
}


void _wrap_ElementDocument_ProcessHeader_rmlui(Rml::ElementDocument *_swig_go_0, Rml::DocumentHeader *_swig_go_1) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::DocumentHeader *arg2 = (Rml::DocumentHeader *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  arg2 = *(Rml::DocumentHeader **)&_swig_go_1; 
  
  (arg1)->ProcessHeader((Rml::DocumentHeader const *)arg2);
  
}


Rml::Context *_wrap_ElementDocument_GetContext_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::Context *result = 0 ;
  Rml::Context *_swig_go_result;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  result = (Rml::Context *)(arg1)->GetContext();
  *(Rml::Context **)&_swig_go_result = (Rml::Context *)result; 
  return _swig_go_result;
}


void _wrap_ElementDocument_SetTitle_rmlui(Rml::ElementDocument *_swig_go_0, _gostring_ _swig_go_1) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::String *arg2 = 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  Rml::String arg2_str(_swig_go_1.p, _swig_go_1.n);
  arg2 = &arg2_str;
  
  
  (arg1)->SetTitle((Rml::String const &)*arg2);
  
}


_gostring_ _wrap_ElementDocument_GetTitle_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::String *result = 0 ;
  _gostring_ _swig_go_result;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  result = (Rml::String *) &((Rml::ElementDocument const *)arg1)->GetTitle();
  _swig_go_result = Swig_AllocateString((*result).data(), (*result).length()); 
  return _swig_go_result;
}


_gostring_ _wrap_ElementDocument_GetSourceURL_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::String *result = 0 ;
  _gostring_ _swig_go_result;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  result = (Rml::String *) &((Rml::ElementDocument const *)arg1)->GetSourceURL();
  _swig_go_result = Swig_AllocateString((*result).data(), (*result).length()); 
  return _swig_go_result;
}


void _wrap_ElementDocument_PullToFront_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  (arg1)->PullToFront();
  
}


void _wrap_ElementDocument_PushToBack_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  (arg1)->PushToBack();
  
}


void _wrap_ElementDocument_Show__SWIG_0_rmlui(Rml::ElementDocument *_swig_go_0, intgo _swig_go_1, intgo _swig_go_2) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::ModalFlag arg2 ;
  Rml::FocusFlag arg3 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  arg2 = (Rml::ModalFlag)_swig_go_1; 
  arg3 = (Rml::FocusFlag)_swig_go_2; 
  
  (arg1)->Show(arg2,arg3);
  
}


void _wrap_ElementDocument_Show__SWIG_1_rmlui(Rml::ElementDocument *_swig_go_0, intgo _swig_go_1) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  Rml::ModalFlag arg2 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  arg2 = (Rml::ModalFlag)_swig_go_1; 
  
  (arg1)->Show(arg2);
  
}


void _wrap_ElementDocument_Show__SWIG_2_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  (arg1)->Show();
  
}


void _wrap_ElementDocument_Hide_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  (arg1)->Hide();
  
}


void _wrap_ElementDocument_Close_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  (arg1)->Close();
  
}


void _wrap_delete_ElementDocument_rmlui(Rml::ElementDocument *_swig_go_0) {
  Rml::ElementDocument *arg1 = (Rml::ElementDocument *) 0 ;
  
  arg1 = *(Rml::ElementDocument **)&_swig_go_0; 
  
  delete arg1;
  
}


intgo _wrap_ModalFlag_None_rmlui() {
  Rml::ModalFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::ModalFlag::None;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}


intgo _wrap_ModalFlag_Modal_rmlui() {
  Rml::ModalFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::ModalFlag::Modal;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}


intgo _wrap_ModalFlag_Keep_rmlui() {
  Rml::ModalFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::ModalFlag::Keep;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}


intgo _wrap_FocusFlag_None_rmlui() {
  Rml::FocusFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::FocusFlag::None;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}


intgo _wrap_FocusFlag_Document_rmlui() {
  Rml::FocusFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::FocusFlag::Document;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}


intgo _wrap_FocusFlag_Keep_rmlui() {
  Rml::FocusFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::FocusFlag::Keep;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}


intgo _wrap_FocusFlag_Auto_rmlui() {
  Rml::FocusFlag result;
  intgo _swig_go_result;
  
  
  result = Rml::FocusFlag::Auto;
  
  _swig_go_result = (intgo)result; 
  return _swig_go_result;
}

#ifdef __cplusplus
}
#endif

