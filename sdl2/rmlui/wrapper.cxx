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

static void _swig_gopanic(const char *p) {
  cgo_panic__rmlui(p);
}

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



void _wrap_delete_SystemInterface_rmlui(Rml::SystemInterface *_swig_go_0) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  delete arg1;
  
}


double _wrap_SystemInterface_GetElapsedTime_rmlui(Rml::SystemInterface *_swig_go_0) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  double result;
  double _swig_go_result;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  result = (double)(arg1)->GetElapsedTime();
  _swig_go_result = result; 
  return _swig_go_result;
}


intgo _wrap_SystemInterface_TranslateString_rmlui(Rml::SystemInterface *_swig_go_0, _gostring_* _swig_go_1, _gostring_ _swig_go_2) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  Rml::String *arg2 = 0 ;
  Rml::String *arg3 = 0 ;
  int result;
  intgo _swig_go_result;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  arg2 = *(Rml::String **)&_swig_go_1; 
  
  Rml::String arg3_str(_swig_go_2.p, _swig_go_2.n);
  arg3 = &arg3_str;
  
  
  result = (int)(arg1)->TranslateString(*arg2,(Rml::String const &)*arg3);
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_SystemInterface_JoinPath_rmlui(Rml::SystemInterface *_swig_go_0, _gostring_* _swig_go_1, _gostring_ _swig_go_2, _gostring_ _swig_go_3) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  Rml::String *arg2 = 0 ;
  Rml::String *arg3 = 0 ;
  Rml::String *arg4 = 0 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  arg2 = *(Rml::String **)&_swig_go_1; 
  
  Rml::String arg3_str(_swig_go_2.p, _swig_go_2.n);
  arg3 = &arg3_str;
  
  
  Rml::String arg4_str(_swig_go_3.p, _swig_go_3.n);
  arg4 = &arg4_str;
  
  
  (arg1)->JoinPath(*arg2,(Rml::String const &)*arg3,(Rml::String const &)*arg4);
  
}


bool _wrap_SystemInterface_LogMessage_rmlui(Rml::SystemInterface *_swig_go_0, Rml::Log::Type *_swig_go_1, _gostring_ _swig_go_2) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  Rml::Log::Type arg2 ;
  Rml::String *arg3 = 0 ;
  Rml::Log::Type *argp2 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  argp2 = (Rml::Log::Type *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null Log::Type");
  }
  arg2 = (Rml::Log::Type)*argp2;
  
  
  Rml::String arg3_str(_swig_go_2.p, _swig_go_2.n);
  arg3 = &arg3_str;
  
  
  result = (bool)(arg1)->LogMessage(arg2,(Rml::String const &)*arg3);
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_SystemInterface_SetMouseCursor_rmlui(Rml::SystemInterface *_swig_go_0, _gostring_ _swig_go_1) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  Rml::String *arg2 = 0 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  Rml::String arg2_str(_swig_go_1.p, _swig_go_1.n);
  arg2 = &arg2_str;
  
  
  (arg1)->SetMouseCursor((Rml::String const &)*arg2);
  
}


void _wrap_SystemInterface_SetClipboardText_rmlui(Rml::SystemInterface *_swig_go_0, _gostring_ _swig_go_1) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  Rml::String *arg2 = 0 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  Rml::String arg2_str(_swig_go_1.p, _swig_go_1.n);
  arg2 = &arg2_str;
  
  
  (arg1)->SetClipboardText((Rml::String const &)*arg2);
  
}


void _wrap_SystemInterface_GetClipboardText_rmlui(Rml::SystemInterface *_swig_go_0, _gostring_* _swig_go_1) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  Rml::String *arg2 = 0 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  arg2 = *(Rml::String **)&_swig_go_1; 
  
  (arg1)->GetClipboardText(*arg2);
  
}


void _wrap_SystemInterface_ActivateKeyboard_rmlui(Rml::SystemInterface *_swig_go_0, Rml::Vector2<float> *_swig_go_1, float _swig_go_2) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  SwigValueWrapper<Rml::Vector2<float>> arg2 ;
  float arg3 ;
  Rml::Vector2f *argp2 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  argp2 = (Rml::Vector2f *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null Rml::Vector2f");
  }
  arg2 = (Rml::Vector2f)*argp2;
  
  arg3 = (float)_swig_go_2; 
  
  (arg1)->ActivateKeyboard(arg2,arg3);
  
}


void _wrap_SystemInterface_DeactivateKeyboard_rmlui(Rml::SystemInterface *_swig_go_0) {
  Rml::SystemInterface *arg1 = (Rml::SystemInterface *) 0 ;
  
  arg1 = *(Rml::SystemInterface **)&_swig_go_0; 
  
  (arg1)->DeactivateKeyboard();
  
}



void _wrap_delete_RenderInterface_rmlui(Rml::RenderInterface *_swig_go_0) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  
  delete arg1;
  
}


void _wrap_RenderInterface_RenderGeometry_rmlui(Rml::RenderInterface *_swig_go_0, Rml::Vertex *_swig_go_1, intgo _swig_go_2, intgo *_swig_go_3, intgo _swig_go_4, Rml::TextureHandle *_swig_go_5, Rml::Vector2< float > *_swig_go_6) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::Vertex *arg2 = (Rml::Vertex *) 0 ;
  int arg3 ;
  int *arg4 = (int *) 0 ;
  int arg5 ;
  Rml::TextureHandle arg6 ;
  Rml::Vector2f *arg7 = 0 ;
  Rml::TextureHandle *argp6 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = *(Rml::Vertex **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  arg4 = *(int **)&_swig_go_3; 
  arg5 = (int)_swig_go_4; 
  
  argp6 = (Rml::TextureHandle *)_swig_go_5;
  if (argp6 == NULL) {
    _swig_gopanic("Attempt to dereference null TextureHandle");
  }
  arg6 = (Rml::TextureHandle)*argp6;
  
  arg7 = *(Rml::Vector2f **)&_swig_go_6; 
  
  (arg1)->RenderGeometry(arg2,arg3,arg4,arg5,arg6,(Rml::Vector2f const &)*arg7);
  
}


Rml::CompiledGeometryHandle *_wrap_RenderInterface_CompileGeometry_rmlui(Rml::RenderInterface *_swig_go_0, Rml::Vertex *_swig_go_1, intgo _swig_go_2, intgo *_swig_go_3, intgo _swig_go_4, Rml::TextureHandle *_swig_go_5) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::Vertex *arg2 = (Rml::Vertex *) 0 ;
  int arg3 ;
  int *arg4 = (int *) 0 ;
  int arg5 ;
  Rml::TextureHandle arg6 ;
  Rml::TextureHandle *argp6 ;
  Rml::CompiledGeometryHandle result;
  Rml::CompiledGeometryHandle *_swig_go_result;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = *(Rml::Vertex **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  arg4 = *(int **)&_swig_go_3; 
  arg5 = (int)_swig_go_4; 
  
  argp6 = (Rml::TextureHandle *)_swig_go_5;
  if (argp6 == NULL) {
    _swig_gopanic("Attempt to dereference null TextureHandle");
  }
  arg6 = (Rml::TextureHandle)*argp6;
  
  
  result = (arg1)->CompileGeometry(arg2,arg3,arg4,arg5,arg6);
  *(Rml::CompiledGeometryHandle **)&_swig_go_result = new Rml::CompiledGeometryHandle(result); 
  return _swig_go_result;
}


void _wrap_RenderInterface_RenderCompiledGeometry_rmlui(Rml::RenderInterface *_swig_go_0, Rml::CompiledGeometryHandle *_swig_go_1, Rml::Vector2< float > *_swig_go_2) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::CompiledGeometryHandle arg2 ;
  Rml::Vector2f *arg3 = 0 ;
  Rml::CompiledGeometryHandle *argp2 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  
  argp2 = (Rml::CompiledGeometryHandle *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null CompiledGeometryHandle");
  }
  arg2 = (Rml::CompiledGeometryHandle)*argp2;
  
  arg3 = *(Rml::Vector2f **)&_swig_go_2; 
  
  (arg1)->RenderCompiledGeometry(arg2,(Rml::Vector2f const &)*arg3);
  
}


void _wrap_RenderInterface_ReleaseCompiledGeometry_rmlui(Rml::RenderInterface *_swig_go_0, Rml::CompiledGeometryHandle *_swig_go_1) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::CompiledGeometryHandle arg2 ;
  Rml::CompiledGeometryHandle *argp2 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  
  argp2 = (Rml::CompiledGeometryHandle *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null CompiledGeometryHandle");
  }
  arg2 = (Rml::CompiledGeometryHandle)*argp2;
  
  
  (arg1)->ReleaseCompiledGeometry(arg2);
  
}


void _wrap_RenderInterface_EnableScissorRegion_rmlui(Rml::RenderInterface *_swig_go_0, bool _swig_go_1) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  bool arg2 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = (bool)_swig_go_1; 
  
  (arg1)->EnableScissorRegion(arg2);
  
}


void _wrap_RenderInterface_SetScissorRegion_rmlui(Rml::RenderInterface *_swig_go_0, intgo _swig_go_1, intgo _swig_go_2, intgo _swig_go_3, intgo _swig_go_4) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  int arg2 ;
  int arg3 ;
  int arg4 ;
  int arg5 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  arg4 = (int)_swig_go_3; 
  arg5 = (int)_swig_go_4; 
  
  (arg1)->SetScissorRegion(arg2,arg3,arg4,arg5);
  
}


bool _wrap_RenderInterface_LoadTexture_rmlui(Rml::RenderInterface *_swig_go_0, Rml::TextureHandle *_swig_go_1, Rml::Vector2< int > *_swig_go_2, _gostring_ _swig_go_3) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::TextureHandle *arg2 = 0 ;
  Rml::Vector2i *arg3 = 0 ;
  Rml::String *arg4 = 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = *(Rml::TextureHandle **)&_swig_go_1; 
  arg3 = *(Rml::Vector2i **)&_swig_go_2; 
  
  Rml::String arg4_str(_swig_go_3.p, _swig_go_3.n);
  arg4 = &arg4_str;
  
  
  result = (bool)(arg1)->LoadTexture(*arg2,*arg3,(Rml::String const &)*arg4);
  _swig_go_result = result; 
  return _swig_go_result;
}


bool _wrap_RenderInterface_GenerateTexture_rmlui(Rml::RenderInterface *_swig_go_0, Rml::TextureHandle *_swig_go_1, char *_swig_go_2, Rml::Vector2< int > *_swig_go_3) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::TextureHandle *arg2 = 0 ;
  Rml::byte *arg3 = (Rml::byte *) 0 ;
  Rml::Vector2i *arg4 = 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = *(Rml::TextureHandle **)&_swig_go_1; 
  arg3 = *(Rml::byte **)&_swig_go_2; 
  arg4 = *(Rml::Vector2i **)&_swig_go_3; 
  
  result = (bool)(arg1)->GenerateTexture(*arg2,(Rml::byte const *)arg3,(Rml::Vector2i const &)*arg4);
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_RenderInterface_ReleaseTexture_rmlui(Rml::RenderInterface *_swig_go_0, Rml::TextureHandle *_swig_go_1) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::TextureHandle arg2 ;
  Rml::TextureHandle *argp2 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  
  argp2 = (Rml::TextureHandle *)_swig_go_1;
  if (argp2 == NULL) {
    _swig_gopanic("Attempt to dereference null TextureHandle");
  }
  arg2 = (Rml::TextureHandle)*argp2;
  
  
  (arg1)->ReleaseTexture(arg2);
  
}


void _wrap_RenderInterface_SetTransform_rmlui(Rml::RenderInterface *_swig_go_0, Rml::Matrix4f *_swig_go_1) {
  Rml::RenderInterface *arg1 = (Rml::RenderInterface *) 0 ;
  Rml::Matrix4f *arg2 = (Rml::Matrix4f *) 0 ;
  
  arg1 = *(Rml::RenderInterface **)&_swig_go_0; 
  arg2 = *(Rml::Matrix4f **)&_swig_go_1; 
  
  (arg1)->SetTransform((Rml::Matrix4f const *)arg2);
  
}

#ifdef __cplusplus
}
#endif

