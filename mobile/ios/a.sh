
# ------------------------------------------------------------------------
# This script build the SDL libs
# ------------------------------------------------------------------------
SDK_DEVICE="-sdk iphoneos"
SDK_SIMULATOR="-sdk iphonesimulator"
OS="iOS"

CONF_DEBUG="-configuration Debug"
CONF_RELEASE="-configuration Release"
PREFIX="../../third-party"
SYMROOT="SYMROOT=../../$PREFIX/SDL/"

# ------------------------------------------------------------------------
# SDL
# ------------------------------------------------------------------------
cd $PREFIX/SDL/Xcode/SDL/
PROJ="-project SDL.xcodeproj"
SCHEME="Static Library-$OS"
#SCHEME="Framework-$OS"

xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 $PROJ $SDK_SIMULATOR  $CONF_RELEASE   -scheme "$SCHEME" build $SYMROOT $PLATFORM IPHONEOS_DEPLOYMENT_TARGET=12.0

