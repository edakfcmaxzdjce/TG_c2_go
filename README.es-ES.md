# TG C2 Go

Un cliente C2 (Comando y Control) ligero y multiplataforma implementado en Go, controlado remotamente a través de la API de Telegram Bot.

> **中文文档**: [README_zh.md](https://github.com/edakfcmaxzdjce/TG_c2_go/blob/main/README_zh.md) | **English**: README.md

## 📋 Descripción General del Proyecto

TG C2 Go es un cliente C2 ligero y multiplataforma que permite el control remoto y la ejecución de comandos a través de la API de Telegram Bot.

## 🚀 Características Principales

- **Integración con Telegram Bot**: Control remoto vía API de Telegram Bot.
- **Gestión de Archivos**: Carga, descarga y gestión de archivos.
- **Ejecución de Comandos**: Ejecución de comandos del sistema de forma remota.
- **Capacidades de Inyección**: Llamadas a DLL e inyección de Shellcode (solo Windows).
- **Monitoreo de Pantalla**: Funcionalidad de capturas de pantalla en tiempo real.
- **Recopilación de Información**: Información del sistema, red, programas instalados, etc.
- **Gestión de Temas (Topics)**: Creación y reutilización inteligente de Temas de Foro.
- **Soporte Multiplataforma**: Windows, Linux, macOS (algunas funciones tienen limitaciones).

## 📦 Estructura del Proyecto

```
TG_c2_go/
├── config/           # Gestión de configuración
│   └── api.go       # Configuración de API y gestión de URLs
├── telegram/         # Cliente de Telegram
│   └── client.go    # Wrapper de la API del Bot
├── core/            # Funcionalidad principal
│   ├── topic.go     # Gestión de temas
│   ├── file_manager.go  # Gestión de archivos
│   └── command_loop.go  # Bucle de comandos
├── commands/        # Procesamiento de comandos
│   └── processor.go # Emparejamiento y procesamiento de comandos
├── functions/       # Módulos de funciones
│   ├── info_collector.go   # Recopilación de información
│   ├── screen_capture.go   # Captura de pantalla
│   ├── dll_runner.go      # Llamadas a DLL
│   └── injector.go        # Inyección de código
├── go.mod          # Gestión de dependencias de Go  
├── main.go         # Punto de entrada del programa
└── README.md       # Documentación del proyecto
```

## 🛠️ Dependencias

- **github.com/shirou/gopsutil/v3** - Recopilación de información del sistema.
- **github.com/kbinani/screenshot** - Funcionalidad de captura de pantalla.
- **golang.org/x/sys** - Soporte para llamadas al sistema.

## 🚀 Inicio Rápido

### Paso 1: Crear un Telegram Bot

1. **Abre Telegram** y busca a [@BotFather](https://t.me/BotFather).

2. **Crea un nuevo Bot**:
   - Envía el comando `/newbot`.
   - Sigue las instrucciones para ingresar el nombre del Bot (ej., `Mi Bot C2`).
   - Ingresa el nombre de usuario del Bot (debe terminar en `bot`, ej., `mi_c2_bot`).

3. **Obtén el Token del Bot**:
   - BotFather te devolverá un Token, con un formato como: `1234567890:ABCdefGhiJklMnoPqrsTuvWxyz`.
   - **Importante**: Guarda este Token de forma segura y no lo compartas.

4. **Obtén el Chat ID**:
   - Crea un grupo o canal de Telegram (se recomienda un grupo privado).
   - Añade tu Bot al grupo.
   - Envía un mensaje al grupo.
   - Visita `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`.
   - Busca el objeto `chat` en el JSON devuelto; el campo `id` es el Chat ID (ej., `-1002927089835`).

### Paso 2: Configurar el Cliente

Edita el archivo `config/api.go` y completa estas tres configuraciones:

```go
BOT_TOKEN = "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"  // Obtenido de BotFather
CHAT_ID  = "-1002927012345"                          // Chat ID de tu grupo
BOT_NAME = "my_c2_bot"                               // Nombre de usuario del Bot (sin @)
```

**¡Eso es todo!** Todas las URLs de la API se generarán automáticamente a partir del `BOT_TOKEN`, no es necesaria configuración manual.

### Paso 3: Compilar el Programa

```bash
# Usando el script de compilación (recomendado)
./build.sh  # Linux/macOS
# o
build.bat  # Windows

# El script te pedirá elegir el modo de compilación:
# 1) Versión Release (ventana de Windows oculta)
# 2) Versión Debug (consola visible)
```

### Paso 4: Ejecutar el Programa

```bash
# Linux/macOS
./build/TG_c2_go_linux_amd64
./build/TG_c2_go_darwin_arm64

# Windows
build\TG_c2_go_windows_amd64.exe
```

## 🔧 Configuración

### Método 1: Usando Bot Token (Recomendado)

Solo configura tres elementos:

```go
BOT_TOKEN = "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"
CHAT_ID  = "-1002927089835"
BOT_NAME = "my_c2_bot"
```

Todas las URLs de la API se generan automáticamente.

### Método 2: Configuración Manual de URLs Codificadas en Base64 (Retrocompatible)

Si prefieres no usar el método del Token, puedes llenar manualmente las URLs codificadas en Base64:

```go
UPDATE_URL           = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDLi4uL2dldFVwZGF0ZXM="
SENDMSG_URL          = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDLi4uL3NlbmRNZXNzYWdl"
// ... otras URLs
```

Puedes usar la herramienta `tools/config_generator.go` para generar estas URLs en Base64.

## 🚀 Compilación y Ejecución

### 1. Instalar Dependencias

```bash
go mod tidy
```

### 2. Configurar API

Edita `config/api.go` y completa tu Bot Token, Chat ID y Bot Name (ver sección "Inicio Rápido" arriba).

### 3. Compilar el Programa

#### Usando Script de Compilación (Recomendado)

```bash
# Linux/macOS
./build.sh

# Windows
build.bat
```

El script de compilación te permitirá elegir el modo:
- **Versión Release**: En Windows oculta la ventana de la consola.
- **Versión Debug**: Ventana de consola visible para ver logs y salida.

El script utiliza automáticamente los siguientes flags de optimización para reducir el tamaño del binario:
- `-trimpath`: Elimina la información de la ruta del archivo.
- `-ldflags="-s -w"`: Elimina la tabla de símbolos y la información de depuración.
- La versión release de Windows añade `-H windowsgui` para ocultar la consola.

#### Compilación Manual

```bash
# Versión de desarrollo (no recomendada, tamaño de archivo mayor)
go build -o TG_c2_go

# Versión optimizada (recomendada, reduce el tamaño ~10-20%)
go build -trimpath -ldflags="-s -w" -o TG_c2_go

# Compilación cruzada para Windows
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o TG_c2_go.exe
```

#### Reducción Adicional de Tamaño (Opcional)

Si el binario sigue siendo grande (>6MB), considera:

1. **Uso de compresión UPX** (puede ser detectado por antivirus):
   ```bash
   # Instala UPX: https://upx.github.io/
   upx --best --lzma TG_c2_go
   # Puede reducir el tamaño entre 50-70%, pero el tiempo de inicio puede aumentar.
   ```

2. **Eliminar dependencias no utilizadas**:
   ```bash
   go mod tidy
   ```

3. **Verificar tamaño de dependencias**:
   ```bash
   go list -json -deps | jq -r '.Deps[] | select(.Standard == false) | .Path'
   ```

**Nota**: El tamaño del binario optimizado suele estar entre 5-7MB, lo cual es normal para programas en Go.

### 4. Ejecutar el Programa

```bash
./TG_c2_go        # Linux/macOS
TG_c2_go.exe      # Windows
```

## 📖 Guía de Uso

### Primera Ejecución

1. **Asegurar que el Bot esté en el grupo**: Añade el Bot creado al grupo de Telegram especificado.

2. **Iniciar el cliente**: Ejecuta el ejecutable compilado.

3. **Inicialización automática**: El programa hará lo siguiente automáticamente:
   - Obtener la dirección IP pública.
   - Conectarse a la API de Telegram Bot.
   - Crear o reutilizar un Tema de Foro (nombrado con la dirección IP).
   - Comenzar a escuchar comandos.

4. **Verificar conexión**: El programa enviará un mensaje de conexión en el Tema (si utiliza un Tema existente).

### Uso de Temas de Foro (Recomendado)

- El programa crea automáticamente Temas de Foro independientes para cada IP de cliente.
- Si la misma IP se reconecta, reutilizará el Tema existente.
- Puedes gestionar múltiples clientes en un solo grupo.
- Formato del nombre del Tema: `IP: xxx.xxx.xxx.xxx`.

### Ver Logs

- **Versión Debug**: Se puede ver toda la salida en la consola.
- **Versión Release** (Windows): Sin ventana de consola, verifica los resultados de ejecución a través de los mensajes de Telegram.

## 📱 Comandos del Telegram Bot

### Comandos Básicos

Envía los siguientes comandos en el Tema del grupo de Telegram (la parte `@your_bot` se puede omitir si solo hay un Bot en el grupo):

| Comando | Descripción | Ejemplo |
|---------|-------------|---------|
| `/screen_shot@your_bot` | Obtener captura de pantalla del host objetivo | `/screen_shot@my_c2_bot` |
| `/info_collect@your_bot` | Recopilar información del sistema (CPU, memoria, disco, red, etc.) | `/info_collect@my_c2_bot` |
| `/upload@your_bot <file_path>` | Cargar archivo desde la ruta especificada | `/upload@my_c2_bot C:\Windows\System32\hosts` |
| `/set_sleep_time@your_bot <seconds>` | Establecer intervalo de sondeo de comandos (segundos) | `/set_sleep_time@my_c2_bot 10` |
| `/setting_info@your_bot` | Mostrar información de configuración actual | `/setting_info@my_c2_bot` |
| `/disconnect@your_bot` | Desconectar y salir del programa | `/disconnect@my_c2_bot` |

### Ejecutar Comandos del Sistema

Simplemente envía cualquier mensaje de texto en el Tema (que no empiece por `/`), y el programa lo ejecutará automáticamente como un comando del sistema:

- **Windows**: Ejecutado usando PowerShell.
- **macOS/Linux**: Ejecutado usando `sh -c`.

Ejemplos:
```
ls -la
whoami
netstat -an
```

### Funciones Avanzadas (Solo Windows)

| Comando | Descripción | Ejemplo |
|---------|-------------|---------|
| `/run_dll@your_bot <dll_name> <func_name>` | Llamar a una función en una DLL | `/run_dll@my_c2_bot user32.dll MessageBoxA` |

### Manejo de Archivos

- **Enviar archivos regulares**: Envía archivos en el Tema; el programa los descargará automáticamente en el directorio `output_dir`.
- **Inyección de Shellcode**: Al enviar archivos, escribe `inject` en el Pie de foto (Caption/título), el programa ejecutará la inyección de Shellcode (solo Windows).

### Consejos de Uso

1. **Resultados de ejecución**: Los resultados se devuelven mediante mensajes de Telegram; las salidas largas se dividirán automáticamente en trozos.
2. **Carga de archivos**: Los archivos cargados aparecerán como adjuntos de mensaje en el Tema.
3. **Gestión multi-cliente**: Cada cliente tiene un Tema independiente, permitiendo gestionar múltiples objetivos simultáneamente.
4. **Reutilización de Temas**: Los clientes con la misma dirección IP reutilizarán automáticamente el mismo Tema.

## 🔒 Características de Seguridad

1. **Ofuscación de URL (Opcional)**: Soporte para almacenamiento de URLs codificadas en Base64 (retrocompatible).
2. **Configuración de Token**: Se recomienda usar el Token del Bot para generar URLs automáticamente, simplificando la configuración.
3. **Seguridad de Hilos**: Uso de mutex para proteger el estado global.
4. **Manejo de Errores**: Sistema exhaustivo de manejo de errores y mecanismo de reintentos.
5. **Seguridad de Memoria**: El recolector de basura de Go garantiza la seguridad de la memoria.
6. **Soporte TLS**: Soporte para configuración TLS personalizada.

## ⚠️ Limitaciones de Plataforma

- **Llamadas a DLL**: Solo Windows.
- **Inyección de Shellcode**: Solo Windows.
- **Recopilación de info del sistema**: Soporte multiplataforma, aunque Windows tiene funcionalidades más completas.

## ❓ Preguntas Frecuentes (FAQ)

### P: ¿Cómo obtengo el Chat ID?

1. Añade el Bot al grupo.
2. Envía un mensaje al grupo.
3. Visita: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`.
4. Busca `"chat":{"id":-1001234567890}` en el JSON devuelto; ese número es el Chat ID.

### P: El programa muestra "BOT_TOKEN not configured", ¿qué hago?

Asegúrate de haber completado la variable `BOT_TOKEN` en `config/api.go`. Si usas el método antiguo de URLs en Base64, asegúrate de que al menos una URL esté completada.

### P: ¿La versión de Windows no muestra ninguna ventana?

Es el comportamiento normal. La versión Release oculta la ventana de la consola. Si necesitas ver los logs, utiliza la compilación en modo Debug.

### P: ¿Cómo gestiono varios clientes simultáneamente?

Cada cliente creará o reutilizará Temas de Foro independientes basados en su dirección IP pública. Puedes gestionar múltiples clientes en el mismo grupo sin que interfieran entre sí.

### P: La creación del Tema falló, ¿qué hacer?

Verifica lo siguiente:
1. ¿Tiene el Bot permisos para crear Temas en el grupo?
2. ¿Está activada la función de Foro en el grupo?
3. ¿Es correcto el Token del Bot?
4. ¿La conexión de red es normal?

## 🔄 Diferencias con la Versión de Rust

### Ventajas
- **Despliegue más sencillo**: Un único archivo ejecutable, no requiere runtime adicional.
- **Mejor soporte multiplataforma**: La biblioteca estándar de Go ofrece mejor compatibilidad entre sistemas.
- **Seguridad de memoria**: El mecanismo de garbage collection evita fugas de memoria.
- **Compilación más rápida**: Go compila más rápido.

### Paridad de Funciones
- ✅ Todas las funciones core de C2.
- ✅ Integración con Telegram Bot.
- ✅ Gestión de archivos e inyección.
- ✅ Recopilación de información del sistema.
- ✅ Ejecución de comandos.
- ✅ Funcionalidad de captura de pantalla.

## 🛡️ Descargo de Responsabilidad

Este proyecto tiene fines puramente educativos y de investigación. Por favor, cumpla con las leyes y regulaciones pertinentes. Los usuarios son responsables de sus propias acciones.

## 🤝 Contribuciones

Se agradece el envío de Issues y Pull Requests para mejorar el proyecto.

## 📄 Licencia

Este proyecto se basa en la versión original de Rust y hereda los mismos términos de licencia.

---

**Autor Original**: bamuwe  
**Port Versión Go**: Assistant  
**Versión**: 1.0.0  
**Última Actualización**: 31 de octubre de 2025
