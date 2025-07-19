#!/bin/bash

set -e

# Запуск Keycloak в фоне
/opt/keycloak/bin/kc.sh start-dev &

# Ждём старта Keycloak через kcadm.sh
echo "Waiting for Keycloak to be ready..."
until /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user admin --password admin; do
  sleep 2
done
sleep 5

# Создание realm
/opt/keycloak/bin/kcadm.sh create realms -s realm=bookshop -s enabled=true || true

# Создание клиента
/opt/keycloak/bin/kcadm.sh create clients -r bookshop -s clientId=bookshop-api -s enabled=true -s publicClient=true -s directAccessGrantsEnabled=true -s 'redirectUris=["*"]' || true

# Создание ролей
/opt/keycloak/bin/kcadm.sh create roles -r bookshop -s name=admin || true
/opt/keycloak/bin/kcadm.sh create roles -r bookshop -s name=user || true

# Создание пользователей и установка паролей
declare -A users=(
  [user1]=user1pass
  [user2]=user2pass
  [admin1]=admin1pass
  [admin2]=admin2pass
)

for username in "${!users[@]}"; do
  /opt/keycloak/bin/kcadm.sh create users -r bookshop \
    -s username=$username \
    -s enabled=true \
    -s email=$username@bookshop.local \
    -s emailVerified=true \
    -s firstName=$username \
    -s lastName=User || true
  /opt/keycloak/bin/kcadm.sh set-password -r bookshop --username $username --new-password ${users[$username]}
  USER_ID=$(/opt/keycloak/bin/kcadm.sh get users -r bookshop -q username=$username --fields id --format csv --noquotes)
  [ -n "$USER_ID" ] && /opt/keycloak/bin/kcadm.sh update users/$USER_ID -r bookshop -s 'requiredActions=[]'
done

# Отключить Verify Profile для всех новых пользователей
/opt/keycloak/bin/kcadm.sh update authentication/required-actions/VERIFY_PROFILE -r bookshop -s enabled=false

# Назначение ролей с ожиданием появления пользователя
for username in admin1 admin2; do
  for i in {1..5}; do
    USER_ID=$(/opt/keycloak/bin/kcadm.sh get users -r bookshop -q username=$username --fields id --format csv --noquotes)
    if [ -n "$USER_ID" ]; then
      /opt/keycloak/bin/kcadm.sh add-roles -r bookshop --uid $USER_ID --rolename admin && break
    else
      sleep 2
    fi
  done
done

for username in user1 user2; do
  for i in {1..5}; do
    USER_ID=$(/opt/keycloak/bin/kcadm.sh get users -r bookshop -q username=$username --fields id --format csv --noquotes)
    if [ -n "$USER_ID" ]; then
      /opt/keycloak/bin/kcadm.sh add-roles -r bookshop --uid $USER_ID --rolename user && break
    else
      sleep 2
    fi
  done
done

echo "Keycloak bookshop realm, users, and passwords initialized!"
