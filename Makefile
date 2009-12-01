#
#  Makefile
#  sed
#
# Copyright (c) 2009 Geoffrey Clements
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
#

# prerequisite: GOROOT and GOARCH must be defined

# defines $(GC) (compiler), $(LD) (linker) and $(O) (architecture)
include $(GOROOT)/src/Make.$(GOARCH)

# name of the package (library) being built
TARG=sed

# source files in package
GOFILES=\
	sed.go \
	cmd.go \
	s_cmd.go \
	d_cmd.go \
	n_cmd.go \
	p_cmd.go \
	q_cmd.go

# test files for this package
GOTESTFILES=\
	sed_test.go

# build "main" executable
sed: package
	$(GC) -I_obj main.go
	$(LD) -L_obj -o $@ main.$O
	@echo "Done. Executable is: $@"

clean:
	rm -rf *.[$(OS)o] *.a [$(OS)].out _obj _test _testmain.go sed

package: _obj/$(TARG).a

fmt:
	gofmt -w -tabwidth=2 -spaces=true main.go
	gofmt -w -tabwidth=2 -spaces=true sed.go
	gofmt -w -tabwidth=2 -spaces=true cmd.go
	gofmt -w -tabwidth=2 -spaces=true d_cmd.go
	gofmt -w -tabwidth=2 -spaces=true n_cmd.go
	gofmt -w -tabwidth=2 -spaces=true p_cmd.go
	gofmt -w -tabwidth=2 -spaces=true q_cmd.go
	gofmt -w -tabwidth=2 -spaces=true s_cmd.go
	gofmt -w -tabwidth=2 -spaces=true sed_test.go


# create a Go package file (.a)
_obj/$(TARG).a: _go_.$O
	@mkdir -p _obj/$(dir)
	rm -f _obj/$(TARG).a
	gopack grc $@ _go_.$O

# create Go package for for tests
_test/$(TARG).a: _gotest_.$O
	@mkdir -p _test/$(dir)
	rm -f _test/$(TARG).a
	gopack grc $@ _gotest_.$O

# compile
_go_.$O: $(GOFILES)
	$(GC) -o $@ $(GOFILES)

# compile tests
_gotest_.$O: $(GOFILES) $(GOTESTFILES)
	$(GC) -o $@ $(GOFILES) $(GOTESTFILES)


# targets needed by gotest

importpath:
	@echo $(TARG)

testpackage: _test/$(TARG).a

testpackage-clean:
	rm -f _test/$(TARG).a _gotest_.$O
