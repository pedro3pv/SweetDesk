# GitHub Actions Workflow Setup

Este repositório usa GitHub Actions para builds automatizados. Como o SweetDesk-core é um repositório privado, o workflow usa o `GITHUB_TOKEN` automático para clonar o repositório durante o build.

## Permissões Configuradas

O workflow já está configurado com as permissões necessárias:

```yaml
permissions:
  contents: write  # Para criar releases e tags
  packages: read   # Para acessar repositórios privados
```

**✅ Nenhuma configuração manual é necessária!** O `GITHUB_TOKEN` é gerado automaticamente pelo GitHub Actions e tem acesso a repositórios privados da mesma organização.

## Como Funciona

1. O workflow clona o SweetDesk (repositório público/privado atual)
2. Usa `GITHUB_TOKEN` para clonar o SweetDesk-core (repositório privado):
   ```bash
   git clone https://x-access-token:${{ secrets.GITHUB_TOKEN }}@github.com/pedro3pv/SweetDesk-core.git
   ```
3. Descomenta automaticamente a linha de replace no go.mod
4. Build usa a versão local do SweetDesk-core em `../SweetDesk-core`
5. Binários são gerados e enviados como artifacts

## Verificar o Workflow

Após fazer commit:

1. Vá para **Actions** no repositório
2. Veja o workflow rodando
3. O step "Clone SweetDesk-core" deve funcionar automaticamente
4. O step "Verify SweetDesk-core cloned" confirma que o clone foi bem-sucedido

## Troubleshooting

### Erro: "Repository not found" ou "Permission denied"

**Causa:** O repositório SweetDesk-core não está acessível ao GITHUB_TOKEN.

**Soluções:**
1. Verifique se ambos os repositórios estão na mesma organização (Molasses-Co)
2. Certifique-se de que o repositório SweetDesk-core existe em: https://github.com/pedro3pv/SweetDesk-core
3. Se os repositórios estão em contas diferentes (Molasses-Co vs pedro3pv), você precisará criar um PAT:
   - Crie um PAT em: https://github.com/settings/tokens
   - Adicione como secret `PAT_TOKEN` no repositório
   - Modifique o workflow para usar `${{ secrets.PAT_TOKEN }}`

### Erro: "Resource not accessible by integration"

**Causa:** Permissões insuficientes no workflow.

**Solução:** O workflow já tem as permissões configuradas:
```yaml
permissions:
  contents: write
  packages: read
```

### Workflow não está rodando

**Causa:** Workflow só roda em push para `main` ou tags `v*`.

**Solução:** Para testar em outras branches, modifique `.github/workflows/main.yml`:

```yaml
on:
  push:
    branches: [ main, fix/*, feature/* ]  # Adicione patterns aqui
```

## Estrutura do Build

O workflow executa os seguintes passos:

1. **Checkout** - Clona o repositório SweetDesk
2. **Clone SweetDesk-core** - Clona o repositório privado SweetDesk-core
3. **Enable replace directive** - Descomenta a linha de replace no go.mod
4. **Setup Bun & Go** - Instala dependências de build
5. **Install system deps** - Instala dependências do sistema (Linux)
6. **Tidy Go modules** - Atualiza dependências Go
7. **Install Wails CLI** - Instala a CLI do Wails
8. **Build** - Compila o app para Linux, Windows e macOS
9. **Upload artifacts** - Faz upload dos binários

## Plataformas Suportadas

- **Linux**: `ubuntu-latest` → `SweetDesk` (ELF binary)
- **Windows**: `windows-latest` → `SweetDesk.exe`
- **macOS**: `macos-latest` → `SweetDesk.app` (Universal binary)

## Release Automático

Quando você cria uma tag de versão (ex: `v1.0.0`), o workflow:

1. Executa todos os builds
2. Cria uma GitHub Release automaticamente
3. Anexa os binários de todas as plataformas
4. Gera release notes automaticamente

Para criar uma release:

```bash
git tag v1.0.0
git push origin v1.0.0
```
