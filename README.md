Histórico de Negociações B3

Este projeto tem como objetivo processar arquivos de negociações da B3 dos últimos dias úteis, armazenar os dados em um banco de dados PostgreSQL e permitir consultas agregadas por meio de uma interface de linha de comando (CLI).

---

## 🚀 Funcionalidades

- Processamento eficiente de arquivos CSV da B3
- Ingestão em paralelo com goroutines
- Inserção em lote com `pgx.CopyFrom` (alta performance)
- Consulta agregada por `ticker` e `data`
- Interface via CLI usando `urfave/cli/v2`
- Logging estruturado com `zap`

---

## 🧱 Estrutura do Projeto

```
.
├── cmd/worker/              # CLI principal
├── internal/
│   ├── domain/              # Lógica de negócio (parser, processor, finder)
│   ├── app/                 # Injeção e orquestração com Fx
├── pkg/
│   ├── repository/          # Repositório de acesso ao banco de dados
│   ├── utils/               # Funções auxiliares (env, etc)
├── migrations/sql/          # Scripts SQL para criação de tabelas
├── mocks/                   # Mocks gerados automaticamente
├── README.md                # Este arquivo
```

---

## ⚙️ Requisitos

- [Go 1.20+](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [sqlc](https://sqlc.dev/) (para gerar código de acesso ao banco)
- Docker (opcional, para ambiente local)

---

## 📦 Instalação

### Clone o repositório

```bash
git clone https://github.com/seu-usuario/nome-do-repo.git
cd nome-do-repo
```

### Gere os arquivos com `sqlc`

```bash
sqlc generate
```

### Configure o banco de dados

Crie um banco PostgreSQL e execute o script de migração inicial:

```bash
psql -U seu_usuario -d seu_banco < migrations/sql/V1__create_trade_table.sql
```

---

## 🛠️ Variáveis de Ambiente

| Variável             | Padrão           | Descrição                                      |
|----------------------|------------------|-----------------------------------------------|
| `DATABASE_URL`       | -                | String de conexão com o PostgreSQL            |
| `BATCH_SIZE`         | `1000`           | Quantidade de registros por lote              |
| `NUM_WORKERS`        | `2 * CPUs`       | Número de goroutines para persistência        |
| `MAX_CHANNEL_BUFFER` | `10000`          | Buffer do canal entre parser e persistência   |
| `TICKER_SECONDS`     | `2`              | Intervalo de flush em segundos                |
| `CSV_DELIMITER`      | `;`              | Separador de campos no CSV                    |
| `SKIP_HEADER`        | `true`           | Ignorar o cabeçalho do CSV                    |

---

## 🧪 Como Rodar

### Processar arquivos CSV:

```bash
go run cmd/worker/main.go process --file ./dados/NEGOCIOS20240601.txt --file ./dados/NEGOCIOS20240602.txt
```

### Consultar agregados:

```bash
go run cmd/worker/main.go query --ticker PETR4 --date 2024-06-01
```

Saída esperada:

```text
Summary Trade Result:
 Ticker: PETR4
 Max Range Value: 32.50
 Max Daily Volume: 438000
```

---

## 🐳 Rodando com Docker

Você pode usar um `docker-compose` como este:

```yaml
version: "3.8"
services:
  db:
    image: postgres:15
    restart: always
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: b3data
    volumes:
      - ./migrations/sql:/docker-entrypoint-initdb.d
```

---

## 🧪 Testes

Rode testes unitários:

```bash
go test ./internal/domain/trade/... ./pkg/repository/... -v
```

## 🧑‍💻 Autor

Desenvolvido por [Lucas de Leão] - [LinkedIn](https://www.linkedin.com/in/lucas-de-le%C3%A3o-999a73156/)  
Desafio técnico realizado para [Nome da Empresa]

---