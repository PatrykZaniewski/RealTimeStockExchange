import dataclasses
from dataclasses import dataclass
from enum import Enum

import requests as requests
from dataclasses_json import dataclass_json, LetterCase


class OrderType(Enum):
    BUY = "BUY"
    SELL = "SELL"


@dataclass_json(letter_case=LetterCase.CAMEL)
@dataclass
class Order:
    asset_name: str
    quantity: int
    order_type: str


def process_order(asset_name: str, amount: int, order_type: OrderType):
    order = Order(asset_name, int(amount), order_type.value)
    requests.post(url="http://localhost:5012/order", data=order.to_json())
    pass
