#!/bin/bash
#
# Testing functions from crebas_misc.sql
#
# Usage:
#  bash test_misc.sh
#  TAG=- SRC=404 bash test_misc.sh
#  TAG=- SRC=arr bash test_misc.sh
#
# Warning: элементы массива должны идти подряд, чтобы корректно получить JSON.
# Поддерживаются только массивы из 2х элементов,
# для массива из одного элемента в конфигу должен быть добавлен он же еще раз с пустым значением.
# Такая пара должна идти последней

# strict mode http://redsymbol.net/articles/unofficial-bash-strict-mode/

# result destination
[[ "$TAG" ]] || TAG="test_misc"

# call source
[[ "$SRC" ]] || SRC="-"

# -o pipefail stops on 404 test
set -eu

HOST=http://localhost:8081/rpc

# NOTE: real GET queries comes without '&name=' suffix
[[ "$AP" ]] || AP=""

Q=$(cat <<EOF
index
index?${AP}lang=ru
func_args?${AP}code=index
func_args?${AP}code=func_args
func_args?${AP}code=func_result
func_result?${AP}code=index
func_result?${AP}code=func_args
func_result?${AP}code=func_result
echo_not_found?${AP}name=test
test_error
EOF
)

# ------------------------------------------------------------------------------

params_json() {
  local p=$1
  set -f                      # avoid globbing (expansion of *).
  array=(${p//&/ })
  echo -n '{'
  local key=""
  local out=""
  local was=""
  local pre_key=""
  for e in "${array[@]}" ; do
    key=${e%=*}

    if [[ "$key" == "$pre_key" ]] ; then
      # array
      val1=${e#*=}
      if [[ "$val1" ]] ; then
        out="\"$key\":[\"$val\",\"$val1\"]"
      else
        # empty elem as array flag
        out="\"$key\":[\"$val\"]"
      fi
    else
      [[ "$was" ]] && echo -n ","
      [[ "$out" ]] && echo -n $out && was=1
      val=${e#*=}
      out="\"$key\":\"$val\""
    fi
    pre_key=$key
  done
  [[ "$was" ]] && echo -n ","
  [[ "$out" ]] && echo -n $out
  echo -n '}'
}

# ------------------------------------------------------------------------------

call() {

  if [[ "$1" == "-" ]] ; then
    shift
  else
    echo $@ && cr
  fi
  local header=1
  local nojson=""
  local s=""

  (stdbuf -o0 $@ | grep -vE "^Date: " | tr -d "\r" |
  while IFS= read -r s
  do
    [[ "$s" ]] || { header="" && cr ; continue ; }
    if [[ "$header" || "$nojson" ]] ; then
      echo $s
      [[ ${s#Content-Type: text/plain} != "$s" ]] && nojson=1
    elif [[ "$s" ]] ; then
      pre && prej &&  echo "$s"
### | jq '.'
    #    [[ "$s" ]] &&  echo "$s"
    fi
    echo -n ""
  done)
  # echo ">>>>" $?
}

cr() {
  echo ""
}

pre() {
  echo '```'
}

prej() {
  echo '```json'
}

# ------------------------------------------------------------------------------

process() {
  local q=$1

  local method=${q%\?*}
  local params=${q#*\?}

  [[ "$method" == "$q" ]] && params=""
  [[ "$method_pre" == "$method" ]] || { echo "## $method" && cr ; }
  method_pre=$method

  # remove name= suffix
  local pg=${params%&*=}

  local pj="{}"
  [[ "$pg" ]] && pj=$(params_json $params)


  [[ "$pg" ]] && echo "### Arguments: $pg" && cr
  echo "#### GET" && cr && pre

  local uri=""
  local arg1="{}"
  local arg2=""

  [[ "$pg" ]] && uri="?$pg" && arg2=",\"params\":$pj"
  call curl -gis "$HOST/$method$uri"

  pre && cr && echo "#### Postgrest" && cr && pre
  call curl -is -d "$pj" -H "Content-type: application/json" "$HOST/$method"

  pre && cr && echo "#### JSON-RPC 2.0" && cr && pre
  local d='{"jsonrpc":"2.0","id":1,"method":"'$method'"'$arg2'}'
  echo "D='$d'"
  echo 'curl -is -d "$D" -H "Content-type: application/json"' "$HOST/" && cr
  call - curl -is -d "$d" -H "Content-type: application/json" "$HOST/"

  pre && cr
}

toc() {
  local q=$1
  local method=${q%\?*}

  [[ "$method_pre" == "$method" ]] || echo "* [$method](#$method)"
  method_pre=$method
}

# ------------------------------------------------------------------------------

if [[ "$SRC" != "-" ]] ; then
  Q=$(cat test_${SRC}.dat)
fi
ERR=""
OUT=""
if [[ "$TAG" != "-" ]] ; then
  ERR="${TAG}_$$.diff"
  [ -f $ERR ] && rm $ERR
  OUT="${TAG}_$$.md"
fi

if [[ "$OUT" ]] ; then
  echo "# crebas_misc.sql testing results" > $OUT
  cr >> $OUT
  echo "Fetching data from server..."

  method_pre=""
  #  for q in $Q ; do
  while IFS= read -r line ; do
    # Skip comments
    q=${line%%#*} # remove endline comments
    [ -n "${q##+([[:space:]])}" ] || continue # ignore line if contains only spaces

    toc $q >> $OUT
  done <<< "$Q"
  cr && cr >> $OUT

fi

method_pre=""
while IFS= read -r line ; do

  # Skip comments
  q=${line%%#*} # remove endline comments
  [ -n "${q##+([[:space:]])}" ] || continue # ignore line if contains only spaces

  if [[ "$OUT" ]] ; then
    echo $q
    process $q >> $OUT
  else
    process $q
  fi
done <<< "$Q"
echo "------------------"
cr
[[ "$OUT" ]] || { echo "no compare" ;  exit ; }
diff -c -I "^X-Elapsed" test_misc.md $OUT > $ERR || echo "***************************************************"

if [ -s $ERR ] ; then
  cat $ERR
  echo "ERRORS found (see $ERR)"
else
  echo "OK"
  rm $ERR
  rm $OUT
fi
