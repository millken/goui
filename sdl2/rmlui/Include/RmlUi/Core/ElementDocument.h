/*
 * This source file is part of RmlUi, the HTML/CSS Interface Middleware
 *
 * For the latest information, see http://github.com/mikke89/RmlUi
 *
 * Copyright (c) 2008-2010 CodePoint Ltd, Shift Technology Ltd
 * Copyright (c) 2019-2023 The RmlUi Team, and contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

#ifndef RMLUI_CORE_ELEMENTDOCUMENT_H
#define RMLUI_CORE_ELEMENTDOCUMENT_H

#include "Element.h"

namespace Rml {

class Context;
class Stream;
class DocumentHeader;
class ElementText;
class StyleSheet;
class StyleSheetContainer;

/**
     ModalFlag used for controlling the modal state of the document.
        None:  Remove modal state.
        Modal: Set modal state, other documents cannot receive focus.
        Keep:  Modal state unchanged.

    FocusFlag used for displaying the document.
        None:     No focus.
        Document: Focus the document.
        Keep:     Focus the element in the document which last had focus.
        Auto:     Focus the first tab element with the 'autofocus' attribute or else the document.
*/
enum class ModalFlag { None, Modal, Keep };
enum class FocusFlag { None, Document, Keep, Auto };

/**
    Represents a document in the dom tree.

    @author Lloyd Weehuizen
 */

class RMLUICORE_API ElementDocument : public Element {
public:


	/// Process given document header
	void ProcessHeader(const DocumentHeader* header);

	/// Returns the document's context.
	Context* GetContext();

	/// Sets the document's title.
	void SetTitle(const String& title);
	/// Returns the title of this document.
	const String& GetTitle() const;

	/// Returns the source address of this document.
	const String& GetSourceURL() const;

	/// Brings the document to the front of the document stack.
	void PullToFront();
	/// Sends the document to the back of the document stack.
	void PushToBack();

	/// Show the document.
	/// @param[in] modal_flag Flags controlling the modal state of the document, see the 'ModalFlag' description for details.
	/// @param[in] focus_flag Flags controlling the focus, see the 'FocusFlag' description for details.
	void Show(ModalFlag modal_flag = ModalFlag::None, FocusFlag focus_flag = FocusFlag::Auto);
	/// Hide the document.
	void Hide();
	/// Close the document.
	/// @note The destruction of the document is deferred until the next call to Context::Update().
	void Close();

};

} // namespace Rml
#endif
