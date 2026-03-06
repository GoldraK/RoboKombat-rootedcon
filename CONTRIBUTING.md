# Contribuyendo a RoboCombat Controller

¡Gracias por tu interés en contribuir! 🎉

## 🚀 Cómo Contribuir

### Reportar Bugs
1. Verifica que el bug no esté ya reportado en [Issues](https://github.com/tu-usuario/robocombat/issues)
2. Crea un nuevo Issue con:
   - Descripción clara del problema
   - Pasos para reproducirlo
   - Sistema operativo y versión
   - Logs con `-debug` si es posible
   - Archivo `config.yaml` usado (sin datos sensibles)

### Proponer Nuevas Características
1. Abre un Issue describiendo la característica
2. Explica el caso de uso y beneficio para la comunidad
3. Espera feedback antes de empezar a codificar

### Pull Requests

#### Antes de Empezar
- Discute cambios grandes en un Issue primero
- Asegúrate de que tu código sigue el estilo del proyecto
- Prueba tu código en Windows, Linux o macOS según corresponda

#### Proceso
1. Fork el repositorio
2. Crea una rama desde `main`:
   ```bash
   git checkout -b feature/mi-nueva-caracteristica
   ```
3. Haz tus cambios:
   - Código limpio y comentado
   - Sigue las convenciones de Go (`gofmt`, `golint`)
   - Actualiza documentación si es necesario
4. Commit con mensajes descriptivos:
   ```bash
   git commit -m "feat: añade soporte para controladores PS4"
   ```
5. Push a tu fork:
   ```bash
   git push origin feature/mi-nueva-caracteristica
   ```
6. Abre un Pull Request explicando:
   - Qué problema resuelve
   - Cómo lo probaste
   - Screenshots/videos si aplica

## 📝 Estilo de Código

### Go
- Usa `gofmt` para formatear
- Sigue [Effective Go](https://go.dev/doc/effective_go)
- Comenta funciones públicas
- Nombres descriptivos (evita `x`, `y`, `temp`)

### YAML
- Indentación: 2 espacios
- Comentarios claros para usuarios finales

### Commits
Usa [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` - Nueva característica
- `fix:` - Corrección de bug
- `docs:` - Cambios en documentación
- `refactor:` - Refactorización sin cambios funcionales
- `test:` - Añadir o modificar tests
- `chore:` - Mantenimiento (deps, build, etc.)

## 🧪 Testing

Aunque el proyecto no tiene tests automatizados aún, prueba manualmente:

### Checklist de Testing
- [ ] Compila sin errores en tu plataforma
- [ ] Funciona con `config.yaml` por defecto
- [ ] Funciona con configuración personalizada
- [ ] Probado con robot real o simulado
- [ ] Modo debug muestra información útil
- [ ] Documentación actualizada

### Pruebas Específicas por Modo
**Xbox Controller:**
- [ ] Detecta mando correctamente
- [ ] Todos los controles responden
- [ ] Deadzone funciona adecuadamente

**Keyboard:**
- [ ] Todas las teclas configuradas funcionan
- [ ] Detención al soltar teclas funciona
- [ ] Teclas especiales (space, arrows) funcionan

## 🏗️ Arquitectura

### Archivos Principales
- `main.go` - Punto de entrada, manejo de config, loops principales
- `xinput_windows.go` - Interfaz XInput para Windows
- `xinput_linux.go` - Stub multiplataforma
- `config.yaml` - Configuración de usuario

### Añadir Nueva Funcionalidad

#### Nuevo Parámetro de Configuración
1. Añade campo a struct `Config` en `main.go`
2. Añade valor por defecto en `loadConfig()`
3. Documenta en `config.yaml.example`
4. Actualiza `README.md`

#### Nuevo Modo de Entrada
1. Crea función `runXXXMode()` similar a `runXboxMode()`
2. Añade opción en `inputMode` del config
3. Añade case en `main()` para tu modo
4. Documenta controles en README

#### Nueva API Endpoint
1. Añade función `sendXXX()` en `main.go`
2. Documenta en sección API del README
3. Prueba con robot real

## 🐛 Debugging

### Habilitar Logs Detallados
```bash
./robocombat -host 192.168.4.1 -debug
```

### Añadir Logs en Código
```go
if debug {
    log.Printf("DEBUG - Mi mensaje: %v", variable)
}
```

### Probar sin Robot
Puedes usar un servidor HTTP mock para probar:
```bash
# Terminal 1: Mock API
python3 -m http.server 8080

# Terminal 2: App apuntando al mock
./robocombat -host localhost:8080
```

## 📚 Recursos

- [Go Documentation](https://go.dev/doc/)
- [YAML Spec](https://yaml.org/spec/)
- [XInput API](https://docs.microsoft.com/en-us/windows/win32/xinput/xinput-game-controller-apis-portal)
- [GitHub Keyboard Library](https://github.com/eiannone/keyboard)

## 💬 Comunicación

- **Issues**: Para bugs y features
- **Pull Requests**: Para código
- **Discusiones**: Para preguntas generales (si está habilitado)

## ⚖️ Código de Conducta

- Sé respetuoso con todos los colaboradores
- Acepta críticas constructivas
- Enfócate en lo mejor para el proyecto y la comunidad
- Ayuda a nuevos contribuidores

## 📋 Checklist Pre-Submit

Antes de abrir tu PR, verifica:

- [ ] Mi código compila sin warnings
- [ ] He probado los cambios manualmente
- [ ] He actualizado documentación relevante
- [ ] Los commits siguen convenciones
- [ ] El mensaje del PR es claro y descriptivo
- [ ] He respondido a comentarios de código reviewers

---

¡Gracias por hacer de RoboCombat Controller un mejor proyecto! 🤖✨
