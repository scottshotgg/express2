# TODO:

---- parser/semantic/transpiler ----

- object - **normal priority**
  this will be a struct essentially on the backend with an integrated map and a constructor function
  objects can't be created at runtime; dynamic json would expand into a map<var:var>
  therefore, just have the compiler use void\* to pass around the object and then have it deref it before it
  their code starts so that it seems like you can pass an object
  --- HOLD OFF ON THIS FOR NOW ---

- fix forof - **normal priority**
- test/fix forstd - **normal priority**

---- parser? ----

- type annotations - **normal priority**
- fix map/struct/object type tokens; they shouldn't be keywords
- allow for named blocks
- adding printf/println from stdlib - **high priority**
- Add typing ability on enums - **not high priority**
- Need to add he enum'ed types to the type scope - **normal priority**
- Need to add all flattened values to scope - **normal priority**
- include and imports need to be looked at more - low priority until imports are needed
- rewrite ParseArrayType to fix multi dimensional arrays and support the direct array types - **high priority**
  - generate types using std::array, will be much easier to generate since the type is completely separate
  - https://www.quora.com/Can-std-array-type-size-be-used-for-a-multidimensional-array-like-myarray-L-M-N-in-c++11
- fix no-main ability for "scripting"

---- semantic ----

- make this
- libs need to be better
- type annotations - **not high priority**
- type checker that will run BEFORE the transpiler and apply types to EVERYTHING - **high priority**
  - this is where block types should be figured out, ident types resolved, etc
- move ALL type checking to this stage
- automatic derefing like go but using the arrow operator behind the scenes
- integrated, automatic map with object
- val - **not high priority**
- let - **not high priority**
- add types to blocks, will make it easier to transpile; use these types to check the blocks for invalid statements - **normal priority**
  - i.e, non-kv pairs in maps, etc. Can also use it for type inference of maps, etc

---- flattener ----

- make the flattener recurse
- make an official stage for it

---- transpiler ----

- need to do vectors - **not high priority**
- char[] should support string literals - **not high priority**

---- infrastructure ----

- make a test folder and put individual tests
- rig up test suite for lex/parse/semantic/flatten/transpile like we had in Express 1
- make an install script
- make a build-and-test-all script
- build a docker image to do the testing automatically
- write documentation; include the type lineage map
- fill out tests
- hook up a build pipeline (CI) to test everything on push
- use circleCI and deploy to there?
- make the cmd folder and implement the commands
- purge current logging
- add more logging
- reorganize code
- organize project folder; move everything into the `stages` folder
- organize the builder folder; try splitting out some stuff
- experiment with custom clang/LLVM compilation
- research custom pass in LLVM
- research generating LLVM code directly

---- finished ----

- var
- access operator
- selection operator
- de/ref operator
