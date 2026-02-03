#!/bin/bash

# Zabbix DNA - Installer Script
# Performance. Observabilidade. Automação. Escalabilidade.

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Iniciando a instalação do Zabbix DNA...${NC}"

# Verificar se é Linux
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo -e "${RED}Erro: Este script suporta apenas sistemas Linux.${NC}"
    exit 1
fi

# Detectar Arquitetura
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        GOARCH="amd64"
        ;;
    aarch64|arm64)
        GOARCH="arm64"
        ;;
    *)
        echo -e "${RED}Erro: Arquitetura $ARCH não suportada.${NC}"
        exit 1
        ;;
esac

# Obter última versão do GitHub
REPO="rsdenck/zabbix-dna"
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo -e "${RED}Erro: Não foi possível obter a última versão do repositório.${NC}"
    exit 1
fi

echo -e "Versão detectada: ${GREEN}$LATEST_TAG${NC}"
echo -e "Arquitetura: ${GREEN}$GOARCH${NC}"

# Download da URL
BINARY_NAME="zabbix-dna-linux-$GOARCH"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$BINARY_NAME"

echo -e "Baixando de: $DOWNLOAD_URL"

# Download temporário
TEMP_DIR=$(mktemp -d)
curl -L -o "$TEMP_DIR/zabbix-dna" "$DOWNLOAD_URL"

# Instalação
echo -e "${GREEN}Instalando em /usr/local/bin/zabbix-dna...${NC}"
chmod +x "$TEMP_DIR/zabbix-dna"
sudo mv "$TEMP_DIR/zabbix-dna" /usr/local/bin/zabbix-dna

# Limpeza
rm -rf "$TEMP_DIR"

echo -e "${GREEN}Instalação concluída com sucesso!${NC}"
echo -e "Execute 'zabbix-dna --help' para começar."
