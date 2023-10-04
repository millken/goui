
#ifndef RMLUI_DEFINED
#define RMLUI_DEFINED

#include <stdint.h>
#include <stddef.h>
#include <string.h>
#ifndef __cplusplus
#include <stdbool.h>
#endif // __cplusplus
#if !defined(RMLUI_API)
#define RMLUI_API
#endif
#define RMLUI_STATIC_LIB
#ifdef __cplusplus
extern "C" {
#endif
RMLUI_API const char* RmlGetVersion();

#ifdef __cplusplus
}
#endif

#endif // RMLUI_DEFINED