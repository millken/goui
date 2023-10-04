#include "rmlui.h"
#include "RmlUi/Core.h"

const char* RmlGetVersion()
{
    return strdup(Rml::GetVersion().c_str());
}
