from enum import Enum


class OrderType(Enum):
    BUY = "BUY"
    SELL = "SELL"


def process_order(asset_name: str, amount: int, order_type: OrderType):
    pass
