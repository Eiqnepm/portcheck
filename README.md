# portcheck
<details>
  <summary>Docker Compose example</summary>

  ```yaml
  version: "3"

  services:
    gluetun:
      cap_add:
        - "NET_ADMIN"
      container_name: "gluetun"
      devices:
        - "/dev/net/tun:/dev/net/tun"
      environment:
        VPN_SERVICE_PROVIDER: "mullvad"
        VPN_TYPE: "wireguard"
        WIREGUARD_PRIVATE_KEY: "👀"
        WIREGUARD_ADDRESSES: "👀"
        SERVER_CITIES: "Amsterdam"
        OWNED_ONLY: "yes"
        FIREWALL_VPN_INPUT_PORTS: "6881"
      image: "qmcgaw/gluetun:latest"
      ports:
        # Gluetun
        - "8000:8000"
        # qBittorrent
        - "8080:8080"
      restart: "always"
      volumes:
        - "./gluetun:/gluetun"

    qbittorrent:
      container_name: "qbittorrent"
      depends_on:
        - "gluetun"
      environment:
        PUID: "1000"
        PGID: "1000"
        TZ: "Etc/UTC"
        WEBUI_PORT: "8080"
      image: "lscr.io/linuxserver/qbittorrent:latest"
      network_mode: "service:gluetun"
      restart: "always"
      volumes:
        - "./qbittorrent:/config"
        - "./torrents:/downloads"

    portcheck:
      container_name: "portcheck"
      depends_on:
        - "gluetun"
        - "qbittorrent"
      environment:
        QBITTORRENT_PORT: "6881"
        QBITTORRENT_WEBUI_PORT: "8080"
        QBITTORRENT_WEBUI_SCHEME: "http"
        QBITTORRENT_USERNAME: "admin"
        QBITTORRENT_PASSWORD: "adminadmin"
        TIMEOUT: "300"
      image: "eiqnepm/portcheck:latest"
      network_mode: "service:gluetun"
      restart: "always"
  ```
</details>

## Environment variables
|Variable|Default|Description|
|-|-|-|
|`QBITTORRENT_PORT`|`6881`|qBittorrent incoming connection port|
|`QBITTORRENT_WEBUI_PORT`||Port of the qBittorrent WebUI|
|`QBITTORRENT_WEBUI_SCHEME`|`http`|Scheme of the qBittorrent WebUI|
|`QBITTORRENT_USERNAME`|`admin`|qBittorrent WebUI username|
|`QBITTORRENT_PASSWORD`|`adminadmin`|qBittorrent WebUI password|
|`TIMEOUT`|`300`|Time in seconds between each port check|
