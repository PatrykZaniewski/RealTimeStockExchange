import uuid
from dataclasses import dataclass
from enum import Enum

from dataclasses_json import dataclass_json, LetterCase


class OrderType(Enum):
    BUY = "BUY"
    SELL = "SELL"


class OrderSubtype(Enum):
    MARKET_ORDER = "MARKET_ORDER"
    LIMIT_ORDER = "LIMIT_ORDER"


@dataclass_json(letter_case=LetterCase.CAMEL)
@dataclass
class Order:
    asset_name: str
    quantity: int
    order_type: str
    order_subtype: str
    order_price: float
    id: str

    @staticmethod
    def create(asset_name: str, quantity: int, price: float, order_type: str, order_subtype: str):
        return Order(asset_name, quantity, order_type, order_subtype, price, str(uuid.uuid4()))

