#!/bin/bash
#
# doc_gen.sh v1.0
# Генерация .md с описанием методов заданной схемы
#
# (c) Алексей Коврижкин, 12.08.2016
#
# Use:
# bash doc_gen.sh | sed -e 's|null||g'  > doc_sample.md
#
# TODO:
# * [ ] Если в описании результата имя поля = '-', заменять на имя функции

. .config
# API сервер
#H=http://iac.tender.pro
H=$APP_SITE
HOST=$H/api

AP=""

# Схема API сервера
SCHEMA=public

toc() {
  local mtd=$1
  def=$(echo $json | jq -r ".result | .[] | select(.code==\"$mtd\")") # "
  anno=$(echo $def | jq -r .anno)
#  [[ "$anno" != "null" ]] || return

  echo "* [$mtd](#$mtd)"
}

describe() {
  local mtd=$1
  echo "Describe method $mtd..." >&2

  def=$(echo $json | jq -r ".result | .[] | select(.code==\"$mtd\")") # "
  anno=$(echo $def | jq -r .anno)
  exam=$(echo $def | jq -r .sample)
#  [[ "$anno" != "null" ]] || return
cat <<EOF

## $mtd

$anno

### Аргументы

Имя | Тип | По умолчанию | Описание
----|-----|--------------|---------
EOF

  local hdr
  while read a ; do
    #[[ "$hdr" ]] || hdr="|Имя|Тип| По умолчанию |Разрешен NULL| Описание|" && echo $hdr
    echo $a
  done < <(curl -gs "$HOST/pg_func_args_ext?${AP}code=$mtd" | jq -r '.result | .[] | " \(.name) | \(.type) | \(.["def"]) | \(.anno)"')
cat <<EOF

### Результат

Имя | Тип | Описание
----|-----|---------
EOF

  while read a ; do
    echo $a
  done < <(curl -gs "$HOST/pg_func_result_ext?${AP}code=$mtd" | jq -r '.result | .[] | " \(.name) | \(.type) | \(.anno) "')

  local result
  if [[ "$exam" != "null" ]] ; then

cat <<EOF

### Пример вызова

\`\`\`
H=$H
Q='$exam'
curl -gsd "\$Q" -H "Content-type: application/json" \$H/rpc/$mtd | jq .[0:2]
EOF
    result=$(curl -gsd "$exam" -H "Content-type: application/json" $HOST/$mtd | jq .)
cat <<EOF
\`\`\`
\`\`\`json
$result
\`\`\`

EOF
  fi
}

cat <<EOF

# Схема $SCHEMA. Методы API

EOF

json=$(curl -s $HOST/index?${AP}nspname=$SCHEMA)

while read a ; do
  toc $a
done < <(echo $json | jq -r '.result | .[] | .code')

while read a ; do
  describe $a
done < <(echo $json | jq -r '.result | .[] | .code')
