#!/bin/bash
#
# Testing functions from crebas_misc.sql
#
# Warning: элементы массива должны идти подряд, чтобы корректно получить JSON.
# Поддерживаются только массивы из 2х элементов,
# для массива из одного элемента в конфигу должен быть добавлен он же еще раз с пустым значением.
# Такая пара должна идти последней

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

Q1=$(cat <<EOF
pg_func_args?code=public.dbsize
echo?name=test&id=1
echo_arr?name=test1&name=test2
echo_arr?name=test1
echo_arr?name=\{test1,t2\}
EOF
)

Q1=$(cat <<EOF
echo?name=test&id=1
echo_arr?name=test1&name=test2
echo_arr?name=test1&name=
EOF
)

Q=$(cat <<EOF
dbsize
dbsize
dbsize
dbsize
dbsize
dbsize
dbsize
dbsize
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
  echo $@ && cr
  $@ | grep -vE "^Date: "
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

  pj=$(params_json $params)

  pg=${params/%&*=/}

  echo "### $pg" && cr
  echo "#### GET" && cr && pre
  call curl -is "$HOST/$method?$pg"

  pre && cr && echo "#### Postgrest" && cr && pre
  call curl -is -d "$pj" -H "Content-type: application/json" "$HOST/$method"

  pre && cr && echo "#### JSON-RPC 2.0" && cr && pre
  d='{"jsonrpc":"2.0","id":1,"method":"'$method'","params":'$pj'}'
  echo "D='$d'"
  echo 'curl -is -d "$D" -H "Content-type: application/json"' "$HOST/" && cr
  curl -is -d "$d" -H "Content-type: application/json" "$HOST/" | grep -vE "^Date: "

  pre && cr
}

# ------------------------------------------------------------------------------

[[ "$TAG" ]] || TAG="test_misc"

if [[ "$TAG" != "-" ]] ; then
 ERR="$TAG_$$.err"
  [ -f $ERR ] && rm $ERR

  OUT="$TAG_$$.out"
fi

if [[ "$OUT" ]] ; then 
  echo "# crebas_misc.sql testing results" > $OUT
  echo "Fetching data from server..."
fi

for q in $Q ; do
  [[ "$q" ]] || continue
  if [[ "$OUT" ]] ; then
    echo $q
    process $q >> $OUT
  else
    process $q
  fi
done
cr
[[ $OUT ]] || exit

diff -c test_misc.md $OUT > $ERR

if [ -s $ERR ] ; then
  cat $ERR
  echo "ERRORS found (see $ERR):"
else
  echo "OK"
  rm $ERR
  rm $OUT
fi
