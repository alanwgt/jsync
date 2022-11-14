jsync
=====

> **Warning** \
> Este programa é fortemente vinculado ao driver do **PostgreSQL**, então tanto a sincronização quanto as migrações só
> funcionam para esse banco. Levante um [issue](https://github.com/alanwgt/jsync/issues/new) para requisitar
> implementações para outros bancos.

O ***jsync*** é um executável que sincroniza os dados de uma imobiliária dentro da [Jetimob](https://www.jetimob.com/)
com um banco de dados *self-hosted*.

São sincronizados: **imóveis**, **condomínios**, **banners** e **corretores**. Desses, apenas os imóveis são atualizados
de forma incremental (a rota de imóveis é a única que aceita o parâmetro para realizar este filtro), isso é, uma variável
é salva localmente com o *timestamp* da última sincronização e posteriormente utilizada nas requisições para baixar
apenas os imóveis modificados após essa data.

<!-- TOC -->
* [Instalação](#instalao)
* [TLDR](#tldr)
* [Uso](#uso)
* [Configurações](#configuraes)
  * [sincronização para uma imobiliária](#sincronizao-para-uma-imobiliria)
  * [sincronização para múltiplas imobiliárias](#sincronizao-para-mltiplas-imobilirias)
  * [configurações obrigatórias](#configuraes-obrigatrias)
  * [Coluna discriminatória para banco de dados *multi-tenancy*](#coluna-discriminatria-para-banco-de-dados-multi-tenancy)
* [Banco de dados (PostgreSQL)](#banco-de-dados)
* [Build local](#build-local)
<!-- TOC -->

## Instalação

Baixe o executável para sua plataforma [na página de *releases*](https://github.com/alanwgt/jsync/releases) e adicione
ao `PATH`.

```bash
# baixar o executável
curl -L -o jsync.tar.gz https://github.com/alanwgt/jsync/releases/latest/download/jsync_[version]_[OS]_[ARCH].tar.gz
# descompactar
tar -xvf jsync.tar.gz
# mover para uma pasta do PATH para poder ser acessado globalmente
sudo mv jsync /usr/local/bin
```

## TLDR

1. Baixe o executável e rode ele uma primeira vez para criar o arquivo de configurações em `$HOME/.jsync.yaml`
2. Modifique o arquivo de configurações conforme a sua necessidade
    1. Remover as linhas comentadas, modificar o valor de conexão com o banco e adicionar a `webservice_key` é suficiente para a maioria dos casos
3. [Execute o SQL](./migrations/000001_create_tables.up.sql) de criação das tabelas (ou utilize a [migration](#banco-de-dados))
4. Execute o programa com os parâmetros: `jsync sync all` para sincronizar todos os dados. Esse comando pode ser posto
numa entrada cron para manter o banco atualizado.

> **Note** \
> Para fins de teste, um banco local pode ser levantado no docker com: `docker run --rm -it -e POSTGRES_PASSWORD=xxx --publish 5432:5432 postgres:latest-alpine`.
> Após a interrupção do processo, todos os dados do container serão removidos.

## Uso

O uso direto da ferramenta se dá através do comando: `jsync sync [recurso]`, onde `[recurso]` pode ser uma opção entre:

- `properties`: imóveis
- `condominiums`: condomínios
- `brokers`: corretores
- `banners`: banners
- `all`: sincroniza todos os recursos

Execute o comando `jsync help` para mais informações sobre os comandos e flags disponíveis.

## Configurações

Ao executar o programa pela primeira vez, um arquivo de configurações base será criado por padrão em `$HOME/.jsync.yaml`
(pode ser sobrescrito utilizando a flag `--config`).

Para iniciar, deixe uma opção dentre os dois próximos blocos no seu arquivo:

##### sincronização para uma imobiliária

- `webservice_key`: chave de integração fornecida pela Jetimob

##### sincronização para múltiplas imobiliárias

> **Note** \
> *multi-tenancy*: Uma aplicação para múltiplos clientes. Cada cliente é denominado *tenant* da aplicação.

- `tenant_column`: nome da coluna de identificação da imobiliária do seu banco
- `tenant_mapping`: vetor de objetos com a estrutura abaixo:
    - `identifier`: identificador da imobiliária
    - `webservice_key`: 

Exemplo:

```yaml
tenant_mapping:
    - identifier: xxx
      webservice_key: xxx
    - identifier: yyy
      webservice_key: yyy
    # ...
```

##### configurações obrigatórias

- `db`:
    - `connection_string`: postgres://[usuário]:[senha]@[host]:[porta]/[database]?sslmode=disable
- `mappings`:
    - `banners_table` (optional,default=*banners*): nome da tabela de banners
    - `brokers_table` (optional,default=*brokers*): nome da tabela de corretores
    - `condominiums_table` (optional,default=*condominiums*): nome da tabela de condomínios
    - `properties_table` (optional,default=*properties*): nome da tabela de imóveis
    - `banners`: mapeamento das colunas disponíveis de banners para colunas do banco de dados
    - `brokers`: mapeamento das colunas disponíveis de corretores para colunas do banco de dados
    - `condominiums`: mapeamento das colunas disponíveis de condomínios para colunas do banco de dados
    - `properties`: mapeamento das colunas disponíveis de imóveis para colunas do banco de dados
- `truncate_all` (bool): remove TODOS os dados da tabela sendo sincronizada. Se for `false` (default), apenas *rows* conflitantes serão removidas

Cofigurações de mapemento de um recurso para a tabela do banco de dados são feitas da forma em que a chave de
configuração representa o nome do dado e o valor o nome da coluna no banco de dados. Chaves removidas não serão
inseridas.

Por exemplo, considere o mapeamento reduzido:

```yaml
mappings:
    brokers_table: corretores
    brokers:
        avatar: avatar
        biografia: biography
        cargo: job_position
        # ... 
```

No mapeamento acima, quando houver a sincronização dos banners, o `jsync` fará o insert: `INSERT INTO corretores (avatar, biography, job_position) VALUES (...)`
com os valores de cada ítem dos banners.

### Coluna discriminatória para banco de dados *multi-tenancy*

> **Note** \
> *multi-tenancy*: Uma aplicação para múltiplos clientes. Cada cliente é denominado *tenant* da aplicação.

Numa aplicação *multi-tenant* com banco de dados também *multi-tenant*, a aplicação aceita vários clientes que compartilham
o mesmo banco de dados. Cada entrada das tabelas desse banco de dados possui uma coluna que identifica o cliente que detém
propriedade sobre aquela informação. Dessa forma cada cliente verá apenas os seus dados (considerando que a aplicação
faça o filtro nas queries de forma adequada, pois se um `where` for esquecido nas queries, um cliente poderá ver
informações que não lhe perctencem).

Isto é, considerando a tabela:

| id  | code | image                                                                                                                                   |
|:---:|:----:|:----------------------------------------------------------------------------------------------------------------------------------------|
|  1  | XXXX | https://s01.jetimgs.com/trvAWQHuYcArjvEQrh93oEZSAxK0Jz8p2OIdekopXlWDY5-MAAMBV0DPcGX3lxwoOeyVrBgSUbpqY-efLaLw_YZiMIVV0qN3gf2D/1660788979 |

Podemos adicionar uma coluna `tenant_id` para identificar o dono daquela informação. Supondo que temos uma tabela chamada
`tenants`, adicionando uma FK de `tenant_id` para `tenants.id` conseguimos vincular uma entrada da tabela a um proprietário
dessa informação. Ficando então:

| id  | tenant_id | code | image                                                                                                                                   |
|:---:|:---------:|:----:|:----------------------------------------------------------------------------------------------------------------------------------------|
|  1  |     1     | XXXX | https://s01.jetimgs.com/trvAWQHuYcArjvEQrh93oEZSAxK0Jz8p2OIdekopXlWDY5-MAAMBV0DPcGX3lxwoOeyVrBgSUbpqY-efLaLw_YZiMIVV0qN3gf2D/1660788979 |

## Banco de dados

> 💡 as migrações criam uma coluna chamada `tenant_id` em todas as tabelas. Essa coluna pode ser desconsiderada ou removida.

Se existirem dúvidas em como construir o banco de dados, utilize [este arquivo](./migrations/000001_create_tables.up.sql)
como base, ou *as is* para uso em produção.

Para executar as *migrations*, primeiro, baixe a ferramenta:

Troque as variáveis **[OS]** e **[ARCH]** para refletir a arquitetura do computador que executará as migrações. Visite o [site](https://github.com/golang-migrate/migrate/releases/latest) para possíveis opções.

```bash
curl -L -o migrate https://github.com/golang-migrate/migrate/releases/latest/download/migrate.[OS]-[ARCH].tar.gz
```

E, finalmente, execute as migrações:

```bash
./migrate -source "github://alanwgt/jsync/migrations" -database "postgres://[usuário]:[senha]@[host]:[porta]/[database]?sslmode=disable" up
```

## Build local

1. Assegure-se que o `go` está [instalado](https://go.dev/dl/) e incluso no [`PATH` global](https://go.dev/doc/install)
2. Instale as dependências: `go mod download`
3. Execute o programa com `go run main.go [params]` ou crie um executável: `go build -o jsync main.go`

jsync sincroniza dados na Jetimob com um banco d
sincronizador de dados jetimob p sites exclusivos
