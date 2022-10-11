from dataclasses import dataclass


@dataclass
class OrderStatus:
    id: str
    status: str
