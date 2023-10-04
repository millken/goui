#include "RmlUi/Core.h"


#include <typeinfo>
#include <stdexcept>
#include <string>
#include <vector>
#include <stdexcept>

#include "wrapper.h"

typedef long long intgo;
typedef unsigned long long uintgo;

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

#ifdef __cplusplus
}
#endif

