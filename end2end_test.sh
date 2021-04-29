#!/usr/bin/env bash
green="\033[32m"
red="\033[31m"
bold="\033[1m"
reset="\033[0m"

curr_dir=$PWD

log_success() {
  echo -e "${green}$@${reset}"
}

log_err() {
  echo -e "${red}$@${reset}"
}

log_bold() {
  echo -e "${bold}$@${reset}"
}

assert_eq() {
  if [ "$2" != "$3" ]; then
    echo -e "  ${red}$1: $2 != $3${reset}"
    log_bold Tearing down env...
    docker-compose down
    exit -1
  fi
}

assert_not_null() {
  if [ "$2" == "" ]; then
    echo -e "  ${red}$1: $2 is empty${reset}"
    log_bold Tearing down env...
    docker-compose down
    exit -1
  fi
}

ensure_tools() {
  set +e
  for tool_check in "${@}"
  do
    # log_bold "  ..testing tool $tool_check"
    $tool_check &> /dev/null
    assert_eq "$tool_check exit code" $? 0
  done
  set -e
}

run_tests() {
  # test: POST articles
  article1_id_created=`curl -X POST localhost:8080/articles -H "Content-Type: application/json" --data '{"title": "t1", "date": "2001-02-03", "body": "b1", "tags": ["tt1", "tt2"]}' | jq -r '.id'`
  article2_id_created=`curl -X POST localhost:8080/articles -H "Content-Type: application/json" --data '{"title": "t2", "date": "2001-02-03", "body": "b2", "tags": ["tt2", "tt3"]}' | jq -r '.id'`
  article3_id_created=`curl -X POST localhost:8080/articles -H "Content-Type: application/json" --data '{"title": "t3", "date": "2001-02-04", "body": "b3", "tags": ["tt3", "tt4"]}' | jq -r '.id'`
  log_bold Created article ids: $article1_id_created, $article3_id_created, $article3_id_created
  assert_not_null "article1_id is empty" $article1_id_created
  assert_not_null "article2_id is empty" $article2_id_created
  assert_not_null "article3_id is empty" $article3_id_created

  # test: GET articles
  article1_id_fetched=`curl -X GET localhost:8080/articles/$article1_id_created | jq -r '.id'`
  article2_id_fetched=`curl -X GET localhost:8080/articles/$article2_id_created | jq -r '.id'`
  article3_id_fetched=`curl -X GET localhost:8080/articles/$article3_id_created | jq -r '.id'`
  log_bold Re-fetched article ids: $article1_id_fetched, $article3_id_fetched, $article3_id_fetched
  assert_eq "fetched alert1_id mismatch" $article1_id_created $article1_id_fetched
  assert_eq "fetched alert2_id mismatch" $article2_id_created $article2_id_fetched
  assert_eq "fetched alert3_id mismatch" $article3_id_created $article3_id_fetched

  # test: GET tags
  related_tags_fetched=`curl -X GET localhost:8080/tags/tt2/20010203 | jq -r -c '.related_tagss'`
  log_bold Fetched tags: $related_tags_fetched
  assert_eq "fetched tags mismatch" $related_tags_fetched '["tt1","tt3"]'
}

ensure_tools "curl --version" "jq --version" "docker --version" "docker-compose --version" 

log_bold Starting up env...
docker-compose up -d
sleep 5

# test: POST articles
article1_id_created=`curl -X POST localhost:8080/articles -H "Content-Type: application/json" --data '{"title": "t1", "date": "2001-02-03", "body": "b1", "tags": ["tt1", "tt2"]}' | jq -r '.id'`
article2_id_created=`curl -X POST localhost:8080/articles -H "Content-Type: application/json" --data '{"title": "t2", "date": "2001-02-03", "body": "b2", "tags": ["tt2", "tt3"]}' | jq -r '.id'`
article3_id_created=`curl -X POST localhost:8080/articles -H "Content-Type: application/json" --data '{"title": "t3", "date": "2001-02-04", "body": "b3", "tags": ["tt3", "tt4"]}' | jq -r '.id'`
log_bold Created article ids: $article1_id_created, $article3_id_created, $article3_id_created
assert_not_null "article1_id is empty" $article1_id_created
assert_not_null "article2_id is empty" $article2_id_created
assert_not_null "article3_id is empty" $article3_id_created

# test: GET articles
article1_id_fetched=`curl -X GET localhost:8080/articles/$article1_id_created | jq -r '.id'`
article2_id_fetched=`curl -X GET localhost:8080/articles/$article2_id_created | jq -r '.id'`
article3_id_fetched=`curl -X GET localhost:8080/articles/$article3_id_created | jq -r '.id'`
log_bold Re-fetched article ids: $article1_id_fetched, $article3_id_fetched, $article3_id_fetched
assert_eq "fetched alert1_id mismatch" $article1_id_created $article1_id_fetched
assert_eq "fetched alert2_id mismatch" $article2_id_created $article2_id_fetched
assert_eq "fetched alert3_id mismatch" $article3_id_created $article3_id_fetched

# test: GET tags
related_tags_fetched=`curl -X GET localhost:8080/tags/tt2/20010203 | jq -r -c '.related_tags'`
log_bold Fetched tags: $related_tags_fetched
assert_eq "fetched tags mismatch" $related_tags_fetched '["tt1","tt3"]'

log_bold Tearing down env...
docker-compose down