// Code auto-generated to provide html and svg element types for 
// trees.
// Documentation source: "HTML element reference" by Mozilla Contributors, 
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element, licensed under CC-BY-SA 2.5.

package ntrees

// SvgAnchor provides Node representation for the element "a" in XML SVG DOM 
// The <a> SVG element creates a hyperlink to other web pages, files, locations within the same page, email addresses, or any other URL.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/a
func SvgAnchor(id string, renders ...Render) (*Node, error) {
	return Element("a", id, renders...)
}


// MustSvgAnchor provides Node representation for the element "a" in XML SVG DOM 
// The <a> SVG element creates a hyperlink to other web pages, files, locations within the same page, email addresses, or any other URL.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/a
func MustSvgAnchor(id string, renders ...Render) *Node {
	var node, err = Element("a", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgAnimate provides Node representation for the element "animate" in XML SVG DOM 
// This element implements the SVGAnimateElement interface.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animate
func SvgAnimate(id string, renders ...Render) (*Node, error) {
	return Element("animate", id, renders...)
}


// MustSvgAnimate provides Node representation for the element "animate" in XML SVG DOM 
// This element implements the SVGAnimateElement interface.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animate
func MustSvgAnimate(id string, renders ...Render) *Node {
	var node, err = Element("animate", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgAnimateMotion provides Node representation for the element "animateMotion" in XML SVG DOM 
// The <animateMotion> element causes a referenced element to move along a motion path.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateMotion
func SvgAnimateMotion(id string, renders ...Render) (*Node, error) {
	return Element("animateMotion", id, renders...)
}


// MustSvgAnimateMotion provides Node representation for the element "animateMotion" in XML SVG DOM 
// The <animateMotion> element causes a referenced element to move along a motion path.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateMotion
func MustSvgAnimateMotion(id string, renders ...Render) *Node {
	var node, err = Element("animateMotion", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgAnimateTransform provides Node representation for the element "animateTransform" in XML SVG DOM 
// The animateTransform element animates a transformation attribute on a target element, thereby allowing animations to control translation, scaling, rotation and/or skewing.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateTransform
func SvgAnimateTransform(id string, renders ...Render) (*Node, error) {
	return Element("animateTransform", id, renders...)
}


// MustSvgAnimateTransform provides Node representation for the element "animateTransform" in XML SVG DOM 
// The animateTransform element animates a transformation attribute on a target element, thereby allowing animations to control translation, scaling, rotation and/or skewing.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/animateTransform
func MustSvgAnimateTransform(id string, renders ...Render) *Node {
	var node, err = Element("animateTransform", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgCircle provides Node representation for the element "circle" in XML SVG DOM 
// The <circle> SVG element is an SVG basic shape, used to create circles based on a center point and a radius.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/circle
func SvgCircle(id string, renders ...Render) (*Node, error) {
	return Element("circle", id, renders...)
}


// MustSvgCircle provides Node representation for the element "circle" in XML SVG DOM 
// The <circle> SVG element is an SVG basic shape, used to create circles based on a center point and a radius.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/circle
func MustSvgCircle(id string, renders ...Render) *Node {
	var node, err = Element("circle", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgClipPath provides Node representation for the element "clipPath" in XML SVG DOM 
// The <clipPath> SVG element defines a clipping path. A clipping path is used/referenced using the clip-path property.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/clipPath
func SvgClipPath(id string, renders ...Render) (*Node, error) {
	return Element("clipPath", id, renders...)
}


// MustSvgClipPath provides Node representation for the element "clipPath" in XML SVG DOM 
// The <clipPath> SVG element defines a clipping path. A clipping path is used/referenced using the clip-path property.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/clipPath
func MustSvgClipPath(id string, renders ...Render) *Node {
	var node, err = Element("clipPath", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgColorProfile provides Node representation for the element "color-profile" in XML SVG DOM 
// The <color-profile> element allows describing the color profile used for the image.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/color-profile
func SvgColorProfile(id string, renders ...Render) (*Node, error) {
	return Element("color-profile", id, renders...)
}


// MustSvgColorProfile provides Node representation for the element "color-profile" in XML SVG DOM 
// The <color-profile> element allows describing the color profile used for the image.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/color-profile
func MustSvgColorProfile(id string, renders ...Render) *Node {
	var node, err = Element("color-profile", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgDefs provides Node representation for the element "defs" in XML SVG DOM 
// The <defs> element is used to store graphical objects that will be used at a later time. Objects created inside a <defs> element are not rendered directly. To display them you have to reference them (with a <use> element for example).
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/defs
func SvgDefs(id string, renders ...Render) (*Node, error) {
	return Element("defs", id, renders...)
}


// MustSvgDefs provides Node representation for the element "defs" in XML SVG DOM 
// The <defs> element is used to store graphical objects that will be used at a later time. Objects created inside a <defs> element are not rendered directly. To display them you have to reference them (with a <use> element for example).
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/defs
func MustSvgDefs(id string, renders ...Render) *Node {
	var node, err = Element("defs", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgDesc provides Node representation for the element "desc" in XML SVG DOM 
// Each container element or graphics element in an SVG drawing can supply a description string using the <desc> element where the description is text-only.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/desc
func SvgDesc(id string, renders ...Render) (*Node, error) {
	return Element("desc", id, renders...)
}


// MustSvgDesc provides Node representation for the element "desc" in XML SVG DOM 
// Each container element or graphics element in an SVG drawing can supply a description string using the <desc> element where the description is text-only.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/desc
func MustSvgDesc(id string, renders ...Render) *Node {
	var node, err = Element("desc", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgDiscard provides Node representation for the element "discard" in XML SVG DOM 
// The <discard> SVG element allows authors to specify the time at which particular elements are to be discarded, thereby reducing the resources required by an SVG user agent. This is particularly useful to help SVG viewers conserve memory while displaying long-running documents.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/discard
func SvgDiscard(id string, renders ...Render) (*Node, error) {
	return Element("discard", id, renders...)
}


// MustSvgDiscard provides Node representation for the element "discard" in XML SVG DOM 
// The <discard> SVG element allows authors to specify the time at which particular elements are to be discarded, thereby reducing the resources required by an SVG user agent. This is particularly useful to help SVG viewers conserve memory while displaying long-running documents.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/discard
func MustSvgDiscard(id string, renders ...Render) *Node {
	var node, err = Element("discard", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgEllipse provides Node representation for the element "ellipse" in XML SVG DOM 
// The <ellipse> element is an SVG basic shape, used to create ellipses based on a center coordinate, and both their x and y radius.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/ellipse
func SvgEllipse(id string, renders ...Render) (*Node, error) {
	return Element("ellipse", id, renders...)
}


// MustSvgEllipse provides Node representation for the element "ellipse" in XML SVG DOM 
// The <ellipse> element is an SVG basic shape, used to create ellipses based on a center coordinate, and both their x and y radius.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/ellipse
func MustSvgEllipse(id string, renders ...Render) *Node {
	var node, err = Element("ellipse", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeBlend provides Node representation for the element "feBlend" in XML SVG DOM 
// The <feBlend> SVG filter primitive composes two objects together ruled by a certain blending mode. This is similar to what is known from image editing software when blending two layers. The mode is defined by the mode attribute.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feBlend
func SvgFeBlend(id string, renders ...Render) (*Node, error) {
	return Element("feBlend", id, renders...)
}


// MustSvgFeBlend provides Node representation for the element "feBlend" in XML SVG DOM 
// The <feBlend> SVG filter primitive composes two objects together ruled by a certain blending mode. This is similar to what is known from image editing software when blending two layers. The mode is defined by the mode attribute.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feBlend
func MustSvgFeBlend(id string, renders ...Render) *Node {
	var node, err = Element("feBlend", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeColorMatrix provides Node representation for the element "feColorMatrix" in XML SVG DOM 
// The <feColorMatrix> SVG filter element changes colors based on a transformation matrix. Every pixel's color value (represented by an [R,G,B,A] vector) is matrix multiplied to create a new color:
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feColorMatrix
func SvgFeColorMatrix(id string, renders ...Render) (*Node, error) {
	return Element("feColorMatrix", id, renders...)
}


// MustSvgFeColorMatrix provides Node representation for the element "feColorMatrix" in XML SVG DOM 
// The <feColorMatrix> SVG filter element changes colors based on a transformation matrix. Every pixel's color value (represented by an [R,G,B,A] vector) is matrix multiplied to create a new color:
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feColorMatrix
func MustSvgFeColorMatrix(id string, renders ...Render) *Node {
	var node, err = Element("feColorMatrix", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeComponentTransfer provides Node representation for the element "feComponentTransfer" in XML SVG DOM 
// Th <feComponentTransfer> SVG filter primitive performs color-component-wise remapping of data for each pixel. It allows operations like brightness adjustment, contrast adjustment, color balance or thresholding.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feComponentTransfer
func SvgFeComponentTransfer(id string, renders ...Render) (*Node, error) {
	return Element("feComponentTransfer", id, renders...)
}


// MustSvgFeComponentTransfer provides Node representation for the element "feComponentTransfer" in XML SVG DOM 
// Th <feComponentTransfer> SVG filter primitive performs color-component-wise remapping of data for each pixel. It allows operations like brightness adjustment, contrast adjustment, color balance or thresholding.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feComponentTransfer
func MustSvgFeComponentTransfer(id string, renders ...Render) *Node {
	var node, err = Element("feComponentTransfer", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeComposite provides Node representation for the element "feComposite" in XML SVG DOM 
// The <feComposite> SVG filter primitive performs the combination of two input images pixel-wise in image space using one of the Porter-Duff compositing operations: over, in, atop, out, xor, and lighter. Additionally, a component-wise arithmetic operation (with the result clamped between [0..1]) can be applied.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feComposite
func SvgFeComposite(id string, renders ...Render) (*Node, error) {
	return Element("feComposite", id, renders...)
}


// MustSvgFeComposite provides Node representation for the element "feComposite" in XML SVG DOM 
// The <feComposite> SVG filter primitive performs the combination of two input images pixel-wise in image space using one of the Porter-Duff compositing operations: over, in, atop, out, xor, and lighter. Additionally, a component-wise arithmetic operation (with the result clamped between [0..1]) can be applied.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feComposite
func MustSvgFeComposite(id string, renders ...Render) *Node {
	var node, err = Element("feComposite", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeConvolveMatrix provides Node representation for the element "feConvolveMatrix" in XML SVG DOM 
// The <feConvolveMatrix> SVG filter primitive applies a matrix convolution filter effect. A convolution combines pixels in the input image with neighboring pixels to produce a resulting image. A wide variety of imaging operations can be achieved through convolutions, including blurring, edge detection, sharpening, embossing and beveling.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feConvolveMatrix
func SvgFeConvolveMatrix(id string, renders ...Render) (*Node, error) {
	return Element("feConvolveMatrix", id, renders...)
}


// MustSvgFeConvolveMatrix provides Node representation for the element "feConvolveMatrix" in XML SVG DOM 
// The <feConvolveMatrix> SVG filter primitive applies a matrix convolution filter effect. A convolution combines pixels in the input image with neighboring pixels to produce a resulting image. A wide variety of imaging operations can be achieved through convolutions, including blurring, edge detection, sharpening, embossing and beveling.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feConvolveMatrix
func MustSvgFeConvolveMatrix(id string, renders ...Render) *Node {
	var node, err = Element("feConvolveMatrix", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeDiffuseLighting provides Node representation for the element "feDiffuseLighting" in XML SVG DOM 
// The <feDiffuseLighting> SVG filter primitive lights an image using the alpha channel as a bump map. The resulting image, which is an RGBA opaque image, depends on the light color, light position and surface geometry of the input bump map.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDiffuseLighting
func SvgFeDiffuseLighting(id string, renders ...Render) (*Node, error) {
	return Element("feDiffuseLighting", id, renders...)
}


// MustSvgFeDiffuseLighting provides Node representation for the element "feDiffuseLighting" in XML SVG DOM 
// The <feDiffuseLighting> SVG filter primitive lights an image using the alpha channel as a bump map. The resulting image, which is an RGBA opaque image, depends on the light color, light position and surface geometry of the input bump map.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDiffuseLighting
func MustSvgFeDiffuseLighting(id string, renders ...Render) *Node {
	var node, err = Element("feDiffuseLighting", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeDisplacementMap provides Node representation for the element "feDisplacementMap" in XML SVG DOM 
// The <feDisplacementMap> SVG filter primitive uses the pixel values from the image from in2 to spatially displace the image from in.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDisplacementMap
func SvgFeDisplacementMap(id string, renders ...Render) (*Node, error) {
	return Element("feDisplacementMap", id, renders...)
}


// MustSvgFeDisplacementMap provides Node representation for the element "feDisplacementMap" in XML SVG DOM 
// The <feDisplacementMap> SVG filter primitive uses the pixel values from the image from in2 to spatially displace the image from in.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDisplacementMap
func MustSvgFeDisplacementMap(id string, renders ...Render) *Node {
	var node, err = Element("feDisplacementMap", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeDistantLight provides Node representation for the element "feDistantLight" in XML SVG DOM 
// The <feDistantLight> filter primitive defines a distant light source that can be used within a lighting filter primitive: <feDiffuseLighting> or <feSpecularLighting>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDistantLight
func SvgFeDistantLight(id string, renders ...Render) (*Node, error) {
	return Element("feDistantLight", id, renders...)
}


// MustSvgFeDistantLight provides Node representation for the element "feDistantLight" in XML SVG DOM 
// The <feDistantLight> filter primitive defines a distant light source that can be used within a lighting filter primitive: <feDiffuseLighting> or <feSpecularLighting>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDistantLight
func MustSvgFeDistantLight(id string, renders ...Render) *Node {
	var node, err = Element("feDistantLight", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeDropShadow provides Node representation for the element "feDropShadow" in XML SVG DOM 
// The <feDropShadow> filter primitive creates a drop shadow of the input image. It is a shorthand filter, and is defined in terms of combinations of other filter primitives.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDropShadow
func SvgFeDropShadow(id string, renders ...Render) (*Node, error) {
	return Element("feDropShadow", id, renders...)
}


// MustSvgFeDropShadow provides Node representation for the element "feDropShadow" in XML SVG DOM 
// The <feDropShadow> filter primitive creates a drop shadow of the input image. It is a shorthand filter, and is defined in terms of combinations of other filter primitives.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feDropShadow
func MustSvgFeDropShadow(id string, renders ...Render) *Node {
	var node, err = Element("feDropShadow", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeFlood provides Node representation for the element "feFlood" in XML SVG DOM 
// The <feFlood> SVG filter primitive fills the filter subregion with the color and opacity defined by flood-color and flood-opacity.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFlood
func SvgFeFlood(id string, renders ...Render) (*Node, error) {
	return Element("feFlood", id, renders...)
}


// MustSvgFeFlood provides Node representation for the element "feFlood" in XML SVG DOM 
// The <feFlood> SVG filter primitive fills the filter subregion with the color and opacity defined by flood-color and flood-opacity.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFlood
func MustSvgFeFlood(id string, renders ...Render) *Node {
	var node, err = Element("feFlood", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeFuncA provides Node representation for the element "feFuncA" in XML SVG DOM 
// The <feFuncA> SVG filter primitive defines the transfer function for the alpha component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncA
func SvgFeFuncA(id string, renders ...Render) (*Node, error) {
	return Element("feFuncA", id, renders...)
}


// MustSvgFeFuncA provides Node representation for the element "feFuncA" in XML SVG DOM 
// The <feFuncA> SVG filter primitive defines the transfer function for the alpha component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncA
func MustSvgFeFuncA(id string, renders ...Render) *Node {
	var node, err = Element("feFuncA", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeFuncB provides Node representation for the element "feFuncB" in XML SVG DOM 
// The <feFuncB> SVG filter primitive defines the transfer function for the blue component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncB
func SvgFeFuncB(id string, renders ...Render) (*Node, error) {
	return Element("feFuncB", id, renders...)
}


// MustSvgFeFuncB provides Node representation for the element "feFuncB" in XML SVG DOM 
// The <feFuncB> SVG filter primitive defines the transfer function for the blue component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncB
func MustSvgFeFuncB(id string, renders ...Render) *Node {
	var node, err = Element("feFuncB", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeFuncG provides Node representation for the element "feFuncG" in XML SVG DOM 
// The <feFuncG> SVG filter primitive defines the transfer function for the green component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncG
func SvgFeFuncG(id string, renders ...Render) (*Node, error) {
	return Element("feFuncG", id, renders...)
}


// MustSvgFeFuncG provides Node representation for the element "feFuncG" in XML SVG DOM 
// The <feFuncG> SVG filter primitive defines the transfer function for the green component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncG
func MustSvgFeFuncG(id string, renders ...Render) *Node {
	var node, err = Element("feFuncG", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeFuncR provides Node representation for the element "feFuncR" in XML SVG DOM 
// The <feFuncR> SVG filter primitive defines the transfer function for the red component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncR
func SvgFeFuncR(id string, renders ...Render) (*Node, error) {
	return Element("feFuncR", id, renders...)
}


// MustSvgFeFuncR provides Node representation for the element "feFuncR" in XML SVG DOM 
// The <feFuncR> SVG filter primitive defines the transfer function for the red component of the input graphic of its parent <feComponentTransfer> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feFuncR
func MustSvgFeFuncR(id string, renders ...Render) *Node {
	var node, err = Element("feFuncR", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeGaussianBlur provides Node representation for the element "feGaussianBlur" in XML SVG DOM 
// The <feGaussianBlur> SVG filter primitive blurs the input image by the amount specified in stdDeviation, which defines the bell-curve.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feGaussianBlur
func SvgFeGaussianBlur(id string, renders ...Render) (*Node, error) {
	return Element("feGaussianBlur", id, renders...)
}


// MustSvgFeGaussianBlur provides Node representation for the element "feGaussianBlur" in XML SVG DOM 
// The <feGaussianBlur> SVG filter primitive blurs the input image by the amount specified in stdDeviation, which defines the bell-curve.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feGaussianBlur
func MustSvgFeGaussianBlur(id string, renders ...Render) *Node {
	var node, err = Element("feGaussianBlur", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeImage provides Node representation for the element "feImage" in XML SVG DOM 
// The <feImage> SVG filter primitive fetches image data from an external source and provides the pixel data as output (meaning if the external source is an SVG image, it is rasterized.)
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feImage
func SvgFeImage(id string, renders ...Render) (*Node, error) {
	return Element("feImage", id, renders...)
}


// MustSvgFeImage provides Node representation for the element "feImage" in XML SVG DOM 
// The <feImage> SVG filter primitive fetches image data from an external source and provides the pixel data as output (meaning if the external source is an SVG image, it is rasterized.)
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feImage
func MustSvgFeImage(id string, renders ...Render) *Node {
	var node, err = Element("feImage", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeMerge provides Node representation for the element "feMerge" in XML SVG DOM 
// The <feMerge> SVG element allows filter effects to be applied concurrently instead of sequentially. This is achieved by other filters storing their output via the result attribute and then accessing it in a <feMergeNode> child.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMerge
func SvgFeMerge(id string, renders ...Render) (*Node, error) {
	return Element("feMerge", id, renders...)
}


// MustSvgFeMerge provides Node representation for the element "feMerge" in XML SVG DOM 
// The <feMerge> SVG element allows filter effects to be applied concurrently instead of sequentially. This is achieved by other filters storing their output via the result attribute and then accessing it in a <feMergeNode> child.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMerge
func MustSvgFeMerge(id string, renders ...Render) *Node {
	var node, err = Element("feMerge", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeMergeNode provides Node representation for the element "feMergeNode" in XML SVG DOM 
// The feMergeNode takes the result of another filter to be processed by its parent <feMerge>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMergeNode
func SvgFeMergeNode(id string, renders ...Render) (*Node, error) {
	return Element("feMergeNode", id, renders...)
}


// MustSvgFeMergeNode provides Node representation for the element "feMergeNode" in XML SVG DOM 
// The feMergeNode takes the result of another filter to be processed by its parent <feMerge>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMergeNode
func MustSvgFeMergeNode(id string, renders ...Render) *Node {
	var node, err = Element("feMergeNode", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeMorphology provides Node representation for the element "feMorphology" in XML SVG DOM 
// The <feMorphology> SVG filter primitive is used to erode or dilate the input image. It's usefulness lies especially in fattening or thinning effects.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMorphology
func SvgFeMorphology(id string, renders ...Render) (*Node, error) {
	return Element("feMorphology", id, renders...)
}


// MustSvgFeMorphology provides Node representation for the element "feMorphology" in XML SVG DOM 
// The <feMorphology> SVG filter primitive is used to erode or dilate the input image. It's usefulness lies especially in fattening or thinning effects.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feMorphology
func MustSvgFeMorphology(id string, renders ...Render) *Node {
	var node, err = Element("feMorphology", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeOffset provides Node representation for the element "feOffset" in XML SVG DOM 
// The <feOffset> SVG filter primitive allows to offset the input image. The input image as a whole is offset by the values specified in the dx and dy attributes.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feOffset
func SvgFeOffset(id string, renders ...Render) (*Node, error) {
	return Element("feOffset", id, renders...)
}


// MustSvgFeOffset provides Node representation for the element "feOffset" in XML SVG DOM 
// The <feOffset> SVG filter primitive allows to offset the input image. The input image as a whole is offset by the values specified in the dx and dy attributes.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feOffset
func MustSvgFeOffset(id string, renders ...Render) *Node {
	var node, err = Element("feOffset", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFePointLight provides Node representation for the element "fePointLight" in XML SVG DOM 
// The <fePointLight> filter primitive defines a light source which allows to create a point light effect. It that can be used within a lighting filter primitive: <feDiffuseLighting> or <feSpecularLighting>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/fePointLight
func SvgFePointLight(id string, renders ...Render) (*Node, error) {
	return Element("fePointLight", id, renders...)
}


// MustSvgFePointLight provides Node representation for the element "fePointLight" in XML SVG DOM 
// The <fePointLight> filter primitive defines a light source which allows to create a point light effect. It that can be used within a lighting filter primitive: <feDiffuseLighting> or <feSpecularLighting>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/fePointLight
func MustSvgFePointLight(id string, renders ...Render) *Node {
	var node, err = Element("fePointLight", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeSpecularLighting provides Node representation for the element "feSpecularLighting" in XML SVG DOM 
// The <feSpecularLighting> SVG filter primitive lights a source graphic using the alpha channel as a bump map. The resulting image is an RGBA image based on the light color. The lighting calculation follows the standard specular component of the Phong lighting model. The resulting image depends on the light color, light position and surface geometry of the input bump map. The result of the lighting calculation is added. The filter primitive assumes that the viewer is at infinity in the z direction.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feSpecularLighting
func SvgFeSpecularLighting(id string, renders ...Render) (*Node, error) {
	return Element("feSpecularLighting", id, renders...)
}


// MustSvgFeSpecularLighting provides Node representation for the element "feSpecularLighting" in XML SVG DOM 
// The <feSpecularLighting> SVG filter primitive lights a source graphic using the alpha channel as a bump map. The resulting image is an RGBA image based on the light color. The lighting calculation follows the standard specular component of the Phong lighting model. The resulting image depends on the light color, light position and surface geometry of the input bump map. The result of the lighting calculation is added. The filter primitive assumes that the viewer is at infinity in the z direction.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feSpecularLighting
func MustSvgFeSpecularLighting(id string, renders ...Render) *Node {
	var node, err = Element("feSpecularLighting", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeSpotLight provides Node representation for the element "feSpotLight" in XML SVG DOM 
// The <feSpotLight> SVG filter primitive defines a light source which allows to create a spotlight effect. It that can be used within a lighting filter primitive: <feDiffuseLighting> or <feSpecularLighting>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feSpotLight
func SvgFeSpotLight(id string, renders ...Render) (*Node, error) {
	return Element("feSpotLight", id, renders...)
}


// MustSvgFeSpotLight provides Node representation for the element "feSpotLight" in XML SVG DOM 
// The <feSpotLight> SVG filter primitive defines a light source which allows to create a spotlight effect. It that can be used within a lighting filter primitive: <feDiffuseLighting> or <feSpecularLighting>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feSpotLight
func MustSvgFeSpotLight(id string, renders ...Render) *Node {
	var node, err = Element("feSpotLight", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeTile provides Node representation for the element "feTile" in XML SVG DOM 
// The <feTile> SVG filter primitive allows to fill a target rectangle with a repeated, tiled pattern of an input image. The effect is similar to the one of a <pattern>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feTile
func SvgFeTile(id string, renders ...Render) (*Node, error) {
	return Element("feTile", id, renders...)
}


// MustSvgFeTile provides Node representation for the element "feTile" in XML SVG DOM 
// The <feTile> SVG filter primitive allows to fill a target rectangle with a repeated, tiled pattern of an input image. The effect is similar to the one of a <pattern>.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feTile
func MustSvgFeTile(id string, renders ...Render) *Node {
	var node, err = Element("feTile", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFeTurbulence provides Node representation for the element "feTurbulence" in XML SVG DOM 
// The <feTurbulence> SVG filter primitive creates an image using the Perlin turbulence function. It allows the synthesis of artificial textures like clouds or marble. The resulting image will fill the entire filter primitive subregion.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feTurbulence
func SvgFeTurbulence(id string, renders ...Render) (*Node, error) {
	return Element("feTurbulence", id, renders...)
}


// MustSvgFeTurbulence provides Node representation for the element "feTurbulence" in XML SVG DOM 
// The <feTurbulence> SVG filter primitive creates an image using the Perlin turbulence function. It allows the synthesis of artificial textures like clouds or marble. The resulting image will fill the entire filter primitive subregion.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/feTurbulence
func MustSvgFeTurbulence(id string, renders ...Render) *Node {
	var node, err = Element("feTurbulence", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgFilter provides Node representation for the element "filter" in XML SVG DOM 
// The <filter> SVG element serves as container for atomic filter operations. It is never rendered directly. A filter is referenced by using the filter attribute on the target SVG element or via the filter CSS property.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/filter
func SvgFilter(id string, renders ...Render) (*Node, error) {
	return Element("filter", id, renders...)
}


// MustSvgFilter provides Node representation for the element "filter" in XML SVG DOM 
// The <filter> SVG element serves as container for atomic filter operations. It is never rendered directly. A filter is referenced by using the filter attribute on the target SVG element or via the filter CSS property.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/filter
func MustSvgFilter(id string, renders ...Render) *Node {
	var node, err = Element("filter", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgForeignObject provides Node representation for the element "foreignObject" in XML SVG DOM 
// The <foreignObject> SVG element allows for inclusion of a different XML namespace. In the context of a browser it is most likely XHTML/HTML.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/foreignObject
func SvgForeignObject(id string, renders ...Render) (*Node, error) {
	return Element("foreignObject", id, renders...)
}


// MustSvgForeignObject provides Node representation for the element "foreignObject" in XML SVG DOM 
// The <foreignObject> SVG element allows for inclusion of a different XML namespace. In the context of a browser it is most likely XHTML/HTML.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/foreignObject
func MustSvgForeignObject(id string, renders ...Render) *Node {
	var node, err = Element("foreignObject", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgGroup provides Node representation for the element "g" in XML SVG DOM 
// The <g> SVG element is a container used to group other SVG elements.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/g
func SvgGroup(id string, renders ...Render) (*Node, error) {
	return Element("g", id, renders...)
}


// MustSvgGroup provides Node representation for the element "g" in XML SVG DOM 
// The <g> SVG element is a container used to group other SVG elements.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/g
func MustSvgGroup(id string, renders ...Render) *Node {
	var node, err = Element("g", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgHatch provides Node representation for the element "hatch" in XML SVG DOM 
// The <hatch> SVG element is used to fill or stroke an object using one or more pre-defined paths that are repeated at fixed intervals in a specified direction to cover the areas to be painted.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/hatch
func SvgHatch(id string, renders ...Render) (*Node, error) {
	return Element("hatch", id, renders...)
}


// MustSvgHatch provides Node representation for the element "hatch" in XML SVG DOM 
// The <hatch> SVG element is used to fill or stroke an object using one or more pre-defined paths that are repeated at fixed intervals in a specified direction to cover the areas to be painted.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/hatch
func MustSvgHatch(id string, renders ...Render) *Node {
	var node, err = Element("hatch", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgHatchpath provides Node representation for the element "hatchpath" in XML SVG DOM 
// The <hatchpath> SVG element defines a hatch path used by the <hatch> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/hatchpath
func SvgHatchpath(id string, renders ...Render) (*Node, error) {
	return Element("hatchpath", id, renders...)
}


// MustSvgHatchpath provides Node representation for the element "hatchpath" in XML SVG DOM 
// The <hatchpath> SVG element defines a hatch path used by the <hatch> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/hatchpath
func MustSvgHatchpath(id string, renders ...Render) *Node {
	var node, err = Element("hatchpath", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgImage provides Node representation for the element "image" in XML SVG DOM 
// The <image> SVG element includes images inside SVG documents. It can display raster image files or other SVG files.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/image
func SvgImage(id string, renders ...Render) (*Node, error) {
	return Element("image", id, renders...)
}


// MustSvgImage provides Node representation for the element "image" in XML SVG DOM 
// The <image> SVG element includes images inside SVG documents. It can display raster image files or other SVG files.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/image
func MustSvgImage(id string, renders ...Render) *Node {
	var node, err = Element("image", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgLine provides Node representation for the element "line" in XML SVG DOM 
// The <line> element is an SVG basic shape used to create a line connecting two points.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/line
func SvgLine(id string, renders ...Render) (*Node, error) {
	return Element("line", id, renders...)
}


// MustSvgLine provides Node representation for the element "line" in XML SVG DOM 
// The <line> element is an SVG basic shape used to create a line connecting two points.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/line
func MustSvgLine(id string, renders ...Render) *Node {
	var node, err = Element("line", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgLinearGradient provides Node representation for the element "linearGradient" in XML SVG DOM 
// The <linearGradient> element lets authors define linear gradients that can be applied to fill or stroke of graphical elements.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/linearGradient
func SvgLinearGradient(id string, renders ...Render) (*Node, error) {
	return Element("linearGradient", id, renders...)
}


// MustSvgLinearGradient provides Node representation for the element "linearGradient" in XML SVG DOM 
// The <linearGradient> element lets authors define linear gradients that can be applied to fill or stroke of graphical elements.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/linearGradient
func MustSvgLinearGradient(id string, renders ...Render) *Node {
	var node, err = Element("linearGradient", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMarker provides Node representation for the element "marker" in XML SVG DOM 
// The <marker> element defines the graphic that is to be used for drawing arrowheads or polymarkers on a given <path>, <line>, <polyline> or <polygon> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/marker
func SvgMarker(id string, renders ...Render) (*Node, error) {
	return Element("marker", id, renders...)
}


// MustSvgMarker provides Node representation for the element "marker" in XML SVG DOM 
// The <marker> element defines the graphic that is to be used for drawing arrowheads or polymarkers on a given <path>, <line>, <polyline> or <polygon> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/marker
func MustSvgMarker(id string, renders ...Render) *Node {
	var node, err = Element("marker", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMask provides Node representation for the element "mask" in XML SVG DOM 
// The <mask> element defines an alpha mask for compositing the current object into the background. A mask is used/referenced using the mask property.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mask
func SvgMask(id string, renders ...Render) (*Node, error) {
	return Element("mask", id, renders...)
}


// MustSvgMask provides Node representation for the element "mask" in XML SVG DOM 
// The <mask> element defines an alpha mask for compositing the current object into the background. A mask is used/referenced using the mask property.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mask
func MustSvgMask(id string, renders ...Render) *Node {
	var node, err = Element("mask", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMesh provides Node representation for the element "mesh" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mesh
func SvgMesh(id string, renders ...Render) (*Node, error) {
	return Element("mesh", id, renders...)
}


// MustSvgMesh provides Node representation for the element "mesh" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mesh
func MustSvgMesh(id string, renders ...Render) *Node {
	var node, err = Element("mesh", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMeshgradient provides Node representation for the element "meshgradient" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/meshgradient
func SvgMeshgradient(id string, renders ...Render) (*Node, error) {
	return Element("meshgradient", id, renders...)
}


// MustSvgMeshgradient provides Node representation for the element "meshgradient" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/meshgradient
func MustSvgMeshgradient(id string, renders ...Render) *Node {
	var node, err = Element("meshgradient", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMeshpatch provides Node representation for the element "meshpatch" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/meshpatch
func SvgMeshpatch(id string, renders ...Render) (*Node, error) {
	return Element("meshpatch", id, renders...)
}


// MustSvgMeshpatch provides Node representation for the element "meshpatch" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/meshpatch
func MustSvgMeshpatch(id string, renders ...Render) *Node {
	var node, err = Element("meshpatch", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMeshrow provides Node representation for the element "meshrow" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/meshrow
func SvgMeshrow(id string, renders ...Render) (*Node, error) {
	return Element("meshrow", id, renders...)
}


// MustSvgMeshrow provides Node representation for the element "meshrow" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/meshrow
func MustSvgMeshrow(id string, renders ...Render) *Node {
	var node, err = Element("meshrow", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMetadata provides Node representation for the element "metadata" in XML SVG DOM 
// The <metadata> SVG element allows to add metadata to SVG content. Metadata is structured information about data. The contents of <metadata> elements should be elements from other XML namespaces such as RDF, FOAF, etc.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/metadata
func SvgMetadata(id string, renders ...Render) (*Node, error) {
	return Element("metadata", id, renders...)
}


// MustSvgMetadata provides Node representation for the element "metadata" in XML SVG DOM 
// The <metadata> SVG element allows to add metadata to SVG content. Metadata is structured information about data. The contents of <metadata> elements should be elements from other XML namespaces such as RDF, FOAF, etc.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/metadata
func MustSvgMetadata(id string, renders ...Render) *Node {
	var node, err = Element("metadata", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgMpath provides Node representation for the element "mpath" in XML SVG DOM 
// The <mpath> sub-element for the <animateMotion> element provides the ability to reference an external <path> element as the definition of a motion path.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mpath
func SvgMpath(id string, renders ...Render) (*Node, error) {
	return Element("mpath", id, renders...)
}


// MustSvgMpath provides Node representation for the element "mpath" in XML SVG DOM 
// The <mpath> sub-element for the <animateMotion> element provides the ability to reference an external <path> element as the definition of a motion path.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/mpath
func MustSvgMpath(id string, renders ...Render) *Node {
	var node, err = Element("mpath", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgPath provides Node representation for the element "path" in XML SVG DOM 
// The <path> SVG element is the generic element to define a shape. All the basic shapes can be created with a path element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/path
func SvgPath(id string, renders ...Render) (*Node, error) {
	return Element("path", id, renders...)
}


// MustSvgPath provides Node representation for the element "path" in XML SVG DOM 
// The <path> SVG element is the generic element to define a shape. All the basic shapes can be created with a path element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/path
func MustSvgPath(id string, renders ...Render) *Node {
	var node, err = Element("path", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgPattern provides Node representation for the element "pattern" in XML SVG DOM 
// The <pattern> element defines a graphics object which can be redrawn at repeated x and y-coordinate intervals ("tiled") to cover an area.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/pattern
func SvgPattern(id string, renders ...Render) (*Node, error) {
	return Element("pattern", id, renders...)
}


// MustSvgPattern provides Node representation for the element "pattern" in XML SVG DOM 
// The <pattern> element defines a graphics object which can be redrawn at repeated x and y-coordinate intervals ("tiled") to cover an area.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/pattern
func MustSvgPattern(id string, renders ...Render) *Node {
	var node, err = Element("pattern", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgPolygon provides Node representation for the element "polygon" in XML SVG DOM 
// The <polygon> element defines a closed shape consisting of a set of connected straight line segments. The last point is connected to the first point. For open shapes see the <polyline> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/polygon
func SvgPolygon(id string, renders ...Render) (*Node, error) {
	return Element("polygon", id, renders...)
}


// MustSvgPolygon provides Node representation for the element "polygon" in XML SVG DOM 
// The <polygon> element defines a closed shape consisting of a set of connected straight line segments. The last point is connected to the first point. For open shapes see the <polyline> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/polygon
func MustSvgPolygon(id string, renders ...Render) *Node {
	var node, err = Element("polygon", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgPolyline provides Node representation for the element "polyline" in XML SVG DOM 
// The <polyline> SVG element is an SVG basic shape that creates straight lines connecting several points. Typically a polyline is used to create open shapes as the last point doesn't have to be connected to the first point. For closed shapes see the <polygon> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/polyline
func SvgPolyline(id string, renders ...Render) (*Node, error) {
	return Element("polyline", id, renders...)
}


// MustSvgPolyline provides Node representation for the element "polyline" in XML SVG DOM 
// The <polyline> SVG element is an SVG basic shape that creates straight lines connecting several points. Typically a polyline is used to create open shapes as the last point doesn't have to be connected to the first point. For closed shapes see the <polygon> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/polyline
func MustSvgPolyline(id string, renders ...Render) *Node {
	var node, err = Element("polyline", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgRadialGradient provides Node representation for the element "radialGradient" in XML SVG DOM 
// The <radialGradient> SVG element lets authors define radial gradients to fill or stroke graphical elements.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/radialGradient
func SvgRadialGradient(id string, renders ...Render) (*Node, error) {
	return Element("radialGradient", id, renders...)
}


// MustSvgRadialGradient provides Node representation for the element "radialGradient" in XML SVG DOM 
// The <radialGradient> SVG element lets authors define radial gradients to fill or stroke graphical elements.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/radialGradient
func MustSvgRadialGradient(id string, renders ...Render) *Node {
	var node, err = Element("radialGradient", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgRect provides Node representation for the element "rect" in XML SVG DOM 
// The <rect> element is a basic SVG shape that creates rectangles, defined by their corner's position, their width, and their height. The rectangles may have their corners rounded.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/rect
func SvgRect(id string, renders ...Render) (*Node, error) {
	return Element("rect", id, renders...)
}


// MustSvgRect provides Node representation for the element "rect" in XML SVG DOM 
// The <rect> element is a basic SVG shape that creates rectangles, defined by their corner's position, their width, and their height. The rectangles may have their corners rounded.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/rect
func MustSvgRect(id string, renders ...Render) *Node {
	var node, err = Element("rect", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgScript provides Node representation for the element "script" in XML SVG DOM 
// A SVG script element is equivalent to the script element in HTML and thus is the place for scripts (e.g., ECMAScript).
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/script
func SvgScript(id string, renders ...Render) (*Node, error) {
	return Element("script", id, renders...)
}


// MustSvgScript provides Node representation for the element "script" in XML SVG DOM 
// A SVG script element is equivalent to the script element in HTML and thus is the place for scripts (e.g., ECMAScript).
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/script
func MustSvgScript(id string, renders ...Render) *Node {
	var node, err = Element("script", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgSet provides Node representation for the element "set" in XML SVG DOM 
// The <set> element provides a simple means of just setting the value of an attribute for a specified duration. It supports all attribute types, including those that cannot reasonably be interpolated, such as string and boolean values. The <set> element is non-additive. The additive and accumulate attributes are not allowed, and will be ignored if specified.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/set
func SvgSet(id string, renders ...Render) (*Node, error) {
	return Element("set", id, renders...)
}


// MustSvgSet provides Node representation for the element "set" in XML SVG DOM 
// The <set> element provides a simple means of just setting the value of an attribute for a specified duration. It supports all attribute types, including those that cannot reasonably be interpolated, such as string and boolean values. The <set> element is non-additive. The additive and accumulate attributes are not allowed, and will be ignored if specified.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/set
func MustSvgSet(id string, renders ...Render) *Node {
	var node, err = Element("set", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgSolidcolor provides Node representation for the element "solidcolor" in XML SVG DOM 
// The <solidcolor> SVG element lets authors define a single color for use in multiple places in an SVG document. It is also useful as away of animating a palette colors.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/solidcolor
func SvgSolidcolor(id string, renders ...Render) (*Node, error) {
	return Element("solidcolor", id, renders...)
}


// MustSvgSolidcolor provides Node representation for the element "solidcolor" in XML SVG DOM 
// The <solidcolor> SVG element lets authors define a single color for use in multiple places in an SVG document. It is also useful as away of animating a palette colors.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/solidcolor
func MustSvgSolidcolor(id string, renders ...Render) *Node {
	var node, err = Element("solidcolor", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgStop provides Node representation for the element "stop" in XML SVG DOM 
// The <stop> SVG element defines the ramp of colors to use on a gradient, which is a child element to either the <linearGradient> or the <radialGradient> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/stop
func SvgStop(id string, renders ...Render) (*Node, error) {
	return Element("stop", id, renders...)
}


// MustSvgStop provides Node representation for the element "stop" in XML SVG DOM 
// The <stop> SVG element defines the ramp of colors to use on a gradient, which is a child element to either the <linearGradient> or the <radialGradient> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/stop
func MustSvgStop(id string, renders ...Render) *Node {
	var node, err = Element("stop", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgStyle provides Node representation for the element "style" in XML SVG DOM 
// The <style> SVG element allows style sheets to be embedded directly within SVG content. SVG's style element has the same attributes as the corresponding element in HTML (see HTML's <style> element).
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/style
func SvgStyle(id string, renders ...Render) (*Node, error) {
	return Element("style", id, renders...)
}


// MustSvgStyle provides Node representation for the element "style" in XML SVG DOM 
// The <style> SVG element allows style sheets to be embedded directly within SVG content. SVG's style element has the same attributes as the corresponding element in HTML (see HTML's <style> element).
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/style
func MustSvgStyle(id string, renders ...Render) *Node {
	var node, err = Element("style", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// Svg provides Node representation for the element "svg" in XML SVG DOM 
// The svg element is a container that defines a new coordinate system and viewport. It is used as the outermost element of any SVG document but it can also be used to embed a SVG fragment inside any SVG or HTML document.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/svg
func Svg(id string, renders ...Render) (*Node, error) {
	return Element("svg", id, renders...)
}


// MustSvg provides Node representation for the element "svg" in XML SVG DOM 
// The svg element is a container that defines a new coordinate system and viewport. It is used as the outermost element of any SVG document but it can also be used to embed a SVG fragment inside any SVG or HTML document.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/svg
func MustSvg(id string, renders ...Render) *Node {
	var node, err = Element("svg", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgSwitch provides Node representation for the element "switch" in XML SVG DOM 
// The <switch> SVG element evaluates the requiredFeatures, requiredExtensions and systemLanguage attributes on its direct child elements in order, and then processes and renders the first child for which these attributes evaluate to true. All others will be bypassed and therefore not rendered. If the child element is a container element such as a <g>, then the entire subtree is either processed/rendered or bypassed/not rendered.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/switch
func SvgSwitch(id string, renders ...Render) (*Node, error) {
	return Element("switch", id, renders...)
}


// MustSvgSwitch provides Node representation for the element "switch" in XML SVG DOM 
// The <switch> SVG element evaluates the requiredFeatures, requiredExtensions and systemLanguage attributes on its direct child elements in order, and then processes and renders the first child for which these attributes evaluate to true. All others will be bypassed and therefore not rendered. If the child element is a container element such as a <g>, then the entire subtree is either processed/rendered or bypassed/not rendered.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/switch
func MustSvgSwitch(id string, renders ...Render) *Node {
	var node, err = Element("switch", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgSymbol provides Node representation for the element "symbol" in XML SVG DOM 
// The <symbol> element is used to define graphical template objects which can be instantiated by a <use> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/symbol
func SvgSymbol(id string, renders ...Render) (*Node, error) {
	return Element("symbol", id, renders...)
}


// MustSvgSymbol provides Node representation for the element "symbol" in XML SVG DOM 
// The <symbol> element is used to define graphical template objects which can be instantiated by a <use> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/symbol
func MustSvgSymbol(id string, renders ...Render) *Node {
	var node, err = Element("symbol", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgText provides Node representation for the element "text" in XML SVG DOM 
// The SVG <text> element defines a graphics element consisting of text. It's possible to apply a gradient, pattern, clipping path, mask, or filter to <text>, just like any other SVG graphics element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/text
func SvgText(id string, renders ...Render) (*Node, error) {
	return Element("text", id, renders...)
}


// MustSvgText provides Node representation for the element "text" in XML SVG DOM 
// The SVG <text> element defines a graphics element consisting of text. It's possible to apply a gradient, pattern, clipping path, mask, or filter to <text>, just like any other SVG graphics element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/text
func MustSvgText(id string, renders ...Render) *Node {
	var node, err = Element("text", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgTextPath provides Node representation for the element "textPath" in XML SVG DOM 
// In addition to text drawn in a straight line, SVG also includes the ability to place text along the shape of a <path> element. To specify that a block of text is to be rendered along the shape of a <path>, include the given text within a <textPath> element which includes an href attribute with a reference to a <path> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/textPath
func SvgTextPath(id string, renders ...Render) (*Node, error) {
	return Element("textPath", id, renders...)
}


// MustSvgTextPath provides Node representation for the element "textPath" in XML SVG DOM 
// In addition to text drawn in a straight line, SVG also includes the ability to place text along the shape of a <path> element. To specify that a block of text is to be rendered along the shape of a <path>, include the given text within a <textPath> element which includes an href attribute with a reference to a <path> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/textPath
func MustSvgTextPath(id string, renders ...Render) *Node {
	var node, err = Element("textPath", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgTitle provides Node representation for the element "title" in XML SVG DOM 
// Each container element or graphics element in an SVG drawing can supply a <title> element containing a description string where the description is text-only. When the current SVG document fragment is rendered as SVG on visual media, <title> element is not rendered as part of the graphics. However, some user agents may, for example, display the <title> element as a tooltip. Alternate presentations are possible, both visual and aural, which display the <title> element but do not display path elements or other graphics elements. The <title> element generally improves accessibility of SVG documents.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/title
func SvgTitle(id string, renders ...Render) (*Node, error) {
	return Element("title", id, renders...)
}


// MustSvgTitle provides Node representation for the element "title" in XML SVG DOM 
// Each container element or graphics element in an SVG drawing can supply a <title> element containing a description string where the description is text-only. When the current SVG document fragment is rendered as SVG on visual media, <title> element is not rendered as part of the graphics. However, some user agents may, for example, display the <title> element as a tooltip. Alternate presentations are possible, both visual and aural, which display the <title> element but do not display path elements or other graphics elements. The <title> element generally improves accessibility of SVG documents.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/title
func MustSvgTitle(id string, renders ...Render) *Node {
	var node, err = Element("title", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgTspan provides Node representation for the element "tspan" in XML SVG DOM 
// Within a <text> element, text and font properties and the current text position can be adjusted with absolute or relative coordinate values by including a <tspan> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/tspan
func SvgTspan(id string, renders ...Render) (*Node, error) {
	return Element("tspan", id, renders...)
}


// MustSvgTspan provides Node representation for the element "tspan" in XML SVG DOM 
// Within a <text> element, text and font properties and the current text position can be adjusted with absolute or relative coordinate values by including a <tspan> element.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/tspan
func MustSvgTspan(id string, renders ...Render) *Node {
	var node, err = Element("tspan", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgUnknown provides Node representation for the element "unknown" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/unknown
func SvgUnknown(id string, renders ...Render) (*Node, error) {
	return Element("unknown", id, renders...)
}


// MustSvgUnknown provides Node representation for the element "unknown" in XML SVG DOM 
// The documentation about this has not yet been written; please consider contributing!
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/unknown
func MustSvgUnknown(id string, renders ...Render) *Node {
	var node, err = Element("unknown", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgUse provides Node representation for the element "use" in XML SVG DOM 
// The <use> element takes nodes from within the SVG document, and duplicates them somewhere else.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/use
func SvgUse(id string, renders ...Render) (*Node, error) {
	return Element("use", id, renders...)
}


// MustSvgUse provides Node representation for the element "use" in XML SVG DOM 
// The <use> element takes nodes from within the SVG document, and duplicates them somewhere else.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/use
func MustSvgUse(id string, renders ...Render) *Node {
	var node, err = Element("use", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}


// SvgView provides Node representation for the element "view" in XML SVG DOM 
// A view is a defined way to view the image, like a zoom level or a detail view.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/view
func SvgView(id string, renders ...Render) (*Node, error) {
	return Element("view", id, renders...)
}


// MustSvgView provides Node representation for the element "view" in XML SVG DOM 
// A view is a defined way to view the image, like a zoom level or a detail view.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Element/view
func MustSvgView(id string, renders ...Render) *Node {
	var node, err = Element("view", id, renders...)
	if err != nil {
		panic(err)
	}
	return node
}

