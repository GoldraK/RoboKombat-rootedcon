# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-06

### Added
- 🎮 Xbox controller support via XInput (Windows native)
- ⌨️ Keyboard control mode with customizable key mappings
- 🏎️ Tank mode control (independent motor control)
- 🎯 Dual servo control with configurable ranges
- 📝 YAML configuration file for easy customization
- 🔄 Motor inversion via config (no need to re-solder)
- 🐛 Debug mode with detailed logging
- 📱 Cross-platform support (Windows, Linux, macOS)
- 🚀 Single binary with no external dependencies
- 🔋 Battery monitoring support
- 🛑 Emergency stop functionality
- 📊 Real-time status monitoring
- 🎨 Configurable servo angle ranges
- ⚡ Automatic motor stop on key release (keyboard mode)
- 🔄 Servo return to center on key release (keyboard mode)

### Features by Mode

#### Xbox Controller Mode
- Left/Right stick Y-axis for tank control
- L1/L2 for servo 1 min/max positions
- R1/R2 for servo 2 min/max positions
- Automatic deadzone application
- Support for up to 4 controllers

#### Keyboard Mode
- Fully customizable key mappings
- Support for letter keys (a-z)
- Support for special keys (space, enter, tab, esc, arrows)
- Hold-to-move behavior
- Immediate stop on key release

### API Endpoints Supported
- `/api/motors` - Motor speed control
- `/api/servo` - Individual servo control
- `/api/control` - Combined motors and servos
- `/api/stop` - Emergency stop
- `/api/status` - Status query
- `/api/battery` - Battery level

### Documentation
- Comprehensive README with quick start guide
- Configuration examples
- API reference
- Troubleshooting section
- Contributing guidelines

## [Unreleased]

### Planned Features
- [ ] PlayStation controller support
- [ ] Web UI for configuration
- [ ] Telemetry recording
- [ ] Multiple robot profiles
- [ ] Macro/combo system
- [ ] Auto-reconnect on disconnect

---

[1.0.0]: https://github.com/tu-usuario/robocombat/releases/tag/v1.0.0
