# RoboCombat Controller

> Control remoto para robots de combate con soporte de mando Xbox y teclado

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com)
[![RoboKombat](https://img.shields.io/badge/🤖-RoboKombat-red?style=flat)](https://rooted.es)

Aplicación de control para robots de combate diseñada para la competición **RoboKombat** en la **Rooted**. Permite controlar motores y servos mediante mando Xbox o teclado, con configuración totalmente personalizable.

---

## 📸 Demo

```
     🎮 Xbox Controller                    ⌨️  Keyboard
     
    ┌─────────────────┐               W I K
    │  ╱═╲       ╱═╲  │               ↑ ↑ ↑
    │ ╱   ╲     ╱   ╲ │            Motors A & B
    │ ╲   ╱     ╲   ╱ │
    │  ╲═╱       ╲═╱  │          A D     J L
    │    ⚫         ⚫  │          ← →     ← →
    │   Left      Right│        Servo1   Servo2
    └─────────────────┘
      L1  L2  R1  R2
       ↓   ↓   ↓   ↓            SPACE = STOP
     Servo 1  Servo 2
```

## ✨ Características

- 🎮 **Doble modo de entrada**: Xbox Controller (XInput) o teclado
- 🏎️ **Control tipo tanque**: Cada motor se controla independientemente para máxima maniobrabilidad
- 🔧 **Configuración YAML**: Sin necesidad de recompilar para cambiar ajustes
- 🔄 **Inversión de motores**: Soluciona motores mal soldados por software
- 🎯 **Control de servos**: 2 servos con rangos personalizables
- ⌨️ **Teclas personalizables**: Define tus propias teclas en modo teclado
- 🚀 **Sin dependencias externas**: Binario único y portable
- 🐧 **Multiplataforma**: Windows, Linux, macOS (teclado), Windows (Xbox)

## 📋 Requisitos

### Hardware
- Robot compatible con la API de RoboKombat (ESP32)
- Mando Xbox (opcional, para modo controller)
- PC/Laptop con conexión WiFi

### Software
- Windows 10/11, Linux, o macOS
- Go 1.21+ (solo para compilar desde código fuente)

## 🚀 Instalación Rápida

### Opción 1: Descargar ejecutable pre-compilado
1. Descarga `robocombat.exe` desde [Releases](https://github.com/GoldraK/RoboKombat-rootedcon/releases)
2. Descarga `config.yaml` de este repositorio
3. Coloca ambos archivos en la misma carpeta
4. ¡Listo para usar!

### Opción 2: Compilar desde código fuente

```bash
# Clonar el repositorio
git clone https://github.com/GoldraK/RoboKombat-rootedcon.git
cd robocombat

# Instalar dependencias
go mod download

# Compilar para Windows
GOOS=windows GOARCH=amd64 go build -o robocombat.exe

# Compilar para Linux
go build -o robocombat

# Compilar para macOS
GOOS=darwin GOARCH=amd64 go build -o robocombat-mac
```

## ⚙️ Configuración

Edita el archivo `config.yaml` antes de usar:

```yaml
# Modo de entrada: "xbox" o "keyboard"
inputMode: xbox

# Inversión de motores (si están mal soldados)
motors:
  invertA: false  # Motor izquierdo
  invertB: false  # Motor derecho

# Rangos de servos (0-180 grados)
servos:
  servo1:
    min: 0
    max: 180
  servo2:
    min: 0
    max: 180

# Teclas personalizadas (solo para modo keyboard)
keyboard:
  motorA:
    forward: w
    backward: s
  motorB:
    forward: i
    backward: k
  servo1:
    min: a
    max: d
  servo2:
    min: j
    max: l
  stop: space
```

### Solucionar motores invertidos
Si un motor gira al revés, simplemente cambia en el config:
```yaml
motors:
  invertA: true  # Ahora motor A gira correctamente
```
**¡No necesitas re-soldar!** 🎉

## 🎮 Uso

### Conectar al robot

1. Conéctate a la red WiFi del robot (AP: `192.168.4.1`)
2. Ejecuta la aplicación:

```bash
# Windows (modo Xbox por defecto)
.\robocombat.exe -host 192.168.4.1

# Con modo debug activado
.\robocombat.exe -host 192.168.4.1 -debug

# Especificar archivo de configuración personalizado
.\robocombat.exe -host 192.168.4.1 -config miconfig.yaml

# Usar controlador específico (0-3)
.\robocombat.exe -host 192.168.4.1 -controller 1
```

### Control con Mando Xbox

**Movimiento (Modo Tanque)**
- **Stick izquierdo** (arriba/abajo) → Motor A (izquierdo)
- **Stick derecho** (arriba/abajo) → Motor B (derecho)

**Servos**
- **LB/L1** → Servo 1 posición mínima
- **LT/L2** → Servo 1 posición máxima
- **RB/R1** → Servo 2 posición mínima
- **RT/R2** → Servo 2 posición máxima

**Ejemplos de movimiento:**
- Ambos sticks arriba → Avanzar recto
- Ambos sticks abajo → Retroceder recto
- Stick izq. arriba + der. abajo → Girar a la derecha
- Stick izq. abajo + der. arriba → Girar a la izquierda

### Control con Teclado

Cambia en `config.yaml`:
```yaml
inputMode: keyboard
```

**Controles por defecto:**
- **W/S** → Motor A adelante/atrás
- **I/K** → Motor B adelante/atrás
- **A/D** → Servo 1 min/max
- **J/L** → Servo 2 min/max
- **Espacio** → Parada de emergencia

> **Nota**: Los motores se detienen automáticamente al soltar las teclas. Los servos vuelven al centro (90°).

## 🔧 API del Robot

La aplicación es compatible con robots que implementen la siguiente API REST:

| Endpoint | Parámetros | Descripción |
|----------|-----------|-------------|
| `/api/motors` | `a`, `b` | Velocidad motores (-255..255) |
| `/api/servo` | `id`, `angle` | Posición servo (id: 1-2, angle: 0-180) |
| `/api/control` | `a`, `b`, `s1`, `s2` | Control completo en una llamada |
| `/api/stop` | - | Parada de emergencia |
| `/api/status` | - | Estado actual (JSON) |
| `/api/battery` | - | Nivel de batería |

## 🐛 Solución de Problemas

### El mando Xbox no se detecta
- ✅ **Windows**: Asegúrate de que el mando esté conectado vía USB o Bluetooth
- ✅ Prueba con `-controller 0`, `-controller 1`, etc.
- ✅ Verifica que el LED del mando esté encendido
- ✅ Comprueba en "Dispositivos de juego" (Windows)

### El robot no responde
- ✅ Verifica la conexión WiFi al AP del robot
- ✅ Prueba hacer ping a `192.168.4.1`
- ✅ Usa el flag `-debug` para ver los comandos enviados
- ✅ Revisa que la API del robot esté funcionando: `http://192.168.4.1/api/status`

### Los motores van al revés
- ✅ Edita `config.yaml` y activa `invertA` o `invertB`
- ✅ No hace falta recompilar ni re-soldar

### El teclado no funciona en WSL
- ✅ Compila para Windows y ejecuta desde PowerShell/CMD
- ✅ WSL no tiene acceso directo al hardware de entrada

## 🏆 Competición RoboKombat - Rooted

Este proyecto fue diseñado específicamente para facilitar el control de robots en la competición de RoboKombat. 

### Consejos para la competición:

1. **Prueba tu configuración antes**: Verifica inversión de motores y rangos de servos
2. **Modo debug es tu amigo**: Usa `-debug` para diagnosticar problemas rápido
3. **Ten un backup**: Lleva el `config.yaml` en un USB por si acaso
4. **Batería**: Monitoriza el nivel con `/api/battery` antes de cada combate
5. **Parada de emergencia**: Familiarízate con el botón/tecla de stop

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! Si quieres mejorar el proyecto:

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Estructura del Proyecto

```
robocombat/
├── main.go              # Lógica principal y modos de control
├── xinput_windows.go    # Implementación XInput para Windows
├── xinput_linux.go      # Stub para compatibilidad multiplataforma
├── config.yaml          # Configuración de usuario
├── go.mod               # Dependencias Go
├── README.md            # Este archivo
└── .instructions.md     # Documentación técnica detallada
```

## 📚 Documentación Técnica

Para más detalles sobre la implementación y las mejores prácticas de desarrollo, consulta [.instructions.md](.instructions.md).

## 📄 Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 🙏 Agradecimientos

- Comunidad de **RoboKombat** y **Rooted**
- Librería [eiannone/keyboard](https://github.com/eiannone/keyboard) para control de teclado
- Librería [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) para parsing de configuración

## 📞 Contacto

¿Problemas o sugerencias? Abre un [Issue](https://github.com/GoldraK/RoboKombat-rootedcon/issues) en GitHub.

---

Hecho con ❤️ para la comunidad de RoboKombat  
**¡Buena suerte en tus combates!** 🤖⚔️🤖
