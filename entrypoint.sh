#!/bin/sh

# Exibir o diretório atual
echo "Current Directory:"
pwd

# Listar arquivos no diretório atual
echo "Files in Current Directory:"
ls -la

# Verificar se o arquivo .env está presente no diretório /app
echo "Files in /app Directory:"
ls -la /app

# Verificar se o diretório de migrações está presente
echo "Files in /app/migrations Directory:"
ls -la /app/migrations

# Espera de alguns segundos para visualizar os logs no console
sleep 10

# Executa a aplicação e, em seguida, as migrações
./ordersystem && migrate -path /app/migrations -database "mysql://root:root@tcp(mysql:3306)/orders" up

# Se algo falhar, manter o contêiner ativo para depuração
if [ $? -ne 0 ]; then
  echo "An error occurred. Keeping the container alive for debugging."
  tail -f /dev/null
fi
