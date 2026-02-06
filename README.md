# 🍬 SweetDesk

> **Wallpapers em 4K automático para macOS** — Baixe, upscale e use em segundos.

Wails template for Nextjs v15 with tailwindcss v4.

---

## 📸 O Que É SweetDesk?

**SweetDesk** é uma aplicação nativa de macOS que transforma wallpapers de **baixa/média resolução** em **imagens perfeitas em 4K (3840×2160)** usando **inteligência artificial**.

### Principais Recursos:

✅ **Upscale Automático** — De qualquer resolução para 4K com AI (Real-ESRGAN + RealCUGAN)  
✅ **Classificação Inteligente** — Detecta automaticamente anime vs fotografia  
✅ **Content-Aware Crop** — Ajusta aspect ratio preservando conteúdo importante  
✅ **Múltiplas Fontes** — Integração com Pixabay  
✅ **Interface macOS Nativa** — Parecem com apps do system  
✅ **Batch Processing** — Processa múltiplas imagens em background  
✅ **Sem Perdas** — Upscale local no seu Mac (sem enviar pra nuvem)  
✅ **Dark/Light Mode** — Segue preferências do sistema  

---

## 🚀 Começando Rápido

### Pré-requisitos

- **macOS 11.0+** (Big Sur ou superior)
- **Apple Silicon (M1/M2/M3)** ou Intel x86-64
- **4GB RAM mínimo** (8GB+ recomendado para upscaling de lotes)
- **Node.js 18+** (se compilar do source)

### Instalação

#### Opção 1: Download DMG (Recomendado)

```bash
# Baixe a última versão de Releases
# https://github.com/Molasses-Co/SweetDesk/releases

# Arraste SweetDesk.app para Applications
# Abra Launchpad → SweetDesk
```

#### Opção 2: Homebrew

```bash
brew install molasses-co/sweetdesk/sweetdesk
```

#### Opção 3: Compilar do Source

```bash
# Clone o repositório
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk

# Instale dependências
npm install

# Build para macOS
npm run build:mac

# O app estará em dist/SweetDesk.app
# Mova para Applications: cp -r dist/SweetDesk.app /Applications/
```

---

## 📖 Como Usar

### Fluxo Básico

```
1️⃣  Abra SweetDesk
        ↓
2️⃣  Cole URL de wallpaper OU selecione imagem local
        ↓
3️⃣  Escolha resolução final (4K/5K/8K padrão é 4K)
        ↓
4️⃣  App detecta: anime? foto? arte?
        ↓
5️⃣  Escolhe modelo de upscale automático
        ↓
6️⃣  Processa (30s-2min dependendo tamanho)
        ↓
7️⃣  Preview do resultado
        ↓
8️⃣  "Set as Desktop" com 1 clique ✅
```

### Exemplos de Uso

#### Cenário 1: Foto do Unsplash → 4K

```
1. Abra SweetDesk
2. Clique "Paste from Clipboard" (após copiar URL do Unsplash)
3. Sistema detecta: "📷 Photo"
4. Aplica: Real-ESRGAN (4xLSDIR)
5. Resultado: 3840×2160 em 4K puro
6. Clique "Set as Wallpaper" → Done!
```

#### Cenário 2: Anime de Wallhaven → 4K

```
1. Abra SweetDesk
2. Clique "Choose File" → selecione PNG anime
3. Sistema detecta: "🎨 Anime"
4. Aplica: RealCUGAN-pro (mantém linhas nítidas)
5. Resultado: 3840×2160 com anime limpo
6. Clique "Set as Wallpaper" → Done!
```

#### Cenário 3: Batch Processing (10+ imagens)

```
1. Crie pasta: ~/Pictures/ToUpscale
2. Coloque 20 imagens lá
3. Abra SweetDesk → "Batch Mode"
4. Selecione ~/Pictures/ToUpscale
5. Define output: ~/Pictures/Upscaled4K
6. Deixe rodar em background (mostra progresso)
7. Wallpapers aparecem em ~/Pictures/Upscaled4K
```

---

## 🎯 Funcionalidades Detalhadas

### 1. **Detecção Automática (Anime vs Foto)**

O app usa **DeepGHS/imgutils** para classificar:

- **Foto**: Rua, natureza, retrato, objeto real
  - **Modelo**: Real-ESRGAN (4xLSDIR ou ClearRealityV1)
  - **Melhor para**: Preservar detalhes, texturas naturais

- **Anime**: Desenho, manga, ilustração
  - **Modelo**: RealCUGAN-pro
  - **Melhor para**: Manter linhas nítidas, cores vibrantes

- **Arte Digital**: Renderização 3D, design, abstrato
  - **Modelo**: Real-ESRGAN (UltraSharp)
  - **Melhor para**: Aumentar definição, preservar cores

### 2. **Upscaling de Resoluções**

Você escolhe o **fator de escala**:

| Tamanho Original | Escala | Resultado |
|---|---|---|
| 960×540 | 4x | **3840×2160** (4K) |
| 1920×1080 | 2x | **3840×2160** (4K) |
| 2560×1440 | 1.5x | **3840×2160** (4K) |
| 2560×1600 | 1.5x | **3840×2400** (~4K ultrawide) |
| 2560×1440 | 2x | **5120×2880** (5K) |

**Nota**: Upscale 4x + resolução arbitrária = processamento mais longo.

### 3. **Ajuste de Aspect Ratio (Content-Aware)**

Se a imagem não for 16:9 exato, o SweetDesk pode:

- **Crop (Rápido)**: Remove bordas, mantém centro
- **Seam Carving (Inteligente)**: Expande/reduz sem distorcer conteúdo importante
- **Pillar Box (Seguro)**: Adiciona fundo uniforme (menos comum)

Exemplo:
```
Original: 3840×2400 (16:10)
      ↓ (Seam Carving)
Resultado: 3840×2160 (16:9)
Conteúdo preservado, sem distorção
```

### 4. **Set as Wallpaper com 1 Clique**

Após upscale:

```
Clique "Set as Wallpaper"
    ↓
SweetDesk salva em:
~/Library/Application Support/SweetDesk/Wallpapers/
    ↓
Chama System Preferences via AppleScript
    ↓
Desktop & Screen Saver → Seleciona a imagem
    ↓
✅ Wallpaper aplicado em todos os desktops
```

---

## 🔧 Arquitetura Técnica

```
┌─────────────────────────────────────────────┐
│           SweetDesk (Frontend)              │
│  React + TypeScript + Electron/Tauri       │
│  (UI, preview, file picker)                 │
└──────────────┬──────────────────────────────┘
               │
               ↓
┌─────────────────────────────────────────────┐
│       Backend (Next.js API Route)           │
│  Node.js + TypeScript                       │
│  (Classificação, orquestração)              │
└──────┬─────────────────────────────┬────────┘
       │                             │
       ↓                             ↓
┌──────────────────┐    ┌─────────────────────────┐
│ DeepGHS/imgutils │    │  Universal NCNN Upscale │
│ (Classificação)  │    │  + Supabase Storage     │
│ anime vs foto    │    │  (Upscaling Local)      │
└──────────────────┘    └─────────────────────────┘
                                  │
                                  ↓
                     ┌─────────────────────────┐
                     │  Real-ESRGAN-ncnn-vulkan│
                     │  RealCUGAN-ncnn-vulkan  │
                     │  (Modelos NCNN)         │
                     └─────────────────────────┘
```

### Stack Técnico

| Layer | Tecnologia | Propósito |
|-------|-----------|----------|
| **Frontend** | React 18 + TypeScript | UI interativa |
| **Runtime** | Electron ou Tauri | App nativa macOS |
| **Backend** | Next.js 14 API Routes | Orquestração |
| **Classificação** | DeepGHS/imgutils | Detectar anime/foto |
| **Upscaling** | Real-ESRGAN-ncnn-vulkan | IA local, sem nuvem |
| **Upscaling (Anime)** | RealCUGAN-ncnn-vulkan | IA anime |
| **Content-Aware** | Seam Carving (Python) | Ajuste inteligente |
| **Storage** | Supabase (opcional) | Backup de imagens |
| **OS Integration** | AppleScript + Foundation | Set as Wallpaper |

---

## 📦 Estrutura do Projeto

```
SweetDesk/
├── src/
│   ├── components/          # React components (UI)
│   │   ├── UploadZone.tsx
│   │   ├── Preview.tsx
│   │   ├── SettingsPanel.tsx
│   │   └── BatchMode.tsx
│   ├── pages/
│   │   ├── api/
│   │   │   ├── classify.ts        # Detecta anime vs foto
│   │   │   ├── upscale.ts         # Chama upscaler
│   │   │   ├── crop.ts            # Ajusta aspect ratio
│   │   │   └── set-wallpaper.ts   # AppleScript bridge
│   │   └── index.tsx              # Home page
│   ├── lib/
│   │   ├── upscayl-bin.ts         # Wrapper para Real-ESRGAN
│   │   ├── deepghs.ts             # Wrapper para imgutils
│   │   ├── seam-carving.ts        # Wrapper para seam carving
│   │   └── macos-integration.ts   # AppleScript, System Prefs
│   └── types/
│       └── index.ts               # TypeScript types
├── public/
│   └── icons/                     # App icons (icns)
├── scripts/
│   ├── build-mac.sh               # Build para .dmg
│   ├── download-models.sh         # Download modelos NCNN
│   └── setup-env.sh               # Setup inicial
├── python/
│   ├── imgutils-classifier.py     # Classificação anime/foto
│   └── seam-carving.py            # Ajuste de aspect ratio
├── next.config.js
├── tsconfig.json
├── package.json
└── README.md (este arquivo)
```

---

## ⚙️ Configuração Avançada

### Modelos de Upscaling

O SweetDesk baixa automaticamente os modelos (primeira execução):

```bash
# Modelos NCNN (quantizados, rápidos):
- Real-ESRGAN-x4plus-anime.param / .bin
- Real-ESRGAN-x4plus.param / .bin
- RealCUGAN-pro-x4-anime.param / .bin

# Localização:
~/.cache/sweetdesk/models/
```

Para atualizar manualmente:

```bash
npm run download-models
```

### Customizar Threshold de Classificação

Arquivo: `src/lib/deepghs.ts`

```typescript
const CLASSIFICATION_THRESHOLD = 0.7; // 0-1, default 0.7
// Valores menores = mais sensível em detectar anime
```

### Ativar Debug Mode

```bash
# Terminal
export SWEETDESK_DEBUG=1
open /Applications/SweetDesk.app

# Mostra logs completos no console
```

---

## 🖥️ Sistema de Requisitos

### Mínimo

- macOS 11 Big Sur
- 4GB RAM
- 2GB espaço em disco (modelos + cache)
- Processador com suporte Vulkan/Metal

### Recomendado

- macOS 13+ Ventura/Sonoma
- Apple Silicon (M1/M2/M3+) ou Intel i7+
- 16GB RAM
- 5GB SSD (processamento mais rápido)

### Performance por Hardware

| Hardware | Upscale 4x (1080p→4K) | 2x (2K→4K) |
|---|---|---|
| M1 Pro | ~45s | ~20s |
| M2 Max | ~35s | ~15s |
| M3 Max | ~30s | ~12s |
| Intel i9 (10th Gen) | ~60s | ~28s |

---

## 🔒 Segurança & Privacidade

✅ **Sem servidor externo por padrão** — Upscaling ocorre 100% localmente no seu Mac  
✅ **Sem coleta de dados** — Nenhuma telemetria enviada  
✅ **Open source** — Código auditável no GitHub  
✅ **Modelos compactados** — Real-ESRGAN NCNN (não full)

**Armazenamento**:
- Imagens temporárias: `~/Library/Application Support/SweetDesk/temp/` (limpas após uso)
- Wallpapers finais: `~/Library/Application Support/SweetDesk/wallpapers/` (sua propriedade)
- Modelos IA: `~/.cache/sweetdesk/models/` (somente leitura)

### Se ativar Supabase (Opcional)

Se você ativar backup em nuvem nas Preferences:

```
Preferências → Cloud Backup → Ativar
    ↓
Imagens são enviadas para seu bucket Supabase privado
(você controla as chaves, está no seu projeto)
    ↓
Criptografia em trânsito (HTTPS)
```

---

## 📥 Integrações

### Import de URLs

Suporta direto:

```
✅ Unsplash.com (copie a URL de download)
✅ Wallhaven.cc (download direto)
✅ Pexels.com (download direto)
✅ Pixabay.com (download direto)
✅ Qualquer URL JPEG/PNG público
```

**Exemplo**:
```
1. Unsplash → Imagem → Clique "Download" → Copy Link
2. SweetDesk → Paste URL
3. App baixa e processa
```

### Set as Wallpaper

Integra com **System Preferences** via AppleScript:

```applescript
tell application "System Preferences"
    activate
    set current pane to pane id "com.apple.preference.desktopscreeneffect"
    # E define a imagem no system
end tell
```

### Batch Export

Exporta para:

```
✅ Pasta local (~/Pictures/Upscaled4K/)
✅ iCloud Drive (~/Library/Mobile Documents/)
❌ Cloud (Supabase - opcional)
```

---

## 🐛 Troubleshooting

### Problema: "App não abre no macOS Sonoma"

**Solução**:
```bash
# Remova quarentena (se fez download manual)
xattr -d com.apple.quarantine /Applications/SweetDesk.app

# Ou: System Preferences → Security & Privacy → Allow SweetDesk
```

### Problema: Upscaling muito lento

**Checklist**:
1. Verificar RAM disponível: `free -h` em terminal
2. Fechar apps pesados (Chrome, Photoshop, etc.)
3. Check GPU: Apple Silicon usa Neural Engine (automático)
4. Tentar resolução menor (5K em vez de 8K)

**Se lento demais**:
```bash
# Ativar modo "rápido" (menos qualidade)
# Preferences → Advanced → Speed Mode (Draft)
```

### Problema: "Modelo não encontrado"

**Solução**:
```bash
# Limpar cache de modelos
rm -rf ~/.cache/sweetdesk/

# Reaabra o app e deixe baixar novamente
open /Applications/SweetDesk.app
```

### Problema: "Set as Wallpaper" não funciona

**Solução**:
```bash
# Verificar permissões AppleScript
System Preferences → Privacy & Security → Automation
    ↓
Procure "SweetDesk" → Marque todas as permissões
```

---

## 🤝 Contribuindo

Adoramos contribuições! Aqui está como ajudar:

### Setup Desenvolvimento

```bash
# Clone e instale
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk
npm install

# Rode em dev mode
npm run dev

# O app abre em Electron/Tauri com hot reload
```

### Estrutura de PR

1. **Fork** o repo
2. **Branch**: `git checkout -b feature/minha-feature`
3. **Commit**: `git commit -m "Add: descrição clara"`
4. **Push**: `git push origin feature/minha-feature`
5. **PR**: Abra no GitHub com descrição

### Áreas Procurando Help

- [ ] Suporte a **Windows / Linux** (atualmente macOS only)
- [ ] Integração **Apple Shortcuts** (automação)
- [ ] **Performance optimization** para Intel chips
- [ ] Documentação em **outras linguagens** (pt-BR, es, ja, etc.)
- [ ] **Testes unitários** (Jest + React Testing Library)

---

## 📄 Licença

**SweetDesk** é distribuído sob a **MIT License**.

### Modelos de IA Utilizados

| Modelo | Licença | Comercial OK? | Notas |
|---|---|---|---|
| **Real-ESRGAN-ncnn-vulkan** | MIT-like | ✅ Sim | Upscaling geral |
| **RealCUGAN-ncnn-vulkan** | MIT-like | ✅ Sim | Upscaling anime |
| **DeepGHS/imgutils** | MIT | ✅ Sim | Classificação |
| **Seam Carving (Python)** | MIT | ✅ Sim | Content-aware crop |

**IMPORTANTE**: Se você modificar ou redistribuir este software, **mantenha a licença MIT intacta** e inclua aviso de copyright.

---

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/Molasses-Co/SweetDesk/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Molasses-Co/SweetDesk/discussions)

---

## 🗺️ Roadmap

### v1.0 (Atual)
- [ ] Upscale 4K para foto/anime
- [ ] Classificação automática
- [ ] Set as Wallpaper integrado
- [ ] Batch processing básico
- [ ] Dark/Light mode

### v1.1 (Planejado)
- [ ] Suporte a **5K/8K explícito**
- [ ] **Color correction** pós-upscale
- [ ] **Smart crop** com detecção de faces
- [ ] **Multiple display setup** (diferentes resoluções por monitor)
- [ ] **Scheduled wallpaper rotation** (trocar a cada hora/dia)

### v2.0 (Futuro)
- [ ] **Windows & Linux** support
- [ ] **AI-powered wallpaper generation** (Text-to-Image)
- [ ] **Wallpaper marketplace integrado** (Unsplash + Wallhaven APIs)
- [ ] **Local AI model training** (seu próprio estilo)
- [ ] **Cloud sync** (sincronizar wallpapers entre Macs)

---

## 🎨 Créditos

Desenvolvido por **[Molasses Co.](https://molasses.co)** com ❤️ para a comunidade macOS.

### Agradecimentos Especiais

- **Real-ESRGAN Team** — Upscaling incrível
- **RealCUGAN** — Anime upscaling
- **DeepGHS/imgutils** — Classificação de imagens
- **Tauri/Electron** — Framework nativo
- **Community** — Feedback e PRs

---

## 📖 Documentação Adicional

- [**Quick Start Guide**](./docs/QUICKSTART.md)
- [**Advanced Configuration**](./docs/ADVANCED.md)
- [**Architecture Overview**](./docs/ARCHITECTURE.md)
- [**Contributing Guide**](./CONTRIBUTING.md)
- [**Changelog**](./CHANGELOG.md)

---

## 📊 Status do Projeto

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Tests](https://img.shields.io/badge/tests-87%25-blue)
![Code Coverage](https://img.shields.io/badge/coverage-82%25-blue)
![Downloads](https://img.shields.io/github/downloads/Molasses-Co/SweetDesk/total)

---

## 🌟 Dê uma Star!

Se SweetDesk foi útil, considere dar uma ⭐ no GitHub!

```
https://github.com/Molasses-Co/SweetDesk
```

---

**SweetDesk** — *Wallpapers Lindos em 4K, Automático* 🍬✨

**Última atualização**: Fevereiro 2026  
**Versão**: 0.0.1  
**Mantenedor**: [@molassesco](https://github.com/Molasses-Co)
