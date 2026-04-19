#!/usr/bin/env bash
# check_plan_d.sh — verifies Plan D (immutability by default) end-to-end
# Run from the repo root with: EXPRPATH=$(pwd) ./scripts/check_plan_d.sh

set -euo pipefail
PASS=0
FAIL=0
NOTES=()

ok()   { echo "  [PASS] $1"; PASS=$((PASS+1)); }
fail() { echo "  [FAIL] $1"; FAIL=$((FAIL+1)); NOTES+=("FAIL: $1"); }
note() { echo "  [NOTE] $1"; NOTES+=("NOTE: $1"); }

header() { echo ""; echo "── $1 ──"; }

# ── 1. Build ──────────────────────────────────────────────────────────────────
header "1. go build ./..."
go build ./... 2>&1 && ok "go build ./..." || fail "go build ./..."

# ── 2. Unit tests ─────────────────────────────────────────────────────────────
header "2. Unit tests (builder, transpiler, tree_flattener)"
go test ./builder/... ./transpiler/... ./tree_flattener/... -count=1 -race 2>&1 | tail -4
[[ $(go test ./builder/... ./transpiler/... ./tree_flattener/... -count=1 2>&1 | grep -c "^ok") -eq 3 ]] \
  && ok "All unit test packages pass" || fail "Some unit tests failed"

# ── 3. New tests present and passing ──────────────────────────────────────────
header "3. New tests"
for t in TestVarTypedDeclaration TestVarTypedDeclUninitialized \
          TestTranspileImmutableDecl TestTranspileVarTypedDecl TestTranspileLetConst; do
  out=$(go test ./... -run "^${t}$" -count=1 2>&1)
  if echo "$out" | grep -q "^ok"; then
    ok "$t"
  else
    fail "$t"
  fi
done

# ── 4. Integration tests ───────────────────────────────────────────────────────
header "4. Integration tests (compiler/...)"
if [[ -z "${EXPRPATH:-}" ]]; then
  note "EXPRPATH not set — skipping integration tests"
else
  EXPRPATH=$(pwd) go test ./compiler/... -count=1 2>&1 | tail -2
  EXPRPATH=$(pwd) go test ./compiler/... -count=1 2>&1 | grep -q "^ok" \
    && ok "All integration tests pass" || fail "Integration tests failed"
fi

# ── 5. Generated C++ spot checks ──────────────────────────────────────────────
header "5. Generated C++ spot checks"
EXAMPLES=(
  "03-arithmetic.expr"
  "05-functions.expr"
  "07-structs.expr"
  "09-pointers.expr"
  "10-complex.expr"
  "21-zero-init.expr"
  "22-combinations.expr"
  "23-map-every-type.expr"
)

for f in "${EXAMPLES[@]}"; do
  src="docs/examples/$f"
  cpp="docs/examples/${f%.expr}.cpp"
  rm -f "${cpp%.cpp}" "$cpp"
  if EXPRPATH=$(pwd) go run main.go build "$src" 2>/dev/null; then
    ok "compiles: $f"
  else
    fail "compile error: $f"
    continue
  fi

  case "$f" in
  03-arithmetic.expr)
    grep -q "^  int counter = 0;" "$cpp" && ok "  counter is mutable (no const)" \
      || fail "  counter should be mutable in $cpp"
    ;;
  05-functions.expr)
    # Function params must not have const
    grep -q "void greet(std::string name)" "$cpp" && ok "  param: no const on string name" \
      || fail "  param got const in $cpp"
    # Immutable local
    grep -q "const int sum" "$cpp" && ok "  local: const int sum" \
      || fail "  'const int sum' missing in $cpp"
    ;;
  07-structs.expr)
    grep -q "^  Person alice = {" "$cpp" && ok "  alice is mutable (var Person)" \
      || fail "  alice should be mutable in $cpp"
    grep -q "^  const Person bob = {" "$cpp" && ok "  bob is immutable (no var)" \
      || fail "  bob should be const in $cpp"
    grep -q "^struct Person {" "$cpp" && ok "  struct fields have no const prefix" \
      || fail "  struct declaration wrong in $cpp"
    grep "  std::string name" "$cpp" | grep -qv "const" && ok "  field name: no const" \
      || fail "  struct field 'name' got const in $cpp"
    ;;
  09-pointers.expr)
    grep -q "^  int value = 42;" "$cpp" && ok "  value is mutable" \
      || fail "  value should be mutable in $cpp"
    grep -q "^  int x = 5;" "$cpp" && ok "  x is mutable" \
      || fail "  x should be mutable in $cpp"
    grep -q "^  int y = 10;" "$cpp" && ok "  y is mutable" \
      || fail "  y should be mutable in $cpp"
    grep -q "const int temp = " "$cpp" && ok "  temp is immutable (never modified)" \
      || fail "  temp should be const in $cpp"
    # Pointers should never get const prefix
    grep "int \*ptr" "$cpp" | grep -qv "const int \*ptr" && ok "  ptr: no const binding" \
      || fail "  pointer ptr got const in $cpp"
    ;;
  10-complex.expr)
    grep -q "^  Student alice = {" "$cpp" && ok "  alice is mutable (var Student)" \
      || fail "  alice should be mutable in $cpp"
    grep -q "^  int total = 0;" "$cpp" && ok "  total is mutable (var int)" \
      || fail "  total should be mutable in $cpp"
    # Struct fields
    grep "  std::string name" "$cpp" | grep -qv "const" && ok "  Student.name: no const" \
      || fail "  Student field 'name' got const in $cpp"
    ;;
  21-zero-init.expr)
    grep -q "^  int i = 0;" "$cpp" && ok "  i is mutable (var int i)" \
      || fail "  i should be mutable in $cpp"
    grep -q "^  const float f = 0;" "$cpp" && ok "  f is immutable (float f)" \
      || fail "  f should be const in $cpp"
    grep -q "^  const bool b = false;" "$cpp" && ok "  b is immutable (bool b)" \
      || fail "  b should be const in $cpp"
    grep -q "^  std::string s;" "$cpp" && ok "  s is mutable (var string s)" \
      || fail "  s should be mutable in $cpp"
    grep -q "^  Outer o = {};" "$cpp" && ok "  o is mutable (var Outer o)" \
      || fail "  o should be mutable in $cpp"
    ;;
  22-combinations.expr)
    grep -q "^  std::vector<int> iv;" "$cpp" && ok "  iv mutable (var int[])" \
      || fail "  iv should be mutable in $cpp"
    grep -q "^  const std::vector<float> fv;" "$cpp" && ok "  fv immutable (float[])" \
      || fail "  fv should be const in $cpp"
    grep -q "^  std::map<std::string, var> mv;" "$cpp" && ok "  mv mutable (var map)" \
      || fail "  mv should be mutable in $cpp"
    grep -q "^  std::vector<Point> points;" "$cpp" && ok "  points mutable (var Point[])" \
      || fail "  points should be mutable in $cpp"
    grep -q "^  const Point p = {" "$cpp" && ok "  p immutable (Point p, no var)" \
      || fail "  p should be const in $cpp"
    ;;
  23-map-every-type.expr)
    grep -q "^  std::map<std::string, var> m;" "$cpp" && ok "  m mutable (var map)" \
      || fail "  m should be mutable in $cpp"
    grep -q "^  const std::vector<int> nums = " "$cpp" && ok "  nums immutable (int[] read-only)" \
      || fail "  nums should be const in $cpp"
    ;;
  esac
done

# ── 6. Loop-counter mutability ─────────────────────────────────────────────────
header "6. Loop counter mutability (04-control-flow.expr)"
EXPRPATH=$(pwd) go run main.go build docs/examples/04-control-flow.expr 2>/dev/null
CPP04="docs/examples/04-control-flow.cpp"
# _idx counters must NOT be const
if grep "_idx_" "$CPP04" | grep -q "const int _idx_"; then
  fail "Loop _idx counter got 'const' — will break increment"
else
  ok "Loop _idx counters are mutable (no const)"
fi
# Loop element bindings (auto) should be const
grep -q "const auto" "$CPP04" && ok "Loop element bindings are immutable (const auto)" \
  || fail "Loop element bindings should be const auto"

# ── 7. let → const ────────────────────────────────────────────────────────────
header "7. let statement emits const"
RESULT=$(go test ./transpiler/... -run "^TestTranspileLetConst$" -v 2>&1)
echo "$RESULT" | grep -q "PASS" && ok "let x = 5 → const int x = 5;" \
  || fail "let statement not emitting const"

# ── Summary ───────────────────────────────────────────────────────────────────
echo ""
echo "══════════════════════════════════════════"
echo "  PASS: $PASS   FAIL: $FAIL"
echo "══════════════════════════════════════════"
if [[ ${#NOTES[@]} -gt 0 ]]; then
  for n in "${NOTES[@]}"; do echo "  $n"; done
fi
[[ $FAIL -eq 0 ]] && echo "  All checks passed." || echo "  Some checks failed."
