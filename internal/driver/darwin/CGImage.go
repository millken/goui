// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

CGImageRef nux_CGImageSourceCreateImageAtIndex(char* name){
    NSURL *url = [NSURL fileURLWithPath: [NSString stringWithUTF8String:name]];
    NSError *err;
    if ([url checkResourceIsReachableAndReturnError:&err] == NO){
        NSLog(@"%@", err);
    }
    CGImageSourceRef source = CGImageSourceCreateWithURL((__bridge CFURLRef)url, nil);
    CGImageRef ref = CGImageSourceCreateImageAtIndex(source, 0, nil);
	CFRelease(source);
	return ref;
}
*/
import "C"

import (
	"unsafe"
)

func CGImageSourceCreateImageAtIndex(filename string) CGImageRef {
	cpath := C.CString(filename)
	defer C.free(unsafe.Pointer(cpath))
	return CGImageRef(C.nux_CGImageSourceCreateImageAtIndex(cpath))
}

func CGImageRelease(image CGImageRef) {
	C.CGImageRelease(C.CGImageRef(image))
}

func CGImageGetSize(image CGImageRef) (width, height int32) {
	return int32(C.CGImageGetWidth(C.CGImageRef(image))), int32(C.CGImageGetHeight(C.CGImageRef(image)))
}
