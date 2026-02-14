# Changelog

All notable changes to SweetDesk will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.2] - 2026-02-14

### Changed
- macOS build now outputs DMG instead of tar.gz for easier installation
- Added ad-hoc codesigning for macOS builds to prevent "damaged app" errors
- Updated GitHub Actions workflow for macOS DMG packaging using create-dmg

### Added
- macOS folder access permission prompts (Desktop, Documents, Downloads, Pictures)
- Entitlements plist for macOS ad-hoc signing with file access and network permissions
- Version info in wails.json (productVersion: 0.0.2)

### Fixed
- macOS permission issues when accessing user folders

## [0.0.1] - 2026-02-07

### Initial Release
- Project structure setup
- Wails + Next.js integration
- Basic application scaffolding

---

## Future Roadmap

### v0.0.1 (Next)
- [ ] Download actual AI binaries (Real-ESRGAN, RealCUGAN)
- [ ] Implement "Set as Wallpaper" functionality
- [ ] Add batch processing support
- [ ] Improve error handling and user feedback
- [ ] Add processing queue for multiple images
- [ ] Implement image caching

### v0.2.0
- [ ] Add Unsplash API integration
- [ ] Add Pexels API integration
- [ ] Implement Supabase storage integration
- [ ] Add user preferences/settings
- [ ] Multi-language support (PT-BR, ES, JA)
- [ ] Add image filters and adjustments

### v0.3.0
- [ ] GPU acceleration for upscaling
- [ ] Advanced classification using DeepGHS/imgutils
- [ ] Face detection for smart cropping
- [ ] Color correction and enhancement
- [ ] HDR processing

### v1.0.0
- [ ] Stable release with all core features
- [ ] Complete test coverage
- [ ] Production-ready AI binaries
- [ ] Installer packages for all platforms
- [ ] macOS App Store release
- [ ] Windows Store release

### v2.0.0 (Future)
- [ ] AI-powered wallpaper generation (Text-to-Image)
- [ ] Wallpaper marketplace integration
- [ ] Cloud sync between devices
- [ ] Mobile app (iOS/Android)
- [ ] Scheduled wallpaper rotation
- [ ] Custom AI model training

---

**Legend:**
- `Added` - New features
- `Changed` - Changes in existing functionality
- `Deprecated` - Soon-to-be removed features
- `Removed` - Removed features
- `Fixed` - Bug fixes
- `Security` - Security improvements
