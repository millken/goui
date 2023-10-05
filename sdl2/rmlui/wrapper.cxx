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

#ifdef __cplusplus
}
#endif

