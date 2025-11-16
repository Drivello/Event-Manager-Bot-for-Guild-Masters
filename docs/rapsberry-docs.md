# 🖥️ Guía de despliegue en Raspberry Pi

Esta guía explica cómo desplegar y mantener el **Discord Event Bot** en una **Raspberry Pi** (por ejemplo, una Raspberry Pi Zero 2 W) usando un binario compilado en tu PC y un servicio `systemd`.

> Nota: Los ejemplos usan rutas genéricas. Adáptalas a tu usuario y estructura (por ejemplo `/home/pi/...` o `/home/Drivello/...`).

---

## 📋 Requisitos

- Raspberry Pi con Linux (Raspberry Pi OS u otra distro basada en Debian).
- Arquitectura ARM64 recomendada (Pi Zero 2 W, 3, 4, etc.).
- Acceso SSH a la Raspberry.
- En tu PC de desarrollo:
  - Go 1.21+ instalado.
  - Este repositorio clonado.

---

## 🧱 Estructura recomendada en la Raspberry

En la Raspberry se recomienda una estructura como:

```bash
/home/<usuario>/event-manager-bot/
  ├── discord-event-bot           # Binario compilado para ARM
  ├── discord-bot.service         # Unit file de systemd
  ├── .env                        # Configuración del bot
  ├── internal/web/templates/     # Templates HTML del panel web
  └── data/
      ├── events/                 # Eventos guardados
      └── templates/              # Templates de eventos
```

Sustituye `<usuario>` por tu usuario real en la Raspberry (`pi`, `Drivello`, etc.).

---

## 1️⃣ Compilar el binario para Raspberry (en tu PC)

Desde tu PC, en la raíz del proyecto:

```bash
./build-pi.sh
```

Este script genera, entre otros, el binario `discord-event-bot-arm64` para ARM64.

Alternativa manual:

```bash
GOOS=linux GOARCH=arm64 go build -o discord-event-bot cmd/main.go
```

El resultado debe ser un binario Linux ARM64 ejecutable en la Raspberry.

---

## 2️⃣ Copiar archivos a la Raspberry

En tu PC, desde la raíz del proyecto:

```bash
# Variables de ejemplo
USER_PI=<usuario_en_pi>          # ej: pi o Drivello
HOST_PI=<ip_raspberry>          # ej: 192.168.0.82
REMOTE_DIR=/home/$USER_PI/event-manager-bot

# Crear estructura remota
ssh $USER_PI@$HOST_PI \ 
  "mkdir -p $REMOTE_DIR/internal/web/templates $REMOTE_DIR/data/events"

# Copiar binario (ajusta nombre si usaste otro)
scp discord-event-bot-arm64 $USER_PI@$HOST_PI:$REMOTE_DIR/discord-event-bot

# Copiar templates HTML del panel
scp -r internal/web/templates/* \ 
    $USER_PI@$HOST_PI:$REMOTE_DIR/internal/web/templates/

# Copiar unit file de systemd
scp discord-bot.service $USER_PI@$HOST_PI:$REMOTE_DIR/

# Copiar .env (o .env.example)
scp .env $USER_PI@$HOST_PI:$REMOTE_DIR/.env
```

En la Raspberry tendrás ahora los archivos necesarios en `$REMOTE_DIR`.

> Si usas el script `deploy-pi.sh`, está pensado para el usuario `pi` y el directorio `/home/pi/discord-event-bot`. Puedes adaptarlo cambiando las variables `PI_HOST` y `PI_DIR` para que apunten a tu usuario y ruta.

---

## 3️⃣ Crear y ajustar el servicio systemd

En la Raspberry, accede por SSH y sitúate en el directorio de despliegue:

```bash
ssh <usuario>@<ip_raspberry>
cd /home/<usuario>/event-manager-bot
chmod +x discord-event-bot
```

Edita el archivo `discord-bot.service` para que apunte a tu usuario y rutas. Un ejemplo típico:

```ini
[Unit]
Description=Discord Event Bot para MMO Guild
After=network.target

[Service]
Type=simple
User=<usuario>
WorkingDirectory=/home/<usuario>/event-manager-bot
ExecStart=/home/<usuario>/event-manager-bot/discord-event-bot
Restart=always
RestartSec=10

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=discord-bot

# Seguridad básica (sin ocultar /home para evitar problemas de rutas)
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

Sustituye `<usuario>` y las rutas según tu caso. Guarda el archivo.

Luego instala el servicio en `systemd`:

```bash
sudo cp discord-bot.service /etc/systemd/system/discord-bot.service
sudo systemctl daemon-reload
sudo systemctl enable discord-bot
```

---

## 4️⃣ Iniciar, detener y reiniciar el bot

### Iniciar el bot como servicio

```bash
sudo systemctl start discord-bot
```

### Ver estado del servicio

```bash
sudo systemctl status discord-bot
```

Debería verse `active (running)` si todo está correcto.

### Ver logs en tiempo real

```bash
sudo journalctl -u discord-bot -f
```

### Reiniciar el bot (por ejemplo tras agregar templates o cambiar `.env`)

```bash
sudo systemctl restart discord-bot
```

> Si tenías el bot corriendo manualmente con `./discord-event-bot`, detén ese proceso (Ctrl+C) antes de usar el servicio `systemd` para evitar tener dos instancias.

### Detener el bot

```bash
sudo systemctl stop discord-bot
```

---

## 5️⃣ Puertos y acceso al panel web

El puerto del panel web lo define la variable `PORT` en tu `.env`.

Ejemplo:

```env
PORT=8080
```

Por defecto verás algo como:

```text
Servidor web iniciado en http://localhost:8080
Servidor web disponible en: http://localhost:8080
```

Para acceder desde otro dispositivo de la red local, usa:

```text
http://IP_DE_TU_RASPBERRY:8080
```

Si cambiaste el puerto, ajusta la URL en consecuencia.

---

## 6️⃣ Logs y monitoreo

Ver logs en tiempo real:

```bash
sudo journalctl -u discord-bot -f
```

Ver logs de la última hora:

```bash
sudo journalctl -u discord-bot --since "1 hour ago"
```

Ver solo los últimos N mensajes:

```bash
sudo journalctl -u discord-bot -n 50
```

---

## 7️⃣ Actualizar el bot en la Raspberry

1. En tu PC, recompila el binario (por ejemplo con `./build-pi.sh`).
2. Copia el nuevo binario a la Raspberry sobre el existente:

   ```bash
   scp discord-event-bot-arm64 <usuario>@<ip_raspberry>:/home/<usuario>/event-manager-bot/discord-event-bot
   ```

3. En la Raspberry, reinicia el servicio:

   ```bash
   sudo systemctl restart discord-bot
   ```

---

## 🧩 Troubleshooting (Problemas comunes)

### 1. Error `status=217/USER` en systemd

**Síntomas:**

```text
systemd[1]: discord-bot.service: Failed at step USER spawning ...: No such process
code=exited, status=217/USER
```

**Causas posibles:**

- El usuario configurado en `User=` no existe.
- Hay un error tipográfico en el nombre de usuario.

**Solución:**

1. Verifica el usuario actual:

   ```bash
   whoami
   ```

2. Edita `/etc/systemd/system/discord-bot.service` y ajusta:

   ```ini
   User=<usuario_correcto>
   ```

3. Recarga y reinicia:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart discord-bot
   ```

---

### 2. Error `status=203/EXEC` o `No such file or directory`

**Síntomas:**

```text
Unable to locate executable '/home/<usuario>/event-manager-bot/discord-event-bot': No such file or directory
code=exited, status=203/EXEC
```

**Causas posibles:**

- La ruta en `ExecStart` no coincide con la ubicación real del binario.
- El binario tiene otro nombre (por ejemplo `discord-event-bot-arm64`).
- El archivo no tiene permisos de ejecución.

**Solución:**

1. Verifica los archivos:

   ```bash
   cd /home/<usuario>/event-manager-bot
   ls -l
   ```

2. Asegúrate de que exista `discord-event-bot` y que sea ejecutable:

   ```bash
   chmod +x discord-event-bot
   ```

3. Si el archivo tiene otro nombre, o ajustas `ExecStart=` en el servicio, o lo renombras:

   ```bash
   mv discord-event-bot-arm64 discord-event-bot
   ```

4. Recarga systemd:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart discord-bot
   ```

---

### 3. Error `exec format error`

**Síntomas:**

```text
bash: ./discord-event-bot: cannot execute binary file: Exec format error
```

**Causa:**

- El binario fue compilado para otra arquitectura (por ejemplo `amd64` en vez de `arm64`).

**Solución:**

1. En tu PC, recompila para ARM64:

   ```bash
   GOOS=linux GOARCH=arm64 go build -o discord-event-bot cmd/main.go
   ```

   o usa `./build-pi.sh`.

2. Vuelve a copiar el binario a la Raspberry y reinicia el servicio.

---

### 4. `panic: html/template: pattern matches no files: "internal/web/templates/*"`

**Síntomas:**

- El bot arranca y se cae inmediatamente con un `panic` como:

  ```text
  panic: html/template: pattern matches no files: "internal/web/templates/*"
  ```

**Causa:**

- Los templates HTML del panel web no se copiaron a la Raspberry.

**Solución:**

1. Asegúrate de tener la ruta en la Raspberry:

   ```bash
   mkdir -p /home/<usuario>/event-manager-bot/internal/web/templates
   ```

2. Desde tu PC, copia los HTML:

   ```bash
   scp -r internal/web/templates/* \ 
       <usuario>@<ip_raspberry>:/home/<usuario>/event-manager-bot/internal/web/templates/
   ```

3. Reinicia el bot:

   ```bash
   sudo systemctl restart discord-bot
   ```

---

### 5. El bot funciona manualmente pero falla bajo systemd

**Síntomas:**

- `./discord-event-bot` funciona bien.
- `systemctl start discord-bot` falla con errores de ruta.

**Causas posibles:**

- Opciones de seguridad de systemd (`ProtectHome`, `ProtectSystem`) impiden acceder a `/home`.

**Solución (simple):**

1. En `/etc/systemd/system/discord-bot.service`, asegúrate de **no** tener:

   ```ini
   ProtectSystem=strict
   ProtectHome=yes
   ```

   (o comenta esas líneas si están presentes).

2. Recarga y reinicia:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart discord-bot
   ```

**Alternativa avanzada:**

- Mover el bot a un directorio de sistema (por ejemplo `/opt/event-manager-bot`) y configurar adecuadamente las opciones de sandboxing. Esto sale del alcance de esta guía básica.

---

### 6. El panel web no carga desde otro dispositivo

**Checklist:**

1. Confirma que el servicio está corriendo:

   ```bash
   sudo systemctl status discord-bot
   ```

2. Verifica el puerto configurado en `.env` (`PORT=`) y que el bot loguea `Servidor web iniciado en http://localhost:<puerto>`.

3. Comprueba conectividad desde tu PC:

   ```bash
   curl http://<ip_raspberry>:<puerto>
   ```

4. Si tienes firewall (ufw, iptables), asegúrate de permitir el puerto.

---

### 7. Error `Permission denied` al arrancar el servicio

**Síntomas:**

```text
discord-bot.service: Unable to locate executable '/home/<usuario>/event-manager-bot/discord-event-bot': Permission denied
discord-bot.service: Failed at step EXEC spawning ...: Permission denied
code=exited, status=203/EXEC
```

**Causas posibles:**

- El binario existe pero **sin bit de ejecución** (`-rw-r--r--` en lugar de `-rwxr-xr-x`).
- El usuario configurado en `User=` no tiene permisos para entrar en el directorio (por ejemplo `/home/<usuario>` con permisos demasiado restrictivos).

**Solución:**

1. Comprueba permisos del binario:

   ```bash
   ls -l /home/<usuario>/event-manager-bot/discord-event-bot
   chmod +x /home/<usuario>/event-manager-bot/discord-event-bot
   ```

2. Comprueba permisos de directorios:

   ```bash
   ls -ld /home/<usuario> /home/<usuario>/event-manager-bot
   ```

   Si el servicio corre como otro usuario distinto, asegúrate de que pueda acceder a esas rutas o mueve el binario a un directorio como `/opt/discord-event-bot` y ajusta `ExecStart`/`WorkingDirectory`.

3. Recarga y reinicia systemd:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart discord-bot
   ```

---

### 8. El bot arranca pero la web no muestra los cambios nuevos

**Síntomas:**

- Has recompilado y desplegado un binario nuevo.
- El servicio arranca correctamente y los logs son normales.
- Pero en el panel web **no aparecen opciones nuevas** (por ejemplo, nuevos campos en el formulario de creación de eventos).

**Causa:**

- Los cambios de interfaz suelen estar en los **templates HTML**, que **no se compilan** dentro del binario. Si solo copias el binario a la Raspberry y no actualizas `internal/web/templates`, verás la UI vieja.

**Solución:**

1. Desde tu PC, copia los templates actualizados a la Raspberry:

   ```bash
   scp -r internal/web/templates/* \
       <usuario>@<ip_raspberry>:/home/<usuario>/event-manager-bot/internal/web/templates/
   ```

2. Reinicia el servicio:

   ```bash
   sudo systemctl restart discord-bot
   ```

---

Con esta guía deberías poder desplegar, actualizar y depurar el bot en una Raspberry Pi de forma fiable. Si encuentras un caso nuevo, documenta el mensaje de error y los pasos realizados para poder extender esta sección de troubleshooting.
