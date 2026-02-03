<div align="center">
  <h1 style="color: #D20000; background-color: #000000; padding: 15px; border: 2px solid #D20000; border-radius: 10px; display: inline-block;">Zabbix CLI | Enterprise Observability</h1>
  <p style="font-size: 1.2em; color: #D20000;">Performance. Observabilidade. AutomaÃ§Ã£o. Escalabilidade.</p>
  <p style="color: #D20000; font-style: italic;">"Redefinindo a interaÃ§Ã£o com o Zabbix atravÃ©s de uma arquitetura nativa em Go de alta performance."</p>

  [![License](https://img.shields.io/badge/License-MIT-D20000?style=for-the-badge&logo=mit&logoColor=black&labelColor=black)](LICENSE)
  [![Go](https://img.shields.io/badge/Engine-Go_1.23-D20000?style=for-the-badge&logo=go&logoColor=black&labelColor=black)](https://go.dev/)
</div>

---

## **ZABBIX-DNA**
O **ZABBIX-DNA** Ã© o propÃ³sito central deste ecossistema: uma plataforma CLI de classe enterprise escrita 100% em Go, focada em performance extrema e observabilidade moderna.

### **PropÃ³sito EstratÃ©gico**
- **EliminaÃ§Ã£o de RuÃ­do**: Transformar milhares de itens em traces acionÃ¡veis.
- **Performance Nativa**: Arquitetura em Go para processamento massivo de dados sem overhead.
- **Conformidade Enterprise**: AutomaÃ§Ã£o de backups, migraÃ§Ãµes e auditoria.
- **Observabilidade First**: IntegraÃ§Ã£o nativa com o stack OpenTelemetry (OTLP).

---

## **Guia de InÃ­cio RÃ¡pido (Linux Only)**

### **CompilaÃ§Ã£o do BinÃ¡rio Nativo**
```bash
go build -o zabbix-dna ./cmd/dna
sudo mv zabbix-dna /usr/local/bin/
```

---

## **ConfiguraÃ§Ã£o**
O sistema utiliza o arquivo de configuraÃ§Ã£o `zabbix-dna.toml`. 

### **DefiniÃ§Ãµes Estruturais (zabbix-dna.toml)**
```toml
[zabbix]
url = "https://zabbix.exemplo.com/api_jsonrpc.php"
token = "seu_token_aqui"
timeout = 30

[otlp]
endpoint = "http://otel-collector:4318"
protocol = "http"
service_name = "zabbix-dna"
```

---

## **OperaÃ§Ãµes**

### **Observabilidade AvanÃ§ada**
ExportaÃ§Ã£o de mÃ©tricas estruturadas para o stack de monitoramento:
```bash
zabbix-dna metrics --endpoint http://localhost:4318 --interval 60s
```

Mapeamento de eventos Zabbix como traces OTLP:
```bash
zabbix-dna traces --endpoint http://localhost:4318 --batch-size 100
```

### **AdministraÃ§Ã£o de Plataforma**
Listagem de recursos:
```bash
zabbix-dna host
zabbix-dna template
zabbix-dna proxy
```

ExecuÃ§Ã£o de backups de configuraÃ§Ã£o:
```bash
zabbix-dna backup
```

---

## **Filosofia**
- *"Se nÃ£o Ã© monitorado, nÃ£o existe."*
- *"Se Ã© repetitivo, deve ser automatizado."*
- *"Infra nÃ£o Ã© arte. Es engenharia."*

---

## **Releases e Pacotes**

### **Release Atual**
**v1.0.0**: LanÃ§amento inicial da plataforma ZABBIX-DNA 100% Go.

### **DistribuiÃ§Ã£o**
**No packages published**: No momento, a distribuiÃ§Ã£o Ã© realizada exclusivamente via cÃ³digo fonte para garantir a integridade enterprise.

---

## **Mantenedor**
**Ranlens Denck** Ã© Analista de Infraestrutura de TI focado em seguranÃ§a, automaÃ§Ã£o e monitoramento. Sua missÃ£o Ã© eliminar o trabalho manual e aumentar a previsibilidade atravÃ©s de uma observabilidade proativa.

<div align="center">
  <p style="color: #D20000; background-color: #000000; padding: 15px; border-top: 2px solid #D20000; border-radius: 0 0 10px 10px;">
    <b>Ranlens Denck | Observabilidade First</b><br>
    Construindo sistemas que nÃ£o acordam pessoas de madrugada.<br>
    <a href="https://www.linkedin.com/in/ranlensdenck/" style="color: #D20000;">LinkedIn</a> |
    <a href="mailto:ranlens.denck@protonmail.com" style="color: #D20000;">Email</a>
  </p>
</div>
