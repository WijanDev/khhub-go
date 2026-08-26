# Despliegue en Hetzner + Dokploy

khhub no instala Dokploy por ti. Cuando el VPS ya tenga Dokploy y un dominio apuntando a la IP:

1. Crea un proyecto Compose en Dokploy y apunta al repositorio (o pega `docker-compose.yml`).
2. Define estas variables de entorno (nunca las subas al git):

   - `POSTGRES_USER` (p. ej. `khhub`)
   - `POSTGRES_PASSWORD` — larga y aleatoria
   - `POSTGRES_DB` (`khhub`)
   - `SESSION_SECRET` — 32+ caracteres aleatorios
   - `ADMIN_EMAIL`
   - `ADMIN_PASSWORD` — mínimo 10 caracteres en producción
   - `APP_ENV=production`
   - `COOKIE_SECURE=true`
   - `CORS_ORIGINS` vacío (el SPA y la API van por el mismo origen vía nginx)

3. En Dokploy, asigna el dominio al servicio `web` (puerto 80) y activa HTTPS (Let's Encrypt).
4. Activa **copias de seguridad de Postgres** hacia un destino externo el mismo día. Esta base guarda datos personales de la congregación.
5. En el firewall del VPS deja solo 22/80/443. No expongas 5432 ni 8080 a Internet.
6. Tras el primer login, cambia la contraseña de administrador en **Congregación**.

No hay envío automático a jw.org. Los totales del panel se copian a mano.
