<!--
Tags: linux, packaging, package-management, development, kpkg
-->

# How to write a package format

In this blog post, I will write about the design decisions and the process that went behind kpkg runFile v3. This blog post will also be available on https://linux.krea.to.

First of all, we will start with a brief introduction to kpkg and its history on its packaging formats.

## History

First iterations of Kreato Linux used completely static packages - that is, the few components that were available were in the source code itself. This makes 0 sense now, but back then it made sense for rapid prototyping and deployment. Eventually, first version of `nyaa` (old name of `kpkg`) introduced runFiles as a simple package spec that wasn't hardcoded.

### runFile v1

v1 was very primitive, I just threw around something that I felt was convenient enough to implement and use.

```sh
NAME="test"
VERSION="0.0.1"
SOURCES="https://test.file/source/testfile.tar.gz"
SHA256SUM="e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  testfile.tar.gz"
DESCRIPTION="Test package"

build() {
    tar -xvf testfile.tar.gz
    echo "Insert additional installation instructions here"
}
```

It didn't have any other functions other than `build`. Back then, builds (and installations!) were made through root. Dependencies were also handled in a seperate file, called `deps` and (later) `build_deps`.

After some time, additional features such as RELEASE, EPOCH, prepare, post install, post uninstall, and custom uninstall (not existent anymore for obvious reasons) were added.

I am suprised this worked so well for so long, but it was obvious a change was required.

### runFile v2

After a while, I've realized that I've pushed POSIX sh enough and I should switch to an actual programming language before it is too late. That was when I discovered Nim and rewrote `nyaa` in it (maybe I'll write more details about this later, but for now I'll leave it at that).

As seen in https://linux.krea.to/blog/new-workflow/ I've centralized the entire source code of Kreato Linux is what is essentially 3 repositories. While doing this major refacotr, I had the time to look at the runFile specification and fix a lot of gripes.

```sh
NAME="test"
VERSION="0.0.1"
RELEASE="1"
SOURCES="https://test.file/source/testfile.tar.gz;git::https://github.com/kreatolinux/src::543ee30eda806029fa9ea16a1f9767eda7cab4d1"
DEPENDS="testpackage1 testpackage3 testpackage4"
BUILD_DEPENDS="testpackage5 testpackage6 testpackage10"
SHA256SUM="e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  testfile.tar.gz"
DESCRIPTION="Test package"

prepare() {
    tar -xvf testfile.tar.gz
}

build() {
    cd testfile
    echo "Insert build instructions here"
}

package() {
    cd testfile
    make install
}

package_test2() {
    cd testfile
    make install_test2
}

postinstall() {
    echo "Insert postinstall instructions here"
}
```

This is one of the later versions of runFile v2, I won't mention other breaking changes as they are boring and not really major (unless you look at strictly compatibility) such as removing `;` to seperate multiple sources/sums in favor of using a space, adding more features, etc.

This time runFile was a proper, comparable format to major Linux distributions.

* Multiple sum support
* In-line dependencies instead of an external file
* check() to test the package beforehand
* preupgrade() and postupgrade() to handle more edgecases
* variables such as BACKUP that allowed the packager to keep certain config files (if preferred)
* Fully nim-based variable and, later function parser
* and more that I probably forgot

This is the current runFile major version we are at now, and it has mostly stayed the same (except a few small breaking issues which can easily be migrated to another format)

So, why the change now?

## Reasons why this is changing

### Simplicity

As mentioned in our motto, on Kreato Linux we are always focused on simplicity. This includes the developer experience.

The new format will be much simpler to use for old users (because of similarity) and for new users (because it is closer to a modern language than a text file.)

* The new format uses arrays extensively, which we couldn't do with POSIX sh (available on bash but read on)
* The new format requires 0 external dependencies, meaning;
  * We eliminate potential external mismatch/issues, making the packages more reproducible as a result
  * The kpkg source code is simpler and easier to maintain

### Speed

Obviously this change will also most definitely bring speed. Currently on runFile v2, we run external commands a LOT (seriously, look at kpkg/commands/buildcmd.nim), and while that isn't going to 0 with this update, it will certainly decrease external commands used, meaning more speed!

### Growing out of sh

Don't get me wrong, POSIX sh is one of the best if not the best scripting language there is (much better than the likes of Powershell) but it simply isn't suited for this task. I could've just wrote a parser for sh and extended it to get most of these benefits but it would have major complications.

The final goal of a(ny) package manager in my opinion should be to decrease the usage of the shell as much as possible, and containing it very well on places where it cannot be done. With this change, we are also bringing proper macro support which works MUCH better than the current implementation (which is basically a glorified alias) which means much lower direct shell usage on packages.

### Conclusion

I can probably add more reasons but these are the main ones that made the push into making runFile v3 a reality.

Now let's talk about this new format;

## runFile v3: third time is the charm

First thing i wanted to do was what to base this off of. I thought to myself, there surely must be a simple, common language format (like YAML but more suited for packaging) that we can use.

Turns out, there is, but the caveats are too costly.

* HCL: This was my first pick because i was on my OpenTofu/Terraform stage at that moment and it made sense. The main issues with it was it was not as easy to read as YAML/sh and it didn't support Nim.
* KDL: This was my second pick mainly because I came across it in a search. It looks pretty nice, it may do the job nicely but there wasn't a library in Nim that supported the KDL V2 spec. And honestly, I *really* don't want major changes after the change this time if we can afford it.
* YAML: This made a lot of sense because there was already a well-established library on Nim and it was a well-known format for packaging. But I believe it lacks the readability and extensibility that POSIX sh had, and I wanted extensibility at all costs.

So, I started doing my own language spec.

### How I started

I started by looking into YAML. Yes, I still didn't want it as the entire format but i felt like for variables it was "good enough". It has first class support for things like arrays which we really need at this point on most of the variables. It also looks nice, is readable, etc. But I didn't like the lack of flexibility it had on things like packaging steps.

Then I had an idea; What if I, just rip YAML just for variables and make my own format for rest?

So that was how I started, a frankenstein of YAML and my random ideas. Of course, there were many questions in my mind.

### Variable expansion

Of course, while YAML for variables is great it doesn't have variable expansion which we so desperately need for things like sources.

For this, I debated a lot about whether I should just copy the syntax of another language or make my own.

I've decided for the former, taking Python as a loose base while making my own changes that made it more v2-like.

An example is `${version.split('.')[0:2].join('.')}`. As you can see, it is *very* similar to Python. And that is by design, because I don't want people to feel like this is entirely new territory like some packaging formats.

### Functions

Okay, we got out of the variables part by just yanking YAML, what about functions? Obviously I didn't want a mess like YAML, so I've got inspired by Go to make the current format.

```sh
func custom_func {
    print "This is a test custom function"
}

prepare {
    custom_func # We can run custom functions just like commands
    exec tar -xvf testfile.tar.gz
    # Or you can use:
    # macro extract --autocd=true
    # to extract all the archives in sources.
}
```

As you can see, its not *too* different from v2, but its clearly not the same either.

The major improvements in the function part are mainly the fact that custom functions are clearly custom (as seen with the `func` keyword) and you cannot use shell commands in them unless you specifically invoke them with `exec`.

### If/else statements and For loops

The story is similar here actually, mostly inspired from Go and some from Python because I was mainly using Go at the time, while making things look simpler.

```go
    if test1 {
        print hello!
    } else {
        print bye!
    }

    for i in test2 {
        print $i
    }
```

### Commands

Since this is not a shell, and it is similar to a shell script where we run "commands" (you can also call them functions, it doesn't really matter) we needed a lot of commands such as macro, cd, exec, print, etc. These commands have their own arguments to make them feel more natural.

```sh
macro package --meson
exec tar -xvf example.tar.gz
#...
```

This seemed like the right approach to me, but suggestions are welcome if theres a major problem with this!

### Conclusion

This is the final product;

```nim
name: "test-v3"
version: "0.0.1"
release: "1"
sources: 
    - "https://test.file/source/testfile.tar.gz"
    - "git::https://github.com/kreatolinux/src::543ee30eda806029fa9ea16a1f9767eda7cab4d1"
    - "https://test.file/sources/v${version.split('.')[0:2].join('.')}/testfile.tar.gz"
depends: 
    - "testpackage1" 
    - "testpackage3" 
    - "testpackage4"
depends_test2: 
    - "testpackage5" 
    - "testpackage6"
no_chkupd: false
replaces: 
    - "test-v2"
backup: 
    - "etc/test-v3/main.conf"
    - "etc/test/settings.conf"
opt_depends:
    - "optional-dependency: This is a test optional dependency"
    - "optional-dependency-2: This is a second optional dependency."
build_depends: 
    - "testpackage5" 
    - "testpackage6" 
    - "testpackage10"
sha256sum: 
    - "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    - "SKIP"
    - "ab37404db60460d548e22546a83fda1eb6061cd8b95e37455149f3baf6c5fd38"
description: "Test package"

func custom_func {
    print "This is a test custom function"
}

prepare {
    custom_func # We can run custom functions just like commands
    exec tar -xvf testfile.tar.gz
    # Or you can use:
    # macro extract --autocd=true
    # to extract all the archives in sources.
}

build {
    cd testfile
    echo "Insert build instructions here"
}

check {
    macro test --ninja
    # exec ninja -C build test
}

preupgrade {
    echo "run before upgrade"
}

preinstall {
    echo "run before first install"
}

package {
    cd testfile
    macro package --meson
}

package_test2 {
    cd testfile
    exec make install_test2 # External commands require exec
}

postinstall {
    echo "Insert postinstall instructions here"
}

postupgrade {
    echo "run after upgrade"
}

postremove {
    echo "Insert postremove instructions here"
}
```

I'd say it turned out pretty well. It exactly feels like the "best of both worlds" to me. Of course any feedback is welcome.

After this, there was another problem;

### Migration

Migration was obviously going to be a pain in the butt. A lot of packages just used meson/configure directly with no decent automated migration path.

For this, I used a simple converter plus manual conversion.

Of course, this also meant that we needed to support both v2 and v3 at the same time on kpkg-repo in atleast one version, or put some kind of automated migration path.

I made a simple but crude runFile converter (that is available [here](https://github.com/kreatolinux/src/raw/refs/heads/master/scripts/run2to3.nim) using LLMs. It kinda sucks but it made my job a lot easier so I am not complaining.

As writing of this post the packages are still not fully migrated yet. Expect that in the few weeks.

And that is it (for now!)

## Thank you

With runFile v3, Kreato Linux source tree also surpassed 2000 commits which is a huge milestone for me! Expect more things to come and thank you to all of my supporters, especially those that I took opinions from to make runFile v3 as readable and easy to use as possible!

The future seems bright with more refactors and features (such as a linter, REPL etc.) on the way!
