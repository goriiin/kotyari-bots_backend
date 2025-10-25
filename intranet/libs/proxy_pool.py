"""
Proxy pool manager for rotating proxies across parsers.
Supports round-robin and random selection strategies.
"""
import logging
import random
from typing import Optional, List

logger = logging.getLogger(__name__)


class ProxyPool:
    """
    Manages a pool of proxies loaded from a file.

    File format (one per line):
        ip:port:username:password
    """

    def __init__(self, filepath: str):
        """
        Initialize proxy pool from file.

        Args:
            filepath: Path to file with proxy list (ip:port:user:pass)
        """
        self.filepath = filepath
        self.proxies: List[str] = []
        self._round_robin_index = 0
        self._load_proxies()

    def _load_proxies(self) -> None:
        """Load and validate proxies from file."""
        try:
            with open(self.filepath, 'r') as f:
                lines = [line.strip() for line in f if line.strip()]

            valid = []
            for line in lines:
                parts = line.split(':')
                if len(parts) == 4:
                    valid.append(line)
                else:
                    logger.warning(f"Invalid proxy format (expected 4 parts): {line}")

            self.proxies = valid
            logger.info(f"Loaded {len(self.proxies)} proxies from {self.filepath}")

        except FileNotFoundError:
            logger.error(f"Proxy file not found: {self.filepath}")
            self.proxies = []
        except Exception as e:
            logger.error(f"Error loading proxies: {e}")
            self.proxies = []

    def get_random_proxy(self) -> Optional[str]:
        """Get random proxy from pool, or None if empty."""
        if not self.proxies:
            logger.warning("Proxy pool is empty, returning None")
            return None
        proxy = random.choice(self.proxies)
        logger.debug(f"Selected random proxy: {proxy.split(':')[0]}:****")
        return proxy

    def get_next_proxy(self) -> Optional[str]:
        """Get next proxy using round-robin, or None if empty."""
        if not self.proxies:
            logger.warning("Proxy pool is empty, returning None")
            return None
        proxy = self.proxies[self._round_robin_index]
        self._round_robin_index = (self._round_robin_index + 1) % len(self.proxies)
        logger.debug(f"Selected round-robin proxy: {proxy.split(':')[0]}:****")
        return proxy
