In this post, I will talk about a Nim library that I am really fond of called Futhark, how it works, and some examples.

## What is Futhark?
[Futhark](https://github.com/PMunch/futhark) is a library that allows you to automatically wrap C libraries, and use them in Nim.

It uses clang to understand C syntax and semantics, and is very accurate compared to c2nim and nimterop for this reason.

Futhark is in a beta state, though it is mostly stable in my experience. Futhark is made by PMunch, a person that is very well-known in the Nim community, and that has a lot of useful libraries for the language.

## My experience with Futhark
As some of you may know, I am the main developer and founder of [Kreato Linux](https://linux.kreato.dev). As an independent Linux distribution, we have our own package manager, build system and the like. Writing these in pure Nim is quite a challenge as the package ecosystem is not mature enough (unlike languages like Rust, Go, etc.) to have native, full-featured libraries for things like Git, archiving/unarchiving files, et cetera. While I was rewriting functions of kpkg in Nim, I have noticed that invoking the shell for such purposes is inefficient, and can result in undefined errors. This was mainly the reason as to why I have started looking for a C interop tool.

There were a couple of options, starting with interop by hand (manual). This was not viable as I wanted to use huge libraries like libarchive and libgit2, which would take days to implement correctly. Another option was using nimterop/c2nim, but I never managed to get a viable, ready-to-use output from them. Then I discovered Futhark, and never looked back.

## Code examples
* On kpkg, there is [modules/libarchive.nim](https://github.com/kreatolinux/src/blob/master/kpkg/modules/libarchive.nim) which is basically a reimplementation minitar/untar.c on libarchive (the createArchive part is still a work-in-progress)
* [The tests of Futhark](https://github.com/PMunch/futhark/tree/master/tests) also has good examples.
* [The Futhark readme](https://github.com/PMunch/futhark) has good examples, and pretty much enough to get started.

## Conclusion
In general, Nim is an incredible language. But without an ecosystem of packages, things can get dire sometimes. With Futhark though, its usable in my opinion. If you are a current Nim user, I would recommend you to check out Futhark, especially if you struggle with ecosystem issues like I was. I would also recommend people that use languages like Zig to check out Futhark, as it is pretty impressive, and if you use these languages for their C interop, this library may make you try out Nim.
