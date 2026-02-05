# Documentação: Agents Zabbix Dimo & Integração SaltStack

Esta documentação detalha a implementação e o uso dos Agents Zabbix personalizados "Dimo" e como automatizar sua implantação usando a CLI `zabbix-dna` integrada ao SaltStack.

## 1. Agents Zabbix Dimo (Custom SDK)

Os Agents Dimo foram desenvolvidos utilizando o SDK oficial da Zabbix (`golang.zabbix.com/sdk`), permitindo uma coleta de métricas de alta performance e descoberta dinâmica via LLD (Low-Level Discovery).

### Métricas Implementadas:
- `dimo.version`: Versão do Agent Dimo.
- `dimo.system.info`: Informações detalhadas do sistema (OS, Arch, CPUs, Goroutines).
- `dimo.discovery.disks`: Descoberta automática de discos/partições (LLD).
- `dimo.discovery.network`: Descoberta automática de interfaces de rede (LLD).
- `dimo.discovery.services`: Descoberta de serviços críticos (LLD).
- `dimo.proc.count`: Contagem total de processos.
- `dimo.mem.usage`: Métricas detalhadas de uso de memória.

### Locais de Instalação:
- **Linux**: `/opt/dimo/`
- **Windows**: `C:\Dimo\`

---

## 2. Integração com SaltStack

A CLI `zabbix-dna` possui um comando dedicado para gerenciar a implantação desses agents em escala, garantindo confiabilidade total através de logs detalhados e validação de passos.

### Comando de Implantação:
```bash
zabbix-dna salt deploy_agent --target "minion-id" --os "linux"
```

### O que o comando faz:
1. **Criação de Diretórios**: Garante que `/opt/dimo/` ou `C:\Dimo\` existam.
2. **Download do Binário**: Baixa a versão correta do agent para a arquitetura alvo.
3. **Permissões**: Configura as permissões necessárias (no Linux).
4. **Reinicialização**: Reinicia o serviço `zabbix-agent2` para carregar o novo plugin Dimo.

---

## 3. Guia de Instalação Remota via CLI

Para instalar o agent em todos os servidores Linux de um grupo:

```bash
zabbix-dna salt deploy_agent --target "L*" --type "glob" --os "linux"
```

Para instalar em servidores Windows:

```bash
zabbix-dna salt deploy_agent --target "W*" --type "glob" --os "windows"
```

### Verificação:
Após a instalação, você pode verificar se o agent está respondendo com:

```bash
zabbix-dna salt run "cmd.run 'zabbix_get -s 127.0.0.1 -k dimo.version'" --target "minion-id"
```

---

## 4. Ponto Crítico: Confiabilidade
A integração foi desenhada para ser transacional. Se qualquer passo (criação de pasta, download ou reinício) falhar, o processo é interrompido imediatamente com um erro claro, evitando configurações parciais ou inconsistentes no parque de servidores.
