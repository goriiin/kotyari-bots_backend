
import random
import re
from typing import Optional, List

class ProxyPool:
    """
    Manages a pool of proxies loaded from a file.
    File format (one per line):
      - 'ip port username password'  (whitespace)
      - 'ip:port:username:password'  (colon)
    """
    def __init__(self, filepath: str):
        self.filepath = filepath
        self.proxies: List[str] = []
        self.round_robin_index = 0
        self.load_proxies()

    def _split_line(self, line: str) -> Optional[str]:
        # Support both whitespace and colon separators; collapse multiple spaces
        parts = re.split(r'[\s:]+', line.strip())
        parts = [p for p in parts if p]
        if len(parts) != 4:
            print(f"Invalid proxy format (expected 4 parts): {line}")
            return None
        host, port, user, pwd = parts
        return f"{host} {port} {user} {pwd}"

    def load_proxies(self) -> None:
        valid: List[str] = []
        try:
            with open(self.filepath, "r") as f:
                for raw in f:
                    if not raw or not raw.strip():
                        continue
                    norm = self._split_line(raw)
                    if norm:
                        valid.append(norm)
            self.proxies = valid
            print(f"Loaded {len(self.proxies)} proxies from {self.filepath}")
        except FileNotFoundError:
            print(f"Proxy file not found: {self.filepath}")
            self.proxies = []
        except Exception as e:
            print(f"Error loading proxies: {e}")
            self.proxies = []

    def get_random_proxy(self) -> Optional[str]:
        if not self.proxies:
            print("Proxy pool is empty, returning None")
            return None
        proxy = random.choice(self.proxies)
        print(f"Selected random proxy: {proxy.split()[0]}")
        return proxy

    def get_next_proxy(self) -> Optional[str]:
        if not self.proxies:
            print("Proxy pool is empty, returning None")
            return None
        proxy = self.proxies[self.round_robin_index]
        self.round_robin_index = (self.round_robin_index + 1) % len(self.proxies)
        print(f"Selected round-robin proxy: {proxy.split()[0]}")
        return proxy
