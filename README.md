# 📈 B3 Processor

Sistema em Go para processar arquivos de negociações da B3 e consultar dados agregados por ticker.

---

## 🚀 Pré-requisitos

- Go 1.21+
- PostgreSQL rodando e acessível
- Variável de ambiente `DATABASE_URL` configurada com a string de conexão do banco:
  ```
  export DATABASE_URL=postgres://user:password@localhost:5432/dbname
  ```

---

## ⚙️ Configuração

Clone o projeto e instale as dependências:

```bash
make tidy
```

Você pode configurar os parâmetros de execução por variáveis de ambiente (valores padrão entre parênteses):

| Variável              | Descrição                              | Padrão            |
|-----------------------|----------------------------------------|-------------------|
| `BATCH_SIZE`          | Tamanho do lote de inserção            | `1000`            |
| `NUM_WORKERS`         | Quantidade de workers                  | `numCPU * 2`      |
| `MAX_CHANNEL_BUFFER`  | Buffer do canal de trades              | `10000`           |
| `TICKER_SECONDS`      | Tempo de espera para flush (segundos)  | `2`               |
| `CSV_DELIMITER`       | Delimitador do arquivo CSV             | `;`               |
| `SKIP_HEADER`         | Ignora cabeçalho do CSV (`true/false`) | `true`            |

---

## 📂 Estrutura dos Arquivos

Os arquivos de entrada devem estar no formato `.txt` com o seguinte layout de colunas:

```
DataReferencia;CodigoInstrumento;AcaoAtualizacao;PrecoNegocio;QuantidadeNegociada;HoraFechamento;CodigoIdentificadorNegocio;TipoSessaoPregao;DataNegocio;CodigoParticipanteComprador;CodigoParticipanteVendedor
```

Exemplo de nome de arquivo:
```
data/01-06-2025_NEGOCIOSAVISTA.txt
```

---

## 📦 Comandos

### 🔄 Processamento de Arquivos

Processa um ou mais arquivos `.txt` e insere os dados no banco de dados:

```bash
go run main.go process --file data/01-06-2025_NEGOCIOSAVISTA.txt
```

Você pode passar múltiplos arquivos:

```bash
go run main.go process --file data/03-06-2025_NEGOCIOSAVISTA.txt --file data/02-06-2025_NEGOCIOSAVISTA.txt
```

---

### 🔍 Consulta por Ticker

Consulta dados agregados de um ticker:

```bash
go run main.go query --ticker=WINQ25 --date=2025-06-01
```

- `--ticker`: **Obrigatório** – símbolo do ticker.
- `--date`: (opcional) – data inicial no formato `YYYY-MM-DD`.

---

## ✅ Exemplo de Saída

```text
Querying data for ticker: WINQ25, starting from: 2025-06-01 00:00:00 -0300 -03
Summary Trade Result 
 Ticker: WINQ25from
 Max Range Value: 140750.00
 Max Daily Volume: 7
```

---

## 🧪 Testes

Execute os testes com:

```bash
go test ./...
```