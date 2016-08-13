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

HOST=http://localhost:8081/api

# NOTE: real GET queries comes without '&name=' suffix

Q=$(cat <<EOF
pg_func_args?code=public.dbsize
pg_func_args?code=public.pg_func_args
pg_func_result?code=public.pg_func_args
pg_func_result?code=public.dbsize

dbsize?name=template1
dbsize?name=template1

echo?name=test&id=1
echo?name=test
echo_jsonb?name=test
echo_single?name=test

echo_arr?name=test1&name=test2
echo_arr?name=test1&name=

echo_not_found?name=test
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
      pre && pre &&  echo "$s" | jq '.'
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

# ------------------------------------------------------------------------------

process() {
  local q=$1

  local method=${q%\?*}
  local params=${q#*\?}

  [[ "$method_pre" == "$method" ]] || { echo "## $method" && cr ; }
  method_pre=$method

  local pj=$(params_json $params)

  # remove name= suffix
  local pg=${params%&*=}

  echo "### $pg" && cr
  echo "#### GET" && cr && pre
  call curl -is "$HOST/$method?$pg"

  pre && cr && echo "#### Postgrest" && cr && pre
  call curl -is -d "$pj" -H "Content-type: application/json" "$HOST/$method"

  pre && cr && echo "#### JSON-RPC 2.0" && cr && pre
  local d='{"jsonrpc":"2.0","id":1,"method":"'$method'","params":'$pj'}'
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
  ERR="${TAG}_$$.err"
  [ -f $ERR ] && rm $ERR
  OUT="${TAG}_$$.out"
fi

if [[ "$OUT" ]] ; then
  echo "# crebas_misc.sql testing results" > $OUT
  cr >> $OUT
  echo "Fetching data from server..."
fi

method_pre=""
for q in $Q ; do
  [[ "$q" ]] || continue
  if [[ "$OUT" ]] ; then
    toc $q >> $OUT
  fi
done

[[ "$OUT" ]] && cr && cr >> $OUT

method_pre=""
for q in $Q ; do
  [[ "$q" ]] || continue
  if [[ "$OUT" ]] ; then
    echo $q
    process $q >> $OUT
  else
    process $q
  fi
done
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
