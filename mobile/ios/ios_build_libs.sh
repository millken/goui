#!/bin/sh
OPTION=$1
OPTION2=$2

PREFIX=$HOME/release/Release_SDL_IOS/

# This script allows to compile libraries of SDL2, SDL_image, SDL_mixer and SDL_ttf from
# sources, for iOS and tvOS.
#
# Create 'PREFIX' directory, and place inside the SDL sources files: 
# 'PREFIX/SDL' contains sources of SDL2
# 'PREFIX/SDL_image' contains sources of SDL_image
# 'PREFIX/SDL_mixer' directory contains sources of SDL_mixer
# 'PREFIX/SDL_ttf' contains sources of SDL_ttf
#
# You can clone from git:
# > cd $PREFIX
# > git clone https://github.com/libsdl-org/SDL.git SDL
# > git clone https://github.com/libsdl-org/SDL_image.git SDL_image
# > git clone https://github.com/libsdl-org/SDL_mixer.git SDL_mixer
# > git clone https://github.com/libsdl-org/SDL_ttf.git SDL_ttf
# 
# Script needs execute permission ("chmod u+x ios_build_libs.sh") and should
# be placed in PREFIX directory
# 
# You need to install 'xcodebuild' console of Xcode
#
# Run the script from PREFIX directory: ./ios_build_libs.sh
#

syntax() 
{
    echo "$0: [--build | --clean | --state | --list] [--tvOS]"
    echo " build: re-compil with current sources"
    echo " clean: remove build directories"
    echo " state: see what's built"
    echo " list:  list xcode schemes & targets"
}

if [ "$OPTION" != "--clean" ] && [ "$OPTION" != "--build" ] && [ "$OPTION" != "--state" ] && [ "$OPTION" != "--list" ]
then
    syntax
    exit -1
fi

if [ "$OPTION2" != "--tvOS" ] && [ "$OPTION2" != "" ]
then
    syntax
    exit -1
fi 

SDL_DIR=$PREFIX/SDL


if [ "$OPTION" = "--list" ] 
then
    cd $PREFIX/SDL/Xcode/SDL/
    PROJ="-project SDL.xcodeproj"
    xcodebuild -list $PROJ

    cd $PREFIX/SDL_image/Xcode/
    PROJ="-project SDL_image.xcodeproj"
    xcodebuild -list $PROJ

    cd $PREFIX/SDL_mixer/Xcode/
    PROJ="-project SDL_mixer.xcodeproj"
    xcodebuild -list $PROJ

    cd $PREFIX/SDL_ttf/Xcode/
    PROJ="-project SDL_ttf.xcodeproj"
    xcodebuild -list $PROJ

    echo "======================= End ======================="
    exit 0
fi

if [ "$OPTION" = "--state" ] 
then
    MODES="Debug-iphonesimulator Release-iphonesimulator Debug-iphoneos Release-iphoneos Debug-appletvsimulator Release-appletvsimulator Debug-appletvos Release-appletvos"
    for mode in $MODES
    do
        echo "======================= List $mode ======================="
        ls -lah  $SDL_DIR/$mode/*
    done
    echo "======================= End ======================="
    exit 0
fi 


if [ "$OPTION" = "--clean" ] 
then
    if [ "$OPTION2" == "--tvOS" ] 
    then
        rm -Rf $SDL_DIR/Debug-appletv* $SDL_DIR/Release-appletv*
    else
        rm -Rf $SDL_DIR/Debug-iphone* $SDL_DIR/Release-iphone*
    fi
        
    rm -Rf $SDL_DIR/SDL.build $SDL_DIR/SDL_image.build $SDL_DIR/SDL_mixer.build $SDL_DIR/SDL_ttf.build

    cd $PREFIX/SDL/Xcode/SDL/
    PROJ="-project SDL.xcodeproj"
    xcodebuild clean $PROJ

    cd $PREFIX/SDL_image/Xcode/
    PROJ="-project SDL_image.xcodeproj"
    xcodebuild clean $PROJ

    cd $PREFIX/SDL_mixer/Xcode/
    PROJ="-project SDL_mixer.xcodeproj"
    xcodebuild clean $PROJ

    cd $PREFIX/SDL_ttf/Xcode/
    PROJ="-project SDL_ttf.xcodeproj"
    xcodebuild clean $PROJ

    exit 0
fi 

echo "======================= build for $OS ======================="

# ------------------------------------------------------------------------
# This script build the SDL libs
# ------------------------------------------------------------------------
SDK_DEVICE="-sdk iphoneos"
SDK_SIMULATOR="-sdk iphonesimulator"
if [ "$OPTION2" == "--tvOS" ] 
then
    SDK_DEVICE="-sdk appletvos"
    SDK_SIMULATOR="-sdk appletvsimulator"
fi

if [ "$OPTION2" == "--tvOS" ] 
then
    OS="tvOS"
else
    OS="iOS"
fi

CONF_DEBUG="-configuration Debug"
CONF_RELEASE="-configuration Release"
SYMROOT="SYMROOT=$PREFIX/SDL/"

# ------------------------------------------------------------------------
# SDL
# ------------------------------------------------------------------------
cd $PREFIX/SDL/Xcode/SDL/
PROJ="-project SDL.xcodeproj"
SCHEME="Static Library-$OS"
#SCHEME="Framework-$OS"

xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM


#SCHEME="Framework"
SCHEME="Static Library"

PLATFORM='PLATFORM=iOS'

# ------------------------------------------------------------------------
# SDL image
# ------------------------------------------------------------------------
cd $PREFIX/SDL_image/Xcode/
PROJ="-project SDL_image.xcodeproj"

xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM

# ------------------------------------------------------------------------
# SDL mixer
# ------------------------------------------------------------------------
cd $PREFIX/SDL_mixer/Xcode/
PROJ="-project SDL_mixer.xcodeproj"

xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM

# ------------------------------------------------------------------------
# SDL ttf
# ------------------------------------------------------------------------
cd $PREFIX/SDL_ttf/Xcode/
PROJ="-project SDL_ttf.xcodeproj"

xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch i386 -arch x86_64 $PROJ $SDK_SIMULATOR  $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_DEBUG   -scheme "$SCHEME" build $SYMROOT $PLATFORM
xcodebuild OTHER_CFLAGS="-fembed-bitcode" ONLY_ACTIVE_ARCH=NO -arch arm64 -arch armv7 $PROJ $SDK_DEVICE     $CONF_RELEASE -scheme "$SCHEME" build $SYMROOT $PLATFORM

